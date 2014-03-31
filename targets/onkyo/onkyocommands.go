package onkyo

import "github.com/cnf/go-claw/clog"

func (r *OnkyoReceiver) onkyoCommand(cmd string, args []string) error {
    var rv string
    var err error
    switch cmd {
    case "PowerOn":
        rv, err = r.sendCmd("PWR01", -1)
    case "PowerOff":
        rv, err = r.sendCmd("PWR00", -1)
    case "TogglePower":
        rv, err = r.sendCmd("PWRQSTN", -1)
        if err != nil {
            clog.Error("ERROR: %s", err.Error())
            return err
        }
        clog.Debug("Power state query: '%s', %d", rv, len(rv))
        if rv == "PWR00" {
            clog.Debug("Sending PWR01")
            r.sendCmd("PWR01", -1)
        } else {
            clog.Debug("Sending PWR00")
            r.sendCmd("PWR00", -1)
        }
    case "MuteOn":
        rv, err = r.sendCmd("AMT01",0)
    case "MuteOff":
        rv, err = r.sendCmd("AMT00",0)
    case "Mute":
        rv, err = r.sendCmd("AMTTG",0)
    case "VolumeUp":
        rv, err = r.sendCmd("MVLUP",0)
    case "VolumeDown":
        rv, err = r.sendCmd("MVLDOWN",0)
    }
    clog.Debug("Onkyo returned: '%s'", rv)
    return err
}
