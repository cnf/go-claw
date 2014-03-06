package main

import "flag"
import "io/ioutil"
import "os/user"
import "fmt"
import "runtime"
import "path/filepath"
import "github.com/cnf/go-claw/clog"
import "encoding/json"

type Config struct {
    cfgfile string
    Home string
    System System
}

var cfg Config

type System struct {
    Modes map[string]map[string][]string
    Targets map[string]map[string]string
    Listeners map[string]map[string]string
}

var defaultcfgfile string
// var cfgfile string
// var cfg string
var Verbose bool

func init() {
    cfg = Config{}
    usr, err := user.Current()
    if err != nil {
        clog.Fatal( err.Error() )
    }

    if runtime.GOOS == "windows" {
        cfg.Home, _ = filepath.Abs(usr.HomeDir)
        // different default for windows
    } else {
        cfg.Home, _ = filepath.Abs(usr.HomeDir)
        defaultcfgfile = fmt.Sprintf("%s/.config/claw/config.json", usr.HomeDir)
    }
}


func (self *Config) Setup() {
    tmp := flag.String("conf", defaultcfgfile, "path to our config file.")
    flag.BoolVar(&Verbose, "v", false, "turn on verbose logging")
    flag.Parse()

    self.cfgfile, _ = filepath.Abs(*tmp)
}

func (self *Config) ReadConfigfile() {
    file, ferr := ioutil.ReadFile(defaultcfgfile)
    if ferr != nil {
        // OOPS!
    }
    err := json.Unmarshal(file, self.System)
    if err != nil {
        //OOPS!
    }
}
