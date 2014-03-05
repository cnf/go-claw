package denon

import "net"
import "fmt"
import "github.com/cnf/go-claw/clog"

type Denon struct {
    name string
    addr *net.TCPAddr
}

func Setup(host string, port int, name string) *Denon {
    clog.Debug("Initializing Denon")
    tmp, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
    if err != nil {
        clog.Error("Could not Initialize Denon: %s", err)
        return nil
    }
    return &Denon{addr: tmp, name: name}

}

func (self *Denon) SendCommand(cmd string, args string) bool {
    return true
}

func (self *Denon) ssendCommand(str string) bool {
    if self.addr == nil {
        clog.Warn("No address to sent Denon command to.")
        return false
    }
    conn, err := net.DialTCP("tcp", nil, self.addr)
    // defer conn.Close()
    if err != nil {
        clog.Error("Connection failed: %s", err)
        defer conn.Close()
        return false
    }
    fmt.Fprintf(conn, "%s\r", str)
    conn.Close()
    return true
}

func (self *Denon) VolumeUp() {
    self.ssendCommand("MVUP")
}

func (self *Denon) Volume(val int) {
    self.ssendCommand(fmt.Sprintf("MV%d", val))
}

const jsonStream = `
{
  "denon": {
    "x2000": [
      {
        "name": "VolumeUp",
        "command": "MVUP",
        "help": "Turn the volume up"
      },
      {
        "name": "Volume",
        "command": "MV<args>",
        "args": "int(0:80)",
        "help": "Set the volume to <args>"
      }
    ]
  }
}
`
