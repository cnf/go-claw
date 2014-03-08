package dispatcher

// import "os"
import "io/ioutil"
import "encoding/json"

import "github.com/cnf/go-claw/clog"

import "github.com/kr/pretty"

func (self *Dispatcher) readConfig() {
    clog.Debug("Reading config file: %s", self.Configfile)
    file, ferr := ioutil.ReadFile(self.Configfile)
    if ferr != nil {
        clog.Error("Failed to open file: %s", ferr.Error())
        // clog.Stop()
        // os.Exit(1)
    }

    err := json.Unmarshal(file, &self.config)
    if err != nil {
        clog.Error("Failed to parse json data: %s", err.Error())
    }
    pretty.Print(self.config)
}
