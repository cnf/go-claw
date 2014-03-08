package dispatcher

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/clog"

type Dispatcher struct {
    Configfile string
    config Config
    cs *listeners.CommandStream
}

func (self *Dispatcher) Start() {
    println("Starting")
    self.cs = listeners.NewCommandStream()
    defer self.cs.Close()

    self.readConfig()

    for key, val := range self.config.Listeners {
        println(key)
        l, ok := listeners.GetListener(val.Module, val.Params)
        if ok {
            clog.Debug("Setting up listener `%s`", key)
            self.cs.AddListener(l)
        }
    }

    var out listeners.RemoteCommand


    for self.cs.Next(&out) {
        if self.cs.HasError() {
            clog.Warn("An error occured somewhere: %v", self.cs.GetError())
            self.cs.ClearError()
        }
        clog.Debug("repeat: %2d - key: %s - source: %s", out.Repeat, out.Key, out.Source)
    }
}

var targetmap map[string]targets.Target


// func Setup(t map[string]setup.Target) {
func Setup() {
    targetmap = make(map[string]targets.Target)
    // for key, val := range t {
    //    t, ok := targets.GetTarget(val.Module, key, val.Params)
    //    if ok {
    //        targetmap[key] = t
    //    }
    // }
}

func Dispatch(rc *listeners.RemoteCommand) bool {
    clog.Debug("repeat: %2d - key: %s - source: %s", rc.Repeat, rc.Key, rc.Source)
    if rc.Key == "KEY_VOLUMEUP" {
        targetmap["myX2000"].SendCommand("VolumeUp")
        clog.Debug("sending VolUP to denon")
    } else if rc.Key == "KEY_OK" {
        targetmap["myX2000"].SendCommand("Volume", "50")
        SetMode("bar")
        clog.Debug("sending Vol50 to denon")
    } else if rc.Key == "KEY_PLAY" {
        clog.Debug("Current mode: %s", GetMode())
    }
    return true
}
