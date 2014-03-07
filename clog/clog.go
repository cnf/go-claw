package clog

import "fmt"
import "strings"
import "io"
import "os"
import "time"

const chsize = 10

const (
    NONE = iota // = 0
    DEBUG
    INFO
    WARN
    ERROR
    FATAL
)

var lvl_names = [...]string{
    "NONE ",
    "DEBUG",
    "INFO ",
    "WARN ",
    "ERROR",
    "FATAL",
}

var ch chan *clogger
var cfgch chan *Config
var stopch chan bool
var cfg *Config

type clogger struct {
    message string
    level int
    time time.Time
    // err error
    // source string
}

type Config struct {
    writer   io.Writer
    loglevel int
}

func init() {
    ch = make(chan *clogger, chsize)
    stopch = make(chan bool)

    cfg = &Config{writer: os.Stderr, loglevel: DEBUG}
    cfgch = make(chan *Config)
    go runlogger(ch, cfgch)
}


func runlogger(cl chan *clogger, cf chan *Config) {
    var buf []byte
    running := true
    for {
        select {
        case newcfg, ok := <- cf:
            println("changing config")
            if !ok {
                // Config channel closed? Terminate logger
                running = false
                cf = nil
                if len(cl) == 0 {
                    stopch <- true
                    return
                }
                continue
            }
            if newcfg.writer != nil {
                cfg.writer = newcfg.writer
            }
            cfg.loglevel = newcfg.loglevel
        case chn, ok := <- cl:
            if !ok {
                continue
            }
            if (cfg.writer == nil) || (cfg.loglevel > chn.level) {
                continue
            }
            buf = buf[:0]

            buf = append(buf, (lvl_names[chn.level] + " - ")... )
            buf = append(buf, (strings.TrimSpace(chn.message))...)
            // now := time.Now()
            // const layout = "Jan 2, 2006 at 3:04pm (MST)"
            // const layout = time.Stamp
            // fmt.Printf("%s - %s\n", now.Format(layout), chn.message)
            if len(buf) > 0 && buf[len(buf)-1] != '\n' {
                buf = append(buf, '\n')
            }
            _, err := cfg.writer.Write(buf)
            if err != nil {
                // OOPS!
            }
            if (!running) && (len(cl) == 0) {
                stopch <- true
                return
            }
        }
    }
}

func Setup(c *Config) {
    cfgch <- c
}

func SetLogLevel(i int) {
    cfgch <- &Config{writer: cfg.writer, loglevel: i}
}

func Stop() {
    if (cfgch == nil) {
        return
    }
    Info("Shutting down logger")
    close(cfgch)
    cfgch = nil
    <- stopch
    close(stopch)
}

func Fatal(format string, a ...interface{}) {
    ch <- &clogger{message: fmt.Sprintf(format, a...), level: FATAL, time: time.Now()}
}

func Error(format string, a ...interface{}) {
    ch <- &clogger{message: fmt.Sprintf(format, a...), level: ERROR, time: time.Now()}
}

func Warn(format string, a ...interface{}) {
    ch <- &clogger{message: fmt.Sprintf(format, a...), level: WARN, time: time.Now()}
}

func Info(format string, a ...interface{}) {
    ch <- &clogger{message: fmt.Sprintf(format, a...), level: INFO, time: time.Now()}
}

func Debug(format string, a ...interface{}) {
    ch <- &clogger{message: fmt.Sprintf(format, a...), level: DEBUG, time: time.Now()}
}
