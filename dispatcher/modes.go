package dispatcher

import "github.com/cnf/go-claw/clog"

var modes map[string]*mode
var active int

type mode struct {
    name string
    id int
}

func init() {
    modes = make(map[string]*mode)
    modes["default"] = &mode{name: "default", id: 1}
    modes["bar"] = &mode{name: "bar", id: 2}
    active = 1
}

func SetMode(name string) bool {
    mode, ok := modes[name]
    if ok {
        active = mode.id
        clog.Debug("Setting active mode to %s", mode.name)
        return true
    }
    return false
}

func GetMode() string {
    for _, value := range modes {
        if value.id == active {
            return value.name
        }
    }
    return ""
}

func (self *Mode) SetActive() bool {
    return false
}
