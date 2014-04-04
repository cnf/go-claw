package modes

import "fmt"

import "github.com/cnf/go-claw/clog"

// Mode holds the data for a single mode
type Mode struct {
    Keys map[string][]string
    Entry []string
    Exit []string
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
    if m.ModeMap[mode] == nil {
        return actions, fmt.Errorf("no such mode found: %s", mode)
    }
    if m.active != nil {
        actions = append(actions, m.active.Exit...)
    }
    m.active = m.ModeMap[mode]
    m.name = mode
    actions = append(actions, m.active.Entry...)

    clog.Debug("Modes: `%s` is now active", mode)

    return actions, nil
}

// Setup sets up a new mode structure
func (m *Modes) Setup(modelist map[string]*Mode) error {
    m.ModeMap = make(map[string]*Mode)
    var err error
    for k, v := range modelist {
        // clog.Info("Setting up mode: %s", k)
        err = m.AddMode(k, v)
        if err != nil { break }
    }
    if m.def == nil {
        m.def = &Mode{}
    }
    if err != nil {
        return fmt.Errorf("not all modes setup: %s", err)
    }
    return nil
}

// AddMode adds a mode to the list
func (m *Modes) AddMode(name string, mode *Mode) error {
    clog.Info("Setting up mode: %s", name)
    // m.ModeMap[name] = &Mode{Keys: mode.Keys, entry: mode.entry, exit: mode.exit}
    m.ModeMap[name] = mode
    if name == "default" {
        m.def = m.ModeMap[name]
    }
    return nil
}

// DelMode removes a mode from the list
func (m *Modes) DelMode(name string) error {
    if name == "default" {
        return fmt.Errorf("can not delete `%s` mode", name)
    }
    delete(m.ModeMap, name)
    return nil
}
