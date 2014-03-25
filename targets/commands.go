package targets

import "errors"
import "strconv"
import "strings"
import "regexp"
import "encoding/json"
//import "fmt"

// Type to validate a value with the given validation string
type ParameterValidator func(value, validation string) (string, error)

type Command struct {
    Name string                   `json"-"`
    Description string            `json:"description"`
    Parameters []*CommandParameter `json:"parameters"`
}

type CommandParameter struct {
    Name string         `json:"name"`
    Description string  `json:"description"`
    Type string         `json:"type"`
    Validation string   `json:"validation"`
    Optional bool       `json:"optional"`

    validation_fnc ParameterValidator
}

func NewCommand(name, desc string, param... *CommandParameter) *Command {
    ret := new(Command)
    ret.Name = name
    ret.Description = desc
    ret.Parameters = param
    return ret
}

func NewParameter(name, desc string, optional bool) *CommandParameter {
    ret := new(CommandParameter)
    ret.Name = name
    ret.Description = desc
    ret.Optional = optional
    ret.Type = "empty"
    ret.Validation = ""
    ret.validation_fnc = nil

    return ret
}

func (c *CommandParameter) Validate(value string) (string, error) {
    var valfnc ParameterValidator

    if c.validation_fnc == nil {
        switch c.Type {
            case "string":
                valfnc = validateString
            case "regex":
                valfnc = validateRegex
            case "numeric":
                valfnc = validateNumeric
            case "range":
                valfnc = validateRange
            case "list":
                valfnc = validateList
            case "custom":
                return "", errors.New("CommandParameter:Validate(): internal error: Custom parameter defined but no function specified!")
            default:
                return value, nil
        }
    } else {
        // If a validation function is defined 
        valfnc = c.validation_fnc
    }
    if valfnc == nil {
        return "", errors.New("Invalid validator function given for parameter " + c.Name)
    }
    return valfnc(value, c.Validation)
}

func (c *CommandParameter) SetRange(start, end int) *CommandParameter {
    c.Type = "range"
    c.Validation = strconv.Itoa(start) + ":" + strconv.Itoa(end)
    c.validation_fnc = validateRange
    return c
}

func (c *CommandParameter) SetList(list... string) *CommandParameter {
    strlist := ""
    for _, s := range list {
        if (len(strlist) != 0) {
            strlist += "|"
        }
        strlist += s
    }
    c.Type = "list"
    c.Validation = strlist
    c.validation_fnc = validateList

    return c
}

func (c *CommandParameter) SetString() *CommandParameter {
    c.Type = "string"
    c.Validation = ""
    c.validation_fnc = validateString

    return c
}
func (c *CommandParameter) SetNumeric() *CommandParameter {
    c.Type = "numeric"
    c.Validation = ""
    c.validation_fnc = validateNumeric

    return c
}
func (c *CommandParameter) SetRegex(regex string) *CommandParameter {
    c.Type = "regex"
    c.Validation = regex
    c.validation_fnc = validateRegex

    return c
}

func (c *CommandParameter) SetCustom(validation string, fnc ParameterValidator) *CommandParameter {
    c.Type = "custom"
    c.Validation = validation
    c.validation_fnc = fnc

    return c
}

func validateString(value, validation string) (string, error) {
    return value, nil
}

func validateRegex(value, validation string) (string, error) {
    if validation == "" {
        return value, nil
    }
    ok, err := regexp.MatchString(validation, value)
    if err != nil {
        return "", err
    }
    if !ok {
        return "", errors.New("Value '" + value + "' did not match regex '" + validation + "'")
    }
    return value, nil
}


func validateNumeric(value, validation string) (string, error) {
    if value == "" {
        return "0", nil
    }
    base := 0

    if validation != "" {
        base_try, err := strconv.ParseInt(validation, 0, 0)
        if err == nil {
            // Valid base specified
            base = int(base_try)
        }
    }
    val_try, err := strconv.ParseInt(value, base, 0)
    if err != nil {
        return "", errors.New("value '" + value + "' is an invalid number: " + err.Error())
    }
    // Return coverted to base 10
    return strconv.Itoa(int(val_try)), nil
}


func validateRange(value, validation string) (string, error) {
    var ispct = false
    var err error
    if value[len(value) - 1] == '%' {
        ispct = true
        value = value[0:len(value)-1]
    }
    value, err = validateNumeric(value, "")
    if err != nil {
        return "", err
    }
    ival, _ := strconv.ParseInt(value, 0, 0)
    
    // Split the validation string
    ranges := strings.Split(validation, ":")
    if len(ranges) != 2 {
        return "", errors.New("invalid validation string for range! Expected 'start:end' format")
    }
    var lval, uval int64
    if ranges[0] != "" {
        lval, err = strconv.ParseInt(ranges[0], 0, 0)
        if err != nil {
            return "", err
        }
        if (!ispct) && (ival < lval) {
            return "", errors.New("value " + value + " too small for range " + validation)
        }
        // % notation assumes 0 as starting point
    }
    if ranges[1] != "" {
        uval, err = strconv.ParseInt(ranges[1], 0, 0)
        if err != nil {
            return "", err
        }
        if (uval < lval) {
            return "", errors.New("range validation error: upper value " + ranges[1] + " < lower value " + ranges[0])
        }
        if (!ispct) && (ival > uval) {
            return "", errors.New("value " + value + " too big for range " + validation)
        }
    } else if (ispct) {
        return "", errors.New("Cannot use % notation for range when no upper bound is specified")
    }


    if (ispct) {
        if (ival < 0) || (ival > 100) {
            return "", errors.New("percentage value has to be in the range 0-100")
        }
        rnge := uval - lval
        return strconv.Itoa( int(float64(lval) + ((float64(rnge) / 100.0) * float64(ival)) ) ), nil
    }
    return value, nil
}

func validateList(value, validation string) (string, error) {
    vallist := strings.Split(validation, "|")
    for _, s := range vallist {
        if (strings.ToUpper(value) == strings.ToUpper(s)) {
            return s, nil
        }
    }
    return "", errors.New("Value '" + value + "' not in " + validation)
}


///////////
type parseStruct struct {
    Commands map[string]*Command `json:"commands"`
}

func ParseCommands(jsonstr string) (map[string]*Command, error) {
    var cmds parseStruct 

    err := json.Unmarshal([]byte(jsonstr), &cmds)
    if err != nil {
        return nil, err
    }
    for v , c:= range cmds.Commands {
        c.Name = v
        cmds.Commands[v] = c
    }
    return cmds.Commands, nil
}

