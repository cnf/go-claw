package clog

import "fmt"
// import "time"

var ch chan *clogger

func init() {
    ch = make(chan *clogger)
    println(">>>>>>>> starting logger")
    go runlogger()
}

type clogger struct {
    message string
    severity string
    // err error
    // source string
}

func runlogger() {
    for {
        chn, ok := <- ch
        if ok {
            // now := time.Now()
            // const layout = "Jan 2, 2006 at 3:04pm (MST)"
            // const layout = time.Stamp
            // fmt.Printf("%s - %s\n", now.Format(layout), chn.message)
            fmt.Printf("%s - %s\n", chn.severity, chn.message)
        }
    }
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
