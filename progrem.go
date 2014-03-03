package main

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/commandstream"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/clog"
import "os"
import "os/signal"

// import "os"

func main() {
    defer clog.Stop()
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

    for cs.Next(&out) {
        if cs.HasError() {
            clog.Warn("An error occured somewhere: %v", cs.GetError())
            cs.ClearError()
        }
        commander(&out)
        // clog.Debug("repeat: %2d - key: %s - source: %s", out.Repeat, out.Key, out.Source)
    }
}

func commander(rc *commandstream.RemoteCommand) bool {
    clog.Debug("repeat: %2d - key: %s - source: %s", rc.Repeat, rc.Key, rc.Source)
    if rc.Key == "KEY_VOLUMEUP" {
        targets.Denon("MVUP")
        clog.Debug("sending VolUP to denon")
    } else if rc.Key == "KEY_OK" {
        targets.Denon("MV50")
        clog.Debug("sending Vol50 to denon")
    }
    return true
}
