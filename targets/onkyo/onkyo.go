package onkyo

//import "strings"
import "net"
import "errors"
import "time"
import "sync"
//import "fmt"
//import "encoding/hex"
//import "github.com/tarm/goserial"
import "github.com/cnf/go-claw/clog"
import "github.com/cnf/go-claw/targets"

// Transport indicates the transport type the Onkyo Reciever uses
type Transport int
const (
    // TransportTCP indicates the useage of TCP
    TransportTCP   Transport = iota
    // TransportSerial indicates the useage of a Serial line
    TransportSerial Transport = iota
)

// OnkyoReceiver structure
type OnkyoReceiver struct {
    Name string
    Transport Transport

    Serialdev string
    Host string
    AutoDetect bool
    Model string
    Identifier string

    rxmu sync.Mutex
    seqnr int64
    rxQchan chan rxCommand
    rxRchan chan rxCommand

    con net.Conn
    mu sync.Mutex
    lastsend time.Time
}

type rxCommand struct {
    msg string
    rxtime time.Time
    txtime time.Time
    seq int64
}


// Register registers the Onkyo Module in the target manager
func Register() {
    targets.RegisterTarget("onkyo", createOnkyoReceiver)
    //targets.RegisterAutoDetect(OnkyoAutoDetect)
}

func (d *OnkyoReceiver) Stop() error {
    return nil
}

func (r *OnkyoReceiver) addRxCommand(msg string) (int64, error) {
    var push rxCommand

    push.msg  = msg
    push.txtime = time.Now()
    // Safeguard the sequence number
    r.rxmu.Lock()
    defer r.rxmu.Unlock()
    push.seq = r.seqnr
    r.seqnr++

    // Push on the channel - but don't block
    select {
    case r.rxQchan <- push:
    default:
        return -1, errors.New("could not push expected message onto channel")
    }
    return r.seqnr - 1, nil
}

func (r *OnkyoReceiver) expectRxCommand(seqnr int64, timeout int) (*rxCommand, error) {
    var tm = time.Now().Add(time.Duration(timeout) * time.Millisecond)
    for {
        // Determine the current timeout
        w := tm.Sub(time.Now())
        if (w <= 0) {
            return nil, errors.New("timeout getting expected command")
        }
        select {
        case msg, ok := <- r.rxRchan:
            if (!ok) {
                // Oops?
                clog.Error("onkyo:expectRxCommand: Could not read from response channel!")
                return nil, errors.New("could not read from response channel")
            } else if msg.seq < seqnr {
                // Older sequence found, skip
                clog.Warn("onkyo:expectRxCommand: Older sequence found: %d:%s, expected %d - discarding", msg.seq, msg.msg, seqnr)
                continue
            } else if (msg.seq == seqnr) {
                // we found our match
                ret := new(rxCommand)
                *ret = msg
                return ret, nil
            } else {
                clog.Error("onkyo:expectRxCommand: sequence number skipped - desynchronized??")
                return nil, errors.New("sequence number skipped, - desynchronized?")
            }
        case <- time.After(w):
            return nil, errors.New("timeout getting expected command")
        }
    }
}

