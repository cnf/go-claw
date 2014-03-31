package main

import "fmt"
//import "time"
import "log"
import "html/template"
import "net"
import "net/http"
import "os"
import "os/signal"
//import "strings"
//import "strconv"
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
        // construct and send the message
        if (pcmd == f) {
            count++
        } else {
            count = 0
        }
        sendmsg := fmt.Sprintf("%s %02X %s %s", "000000037ff07bef", count, f, "PH00SBLe")
        log.Println("Sending command " + sendmsg)
        _, err := buffrw.WriteString(sendmsg + "\n")
        if err != nil {
            log.Println("ERROR: Could not write on socket!")
            return
        }
        buffrw.Flush()
        pcmd = f
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
            for ri := 0; ri < len(chanarr); ri++ {
                if chanarr[ri] == rc {
                    // Remove the channel
                    log.Printf("Closing socket %d...", ri)
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
var connected = false;
var lastsent = ""

// Javascript keycodes: http://www.cambiaresearch.com/articles/15/javascript-char-codes-key-codes
var buttons = {
    "Up"               : { command: "KEY_UP"           , keycode: [ 38 ] },
    "Down"             : { command: "KEY_DOWN"         , keycode: [ 40 ] },
    "Left"             : { command: "KEY_LEFT"         , keycode: [ 37 ] },
    "Right"            : { command: "KEY_RIGHT"        , keycode: [ 39 ] },
    "OK (Enter)"       : { command: "KEY_OK"           , keycode: [ 13 ] },
    "Back (backspace)" : { command: "KEY_BACK"         , keycode: [  8 ] },
    "Power toggle (p)" : { command: "KEY_POWER"        , keycode: [ 80 ] },
    "Exit (Esc)"       : { command: "KEY_EXIT"         , keycode: [ 27 ] },
    "Play (space)"     : { command: "KEY_PLAY"         , keycode: [ 32 ] },
    "Volume Up (+)"    : { command: "KEY_VOLUMEUP"     , keycode: [ 107, 187, 33 ] },
    "Volume Down (-)"  : { command: "KEY_VOLUMEDOWN"   , keycode: [ 109, 189, 34 ] },
    "Mute Toggle (m)"  : { command: "KEY_MUTE"         , keycode: [ 77 ] },
    "Mute On (<)"      : { command: "KEY_MUTEON"       , keycode: [ 188 ] },
    "Mute Off (>)"     : { command: "KEY_MUTEOFF"      , keycode: [ 190 ] },
    "Key 0 (0)"        : { command: "KEY_0"            , keycode: [ 48 ] },
    "Key 1 (1)"        : { command: "KEY_1"            , keycode: [ 49 ] },
    "Key 2 (2)"        : { command: "KEY_2"            , keycode: [ 50 ] },
    "Key 3 (3)"        : { command: "KEY_3"            , keycode: [ 51 ] },
    "Key 4 (4)"        : { command: "KEY_4"            , keycode: [ 52 ] },
    "Key 5 (5)"        : { command: "KEY_5"            , keycode: [ 53 ] },
    "Key 6 (6)"        : { command: "KEY_6"            , keycode: [ 54 ] },
    "Key 7 (7)"        : { command: "KEY_7"            , keycode: [ 55 ] },
    "Key 8 (8)"        : { command: "KEY_8"            , keycode: [ 56 ] },
    "Key 9 (9)"        : { command: "KEY_9"            , keycode: [ 57 ] },
};

function sendCmd(button) {
    if ((ws == null) || (!connected)) {
        return
    }
    if ( !(button in buttons)) {
        return
    }
    ws.send(buttons[button].command)
    log("Sent command: " + buttons[button].command)

}

function stopWs(evt) {
    if (!connected) {
        return
    }
    connected = false
    removeButtons()
    log("")
    log("")
    log("Connection closed.")

    ws.close()
    ws = null
    connectWs()
}

function onMessage(evt) {
    // Message received
}

function connectWs() {
    if (ws != null) {
        return
    }
    ws = new WebSocket(wsURL);
    if (ws == null) {
        return
    }
    ws.onopen = function() {
        log("Connected.")
        addButtons()
        connected = true
    }
    ws.onmessage = onMessage
    ws.onerror   = function(evt) {
        log("Error occured: " + evt)
    }
    ws.onclose   = stopWs
}

document.addEventListener("DOMContentLoaded", connectWs)
document.addEventListener('keydown', function(event) {
    for (var i in buttons) {
        for (var b in buttons[i].keycode) {
            if (buttons[i].keycode[b] == event.keyCode) {
                sendCmd(i)
                event.preventDefault();
                return false
            }
        }
    }
    log("Unknown key: " + event.keyCode )
    return true
}, true);

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
