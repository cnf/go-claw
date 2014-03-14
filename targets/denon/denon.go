package denon

import "net"
import "fmt"
import "time"

import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

type Denon struct {
    name string
    addr *net.TCPAddr
    commands map[string]Commander
    last time.Time
    wait time.Duration
}

func Register() {
    targets.RegisterTarget("denon", Create)
}

func Create(name string, params map[string]string) (t targets.Target, ok bool) {
    // TODO: VALIDATE PARAMS
    if val, ok := params["address"]; ok {
        d := setup(name, val, 23)
        d.commands = AVRX2000
        d.wait = time.Duration(100 * time.Millisecond)
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
    default:
        clog.Debug("Looking up %s in the map", cmd)
        // clog.Debug(">>>>>>>>> %# v", self.commands)
        if val, ok := self.commands[cmd]; ok {
            cstr, err := val.Command(args...)
            if err != nil {
                return false
            }
            clog.Debug(">>>>>>>>> %# v", cstr)
            return self.socketSend(cstr)
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

    tdiff := time.Since(self.last)
    if tdiff < self.wait {
        // time.Sleep(self.wait)
        clog.Debug("+++++ Waiting %# v", self.wait - tdiff)
        time.Sleep(self.wait - tdiff)
    }

    conn, err := net.DialTCP("tcp", nil, self.addr)
    if err != nil {
        clog.Error("Connection failed: %s", err)
        if conn != nil {
            conn.Close()
        }
        return false
    }
    clog.Debug("Sending %s to %s", str, self.name)
    fmt.Fprintf(conn, "%s\r", str)
    conn.Close()
    self.last = time.Now()
    return true
}
