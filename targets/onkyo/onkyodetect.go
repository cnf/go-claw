package onkyo

import "net"
import "time"
import "strings"
import "strconv"
import "github.com/cnf/go-claw/tools"
import "github.com/cnf/go-claw/clog"

const onkyoPort = 60128
const onkyoDetectMagic = "!xECNQSTN"

// TargetDevice is the structure the autodetect code returns
type TargetDevice struct {
    Name string
    Model string
    Params map[string]string   // Parameters that should be saved in the config file
    Detected map[string]string // Detected parameters that should be re-evaluated each connect
}

func runOnkyoDetect(ch chan *TargetDevice, timeout int) {
    defer close(ch)
    addrs := tools.BroadcastAddrs()
    port := strconv.Itoa(onkyoPort)

    cmd := &OnkyoFrameTCP{onkyoDetectMagic}
    transmitStr := cmd.Bytes()

    udpaddr, err := net.ResolveUDPAddr("udp4", ":0")
    if err != nil {
        return
    }
    conn, err := net.ListenUDP("udp4", udpaddr)
    if err != nil {
        clog.Error("targets/onkyo: Could not listen on UDP Address %v: %s", udpaddr, err.Error())
        return
    }
    defer conn.Close()
    conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))

    // Send broadcast
    for _, a := range addrs {
        udpdest, err := net.ResolveUDPAddr("udp4", a+":"+port)
        if err != nil {
            clog.Warn("targets/onkyo: ResolveUDPAddr('udp4','%s:%s'): %s", a, port, err.Error())
            continue
        }
        conn.WriteToUDP(transmitStr, udpdest)
        // Ignore errors - just continue
    }

    // Now loop listening to responses with timeout
    buf := make([]byte, 1024)
    for {
        rlen, raddr, err := conn.ReadFromUDP(buf);
        if err != nil {
            if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() || nerr.Timeout() {
                if (!nerr.Timeout()) {
                    // Unexpected!
                    clog.Error("targets/onkyo: ReadFromUDP: %s", err.Error())
                }
                return
            }
            continue
        }
        //fmt.Printf("Got response from: %v (len: %d):\n", raddr, rlen)
        resp := new(OnkyoFrameTCP)
        if err := resp.Parse(buf[0:rlen]); err != nil {
            // fmt.Printf("Error parsing Onkyo response: %s", err.Error())
            clog.Warn("targets/onkyo: Parse Onkyo Response: %s", err.Error())
            continue
        }
        // Parse the message:
        // ECN<model name>/<ISCP port>/<region:DX|XX|JJ>/<id>
        rmsg := resp.Message()
        if len(rmsg) <= 3 {
            clog.Warn("targets/onkyo: Autodetect response too short: %d", len(rmsg))
            continue
        }
        if rmsg[0:3] != "ECN" {
            clog.Warn("targets/onkyo: Autodetect Unexpected response: %s", rmsg)
            continue
        }
        splitmsg := strings.Split(rmsg[3:], "/")
        if len(splitmsg) != 4 {
            clog.Warn("targets/onkyo: Autodetect invalid response format, expected 4 fields: %s", rmsg)
            continue
        }
        tgt := &TargetDevice{}
        tgt.Name = splitmsg[0] + " (" + splitmsg[3] + ")"
        tgt.Model = splitmsg[0]
        tgt.Params = make(map[string]string)
        tgt.Params["id"] = splitmsg[3]
        tgt.Params["model"] = splitmsg[0]
        tgt.Detected = make(map[string]string)
        tgt.Detected["host"] = raddr.String()

        ch <- tgt

        // Once we got a response, wait maximum 50ms for other responses
        conn.SetReadDeadline(time.Now().Add(time.Duration(20) * time.Millisecond))
    }
}

// OnkyoFind is used to get a list of Onkyo receivers detected on the network
func OnkyoFind(model, id string, timeout int) *TargetDevice {
    ch := make(chan *TargetDevice)
    go runOnkyoDetect(ch, timeout)
    for {
        select {
            case s, ok := <- ch:
                if (!ok) {
                    return nil
                }
                if (model == "") && (id == "") { 
                    clog.Warn("No onkyo receiver specified, accepting first match: %s (%s)",
                            s.Params["model"],
                            s.Params["id"],
                        )
                    return s
                }
                if (s.Model != model) {
                    continue
                }
                if (id == "") {
                    clog.Warn("no identifier speficier, picking first available with id '%s'",
                            s.Params["id"],
                        )
                    return s
                }
                if (id == s.Params["id"]) {
                    return s
                }
        }
    }
    return nil
}

// OnkyoAutoDetect is used to get a list of Onkyo receivers detected on the network
func OnkyoAutoDetect(timeout int) []TargetDevice {
    ch := make(chan *TargetDevice)
    go runOnkyoDetect(ch, timeout)
    var ret []TargetDevice
    //timer := time.NewTimer(time.Millisecond * time.Duration(timeout + 200))
    for {
        select {
            case s, ok := <- ch:
                if !ok {
                    return ret
                }
                ret = append(ret, *s)
        }
    }
    return ret
}

