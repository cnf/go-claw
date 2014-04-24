package config

// import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/modes"
// import "github.com/cnf/go-claw/targets"

// Config is the main configuration object
type Config struct {
    Listeners map[string]*ConfigListener `json:"listeners"`
    Modes map[string]*modes.Mode `json:"modes"`
    Targets map[string]*ConfigTarget `json:"targets"`
    cfgfile string
    verbose bool
}

// ConfigListener holds Listener configuration
type ConfigListener struct {
    Module string
    Params map[string]string
}

// ConfigMode holds Mode configuration
type ConfigMode map[string]Actionlist

// Actionlist has a list of actions
type Actionlist []string

// ConfigTarget holds Target configuration
type ConfigTarget struct {
    Module string
    Params map[string]string
}
