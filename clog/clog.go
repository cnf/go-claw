package clog

import "fmt"
import "strings"
// import "time"

var ch chan *clogger
const chsize = 10
var cfgch chan *Config
var cfg *Config

func init() {
    ch = make(chan *clogger, chsize)
    cfg = &Config{Path: "stdout"}
    cfgch = make(chan *Config)
    go runlogger()
}

type clogger struct {
    message string
    severity string
    // err error
    // source string
}

type Config struct {
    Path string
}

func runlogger() {
    for {
        select {
        case cfg, ok := <- cfgch:
            if ok {
                println("new config")
                println(cfg.Path)
            }
        case chn, ok := <- ch:
            if ok {
                // now := time.Now()
                // const layout = "Jan 2, 2006 at 3:04pm (MST)"
                // const layout = time.Stamp
                // fmt.Printf("%s - %s\n", now.Format(layout), chn.message)
                fmt.Printf("%s - %s\n", chn.severity, strings.TrimSpace(chn.message))
            }
        }
    }
}

func Setup(c *Config) {
    cfgch <- c
}

func Info(format string, a ...interface{}) {
    ch <- &clogger{message: fmt.Sprintf(format, a...), severity: "INFO"}
}

func Warn(format string, a ...interface{}) {
    ch <- &clogger{message: fmt.Sprintf(format, a...), severity: "WARNING"}
}

func Error(format string, a ...interface{}) {
    ch <- &clogger{message: fmt.Sprintf(format, a...), severity: "ERROR"}
}

func Debug(format string, a ...interface{}) {
    ch <- &clogger{message: fmt.Sprintf(format, a...), severity: "DEBUG"}
}
