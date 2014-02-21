package main

import "fmt"
import "github.com/cnf/progrem/listeners"
import "github.com/cnf/progrem/dispatcher"

func main() {
    println("running")
    cs := &dispatcher.CommandStream{Ch: make(chan *dispatcher.RemoteCommand)}
    defer close(cs.Ch)
    var out dispatcher.RemoteCommand
    cs.AddListener(&listeners.LircSocketListener{"/var/run/lirc/lircd"})
    for cs.Next(&out) {
        fmt.Printf("code: %s - repeat: %3d - key: %s - name: %s\r\n", out.Code, out.Repeat, out.Key, out.Source)
    }
}
