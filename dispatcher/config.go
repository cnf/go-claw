package dispatcher

// type Config struct {
//     cfgfile string
//     Home string
//     System System
// }

type Config struct {
    // Modes map[string]map[string][]string `json:"mode"`
    Listeners map[string]ConfigListener
    Modes map[string]*Mode
    Targets map[string]ConfigTarget
}

type ConfigListener struct {
    Module string
    Params map[string]string
}

type ConfigMode map[string]Actionlist

type Actionlist []string

type ConfigTarget struct {
    Module string
    Params map[string]string
}
