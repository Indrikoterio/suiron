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
// Note: An infix must be preceded and followed by a space.
// Don't like this:  $X<100
//
// Params: string to parse (runes)
// Return: identifier (int)
//         index
//
func identifyInfix(runestring []rune) (int, int) {

    length := len(runestring)
    prev   := '#'  // not a space.

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
            // Can't be first or last character.
            // Previous character must be space.
            if prev != ' ' || i >= (length - 2) {
                prev = c1
                continue
            }
            if c1 == '<' {
                c2 := runestring[i+1]
                if c2 == '=' && i < length - 2 {
                    c3 := runestring[i+2]
                    if c3 == ' ' {
                        return LESS_THAN_OR_EQUAL, i
                    }
                } else if c2 == ' ' {
                    return LESS_THAN, i
                }
            } else if c1 == '>' {
                c2 := runestring[i+1]
                if c2 == '=' && i < length - 2 {
                    c3 := runestring[i+2]
                    if c3 == ' ' {
                        return GREATER_THAN_OR_EQUAL, i
                    }
                } else if c2 == ' ' {
                    return GREATER_THAN, i
                }
            } else if c1 == '=' {
                c2 := runestring[i+1]
                if c2 == '=' &&  i < length - 2 {
                    c3 := runestring[i+2]
                    if c3 == ' ' {
                        return EQUAL, i
                    }
                } else if c2 == ' ' {
                    return UNIFY, i
                }
            }
        } // else

        prev = c1

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

    if countLeft != countRight {
        s := string(chars)
        err := fmt.Errorf("indicesOfParentheses() - Unbalanced parentheses: %v", s)
        return first, second, err
    }

    if second < first {
        s := string(chars)
        err := fmt.Errorf("indicesOfParentheses() - Invalid parentheses: %v", s)
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

    // not() looks like a built-in predicate
    // but it's actually an operator.
    if strings.HasPrefix(s, "not(") {
        s2 := s[4: length - 1]
        operand, err := ParseSubgoal(s2)
        if err != nil { return nil, err }
        return Not(operand), nil
    }

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

    if strFunctor == "time" {
        goal, err := ParseComplex(strArgs)
        if err != nil { return goal, err }
        return Time(goal), nil
    }

    args, err := parseArguments(strArgs)
    if err != nil { return nil, err }

    if strFunctor == "append"  { return Append(args...), nil }
    if strFunctor == "print"   { return Print(args...), nil }
    if strFunctor == "functor" { return Functor(args...), nil }
    if strFunctor == "include" { return Include(args...), nil }
    if strFunctor == "exclude" { return Exclude(args...), nil }
    if strFunctor == "print_list" { return PrintList(args...), nil }

    // Create a complex term.
    f := Atom(strFunctor)
    unifiables := append([]Unifiable{f}, args...)
    return Complex(unifiables), nil

} // ParseSubgoal


// ParseFunction - parses a string to produce a built-in Suiron function.
// ParseFunction is similar to ParseComplex in complex.go.
// Perhaps some consolidation could be done in future.
//
// Example of usage:
//     c := ParseFunction("add(7, 9, 4)")
//
// Params: string representation
// Return: built-in suiron function
//         error
//
func ParseFunction(str string) (Function, error) {

    s := strings.TrimSpace(str)
    length := len(s)

    if length > 1000 {
        err := fmt.Errorf("ParseFunction - String is too long: %v", s)
        return nil, err
    }

    // Get indices.
    left, right, err := indicesOfParentheses([]rune(s))
    if err != nil { return nil, err }

    functor := strings.TrimSpace(s[0: left])
    args    := strings.TrimSpace(s[left + 1: right])

    t, err := parseArguments(args)
    if err != nil { return nil, err }

    unifiables := append([]Unifiable{}, t...)

    if functor == "join" { return Join(unifiables...), nil }
    if functor == "add" { return Add(unifiables...), nil }
    if functor == "subtract" { return Subtract(unifiables...), nil }
    if functor == "multiply" { return Multiply(unifiables...), nil }
    if functor == "divide"   { return Divide(unifiables...), nil }

    err = fmt.Errorf("ParseFunction - Unknown function: %v", functor)
    return nil, err

} // ParseFunction
