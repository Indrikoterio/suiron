package suiron

// parse_args
//
// A complex term consists of a functor followed by a list of arguments
// (terms) between parentheses: functor(arg1, arg2, arg3). These arguments
// might be Atoms, Integers, Floats, Variables or LinkedLists.
//
// This file contains functions to parse a list of arguments, and return
// a slice of unifiable terms (Atoms, Integers, etc.). Arguments are comma
// separated, but if a comma is between double quotes, the comma is included
// in the argument. For example, in the complex term:
//
//      address(London, Baker St., 221B)
//
// there are 3 arguments (arity = 3): London, 'Baker St.' and 221B.
//
// In the complex term,
//
//      address("London, Baker St., 221B")
//
// there is one argument (arity = 1): ->London, Baker St., 221B<-
//
// Another way to define a term which contains commas is to escape the
// commas with backslashes:
//
//      address(London\, Baker St.\, 221B).
//
// This complex term above has one argument, as in the previous example.
//
// Any character can be escaped with a backslash, including the double
// quote mark.
//
// Note:
// Complex terms can be parsed from a Suiron program file (a text file),
// as in the examples above, or they can be created from a Go program file.
// For example;
//
//     loves   := Atom("loves")
//     leonard := Atom("Leonard")
//     penny   := Atom("Penny")
//     c1 := Complex{loves, leonard, penny}
//
// C1 above represents the complex term: loves(Leonard, Penny)
//
// In a text file, a complex term which defines a double quote mark might
// be defined as follows:
//
//      double_quote(\")
//
// The quotation mark is escaped so that it will be interpreted as an Atom,
// not as the start of a quoted term.
//
// The backslash character is also used as an escape character in the Go
// programming language. In a Go program source file, in order to define
// the above term, both the quotation mark and the backslash would need
// to be escaped to pass through the Go compiler, as follows:
//
//      c := ParseComplex("double_quote(\\\")")
//
// Another example,
//
//      inspiration("\"Do or do not. There is no try.\" -- Yoda")
//
// In this complex term, there is one argument, which includes quote marks.
//
//      "Do or do not. There is no try." -- Yoda
//
// In a Go program file:
//
//      c := ParseComplex("\\\"Do or do not. There is no try.\\\" -- Yoda")
//
// There are 3 kinds of constants in Suiron: Atoms, Integers, Floats.
// In the following complex term,
//
//      person(Cleve Lendon, 1961, 1.78)
//
// 'Cleve Lendon' is an Atom (= string). 1961 is a 64-bit Integer, and
// 1.78 is a 64-bit Float.
//
// Any term between double quotation marks will be parsed as an Atom.
//
//      person("Cleve Lendon", "1961", "1.78")
//
// In the above complex term, all three arguments are Atoms.
//
// Variables start with a dollar sign, followed by at least one letter:
//
//      father(Alfred, $Son)
//
// Of course, the dollar sign can be escaped with a backslash, or put
// inside quotes:
//
//      dollar_sign("$")
//      dollar_sign(\$)
//
// Cleve Lendon

import (
    "strings"
    "strconv"
    "fmt"
)

