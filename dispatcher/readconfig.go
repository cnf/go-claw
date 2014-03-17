package dispatcher

// import "os"
import "io/ioutil"
import "encoding/json"

import "github.com/cnf/go-claw/clog"

func (d *Dispatcher) readConfig() {
    clog.Info("Reading config file: %s", d.Configfile)
    file, ferr := ioutil.ReadFile(d.Configfile)
    if ferr != nil {
        clog.Error("Failed to open file: %s", ferr.Error())
        // clog.Stop()
        // os.Exit(1)
    }

    err := json.Unmarshal(file, &d.config)
    if err != nil {
        clog.Error("Failed to parse json data: %s", err.Error())
    }
}
