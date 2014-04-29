package config

import "github.com/cnf/go-claw/modes"

// Config is the main configuration object
type Config struct {
    Listeners map[string]*ListenerConfig `json:"listeners"`
    Modes map[string]*modes.Mode `json:"modes"`
    Targets map[string]*TargetConfig `json:"targets"`
    cfgfile string
    verbose bool
    httpport int
}

// ListenerConfig holds Listener configuration
type ListenerConfig struct {
    Module string
    Params map[string]string
}

// Actionlist has a list of actions
type Actionlist []string

// TargetConfig holds Target configuration
type TargetConfig struct {
    Module string
    Params map[string]string
}
