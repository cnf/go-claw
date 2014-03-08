package dispatcher

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/clog"

type Dispatcher struct {
    Configfile string
    config Config
    listenermap map[string]*listeners.Listener
    targetmap map[string]targets.Target
    modemap map[string]*Mode
    cs *listeners.CommandStream
}

func (self *Dispatcher) Start() {
    println("Starting")
    defer self.cs.Close()
    self.readConfig()
    println("================")
    self.setupListeners()
    println("================")
    self.setupModes()
    println("================")
    self.setupTargets()
    println("================")

    var out listeners.RemoteCommand

    for self.cs.Next(&out) {
        if self.cs.HasError() {
            clog.Warn("An error occured somewhere: %v", self.cs.GetError())
            self.cs.ClearError()
        }
        clog.Debug("repeat: %2d - key: %s - source: %s", out.Repeat, out.Key, out.Source)
        self.Dispatch(&out)
    }
}

func (self *Dispatcher) setupListeners() {
    self.listenermap = make(map[string]*listeners.Listener)
    self.cs = listeners.NewCommandStream()

    for k, v := range self.config.Listeners {
        l, ok := listeners.GetListener(v.Module, v.Params)
        if ok {
            clog.Debug("Setting up listener `%s`", k)
            self.listenermap[k] = &l
            self.cs.AddListener(l)
        }
    }

}

func (self *Dispatcher) setupModes() {
    for k, _ := range self.config.Modes {
        println(k)
    }
}

func (self *Dispatcher) setupTargets() {
    self.targetmap = make(map[string]targets.Target)
    for k, v := range self.config.Targets {
        t, ok := targets.GetTarget(v.Module, k, v.Params)
        if ok {
            self.targetmap[k] = t
        }
        println(k)
    }
}



// func Setup(t map[string]setup.Target) {
// func Setup() {
    // targetmap = make(map[string]targets.Target)
    // for key, val := range t {
    //    t, ok := targets.GetTarget(val.Module, key, val.Params)
    //    if ok {
    //        targetmap[key] = t
    //    }
    // }
// }

func (self *Dispatcher) Dispatch(rc *listeners.RemoteCommand) bool {
    clog.Debug("repeat: %2d - key: %s - source: %s", rc.Repeat, rc.Key, rc.Source)
    if rc.Key == "KEY_VOLUMEUP" {
        self.targetmap["myX2000"].SendCommand("VolumeUp")
        clog.Debug("sending VolUP to denon")
    } else if rc.Key == "KEY_OK" {
        self.targetmap["myX2000"].SendCommand("Volume", "50")
        SetMode("bar")
        clog.Debug("sending Vol50 to denon")
    } else if rc.Key == "KEY_PLAY" {
        clog.Debug("Current mode: %s", GetMode())
    }
    return true
}
