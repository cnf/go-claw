package main

import "github.com/cnf/go-claw/targets/denon"
import "github.com/cnf/go-claw/targets/plex"

func RegisterAllTargets() {
    denon.Register()
    plex.Register()
}
