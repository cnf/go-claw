package targets

type Targets interface {
    SendCommand(cmd string, args string) bool
}
