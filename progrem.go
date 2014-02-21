package main

import "fmt"
import "github.com/cnf/progrem/listeners"
import "github.com/cnf/progrem/dispatcher"

func main() {
    println("running")
    cs := &dispatcher.CommandStream{Ch: make(chan *dispatcher.RemoteCommand)}
    defer close(cs.Ch)
    var out dispatcher.RemoteCommand
    // go watch_socket(cs)
    cs.AddListener(&listeners.LircSocketListener{"/var/run/lirc/lircd"})
    for cs.Next(&out) {
        // if !cs.Next(&out) { if cs.Error() != nil { return cs.Error() } }
        // fmt.Printf("Got message? %v\n", <-ch);
        // out := <-cs.ch
        // fmt.Printf("%s", cs

        fmt.Printf("code: %s - repeat: %3d - key: %s - name: %s\r\n", out.Code, out.Repeat, out.Key, out.Source)
    }
}
