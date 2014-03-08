package listeners

type RemoteListener interface {
    RunListener(cmd *CommandStream)
}
