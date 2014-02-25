package listeners
// package main

import "net"
import "fmt"
import "io"
import "strings"
import "bufio"
import "log"
import "strconv"
import "time"

import "github.com/cnf/progrem/dispatcher"

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
                fmt.Printf("ERROR: Unknown error occured!\n")
            } else {
                fmt.Printf("========================================\n")
                fmt.Printf("ERROR: Socket closed by remote host!\n")
                time.Sleep(1000 * time.Millisecond)
                // var err error
                // self.conn, err = net.Dial("unix", self.Path)
                if (!self.Setup(cs)) {
                    fmt.Printf("DEBUG: setup failed\n")
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
            log.Println(fmt.Sprintf("ERROR: Length of split '%v' is not 4!\n", str))
            continue
        }
        rpt, err := strconv.ParseInt(out[1], 16, 0)
        if (err != nil) {
            fmt.Printf("ERROR: Could not parse %v, not a number? \n", out[1])
            continue
        }
        cs.Ch <- &dispatcher.RemoteCommand{ Code: out[0], Repeat: int(rpt), Key: out[2], Source: out[3], Time: now }

    }
}
