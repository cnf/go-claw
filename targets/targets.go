package targets

import "strings"

import "github.com/cnf/go-claw/clog"

// Target is an interface which every Target must implement
type Target interface {
    SendCommand(repeated int, cmd string, args ...string) error
    Stop() error
    Commands() map[string]*Command
}


func init() {
    RegisterTarget("claw", createClawTarget);
}

// CreateTarget is a function definition each target must provide during registration
type CreateTarget func(name string, params map[string]string) (t Target, err error)


// targetlist is the global internal list of all registered targets
var targetlist = make(map[string]CreateTarget)

// RegisterTarget is used by targets to register itself
func RegisterTarget(name string, creator CreateTarget) {
    name = strings.ToLower(name)
    clog.Info("Registering target: %s", name)

    if targetlist[name] != nil {
        panic("RegisterTarget: target name '" + name + "' already exists! Pick another name for your module!")
    }

    targetlist[name] = creator
}

