package dispatcher

import "time"

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/clog"

type Dispatcher struct {
    Configfile string
    config Config
    keytimeout time.Duration
    listenermap map[string]*listeners.Listener
    targetmanager *targets.TargetManager
    modemap map[string]*Mode
    activemode string
    cs *listeners.CommandStream
}

func (d *Dispatcher) Start() {
    defer d.cs.Close()
    d.keytimeout = time.Duration(120 * time.Millisecond)
    d.readConfig()
    d.setupListeners()
    d.setupModes()
    d.setupTargets()

    var out listeners.RemoteCommand

    for d.cs.Next(&out) {
        if d.cs.HasError() {
            clog.Warn("An error occured somewhere: %v", d.cs.GetError())
            d.cs.ClearError()
        }
        d.dispatch(&out)
    }
}

func (d *Dispatcher) setupListeners() {
    d.listenermap = make(map[string]*listeners.Listener)
    d.cs = listeners.NewCommandStream()

    for k, v := range d.config.Listeners {
        l, ok := listeners.GetListener(v.Module, v.Params)
        if ok {
            clog.Info("Setting up listener: %s", k)
            d.listenermap[k] = &l
            d.cs.AddListener(l)
        }
    }

}

func (d *Dispatcher) setupModes() {
    d.modemap = make(map[string]*Mode)

    for k, v := range d.config.Modes {
        clog.Info("Setting up mode: %s", k)
        d.modemap[k] = &Mode{Keys: make(map[string][]string)}
        for kk, kv := range v {
            d.modemap[k].Keys[kk] = make([]string, len(kv))
            i := 0
            for _, av := range kv {
                d.modemap[k].Keys[kk][i] = av
                i++
            }
        }
    }

    // Ensure the "default" mode map exists
    if _, ok := d.modemap["default"]; !ok {
        clog.Warn("Warning: 'default' modemap did not exist - adding dummy")
        d.modemap["default"] = &Mode{Keys: make(map[string][]string)}
    }
    d.setMode("default")
}

func (d *Dispatcher) setupTargets() {
    if d.targetmanager == nil {
        d.targetmanager = targets.NewTargetManager()
    }
    // Stop and remove all targets if needed
    d.targetmanager.Stop()

    for k, v := range d.config.Targets {
        clog.Info("Setting up target: %s", k)
        if err := d.targetmanager.Add(v.Module, k, v.Params); err != nil {
            clog.Warn("Could not add target '%s:%s': %s", v.Module, k, err.Error())
        }
    }
}

func (d *Dispatcher) dispatch(rc *listeners.RemoteCommand) bool {
    clog.Debug("Dispatch: repeat `%2d` - key `%s` - source `%s`", rc.Repeat, rc.Key, rc.Source)
    tdiff := time.Since(rc.Time)
    clog.Debug("Dispatch: --> t: %s", tdiff.String())
    if tdiff > d.keytimeout {
        clog.Info("Dispatch: Key timeout reached: %# v", tdiff.String())
        return false
    }
    var rok = true
    // FIXME: NEED NEW MODES!
    if _, ok := d.modemap[d.activemode]; !ok {
        d.setMode("default")
    }
    if val, ok := d.modemap[d.activemode].Keys[rc.Key]; ok {
        for _, v := range val {
            if err := d.targetmanager.RunCommand(v); err != nil {
                clog.Warn("Command failed: %s", v)
                rok = false
            }
        }
    } else if val, ok := d.modemap["default"].Keys[rc.Key]; ok {
        for _, v := range val {
            if err := d.targetmanager.RunCommand(v); err != nil {
                clog.Warn("Command failed: %s", v)
                rok = false
            }
        }
    } else {
        clog.Warn("Dispatch: key `%s` Not found in any mode.", rc.Key)
        rok = false
    }

    return rok
}

func (d *Dispatcher) setMode(mode string) bool {
    if _, ok := d.modemap[mode]; ok {
        clog.Info("Dispatch: mode changed to `%s`", mode)
        d.activemode = mode
        return true
    }
    return false
}
