package dispatcher

import "github.com/cnf/go-claw/commandstream"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/targets/denon"
import "github.com/cnf/go-claw/setup"
import "github.com/cnf/go-claw/clog"

var mydenon *denon.Denon
var targetmap map[string]targets.Target

func Setup(t map[string]setup.Target) {
    targetmap = make(map[string]targets.Target)
    for key, val := range t {
       t, ok := targets.GetTarget(val.Module, key, val.Params)
       if ok {
           targetmap[key] = t
       }
    }
    // targetmap["mydenon"] = denon.Setup("192.168.178.58", 23, "X2000")
}

func Dispatch(rc *commandstream.RemoteCommand) bool {
    clog.Debug("repeat: %2d - key: %s - source: %s", rc.Repeat, rc.Key, rc.Source)
    if rc.Key == "KEY_VOLUMEUP" {
        targetmap["myX2000"].SendCommand("VolumeUp")
        clog.Debug("sending VolUP to denon")
    } else if rc.Key == "KEY_OK" {
        targetmap["myX2000"].SendCommand("Volume", "50")
        SetMode("bar")
        clog.Debug("sending Vol50 to denon")
    } else if rc.Key == "KEY_PLAY" {
        clog.Debug("Current mode: %s", GetMode())
    }
    return true
}
