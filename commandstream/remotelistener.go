package commandstream

type RemoteListener interface {
    RunListener(cmd *CommandStream)
}
