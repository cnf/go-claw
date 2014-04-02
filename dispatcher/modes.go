package dispatcher

import "fmt"
import "github.com/cnf/go-claw/clog"

import "github.com/kr/pretty"

// var active int

// Mode holds the data for a single mode
type Mode struct {
    Keys map[string][]string
    entry []string
    exit []string
}

// Modes holds all the modes data
type Modes struct {
    name string
    active *Mode
    def *Mode
    ModeMap map[string]*Mode
}

// ActionsFor returns a list of actions for a specific key
func (m *Modes) ActionsFor(key string) ([]string, error) {
    var actions []string
    if (m.active == nil) && (m.def == nil) {
        return nil, fmt.Errorf("no modes found")
    }
    if m.active != nil && m.active.Keys[key] != nil {
        actions = m.active.Keys[key]
    } else if m.def.Keys[key] != nil {
        actions = m.def.Keys[key]
    } else {
        return nil, fmt.Errorf("key `%s` not found", key)
    }
    return actions, nil
}

// SetActive text
func (m *Modes) SetActive(mode string) ([]string, error) {
    var actions []string
    if m.active != nil {
        actions = append(actions, m.active.exit...)
    }
    m.active = m.ModeMap[mode]
    m.name = mode
    actions = append(actions, m.active.entry...)

    return actions, nil
}

// Setup sets up a new mode structure
func (m *Modes) Setup(modelist map[string]*Mode) bool {
    m.ModeMap = make(map[string]*Mode)
    clog.Debug("%# v", modelist)
    for k, v := range modelist {
        // clog.Info("Setting up mode: %s", k)
        m.AddMode(k, v)
    }
    if m.def == nil {
        clog.Debug("No default mode set!")
    }
    return false
}

// AddMode adds a mode to the list
func (m *Modes) AddMode(name string, mode *Mode) {
    // if name != "default" {return}
    clog.Debug("Adding mode %s", name)
    // pretty.Print(mode)
    m.ModeMap[name] = &Mode{Keys: mode.Keys, entry: mode.entry, exit: mode.exit}
    pretty.Print(mode)
    if name == "default" {
        m.def = m.ModeMap[name]
    }
    // d.modemap[k] = &Mode{Keys: make(map[string][]string)}
    // for k, _ := range mode.Keys {
        // println("++++++", k)
    // }
}
