package targets

//import "fmt"
import "github.com/cnf/go-claw/clog"

type Target interface {
    SendCommand(cmd string, args ...string) error
    Stop() error
    Commands() map[string]*Command
}


func init() {
    RegisterTarget("mode", createModeHandler);
}

type CreateTarget func(name string, params map[string]string) (t Target, err error)
var targetlist = make(map[string]CreateTarget)

func RegisterTarget(name string, creator CreateTarget) {
    clog.Info("Registering target: %s", name)

    if targetlist[name] != nil {
        panic("RegisterTarget: error " + name + " already exists! Pick another name for your module!")
    }

    targetlist[name] = creator
}

