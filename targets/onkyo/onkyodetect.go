package onkyo

import "github.com/cnf/go-claw/tools"
import "strconv"
import "net"
import "fmt"
import "time"

const ONKYO_PORT = 60128
const ONKYO_MAGIC = "!xECNQSTN"

func runOnkyoDetect(ch chan string, timeout int) {
    defer close(ch)
    addrs := tools.BroadcastAddrs()
    port := strconv.Itoa(ONKYO_PORT)

    cmd := &OnkyoCommandTCP{ONKYO_MAGIC}
    transmit_str, err := cmd.Bytes()

    if err != nil {
        return
    }

    udpaddr, err := net.ResolveUDPAddr("udp4", ":0")
    if err != nil {
        return
    }
    conn, err := net.ListenUDP("udp4", udpaddr)
    if err != nil {
        fmt.Printf("ERROR: %s\n", err.Error())
        return
    }
    defer conn.Close()
    conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))

    // Send broadcast
    for _, a := range addrs {
        udpdest, err := net.ResolveUDPAddr("udp4", a+":"+port)
        if err != nil {
            continue
        }
        conn.WriteToUDP(transmit_str, udpdest)
        // Ignore errors - just continue
    }

    // Now loop listening to responses with timeout
    buf := make([]byte, 1024)
    for {
        rlen, raddr, err := conn.ReadFromUDP(buf);
        if err != nil {
            if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() || nerr.Timeout() {
                return
            }
            continue
        }
        //fmt.Printf("Got response from: %v (len: %d):\n", raddr, rlen)
        resp := &OnkyoCommandTCP{}
        if err := resp.Parse(buf[0:rlen]); err != nil {
            continue
        }
        ch <- raddr.String() + ":" + resp.Message()

        // Once we got a response, wait maximum 50ms for other responses
        conn.SetReadDeadline(time.Now().Add(time.Duration(50) * time.Millisecond))
    }
}

func OnkyoAutoDetect(timeout int) []string {
    ch := make(chan string)
    go runOnkyoDetect(ch, timeout)
    ret := make([]string, 0)
    //timer := time.NewTimer(time.Millisecond * time.Duration(timeout + 200))
    for {
        select {
            case s, ok := <- ch:
                if !ok {
                    fmt.Printf("Go channel closed - stop!\n");
                    return ret
                }
                ret = append(ret, s)
                //fmt.Printf("Got response: '%s'\n", s)
            //case <- timer.C:
                //fmt.Printf("Timeout!\n")
            //    return ret
        }
    }
    return ret
}

