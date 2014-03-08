package targets

import "github.com/cnf/go-claw/clog"

type Target interface {
    SendCommand(cmd string, args ...string) bool
}

type CreateTarget func(name string, params map[string]string) (t Target, ok bool)

var list = make(map[string]CreateTarget)

func RegisterTarget(name string, creator CreateTarget) {
    list[name] = creator
}

func GetTarget(module string, name string, params map[string]string) (t Target, ok bool) {
    if _, ok := list[module]; ok {
        return list[module](name, params)
    }
    clog.Warn("Tried to use Target module %s, but module does not exist.", module)
    return nil, false
}
