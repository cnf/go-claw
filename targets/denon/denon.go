package denon

import "net"
import "fmt"
import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

type Denon struct {
    name string
    addr *net.TCPAddr
}

func Register() {
    targets.RegisterTarget("denon", Create)
}

func Create(name string, params map[string]string) (t targets.Target, ok bool) {
    // TODO: VALIDATE PARAMS
    if val, ok := params["address"]; ok {
        d := setup(name, val, 23)
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
    case "VolumeUp":
      return self.socketSend("MVUP")
    case "Volume":
      return self.socketSend(fmt.Sprintf("MV%s", args[0]))
    }
    return false
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
        defer conn.Close()
        return false
    }
    fmt.Fprintf(conn, "%s\r", str)
    conn.Close()
    return true
}
