package dispatcher

type CommandStream struct {
    Ch chan *RemoteCommand
    ChErr chan error
    count int
    err error
}

func NewCommandStream() *CommandStream {
    cs := &CommandStream{ Ch: make(chan *RemoteCommand), ChErr: make(chan error), count: 0, err: nil}
    return cs
}

func (self *CommandStream) Close() {
    close(self.Ch)
    close(self.ChErr)
}

func (self *CommandStream) AddListener(l RemoteListener) bool {
    go l.RunListener(self)
    self.count++
    return true
}

func (self *CommandStream) HasError() bool {
    return (self.err != nil)
}

func (self *CommandStream) GetError() error {
    return self.err
}

func (self *CommandStream) ClearError() {
    self.err = nil
}

func (self *CommandStream) Next(cmd *RemoteCommand) bool {
    var ok bool
    tmp, ok := <-self.Ch
    *cmd = *tmp
    return ok
}

func (self *CommandStream) Error() error {
    return self.err
}
