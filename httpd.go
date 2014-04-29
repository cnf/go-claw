package main

import "fmt"
import "net/http"

import "github.com/cnf/go-claw/clog"

func startHTTPListener(port int) {
    go func(port int) {
        clog.Info("Starting the http listener on port %d", port)
        portstr := fmt.Sprintf(":%d", port)
        err := http.ListenAndServe(portstr, nil)
        if err != nil {
            clog.Error("could not start httpd: %s", err.Error)
            return
        }
    }(port)
}
