package targets

//import "fmt"
import "github.com/cnf/go-claw/clog"

type Target interface {
    SendCommand(cmd string, args ...string) error
    Stop() error
    Commands() map[string]*Command
}

type CreateTarget func(name string, params map[string]string) (t Target, err error)
var targetlist = make(map[string]CreateTarget)

func RegisterTarget(name string, creator CreateTarget) {
    clog.Info("Registering target: %s", name)
    targetlist[name] = creator
}

