package main

import "net"
import "bufio"
import "time"
import "os/signal"
import "os"
import "fmt"

func echoServer(c net.Conn) {
    buffrw := bufio.NewWriter(c)
    defer c.Close()
    var list = []string{
            "000000037ff07be9 00 KEY_POWER PH00SBLe",
            //"000000037ff07be9 00 KEY_PLAY PH00SBLe",
            //"000000037ff07bdd 00 KEY_OK PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEDOWN PH00SBLe",

            /*
            "000000037ff07bef 00 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 02 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 03 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 04 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 05 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 06 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 07 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 00 KEY_VOLUMEUP PH00SBLe",
            "000000037ff07bef 01 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 02 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 03 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 04 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 05 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 06 KEY_VOLUMEDOWN PH00SBLe",
            "000000037ff07bef 07 KEY_VOLUMEDOWN PH00SBLe",
            */
            //"000000037ff07bdd 00 KEY_OK PH00SBLe",
            "000000037ff07be9 00 KEY_POWER PH00SBLe",
        }
    for {
        fmt.Printf("Sending.")
        for _, element := range list {
            _, err := buffrw.WriteString(element+"\n")
            if err != nil {
                //panic("ERROR: " + err.Error())
                println("Error: " + err.Error())
                return
            }
            buffrw.Flush()
            fmt.Printf(".")
            time.Sleep(1000 * time.Millisecond)
        }
        fmt.Printf("\nWaiting...\n")
        time.Sleep(4000 * time.Millisecond)
    }
}

func main() {
    l, err := net.Listen("unix", "/tmp/echo.sock")
    if err != nil {
        println("listen error", err.Error())
        return
    }
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func() {
        <- sigc
        l.Close()
        os.Exit(1)
    }()

    for {
        fd, err := l.Accept()
        if err != nil {
            println("accept error", err.Error())
            return
        }
        go echoServer(fd)
    }
}
