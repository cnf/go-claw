package main

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/dispatcher"
import "github.com/cnf/go-claw/clog"

func main() {
    clog.Info("running")
    // cs := &dispatcher.CommandStream{Ch: make(chan *dispatcher.RemoteCommand)}
    cs := dispatcher.NewCommandStream()
    defer cs.Close()
    var out dispatcher.RemoteCommand

    cs.AddListener(&listeners.LircSocketListener{Path: "/var/run/lirc/lircd"})
    cs.AddListener(&listeners.LircSocketListener{Path: "/tmp/echo.sock"})

    for cs.Next(&out) {
        if cs.HasError() {
            clog.Warn("An error occured somewhere: %v\n", cs.GetError())
            cs.ClearError()
        }
        clog.Info("code: %s - repeat: %3d - key: %s - name: %s\n", out.Code, out.Repeat, out.Key, out.Source)
    }
}
