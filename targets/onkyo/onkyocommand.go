package onkyo

import "fmt"
import "bytes"
import "errors"
import "encoding/binary"

// OnkyoCommand describes the main interface for an object to parse and
// construct an Onkyo remote control message
type OnkyoCommand interface {
    SetMessage(string)
    Message() string
    Bytes() []byte
    Parse([]byte) error
}

// OnkyoCommandSerial implements the OnkyoCommand for serial communication
type OnkyoCommandSerial struct {
    Msg string
}

// OnkyoCommandTCP implements the OnkyoCommand for network/TCP communication
type OnkyoCommandTCP struct {
    Msg string
}

// SetMessage sets the message to construct the frame with
func (c *OnkyoCommandTCP) SetMessage(msg string) {
    c.Msg = msg
}
func NewOnkyoCommandTCP(msg string) *OnkyoCommandTCP {
    return &OnkyoCommandTCP{Msg: msg}
}

// Bytes returns the []byte of the constructed message
func (c *OnkyoCommandTCP) Bytes() []byte {
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
    binary.Write(buf, binary.BigEndian, uint8(0x19)) // EOF
    binary.Write(buf, binary.BigEndian, uint8(0x0D)) // Carriage return
    binary.Write(buf, binary.BigEndian, uint8(0x0A)) // Line feed
    return buf.Bytes()
}

// Parse parses an incoming []byte, validates and extracts the message
func (c *OnkyoCommandTCP) Parse(buf []byte) (error) {
    var magic [4]byte
    var headersize uint32
    var datalen uint32
    var version uint8
    var rfu [3]byte

    if (len(buf) < 16) {
        // Smaller than header
        return errors.New("buffer length smaller than header size of an onkyo message")
    }
    // Determine endpos
    endpos := bytes.IndexByte(buf[16:], 0x19) + 16
    if endpos < 16 {
        // No end position
        return errors.New("missing EOF character to terminate the onkyo message")
    }
    nlpos := bytes.IndexByte(buf[endpos:], 0x0A) + endpos
    crpos := bytes.IndexByte(buf[endpos:], 0x0D) + endpos

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
    rxdatalen := uint32(len(buf[16:intMax(endpos, nlpos, crpos)])) + 1
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
    // set the message - strip the "!1" start
    c.Msg = string(buf[18:endpos])

    return nil
}

// Message returns the message associated with the command
func (c *OnkyoCommandTCP) Message() (string) {
    return c.Msg
}

/////////////////////////////////////////////////////////////////////////////
// TODO: Serial implementation of the messages

func NewOnkyoCommandSerial(msg string) *OnkyoCommandSerial {
    return &OnkyoCommandSerial{Msg: msg}
}
// SetMessage sets the message to construct the frame with
func (c *OnkyoCommandSerial) SetMessage(msg string) {
    c.Msg = msg
}

// Bytes returns the []byte of the constructed message
func (c *OnkyoCommandSerial) Bytes() []byte {
    return nil
}

// Parse parses an incoming []byte, validates and extracts the message
func (c *OnkyoCommandSerial) Parse(buf []byte) (error) {
    return errors.New("not implemented: OnkyoCommandSerial")
}

// Message returns the message associated with the command
func (c *OnkyoCommandSerial) Message() (string) {
    return c.Msg
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

