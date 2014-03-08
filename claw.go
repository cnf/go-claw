package main

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/dispatcher"
import "github.com/cnf/go-claw/setup"
import "github.com/cnf/go-claw/clog"
import "os"
import "os/signal"

// import "os"

func main() {
    defer clog.Stop()

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func() {
        <- sigc
        clog.Stop()
        os.Exit(1)
    }()

    cfg := setup.Setup()
    cfg.ReadConfigfile()

    if setup.Verbose {
        clog.SetLogLevel(clog.DEBUG)
    } else {
        clog.SetLogLevel(clog.WARN)
    }
    RegisterAllListeners()
    RegisterAllTargets()
    // Parse command line

    cs := listeners.NewCommandStream()
    defer cs.Close()
    var out listeners.RemoteCommand

    for key, val := range cfg.System.Listeners {
        l, ok := listeners.GetListener(val.Module, val.Params)
        if ok {
            clog.Debug("Setting up listener `%s`", key)
            cs.AddListener(l)
        }
    }

    dispatcher.Setup(cfg.System.Targets)

    for cs.Next(&out) {
        if cs.HasError() {
            clog.Warn("An error occured somewhere: %v", cs.GetError())
            cs.ClearError()
        }
        dispatcher.Dispatch(&out)
        // clog.Debug("repeat: %2d - key: %s - source: %s", out.Repeat, out.Key, out.Source)
    }
}
