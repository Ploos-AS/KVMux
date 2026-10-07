package controller

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

type scriptedRW struct {
	w bytes.Buffer
	r *strings.Reader
}
func (s *scriptedRW) Write(p []byte) (int,error) { return s.w.Write(p) }
func (s *scriptedRW) Read(p []byte) (int,error) { return s.r.Read(p) }

func TestPing(t *testing.T) {
	rw:=&scriptedRW{r:strings.NewReader("1 OK PONG\n")}
	c:=New(rw)
	if err:=c.Ping(); err!=nil { t.Fatal(err) }
	if got:=rw.w.String(); got!="1 PING\n" { t.Fatalf("wire=%q",got) }
}

func TestPowerCycleEncoding(t *testing.T) {
	rw:=&scriptedRW{r:strings.NewReader("1 OK DONE\n")}
	c:=New(rw)
	if err:=c.PowerCycle(2,250*time.Millisecond); err!=nil { t.Fatal(err) }
	if got:=rw.w.String(); got!="1 POWER_CYCLE 2 250\n" { t.Fatalf("wire=%q",got) }
}

func TestRejectMismatchedID(t *testing.T) {
	rw:=&scriptedRW{r:strings.NewReader("99 OK PONG\n")}
	if err:=New(rw).Ping(); err==nil { t.Fatal("expected id mismatch") }
}

func TestControllerError(t *testing.T) {
	rw:=&scriptedRW{r:strings.NewReader("1 ERR INVALID_PORT\n")}
	if err:=New(rw).PowerOn(9); err==nil { t.Fatal("expected controller error") }
}

var _ io.ReadWriter = (*scriptedRW)(nil)
