package listeners

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

func (cs *CommandStream) Count() int {
    return cs.count
}

func (cs *CommandStream) Close() {
    if cs == nil {
        return
    }
    close(cs.Ch)
    close(cs.ChErr)
}

func (cs *CommandStream) AddListener(l RemoteListener) bool {
    go l.RunListener(cs)
    cs.count++
    return true
}

func (cs *CommandStream) HasError() bool {
    return (cs.err != nil)
}

func (cs *CommandStream) GetError() error {
    return cs.err
}

func (cs *CommandStream) ClearError() {
    cs.err = nil
    cs.Fatal = false
}

func (cs *CommandStream) Next(cmd *RemoteCommand) bool {
    if (cs.count <= 0) {
        clog.Warn("No listeners, shutting down")
        return false
    }
    for {
        select {
        case tmp, ok := <- cs.Ch:
            if (!ok) {
                clog.Warn("Error encountered while reading the next command")
                return false
            }
            *cmd = *tmp
            return true
        case err := <- cs.ChErr:
            cs.err = err
            if (cs.Fatal) {
                clog.Error("Fatal error, listener shutting down")
                cs.count--
            }
            clog.Error("Listener exited and reported an error: %v", err)
            if (cs.count > 0) {
                continue
            }
            clog.Warn("Nothing to listen to!")
            return false
        }
    }
    return false
}

func (cs *CommandStream) Error() error {
    return cs.err
}
