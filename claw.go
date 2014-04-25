package main

import "os"
import "os/signal"

import "github.com/cnf/go-claw/config"
import "github.com/cnf/go-claw/dispatcher"
import "github.com/cnf/go-claw/clog"

func main() {
    defer clog.Stop()

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func() {
        <- sigc
        clog.Stop()
        os.Exit(1)
    }()

    cfg := &config.Config{}
    cfg.Setup()
    cfg.ReadConfig()

    clog.SetFlags(clog.Lshortlevel | clog.Ltimebetween | clog.Ltime)
    if cfg.Verbose() {
        clog.SetLogLevel(clog.DEBUG)
    } else {
        clog.SetLogLevel(clog.WARN)
    }


    registerAllListeners()
    registerAllTargets()

    dispatch := dispatcher.Dispatcher{}
    dispatch.Setup(cfg)

    dispatch.Start()
}
