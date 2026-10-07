package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
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

type Server struct {
	mu sync.Mutex
	nodes map[string]*Node
	images map[string]Image
}

func newServer() *Server {
	s:=&Server{nodes:map[string]*Node{},images:map[string]Image{}}
	for i:=1;i<=8;i++ {
		id:=fmt.Sprintf("node-%02d",i)
		s.nodes[id]=&Node{ID:id,Name:id,Power:"off",Health:"unknown",
			Capabilities:[]string{"video","hid","serial","power","reset","virtual_media"}}
	}
	s.images["demo-os"]=Image{ID:"demo-os",Name:"KVMux demo installer",MediaType:"disk",ReadOnly:true}
	return s
}

func writeJSON(w http.ResponseWriter,status int,v any) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) nodesHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method!="GET" { http.Error(w,"method not allowed",405); return }
	s.mu.Lock(); defer s.mu.Unlock()
	out:=make([]*Node,0,len(s.nodes))
	for _,n:=range s.nodes { out=append(out,n) }
	writeJSON(w,200,out)
}

func (s *Server) imagesHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method!="GET" { http.Error(w,"method not allowed",405); return }
	s.mu.Lock(); defer s.mu.Unlock()
	out:=make([]Image,0,len(s.images))
	for _,im:=range s.images { out=append(out,im) }
	writeJSON(w,200,out)
}

func (s *Server) nodeHandler(w http.ResponseWriter,r *http.Request) {
	parts:=strings.Split(strings.Trim(r.URL.Path,"/"),"/")
	if len(parts)<4 { http.NotFound(w,r); return }
	id:=parts[3]
	s.mu.Lock(); defer s.mu.Unlock()
	n,ok:=s.nodes[id]
	if !ok { http.Error(w,"node not found",404); return }
	if len(parts)==4 && r.Method=="GET" { writeJSON(w,200,n); return }
	if r.Method!="POST" { http.Error(w,"method not allowed",405); return }

	action:=strings.Join(parts[4:],"/")
	switch action {
	case "power/on":
		n.Power="on"; n.Health="unknown"
	case "power/off":
		n.Power="off"; n.Health="unknown"
	case "power/cycle":
		n.Power="cycling"
		time.Sleep(100*time.Millisecond)
		n.Power="on"; n.Health="unknown"
	case "reset":
		if n.Power!="on" { http.Error(w,"cannot reset powered-off node",409); return }
		n.Health="unknown"
	case "media/attach":
		var req struct{ Image string `json:"image"` }
		if json.NewDecoder(r.Body).Decode(&req)!=nil { http.Error(w,"invalid json",400); return }
		if _,ok:=s.images[req.Image]; !ok { http.Error(w,"image not found",404); return }
		n.AttachedImage=req.Image
	case "media/detach":
		n.AttachedImage=""
	case "provision":
		var req struct{ Image string `json:"image"` }
		if json.NewDecoder(r.Body).Decode(&req)!=nil { http.Error(w,"invalid json",400); return }
		if _,ok:=s.images[req.Image]; !ok { http.Error(w,"image not found",404); return }
		n.Power="off"; n.AttachedImage=req.Image
		time.Sleep(100*time.Millisecond)
		n.Power="on"; n.Health="installing"
		time.Sleep(100*time.Millisecond)
		n.AttachedImage=""; n.Health="healthy"
	default:
		http.NotFound(w,r); return
	}
	writeJSON(w,200,n)
}

func main() {
	s:=newServer()
	mux:=http.NewServeMux()
	mux.HandleFunc("/api/v1/nodes",s.nodesHandler)
	mux.HandleFunc("/api/v1/nodes/",s.nodeHandler)
	mux.HandleFunc("/api/v1/images",s.imagesHandler)
	addr:=os.Getenv("KVMUX_LISTEN")
	if addr=="" { addr=":8080" }
	log.Printf("KVMux M0 simulator listening on %s",addr)
	log.Fatal(http.ListenAndServe(addr,mux))
}
