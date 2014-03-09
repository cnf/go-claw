package dispatcher

import "strings"

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/clog"

type Dispatcher struct {
    Configfile string
    config Config
    listenermap map[string]*listeners.Listener
    targetmap map[string]targets.Target
    modemap map[string]*Mode
    activemode string
    cs *listeners.CommandStream
}

func (self *Dispatcher) Start() {
    defer self.cs.Close()
    self.activemode = "default"
    self.activemode = "spotify"
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
        // clog.Debug("repeat: %2d - key: %s - source: %s", out.Repeat, out.Key, out.Source)
        self.dispatch(&out)
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
    self.modemap = make(map[string]*Mode)
    for k, v := range self.config.Modes {
        println(k)
        self.modemap[k] = &Mode{Keys: make(map[string][]string)}
        for kk, kv := range v {
            println(kk)
            self.modemap[k].Keys[kk] = make([]string, len(kv))
            i := 0
            for _, av := range kv {
                self.modemap[k].Keys[kk][i] = av
                i++
                println(av)
            }
        }
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

func (self *Dispatcher) dispatch(rc *listeners.RemoteCommand) bool {
    clog.Debug("repeat: %2d - key: %s - source: %s", rc.Repeat, rc.Key, rc.Source)
    var mod string
    var cmd string
    if val, ok := self.modemap[self.activemode].Keys[rc.Key]; ok {
        println("FOUND!")
        mod, cmd = self.resolve(val[0])
    } else if val, ok := self.modemap["default"].Keys[rc.Key]; ok {
        println("FOUND in default!")
        mod, cmd = self.resolve(val[0])
    } else {
        println("Not found")
        return false
    }

    if t, ok := self.targetmap[mod]; ok {
        println(cmd)
        t.SendCommand(cmd)
        return true
    }

    // if rc.Key == "KEY_VOLUMEUP" {
    //     self.targetmap["myX2000"].SendCommand("VolumeUp")
    //     clog.Debug("sending VolUP to denon")
    // } else if rc.Key == "KEY_OK" {
    //     self.targetmap["myX2000"].SendCommand("Volume", "50")
    //     // SetMode("bar")
    //     clog.Debug("sending Vol50 to denon")
    // } else if rc.Key == "KEY_PLAY" {
    //     // clog.Debug("Current mode: %s", GetMode())
    // }
    return true
}

func (self *Dispatcher) resolve(input string) (mod string, cmd string) {
    foo := strings.Split(input, ":")
    println(foo[0], foo[1])
    return foo[0], foo[1]
}
