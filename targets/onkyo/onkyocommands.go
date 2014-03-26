package onkyo

import "github.com/cnf/go-claw/clog"

func (r *OnkyoReceiver) onkyoCommand(cmd string, args []string) error {
    var rv string
    var err error
    switch cmd {
    case "PowerOn":
        rv, err = r.sendCmd("PWR01")
    case "PowerOff":
        rv, err = r.sendCmd("PWR00")
    case "TogglePower":
        rv, err = r.sendCmd("PWRQSTN")
        if rv == "PWR00" {
            r.sendCmd("PWR01")
        } else {
            r.sendCmd("PWR00")
        }
    case "MuteOn":
        rv, err = r.sendCmd("AMT01")
    case "MuteOff":
        rv, err = r.sendCmd("AMT00")
    case "Mute":
        rv, err = r.sendCmd("AMTTG")
    case "VolumeUp":
        rv, err = r.sendCmd("MVLUP")
    case "VolumeDown":
        rv, err = r.sendCmd("MVLDOWN")
    }
    clog.Debug("Onkyo returned: '%s'", rv)
    return err
}
