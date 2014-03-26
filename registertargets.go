package main

import "github.com/cnf/go-claw/targets/denon"
import "github.com/cnf/go-claw/targets/plex"
import "github.com/cnf/go-claw/targets/linux"
import "github.com/cnf/go-claw/targets/onkyo"

func registerAllTargets() {
    denon.Register()
    plex.Register()
    linux.Register()
    onkyo.Register()
}
