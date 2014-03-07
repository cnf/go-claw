package dispatcher

import "github.com/cnf/go-claw/commandstream"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/targets/denon"
import "github.com/cnf/go-claw/setup"
import "github.com/cnf/go-claw/clog"

var mydenon *denon.Denon
var targetmap map[string]targets.Targets

func Setup(t map[string]setup.Target) {
    targetmap = make(map[string]targets.Targets)
    for key, value := range t {
       println(key, value.Module)
    }
    targetmap["mydenon"] = denon.Setup("192.168.178.58", 23, "X2000")
}

func Dispatch(rc *commandstream.RemoteCommand) bool {
    clog.Debug("repeat: %2d - key: %s - source: %s", rc.Repeat, rc.Key, rc.Source)
    if rc.Key == "KEY_VOLUMEUP" {
        targetmap["mydenon"].SendCommand("VolumeUp")
        clog.Debug("sending VolUP to denon")
    } else if rc.Key == "KEY_OK" {
        targetmap["mydenon"].SendCommand("Volume", "50")
        SetMode("bar")
        clog.Debug("sending Vol50 to denon")
    } else if rc.Key == "KEY_PLAY" {
        clog.Debug("Current mode: %s", GetMode())
    }
    return true
}
