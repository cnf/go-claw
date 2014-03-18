package denon

import "net"
import "fmt"
import "time"
import "errors"
import "strings"

import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

type Denon struct {
    name string
    addr *net.TCPAddr
    commands map[string]Commander
    last time.Time
    wait time.Duration
}

func Register() {
    targets.RegisterTarget("denon", Create)
}

func Create(name string, params map[string]string) (t targets.Target, ok bool) {
    // TODO: VALIDATE PARAMS
    if val, ok := params["address"]; ok {
        d := setup(name, val, 23)
        d.commands = AVRX2000
        d.wait = time.Duration(110 * time.Millisecond)
        return d, true
    }
    return nil, false
}

func setup(name string, host string, port int) *Denon {
    tmp, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
    if err != nil {
        clog.Error("Could not Initialize Denon: %s", err)
        return nil
    }
    return &Denon{addr: tmp, name: name}

}

func (d *Denon) SendCommand(cmd string, args ...string) bool {
    switch cmd {
    case "PowerOn":
        return d.powerOn()
    case "Mute":
        return d.toggleMute()
    default:
        cstr, err := d.getCommand(cmd, args...)
        if err != nil { return false }
        _, serr := d.socketSend(cstr)
        if serr != nil { return false }
        return true
    }
    return false
}

func (d *Denon) Capabilities() []string {
    return []string{}
}

func (d *Denon) getCommand(cmd string, args ...string) (string, error) {
    if val, ok := d.commands[cmd]; ok {
        cstr, err := val.Command(args...)
        if err != nil {
            return "", err
        }
        return cstr, nil
    }
    return "", errors.New("could not get command")
}

func (d *Denon) socketSend(str string) (cmd string, err error) {
    if d.addr == nil {
        clog.Warn("No address to sent Denon command to.")
        return "", errors.New("no address set")
    }

    tdiff := time.Since(d.last)
    if tdiff < d.wait {
        // time.Sleep(d.wait)
        clog.Debug("Denon: Waiting %s", (d.wait - tdiff).String())
        time.Sleep(d.wait - tdiff)
    }

    conn, err := net.DialTCP("tcp", nil, d.addr)
    if err != nil {
        clog.Error("Connection failed: %s", err)
        if conn != nil {
            conn.Close()
        }
        return "", err
    }
    conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
    clog.Debug("Sending %s to %s", str, d.name)
    fmt.Fprintf(conn, "%s\r", str)
    reply := make([]byte, 32)
    l, err := conn.Read(reply)
    conn.Close()
    d.last = time.Now()
    if err != nil { return "", err }
    return string(reply[0:l]), nil
}

func (d *Denon) toggleMute() bool {
    r, err := d.socketSend("MU?")
    if err != nil { return false }
    r = strings.TrimSpace(r)
    if r == "MUOFF" {
        cmd, err := d.getCommand("MuteOn")
        if err != nil { return false }
        _, serr := d.socketSend(cmd)
        if serr != nil { return false }
    } else if r == "MUON" {
        cmd, err := d.getCommand("MuteOff")
        if err != nil { return false }
        _, serr := d.socketSend(cmd)
        if serr != nil { return false }
        // _, serr := d.socketSend("MUOFF")
    }
    return true
}

func (d *Denon) powerOn() bool {
    pstr, err := d.getCommand("PowerOn")
    if err != nil { return false }
    rtrn, serr := d.socketSend(pstr)
    if serr != nil { return false }
    rtrn = strings.TrimSpace(rtrn)
    if rtrn != "PWON" { return false }
    time.Sleep(10 * time.Second)

    return true
}
