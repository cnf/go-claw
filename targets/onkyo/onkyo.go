package onkyo

//import "strings"
import "net"
import "errors"
import "time"
import "sync"
//import "encoding/hex"
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
    Name string
    Transport Transport

    Serialdev string
    Host string
    AutoDetect bool
    Model string
    Identifier string

    con net.Conn
    mu sync.Mutex
    lastsend time.Time
}


// Register registers the Onkyo Module in the target manager
func Register() {
    targets.RegisterTarget("onkyo", createOnkyoReceiver)
    //targets.RegisterAutoDetect(OnkyoAutoDetect)
}

/*
params["model"]       = "TX-NR509"
params["id"]         = "CAFFEE"
params["connection"] = "serial|tcp*"
params["device"]     = "/dev/ttyS1"
params["host"]       = "192.168.0.1:123"
*/

func (r *OnkyoReceiver) doConnect() bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    if (r.Transport == TransportSerial) {
        clog.Error("onkyo serial connection is not implemented!")
        return false
    }
    if (r.con != nil) {
        return true
    }
    var autodetected = false
    for {
        if (r.Host == "") && (r.AutoDetect) {
            clog.Debug("Autodetecting Onkyo receiver: %s (%s)", r.Model, r.Identifier)
            if t := OnkyoFind(r.Model, r.Identifier, 3000); t != nil {
                r.Host = t.Detected["host"]
                autodetected = true
            }
        }
        if r.Host == "" {
            clog.Error("No host setting found!")
            return false
        }
        var err error
        clog.Debug("Connecting to %s", r.Host)
        r.con, err = net.DialTimeout("tcp", r.Host, time.Duration(5000) * time.Millisecond)
        if err != nil {
            clog.Error("error connecting to Onkyo Receiver: %s", err.Error());
            if r.con != nil {
                // Should not happen?
                r.con.Close()
                r.con = nil
            }
            if autodetected {
                // Already tried to autodetect, but failed?
                break
            } else if r.AutoDetect {
                // Retry autodetection
                r.Host = ""
                continue
            }
        } else {
            return true
        }
    }
    return false
}

func (r *OnkyoReceiver) processparams(pname string, params map[string]string) error {
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
            return errors.New("no 'device' parameter specified for serial Onkyo receiver")
        }
        r.Serialdev = params["device"]
        if _, ok := params["type"]; !ok {
            return errors.New("no 'type' parameter specified for serial Onkyo receiver")
        }
        // Baudrate is fixed: 9600
    case TransportTCP:
        if _, ok := params["host"]; !ok {
            // No host specified - attempt auto discovery
            var ok bool
            if r.Model, ok = params["model"]; !ok {
                return errors.New("no 'host' or 'type' parameter specified for TCP Onkyo receiver")
            }
            r.AutoDetect = true
            if r.Identifier, ok = params["id"]; !ok {
                clog.Warn("no 'id' specified for onkyo type %s", params["type"])
            }
            if t := OnkyoFind(r.Model, r.Identifier, 3000); t != nil {
                clog.Debug("Found OnkyoReceiver: %v", t)
                r.Host = t.Detected["host"]
            } else {
                // This is not an error? Try again later
                clog.Warn("could not find Onkyo receiver model '%s' id '%s'", r.Model, r.Identifier)
                r.Host = ""
            }
        } else {
            // Test if the host is correct
            _, _, err := net.SplitHostPort(params["host"])
            if (err != nil) {
                return errors.New("not a valid host:port notation in host parameter")
            }
            r.AutoDetect = false
            r.Host = params["host"]
        }
    }
    return nil
}

func (r *OnkyoReceiver) sendCmd(cmd string) (string, bool) {
    errcnt := 0
    buf := make([]byte, 1024)
    for {
        if (errcnt >= 3) {
            return "", false
        }
        if !r.doConnect() {
            clog.Error("Connect failed, giving up")
            return "", false
        }
        switch r.Transport {
        case TransportTCP:
            // Prevent sending a next command within 50ms
            tdiff := time.Since(r.lastsend)
            if tdiff < (time.Duration(50) * time.Millisecond) {
                time.Sleep((time.Duration(50) * time.Millisecond) - tdiff)
            }
            //clog.Debug("Sending command to Onkyo: %s", cmd)
            r.con.SetWriteDeadline(time.Now().Add(time.Duration(300) * time.Millisecond))
            _, err := r.con.Write(NewOnkyoFrameTCP(cmd).Bytes())
            r.lastsend = time.Now()
            if (err != nil) {
                // check error type
                if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
                    r.con.Close()
                    r.con = nil
                } else if (errcnt == 1) {
                    // Second retry that failed - reconnect
                    r.con.Close()
                    r.con = nil
                }
                errcnt++
                continue;
            }
            // Read the response
            r.con.SetReadDeadline(time.Now().Add(time.Duration(100) * time.Millisecond))
            rlen, err := r.con.Read(buf)
            if err != nil {
                nerr, ok := err.(net.Error)
                if ok && nerr.Timeout() {
                    clog.Warn("timeout receiving Onkyo response for command '%s'", cmd)
                } else if ok && nerr.Temporary() {
                    clog.Warn("temporary error reading Onkyo command response: %s", err.Error())
                } else {
                    clog.Error("error reading Onkyo command response: %s", err.Error())
                    r.con.Close()
                    r.con = nil
                }
                errcnt++
                continue
            }
            rcmd := &OnkyoFrameTCP{}
            if err := rcmd.Parse(buf[0:rlen]); err != nil {
                clog.Error("Could not parse Onkyo response: %s", err.Error())
                return "", false
            }
            r.lastsend = time.Now()
            return rcmd.Message(), true
        case TransportSerial:
            return "", false
        }
        break
    }
    return "", false
}

func createOnkyoReceiver(name string, params map[string]string) (targets.Target, bool) {
    clog.Debug("Creating Onkyo Receiver %s", name)
    var ret OnkyoReceiver

    // Process incoming parameters
    if err := ret.processparams(name, params); err != nil {
        clog.Error(err.Error())
        return nil, false
    }
    // 5 seconds in the past
    ret.lastsend = time.Now().Add(time.Duration(-5) * time.Second)
    if !ret.doConnect() {
        clog.Error("could not connect to Onkyo Reciever!")
    }
    return &ret, true
}

// SendCommand sends a command to the receiver
func (r *OnkyoReceiver) SendCommand(cmd string, args ...string) bool {
    clog.Debug("Sending command: %s (%v)", cmd, args)
    // Look up command
    var rv string
    var ok bool
    switch cmd {
    case "PowerOn":
        rv, ok = r.sendCmd("PWR01")
    case "PowerOff":
        rv, ok = r.sendCmd("PWR00")
    case "TogglePower":
        rv, ok = r.sendCmd("PWRQSTN")
        if rv == "PWR00" {
            r.sendCmd("PWR01")
        } else {
            r.sendCmd("PWR00")
        }
    case "MuteOn":
        rv, ok = r.sendCmd("AMT01")
    case "MuteOff":
        rv, ok = r.sendCmd("AMT00")
    case "Mute":
        rv, ok = r.sendCmd("AMTTG")
    case "VolumeUp":
        rv, ok = r.sendCmd("MVLUP")
    case "VolumeDown":
        rv, ok = r.sendCmd("MVLDOWN")
    }
    clog.Debug("Onkyo returned: '%s'", rv)
    return ok
}

