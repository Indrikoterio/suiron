package suiron

// parsegoals
//
//
// Cleve Lendon

import (
    "strings"
    "fmt"
)

// specialIndexOf - searches for a string of runes within
// a larger string of runes, and returns the index if found,
// but not if the runes are found between quotes or parentheses.
// For example, if the source string has:
//    print("Is x > 7?")
// ...then searching for '>' should return -1.
//
// Params: larger string of runes
//         runes to find
// Return: index if found, or -1
//
func specialIndexOf(runestring []rune, find []rune) int {
    length := len(runestring)
    for i := 0; i < length; i++ {
        c1 := runestring[i]
        if c1 == '"' {
            for j := i + 1; j < length; j++ {
                c2 := runestring[j]
                if c2 == '"' {
                    i = j
                    break
                }
            }
        } else if c1 == '(' {
            for j := i + 1; j < length; j++ {
                c2 := runestring[j]
                if c2 == ')' {
                    i = j
                    break
                }
            }
        } else {
            if foundAtIndex(runestring, i, find) { return i }
        }
    } // for
    return -1
} // specialIndexOf

// foundAtIndex - returns true if the runes in 'match' are
// found within the runestring at the given index.
// Params: runestring
//         index into runestring
//         runes to match
// Return: true or false
func foundAtIndex(runestring []rune, index int, match []rune) bool {
    length := len(match)
    if index + length > len(runestring) { return false }
    for i := 0; i < length; i++ {
        if runestring[index + i] != match[i] { return false }
    }
    return true
} // foundAtIndex

// identifyInfix - Used for parsing. Determines whether the given
// string contains an infix. If it does, returns the type and the
// index.
//    $X < 6
// ...contains the LESS_THAN infix, index 3.
//
// Params: string to parse
// Return: identifier (int)
//         index
//
func identifyInfix(runestring []rune) (int, int) {

    length := len(runestring)

    for i := 0; i < length; i++ {
        c1 := runestring[i]
        if c1 == '"' {
            for j := i + 1; j < length; j++ {
                c2 := runestring[j]
                if c2 == '"' {
                    i = j
                    break
                }
            }
        } else if c1 == '(' {
            for j := i + 1; j < length; j++ {
                c2 := runestring[j]
                if c2 == ')' {
                    i = j
                    break
                }
            }
        } else {
            // Can't be last character.
            if i == (length - 1) { break }
            if c1 == '<' {
                c2 := runestring[i+1]
                if c2 == '=' { return LESS_THAN_OR_EQUAL, i }
                return LESS_THAN, i
            }
            if c1 == '>' {
                c2 := runestring[i+1]
                if c2 == '=' { return GREATER_THAN_OR_EQUAL, i }
                return GREATER_THAN, i
            }
            if c1 == '=' {
                c2 := runestring[i+1]
                if c2 == '=' { return EQUAL, i }
                return UNIFY, i
            }
        }
    } // for

    return NONE, -1  // failed to find infix

} // identifyInfix

// separateTwoTerms - This function is used to parse built-in predicates,
// which are represented with an infix, such as "$X = verb" or "$X <= 47".
// It separates the two terms. If there is an error in parsing a term,
// the function throws a panic.
// Params: string to parse (runes)
//         index of infix
//         size of infix
// Return: term1, term2
func separateTwoTerms(runes []rune, index int, size int) (Unifiable, Unifiable) {
   arg1 := runes[0: index]
   arg2 := runes[index + size:]
   term1, err := parseTerm(string(arg1))
   if err != nil { panic(err.Error()) }
   term2, err := parseTerm(string(arg2))
   if err != nil { panic(err.Error()) }
   return term1, term2
} // separateTwoTerms

// splitComplexTerm - splits a string representation of a complex
// term into its functor and terms. For example, if the complex
// term is:
//
//    "father(Philip, Alize)"
//
// and the indices (index1, index2) are 6 and 20, the function will
// return: "father", "Philip, Alize"
//
// This method assumes that index1 and index2 are valid.
//
// Params: complex term (string)
//         index1
//         index2
// Return: functor (string)
//         terms   (string)
//
func splitComplexTerm(comp []rune, index1 int, index2 int) (string, string) {

      functor := comp[0: index1]
      terms   := comp[index1 + 1: index2]
      return string(functor), string(terms)

} // splitComplexTerm


