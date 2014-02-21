package structures

// import "github.com/cnf/progrem/listeners"

type RemoteCommand struct {
    code   string
    repeat int
    key    string
    source   string
}

type CommandStream struct {
    ch chan *RemoteCommand
    err error
}

func (self *CommandStream) AddListener(l RemoteListener) bool {
    go l.RunListener(self.ch)
    return true
}

func (self *CommandStream) Next(cmd *RemoteCommand) bool {
    var ok bool
    tmp, ok := <-self.ch
    *cmd = *tmp
    return ok
}

func (self *CommandStream) Error() error {
    return self.err
}
