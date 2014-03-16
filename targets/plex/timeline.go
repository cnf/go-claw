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

func (self *Plex) subscribe() {
    go self.listen()
    time.Sleep(3 * time.Second)
    go self.subscriberLoop()
}

func (self *Plex) listen() {
    http.HandleFunc("/", self.handler)
    l, err := net.Listen("tcp", ":0")
    if err != nil { return }
    lport := l.Addr().(*net.TCPAddr).Port
    self.listenport = lport
    clog.Debug("Plex subscription listener on port `%d`", lport)

    http.Serve(l, nil)
}

func (self *Plex) handler(w http.ResponseWriter, r *http.Request) {
    clog.Debug("Incoming subscription")
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
    self.setTimeline(loc, tls)
}

func (self *Plex) subscriberLoop() {
    for {
        burl := self.getUrl()
        if burl == "" {
            // clog.Debug("Plex: no client found to subscribe to")
            time.Sleep(3 * time.Second)
            continue
        }
        if !self.hasCapability("timeline") {
            // clog.Debug("Plex: Client doesn't support timeline")
            time.Sleep(3 * time.Second)
            continue
        }
        // clog.Debug("Sending subscribe")
        surl := fmt.Sprintf("%s%s", burl, "/player/timeline/subscribe")
        u, _ := url.Parse(surl)
        q := u.Query()
        q.Set("commandID", strconv.Itoa(self.getCommandID()))
        q.Set("port", strconv.Itoa(self.listenport))
        q.Set("protocol", "http")
        u.RawQuery = q.Encode()

        request, _ := http.NewRequest("GET", u.String(), nil)
        request.Header.Add("X-Plex-Client-Identifier", self.uuid)
        request.Header.Add("X-Plex-Device-Name", "Claw")

        // client := &http.Client{}
        client := &http.Client{ Transport: &http.Transport{Dial: dialTimeout}, }

        resp, err := client.Do(request)
        if err != nil {
            if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
                clog.Debug("Sub ERR: %s", err.Error())
                self.mu.Lock()
                self.url = ""
                self.capabilities = []string{}
                self.mu.Unlock()
            } else {
                clog.Debug("Sub warn: %s", err.Error())
            }
            time.Sleep(5 * time.Second)
            continue
        }
        defer resp.Body.Close()
        // FIXME: do something useful
        // body, err := ioutil.ReadAll(resp.Body)
        // clog.Debug("GOT: %s", string(body))
        time.Sleep(30 * time.Second)
    }
}
