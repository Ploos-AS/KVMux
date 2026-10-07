package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Power string `json:"power_state"`
	Health string `json:"health"`
	Capabilities []string `json:"capabilities"`
	AttachedImage string `json:"attached_image,omitempty"`
}

type Image struct {
	ID string `json:"id"`
	Name string `json:"name"`
	MediaType string `json:"media_type"`
	ReadOnly bool `json:"read_only"`
}

type Job struct {
	ID string `json:"id"`
	Type string `json:"type"`
	NodeID string `json:"node_id"`
	Image string `json:"image,omitempty"`
	State string `json:"state"`
	Step string `json:"step"`
	Error string `json:"error,omitempty"`
}

type Server struct {
	mu sync.Mutex
	nodes map[string]*Node
	images map[string]Image
	jobs map[string]*Job
	jobSeq uint64
	power PowerBackend
	reset ResetBackend
	media MediaBackend
}

func newServer() *Server {
	b:=SimBackend{}
	s:=&Server{nodes:map[string]*Node{},images:map[string]Image{},jobs:map[string]*Job{},
		power:b,reset:b,media:b}
	for i:=1;i<=8;i++ {
		id:=fmt.Sprintf("node-%02d",i)
		s.nodes[id]=&Node{ID:id,Name:id,Power:"off",Health:"unknown",
			Capabilities:[]string{"video","hid","serial","power","reset","virtual_media"}}
	}
	s.images["demo-os"]=Image{ID:"demo-os",Name:"KVMux demo installer",MediaType:"disk",ReadOnly:true}
	return s
}

func writeJSON(w http.ResponseWriter,status int,v any) {
	w.Header().Set("Content-Type","application/json"); w.WriteHeader(status); _=json.NewEncoder(w).Encode(v)
}

func (s *Server) nodesHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method!="GET" { http.Error(w,"method not allowed",405); return }
	s.mu.Lock(); defer s.mu.Unlock()
	out:=make([]*Node,0,len(s.nodes)); for _,n:=range s.nodes { out=append(out,n) }; writeJSON(w,200,out)
}
func (s *Server) imagesHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method!="GET" { http.Error(w,"method not allowed",405); return }
	s.mu.Lock(); defer s.mu.Unlock()
	out:=make([]Image,0,len(s.images)); for _,im:=range s.images { out=append(out,im) }; writeJSON(w,200,out)
}
func (s *Server) jobsHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method!="GET" { http.Error(w,"method not allowed",405); return }
	id:=strings.TrimPrefix(r.URL.Path,"/api/v1/jobs/")
	s.mu.Lock(); defer s.mu.Unlock()
	j,ok:=s.jobs[id]; if !ok { http.Error(w,"job not found",404); return }; writeJSON(w,200,j)
}

func (s *Server) nodeHandler(w http.ResponseWriter,r *http.Request) {
	parts:=strings.Split(strings.Trim(r.URL.Path,"/"),"/"); if len(parts)<4 { http.NotFound(w,r); return }
	id:=parts[3]
	s.mu.Lock()
	n,ok:=s.nodes[id]; if !ok { s.mu.Unlock(); http.Error(w,"node not found",404); return }
	if len(parts)==4 && r.Method=="GET" { writeJSON(w,200,n); s.mu.Unlock(); return }
	if r.Method!="POST" { s.mu.Unlock(); http.Error(w,"method not allowed",405); return }
	action:=strings.Join(parts[4:],"/")
	switch action {
	case "power/on":
		err:=s.power.On(n); if err!=nil { s.mu.Unlock(); http.Error(w,err.Error(),409); return }
	case "power/off":
		err:=s.power.Off(n); if err!=nil { s.mu.Unlock(); http.Error(w,err.Error(),409); return }
	case "power/cycle":
		err:=s.power.Cycle(n,100*time.Millisecond); if err!=nil { s.mu.Unlock(); http.Error(w,err.Error(),409); return }
	case "reset":
		err:=s.reset.Reset(n); if err!=nil { s.mu.Unlock(); http.Error(w,err.Error(),409); return }
	case "media/attach":
		var req struct{ Image string `json:"image"` }; if json.NewDecoder(r.Body).Decode(&req)!=nil { s.mu.Unlock(); http.Error(w,"invalid json",400); return }
		im,ok:=s.images[req.Image]; if !ok { s.mu.Unlock(); http.Error(w,"image not found",404); return }
		if err:=s.media.Attach(n,im); err!=nil { s.mu.Unlock(); http.Error(w,err.Error(),409); return }
	case "media/detach":
		if err:=s.media.Detach(n); err!=nil { s.mu.Unlock(); http.Error(w,err.Error(),409); return }
	case "provision":
		var req struct{ Image string `json:"image"` }; if json.NewDecoder(r.Body).Decode(&req)!=nil { s.mu.Unlock(); http.Error(w,"invalid json",400); return }
		if _,ok:=s.images[req.Image]; !ok { s.mu.Unlock(); http.Error(w,"image not found",404); return }
		jid:=fmt.Sprintf("job-%06d",atomic.AddUint64(&s.jobSeq,1))
		j:=&Job{ID:jid,Type:"provision",NodeID:id,Image:req.Image,State:"queued",Step:"queued"}
		s.jobs[jid]=j; s.mu.Unlock(); go s.runProvision(jid); writeJSON(w,http.StatusAccepted,j); return
	default:
		s.mu.Unlock(); http.NotFound(w,r); return
	}
	writeJSON(w,200,n); s.mu.Unlock()
}

func (s *Server) jobStep(id,state,step string) {
	s.mu.Lock(); defer s.mu.Unlock(); if j:=s.jobs[id]; j!=nil { j.State=state; j.Step=step }
}
func (s *Server) jobFail(id string,err error) {
	s.mu.Lock(); defer s.mu.Unlock(); if j:=s.jobs[id]; j!=nil { j.State="failed"; j.Error=err.Error() }
}
func (s *Server) runProvision(id string) {
	s.mu.Lock(); j:=s.jobs[id]; n:=s.nodes[j.NodeID]; im:=s.images[j.Image]; s.mu.Unlock()
	s.jobStep(id,"running","power-off"); s.mu.Lock(); err:=s.power.Off(n); s.mu.Unlock(); if err!=nil { s.jobFail(id,err); return }
	s.jobStep(id,"running","attach-media"); s.mu.Lock(); err=s.media.Attach(n,im); s.mu.Unlock(); if err!=nil { s.jobFail(id,err); return }
	s.jobStep(id,"running","boot-installer"); s.mu.Lock(); err=s.power.On(n); n.Health="installing"; s.mu.Unlock(); if err!=nil { s.jobFail(id,err); return }
	time.Sleep(100*time.Millisecond)
	s.jobStep(id,"running","detach-media"); s.mu.Lock(); err=s.media.Detach(n); s.mu.Unlock(); if err!=nil { s.jobFail(id,err); return }
	s.jobStep(id,"running","verify"); s.mu.Lock(); n.Health="healthy"; s.mu.Unlock()
	s.jobStep(id,"completed","done")
}

func main() {
	s:=newServer(); mux:=http.NewServeMux()
	mux.HandleFunc("/api/v1/nodes",s.nodesHandler); mux.HandleFunc("/api/v1/nodes/",s.nodeHandler)
	mux.HandleFunc("/api/v1/images",s.imagesHandler); mux.HandleFunc("/api/v1/jobs/",s.jobsHandler)
	addr:=os.Getenv("KVMUX_LISTEN"); if addr=="" { addr=":8080" }
	log.Printf("KVMux M0 simulator listening on %s",addr); log.Fatal(http.ListenAndServe(addr,mux))
}
