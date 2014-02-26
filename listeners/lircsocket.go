package listeners
// package main

import "net"
import "io"
import "strings"
import "bufio"
import "strconv"
import "time"

import "github.com/cnf/go-claw/dispatcher"
import "github.com/cnf/go-claw/clog"

type LircSocketListener struct {
    Path string
    conn net.Conn
    reader *bufio.Reader
}

func (self *LircSocketListener) Setup(cs *dispatcher.CommandStream) bool {
    // var err error
    // self.conn, err = net.Dial("unix", self.Path)
    c, err := net.Dial("unix", self.Path)
    if err != nil {
        cs.ChErr <- err
        return false
    }
    self.reader = bufio.NewReader(c)
    return true
}

func (self *LircSocketListener) RunListener(cs *dispatcher.CommandStream) {
    // var err error
    // self.conn, err = net.Dial("unix", self.Path)
    // if err != nil {
        // cs.ChErr <- err
        // return
    // }
    if (!self.Setup(cs)) {
        cs.Fatal = true
        return
    }
    for {
        // fmt.Printf("DEBUG: enter for - 1\n")
        // reader := bufio.NewReader(self.conn)
        // fmt.Printf("DEBUG: enter for - 2\n")
        now := time.Now()
        str, err := self.reader.ReadString('\n')
        // fmt.Printf("DEBUG: enter for - 3\n")
        if err != nil {
            if err != io.EOF {
                // Remote end closed socket
                clog.Error("Unknown error occured: %s", err.Error())
            } else {
                clog.Error("Socket closed by remote host: %s", err.Error())
                time.Sleep(1000 * time.Millisecond)
                // var err error
                // self.conn, err = net.Dial("unix", self.Path)
                if (!self.Setup(cs)) {
                    clog.Debug("DEBUG: setup failed")
                    time.Sleep(3000 * time.Millisecond)
                    continue
                }
                continue
            }
            cs.ChErr <- err
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
        cs.Ch <- &dispatcher.RemoteCommand{ Code: out[0], Repeat: int(rpt), Key: out[2], Source: out[3], Time: now }

    }
}
