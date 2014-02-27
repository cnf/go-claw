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

    cfg = &Config{writer: os.Stderr, loglevel: 5}
    cfgch = make(chan *Config)
    go runlogger()
}


func runlogger() {
    var buf []byte
    running := true
    for {
        select {
        case newcfg, ok := <- cfgch:
            if !ok {
                // Config channel closed? Terminate logger
                running = false
                continue
            }
            Debug("Changing logging config: %s, %d", newcfg.writer, newcfg.loglevel)
            if newcfg.writer != nil {
                cfg.writer = newcfg.writer
            }
            cfg.loglevel = newcfg.loglevel
        case chn, ok := <- ch:
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
            if (!running) && (len(ch) == 0) {
                stopch <- true
                cfgch = nil
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
    Info("Shutting down logger")
    close(cfgch)
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
