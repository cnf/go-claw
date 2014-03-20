package tools

import "net"
import "fmt"
import "strings"

// BroadcastAddrs returns a list of all broadcast addresses on your
// network interfaces
// FIXME: add IPv6 support
func BroadcastAddrs() (baddr []string) {
    ifaces, _ := net.Interfaces()
    for _, iface := range ifaces {
        // only send on interfaces that are up and can do broadcast
        upandbcast := net.FlagBroadcast | net.FlagUp
        if iface.Flags & upandbcast != upandbcast {
            continue
        }
        addrs, _ := iface.Addrs()
        if len(addrs) == 0 {
            continue
        }
        for _, addr := range addrs {
            var bcast []string
            ip, ipn, err := net.ParseCIDR(addr.String())
            if err != nil {
                continue
            }
            // we only handle ipv4 addresses for now
            ipv4 := ip.To4()
            if ipv4 == nil {
                continue
            }
            for i, m := range ipn.Mask {
                if m == 0 {
                    bcast = append(bcast, "255")
                } else {
                    bcast = append(bcast, fmt.Sprintf("%d", ipv4[i]))
                }
            }
            baddr = append(baddr, strings.Join(bcast, "."))
        }
    }
    return baddr
}
