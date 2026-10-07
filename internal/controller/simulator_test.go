package controller

import ("strings";"testing")

func TestSimulatorTwoPorts(t *testing.T){
 s:=&Simulator{}
 for _,tc:=range []struct{in,want string}{
 {"1 PING","1 OK PONG"},
 {"2 POWER_ON 1","2 OK DONE"},
 {"3 POWER_ON 2","3 OK DONE"},
 {"4 RESET 1 50","4 OK DONE"},
 {"5 SERVICE_SELECT 2","5 OK DONE"},
 {"6 POWER_OFF 1","6 OK DONE"},
 {"7 RESET 1 50","7 ERR POWERED_OFF"},
 {"8 POWER_ON 3","8 ERR INVALID_PORT"},
 {"9 POWER_CYCLE 2 10","9 ERR INVALID_INTERVAL"},
 {"10 ALL_SAFE","10 OK SAFE"},
 } {if got:=s.Execute(tc.in);got!=tc.want{t.Fatalf("%s: %q != %q",tc.in,got,tc.want)}}
 if s.Power[0]||!s.Power[1]||s.Selected!=0{t.Fatalf("unexpected state: %+v",s)}
}
func TestSimulatorStream(t *testing.T){
 s:=&Simulator{};var out strings.Builder
 if err:=s.Serve(strings.NewReader("1 PING\n2 POWER_ON 1\n"),&out);err!=nil{t.Fatal(err)}
 if out.String()!="1 OK PONG\n2 OK DONE\n"{t.Fatalf("output=%q",out.String())}
}
