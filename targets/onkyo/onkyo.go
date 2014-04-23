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

// Starts the Onkyo target instance
func (o *OnkyoReceiver) Start() error {
    return nil
}

// Stop stops the current onkyo target instance
func (o *OnkyoReceiver) Stop() error {
    switch o.Transport {
    case TransportSerial:
    case TransportTCP:
        o.rxmu.Lock()
        defer o.rxmu.Unlock()
        if o.con != nil {
            o.con.Close()
            o.con = nil
        }
        if o.rxQchan != nil {
            close(o.rxQchan)
            o.rxQchan = nil
        }
    }
    return nil
}

func (o *OnkyoReceiver) addRxCommand(msg string) (int64, error) {
    var push rxCommand

    push.msg  = msg
    push.txtime = time.Now()
    // Safeguard the sequence number
    o.rxmu.Lock()
    defer o.rxmu.Unlock()
    if o.rxQchan == nil {
        return -1, errors.New("onkyo disconnected?")
    }
    push.seq = o.seqnr
    o.seqnr++

    // Push on the channel - but don't block
    select {
    case o.rxQchan <- push:
    default:
        return -1, errors.New("could not push expected message onto channel")
    }
    return o.seqnr - 1, nil
}


// Gets an expected response from the 
func (o *OnkyoReceiver) expectRxCommand(seqnr int64, timeout int) (*rxCommand, error) {
    var tm = time.Now().Add(time.Duration(timeout) * time.Millisecond)
    for {
        // Determine the current timeout
        w := tm.Sub(time.Now())
        if (w <= 0) {
            return nil, errors.New("timeout getting expected command")
        }
        select {
        case msg, ok := <- o.rxRchan:
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
            }
            clog.Error("onkyo:expectRxCommand: sequence number skipped - desynchronized??")
            return nil, errors.New("sequence number skipped, - desynchronized?")
        case <- time.After(w):
            return nil, errors.New("timeout getting expected command")
        }
    }
}

