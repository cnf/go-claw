package denon

import "fmt"
import "strconv"
import "errors"

// commanders[cmd].Command(cmd, args...)

type Commander interface {
    Command(args ...string) (string, error)
}

type PlainCommand struct {
    Send string
}

type RangeCommand struct {
    Send string
    Min, Max int
}

type VolumeCommand struct {
    Send string
    Min, Max int
}

type ToggleCommand struct {
    Query string
    OnString string
    OffString string
}

func (d PlainCommand) Command(args ...string) (string, error) {
    return d.Send, nil
}

func (d RangeCommand) Command(args ...string) (string, error) {
    vol, err := strconv.Atoi(args[0])
    if err != nil {return "", err}
    if (int(vol) >= d.Min) && (int(vol) <= d.Max) {
        return fmt.Sprintf(d.Send, vol), nil
    }
    return "", errors.New("Could not construct denon command")
}

func (d VolumeCommand) Command(args ...string) (string, error) {
    vol, err := strconv.Atoi(args[0])
    if err != nil {return "", err}
    if (int(vol) >= 0) && (int(vol) <= 100) {
        i := percentageOfRange(int(vol), d.Min, d.Max)
        return fmt.Sprintf(d.Send, i), nil
    }
    return "", errors.New("Could not construct denon command")
}

func (d ToggleCommand) Command(args ...string) (string, error) {
    // compare stuff here
    return "", errors.New("Could not toggle")
}

func percentageOfRange(pct int, min int, max int) int {
    // Return a percentage of a range
    return int(float32(pct*max + 100*min - pct*min) / 100)
}
