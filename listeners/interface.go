package listeners

import "github.com/cnf/progrem/structures"

type RemoteListener interface {
    RunListener(ch chan *structures.RemoteCommand)
}

func (self *structures.CommandStream) AddListener(l RemoteListener) bool {
    go l.RunListener(self.ch)
    return true
}

func (self *structures.CommandStream) Next(cmd *structures.RemoteCommand) bool {
    var ok bool
    tmp, ok := <-self.ch
    *cmd = *tmp
    return ok
}

func (self *structures.CommandStream) Error() error {
    return self.err
}
