package controller

import (
 "bufio"
 "fmt"
 "io"
 "strconv"
 "strings"
 "sync"
)

type Simulator struct { mu sync.Mutex; Power [2]bool; Selected int }
func (s *Simulator) Serve(r io.Reader,w io.Writer) error {
 scan:=bufio.NewScanner(r)
 for scan.Scan() { fmt.Fprintln(w,s.Execute(scan.Text())) }
 return scan.Err()
}
func (s *Simulator) Execute(line string) string {
 s.mu.Lock(); defer s.mu.Unlock()
 f:=strings.Fields(line)
 if len(f)<2 { return "0 ERR MALFORMED" }
 id,cmd:=f[0],f[1]
 reply:=func(v string)string{return id+" "+v}
 if cmd=="PING" && len(f)==2 { return reply("OK PONG") }
 if cmd=="ALL_SAFE" && len(f)==2 { s.Selected=0;return reply("OK SAFE") }
 if len(f)<3 { return reply("ERR INVALID_ARGS") }
 port,err:=strconv.Atoi(f[2]);if err!=nil||port<1||port>2{return reply("ERR INVALID_PORT")}
 switch cmd {
 case "POWER_ON":if len(f)!=3{return reply("ERR INVALID_ARGS")};s.Power[port-1]=true
 case "POWER_OFF":if len(f)!=3{return reply("ERR INVALID_ARGS")};s.Power[port-1]=false
 case "POWER_CYCLE":
  if len(f)!=4{return reply("ERR INVALID_ARGS")}
  ms,e:=strconv.Atoi(f[3]);if e!=nil||ms<100||ms>60000{return reply("ERR INVALID_INTERVAL")}
  // Logical simulation only: physical firmware must enforce real off-time.
  s.Power[port-1]=true
 case "RESET":
  if len(f)!=4{return reply("ERR INVALID_ARGS")}
  ms,e:=strconv.Atoi(f[3]);if e!=nil||ms<1||ms>5000{return reply("ERR INVALID_PULSE")}
  if !s.Power[port-1]{return reply("ERR POWERED_OFF")}
 case "SERVICE_SELECT":if len(f)!=3{return reply("ERR INVALID_ARGS")};s.Selected=port
 default:return reply("ERR UNKNOWN_COMMAND")
 }
 return reply("OK DONE")
}
