package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEightNodes(t *testing.T) {
	s:=newServer(); if len(s.nodes)!=8 { t.Fatalf("got %d nodes, want 8",len(s.nodes)) }
}
func TestPowerCycle(t *testing.T) {
	s:=newServer(); req:=httptest.NewRequest(http.MethodPost,"/api/v1/nodes/node-01/power/cycle",nil); w:=httptest.NewRecorder()
	s.nodeHandler(w,req); if w.Code!=200 || s.nodes["node-01"].Power!="on" { t.Fatalf("status=%d power=%s",w.Code,s.nodes["node-01"].Power) }
}
func TestVirtualMedia(t *testing.T) {
	s:=newServer(); req:=httptest.NewRequest(http.MethodPost,"/api/v1/nodes/node-01/media/attach",strings.NewReader(`{"image":"demo-os"}`)); w:=httptest.NewRecorder()
	s.nodeHandler(w,req); if w.Code!=200 || s.nodes["node-01"].AttachedImage!="demo-os" { t.Fatal("attach failed") }
}
func TestProvisionIsAsyncAndCompletes(t *testing.T) {
	s:=newServer(); req:=httptest.NewRequest(http.MethodPost,"/api/v1/nodes/node-01/provision",strings.NewReader(`{"image":"demo-os"}`)); w:=httptest.NewRecorder()
	s.nodeHandler(w,req); if w.Code!=http.StatusAccepted { t.Fatalf("status=%d",w.Code) }
	var j Job; if err:=json.Unmarshal(w.Body.Bytes(),&j); err!=nil { t.Fatal(err) }
	deadline:=time.Now().Add(2*time.Second)
	for time.Now().Before(deadline) {
		s.mu.Lock(); state:=s.jobs[j.ID].State; s.mu.Unlock()
		if state=="completed" {
			s.mu.Lock(); n:=*s.nodes["node-01"]; s.mu.Unlock()
			if n.Power!="on" || n.Health!="healthy" || n.AttachedImage!="" { t.Fatalf("bad final state: %+v",n) }
			return
		}
		time.Sleep(10*time.Millisecond)
	}
	t.Fatal("provision job did not complete")
}

func TestBusyNodeRejectsSecondProvision(t *testing.T) {
	s:=newServer()
	req:=httptest.NewRequest(http.MethodPost,"/api/v1/nodes/node-01/provision",strings.NewReader(`{"image":"demo-os"}`))
	w:=httptest.NewRecorder(); s.nodeHandler(w,req)
	if w.Code!=http.StatusAccepted { t.Fatalf("first status=%d",w.Code) }
	req2:=httptest.NewRequest(http.MethodPost,"/api/v1/nodes/node-01/provision",strings.NewReader(`{"image":"demo-os"}`))
	w2:=httptest.NewRecorder(); s.nodeHandler(w2,req2)
	if w2.Code!=http.StatusLocked { t.Fatalf("second status=%d, want 423",w2.Code) }
}

func TestFaultInjectionFailsJob(t *testing.T) {
	s:=newServer(); b:=SimBackend{FailMediaAttach:true}; s.power=b; s.reset=b; s.media=b
	req:=httptest.NewRequest(http.MethodPost,"/api/v1/nodes/node-01/provision",strings.NewReader(`{"image":"demo-os"}`))
	w:=httptest.NewRecorder(); s.nodeHandler(w,req)
	var j Job; _=json.Unmarshal(w.Body.Bytes(),&j)
	deadline:=time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		s.mu.Lock(); state:=s.jobs[j.ID].State; busy:=s.busy["node-01"]; s.mu.Unlock()
		if state=="failed" {
			if busy!="" { t.Fatal("node lock not released after failed job") }
			return
		}
		time.Sleep(10*time.Millisecond)
	}
	t.Fatal("fault-injected job did not fail")
}

func TestAuditLog(t *testing.T) {
	s:=newServer()
	req:=httptest.NewRequest(http.MethodPost,"/api/v1/nodes/node-01/power/on",nil)
	w:=httptest.NewRecorder(); s.nodeHandler(w,req)
	s.mu.Lock(); defer s.mu.Unlock()
	if len(s.events)==0 || s.events[0].Action!="power/on" { t.Fatalf("missing audit event: %+v",s.events) }
}
