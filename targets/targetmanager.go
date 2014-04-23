package targets

import "fmt"
import "errors"
import "strings"
import "time"
import "io"
import "bufio"

import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/modes"
import "github.com/cnf/go-claw/core"

// TargetManager is the structure which manages all targets
type TargetManager struct {
    targets map[string]Target
    targetCmds map[string]map[string]*Command
    modes *modes.Modes
}

// NewTargetManager creates and initialize a new TargetManager object
func NewTargetManager(m *modes.Modes) *TargetManager {
    ret := &TargetManager{ targets: nil, targetCmds: nil, modes: m }
    //clog.Debug("Adding internal mode target...")
    //ret.Add("mode", "mode", nil)
    ret.Stop()
    return ret
}

// Add adds a new target using the specified module, given the specified name and parameters
func (t *TargetManager) Add(module, name string, params map[string]string) error {
    // Validate name
    var err error

    // Everything lower case internally
    module = strings.ToLower(module)
    name  = strings.ToLower(name)

    if err := core.ValidateName(name); err != nil {
        return err
    }
    // check if target already exists
    if _, ok := t.targets[name]; ok {
        clog.Warn("TargetManager::Add(): Target name already existed - removing first")
        t.Remove(name)
    }
    // Check if the requested module exists
    if _, ok := targetlist[module]; !ok {
        return errors.New("could not create target '" + name + "': module '" + module + "' is not registered")
    }
    // Create the target instance
    var tgt Target
    tgt, err = targetlist[module](name, params)
    if err != nil {
        clog.Warn("Could not create %s::%s: %s", module, name, err.Error())
        return err
    }
    t.targets[name] = tgt

    // Special case - test if this is the modes target
    if mt, ok := tgt.(*clawTarget); ok {
        mt.setTargetManager(t)
    }

    // Fetch the command list
    tcmdlist := tgt.Commands()
    if tcmdlist == nil {
        clog.Warn("warning: %s::%s returned an empty command list!", module, name)
    } else {
        t.targetCmds[name] = make(map[string]*Command, len(tcmdlist))
        for r := range(tcmdlist) {
            t.targetCmds[name][strings.ToLower(r)] = tcmdlist[r]
            t.targetCmds[name][r].Name = strings.ToLower(r)
        }
    }

    return nil
}

func (t *TargetManager) PrintCommands(w io.Writer, json bool) {
    wt := bufio.NewWriter(w)
    if json { wt.WriteString("{\n") }
    for tgt := range(t.targetCmds) {
        if !json { fmt.Fprintf(wt, "## target: %s\n") }
        for cmd := range(t.targetCmds[tgt]) {
            if json {
                fmt.Fprintf(wt, `"%s::%s": "%s"\n`, tgt, cmd)
            } else {
                fmt.Fprintf(wt, `%s::%s %s\n`, tgt, cmd)
            }
        }
    }
    if json { wt.WriteString("}\n") }
    wt.Flush()
}

// Remove removes a target instance from the list
func (t *TargetManager) Remove(name string) error {
    if _, ok := t.targets[name]; !ok {
        return errors.New("cannot remove " + name + ": does not exist")
    }
    if err := t.targets[name].Stop(); err != nil {
        return err
    }
    delete(t.targets, name)
    if _, ok := t.targetCmds[name]; ok {
        delete(t.targetCmds, name)
    }
    return nil
}

// Starts all target's background processes if needed
func (t *TargetManager) Start() error {
    for tgt := range(t.targets) {
        err := t.targets[tgt].Start()
        if err != nil {
            return err
        }
    }
    return nil
}

// Stop stops all target instances and removes them
func (t *TargetManager) Stop() error {
    for k := range t.targets {
        if err := t.Remove(k); err != nil {
            return err
        }
    }
    t.targets    = make(map[string]Target)
    t.targetCmds = make(map[string]map[string]*Command)
    clog.Debug("TargetManager::Stop(): Adding internal claw target...")
    t.Add("claw", "claw", nil)

    return nil
}

// RunCommand parses a given command, determines which target should run it,
// checks the provided parameters, and if all is good - run the command.
func (t *TargetManager) RunCommand(repeated int, cmdstring string) error {
    splitstr := strings.SplitN(cmdstring, "::", 2)
    tstart := time.Now()
    if len(splitstr) != 2 {
        return fmt.Errorf("invalid command string '%s', expected it to contain '::'", cmdstring)
    }
    // Validate the name of the target we just parsed out
    if err := core.ValidateName(splitstr[0]); err != nil {
        return err
    }
    tgtname := strings.ToLower(splitstr[0])

    if _, ok := t.targets[tgtname]; !ok {
        //return fmt.Errorf("command '%s' uses a target '%s' that does not exist", cmdstring, tgtname)
        return NewCommandError(tgtname, false, splitstr[1], false, nil)
    }

    // Split the command
    splitcmd := core.SplitQuoted(splitstr[1])
    if len(splitcmd)  == 0 {
        //return fmt.Errorf("empty target command in '%s'", cmdstring)
        return NewCommandError(tgtname, true, splitstr[1], false, nil)
    }
    tcommand := strings.ToLower(splitcmd[0])
    tparams := splitcmd[1:]

    // Check if the instance provided a commands list to check
    if _, ok := t.targetCmds[tgtname]; ok && t.targetCmds[tgtname] != nil {
        // Check if the command exists for this target
        if _, ok := t.targetCmds[tgtname][tcommand]; !ok {
            //return fmt.Errorf("command '%s' not recognized by target '%s'", tcommand, tgtname)
            return NewCommandError(tgtname, true, tcommand, false, tparams)
        }
        // Validate all parameters
        pc := 0
        var tparamsArr []string
        for prm := 0; prm < len(t.targetCmds[tgtname][tcommand].Parameters); prm++ {
            if pc >= len(tparams) {
                // Parameter not present, check if required
                if !t.targetCmds[tgtname][tcommand].Parameters[prm].Optional {
                    clog.Error("Non-optional parameter %s missing for command %s, target %s",
                            t.targetCmds[tgtname][tcommand].Parameters[prm].Name,
                            tcommand,
                            tgtname,
                        )
                    return fmt.Errorf("non-optional parameter '%s' missing for command '%s', target '%s'",
                            t.targetCmds[tgtname][tcommand].Parameters[prm].Name,
                            tcommand,
                            tgtname,
                        )
                }
            } else {
                pval, err := t.targetCmds[tgtname][tcommand].Parameters[prm].Validate(tparams[pc])
                if (err != nil) {
                    // validation returned an error
                    clog.Error("Parameter validation %s failed for command %s, target %s",
                            t.targetCmds[tgtname][tcommand].Parameters[prm].Name,
                            tcommand,
                            tgtname,
                        )
                    return err
                }
                tparamsArr = append(tparamsArr, pval)
            }
            pc++
        }
        // replace the original parameters with the validated parameters
        tparams = tparamsArr
    }
    // Run the command
    //clog.Debug("--> Process cmd '%s' took: %s", cmdstring, time.Since(tstart).String())
    //tstart = time.Now()
    err := t.targets[tgtname].SendCommand(repeated, tcommand, tparams...)
    clog.Debug("--> Execute cmd '%s' took: %s", cmdstring, time.Since(tstart).String())
    return err
}


