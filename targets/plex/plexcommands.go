package plex

import "fmt"

import "github.com/cnf/go-claw/targets"

// Commands returns a list of commands
func (p *Plex) Commands() map[string]*targets.Command {
    cmds := map[string]*targets.Command {
        "power"       : targets.NewCommand("Power state",
            targets.NewParameter("powerstate", "The power state").SetList("on"),
        ),
        "nav"         : targets.NewCommand("Navigation",
            targets.NewParameter("direction", "Navigation direction").SetList("up", "down", "left", "right", "select", "back", "home"),
        ),
        "playback"    : targets.NewCommand("Playback",
            targets.NewParameter("Play state", "Change the playback state").SetList("play", "pauze", "stop", "next", "previous", "forward", "backwards"),
        ),
    }
    return cmds
}

func (p *Plex) navigation(cmd string) error {
    navbase := "/player/navigation/"
    playbase := "/player/playback/"

    var path string
    switch cmd {
    case "up":
        if p.isNav() {
            path = navbase + "moveUp"
        } else {
            path = playbase + "skipNext"
        }
    case "down":
        if p.isNav() {
            path = navbase + "moveDown"
        } else {
            path = playbase + "skipPrevious"
        }
    case "left":
        if p.isNav() {
            path = navbase + "moveLeft"
        } else {
            path = playbase + "stepBack"
        }
    case "right":
        if p.isNav() {
            path = navbase + "moveRight"
        } else {
            path = playbase + "stepForward"
        }
    case "select":
        if p.isNav() {
            path = navbase + "select"
        } else {
            path = playbase + "play"
        }
    case "back":
        path = navbase + "back"
    case "home":
        path = navbase + "home"
    case "music":
        path = navbase + "music"
    default:
        return fmt.Errorf("could not send `%s` to `%s`", cmd, p.name)
    }
    return p.plexGet(path)
}

func (p *Plex) playback(cmd string) error {
    base := "/player/playback/"
    var path string
    switch cmd {
    case "play":
        path = base + "play"
    case "pause":
        path = base + "pause"
    case "stop":
        path = base + "stop"
    case "next":
        path = base + "skipNext"
    case "previous":
        path = base + "skipPrevious"
    case "forward":
        path = base + "stepForward"
    case "backwards":
        path = base + "stepBack"
    default:
        return fmt.Errorf("could not send `%s` to `%s`", cmd, p.name)
    }
    return p.plexGet(path)
}