func (r *OnkyoReceiver) readOnkyoResponses(qchan, rchan chan rxCommand, conn net.Conn) {
    // Make sure to close the response channel
    defer close(rchan)
    var rcmd *OnkyoFrameTCP
    var err error
    var expectlist = make([]rxCommand, 0)
    var currseq int64
    currseq = -1
    for {
        // Every 100ms, check everything
        conn.SetReadDeadline(time.Now().Add(time.Duration(100) * time.Millisecond))
        rcmd, err = ReadOnkyoFrameTCP(conn)
        if err != nil {
            nerr, ok := err.(net.Error)
            if !ok {
                clog.Warn("onkyo:readOnkyoResponses frame error: %s", err.Error())
            } else if !nerr.Temporary() && !nerr.Timeout() {
                clog.Error("onkyo:readOnkyoResponses: %s - exiting go-routine", err.Error())
                // Close the response channel
                return
            }
        }
        // Check if we have expected commands waiting for us
        for {
            select {
            case cmd, ok := <-qchan:
                if !ok {
                    // Response channel closed?
                    clog.Error("onkyo:readOnkyoResponses: error reading repsonse channel - aborting")
                    return
                }
                //clog.Debug("Added expected command: %s (%d)", cmd.msg, cmd.seq)
                expectlist = append(expectlist, cmd)
                if currseq < 0 {
                    currseq = cmd.seq - 1
                }
                continue
            default:
            }
            // No responses - break out of the loop
            break
        }
        if (rcmd == nil) {
            continue
        }
        // Walk backward, only respond to latest request
        for i := len(expectlist) - 1; i >= 0; i-- {
            // Remove frames older than 16 seconds
            if (time.Since(expectlist[i].txtime) > (time.Duration(16000) * time.Millisecond)) ||
               (expectlist[i].seq <= currseq) {
                // remove from list
                clog.Debug("onkyo:readOnkyoResponses: removing %d:%s from list...", expectlist[i].seq, expectlist[i].msg)
                expectlist = append(expectlist[:i], expectlist[i+1:]...)
                continue
            }
            if (expectlist[i].msg[0:3] == rcmd.Message()[0:3]) {
                var rv rxCommand
                // Synchronize sequences
                //clog.Debug("onkyo:readOnkyoResponses: found matching cmd: %d:%s ", expectlist[i].seq, expectlist[i].msg)
                currseq = expectlist[i].seq
                rv.rxtime = time.Now()
                rv.txtime = expectlist[i].txtime
                rv.seq = currseq
                rv.msg = rcmd.Message()
                rchan <- rv
                // New sequences are added to the end - so we can remove all previous entries
                expectlist = expectlist[i+1:]
                break
            }
        }
    }
}

func (r *OnkyoReceiver) doConnect() error {
    if (r.Transport == TransportSerial) {
        return errors.New("onkyo: serial connection is not implemented!")
    }
    if (r.con != nil) {
        return nil
    }
    var autodetected = false
    for {
        if (r.Host == "") && (r.AutoDetect) {
            if t := OnkyoFind(r.Model, r.Identifier, 3000); t != nil {
                r.Host = t.Detected["host"]
                autodetected = true
                clog.Info("onkyo:detected receiver: %s (%s)", r.Model, r.Identifier)
            }
        }
        if r.Host == "" {
            return errors.New("onkyo:doConnect: no host setting found")
        }
        var err error
        r.con, err = net.DialTimeout("tcp", r.Host, time.Duration(5000) * time.Millisecond)
        if err != nil {
            clog.Error("onkyo:doConnect: error sending receiver: %s", err.Error());
            if r.con != nil {
                // Should not happen?
                r.con.Close()
                r.con = nil
            }
            if autodetected {
                // Already tried to autodetect, but failed?
                break
            } else if r.AutoDetect {
                // Retry autodetection
                r.Host = ""
                continue
            }
        } else {
            clog.Info("onkyo: connected to %s", r.Host)
            // All ok - create response channel and launch go-routine
            if r.rxQchan != nil {
                close(r.rxQchan)
            }
            r.rxQchan = make(chan rxCommand, 10) // Buffered channel
            r.rxRchan = make(chan rxCommand, 10) // Buffered channel
            go r.readOnkyoResponses(r.rxQchan, r.rxRchan, r.con)
            return nil
        }
    }
    return errors.New("onkyo:doConnect: unknown error")
}