// indicesOfParentheses - if a string has parentheses, this function
// will return their indices. If there are no parentheses, the indices
// will be -1.
//
// Params: chars (runes)
// Return: index of left parenthesis  (, or -1
//         index of right parenthesis ), or -1
//         error
func indicesOfParentheses(chars []rune) (int, int, error) {

    first  := -1  // index of first parenthesis
    second := -1
    countLeft  := 0
    countRight := 0

    for i, ch := range chars {
        if ch == '(' {
            if first == -1 { first = i }
            countLeft++
        } else if ch == ')' {
            second = i
            countRight++
        }
    }

    if second < first {
        s := string(chars)
        err := fmt.Errorf("indicesOfParentheses() - Invalid parentheses: %v", s)
        return first, second, err
    }

    if countLeft != countRight {
        s := string(chars)
        err := fmt.Errorf("indicesOfParentheses() - Unbalanced parentheses: %v", s)
        return first, second, err
    }

    return first, second, nil

} // indicesOfParentheses

// ParseSubgoal
//
// This function accepts a string which represents a subgoal, and
// creates its corresponding Goal object. It parses complex terms,
// the Unify operator (=), the Cut (!), and others.
//
// The Not and Time operators are dealt with first, because they
// enclose subgoals. Eg.
//
//    not($X = $Y)
//    time(qsort)
//
// Params: subgoal as string
// Return: subgoal as Goal object
//         error
//
func ParseSubgoal(subgoal string) (Goal, error) {

    s := strings.TrimSpace(subgoal)
    r := []rune(s)
    length := len(r)

    if length == 0 {
        err := fmt.Errorf("ParseSubgoal() - Empty string.")
        return nil, err
    }

    if length > 5 {
        last := r[length - 1]
        if last == ')' {
            // If the string starts with not(
            if strings.HasPrefix(s, "not(") {
//                s2 := s[4: length - 1]
//                return Not(subgoal(s2))
            }
            // If the string starts with time(
        }
    } // if length > 5

    if s == "!" {  // cut
        return Cut(), nil
    } else if s == "fail" {
        return Fail(), nil
    } else if s == "nl" {
        return NL(), nil
    }

    //--------------------------------------
    // Handle infixes: = > < >= <= == =

    infix, index := identifyInfix(r)
    if infix != NONE {
        if infix == UNIFY {
            term1, term2 := separateTwoTerms(r, index, 1)
            return Unify(term1, term2), nil
        }
        if infix == LESS_THAN {
            term1, term2 := separateTwoTerms(r, index, 2)
            return LessThan(term1, term2), nil
        }
        if infix == LESS_THAN_OR_EQUAL {
            term1, term2 := separateTwoTerms(r, index, 2)
            return LessThanOrEqual(term1, term2), nil
        }
        if infix == GREATER_THAN {
            term1, term2 := separateTwoTerms(r, index, 2)
            return GreaterThan(term1, term2), nil
        }
        if infix == GREATER_THAN_OR_EQUAL {
            term1, term2 := separateTwoTerms(r, index, 2)
            return GreaterThanOrEqual(term1, term2), nil
        }
        if infix == EQUAL {
            term1, term2 := separateTwoTerms(r, index, 2)
            return Equal(term1, term2), nil
        }
        panic("identifyInfix() - Missing an infix?")
    }

    // Check for parentheses.
    leftIndex, rightIndex, err := indicesOfParentheses(r)
    if err != nil { return nil, err }

    if leftIndex == -1 {   // If left is -1, right is too.
        // This is OK.
        // A 'goal' can be a simple word, without parentheses.
        return parseFunctorTerms(s, "")
    }

    strFunctor, strArgs := splitComplexTerm(r, leftIndex, rightIndex)

    // Check for built-in predicates.

    if strFunctor == "append" {
        args, err := parseArguments(strArgs)
        if err != nil { return nil, err }
        return Append(args...), nil
    }

    if strFunctor == "print" {
        args, err := parseArguments(strArgs)
        if err != nil { return nil, err }
        return Print(args...), nil
    }

    if strFunctor == "time" {
        goal, err := ParseComplex(strArgs)
        if err != nil { return goal, err }
        return Time(goal), nil
    }

    // Create a complex term.
    return parseFunctorTerms(strFunctor, strArgs)

} // ParseSubgoal
