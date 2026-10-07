package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEightNodes(t *testing.T) {
	s:=newServer()
	if len(s.nodes)!=8 { t.Fatalf("got %d nodes, want 8",len(s.nodes)) }
}

func TestPowerCycle(t *testing.T) {
	s:=newServer()
	req:=httptest.NewRequest(http.MethodPost,"/api/v1/nodes/node-01/power/cycle",nil)
	w:=httptest.NewRecorder()
	s.nodeHandler(w,req)
	if w.Code!=200 { t.Fatalf("status %d",w.Code) }
	if s.nodes["node-01"].Power!="on" { t.Fatalf("power=%s",s.nodes["node-01"].Power) }
}

func TestVirtualMedia(t *testing.T) {
	s:=newServer()
	req:=httptest.NewRequest(http.MethodPost,"/api/v1/nodes/node-01/media/attach",
		strings.NewReader(`{"image":"demo-os"}`))
	w:=httptest.NewRecorder()
	s.nodeHandler(w,req)
	if w.Code!=200 || s.nodes["node-01"].AttachedImage!="demo-os" {
		t.Fatalf("attach failed: status=%d image=%q",w.Code,s.nodes["node-01"].AttachedImage)
	}
}

func TestProvision(t *testing.T) {
	s:=newServer()
	req:=httptest.NewRequest(http.MethodPost,"/api/v1/nodes/node-01/provision",
		strings.NewReader(`{"image":"demo-os"}`))
	w:=httptest.NewRecorder()
	s.nodeHandler(w,req)
	n:=s.nodes["node-01"]
	if w.Code!=200 || n.Power!="on" || n.Health!="healthy" || n.AttachedImage!="" {
		t.Fatalf("unexpected final node state: %+v status=%d",n,w.Code)
	}
}