func (r *OnkyoReceiver) processparams(pname string, params map[string]string) error {
    if params["connection"] == "serial" {
        r.Transport = TransportSerial
    } else {
        // By default assume TCP
        r.Transport = TransportTCP
    }
    r.Name = pname
    switch r.Transport {
    case TransportSerial:
        if _, ok := params["device"]; !ok {
            return errors.New("onkyo: missing 'device' parameter for serial receiver")
        }
        r.Serialdev = params["device"]
        if _, ok := params["type"]; !ok {
            return errors.New("onkyo: missing 'type' parameter for serial receiver")
        }
        // Baudrate is fixed: 9600
    case TransportTCP:
        if _, ok := params["host"]; !ok {
            // No host specified - attempt auto discovery
            var ok bool
            if r.Model, ok = params["model"]; !ok {
                return errors.New("onkyo: missing 'host' or 'type' parameter for TCP receiver")
            }
            r.AutoDetect = true
            if r.Identifier, ok = params["id"]; !ok {
                clog.Warn("onkyo:processparams: missing 'id' parmaeter for type '%s'", params["type"])
            }
            if t := OnkyoFind(r.Model, r.Identifier, 3000); t != nil {
                clog.Info("onkyo: detected receiver: %s (%s)", r.Model, r.Identifier)
                r.Host = t.Detected["host"]
            } else {
                // This is not an error? Try again later
                clog.Warn("onkyo:processparams: could not find receiver model '%s' id '%s'", r.Model, r.Identifier)
                r.Host = ""
            }
        } else {
            // Test if the host is correct
            _, _, err := net.SplitHostPort(params["host"])
            if (err != nil) {
                return errors.New("onkyo: invalid 'host' parameter: not a valid host:port notation")
            }
            r.AutoDetect = false
            r.Host = params["host"]
        }
    }
    return nil
}

// Send a command to the onkyo.
// timeout = timeout to wait for response in ms
//   timeout = 0 -> no response expected.
//   timeout < 0 -> default timeout (15 seconds)
//   timeout > 0 -> timeout in ms
func (r *OnkyoReceiver) sendCmd(cmd string, timeout int) (string, error) {
    // Don't allow commands to be sent simultaneously
    r.mu.Lock()
    defer r.mu.Unlock()
    errcnt := 0
    var waitseq int64
    var err error

    if (timeout != 0) {
        waitseq, err = r.addRxCommand(cmd)
        if err != nil {
            return "", err
        }
    } else {
        waitseq = -1
    }
    for {
        if (errcnt >= 2) {
            return "", errors.New("onkyo: could not send command, retry count exceeded")
        }
        if err := r.doConnect(); err != nil {
            return "", err
        }
        switch r.Transport {
        case TransportTCP:
            // Prevent sending a next command within 50ms
            tdiff := time.Since(r.lastsend)
            if tdiff < (time.Duration(50) * time.Millisecond) {
                time.Sleep((time.Duration(50) * time.Millisecond) - tdiff)
            }
            r.con.SetWriteDeadline(time.Now().Add(time.Duration(500) * time.Millisecond))
            b := NewOnkyoFrameTCP(cmd).Bytes()
            //print(hex.Dump(b))
            _, err := r.con.Write(b)
            r.lastsend = time.Now()
            if (err != nil) {
                // check error type
                if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
                    // Socket error - close, and retry
                    r.con.Close()
                    r.con = nil
                } else if (errcnt == 1) {
                    // Second retry that failed - reconnect
                    r.con.Close()
                    r.con = nil
                }
                errcnt++
                continue;
            }
            if (waitseq >= 0 ) {
                // default timeout = 15 seconds
                if timeout < 0 {
                    timeout = 15000
                }
                cmd, err := r.expectRxCommand(waitseq, timeout)
                if (err != nil) {
                    return "", err
                }
                return cmd.msg, nil
            }
            return "", nil
        case TransportSerial:
            return "", errors.New("onkyo: serial protocol not implemented")
        }
        break
    }
    return "", errors.New("onkyo: unknown error sending a command")
}

func createOnkyoReceiver(name string, params map[string]string) (targets.Target, error) {
    clog.Debug("onkyo: creating receiver '%s'", name)
    var ret OnkyoReceiver

    // Process incoming parameters
    if err := ret.processparams(name, params); err != nil {
        clog.Error(err.Error())
        return nil, err
    }
    // 5 seconds in the past
    ret.lastsend = time.Now().Add(time.Duration(-5) * time.Second)
    if err := ret.doConnect(); err != nil {
        clog.Warn("onkyo: could not connect to reciever: %s", err.Error())
    }
    return &ret, nil
}

// SendCommand sends a command to the receiver
func (r *OnkyoReceiver) SendCommand(cmd string, args ...string) error {
    return r.onkyoCommand(cmd, args)
}

