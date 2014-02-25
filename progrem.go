package main

import "fmt"
import "github.com/cnf/progrem/listeners"
import "github.com/cnf/progrem/dispatcher"

func main() {
    println("running")
    // cs := &dispatcher.CommandStream{Ch: make(chan *dispatcher.RemoteCommand)}
    cs := dispatcher.NewCommandStream()
    defer cs.Close()
    var out dispatcher.RemoteCommand

    cs.AddListener(&listeners.LircSocketListener{Path: "/var/run/lirc/lircd"})
    cs.AddListener(&listeners.LircSocketListener{Path: "/tmp/echo.sock"})

    for cs.Next(&out) {
        if cs.HasError() {
            fmt.Printf("An error occured somewhere: %v\n", cs.GetError())
            cs.ClearError()
        }
        fmt.Printf("code: %s - repeat: %3d - key: %s - name: %s\n", out.Code, out.Repeat, out.Key, out.Source)
    }
}
