package onkyo

import "fmt"
import "strconv"
import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

func (d *OnkyoReceiver) Commands() map[string]*targets.Command {
    cmds := map[string]*targets.Command {
        "Power"       : targets.NewCommand("Powers off the receiver",
                targets.NewParameter("powerstate", "The power state").SetList("on", "off", "toggle"),
                ),
        "Mute"        : targets.NewCommand("Controls the Mute state",
                targets.NewParameter("mutestate", "The mute state").SetList("on", "off", "toggle"),
                ),
        "VolumeUp"    : targets.NewCommand("Turns up the volume"),
        "VolumeDown"  : targets.NewCommand("Turns down the volume"),
        "Volume"      : targets.NewCommand("Sets the volume",
                targets.NewParameter("volumelevel", "The volume level").SetRange(0, 77),
                ),
        "Input"       : targets.NewCommand("selects an input",
                targets.NewParameter("input", "the input to select").SetList("test|test2"),
                ),
        "InputRaw"    : targets.NewCommand("Selects raw input number",
                targets.NewParameter("input", "the input number").SetNumeric(),
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
        rv, err = r.sendCmd("AMT01", 0)
    case "off":
        rv, err = r.sendCmd("AMT00", 0)
    case "toggle":
        rv, err = r.sendCmd("AMTTG", 0)
    }
    return rv, err
}

func (r *OnkyoReceiver) Power(state string) (string, error) {
    var rv string
    var err error

    switch state {
    case "on":
        rv, err = r.sendCmd("PWR01", -1)
    case "off":
        rv, err = r.sendCmd("PWR00", -1)
    case "toggle":
        rv, err = r.sendCmd("PWRQSTN", -1)
        if err != nil {
            clog.Error("ERROR: %s", err.Error())
            return "", err
        }
        clog.Debug("Power state query: '%s', %d", rv, len(rv))
        if rv == "PWR00" {
            clog.Debug("Sending PWR01")
            r.sendCmd("PWR01", -1)
        } else {
            clog.Debug("Sending PWR00")
            r.sendCmd("PWR00", -1)
        }
    }
    return rv, err
}

func (r *OnkyoReceiver) setInput(input string) error {
    _, err := r.sendCmd(fmt.Sprintf("SLI%s", input), 0)
    return err
}

func (r *OnkyoReceiver) onkyoCommand(cmd string, args []string) error {
    var err error
    switch cmd {
    case "Power":
        _, err = r.Power(args[0])
    case "Mute":
        _, err = r.Mute(args[0])
    case "VolumeUp":
        _, err = r.sendCmd("MVLUP",0)
    case "VolumeDown":
        _, err = r.sendCmd("MVLDOWN",0)
    case "Volume":
        ml, _ := strconv.Atoi(args[0])
        // TODO: Most models require hex volume level, some require decimal!
        _, err = r.sendCmd(fmt.Sprintf("MVL%02X", ml), 0)
    case "InputRaw":
        ml, _ := strconv.Atoi(args[0])
        _, err = r.sendCmd(fmt.Sprintf("SLI%02X", ml), 0)
    case "Input":
        err = r.setInput(args[0])
    default:
        err = fmt.Errorf("unknown command for onkyo module: '%s'", cmd)
    }
    return err
}
