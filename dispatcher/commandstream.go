package dispatcher

import "fmt"

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
        return false
    }
    for {
        fmt.Printf("Nr of listeners: %d\n", self.count)
        select {
        case tmp, ok := <- self.Ch:
            if (!ok) {
                fmt.Printf("Error encountered while reading the next command\n")
                return false
            }
            *cmd = *tmp
            return true
        case err := <- self.ChErr:
            self.err = err
            if (self.Fatal) {
                self.count--
            }
            fmt.Printf("Listener exited and reported an error: %v\n", err)
            if (self.count > 0) {
                continue
            }
            return false
        }
    }
    return false
}

func (self *CommandStream) Error() error {
    return self.err
}
