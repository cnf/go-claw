package dispatcher

type RemoteListener interface {
    RunListener(cmd *CommandStream)
}
