package targets

import "strings"
import "fmt"

type CommandError struct {
    target string
    targetfound bool
    command string
    commandfound bool
    params []string
}

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
func (c *CommandError) Target()       string   { return c.target }
func (c *CommandError) TargetFound()  bool     { return c.targetfound }
func (c *CommandError) Command()      string   { return c.command }
func (c *CommandError) CommandFound() bool     { return c.commandfound }
func (c *CommandError) Params()       []string { return c.params }

func (c *CommandError) Error() string {
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
