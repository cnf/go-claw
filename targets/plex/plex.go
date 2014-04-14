package plex

import "net"
import "fmt"
import "net/http"
import "time"
import "sync"
import "strings"
import "net/url"
import "strconv"

import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"
import "github.com/cnf/go-claw/tools"
import "github.com/cnf/go-gdm"

// Plex struct hold all the Plex information
type Plex struct {
    name string
    wol string
    cname string
    proto string
    commands map[string]commander
    capabilities []string
    uuid string
    listenport int
    // mutex set
    url string
    commandID int
    mu sync.Mutex
    // mutex set
    timelines map[string]timelineXML
    location string
    tlmu sync.Mutex

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

// Register this package in the target list
func Register() {
    targets.RegisterTarget("plex", Create)
}

// Create a new instance of this target
func Create(name string, params map[string]string) (targets.Target, error) {
    p := &Plex{name: name, }
    p.proto = "http"
    if val, ok := params["name"]; ok {
        p.cname = val
    }
    if val, ok := params["wol"]; ok {
        p.wol = val
    }
    go p.plexWatcher()
    p.commands = pht
    // FIXME: generate an actual ID
    p.uuid = "1A5C18A3-C398-4A50-A6CE-FCFDDD7FC1F2"
    p.commandID = 1
    go p.subscribe()
    return p, nil
}

func (d *Plex) Commands() map[string]*targets.Command {
    return nil
}
func (d *Plex) Stop() error {
    return nil
}


// SendCommand receives the command from the dispatcher
func (p *Plex) SendCommand(cmd string, args ...string) error {
    var path string
    var err error
    var val commander
    var ok bool
    switch cmd {
    case "poweron":
        return p.powerOn()
    case "smartup":
        if p.isNav() {
            val, ok = p.commands["moveup"]
        } else {
            val, ok = p.commands["skipnext"]
        }
    case "smartdown":
        if p.isNav() {
            val, ok = p.commands["movedown"]
        } else {
            val, ok = p.commands["skipprevious"]
        }
    case "smartleft":
        if p.isNav() {
            val, ok = p.commands["moveleft"]
        } else {
            val, ok = p.commands["stepback"]
        }
    case "smartright":
        if p.isNav() {
            val, ok = p.commands["moveright"]
        } else {
            val, ok = p.commands["stepforward"]
        }
    case "smartselect":
        if p.isNav() {
            val, ok = p.commands["select"]
        } else {
            val, ok = p.commands["play"]
        }
    default:
        val, ok = p.commands[cmd]
    }
    if !ok {
        err = fmt.Errorf("could not send `%s` to `%s`", cmd, p.name)
        return err
    }
    path, err = val.command(args...)
    if err != nil { return err }
    return p.plexGet(path)
}

func (p *Plex) plexWatcher() {
    w, err := gdm.WatchPlayers(5)
    if err != nil {
        clog.Error("!!!! Can't watch for plex: %s", err.Error())
        return
    }
    for gdm := range w.Watch {
        if gdm.Props["Name"] != p.cname {
            continue
        }
        url := fmt.Sprintf("%s://%s:%s", p.proto, gdm.Address.IP.String(), gdm.Props["Port"])
        caps := strings.Split(gdm.Props["Protocol-Capabilities"], ",")
        p.mu.Lock()
        p.url = url
        p.capabilities = caps
        p.mu.Unlock()
    }
    //
}

func (p *Plex) plexPlaying() {
}

func (p *Plex) getURL() string {
    p.mu.Lock()
    url := p.url
    p.mu.Unlock()
    return url
}

func (p *Plex) getCommandID() int {
    p.mu.Lock()
    id := p.commandID
    p.commandID++
    p.mu.Unlock()
    return id
}

func (p *Plex) hasCapability(c string) bool {
    p.mu.Lock()
    caps := p.capabilities
    p.mu.Unlock()
    for _, v := range caps {
        if c == v {
            return true
        }
    }
    return false
}

func (p *Plex) plexGet(str string) error {
    burl := p.getURL()
    if burl == "" {
        clog.Info("Plex: no url set, client not running?")
        return fmt.Errorf("no url set, client not running?")
    }
    purl := fmt.Sprintf("%s%s", burl, str)
    clog.Debug("Plex: GET %s", purl)
    u, _ := url.Parse(purl)
    q := u.Query()
    q.Set("commandID", strconv.Itoa(p.getCommandID()))
    u.RawQuery = q.Encode()

    request, _ := http.NewRequest("GET", u.String(), nil)
    request.Header.Add("X-Plex-Client-Identifier", p.uuid)
    request.Header.Add("X-Plex-Device-Name", "Claw")

    // FIXME: cleaner timeouts in go1.3
    client := http.Client{ Transport: &http.Transport{Dial: dialTimeout}, }
    resp, err := client.Do(request)
    if err != nil {
        clog.Error("FIXME: go1.3 - %s", err.Error())
        return err
    }
    resp.Body.Close()
    return nil
}

func (p *Plex) setTimeline(loc string, tls map[string]timelineXML) {
    p.tlmu.Lock()
    p.location = loc
    p.timelines = tls
    p.tlmu.Unlock()
}

func (p *Plex) getLocation() string {
    p.tlmu.Lock()
    loc := p.location
    p.tlmu.Unlock()
    return loc
}

func (p *Plex) isNav() bool {
    loc := p.getLocation()
    p.tlmu.Lock()
    tls := p.timelines
    p.tlmu.Unlock()
    if loc == "navigation" {
        return true
    }
    // navigation,fullScreenVideo,fullScreenPhoto,fullScreenMusic

    if (loc == "fullScreenVideo") && (tls["video"].State == "playing") {
        return false
    }
    if (loc == "fullScreenPhoto") && (tls["photo"].State == "playing") {
        return false
    }
    if (loc == "fullScreenMusic") && (tls["music"].State == "playing") {
        return false
    }
    return true
}

func (p *Plex) powerOn() error {
    if p.wol != "" {
        ok := tools.Wol(p.wol)
        if !ok { return fmt.Errorf("can not power on %s", p.name) }
    }
    return fmt.Errorf("do not know how to power on %s", p.name)
}

func dialTimeout(network, addr string) (net.Conn, error) {
    return net.DialTimeout(network, addr, time.Duration(1 * time.Second))
}

