package main

import "fmt"
import "time"
import "log"
import "html/template"
import "net"
import "net/http"
import "os"
import "os/signal"
import "strings"
import "strconv"
import "bufio"
import "code.google.com/p/go.net/websocket"

//import "os"

// Channel used by websockets to send to broadcast function
var incmds = make(chan string)
// Channels used to add and remove 'output socket' channels in broadcast function
var chanadd = make(chan chan string)
var chanrm = make(chan chan string)
var exitch = make(chan bool)

func main() {
    sock, err := net.Listen("unix", "/tmp/echo.sock")
    if err != nil {
        log.Println("ERROR: " + err.Error())
        return
    }
    defer sock.Close()

    // Handle ctrl-c
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt)
    go func() {
        <- sigc
        exitch <- true
    }()

    // Listens on unix-socket and spawns new go-routine per socket
    go listenUnixSocket(sock)
    // Handles input from websockets and broadcasts it to registered unix socket handlers
    go handlebroadcast()

    // Add http handlers
    http.HandleFunc("/", httpRoot)
    http.Handle("/sock", websocket.Handler(processWebsocket))
    // Now listen on http 
    go http.ListenAndServe(":8080", nil)
    <- exitch
    sock.Close()
    os.Exit(0)
}

func handleUnixConn(c net.Conn, ch chan string) {
    // Register channel in broadcaster
    chanadd <- ch
    defer func() {
        chanrm <- ch;
        c.Close()
    }()

    buffrw := bufio.NewWriter(c)
    count := 0
    pcmd := ""
    for f := range ch {
        // Parse the incoming message
        split := strings.Split(f, ":")
        if len(split) != 2 {
            continue
        }
        sleep, err := strconv.Atoi(split[1])
        if err != nil {
            sleep = 100
        }
        if pcmd == split[0] {
            count++
        } else {
            count = 0
        }
        // construct and send the message
        sendmsg := fmt.Sprintf("%s %02X %s %s", "000000037ff07bef", count, split[0], "PH00SBLe")
        log.Println("Sending command " + sendmsg)
        buffrw.WriteString(sendmsg + "\n")
        pcmd = split[0]

        // Now sleep
        log.Println("Sleeping " + strconv.Itoa(sleep) + "ms")
        time.Sleep(time.Duration(sleep) * time.Millisecond)
    }
}

func listenUnixSocket(l net.Listener) {
    defer l.Close()
    for {
        fd, err := l.Accept()
        if err != nil {
            log.Println("Error: handlesockets:", err.Error())
            return
        }
        fchan := make(chan string)
        go handleUnixConn(fd, fchan)
    }
    exitch <- true
}

// This sends all incoming messages to all output sockets
func handlebroadcast() {
    var chanarr = make([]chan string,0)
    defer func(){
        log.Println("Broadcast cleaning up...")
        for ri := range(chanarr) {
            close(chanarr[ri])
        }
    }()
    log.Println("Broadcaster ready.")
    for {
        select {
        case nc, ok := <- chanadd:
            // Output socket added, new channel
            if !ok {
                log.Println("ERROR: chanadd read (handlebroadcast)")
                return
            }
            chanarr = append(chanarr, nc)
            log.Printf("Broadcast: added listener, %d in total now\n", len(chanarr))
        case rc, ok := <- chanrm:
            if !ok {
                log.Println("ERROR: chanrm read (handlebroadcast)")
                return
            }
            for ri := range(chanarr) {
                if chanarr[ri] == rc {
                    // Remove the channel
                    chanarr = append(chanarr[:ri], chanarr[ri+1:]...)
                    close(rc)
                }
            }
            log.Printf("Broadcast: removed listener, %d left\n", len(chanarr))
        case msg, ok := <- incmds:
            if !ok {
                log.Println("ERROR: msg read (handlebroadcast)")
                return
            }
            log.Printf("Broadcasting message %s to sockets...\n", msg)
            for si := range(chanarr) {
                chanarr[si] <- msg
            }
        }
    }
}

func httpRoot(w http.ResponseWriter, r *http.Request) {
    indexpg.Execute(w, nil)
}

func processWebsocket(conn *websocket.Conn) {
    var msg string
	log.Println("Websocket connected: ", conn.RemoteAddr().String())
	defer log.Println("Websocket closed")
    for {
        if err := websocket.Message.Receive(conn, &msg); err != nil {
            log.Println("processWebsocket: aborting:", err)
            return
        }
        //log.Println("Received message from websocket: " + msg)
        incmds <- msg
    }
}

var indexpg = template.Must(template.New("index").Parse(`
<html>
    <head>
        <title>Web remote</title>
        <script>
var path = window.location.pathname;
var wsURL = "ws://" + window.location.host +  path.substring(0, path.lastIndexOf('/')) + "/sock";
var ws;
var buttons = {
    "Up"   : { command: "KEY_UP"   , wait: 100, keycode: 38 },
    "Down" : { command: "KEY_DOWN" , wait: 100, keycode: 40 },
    "Left" : { command: "KEY_LEFT" , wait: 100, keycode: 37 },
    "Right": { command: "KEY_RIGHT", wait: 100, keycode: 39 }
};
var keyqueue = [];

function processQueue() {
    if (ws == null) {
        return
    }
    wait = 50
    if ((keyqueue != null) && (keyqueue.length > 0)) {
        v = keyqueue.pop()
        if ("wait" in v) {
            wait = v.wait
        }
        ws.send(v.command + ":" + wait)
        log("Sent command: " + v.command)
    }
    if (keyqueue != null) {
        setTimeout(processQueue, wait)
    } else {
        ws.close()
        ws = null
    }
}

function stopQueue(evt) {
    if (ws != null) {
        ws.close()
        ws = null
    }
    keyqueue = null
}

function onMessage(evt) {
    // Message received
}

document.addEventListener("DOMContentLoaded", function() {
    ws = new WebSocket(wsURL);
    if (ws == null) {
        return
    }
    ws.onopen = function() {
        addButtons()
        processQueue()
    }
    ws.onmessage = onMessage
    ws.onerror   = stopQueue
    ws.onclose   = stopQueue
})
document.addEventListener('keydown', function(event) {
    for (var i in buttons) {
        if (buttons[i].keycode == event.keyCode) {
            sendCmd(i)
            event.preventDefault();
            return false
        }
    }
    return true
}, true);

function sendCmd(button) {
    if ( !(button in buttons)) {
        return
    }
    keyqueue.push(buttons[button])
    log("Queued command: " + buttons[button].command)
}

function log(str) {
    document.getElementById("messages").innerHTML = str + "<br />\n" + document.getElementById("messages").innerHTML
}

function addButtons() {
    for (var i in buttons) {
        // Add the button
        document.getElementById("buttons").innerHTML += "<input type='button' onclick='javascript:sendCmd(\"" + i + "\")' name=\""+i+"\" value=\""+i+"\">\n";
    }
}
function removeButtons() {
    document.getElementById("buttons").innerHTML = ""
}

        </script>
    </head>
    <body>
    <span id="buttons"></span>
    <hr />
    <span id="messages"></span>
    </body>
</html>
`))
