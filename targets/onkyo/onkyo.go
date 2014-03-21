package onkyo

import "strings"
//import "github.com/tarm/goserial"
import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

// Transport indicates the transport type the Onkyo Reciever uses
type Transport int
const (
    // TransportTCP indicates the useage of TCP
    TransportTCP   Transport = iota
    // TransportSerial indicates the useage of a Serial line
    TransportSerial Transport = iota
)

// OnkyoReceiver structure
type OnkyoReceiver struct {
    Modelname string
    Name string
    Transport Transport

    Serialdev string
    Host string
    Port int
    Model string
    DestArea string // DX: North America; XX: Europe/Asia; JJ: Japan
    Identifier string
}


// Register registers the Onkyo Module in the target manager
func Register() {
    targets.RegisterTarget("onkyo", createOnkyoReceiver)
    //targets.RegisterAutoDetect(OnkyoAutoDetect)
}

func (r OnkyoReceiver) processparams(pname string, params map[string]string) bool {
    if params["connection"] == "serial" {
        r.Transport = TransportSerial
    } else {
        // By default assume TCP
        r.Transport = TransportTCP
    }
    r.Name = pname
    switch r.Transport {
    case TransportSerial:
        if _, ok := params["device"]; !ok {
            return false
        }
        r.Serialdev = params["device"]
        // Baudrate is fixed: 9600
    case TransportTCP:
        if _, ok := params["host"]; !ok {
            // No host specified - attempt auto discovery
            if _, ok := params["devname"]; !ok {
                clog.Error("No 'host' or 'devname' parameter specified")
                return false
            }
            r.Identifier = params["id"]
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

func createOnkyoReceiver(name string, params map[string]string) (t targets.Target, ok bool) {
    clog.Debug("Creating Onkyo Receiver %s", name)
    var ret OnkyoReceiver

    // Process incoming parameters

    t = ret
    ok = true
    return
}

// SendCommand sends a command to the receiver
func (r OnkyoReceiver) SendCommand(cmd string, args ...string) bool {
    clog.Debug("Sending command: %s (%v)", cmd, args)
    return false
}

