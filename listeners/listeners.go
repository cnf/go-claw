package listeners

import "github.com/cnf/go-claw/clog"

type Listener interface {
    RunListener(cs *CommandStream)
}

type CreateListener func(params map[string]string) (l Listener, ok bool)

var list = make(map[string]CreateListener)

func RegisterListener(name string, creator CreateListener) {
    list[name] = creator
}

func GetListener(name string, params map[string]string) (l Listener, ok bool) {
    if _, ok := list[name]; ok {
        return list[name](params)
    }
    clog.Warn("Listener `%s` does not exist", name)
    return nil, false
}
