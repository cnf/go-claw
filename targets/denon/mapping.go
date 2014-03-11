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
    Send     string
    Min, Max int
}

type VolumeCommand struct {
    Send string
    Min, Max int
}

func (self PlainCommand) Command(args ...string) (string, error) {
    return self.Send, nil
}

func (self RangeCommand) Command(args ...string) (string, error) {
    vol, err := strconv.Atoi(args[0])
    if err != nil {return "", err}
    if (int(vol) >= self.Min) && (int(vol) <= self.Max) {
        return fmt.Sprintf(self.Send, vol), nil
    }
    return "", errors.New("Could not construct denon command")
}

func (self VolumeCommand) Command(args ...string) (string, error) {
    vol, err := strconv.Atoi(args[0])
    if err != nil {return "", err}
    if (int(vol) >= 0) && (int(vol) <= 100) {
        i := percentageOfRange(int(vol), self.Min, self.Max)
        return fmt.Sprintf(self.Send, i), nil
    }
    return "", errors.New("Could not construct denon command")
}

func percentageOfRange(pct int, min int, max int) int {
    // Return a percentage of a range
    return int(float32(pct*max + 100*min - pct*min) / 100)
}
