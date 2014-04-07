package targets

import "testing"
import "fmt"

var paramstr = `
{
    "commands": {
        "VolDown": {
            "description": "Volume Down",
            "parameters": []
        },
        "VolUp": {
            "description": "Volume Up",
            "parameters": []
        },
        "SetVolume": {
            "description": "Set the volume level",
            "parameters": [
                {
                    "name": "volumelevel",
                    "description": "the volume level in percentage to set",
                    "type": "range",
                    "validation": "0:60",
                    "required": true
                }
            ]
        },
        "SelectOutput": {
            "description": "selects the specified input",
            "parameters": [
                {
                    "name": "outputname",
                    "description": "the name of the input to select",
                    "type": "list",
                    "validation": "hdmi1|hdmi2|hdmi3|hdmi4|composite1|composite2",
                    "required": true
                }
            ]
        }
    }
}
`

func printCommands(cmds map[string]*Command) {
    for k, v := range cmds {
        fmt.Printf("Command: '%s':\n", k)
        fmt.Printf("    Name       : '%s'\n", v.Name)
        fmt.Printf("    Description: '%s'\n", v.Description)

        for _, p := range v.Parameters {
            fmt.Printf("    Parameter: %s\n", p.Name)
            fmt.Printf("       Description: '%s'\n", p.Description)
            fmt.Printf("       Optional   : '%v'\n", p.Optional)
            fmt.Printf("       Type       : '%s'\n", p.Type)
            fmt.Printf("       Validation : '%s'\n", p.Validation)
        }
    }
}

func Test_CommandList(t *testing.T) {
    cmds := map[string]*Command {
        "Test": NewCommand("Test Command",
                NewParameter("Amount", "Test parameter").SetRange(0,100).SetOptional(),
            ),
        "Test2": NewCommand("Another test command",
                NewParameter("prm1", "blabla").SetList("option1", "option2"),
                NewParameter("prm2", "blabla").SetRegex("[a-z]*"),
            ),
    }
    printCommands(cmds)
}

func Test_CommandParser(t *testing.T) {
    cmds, err := ParseCommands(paramstr)
    if err != nil {
        t.Errorf("Got error: %s", err.Error())
    }
    printCommands(cmds)
}

func testValidation(t *testing.T, ptype string, validstr, val, expect string, expecterr bool) {
    cmd := CommandParameter{
        Name: "validtest",
        Description: "parameter validation",
        Type: ptype,
        Validation: validstr,
    }
    v, err := cmd.Validate(val)
    if expecterr {
        if err == nil {
            t.Errorf("Expected error for %s '%s' value '%s'", ptype, validstr, val)
        }
        return
    } 
    if err != nil {
        t.Errorf("Got an unexpected error for %s '%s' value '%s': %s", ptype, validstr, val, err.Error())
        return
    }
    if (v != expect) {
        t.Errorf("Expected value '%s' for value '%s' in %s '%s' - got '%s'", expect, val, ptype, validstr, v)
    }
}

func Test_Regex(t *testing.T) {
    //testValidation(t, paramRegex, ".*", "test1", "test1", false)
}

func Test_List(t *testing.T) {
    testValidation(t, "list", "test1", "test1", "test1", false)
    testValidation(t, "list", "test1|test2", "test1", "test1", false)
    testValidation(t, "list", "test1|test2", "test2", "test2", false)
    testValidation(t, "list", "test1|test2|bla", "test2", "test2", false)
    testValidation(t, "list", "test1|test2|bla", "bla", "bla", false)
    
    // Whe should always get the value specified in the list, original case is ignored
    testValidation(t, "list", "TeSt1|test2|bla", "test1", "TeSt1", false)

    // Test error case
    testValidation(t, "list", "test1|test2", "test3", "test3", true)
}

func Test_String(t *testing.T) {
    testValidation(t, "string", "", "bla", "bla", false)
}

func Test_Numeric(t *testing.T) {
    testValidation(t, "numeric", "", "10", "10", false)
    testValidation(t, "numeric", "", "0x5", "5", false)
    testValidation(t, "numeric", "", "abc", "", true)

    // Empty string should return a "0" value
    testValidation(t, "numeric", "", "", "0", false)
}

func Test_Range(t *testing.T) {
    // Test valid ranges
    testValidation(t, "range", "0:10", "0", "0", false)
    testValidation(t, "range", "0:10", "5", "5", false)
    testValidation(t, "range", "0:10", "10", "10", false)
    // Test values outside of the specified range
    testValidation(t, "range", "0:10", "-1", "", true)
    testValidation(t, "range", "0:10", "11", "", true)

    // Test open-ended ranges
    testValidation(t, "range", "0:", "11", "11", false)
    testValidation(t, "range", "0:", "-1", "", true)

    testValidation(t, "range", ":10", "11", "", true)
    testValidation(t, "range", ":10", "10", "10", false)
    testValidation(t, "range", ":10", "1", "1", false)
    testValidation(t, "range", ":10", "-1", "-1", false)

    // Test hexadecimal values
    testValidation(t, "range", "0x0:0x10", "0x05", "5", false)

    // Test % notation
    testValidation(t, "range", "0:200", "20%", "40", false)
    testValidation(t, "range", "-10:10", "20%", "-6", false)
    testValidation(t, "range", "10:210", "20%", "50", false)
    testValidation(t, "range", "0:10", "100%", "10", false)
    testValidation(t, "range", "0:10", "0%", "0", false)

    // These should fail
    testValidation(t, "range", "100:10", "100%", "", true)
    testValidation(t, "range", "0:10", "101%", "", true)
    testValidation(t, "range", "0:10", "-1%", "", true)
}

func Test_Rangedef(t *testing.T) {
    // These should return an error on the range definition
    testValidation(t, "range", "10-123", "0", "0", true)
    testValidation(t, "range", "10", "0", "0", true)
    testValidation(t, "range", "", "0", "0", true)
    testValidation(t, "range", "+10", "0", "0", true)
    testValidation(t, "range", "0:10", "5", "5", false)
}

