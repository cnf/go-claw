package plex

import "fmt"

import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/clog"

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

func (p *Plex) navigation(cmd string, repeated int) error {
    navbase := "/player/navigation/"
    playbase := "/player/playback/"

    if repeated == 0 {
        p.setLast("")
    }

    clog.Debug("nav:::::: %s", p.last)
    clog.Debug("plexcmd:nav: cmd %s", cmd)

    var path string
    var err error
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
        if p.isNav() && !p.isLast("seek") {
            path = navbase + "moveLeft"
        } else {
            if repeated == 0 {
                path = playbase + "stepBack"
            } else if repeated % 5 == 0 {
                 p.setLast("seek")
                path, err = p.seekTo("back")
            } else {
                clog.Debug("plexcmd:nav: retun nil")
                return nil
            }
        }
    case "right":
        if p.isNav() && !p.isLast("seek") {
            path = navbase + "moveRight"
        } else {
            if repeated == 0 {
                path = playbase + "stepForward"
            } else if repeated % 5 == 0 {
                p.setLast("seek")
                path, err = p.seekTo("fwd")
            } else {
                clog.Debug("plexcmd:nav: retun nil")
                return nil
            }
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
    if err != nil {
        clog.Debug("plexcmd:nav: ERROR: %v", err)
        return err
    }
    clog.Debug("plexcmd:nav: path: %s", path)
    return p.plexGet(path)
}

func (p *Plex) playback(cmd string, repeated int) error {
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

func (p *Plex) seekTo(dir string) (string, error) {
    // base := "/player/playback/seekTo?offset="
    time, duration, err := p.getOffset()
    if err != nil {
        return "", err
    }
    var offset int
    if dir == "fwd" {
        offset = time + 301000
        if offset >= duration {
            offset = duration
        }
    }
    if dir == "back" {
        offset = time - 120000
        if offset <= 0 {
            offset = 0
        }
    }
    path := fmt.Sprintf("/player/playback/seekTo?offset=%d", offset)
    return path, nil
}
