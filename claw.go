package main

import "os"
import "os/signal"

import "github.com/cnf/go-claw/config"
import "github.com/cnf/go-claw/dispatcher"
import "github.com/cnf/go-claw/clog"

func main() {
    defer clog.Stop()

    // Handle os signals
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func() {
        <- sigc
        clog.Stop()
        os.Exit(1)
    }()

    // Setup configuration
    cfg := &config.Config{}
    cfg.Setup()
    cfg.ReadConfig()

    // Setup claw logger
    clog.SetFlags(clog.Lshortlevel | clog.Ltimebetween | clog.Ltime)
    if cfg.Verbose() {
        clog.SetLogLevel(clog.DEBUG)
    } else {
        clog.SetLogLevel(clog.WARN)
    }

    startHTTPListener(cfg.HTTPPort())

    // Register all listener / target modules
    registerAllListeners()
    registerAllTargets()

    // Start dispatching
    dispatch := dispatcher.Dispatcher{}
    dispatch.Setup(cfg)

    dispatch.Start()
}
