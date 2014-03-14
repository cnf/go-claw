package denon

import "net"
import "fmt"
import "time"
import "errors"
import "strings"

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
    case "Mute":
        return self.toggleMute()
    default:
        // clog.Debug(">>>>>>>>> %# v", self.commands)
        cstr, err := self.getCommand(cmd, args...)
        if err != nil { return false }
        clog.Debug(">>>>>>>>> %# v", cstr)
        _, serr := self.socketSend(cstr)
        if serr != nil { return false }
    }
    return false
}

func (self *Denon) Capabilities() []string {
    return []string{}
}

func (self *Denon) getCommand(cmd string, args ...string) (string, error) {
    clog.Debug("Looking up %s in the map", cmd)
    if val, ok := self.commands[cmd]; ok {
        cstr, err := val.Command(args...)
        if err != nil {
            return "", err
        }
        return cstr, nil
    }
    return "", errors.New("Could not get command")
}

func (self *Denon) socketSend(str string) (cmd string, err error) {
    if self.addr == nil {
        clog.Warn("No address to sent Denon command to.")
        return "", errors.New("No address set")
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
        return "", err
    }
    conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
    clog.Debug("Sending %s to %s", str, self.name)
    fmt.Fprintf(conn, "%s\r", str)
    reply := make([]byte, 32)
    l, err := conn.Read(reply)
    if err != nil { return "", err }
    conn.Close()
    self.last = time.Now()
    return string(reply[0:l]), nil
}

func (self *Denon) toggleMute() bool {
    r, err := self.socketSend("MU?")
    // MuteOn
    if err != nil { return false }
    r = strings.TrimSpace(r)
    clog.Debug("Found: %# v", r)
    // _, serr := self.socketSend("MU?")
    // if serr != nil { return false }
    if r == "MUOFF" {
        clog.Debug("MU OFF!")
        cmd, err := self.getCommand("MuteOn")
        if err != nil { return false }
        _, serr := self.socketSend(cmd)
        if serr != nil { return false }
    } else if r == "MUON" {
        clog.Debug("MU ON!")
        cmd, err := self.getCommand("MuteOff")
        if err != nil { return false }
        _, serr := self.socketSend(cmd)
        if serr != nil { return false }
        // _, serr := self.socketSend("MUOFF")
    }
    return true
}
