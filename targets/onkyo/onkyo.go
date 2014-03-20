package onkyo

import "strings"
import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"
//import "github.com/tarm/goserial"

type Transport int
const (
    TRANSPORT_TCP    Transport = iota
    TRANSPORT_SERIAL Transport = iota
)

type OnkyoReceiver struct {
    modelname string
    name string
    transport Transport

    serialdev string
    host string
    port int
    model string
    dest_area string // DX: North America; XX: Europe/Asia; JJ: Japan
    identifier string
}


func Register() {
    targets.RegisterTarget("onkyo", CreateOnkyoReceiver)
    //targets.RegisterAutoDetect(OnkyoAutoDetect)
}

func (r *OnkyoReceiver) processparams(params map[string]string) bool {
    if params["connection"] == "serial" {
        r.transport = TRANSPORT_SERIAL
    } else {
        // By default assume TCP
        r.transport = TRANSPORT_TCP
    }

    switch r.transport {
    case TRANSPORT_SERIAL:
        if _, ok := params["device"]; !ok {
            return false
        }
        r.serialdev = params["device"]
        // Baudrate is fixed: 9600
    case TRANSPORT_TCP:
        if _, ok := params["host"]; !ok {
            // No host specified - attempt auto discovery
            if _, ok := params["devname"]; !ok {
                clog.Error("No 'host' or 'devname' parameter specified")
                return false
            }
            r.name = params["devname"]
        } else if _, ok := params["port"]; ok {
            // Host and port specified
        } else if strings.Contains(params["ip"], ":") {
            // Port specified in IP string
        } else {
            return false
        }
    }
    return true
}

func CreateOnkyoReceiver(name string, params map[string]string) (t targets.Target, ok bool) {
    clog.Debug("Creating Onkyo Receiver %s", name)
    var ret OnkyoReceiver

    // Process incoming parameters

    t = ret
    ok = true
    return
}

func (o OnkyoReceiver) SendCommand(cmd string, args ...string) bool {
    clog.Debug("Sending command: %s (%v)", cmd, args)
    return false
}

