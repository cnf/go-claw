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
    self.activemode = "plex"
    self.readConfig()
    self.setupListeners()
    self.setupModes()
    self.setupTargets()

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
        self.modemap[k] = &Mode{Keys: make(map[string][]string)}
        for kk, kv := range v {
            self.modemap[k].Keys[kk] = make([]string, len(kv))
            i := 0
            for _, av := range kv {
                self.modemap[k].Keys[kk][i] = av
                i++
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
            println(k)
        }
    }
}

func (self *Dispatcher) dispatch(rc *listeners.RemoteCommand) bool {
    clog.Debug("repeat: %2d - key: %s - source: %s", rc.Repeat, rc.Key, rc.Source)
    var mod string
    var cmd string
    var args string
    var rok bool
    if val, ok := self.modemap[self.activemode].Keys[rc.Key]; ok {
        clog.Debug("+ Found `%s` in %s", rc.Key, self.activemode)
        for _, v := range val {
            clog.Debug(v)
            mod, cmd, args, rok = self.resolve(v)
            self.sender(mod, cmd, args)
        }
        return true
    } else if val, ok := self.modemap["default"].Keys[rc.Key]; ok {
        clog.Debug("+ Found `%s` in default!", rc.Key)
        for _, v := range val {
            mod, cmd, args, rok = self.resolve(v)
            self.sender(mod, cmd, args)
        }
        return true
    } else {
        clog.Debug("+ `%s` Not found.")
        return false
    }
    if !rok {
        return false
    }

    return true
}

func (self *Dispatcher) resolve(input string) (mod string, cmd string, args string, ok bool) {
    clog.Debug("++ Resolving input for %s", input)
    foo := strings.SplitN(input, "::", 2)
    if len(foo) < 2 {
        clog.Warn("%s is not a well formed command", input)
        return "", "", "", false
    }
    bar := strings.SplitN(foo[1], " ", 2)
    baz := ""
    if len(bar) > 1 {
        baz = bar[1]
    }

    return foo[0], bar[0], baz, true
}

func (self *Dispatcher) sender(mod string, cmd string, args string) bool {
    if mod == "mode" {
        clog.Debug("++++ %s - %s", mod, cmd)
        return self.setMode(cmd)
    }
    if t, ok := self.targetmap[mod]; ok {
        sok := t.SendCommand(cmd, args)
        if !sok {
            clog.Debug("- Failed to send command `%s` for `%s`", cmd, mod)
        }
        return true
    }
    return false
}

func (self *Dispatcher) setMode(mode string) bool {
    if _, ok := self.modemap[mode]; ok {
        clog.Debug("+ Mode changed to `%s`", mode)
        self.activemode = mode
    } else {
        for k, _ := range self.modemap {
            clog.Debug("---- %s", k)
        }
        return false
    }
    return true
}
