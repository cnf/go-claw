package linux

import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/tools"

type Linux struct {
    name string
    wol string
}

func Register() {
    targets.RegisterTarget("linux", Create)
}

func Create(name string, params map[string]string) (t targets.Target, ok bool) {
    l := &Linux{name: name}
    if wol, ok := params["wol"]; ok {
        l.wol = wol
    }
    return l, true
}

func (l *Linux) SendCommand(cmd string, args ...string) bool {
    switch cmd {
    case "PowerOn":
        clog.Debug("Power on %s", l.name)
        return l.powerOn()
    }
    return false
}

func (l *Linux) powerOn() bool {
    if l.wol != "" {
        return tools.Wol(l.wol)
    }
    clog.Debug("Can not power on %s", l.name)
    return false
}
