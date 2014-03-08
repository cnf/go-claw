package listeners

import "github.com/cnf/go-claw/commandstream"

type Listener interface {
    RunListener(cs *commandstream.CommandStream)
}

type CreateListener func(params map[string]string) (l Listener, ok bool)

var list = make(map[string]CreateListener)

func RegisterListener(name string, creator CreateListener) {
    list[name] = creator
}

func GetListener(name string, params map[string]string) (l Listener, ok bool) {
    if _, ok := list[name]; ok {
        println(name, "exists")
        return list[name](params)
    }
    println(name, "does not exist")
    return nil, false
}
