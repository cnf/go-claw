package targets

import "fmt"

import "github.com/cnf/go-claw/clog"

type Target interface {
    SendCommand(cmd string, args ...string) error
}

type CreateTarget func(name string, params map[string]string) (t Target, err error)

var list = make(map[string]CreateTarget)

func RegisterTarget(name string, creator CreateTarget) {
    list[name] = creator
}

func GetTarget(module string, name string, params map[string]string) (Target, error) {
    if _, ok := list[module]; ok {
        t, err := list[module](name, params)
        if err != nil { return nil, err }
        return t, nil
    }
    clog.Warn("Tried to use Target module %s, but module does not exist.", module)
    return nil, fmt.Errorf("target module %s does not exist", module)
}
