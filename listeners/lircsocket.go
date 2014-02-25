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
    Path string //:= /var/run/lirc/lircd
}

func (self *LircSocketListener) RunListener(cs *dispatcher.CommandStream) {
    c, err := net.Dial("unix", self.Path)
    if err != nil {
        cs.ChErr <- err
        return
    }
    for {
        reader := bufio.NewReader(c)
        now := time.Now()
        str, err := reader.ReadString('\n')
        if err != nil {
           if err != io.EOF {
               // Remote end closed socket
               fmt.Printf("ERROR: Unknown error occured!\n")
           } else {
               fmt.Printf("ERROR: Socket closed by remote host!\n")
               time.Sleep(1000 * time.Millisecond)
               continue
           }
            cs.ChErr <- err
            return
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
