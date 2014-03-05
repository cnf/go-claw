package main

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/commandstream"
import "github.com/cnf/go-claw/dispatcher"
import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/setup"
import "os"
import "os/signal"

// import "os"

func main() {
    defer clog.Stop()
    setup.Clargs()
    clog.SetLogLevel(clog.DEBUG)

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func() {
        <- sigc
        clog.Stop()
        os.Exit(1)
    }()

    cs := commandstream.NewCommandStream()
    defer cs.Close()
    var out commandstream.RemoteCommand

    cs.AddListener(&listeners.LircSocketListener{Path: "/var/run/lirc/lircd"})
    cs.AddListener(&listeners.LircSocketListener{Path: "/tmp/echo.sock"})

    dispatcher.Setup()

    for cs.Next(&out) {
        if cs.HasError() {
            clog.Warn("An error occured somewhere: %v", cs.GetError())
            cs.ClearError()
        }
        dispatcher.Dispatch(&out)
        clog.Debug("repeat: %2d - key: %s - source: %s", out.Repeat, out.Key, out.Source)
    }
}
