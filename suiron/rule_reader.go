package suiron

// RuleReader - reads Suiron facts and rules from a file.
//
// Cleve Lendon

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// separateRules - divides a text string to create an array of rules.
// Each rule or fact ends with a period.
//
// Params:
//    original text
// Return
//    rules (slice)
//    error
func separateRules(text string) ([]string, error) {

    var str string
    var sb strings.Builder
    rules := []string{}

    roundDepth  := 0
    squareDepth := 0
    quoteDepth  := 0

    for _, ch := range text {
        sb.WriteRune(ch)
        if ch == '.' && roundDepth == 0 &&
           squareDepth == 0 && quoteDepth % 2 == 0 {
            str = sb.String()
            rules = append(rules, str)
            sb.Reset()
        } else if ch == '(' {
            roundDepth++
        } else if ch == '[' {
            squareDepth++
        } else if ch == ')' {
            roundDepth--
        } else if ch == ']' {
            squareDepth--
        } else if ch == '"' {
            quoteDepth++
        }
    } // for

    // Check for unmatched brackets here.
    err := unmatchedBracket(str, roundDepth, squareDepth)

    return rules, err

} // separateRules


// unmatchedBracket - checks for an unmatched bracket. If there is an
// unmatched bracket, returns an error. Otherwise returns nil.
// If there is no error, return nil.
// Params:  previous string, for context
//          depth of round brackets (int)
//          depth of square brackets (int)
// Return:  error or nil
func unmatchedBracket(str string, roundDepth int, squareDepth int) error {

    // If no error, return nil.
    if roundDepth == 0 && squareDepth == 0 { return nil }

    var msg  string
    var msg2 string

    if roundDepth > 0 {
        msg = "unmatched parenthesis: (\n"
    } else if roundDepth < 0 {
        msg = "unmatched parenthesis: )\n"
    } else if squareDepth > 0 {
        msg = "unmatched bracket: [\n"
    } else if squareDepth < 0 {
        msg = "unmatched bracket: ]\n"
    }

    if len(str) == 0 {
        msg2 = "Check start of file."
    } else {
        if len(str) > 60 {
            str = str[0: 60]
        }
        msg2 = "Error occurs after: " + str
    }
    return fmt.Errorf("Error - " + msg + msg2)

} // unmatchedBracket


// stripComments - strips comments from a line. In Suiron, valid
// comment delimiters are:
//
//   %  Comment
//   #  Comment
//   // Comment
//
// Any text which occurs after these delimiters is considered
// a comment, and removed from the line. But, if these delimiters
// occur within braces, they are not treated as comment delimiters.
// For example, in the line
//
//    print(Your rank is %s., $Rank),   % Print user's rank.
//
// the first percent sign does not start a comment, but the second
// one does.
//
// Params:
//    original line
// Return:
//    trimmed line
//
func stripComments(line string) string {

    previous    := 'x'
    roundDepth  := 0
    squareDepth := 0

    index := -1
    for i, ch := range line {
        if ch == '(' {
            roundDepth++
        } else if ch == '[' {
            squareDepth++
        } else if ch == ')' {
            roundDepth--
        } else if ch == ']' {
            squareDepth--
        } else if roundDepth == 0 && squareDepth == 0 {
            if ch == '#' || ch == '%' {
                index = i
                break
            } else if ch == '/' && previous == '/' {
               index = i - 1
               break
            }
        }
        previous = ch
    } // for

    if index >= 0 {
        return strings.TrimSpace(line[0: index])
    } else {
        return strings.TrimSpace(line)
    }

}  // stripComments


// ReadFactsAndRules - reads Suiron facts and rules from a text file.
// Strips out all comments. (Comments are preceded by #, % or // .)
// Param:  file name
// Return: array (slice) of rules
//         error
func ReadFactsAndRules(fileName string) ([]string, error) {

    file, err := os.Open(fileName)
    if err != nil { return []string{}, err }
    defer file.Close()

    var sb strings.Builder
    scanner := bufio.NewScanner(file)

    lineNum := 1
    for scanner.Scan() {
        aLine := scanner.Text()
        strippedLine := stripComments(aLine)
        if len(strippedLine) > 0 {
            err := checkEndOfLine(strippedLine, lineNum)
            if err != nil { return []string{}, err }
            sb.WriteString(strippedLine)
            sb.WriteString(" ")
        }
        lineNum++
    }

    roolz, err := separateRules(sb.String())
    return roolz, err

} // ReadFactsAndRules


// StringToRules - Divides string into an array of facts and rules.
// Strips out all comments. (Comments are preceded by #, % or // .)
//
// Param:  string
// Return: array (slice) of facts and rules
//         error
func StringToRules(str string) ([]string, error) {

    var sb strings.Builder
    scanner := bufio.NewScanner(strings.NewReader(str))

    lineNum := 1
    for scanner.Scan() {
        aLine := scanner.Text()
        strippedLine := stripComments(aLine)
        if len(strippedLine) > 0 {
            err := checkEndOfLine(strippedLine, lineNum)
            if err != nil { return []string{}, err }
            sb.WriteString(strippedLine)
            sb.WriteString(" ")
        }
        lineNum++
    }

    roolz, err := separateRules(sb.String())
    return roolz, err

} // StringToRules


// checkEndOfLine - checks to ensure that a line read from a file
// ends with a valid character. Why?
// Rules and facts can be split over several lines. For example,
// it is valid to write a rule as:
//
//   parse($In, $Out) :-
//         words_to_pos($In, $In2),
//         remove_punc($In2, $In3),
//         sentence($In3, $Out).
//
// The lines above end in dash, comma, comma and period.
// If a line were a simple word, such as:
//
//   sentence
//
// That would indicate an error in the source.
//
// Currently, valid characters are dash, comma, semicolon,
// period and the equal sign.
//
// Params: line of text to check
//         number of line to check
// Return: error or nil
func checkEndOfLine(line string, num int) error {
    length := len(line)
    if length > 0 {
        last := line[length - 1:]
        if last != "-" && last != "," &&
           last != ";" && last != "." && last != "=" {
            return fmt.Errorf("Check line %d: %v", num, line)
        }
    }
    return nil
}

// LoadKBFromFile - reads rules and facts from a text file, parses
// them, then adds them to the knowledge base. If a parsing error is
// generated, add the previous line to the error message.
//
// Params:  knowledge base
//          filename
// Return:  error or nil
//
func LoadKBFromFile(kb KnowledgeBase, fileName string) error {
    factsAndRules, err := ReadFactsAndRules(fileName)
    if err != nil { return err }
    var previous string
    for _, str := range factsAndRules {
        factOrRule, err := ParseRule(str)
        if err != nil {
            return LoadParseError(previous, err)
        }
        kb.Add(factOrRule)
        previous = str
    }
    return nil
} // LoadKBFromFile

// LoadParseError - If a parse error occurs while loading rules,
// this function adds the previous line for context.
// Params: previous line
//         parsing error
// Return: new error
func LoadParseError(previous string, err error) error {
    strError := err.Error()
    if len(previous) == 0 {
        return fmt.Errorf(strError + "Check start of file.")
    } else {
        return fmt.Errorf(strError + "Error occurs after: " + previous)
    } 
} // LoadParseError
