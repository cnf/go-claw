package dispatcher

// import "github.com/cnf/progrem/listeners"

type RemoteCommand struct {
    Code   string
    Repeat int
    Key    string
    Source   string
}

type CommandStream struct {
    Ch chan *RemoteCommand
    Err error
}

type RemoteListener interface {
    //RunListener(ch chan *RemoteCommand)
    RunListener(cs *CommandStream)
}

func (self *CommandStream) AddListener(l RemoteListener) bool {
    //go l.RunListener(self.Ch)
    go l.RunListener(self)
    return true
}

func (self *CommandStream) Next(cmd *RemoteCommand) bool {
    var ok bool
    tmp, ok := <-self.Ch
    *cmd = *tmp
    return ok
}

func (self *CommandStream) Error() error {
    return self.Err
}
