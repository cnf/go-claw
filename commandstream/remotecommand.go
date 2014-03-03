package commandstream

import "time"

type RemoteCommand struct {
    Code    string
    Repeat  int
    Key     string
    Source  string
    Time    time.Time
}
