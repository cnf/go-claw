package main

import "fmt"
import "github.com/cnf/progrem/listeners"
import "github.com/cnf/progrem/structures"

func main() {
    println("running")
    cs := &CommandStream{ch: make(chan *RemoteCommand)}
    defer close(cs.ch)
    var out RemoteCommand
    // go watch_socket(cs)
    cs.AddListener(&LircSocketListener{"/var/run/lirc/lircd"})
    for cs.Next(&out) {
        // if !cs.Next(&out) { if cs.Error() != nil { return cs.Error() } }
        // fmt.Printf("Got message? %v\n", <-ch);
        // out := <-cs.ch
        // fmt.Printf("%s", cs

        fmt.Printf("code: %s - repeat: %3d - key: %s - name: %s\r\n", out.code, out.repeat, out.key, out.source)
    }
}
