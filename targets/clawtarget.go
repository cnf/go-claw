package targets

import "fmt"
import "strings"

import "github.com/cnf/go-claw/clog"
//import "github.com/cnf/go-claw/modes"

type clawTarget struct {
    targetmanager *TargetManager
    modeactive string
}

// RegisterTarget("modes", createModes)
func (t *clawTarget) Commands() map[string]*Command {
    if t.targetmanager == nil {
        return nil
    }
    cmds := make(map[string]*Command)


    // Add the mode command
    modelist := make([]string, len(t.targetmanager.modes.ModeMap))
    i := 0
    for m := range t.targetmanager.modes.ModeMap {
        modelist[i] = m
        i++
    }
    cmds["mode"] = NewCommand("Selects a mode", 
                       NewParameter("mode", "the mode to select").SetList(strings.Join(modelist, "|")),
                   )
    // Add other internal modes
    return cmds
}

func (t *clawTarget) Stop() error {
    return nil
}

func (t *clawTarget) setMode(cmd string, args ...string) error {
    newmode := args[0]

    if (t.modeactive != "") {
        return fmt.Errorf("aborted: attempting to recursively set mode '%s' while still setting mode '%s'", cmd, t.modeactive)
    }
    t.modeactive = newmode
    defer func() { t.modeactive = "" }()

    clog.Debug("Setting mode to: '%s'", newmode)
    str, err := t.targetmanager.modes.SetActive(newmode)
    if err != nil {
        return err
    }
    var ret error
    ret = nil
    for i := 0; i < len(str); i++ {
        err := t.targetmanager.RunCommand(str[i])
        if err != nil {
            clog.Error("Command error while switching to mode '%s': %s", newmode, err.Error())
            // Return last error?
            ret = err
        }
    }

    return ret
}

func (t *clawTarget) SendCommand(cmd string, args ...string) error {
    switch(cmd) {
    case "mode":
        return t.setMode(cmd, args...)
    default:
        return fmt.Errorf("clawtarget does not have a command %s", cmd)
    }
}

func createClawTarget(name string, params map[string]string) (Target, error) {
    ret := &clawTarget{targetmanager: nil, modeactive: ""}
    return ret, nil
}

func (t *clawTarget) setTargetManager(tm *TargetManager) {
    t.targetmanager = tm
}

