package plex

import "path"

//
type commander interface {
    command(args ...string) (string, error)
}

type plainCommand struct {
    Path string
}

func (p plainCommand) command(args ...string) (string, error) {
    // validate path?
    return path.Clean("/" + p.Path), nil
}
