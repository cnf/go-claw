package plex

import "net"

import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

type Plex struct {
    name string
    addr *net.TCPAddr
}

func Register() {
    targets.RegisterTarget("plex", Create)
}

func Create(name string, params map[string]string) (t targets.Target, ok bool) {
    clog.Debug("Plex Create called")
    return &Plex{name: name}, true
    return nil, false
}

func (self *Plex) SendCommand(cmd string, args ...string) bool {
    switch cmd {
    case "MenuUp":
        // return self.restSend()
    }
    return false
}
