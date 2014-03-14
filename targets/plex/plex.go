package plex

import "net"
import "fmt"
import "net/http"
import "time"
import "sync"
// import "net/url"

import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-gdm"

type Plex struct {
    name string
    cname string
    url string
    proto string
    commands map[string]Commander
    capabilities []string
    uuid string
    mu sync.Mutex
    // Content-Type v: plex/media-player
    // Resource-Identifier  v: 87615ee6-5b86-4a8d-abf6-e3b4f0e72311
    // Protocol v: plex
    // Version  v: 1.0.10.199-939d4f2b
    // Device-Class v: HTPC
    // Name v: yBox
    // Port v: 3005
    // Product  v: Plex Home Theater
    // Protocol-Capabilities    v: navigation,playback,timeline
    // Protocol-Version v: 1
}

func Register() {
    targets.RegisterTarget("plex", Create)
}

func Create(name string, params map[string]string) (t targets.Target, ok bool) {
    clog.Debug("Plex Create called")
    p := &Plex{name: name, }
    p.proto = "http"
    if val,ok := params["name"]; ok {
        p.cname = val
    }
    go p.plexWatcher()
    p.commands = PHT
    p.uuid = "1A5C18A3-C398-4A50-A6CE-FCFDDD7FC1F2"
    return p, true
}

func (self *Plex) plexWatcher() {
    w, err := gdm.WatchPlayers(5)
    if err != nil {
        clog.Error("Can't watch for plex: %s", err.Error())
    }
    for gdm := range w.Watch {
        if gdm.Props["Name"] != self.cname {
            continue
        }
        url := fmt.Sprintf("%s://%s:%s", self.proto, gdm.Address.IP.String(), gdm.Props["Port"])
        self.mu.Lock()
        self.url = url
        self.mu.Unlock()
    }
    //
}

func (self *Plex) getUrl() (url string) {
    self.mu.Lock()
    url = self.url
    self.mu.Unlock()
    return
}

func (self *Plex) SendCommand(cmd string, args ...string) bool {
    switch cmd {
    case "MenuUp":
        // return self.restSend()
    default:
        clog.Debug("Looking up %s in the Plex map", cmd)
        if val, ok := self.commands[cmd]; ok {
            p, err := val.Command(args...)
            if err != nil {
                return false
            }
            return self.plexGet(p)
        }
    }
    return false
}

func (self *Plex) plexGet(str string) bool {
    purl := fmt.Sprintf("%s%s", self.getUrl(), str)
    // clean, err := url.Parse(purl)
    clog.Debug(">>> Plex get %s", purl)
    // FIXME: cleaner timeouts in go1.3
    client := http.Client{ Transport: &http.Transport{Dial: dialTimeout}, }
    resp, err := client.Get(purl)

    if err != nil {
        clog.Error("FIXME: go1.3 - %s", err.Error())
        return false
    } else {
        resp.Body.Close()
    }
    return true
}

func dialTimeout(network, addr string) (net.Conn, error) {
    return net.DialTimeout(network, addr, time.Duration(1 * time.Second))
}
