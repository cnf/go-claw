package dispatcher

import "time"

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/modes"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/clog"

// Dispatcher holds all the dispatcher info
type Dispatcher struct {
    configfile string
    config Config
    keytimeout time.Duration
    listenermap map[string]*listeners.Listener
    targetmanager *targets.TargetManager
    modes *modes.Modes
    activemode string
    cs *listeners.CommandStream
}


func (d *Dispatcher) Setup(configfile string) error {
    d.configfile = configfile
    d.activemode = "default"
    d.keytimeout = time.Duration(120 * time.Millisecond)
    d.readConfig()
    d.setupListeners()
    d.setupModes()
    d.setupTargets()

    return nil
}

func (d *Dispatcher) Start() {
    defer d.cs.Close()

    d.startListeners()
    d.startTargets()

    var out listeners.RemoteCommand

    for d.cs.Next(&out) {
        if d.cs.HasError() {
            clog.Warn("An error occured somewhere: %v", d.cs.GetError())
            d.cs.ClearError()
        }
        d.dispatch(&out)
    }
}

func (d *Dispatcher) startTargets() error {
    return d.targetmanager.Start()
}

func (d *Dispatcher) startListeners() error {
    return nil
}

func (d *Dispatcher) setupListeners() error {
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
    return nil
}

func (d *Dispatcher) setupModes() error {
    d.modes = &modes.Modes{}
    err := d.modes.Setup(d.config.Modes)
    if err != nil {
        clog.Error("Dispatcher: could not set up modes: %s", err)
    }
    return nil
}

func (d *Dispatcher) setupTargets() error {
    if d.targetmanager == nil {
        d.targetmanager = targets.NewTargetManager(d.modes)
    } else {
        // Stop and remove all targets if needed
        d.targetmanager.Stop()
    }

    for k, v := range d.config.Targets {
        clog.Info("Setting up target: %s", k)
        if err := d.targetmanager.Add(v.Module, k, v.Params); err != nil {
            clog.Warn("Could not add target '%s:%s': %s", v.Module, k, err.Error())
        }
    }
    return nil
}

func (d *Dispatcher) dispatch(rc *listeners.RemoteCommand) bool {
    clog.Debug("Dispatch: repeat `%2d` - key `%s` - source `%s`", rc.Repeat, rc.Key, rc.Source)
    tdiff := time.Since(rc.Time)
    clog.Debug("Dispatch: --> t: %s", tdiff.String())
    if tdiff > d.keytimeout {
        clog.Info("dispatch: Key timeout reached: %# v", tdiff.String())
        return false
    }
    var rok = true
    actions, err := d.modes.ActionsFor(rc.Key)
    if err != nil {
        clog.Debug("dispatch:ActionsFor: %s", err)
        return false
    }
    for _, v := range actions {
        err := d.targetmanager.RunCommand(rc.Repeat, v)
        if err != nil {
            rok = false
            clog.Debug("dispatch:RunCommand: %s", err)
        }
    }
    return rok
}
