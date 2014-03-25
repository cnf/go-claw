package onkyo

import "github.com/cnf/go-claw/clog"

func (r *OnkyoReceiver) onkyoCommand(cmd string, args []string) bool {
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
