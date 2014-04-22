package core

import "errors"
import "fmt"
import "strings"

func ValidateName(name string) error {
    if name == "" {
        return errors.New("a target name cannot be empty")
    }
    if strings.ContainsAny(name, "\t\n\r :@!+=*") {
        return fmt.Errorf("target name '%s' cannot contain whitespace, ':', '@', '!', '+', '=' or '*' characters")
    }
    return nil
}