// Go routine which reads responses from the sockets and if necessary pushes them back
func (o *OnkyoReceiver) readOnkyoResponses(qchan, rchan chan rxCommand, conn net.Conn) {
    // Make sure to close the response channel
    defer close(rchan)
    var rcmd *OnkyoFrameTCP
    var err error
    var expectlist = make([]rxCommand, 0)
    var currseq int64
    currseq = -1
    for {
        // Every 100ms, check everything
        //conn.SetReadDeadline(time.Now().Add(time.Duration(1000) * time.Millisecond))
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
        clog.Debug("onkyo:readOnkyoResponses: Got '%s'", rcmd.Message())
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

func (o *OnkyoReceiver) doConnect() error {
    if (o.Transport == TransportSerial) {
        return errors.New("onkyo: serial connection is not implemented")
    }
    if (o.con != nil) {
        return nil
    }
    var autodetected = false
    for {
        if (o.Host == "") && (o.AutoDetect) {
            if t := OnkyoFind(o.Model, o.Identifier, 3000); t != nil {
                o.Host = t.Detected["host"]
                autodetected = true
                clog.Info("onkyo:detected receiver: %s (%s)", o.Model, o.Identifier)
            }
        }
        if o.Host == "" {
            return errors.New("onkyo:doConnect: no host setting found")
        }
        var err error
        o.con, err = net.DialTimeout("tcp", o.Host, time.Duration(5000) * time.Millisecond)
        if err != nil {
            clog.Error("onkyo:doConnect: error sending receiver: %s", err.Error());
            if o.con != nil {
                // Should not happen?
                o.con.Close()
                o.con = nil
            }
            if autodetected {
                // Already tried to autodetect, but failed?
                break
            } else if o.AutoDetect {
                // Retry autodetection
                o.Host = ""
                continue
            }
        } else {
            clog.Info("onkyo: connected to %s", o.Host)
            // All ok - create response channel and launch go-routine
            if o.rxQchan != nil {
                close(o.rxQchan)
            }
            o.rxQchan = make(chan rxCommand, 10) // Buffered channel
            o.rxRchan = make(chan rxCommand, 10) // Buffered channel
            go o.readOnkyoResponses(o.rxQchan, o.rxRchan, o.con)
            return nil
        }
    }
    return errors.New("onkyo:doConnect: unknown error")
}

func (o *OnkyoReceiver) processparams(pname string, params map[string]string) error {
    if params["connection"] == "serial" {
        o.Transport = TransportSerial
    } else {
        // By default assume TCP
        o.Transport = TransportTCP
    }
    o.Name = pname
    switch o.Transport {
    case TransportSerial:
        if _, ok := params["device"]; !ok {
            return errors.New("onkyo: missing 'device' parameter for serial receiver")
        }
        o.Serialdev = params["device"]
        if _, ok := params["type"]; !ok {
            return errors.New("onkyo: missing 'type' parameter for serial receiver")
        }
        // Baudrate is fixed: 9600
    case TransportTCP:
        if _, ok := params["host"]; !ok {
            // No host specified - attempt auto discovery
            var ok bool
            if o.Model, ok = params["model"]; !ok {
                return errors.New("onkyo: missing 'host' or 'type' parameter for TCP receiver")
            }
            o.AutoDetect = true
            if o.Identifier, ok = params["id"]; !ok {
                clog.Warn("onkyo:processparams: missing 'id' parmaeter for type '%s'", params["type"])
            }
            if t := OnkyoFind(o.Model, o.Identifier, 3000); t != nil {
                clog.Info("onkyo: detected receiver: %s (%s)", o.Model, o.Identifier)
                o.Host = t.Detected["host"]
            } else {
                // This is not an error? Try again later
                clog.Warn("onkyo:processparams: could not find receiver model '%s' id '%s'", o.Model, o.Identifier)
                o.Host = ""
            }
        } else {
            // Test if the host is correct
            _, _, err := net.SplitHostPort(params["host"])
            if (err != nil) {
                return errors.New("onkyo: invalid 'host' parameter: not a valid host:port notation")
            }
            o.AutoDetect = false
            o.Host = params["host"]
        }
    }
    return nil
}

// Send a command to the onkyo.
// timeout = timeout to wait for response in ms
//   timeout = 0 -> no response expected.
//   timeout < 0 -> default timeout (15 seconds)
//   timeout > 0 -> timeout in ms
func (o *OnkyoReceiver) sendCmd(cmd string, timeout int) (string, error) {
    // Don't allow commands to be sent simultaneously
    o.mu.Lock()
    defer o.mu.Unlock()
    errcnt := 0
    var waitseq int64
    var err error

    if (timeout != 0) {
        waitseq, err = o.addRxCommand(cmd)
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
        if err := o.doConnect(); err != nil {
            return "", err
        }
        switch o.Transport {
        case TransportTCP:
            // Prevent sending a next command within 50ms
            tdiff := time.Since(o.lastsend)
            if tdiff < (time.Duration(50) * time.Millisecond) {
                time.Sleep((time.Duration(50) * time.Millisecond) - tdiff)
            }
            o.con.SetWriteDeadline(time.Now().Add(time.Duration(500) * time.Millisecond))
            b := NewOnkyoFrameTCP(cmd).Bytes()
            //print(hex.Dump(b))
            _, err := o.con.Write(b)
            o.lastsend = time.Now()
            if (err != nil) {
                // check error type
                if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
                    // Socket error - close, and retry
                    o.con.Close()
                    o.con = nil
                } else if (errcnt == 1) {
                    // Second retry that failed - reconnect
                    o.con.Close()
                    o.con = nil
                }
                errcnt++
                continue;
            }
            if (waitseq >= 0 ) {
                // default timeout = 15 seconds
                if timeout < 0 {
                    timeout = 15000
                }
                cmd, err := o.expectRxCommand(waitseq, timeout)
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
func (o *OnkyoReceiver) SendCommand(repeated int, cmd string, args ...string) error {
    return o.onkyoCommand(cmd, args)
}

