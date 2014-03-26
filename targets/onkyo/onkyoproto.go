package onkyo

import "fmt"
import "bytes"
import "errors"
import "encoding/binary"
//import "encoding/hex"
//import "github.com/cnf/go-claw/clog"

// OnkyoFrame describes the main interface for an object to parse and
// construct an Onkyo remote control message
type OnkyoFrame interface {
    SetMessage(string)
    Message() string
    Bytes() []byte
    Parse([]byte) error
}

// OnkyoFrameSerial implements the OnkyoFrame for serial communication
type OnkyoFrameSerial struct {
    Msg string
}

// OnkyoFrameTCP implements the OnkyoFrame for network/TCP communication
type OnkyoFrameTCP struct {
    Msg string
}

// SetMessage sets the message to construct the frame with
func (c *OnkyoFrameTCP) SetMessage(msg string) {
    c.Msg = msg
}
func NewOnkyoFrameTCP(msg string) *OnkyoFrameTCP {
    return &OnkyoFrameTCP{Msg: msg}
}

// Bytes returns the []byte of the constructed message
func (c *OnkyoFrameTCP) Bytes() []byte {
    if (c.Msg == "") {
        return make([]byte, 0)
    }
    buf := new(bytes.Buffer)
    msg := c.Msg
    if msg[0] != '!' {
        msg = "!1" + msg
    }
    // Build the ISCP packet
    binary.Write(buf, binary.BigEndian, []byte("ISCP")) // ISCP Magic
    binary.Write(buf, binary.BigEndian, uint32(16)) // Header size
    binary.Write(buf, binary.BigEndian, uint32(len(msg))) // Data length
    binary.Write(buf, binary.BigEndian, uint8(1)) // Version
    binary.Write(buf, binary.BigEndian, uint8(0)) // Reserved
    binary.Write(buf, binary.BigEndian, uint8(0)) // Reserved
    binary.Write(buf, binary.BigEndian, uint8(0)) // Reserved
    binary.Write(buf, binary.BigEndian, []byte(msg)) // Data
    binary.Write(buf, binary.BigEndian, uint8(0x0D)) // Carriage return
    //binary.Write(buf, binary.BigEndian, uint8(0x19)) // EOF
    //binary.Write(buf, binary.BigEndian, uint8(0x0A)) // Line feed
    //clog.Debug(hex.Dump(buf.Bytes()))
    return buf.Bytes()
}

// Parse parses an incoming []byte, validates and extracts the message
func (c *OnkyoFrameTCP) Parse(buf []byte) (error) {
    var magic [4]byte
    var headersize uint32
    var datalen uint32
    var version uint8
    var rfu [3]byte
    //clog.Debug("----- PARSE ----")
    //clog.Debug(hex.Dump(buf))
    if (len(buf) < 16) {
        // Smaller than header
        return errors.New("buffer length smaller than header size of an onkyo message")
    }
    // Determine endpos
    endpos := bytes.IndexByte(buf[16:], 0x19)
    nlpos := bytes.IndexByte(buf[16:], 0x0A)
    crpos := bytes.IndexByte(buf[16:], 0x0D)
    if (endpos < 0) {
        endpos = intMinPositive(endpos, nlpos, crpos)
    }

    // parse the header
    b := bytes.NewReader(buf[0:16])
    if err := binary.Read(b, binary.BigEndian, &magic); err != nil {
        return err
    }
    if string(magic[0:4]) != "ISCP" {
        return errors.New("onkyo message magic mismatch")
    }
    if err := binary.Read(b, binary.BigEndian, &headersize); err != nil {
        return err
    }
    if headersize != 16 {
        return errors.New("onkyo message header length not 16")
    }
    if err := binary.Read(b, binary.BigEndian, &datalen); err != nil {
        return err
    }
    if err := binary.Read(b, binary.BigEndian, &version); err != nil {
        return err
    }
    if version != 1 {
        return fmt.Errorf("unknown onkyo message version, expected 1, got %d", version)
    }
    if err := binary.Read(b, binary.BigEndian, &rfu); err != nil {
        return err
    }
    rxdatalen := uint32(len(buf[16:intMax(endpos, nlpos, crpos)+16])) + 1
    if rxdatalen != datalen {
        return fmt.Errorf("onkyo message data length mismatch: %d != expected %d",
                rxdatalen, datalen,
            )
    }
    if datalen < 2 {
        return fmt.Errorf("onkyo message too short, expected minimum length of 2, got %d", datalen)
    }
    // Get the message
    if buf[16] != '!' {
        return errors.New("onkyo message does not start with expected '!'")
    }
    if buf[17] != '1' {
        return errors.New("onkyo message not coming from receiver, don't know how to handle")
    }
    //clog.Debug("End position: %d", endpos+16)
    // set the message - strip the "!1" start
    c.Msg = string(buf[18:endpos+15])

    return nil
}

// Message returns the message associated with the command
func (c *OnkyoFrameTCP) Message() (string) {
    return c.Msg
}

/////////////////////////////////////////////////////////////////////////////
// TODO: Serial implementation of the messages

func NewOnkyoFrameSerial(msg string) *OnkyoFrameSerial {
    return &OnkyoFrameSerial{Msg: msg}
}
// SetMessage sets the message to construct the frame with
func (c *OnkyoFrameSerial) SetMessage(msg string) {
    c.Msg = msg
}

// Bytes returns the []byte of the constructed message
func (c *OnkyoFrameSerial) Bytes() []byte {
    return nil
}

// Parse parses an incoming []byte, validates and extracts the message
func (c *OnkyoFrameSerial) Parse(buf []byte) (error) {
    return errors.New("not implemented: OnkyoFrameSerial")
}

// Message returns the message associated with the command
func (c *OnkyoFrameSerial) Message() (string) {
    return c.Msg
}

func intMinPositive(i int, ints... int) int {
    min := i
    for _, ci := range ints {
        if (min < 0) && (ci >= 0) {
            min = ci
        } else if (ci >= 0) && (ci < min) {
            min = ci
        }
    }
    return min
}

func intMin(i int, ints... int) int {
    min := i
    for _, ci := range ints {
        if ci < min {
            min = ci
        }
    }
    return min
}

func intMax(i int, ints... int) int {
    max := i
    for _, ci := range ints {
        if ci > max {
            max = ci
        }
    }
    return max
}

