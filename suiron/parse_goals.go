package suiron

// parse_goals
//
// This file contains functions which are used to parse goals:
//
//    identifyInfix(runestring []rune) (int, int)
//    getLeftAndRight(runes []rune, index int, size int) (Unifiable, Unifiable)
//    splitComplexTerm(comp []rune, index1 int, index2 int) (string, string)
//    ParseSubgoal(subgoal string) (Goal, error)
//
// Cleve Lendon

import (
    "strings"
    "fmt"
)

// identifyInfix - Determines whether the given string contains an infix.
// If it does, returns the type and the index. For example,
//    $X < 6
// ...contains the LESS_THAN infix, index 3.
//
// Params: string to parse (runes)
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

// getLeftAndRight - This function is used to parse built-in predicates,
// which are represented with an infix, such as "$X = verb" or "$X <= 47".
// It separates the two terms. If there is an error in parsing a term,
// the function throws a panic.
// Params: string to parse (runes)
//         index of infix
//         size of infix
// Return: term1, term2
func getLeftAndRight(runes []rune, index int, size int) (Unifiable, Unifiable) {
   arg1 := runes[0: index]
   arg2 := runes[index + size:]
   term1, err := parseTerm(string(arg1))
   if err != nil { panic(err.Error()) }
   term2, err := parseTerm(string(arg2))
   if err != nil { panic(err.Error()) }
   return term1, term2
} // getLeftAndRight

// splitComplexTerm - splits a string representation (runes) of a complex
// term into its functor and terms. For example, if the complex term is:
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
// This function parses all subgoals. It returns a goal object, and an error.
//
// The Not and Time operators are dealt with first, because they enclose subgoals.
// Eg.
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
            term1, term2 := getLeftAndRight(r, index, 1)
            return Unify(term1, term2), nil
        }
        if infix == LESS_THAN {
            term1, term2 := getLeftAndRight(r, index, 2)
            return LessThan(term1, term2), nil
        }
        if infix == LESS_THAN_OR_EQUAL {
            term1, term2 := getLeftAndRight(r, index, 2)
            return LessThanOrEqual(term1, term2), nil
        }
        if infix == GREATER_THAN {
            term1, term2 := getLeftAndRight(r, index, 2)
            return GreaterThan(term1, term2), nil
        }
        if infix == GREATER_THAN_OR_EQUAL {
            term1, term2 := getLeftAndRight(r, index, 2)
            return GreaterThanOrEqual(term1, term2), nil
        }
        if infix == EQUAL {
            term1, term2 := getLeftAndRight(r, index, 2)
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

    if strFunctor == "functor" {
        args, err := parseArguments(strArgs)
        if err != nil { return nil, err }
        return Functor(args...), nil
    }

    // Create a complex term.
    return parseFunctorTerms(strFunctor, strArgs)

} // ParseSubgoal
