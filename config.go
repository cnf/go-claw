package main

import "os"
import "flag"
import "io/ioutil"
import "os/user"
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
    // Modes map[string]map[string][]string `json:"mode"`
    Listeners map[string]Listener
    Modes map[string]Mode
    Targets map[string]Target
}

type Listener struct {
    Module string
    Params map[string]string
}
type Mode map[string]Actionlist
type Actionlist []string
type Target struct {
    Module string
    Params map[string]string
}

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
        cfg.cfgfile = filepath.Join(usr.HomeDir, ".config/claw/config.json")
    }
}


func (self *Config) Setup() {
    flag.StringVar(&self.cfgfile, "conf", self.cfgfile, "path to our config file.")
    flag.BoolVar(&Verbose, "v", Verbose, "turn on verbose logging")
    flag.Parse()

    self.cfgfile, _ = filepath.Abs(self.cfgfile)
}

func (self *Config) ReadConfigfile() {
    clog.Debug("Reading config file: %s", self.cfgfile)
    file, ferr := ioutil.ReadFile(self.cfgfile)
    if ferr != nil {
        clog.Error("Failed to open file: %s", ferr.Error())
        clog.Stop()
        os.Exit(1)
    }
    err := json.Unmarshal(file, &self.System)
    if err != nil {
        clog.Error("Failed to parse json data: %s", err.Error())
        clog.Stop()
        os.Exit(1)
    }
}
