package onkyo

import "encoding/binary"
import "bytes"
import "errors"

type OnkyoCommand interface {
    SetMessage(string)
    Message() (string)
    Bytes() ([]byte, error)
    Parse([]byte) (error)
}

type OnkyoCommandSerial struct {
    Msg string
}

type OnkyoCommandTCP struct {
    Msg string
}

func (c *OnkyoCommandTCP) SetMessage(msg string) {
    c.Msg = msg
}

func (c *OnkyoCommandTCP) Bytes() ([]byte, error) {
    buf := new(bytes.Buffer)
    if c.Msg == "" {
        return nil, errors.New("OnkyoCommandTCP:Bytes(): Empty message, cannot construct emtpy command")
    }
    msg := c.Msg
    if msg[0] != '1' {
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
    return buf.Bytes(), nil
}

func (c *OnkyoCommandTCP) Parse(buf []byte) (error) {
    var magic [4]byte
    var headersize uint32
    var datalen uint32
    var version uint8
    var rfu [3]byte

    if (len(buf) < 16) {
        // Smaller than header
        return errors.New("OnkyoCommandTCP:Parse(): buffer length smaller than header size")
    }
    // Determine endpos
    endpos := bytes.IndexByte(buf[16:], 0x19) + 16
    if endpos < 16 {
        // No end position
        return errors.New("OnkyoCommandTCP:Parse(): EOF byte missing!")
    }

    // parse the header
    b := bytes.NewReader(buf[0:16])
    if err := binary.Read(b, binary.BigEndian, &magic); err != nil {
        return err
    }
    if string(magic[0:4]) != "ISCP" {
        return errors.New("OnkyoCommandTCP:Parse(): magic mismatch")
    }
    if err := binary.Read(b, binary.BigEndian, &headersize); err != nil {
        return err
    }
    if headersize != 16 {
        return errors.New("OnkyoCommandTCP:Parse(): header length not 16")
    }
    if err := binary.Read(b, binary.BigEndian, &datalen); err != nil {
        return err
    }
    if err := binary.Read(b, binary.BigEndian, &version); err != nil {
        return err
    }
    if version != 1 {
        return errors.New("OnkyoCommandTCP:Parse(): unknown version, expected 1")
    }
    if err := binary.Read(b, binary.BigEndian, &rfu); err != nil {
        return err
    }
    if uint32(len(buf[16:endpos])) != datalen {
        return errors.New("OnkyoCommandTCP:Parse(): data length mismatch")
    }
    if datalen < 2 {
        return errors.New("OnkyoCommandTCP:Parse(): data too short, expected minimum length of 2")
    }
    // Get the message
    if buf[16] != '!' {
        return errors.New("OnkyoCommandTCP:Parse(): does not start with expected '!'")
    }
    if buf[17] != '1' {
        return errors.New("OnkyoCommandTCP:Parse(): Message not coming from receiver, don't know how to handle.")
    }
    // set the message - strip the "!1" start
    c.msg = string(buf[18:endpos])

    return nil
}

func (c *OnkyoCommandTCP) Message() (string) {
    return c.Msg
}

/////////////////////////////////////////////////////////////////////////////
// TODO: Serial implementation of the messages

func (c *OnkyoCommandSerial) SetMessage(msg string) {
    c.Msg = msg
}

func (c *OnkyoCommandSerial) Bytes() ([]byte, error) {
    return nil, errors.New("OnkyoCommandSerial not implemented")
}

func (c *OnkyoCommandSerial) Parse(buf []byte) (error) {
    return errors.New("OnkyoCommandSerial not implemented")
}

func (c *OnkyoCommandSerial) Message() (string) {
    return c.Msg
}
