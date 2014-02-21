package structures

// import "github.com/cnf/progrem/listeners"

type RemoteCommand struct {
    code   string
    repeat int
    key    string
    source   string
}

type CommandStream struct {
    ch chan *RemoteCommand
    err error
}
