package commandstream

import "github.com/cnf/go-claw/clog"

type CommandStream struct {
    Ch chan *RemoteCommand
    ChErr chan error
    Fatal bool
    count int
    err error
}

func NewCommandStream() *CommandStream {
    cs := &CommandStream{ Ch: make(chan *RemoteCommand), ChErr: make(chan error), count: 0, err: nil}
    return cs
}

func (self *CommandStream) Count() int {
    return self.count
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
    self.Fatal = false
}

func (self *CommandStream) Next(cmd *RemoteCommand) bool {
    if (self.count <= 0) {
        clog.Debug("No listeners, shutting down")
        return false
    }
    for {
        select {
        case tmp, ok := <- self.Ch:
            if (!ok) {
                clog.Warn("Error encountered while reading the next command")
                return false
            }
            *cmd = *tmp
            return true
        case err := <- self.ChErr:
            self.err = err
            if (self.Fatal) {
                clog.Debug("Fatal error, listener shutting down")
                self.count--
            }
            clog.Error("Listener exited and reported an error: %v", err)
            if (self.count > 0) {
                continue
            }
            clog.Error("Nothing to listen to!")
            return false
        }
    }
    return false
}

func (self *CommandStream) Error() error {
    return self.err
}
