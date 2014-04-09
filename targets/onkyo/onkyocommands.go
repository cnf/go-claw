package onkyo

import "fmt"
import "strconv"
import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

// Commands returns the list of accepted onkyo commands to comply with the Target interface
func (o *OnkyoReceiver) Commands() map[string]*targets.Command {
    cmds := map[string]*targets.Command {
        "power"       : targets.NewCommand("Powers off the receiver",
                targets.NewParameter("powerstate", "The power state").SetList("on", "off", "toggle"),
                ),
        "mute"        : targets.NewCommand("Controls the Mute state",
                targets.NewParameter("mutestate", "The mute state").SetList("on", "off", "toggle"),
                ),
        "volumeup"    : targets.NewCommand("Turns up the volume"),
        "volumedown"  : targets.NewCommand("Turns down the volume"),
        "volume"      : targets.NewCommand("Sets the volume",
                targets.NewParameter("volumelevel", "The volume level").SetRange(0, 77),
                ),
        "input"       : targets.NewCommand("selects an input",
                targets.NewParameter("input", "the input to select").SetList("test|test2"),
                ),
        "inputraw"    : targets.NewCommand("Selects raw input number",
                targets.NewParameter("input", "the input number").SetNumeric(),
                ),
    }
    if (true) {
        // Fix the volume range for specific models
        cmds["volume"].Parameters[0].SetRange(0,77)
    }
    return cmds
}

// Mute executes the mute command on a receiver. Accepts "on", "off" and "toggle" as parameters.
func (o *OnkyoReceiver) Mute(state string) (string, error) {
    var rv string
    var err error
    switch state {
    case "on":
        rv, err = o.sendCmd("AMT01", 0)
    case "off":
        rv, err = o.sendCmd("AMT00", 0)
    case "toggle":
        rv, err = o.sendCmd("AMTTG", 0)
    }
    return rv, err
}

// Power controls the power state of the receiver. Accepts "on", "off" and "toggle" as parameters.
func (o *OnkyoReceiver) Power(state string) (string, error) {
    var rv string
    var err error

    switch state {
    case "on":
        rv, err = o.sendCmd("PWR01", -1)
    case "off":
        rv, err = o.sendCmd("PWR00", -1)
    case "toggle":
        rv, err = o.sendCmd("PWRQSTN", -1)
        if err != nil {
            clog.Error("ERROR: %s", err.Error())
            return "", err
        }
        clog.Debug("Power state query: '%s', %d", rv, len(rv))
        if rv == "PWR00" {
            clog.Debug("Sending PWR01")
            o.sendCmd("PWR01", -1)
        } else {
            clog.Debug("Sending PWR00")
            o.sendCmd("PWR00", -1)
        }
    }
    return rv, err
}

// SetInput sets the input to the specified value. The suported inputs are type-specific
func (o *OnkyoReceiver) SetInput(input string) error {
    _, err := o.sendCmd(fmt.Sprintf("SLI%s", input), 0)
    return err
}

func (o *OnkyoReceiver) onkyoCommand(cmd string, args []string) error {
    var err error
    switch cmd {
    case "power":
        _, err = o.Power(args[0])
    case "mute":
        _, err = o.Mute(args[0])
    case "volumeup":
        _, err = o.sendCmd("MVLUP",0)
    case "volumedown":
        _, err = o.sendCmd("MVLDOWN",0)
    case "volume":
        ml, _ := strconv.Atoi(args[0])
        // TODO: Most models require hex volume level, some require decimal!
        _, err = o.sendCmd(fmt.Sprintf("MVL%02X", ml), 0)
    case "inputraw":
        ml, _ := strconv.Atoi(args[0])
        _, err = o.sendCmd(fmt.Sprintf("SLI%02X", ml), 0)
    case "input":
        err = o.SetInput(args[0])
    default:
        err = fmt.Errorf("unknown command for onkyo module: '%s'", cmd)
    }
    return err
}
