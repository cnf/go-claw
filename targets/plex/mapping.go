package plex

import "path"

type Commander interface {
    Command(args ...string) (string, error)
}

type PlainCommand struct {
    Path string
}

func (self PlainCommand) Command(args ...string) (string, error) {
    // validate path?
    return path.Clean("/" + self.Path), nil
}
