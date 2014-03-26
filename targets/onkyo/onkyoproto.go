package onkyo

import "fmt"
import "bytes"
import "errors"
import "io"
import "bufio"
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

const OnkyoMagic = "ISCP"


func parseHeader(h []byte) (datalen uint32, err error) {
    var magic [4]byte
    var headersize uint32
    var version uint8
    var rfu [3]byte
    datalen = 0

    if len(h) < 16 {
        return datalen, errors.New("expected 16 bytes onkyo frame header")
    }
    b := bytes.NewReader(h[0:16])
    if err := binary.Read(b, binary.BigEndian, &magic); err != nil {
        return datalen, err
    }
    if string(magic[0:4]) != "ISCP" {
        return datalen, errors.New("onkyo message magic mismatch")
    }
    if err := binary.Read(b, binary.BigEndian, &headersize); err != nil {
        return datalen, err
    }
    if headersize != 16 {
        return datalen, errors.New("onkyo message header length not 16")
    }
    if err := binary.Read(b, binary.BigEndian, &datalen); err != nil {
        return datalen, err
    }
    if err := binary.Read(b, binary.BigEndian, &version); err != nil {
        return datalen, err
    }
    if version != 1 {
        return datalen, fmt.Errorf("unknown onkyo message version, expected 1, got %d", version)
    }
    if err := binary.Read(b, binary.BigEndian, &rfu); err != nil {
        return datalen, err
    }
    return datalen, nil
}

func (c *OnkyoFrameTCP) parseData(buf []byte, datalen uint32) error {
    // Determine endpos
    endpos := intMinPositive(
            bytes.IndexByte(buf, 0x19),
            bytes.IndexByte(buf, 0x0A),
            bytes.IndexByte(buf, 0x0D),
        )
    if endpos < 0 {
        return fmt.Errorf("onkyo message is missing data terminator")
    }

    msgdata := buf[0:endpos]

    if len(msgdata) < 2 {
        return fmt.Errorf("onkyo message too short, expected minimum length of 2, got %d", datalen)
    }
    // Get the message
    if msgdata[0] != '!' {
        return errors.New("onkyo message does not start with expected '!'")
    }
    if msgdata[1] != '1' {
        return errors.New("onkyo message not coming from receiver, don't know how to handle")
    }
    // set the message - strip the "!1" start
    c.Msg = string(msgdata[2:])

    return nil
}

// Reads the frame from an io.Reader instance
func (c *OnkyoFrameTCP) ReadFrom(r io.Reader) error {
    br := bufio.NewReader(r)
    // Scan for the magic
    for {
        c, err := br.ReadByte()
        if err != nil {
            return err
        }
        // Not the start of the magic
        if c != OnkyoMagic[0] {
            continue
        }
        for i := 0; i < (len(OnkyoMagic) - 1); i++ {
            c, err := br.ReadByte()
            if err != nil {
                return err
            }
            if (c != OnkyoMagic[i]) {
                // Not the magic
                continue
            }
        }
        // If we get here, this is a valid frame!
        break
    }
    var frame [16]byte
    for i := 0; i < len(OnkyoMagic); i++ {
        frame[i] = OnkyoMagic[i]
    }
    rlen, err := br.Read(frame[len(OnkyoMagic):])
    if err != nil {
        return err
    }
    if (rlen + len(OnkyoMagic)) != 16 {
        return errors.New("error reading onkyo frame header")
    }
    datalen, err := parseHeader(frame[0:])
    if err != nil {
        return err
    }
    var data = make([]byte, datalen)
    rlen, err = br.Read(data)
    if err != nil {
        return err
    }
    if rlen != len(data) {
        return errors.New("received data length mismatch")
    }
    return c.parseData(data, datalen)
}

// Parse parses an incoming []byte, validates and extracts the message
func (c *OnkyoFrameTCP) Parse(buf []byte) (error) {

    if (len(buf) < 16) {
        // Smaller than header
        return errors.New("buffer length smaller than header size of an onkyo message")
    }
    // parse the header
    datalen, err := parseHeader(buf)
    if (err != nil) {
        return err
    }
    return c.parseData(buf[16:], datalen)
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

