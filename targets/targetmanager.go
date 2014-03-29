package targets

import "fmt"
import "errors"
import "strings"
import "unicode"
import "time"
import "github.com/cnf/go-claw/clog"

type TargetManager struct {
    targets map[string]Target
    target_cmds map[string]map[string]*Command
}

// Create and initialize a new TargetManager object
func NewTargetManager() *TargetManager {
    ret := &TargetManager{ 
            targets    : make(map[string]Target),
            target_cmds: make(map[string]map[string]*Command),
        }
    return ret
}

// Adds a new target using the specified module, given the specified name and parameters
func (t *TargetManager) Add(module, name string, params map[string]string) error {
    // Validate name
    var err error
    if err := validateTargetName(name); err != nil {
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

    // Fetch the command list
    tcmdlist := tgt.Commands()
    if tcmdlist == nil {
        clog.Warn("warning: %s::%s returned an empty command list!", module, name)
    } else {
        t.target_cmds[name] = tcmdlist
    }

    return nil
}

// Remove a target instance from the list
func (t *TargetManager) Remove(name string) error {
    if _, ok := t.targets[name]; !ok {
        return errors.New("cannot remove " + name + ": does not exist")
    }
    if err := t.targets[name].Stop(); err != nil {
        return err
    }
    delete(t.targets, name)
    if _, ok := t.target_cmds[name]; ok {
        delete(t.target_cmds, name)
    }
    return nil
}

// Stops all target instances and removes them
func (t *TargetManager) Stop() error {
    for k := range t.targets {
        if err := t.Remove(k); err != nil {
            return err
        }
    }
    t.targets    = make(map[string]Target)
    t.target_cmds = make(map[string]map[string]*Command)
    return nil
}

// Parses command, determines which target should run it, checks the provided parameters,
// and if all is good - run the command.
func (t *TargetManager) RunCommand(cmdstring string) error {
    splitstr := strings.SplitN(cmdstring, "::", 2)
    tstart := time.Now()
    if len(splitstr) != 2 {
        return fmt.Errorf("invalid command string '%s', expected it to contain '::'", cmdstring)
    }
    // Validate the name of the target we just parsed out
    if err := validateTargetName(splitstr[0]); err != nil {
        return err
    }
    tgtname := splitstr[0]

    if _, ok := t.targets[tgtname]; !ok {
        return fmt.Errorf("command '%s' uses a target '%s' that does not exist", cmdstring, tgtname)
    }

    // Split the command
    splitcmd := splitQuoted(splitstr[1])
    if len(splitcmd)  == 0 {
        return fmt.Errorf("empty target command in '%s'", cmdstring)
    }
    tcommand := splitcmd[0]
    tparams := splitcmd[1:]

    // Check if the instance provided a commands list to check
    if _, ok := t.target_cmds[tgtname]; ok && t.target_cmds[tgtname] != nil {
        // Check if the command exists for this target
        if _, ok := t.target_cmds[tgtname][tcommand]; !ok {
            return fmt.Errorf("command '%s' not recognized by target '%s'", tcommand, tgtname)
        }
        // Validate all parameters
        pc := 0
        tparams_n := make([]string, 0)
        for prm := 0; prm < len(t.target_cmds[tgtname][tcommand].Parameters); prm++ {
            if pc >= len(tparams) {
                // Parameter not present, check if required
                if !t.target_cmds[tgtname][tcommand].Parameters[prm].Optional {
                    clog.Error("Non-optional parameter %s missing for command %s, target %s",
                            t.target_cmds[tgtname][tcommand].Parameters[prm].Name,
                            tcommand,
                            tgtname,
                        )
                    return fmt.Errorf("non-optional parameter '%s' missing for command '%s', target '%s'",
                            t.target_cmds[tgtname][tcommand].Parameters[prm].Name,
                            tcommand,
                            tgtname,
                        )
                }
            } else {
                pval, err := t.target_cmds[tgtname][tcommand].Parameters[prm].Validate(tparams[pc])
                if (err != nil) {
                    // validation returned an error
                    clog.Error("Parameter validation %s failed for command %s, target %s",
                            t.target_cmds[tgtname][tcommand].Parameters[prm].Name,
                            tcommand,
                            tgtname,
                        )
                    return err
                }
                tparams_n = append(tparams_n, pval)
            }
            pc++
        }
        // replace the original parameters with the validated parameters
        tparams = tparams_n
    }
    // Run the command
    clog.Debug("--> Process cmd '%s' took: %s", cmdstring, time.Since(tstart).String())
    tstart = time.Now()
    err := t.targets[tgtname].SendCommand(tcommand, tparams...)
    clog.Debug("--> Execute cmd '%s' took: %s", cmdstring, time.Since(tstart).String())
    return err
}

// Split a string containing quoted strings on newlines, quotes, ... 
// Supports escaping of space, newline, ...
func splitQuoted(s string) []string {
    var ret []string
    var curr = make([]rune, len(s))
    var cpos = 0
    var quoted = ' '
    var escaped = false
    sr := strings.NewReader(s)

    for {
        r, _, err := sr.ReadRune()
        if err != nil {
            // Append last
            if cpos != 0 {
                ret = append(ret, string(curr[0:cpos]))
                cpos = 0
            }
            break
        }
        switch r {
        case ' ':
            if quoted != ' ' {
                if escaped {
                    curr[cpos] = '\\'
                    cpos++
                    escaped = false
                }
                curr[cpos] = ' '
                cpos++
            } else if escaped {
                curr[cpos] = ' '
                cpos++
                escaped = false
            } else if cpos != 0 {
                ret = append(ret, string(curr[0:cpos]))
                cpos = 0
            }
        case '"', '\'':
            if escaped {
                curr[cpos] = r
                cpos++
                escaped = false
            } else if quoted == r {
                // Quoted string closed
                // Don't add to list yet, whitespace should follow if 
                // it's a new string/parameter, otherwise treat it as
                // the same
                //ret = append(ret, string(curr[0:cpos]))
                //cpos = 0
                quoted = ' '
            } else if quoted == ' ' {
                // New quote, start new entry
                quoted = r
            } else {
                curr[cpos] = r
                cpos++
            }
        case '\\':
            if escaped {
                curr[cpos] = '\\'
                cpos++
                escaped = false
            } else {
                escaped = true
            }
        default:
            if unicode.IsSpace(r) {
                // the other white space - cannot escape this!
                if quoted != ' ' {
                    curr[cpos] = r
                    cpos++
                } else if cpos != 0 {
                    // Add to lst
                    ret = append(ret, string(curr[0:cpos]))
                    cpos = 0
                }
            } else {
                if escaped {
                    curr[cpos] = '\\'
                    cpos++
                    escaped = false
                }
                curr[cpos] = r
                cpos++
            }
        }
    }
    return ret
}

func validateTargetName(name string) error {
    if name == "" {
        return errors.New("a target name cannot be empty")
    }
    if strings.ContainsAny(name, "\t\n\r :@!+=*") {
        return fmt.Errorf("target name '%s' cannot contain whitespace, ':', '@', '!', '+', '=' or '*' characters")
    }
    return nil
}

