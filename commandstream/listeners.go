package commandstream

import "github.com/cnf/go-claw/setup"
import "fmt"

type Listener interface {
    Setup(cs *CommandStream) bool
    RunListener(cs *CommandStream)
}

type CreateListener func(ptype string, params map[string]string) (l Listener, ok bool)

var list = make(map[string]CreateListener, 5)

func RegisterListener(name string, creator CreateListener) {
    list[name] = creator
}

func Testing() {
    for key, value := range list {
        fmt.Printf("%s -> %# v\n", key, value)
    }
}

func MakeListener(name string, params map[string]string, cs *CommandStream) bool {
    list[name]("foo", params)
    return true
}

func ProcessListenerConfig(cs *CommandStream, config map[string]setup.Listener) {
}
