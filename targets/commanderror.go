package targets

import "strings"
import "fmt"

// CommandError is a structure representing an error when executing a command
type CommandError struct {
    target string
    targetfound bool
    command string
    commandfound bool
    params []string
}

// NewCommandError creates a new commanderror
func NewCommandError(tgt string, tgtfound bool, cmd string, cmdfound bool, prms []string) *CommandError {
    if prms == nil {
        prms = make([]string, 0)
    }
    ret := &CommandError{
        target: tgt,
        targetfound: tgtfound,
        command: cmd,
        commandfound: cmdfound,
        params: prms,
    }
    return ret
}
// Target returns the target of the command that failed
func (c *CommandError) Target()       string   { return c.target }
// TargetFound returns if the target of the command that failed existed or not
func (c *CommandError) TargetFound()  bool     { return c.targetfound }
// Command returns the command that failed
func (c *CommandError) Command()      string   { return c.command }
// CommandFound returns if the command that failed existed or not
func (c *CommandError) CommandFound() bool     { return c.commandfound }
// Params returns the parameters passed to the command that failed
func (c *CommandError) Params()       []string { return c.params }

// Error returns the error description string for the command that failed
func (c CommandError) Error() string {
    if (!c.targetfound) {
        return fmt.Sprintf("could not execute '%s::%s \"%s\"': target not found",
                c.target, c.command, strings.Join(c.params, "\", \""),
            )
    } else if (!c.commandfound) {
        return fmt.Sprintf("could not execute '%s::%s \"%s\"': command not found in target",
                c.target, c.command, strings.Join(c.params, "\", \""),
            )
    } else {
        return fmt.Sprintf("could not execute '%s::%s \"%s\"'",
                c.target, c.command, strings.Join(c.params, "\", \""),
            )
    }
}
