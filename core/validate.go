package core

import "fmt"
import "regexp"


var regexptr *regexp.Regexp

func ValidateName(name string) error {
    if (regexptr == nil) {
        regexptr = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
    }
    if (!regexptr.MatchString(name)) {
        return fmt.Errorf("name '%s' can only contain A-Z, a-z, 0-9, '_' and '-' characters!", name)
    }
    return nil
}


