package dispatcher

// import "github.com/cnf/go-claw/clog"

// var modes map[string]*Mode
var active int

type Mode struct {
    Keys map[string][]string
}

// type Key []Action

// type Action string

// func init() {
    // modes = make(map[string]*Mode)
    // modes["default"] = &Mode{name: "default", id: 1}
    // modes["bar"] = &Mode{name: "bar", id: 2}
    // active = 1
// }

// func SetMode(name string) bool {
//     mode, ok := modes[name]
//     if ok {
//         active = mode.id
//         clog.Debug("Setting active mode to %s", mode.name)
//         return true
//     }
//     return false
// }

// func GetMode() string {
//     for _, value := range modes {
//         if value.id == active {
//             return value.name
//         }
//     }
//     return ""
// }

// func (self *Mode) SetActive() bool {
//     return false
// }
