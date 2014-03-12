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

var ch = make(chan *clogger, chsize)
var cfgch = make(chan *Config)
var stopch = make(chan bool)
var cfg = &Config{writer: os.Stderr, loglevel: DEBUG}

type clogger struct {
    message string
    level   int
    time    time.Time
    // err error
    // source string
}

type Config struct {
    writer   io.Writer
    loglevel int
}

func init() {
    go runlogger(ch, cfgch)
}

func runlogger(cl chan *clogger, cf chan *Config) {
    var buf []byte
    defer close(stopch)
    for {
        select {
        case newcfg := <-cf:
            println("changing config")
            if newcfg.writer != nil {
                cfg.writer = newcfg.writer
            }
            cfg.loglevel = newcfg.loglevel
        case chn, ok := <-cl:
            if !ok {
                return
            }
            if cfg.writer == nil || cfg.loglevel > chn.level {
                continue
            }
            buf = buf[:0]

            buf = append(buf, (lvl_names[chn.level] + " - ")...)
            buf = append(buf, (strings.TrimSpace(chn.message))...)
            if len(buf) > 0 && buf[len(buf)-1] != '\n' {
                buf = append(buf, '\n')
            }
            _, err := cfg.writer.Write(buf)
            if err != nil {
                // OOPS!
            }
        }
    }
}

func Setup(c *Config) {
    cfgch <- c
}

func SetLogLevel(i int) {
    if (i >= DEBUG) && (i <= FATAL) {
        cfgch <- &Config{writer: cfg.writer, loglevel: i}
    }
}

func Stop() {
    Info("Shutting down logger")
    close(ch)
    if cfgch == nil {
        return
    }
    <-stopch
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
