package plex

import "fmt"
import "net"
import "net/http"
import "net/url"
import "strconv"
import "time"
import "io/ioutil"
import "encoding/xml"

import "github.com/cnf/go-claw/clog"

type timelineXML struct {
    Address      string `xml:"address,attr"`
    ContainerKey string `xml:"containerKey,attr"`
    Controllable string `xml:"controllable,attr"`
    Continuing   bool   `xml:"continuing,attr"`
    Duration     string `xml:"duration,attr"`
    GUID         string `xml:"guid,attr"`
    Key          string `xml:"key,attr"`
    Location     string `xml:"location,attr"`
    MachineID    string `xml:"machineIdentifier,attr"`
    Mute         string `xml:"mute,attr"`
    Port         string `xml:"port,attr"`
    RatingKey    string `xml:"ratingKey,attr"`
    Repeat       string `xml:"repeat,attr"`
    SeekRange    string `xml:"seekRange,attr"`
    Shuffle      string `xml:"shuffle,attr"`
    State        string `xml:"state,attr"`
    Time         string `xml:"time,attr"`
    Type         string `xml:"type,attr"`
    Volume       string `xml:"volume,attr"`
}

type mediaContainerXML struct {
    CommandID string        `xml:"commandID,attr"`
    Location  string        `xml:"location,attr"`
    Timelines []timelineXML `xml:"Timeline"`
}

func (p *Plex) subscribe() {
    go p.listen()
    time.Sleep(3 * time.Second)
    go p.subscriberLoop()
}

func (p *Plex) listen() {
    l, err := net.Listen("tcp", ":0")
    if err != nil { return }
    lport := l.Addr().(*net.TCPAddr).Port
    addr := fmt.Sprintf(":%d", lport)

    s := &http.Server{Addr: addr, Handler: p}
    p.listenport = lport
    clog.Debug("Plex: subscription listener on port `%d`", lport)

    clog.Debug("Plex: %$ v", s.Serve(l))
}

func (p *Plex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    body, rerr := ioutil.ReadAll(r.Body)
    if rerr != nil { return }
    var mc mediaContainerXML
    xmlerr := xml.Unmarshal(body, &mc)
    if xmlerr != nil { return }
    loc := mc.Location
    tls := make(map[string]timelineXML)
    for _, tl := range mc.Timelines {
        tls[tl.Type] = tl
    }
    p.setTimeline(loc, tls)
}

func (p *Plex) subscriberLoop() {
    for {
        burl := p.getURL()
        if burl == "" {
            time.Sleep(3 * time.Second)
            continue
        }
        if !p.hasCapability("timeline") {
            time.Sleep(5 * time.Second)
            continue
        }
        surl := fmt.Sprintf("%s%s", burl, "/player/timeline/subscribe")
        u, _ := url.Parse(surl)
        q := u.Query()
        q.Set("commandID", strconv.Itoa(p.getCommandID()))
        q.Set("port", strconv.Itoa(p.listenport))
        q.Set("protocol", "http")
        u.RawQuery = q.Encode()

        request, _ := http.NewRequest("GET", u.String(), nil)
        request.Header.Add("X-Plex-Client-Identifier", p.uuid)
        request.Header.Add("X-Plex-Device-Name", "Claw")

        // FIXME: cleaner timeouts in go1.3
        client := &http.Client{ Transport: &http.Transport{Dial: dialTimeout}, }

        resp, err := client.Do(request)
        if err != nil {
            if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
                clog.Warn("Plex: Sub ERR: %s", err.Error())
                p.mu.Lock()
                p.url = ""
                p.capabilities = []string{}
                p.mu.Unlock()
                p.tlmu.Lock()
                p.timelines = nil
                p.location = ""
                p.tlmu.Unlock()
            } else {
                clog.Warn("Plex: Sub warn: %s", err.Error())
            }
            time.Sleep(5 * time.Second)
            continue
        }
        // FIXME: do something useful
        // body, err := ioutil.ReadAll(resp.Body)
        resp.Body.Close()
        time.Sleep(30 * time.Second)
    }
}
