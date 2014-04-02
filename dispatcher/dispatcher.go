package dispatcher

import "strings"
import "time"

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/clog"

// Dispatcher holds all the dispatcher info
type Dispatcher struct {
    Configfile string
    config Config
    keytimeout time.Duration
    listenermap map[string]*listeners.Listener
    targetmap map[string]targets.Target
    modemap map[string]*Mode
    modes *Modes
    activemode string
    cs *listeners.CommandStream
}

func (d *Dispatcher) Start() {
    defer d.cs.Close()
    d.activemode = "default"
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
    d.modes = &Modes{}
    err := d.modes.Setup(d.config.Modes)
    if err != nil {
        clog.Error("Dispatcher: could not set up modes: %s", err)
    }
}

func (d *Dispatcher) setupTargets() {
    d.targetmap = make(map[string]targets.Target)
    for k, v := range d.config.Targets {
        t, err := targets.GetTarget(v.Module, k, v.Params)
        if err != nil {
            clog.Warn("Dispatcher: %s", err)
        }
        d.targetmap[k] = t
        clog.Info("Setting up target: %s", k)
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
    var mod, cmd, args string
    var rok bool
    actions, err := d.modes.ActionsFor(rc.Key)
    if err != nil {
        clog.Debug("Dispatcher: %s", err)
        return false
    }
    for _, v := range actions {
        mod, cmd, args, rok = d.resolve(v)
        if mod == "mode" {
            cmdlist, err := d.modes.SetActive(cmd)
            if err != nil {
                return false
            }
            for _, v := range cmdlist {
                clog.Debug("Dispatcher: running exit/entry")
                mod, cmd, args, rok = d.resolve(v)
                d.sender(mod, cmd, args)
            }
        }
        d.sender(mod, cmd, args)
    }
    if !rok {
        return false
    }

    return true
}

func (d *Dispatcher) resolve(input string) (mod string, cmd string, args string, ok bool) {
    foo := strings.SplitN(input, "::", 2)
    if len(foo) < 2 {
        clog.Warn("Dispatch: `%s` is not a well formed command", input)
        return "", "", "", false
    }
    bar := strings.SplitN(foo[1], " ", 2)
    baz := ""
    if len(bar) > 1 {
        baz = bar[1]
    }

    return foo[0], bar[0], baz, true
}

func (d *Dispatcher) sender(mod string, cmd string, args string) bool {
    if mod == "mode" {
        return d.setMode(cmd)
    }
    if t, ok := d.targetmap[mod]; ok {
        err := t.SendCommand(cmd, args)
        if err != nil {
            clog.Debug("Dispatch: failed to send command `%s` for `%s` with error: %s", cmd, mod, err)
        }
        return true
    }
    return false
}

func (d *Dispatcher) setMode(mode string) bool {
    if _, ok := d.modemap[mode]; ok {
        clog.Info("Dispatch: mode changed to `%s`", mode)
        d.activemode = mode
        return true
    }
    return false
}
