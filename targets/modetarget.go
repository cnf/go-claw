package targets

import "fmt"

import "github.com/cnf/go-claw/clog"
//import "github.com/cnf/go-claw/modes"

type ModeTarget struct {
    targetmanager *TargetManager
    modeactive string
}

// RegisterTarget("modes", createModes)
func (m *ModeTarget) Commands() map[string]*Command {
    if m.targetmanager == nil {
        return nil
    }
    cmds := make(map[string]*Command)
    for m := range(m.targetmanager.modes.ModeMap) {
        clog.Debug("ModeTarget: Adding 'mode::%s'", m)
        cmds[m] = NewCommand(m, "Selects the mode '" + m + "'")
    }
    /* For future reference - add "set" command?
    cmds["set"] = NewCommand("set", "select a mode",
                    NewParameter("mode", "the mode to select", false).SetString(),
                )
    */
    return cmds
}

func (t *ModeTarget) Stop() error {
    return nil
}

func (t *ModeTarget) SendCommand(cmd string, args ...string) error {
    // If we get here, the 'cmd' should have been validated
    if (t.modeactive != "") {
        return fmt.Errorf("aborted: attempting to recursively set mode '%s' while still setting mode '%s'", cmd, t.modeactive)
    }
    t.modeactive = cmd
    defer func() { t.modeactive = "" }()

    clog.Debug("Setting mode to: '%s'", cmd)
    str, err := t.targetmanager.modes.SetActive(cmd)
    if err != nil {
        return err
    }
    var ret error
    ret = nil
    for i := 0; i < len(str); i++ {
        err := t.targetmanager.RunCommand(str[i])
        if err != nil {
            clog.Error("Command error while switching to mode '%s': %s", cmd, err.Error())
            // Return last error?
            ret = err
        }
    }

    return ret
}

func createModeHandler(name string, params map[string]string) (Target, error) {
    ret := &ModeTarget{targetmanager: nil, modeactive: ""}
    return ret, nil
}

func (t *ModeTarget) setTargetManager(tm *TargetManager) {
    t.targetmanager = tm
}

