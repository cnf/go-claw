package main

import "os"
import "os/user"
import "path/filepath"
import "runtime"
import "flag"

import "github.com/cnf/go-claw/config"
import "github.com/cnf/go-claw/clog"

func setup(cfg *config.Config) {
    var cfgfile string
    var verbose bool
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

    cfg.SetVerbose(verbose)
    cfg.SetCfgFile(cfgfile)
}
