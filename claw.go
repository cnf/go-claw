package main

import "os"
import "os/signal"
import "os/user"
import "path/filepath"
import "runtime"
import "flag"

// import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/dispatcher"
// import "github.com/cnf/go-claw/setup"
import "github.com/cnf/go-claw/clog"

var cfgfile string
var verbose bool

func main() {
    defer clog.Stop()

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func() {
        <- sigc
        clog.Stop()
        os.Exit(1)
    }()

    setup()
    if verbose {
        clog.SetLogLevel(clog.DEBUG)
    } else {
        clog.SetLogLevel(clog.WARN)
    }


    registerAllListeners()
    registerAllTargets()

    dispatch := dispatcher.Dispatcher{}
    dispatch.Configfile = cfgfile

    dispatch.Start()
}

func setup() {
    var home string
    usr, uerr := user.Current()
    if uerr != nil {
        clog.Warn( uerr.Error() )
        home = os.ExpandEnv("$HOME")
    } else {
        home = usr.HomeDir
    }

    if runtime.GOOS == "windows" {
        // cfg.Home, _ = filepath.Abs(usr.HomeDir)
        // TODO: Different defaults for windows
    } else {
        // cfg.Home, _ = filepath.Abs(usr.HomeDir)
        cfgfile = filepath.Join(home, ".config/claw/config.json")
    }
    flag.StringVar(&cfgfile, "conf", cfgfile, "path to our config file.")
    flag.BoolVar(&verbose, "v", verbose, "turn on verbose logging")
    flag.Parse()
    cfgfile, _ = filepath.Abs(cfgfile)
}
