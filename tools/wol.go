package tools

import "net"
import "errors"
import "strings"
import "encoding/hex"

import "github.com/cnf/go-claw/clog"

func Wol(mac string) bool {
    mp, cerr := craftMagicPacket(mac)
    if cerr != nil {
        clog.Debug("Wol failed: %s", cerr.Error())
        return false
    }
    serr := send(mp)
    if serr != nil {
        clog.Debug("Wol failed: %s", serr.Error())
        return false
    }
    return true
}

type magicPacket []byte

func craftMagicPacket(mac string) (magicPacket, error) {
    if len(mac) != 17 {
        println("kaka")
        return nil, errors.New("Invalid MAC Address String: " + mac)
    }
    macBytes, err := hex.DecodeString(strings.Join(strings.Split(mac, ":"), ""))
    if err != nil {
        return nil, err
    }
    b := []uint8{255, 255, 255, 255, 255, 255}
    for i := 0; i < 16; i++ {
        b = append(b, macBytes...)
    }

    return magicPacket(b), nil

}

func send(mp magicPacket) error {
    t, err := net.ResolveUDPAddr("udp", "255.255.255.255:7")
    c, err := net.DialUDP("udp", nil, t)
    if err != nil {
        return err
    }
    written, err := c.Write(mp)
    c.Close()

    if written != 102 {
        return err
    }

    return nil
}
