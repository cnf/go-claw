package denon

import "net"
import "fmt"
import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

type Denon struct {
    name string
    addr *net.TCPAddr
    commands map[string]Commander
}

func Register() {
    targets.RegisterTarget("denon", Create)
}

func Create(name string, params map[string]string) (t targets.Target, ok bool) {
    // TODO: VALIDATE PARAMS
    if val, ok := params["address"]; ok {
        d := setup(name, val, 23)
        d.commands = AVRX2000
        return d, true
    }
    return nil, false
}

func setup(name string, host string, port int) *Denon {
    clog.Debug("Initializing Denon")
    tmp, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
    if err != nil {
        clog.Error("Could not Initialize Denon: %s", err)
        return nil
    }
    return &Denon{addr: tmp, name: name}

}

func (self *Denon) SendCommand(cmd string, args ...string) bool {
    switch cmd {
    // case "VolumeUp":
    //   return self.socketSend("MVUP")
    // case "VolumeDown":
    //   return self.socketSend("MVDOWN")
    // case "Volume":
      // return self.socketSend(fmt.Sprintf("MV%s", args[0]))
    default:
        clog.Debug("Looking up %s in the map", cmd)
        // clog.Debug(">>>>>>>>> %# v", self.commands)
        if val, ok := self.commands[cmd]; ok {
            cmd, err := val.Command(args...)
            if err != nil {
                return false
            }
            clog.Debug(">>>>>>>>> %# v", cmd)
            // return self.socketSend(val.command["send"])
        }
    }
    return false
}

func (self *Denon) Capabilities() []string {
    return []string{}
}

func (self *Denon) socketSend(str string) bool {
    if self.addr == nil {
        clog.Warn("No address to sent Denon command to.")
        return false
}
    conn, err := net.DialTCP("tcp", nil, self.addr)
    // defer conn.Close()
    if err != nil {
        clog.Error("Connection failed: %s", err)
        if conn != nil {
            conn.Close()
        }
        return false
    }
    fmt.Fprintf(conn, "%s\r", str)
    conn.Close()
    return true
}

// var commands map[string]command
// 
// type Cmd struct {
//     Send string
//     Args []string
// }
// 
// 
// type command struct {
//     send string
//     params []string
// }
