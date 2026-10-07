package controller

import (
 "io"
 "net"
 "testing"
 "time"
)

func TestClientSimulatorEndToEnd(t *testing.T) {
 host,device:=net.Pipe()
 defer host.Close()
 defer device.Close()
 sim:=&Simulator{}
 done:=make(chan error,1)
 go func(){done<-sim.Serve(device,device)}()
 client:=New(host)
 if err:=client.Ping();err!=nil{t.Fatal(err)}
 if err:=client.PowerOn(1);err!=nil{t.Fatal(err)}
 if err:=client.PowerOn(2);err!=nil{t.Fatal(err)}
 if err:=client.Reset(1,50*time.Millisecond);err!=nil{t.Fatal(err)}
 if err:=client.PowerCycle(2,250*time.Millisecond);err!=nil{t.Fatal(err)}
 if err:=client.ServiceSelect(2);err!=nil{t.Fatal(err)}
 if err:=client.AllSafe();err!=nil{t.Fatal(err)}
 if err:=client.PowerOff(1);err!=nil{t.Fatal(err)}
 if err:=client.Reset(1,50*time.Millisecond);err==nil{t.Fatal("reset on powered-off port accepted")}
 if err:=client.PowerOn(3);err==nil{t.Fatal("invalid port accepted")}
 host.Close()
 select {case err:=<-done:if err!=nil&&err!=io.EOF{t.Fatal(err)}
 case <-time.After(time.Second):t.Fatal("simulator did not exit")}
 if sim.Power[0]||!sim.Power[1]||sim.Selected!=0{t.Fatalf("unexpected final state: %+v",sim)}
}
