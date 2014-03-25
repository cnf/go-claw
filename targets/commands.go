package targets

import "errors"
import "strconv"
import "strings"
import "regexp"
import "encoding/json"

type commandwrapper struct {
    Commands map[string]Command
}
type Command struct {
    Name string                   `json:name`
    Description string            `json:description`
    Parameters []CommandParameter `json:parameters`
}

// ParameterType indicates the type of the parameter internally in the 
type parameterType int
// Type to validate a value with the given validation string
type ParameterValidator func(value, validation string) (string, error)

type CommandParameterType interface {
    Name() string
    Description() string
    Validate(value string) (string, error)
    Required() bool
}
type CommandParameter struct {
    NameStr string         `json:name`
    DescriptionStr string  `json:description`
    Paramtype string       `json:type`
    Validation string      `json:validation`
    Optional bool          `json:optional`

    validation_fnc ParameterValidator
}

func ParseCommands(jsonstr string) (map[string]Command, error) {
    cmds := new(commandwrapper)

    err := json.Unmarshal([]byte(jsonstr), &cmds)
    if err != nil {
        return nil, err
    }
    for k, v := range cmds.Commands {
        v.Name = k
    }
    return cmds.Commands, nil
}

func (c *CommandParameter) Name() string {
    return c.NameStr
}
func (c *CommandParameter) Description() string {
    return c.DescriptionStr
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


func (c *CommandParameter) Validate(value string) (string, error) {
    var valfnc ParameterValidator

    if c.validation_fnc == nil {
        switch c.Paramtype {
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
        return "", errors.New("Invalid validator function given for parameter " + c.NameStr)
    }
    return valfnc(value, c.Validation)
}

