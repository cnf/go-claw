package clog

import "fmt"
import "strings"
import "io"
import "os"
import "time"
import "bytes"
import "log"

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

const (
    Ldate = 1 << iota
    Ltime
    Lmicroseconds
    Llongfile
    Lshortfile
    LstdFlags = Ldate | Ltime
)


var ch = make(chan *clogger, chsize)
var cfgch = make(chan *Config)
var stopch = make(chan bool)
var cfg = &Config{
    Writer: os.Stderr,
    update_loglevel: true,
    Loglevel: DEBUG,
    update_flag: true,
    Flag: LstdFlags,
}
var cw = &clogwriter{bf: new(bytes.Buffer), laststr: ""}


type clogger struct {
    message string
    level   int
    time    time.Time
    // err error
    // source string
}

type clogwriter struct {
    bf *bytes.Buffer
    laststr string
    flag int
}

type Config struct {
    Writer   io.Writer
    Loglevel int
    update_loglevel bool
    Flag int
    update_flag bool
}

func init() {
    // Create writer for the default logger
    cw.flag = log.Flags()
    log.SetOutput(cw)
    go runlogger(ch, cfgch)
}

func (w *clogwriter) Write(p []byte) (n int, err error) {
    rl, err := w.bf.Write(p)
    for {
        str, _ := w.bf.ReadString('\n')
        w.laststr += str
        if strings.IndexByte(w.laststr, '\n') < 0 {
            // no newline in read string -> stop
            break
        }
        // Parse the string to extract the message
        fc := 0
        lf := w.flag
        if (lf & log.Ldate) == log.Ldate {
            fc++
        }
        if (lf & (log.Ltime | log.Lmicroseconds)) != 0 {
            fc++
        }
        if (fc > 0) {
            liststr := strings.Split(w.laststr, " ")
            w.laststr = strings.Join(liststr[fc:], " ")
        }
        /*
        // Uncomment if you want to strip the long and shortfile logging
        if (lf & (log.Lshortfile | log.Llongfile)) != 0 {
            pos := strings.Index(w.laststr, ": ")
            if (pos >= 0) {
                w.laststr = w.laststr[pos+2:]
            }
        }
        */
        // Send to warning level
        dolog(&clogger{message: w.laststr, level: WARN, time: time.Now()})
        w.laststr = ""
    }
    return rl, err
}

// Shameless copy from log/log.go
func itoa(buf *[]byte, i int, wid int) {
    var u uint = uint(i)
    if u == 0 && wid <= 1 {
        *buf = append(*buf, '0')
        return
    }

    // Assemble decimal in reverse order.
    var b [32]byte
    bp := len(b)
    for ; u > 0 || wid > 0; u /= 10 {
        bp--
        wid--
        b[bp] = byte(u%10) + '0'
    }
    *buf = append(*buf, b[bp:]...)
}
func (l *Config) formatHeader(buf *[]byte, t time.Time, file string, line int) {
    //*buf = append(*buf, l.prefix...)
    if l.Flag&(Ldate|Ltime|Lmicroseconds) != 0 {
        *buf = append(*buf, '[')
        if l.Flag&Ldate != 0 {
            year, month, day := t.Date()
            itoa(buf, year, 4)
            *buf = append(*buf, '/')
            itoa(buf, int(month), 2)
            *buf = append(*buf, '/')
            itoa(buf, day, 2)
        }
        if l.Flag&(Ltime|Lmicroseconds) != 0 {
            if l.Flag&Ldate != 0 {
                *buf = append(*buf, ' ')
            }
            hour, min, sec := t.Clock()
            itoa(buf, hour, 2)
            *buf = append(*buf, ':')
            itoa(buf, min, 2)
            *buf = append(*buf, ':')
            itoa(buf, sec, 2)
            if l.Flag&Lmicroseconds != 0 {
                *buf = append(*buf, '.')
                itoa(buf, t.Nanosecond()/1e3, 6)
            }
        }
        *buf = append(*buf, ']')
        *buf = append(*buf, ' ')
    }
    /*
    if l.flag&(Lshortfile|Llongfile) != 0 {
        if l.flag&Lshortfile != 0 {
            short := file
            for i := len(file) - 1; i > 0; i-- {
                if file[i] == '/' {
                    short = file[i+1:]
                    break
                }
            }
            file = short
        }
        *buf = append(*buf, file...)
        *buf = append(*buf, ':')
        itoa(buf, line, -1)
        *buf = append(*buf, ": "...)
    }
    */
}

func runlogger(cl chan *clogger, cf chan *Config) {
    var buf []byte
    defer close(stopch)
    for {
        select {
        case newcfg := <-cf:
            if newcfg == nil {
                // Error reading from config channel -> Abort
                return
            }
            if newcfg.Writer != nil {
                cfg.Writer = newcfg.Writer
            }
            if (cfg.update_loglevel) {
                cfg.Loglevel = newcfg.Loglevel
            }
            if (cfg.update_flag) {
                cfg.Flag = newcfg.Flag
            }
        case chn, ok := <-cl:
            if !ok {
                return
            }
            if cfg.Writer == nil || cfg.Loglevel > chn.level {
                continue
            }
            buf = buf[:0]
            buf = append(buf, (lvl_names[chn.level] + " ")...)
            // No support for line number/filename at the moment
            cfg.formatHeader(&buf, chn.time, "", 0)

            // Append the real message
            buf = append(buf, (strings.TrimSpace(chn.message))...)
            if len(buf) > 0 && buf[len(buf)-1] != '\n' {
                buf = append(buf, '\n')
            }
            _, err := cfg.Writer.Write(buf)
            if err != nil {
                // OOPS!
            }
        }
    }
}

func Setup(c *Config) {
    cw.flag = log.Flags()
    if (cfgch == nil) || (c == nil) {
        return
    }
    c.update_flag = true
    c.update_loglevel = true
    cfgch <- c
}

func SetLogLevel(i int) {
    cw.flag = log.Flags()
    if cfgch == nil {
        return
    }
    if (i >= DEBUG) && (i <= FATAL) {
        cfgch <- &Config{Writer: nil, update_loglevel: true, Loglevel: i, update_flag: false}
    }
}

func SetFlags(flags int) {
    cw.flag = log.Flags()
    cfgch <- &Config{Writer: nil, update_loglevel: false, update_flag: true, Flag: flags }
}

func Stop() {
    Info("Shutting down logger")
    if (ch != nil) {
        close(ch)
        ch = nil
    }
    // Wait for goroutine to exit
    <-stopch

    // Now kill the config channel
    if (cfgch != nil) {
        close(cfgch)
        cfgch = nil
    }
}

func dolog(l *clogger) {
    if (ch == nil) {
        return
    }
    ch <- l
}

func Fatal(format string, a ...interface{}) {
    dolog(&clogger{message: fmt.Sprintf(format, a...), level: FATAL, time: time.Now()})
}

func Error(format string, a ...interface{}) {
    dolog(&clogger{message: fmt.Sprintf(format, a...), level: ERROR, time: time.Now()})
}

func Warn(format string, a ...interface{}) {
    dolog(&clogger{message: fmt.Sprintf(format, a...), level: WARN, time: time.Now()})
}

func Info(format string, a ...interface{}) {
    dolog(&clogger{message: fmt.Sprintf(format, a...), level: INFO, time: time.Now()})
}

func Debug(format string, a ...interface{}) {
    dolog(&clogger{message: fmt.Sprintf(format, a...), level: DEBUG, time: time.Now()})
}