// parseArguments - parses a comma separated list of arguments.
// Returns a slice of Unifiables terms and an error, if any.
func parseArguments(str string) ([]Unifiable, error) {

    s := strings.TrimSpace(str)
    r := []rune(s)
    length := len(r)

    if len(s) == 0 {
        err := makeParseError("Empty argument list", str)
        return []Unifiable{}, err
    }

    first := r[0]
    if first == ',' {
        err := makeParseError("Missing first argument", str)
        return []Unifiable{}, err
    }

    // A comma at the end probably indicates a missing argument, but...
    // make sure comma is not escaped, because this is valid: "term1, term2, \,"
    last := r[length - 1]
    if last == ',' {
        // Length must be longer than 1, because comma
        // is not the first character.
        prev := r[length - 2]
        if prev != '\\' {   // escape character
            err := makeParseError("Missing last argument", str)
            return []Unifiable{}, err
        }
    }

    hasDigit    := false
    hasNonDigit := false
    hasPeriod   := false
    openQuote   := false

    numQuotes   := 0
    roundDepth  := 0   // depth of round parentheses (())
    squareDepth := 0   // depth of square brackets [[]]

    var argument  []rune
    var arguments []Unifiable

    start := 0

    for i := start; i < length; i++ {

        ch := r[i]

        // If this argument is between double quotes,
        // it must be an Atom.
        if openQuote {
            argument = append(argument, ch)
            if ch == '"' {
                openQuote = false
                numQuotes += 1
            }
        } else {
            if ch == '[' {
                argument = append(argument, ch)
                squareDepth++
            } else if ch == ']' {
                argument = append(argument, ch)
                squareDepth--
            } else if ch == '(' {
                argument = append(argument, ch)
                roundDepth++
            } else if ch == ')' {
                argument = append(argument, ch)
                roundDepth--
            } else if roundDepth == 0 && squareDepth == 0 {
                if ch == ',' {
                    s2 := strings.TrimSpace(string(argument))
                    err := checkQuotes(s2, numQuotes)
                    if err != nil { return []Unifiable{}, err }
                    numQuotes   = 0
                    c, err := makeTerm(s2, hasDigit, hasNonDigit, hasPeriod)
                    if err != nil { return []Unifiable{}, err }
                    arguments   = append(arguments, c)
                    argument    = nil
                    hasDigit    = false
                    hasNonDigit = false
                    hasPeriod   = false
                    start = i + 1    // past comma
                } else if ch >= '0' && ch <= '9' {
                    argument = append(argument, ch)
                    hasDigit = true
                } else if ch == '.' {
                    argument = append(argument, ch)
                    hasPeriod = true
                } else if ch == '\\' {  // escape character, must include next character
                    if i < length - 1 {
                        i += 1
                        argument = append(argument, r[i])
                    } else {  // must be at end of argument string
                        argument = append(argument, ch)
                    }
                } else if ch == '"' {
                    argument = append(argument, ch)
                    openQuote = true  // first quote
                    numQuotes += 1
                } else {
                    argument = append(argument, ch)
                    if ch > ' ' { hasNonDigit = true }
                }
            } else {
                // Must be between () or []. Just add character.
                argument = append(argument, ch)
            }
        } // not openQuote
    } // for

    if start < length {
        s2 := strings.TrimSpace(string(argument))
        err := checkQuotes(s2, numQuotes)
        if err != nil { return []Unifiable{}, err }
        c, err := makeTerm(s2, hasDigit, hasNonDigit, hasPeriod)
        if err != nil { return []Unifiable{}, err }
        arguments = append(arguments, c)
    }

    if roundDepth != 0 {
        err := makeParseError("Unmatched parentheses", str)
        return arguments, err
    }

    if squareDepth != 0 {
        err := makeParseError("Unmatched brackets", str)
        return arguments, err
    }

    return arguments, nil

} // parseArguments()


