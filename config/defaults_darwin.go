package config

import "os"
import "os/user"
import "path/filepath"

import "github.com/cnf/go-claw/clog"

func getConfigPath() string {
    var cfgfile string
    var home string
    usr, uerr := user.Current()
    if uerr != nil {
        clog.Warn( uerr.Error() )
        home = os.ExpandEnv("$HOME")
    } else {
        home = usr.HomeDir
    }
    cfgfile = filepath.Join(home, ".config/claw/config.json")
    return cfgfile
}
