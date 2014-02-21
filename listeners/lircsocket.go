package listeners
// package main

import "net"
import "fmt"
import "io"
import "strings"
import "bufio"
import "log"
import "strconv"

import "github.com/cnf/progrem/dispatcher"

type LircSocketListener struct {
    Path string //:= /var/run/lirc/lircd
}

//func (self *LircSocketListener) RunListener(ch chan *dispatcher.RemoteCommand) {
func (self *LircSocketListener) RunListener(cs *dispatcher.CommandStream) {
    c,err := net.Dial("unix", self.Path)
    if err != nil {
        panic(err.Error())
    }
    for {
        reader := bufio.NewReader(c)
        str, err := reader.ReadString('\n')
        if err != nil && err != io.EOF { panic(err.Error()) }

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
        cs.Ch <- &dispatcher.RemoteCommand{ Code: out[0], Repeat: int(rpt), Key: out[2], Source: out[3] }

        // ch <- &CommandStream{ch: make(chan *RemoteCommand), err: err}
    }
}
