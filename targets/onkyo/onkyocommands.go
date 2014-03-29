package onkyo

import "fmt"
import "strconv"
import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

func (d *OnkyoReceiver) Commands() map[string]*targets.Command {
    cmds := map[string]*targets.Command {
        "PowerOn"     : targets.NewCommand("PowerOn", "Powers on the receiver"),
        "PowerOff"    : targets.NewCommand("PowerOff", "Powers off the receiver"),
        "PowerToggle" : targets.NewCommand("PowerOff", "Powers off the receiver"),
        "Power"       : targets.NewCommand("PowerOff", "Powers off the receiver",
                targets.NewParameter("powerstate", "The power state", false).SetList("on", "off", "toggle"),
                ),
        "MuteOn"      : targets.NewCommand("MuteOn", "Mutes the sound"),
        "MuteOff"     : targets.NewCommand("MuteOff", "Unmutes the sound"),
        "MuteToggle"  : targets.NewCommand("MuteToggle", "Toggles the muting of the sound"),
        "Mute"        : targets.NewCommand("Mute", "Controls the Mute state",
                targets.NewParameter("mutestate", "The mute state", false).SetList("on", "off", "toggle"),
                ),
        "VolumeUp"    : targets.NewCommand("VolumeUp", "Turns up the volume"),
        "VolumeDown"  : targets.NewCommand("VolumeDown", "Turns down the volume"),
        "Volume"      : targets.NewCommand("Volume", "Sets the volume",
                targets.NewParameter("volumelevel", "The volume level", false).SetRange(0, 77),
                ),
    }
    if (true) {
        // Fix the volume range for specific models
        cmds["Volume"].Parameters[0].SetRange(0,77)
    }
    return cmds
}

func (r *OnkyoReceiver) Mute(state string) (string, error) {
    var rv string
    var err error

    switch state {
    case "on":
        rv, err = r.sendCmd("AMT01")
    case "off":
        rv, err = r.sendCmd("AMT00")
    case "toggle":
        rv, err = r.sendCmd("AMTTG")
    }
    return rv, err
}
func (r *OnkyoReceiver) Power(state string) (string, error) {
    var rv string
    var err error

    switch state {
    case "on":
        rv, err = r.sendCmd("PWR01")
    case "off":
        rv, err = r.sendCmd("PWR00")
    case "toggle":
        rv, err = r.sendCmd("PWRQSTN")
        if err != nil {
            clog.Error("ERROR: %s", err.Error())
            return "", err
        }
        clog.Debug("Power state query: '%s', %d", rv, len(rv))
        if rv == "PWR00" {
            clog.Debug("Sending PWR01")
            r.sendCmd("PWR01")
        } else {
            clog.Debug("Sending PWR00")
            r.sendCmd("PWR00")
        }
    }
    return rv, err
}
func (r *OnkyoReceiver) onkyoCommand(cmd string, args []string) error {
    var rv string
    var err error
    switch cmd {
    case "PowerOn":
        r.Power("on")
    case "PowerOff":
        r.Power("off")
    case "PowerToggle":
        r.Power("toggle")
    case "Power":
        r.Power(args[0])
    case "MuteOn":
        r.Mute("on")
    case "MuteOff":
        r.Mute("off")
    case "MuteToggle":
        r.Mute("toggle")
    case "Mute":
        r.Mute(args[0])
    case "VolumeUp":
        rv, err = r.sendCmd("MVLUP")
    case "VolumeDown":
        rv, err = r.sendCmd("MVLDOWN")
    case "Volume":
        ml, _ := strconv.Atoi(args[0])
        // TODO: Most models require hex volume level, some require decimal!
        rv, err = r.sendCmd(fmt.Sprintf("MVL%02X", ml))
    }
    clog.Debug("Onkyo returned: '%s'", rv)
    return err
}
