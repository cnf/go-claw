package lircsocket

import "net"
import "io"
import "strings"
import "bufio"
import "strconv"
import "time"

import "github.com/cnf/go-claw/listeners"
import "github.com/cnf/go-claw/clog"

type LircSocketListener struct {
    Path string
    // conn net.Conn
    reader *bufio.Reader
}

func Register() {
    listeners.RegisterListener("lircsocket", Create)
}

func Create(params map[string]string) (l listeners.Listener, ok bool) {
    // TODO: VALIDATE PARAMS
    sl := &LircSocketListener{}
    if val, ok := params["path"]; ok {
        sl.Path = val
    } else {
        clog.Warn("Incorrect parameters")
        return nil, false
    }
    return sl, true
}

func (self *LircSocketListener) setup(cs *listeners.CommandStream) bool {
    clog.Debug("Opening socket: %s", self.Path)
    c, err := net.Dial("unix", self.Path)
    // If there is no socket to bind to during setup, we fail.
    if err != nil {
        clog.Warn("Socket setup failed for %s", self.Path)
        cs.ChErr <- err
        return false
    }
    self.reader = bufio.NewReader(c)
    return true
}

func (self *LircSocketListener) RunListener(cs *listeners.CommandStream) {
    // var err error
    // self.conn, err = net.Dial("unix", self.Path)
    // if err != nil {
        // cs.ChErr <- err
        // return
    // }
    if (!self.setup(cs)) {
        cs.Fatal = true
        return
    }
    for {
        str, err := self.reader.ReadString('\n')
        if err != nil {
            if err != io.EOF {
                // Remote end closed socket
                clog.Error("Unknown error occured: %s", err.Error())
            } else {
                clog.Error("Socket closed by remote host: %s", err.Error())
                time.Sleep(1000 * time.Millisecond)
                // var err error
                // self.conn, err = net.Dial("unix", self.Path)
                if (!self.setup(cs)) {
                    time.Sleep(3000 * time.Millisecond)
                    continue
                }
                continue
            }
            // cs.ChErr <- err
            // return
            time.Sleep(1000 * time.Millisecond)
            continue
        }

        out := strings.Split(strings.TrimSpace(str), " ")
        if (len(out) != 4) {
            clog.Error("Length of split '%v' is not 4!", str)
            continue
        }
        rpt, err := strconv.ParseInt(out[1], 16, 0)
        if (err != nil) {
            clog.Error("Could not parse %v, not a number?", out[1])
            continue
        }
        now := time.Now()
        cs.Ch <- &listeners.RemoteCommand{ Code: out[0], Repeat: int(rpt), Key: out[2], Source: out[3], Time: now }

    }
}
