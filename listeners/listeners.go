package listeners

// import "github.com/cnf/go-claw/setup"
import "github.com/cnf/go-claw/commandstream"
// import "fmt"

type Listener interface {
    RunListener(cs *commandstream.CommandStream)
}

type CreateListener func(params map[string]string) (l Listener, ok bool)

var list = make(map[string]CreateListener, 5)

func RegisterListener(name string, creator CreateListener) {
    list[name] = creator
}

// func Testing() {
//     for key, value := range list {
//         fmt.Printf("%s -> %# v\n", key, value)
//     }
// }

func GetListener(name string, params map[string]string) (l Listener, ok bool) {
    if _, ok := list[name]; ok {
        println(name, "exists")
        return list[name](params)
        // l := val(params)
        // return l, true
    }
    println(name, "does not exist")
    return nil, false
}

// func MakeListener(name string, params map[string]string, cs *CommandStream) bool {
    // list[name]("foo", params)
    // return true
// }

// func ProcessListenerConfig(cs *CommandStream, config map[string]setup.Listener) {
// }
