package config

// import "os"
import "io/ioutil"
import "encoding/json"

import "github.com/cnf/go-claw/clog"

// ReadConfig reads the config file
func (c *Config) ReadConfig() {
    clog.Info("Reading config file: %s", c.cfgfile)
    file, ferr := ioutil.ReadFile(c.cfgfile)
    if ferr != nil {
        clog.Error("Failed to open file: %s", ferr.Error())
        // clog.Stop()
        // os.Exit(1)
    }

    err := json.Unmarshal(file, &c)
    if err != nil {
        clog.Error("Failed to parse json data: %s", err.Error())
    }
}