// makeTerm - determines whether the given string represents an integer,
// a floating point number, an atom (text string), a logic variable or
// a linked list, and returns the appropriate term.
// If the programmer makes a coding error, for example, typing $1X when
// he/she intended $X1, this function will return $1X as an Atom.
// 
//    str - string to parse
//    hasDigit    - argument has digit
//    hasNonDigit - argument has non-digit
//    hasPeriod   - avoid argument
func makeTerm(str string,
              hasDigit bool,
              hasNonDigit bool,
              hasPeriod bool) (Unifiable, error) {

    s := strings.TrimSpace(str)

    length := len(s)
    if length == 0 {
        e := makeTermError("Length of term is 0.", s)
        return Atom(s), e
    }

    first := s[0:1]
    if first == "\\" && length > 1 {
        s = s[1:]
        length = len(s)
    } else if first == "$" {

        // Anonymous variable.
        if s == "$_" { return Anon(), nil }

        logicVariable, err := LogicVar(s)
        // If the string is not a valid variable ($, $10), make it an Atom.
        if err != nil { return Atom(s), nil }
        return logicVariable, nil
    }

    // If the argument begins and ends with a quotation mark,
    // the argument is an Atom. Strip off quotation marks.
    if length >= 2 {
        last := s[length-1:]
        if first == "\"" {
            if last == "\"" {
                s2 := s[1: length -1]
                if len(s2) == 0 {
                    err := makeTermError("Invalid term. Length is 0.", str)
                    return Atom(s), err
                }
                return Atom(s2), nil
            } else {
                err := makeTermError("Invalid term. Unmatched quote mark.", str)
                return Atom(str), err
            }
        } else if first == "[" && last == "]" {
            term, err := ParseLinkedList(s)
            return term, err

        // Try complex terms, eg.:  job(programmer)
        } else if first != "(" && last == ")" {
            // Check for built-in functions.
            if strings.HasPrefix(s, "join(")     { return ParseFunction(s) }
            if strings.HasPrefix(s, "add(")      { return ParseFunction(s) }
            if strings.HasPrefix(s, "subtract(") { return ParseFunction(s) }
            if strings.HasPrefix(s, "multiply(") { return ParseFunction(s) }
            if strings.HasPrefix(s, "divide(")   { return ParseFunction(s) }
            c, err := ParseComplex(s)
            return c, err
        }
    } // length > 2

    if hasDigit && !hasNonDigit { // Must be Integer or Float.
        if hasPeriod {
            f, err := strconv.ParseFloat(s, 64)
            if err == nil { return Float(f), nil }
        } else {
            i, err := strconv.ParseInt(s, 10, 64)
            if err == nil { return Integer(i), nil }
        }
    }
    return Atom(s), nil

}  // makeTerm


// checkQuotes - checks syntax of double quote marks (") and
// produces an error if there is a problem. An argument may be
// enclosed in double quotation marks at the beginning and end,
// eg. "Sophie". Quotation marks which have been escaped with
// a backslash, (\"), are not counted.
// Arguments such as "Hello"" or "Hello are quite wrong.
func checkQuotes(str string, count int) error {
    if count == 0 { return nil }
    if count != 2 {
        return fmt.Errorf("Unmatched quotes: %v", str)
    }
    first := str[0:1]
    if first != "\"" {
        return fmt.Errorf("Text before opening quote: %v", str)
    }
    last  := str[len(str) - 1:]
    if last != "\"" {
        return fmt.Errorf("Text after closing quote: %v", str)
    }
    return nil
}

// makeParseError - creates an error for parseArguments().
// msg - error message
// str - string which caused the error
func makeParseError(msg string, str string) error {
    return fmt.Errorf("parseArguments() - %v: >%v<\n", msg, str)
}

// parseTerm - determines whether the given string represents a floating
// point number, an integer, a logic variable, etc., and returns the
// appropriate term. If the programmer makes a coding error, for example,
// typing $1X when he/she intended $X1, the function will return $1X as
// an Atom (= string).
//
// Params:
//    str - string to parse
// Return:
//    unifiable term and error
func parseTerm(str string) (Unifiable, error) {

    s := strings.TrimSpace(str)

    hasDigit    := false
    hasNonDigit := false
    hasPeriod   := false

    for _, ch := range s {
        if ch >= '0' && ch <= '9' {
            hasDigit = true
        } else if ch == '.' {
            hasPeriod = true
        } else {
            hasNonDigit = true
        }
    }
    return makeTerm(str, hasDigit, hasNonDigit, hasPeriod)

}  // parseTerm

// makeTermError - creates an error for makeTerm().
// msg - error message
// str - string which caused the error
func makeTermError(msg string, str string) error {
    return fmt.Errorf("makeTerm() - %v: >%v<", msg, str)
}
