package main

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/dispatcher"
import "github.com/cnf/go-claw/clog"

// import "os"

func main() {
    defer clog.Stop()
    clog.Debug("debug running")
    clog.Info("running")
    clog.SetLogLevel(clog.DEBUG)
    // clog.Setup(&clog.Config{Path: "/tmp/clog.log"})
    // clog.Setup(os.Stdout)
    // cs := &dispatcher.CommandStream{Ch: make(chan *dispatcher.RemoteCommand)}
    cs := dispatcher.NewCommandStream()
    defer cs.Close()
    var out dispatcher.RemoteCommand

    cs.AddListener(&listeners.LircSocketListener{Path: "/var/run/lirc/lircd"})
    cs.AddListener(&listeners.LircSocketListener{Path: "/tmp/echo.sock"})

    for cs.Next(&out) {
        if cs.HasError() {
            clog.Warn("An error occured somewhere: %v", cs.GetError())
            cs.ClearError()
        }
        clog.Debug("repeat: %2d - key: %s - source: %s", out.Repeat, out.Key, out.Source)
    }
}
