package dispatcher

// type Config struct {
//     cfgfile string
//     Home string
//     System System
// }

type Config struct {
    // Modes map[string]map[string][]string `json:"mode"`
    Listeners map[string]Listener
    Modes map[string]Mode
    Targets map[string]Target
}

type Listener struct {
    Module string
    Params map[string]string
}

type Mode map[string]Actionlist

type Actionlist []string

type Target struct {
    Module string
    Params map[string]string
}
