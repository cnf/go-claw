package core

import "strings"
import "unicode"

// Split a string containing quoted strings on newlines, quotes, ... 
// Supports escaping of space, newline, ...
func SplitQuoted(s string) []string {
    var ret []string
    var curr = make([]rune, len(s))
    var cpos = 0
    var quoted = ' '
    var escaped = false
    sr := strings.NewReader(s)

    for {
        r, _, err := sr.ReadRune()
        if err != nil {
            // Append last
            if cpos != 0 {
                ret = append(ret, string(curr[0:cpos]))
                cpos = 0
            }
            break
        }
        switch r {
        case ' ':
            if quoted != ' ' {
                if escaped {
                    curr[cpos] = '\\'
                    cpos++
                    escaped = false
                }
                curr[cpos] = ' '
                cpos++
            } else if escaped {
                curr[cpos] = ' '
                cpos++
                escaped = false
            } else if cpos != 0 {
                ret = append(ret, string(curr[0:cpos]))
                cpos = 0
            }
        case '"', '\'':
            if escaped {
                curr[cpos] = r
                cpos++
                escaped = false
            } else if quoted == r {
                // Quoted string closed
                // Don't add to list yet, whitespace should follow if 
                // it's a new string/parameter, otherwise treat it as
                // the same
                //ret = append(ret, string(curr[0:cpos]))
                //cpos = 0
                quoted = ' '
            } else if quoted == ' ' {
                // New quote, start new entry
                quoted = r
            } else {
                curr[cpos] = r
                cpos++
            }
        case '\\':
            if escaped {
                curr[cpos] = '\\'
                cpos++
                escaped = false
            } else {
                escaped = true
            }
        default:
            if unicode.IsSpace(r) {
                // the other white space - cannot escape this!
                if quoted != ' ' {
                    curr[cpos] = r
                    cpos++
                } else if cpos != 0 {
                    // Add to lst
                    ret = append(ret, string(curr[0:cpos]))
                    cpos = 0
                }
            } else {
                if escaped {
                    curr[cpos] = '\\'
                    cpos++
                    escaped = false
                }
                curr[cpos] = r
                cpos++
            }
        }
    }
    return ret
}

