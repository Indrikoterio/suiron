package suiron

// Tokenizer - parses Suiron's facts and rules.
//
// Cleve Lendon

import (
    "strings"
    "fmt"
)

// LetterNumberHyphen - determines whether the given character (rune)
// is a letter, a number, or a hyphen. This excludes punctuation.
//
// Param:  character (rune)
// Return: true/false
//
func LetterNumberHyphen(ch rune) bool {
    if ch >= 'a' && ch <= 'z' { return true }
    if ch >= 'A' && ch <= 'Z' { return true }
    if ch >= '0' && ch <= '9' { return true }
    // hyphen or soft hyphen
    if ch == '-' || ch == 0xAD { return true }
    if ch == '_' { return true }
    if ch >= 0xC0  && ch < 0x2C0 { return true }
    if ch >= 0x380 && ch < 0x510 { return true }
    return false
}

// invalidBetweenTerms - Tests for an invalid character.
// Quote, hash and 'at' are invalid between terms.
func invalidBetweenTerms(ch rune) bool {
    if ch == '"' { return true }
    if ch == '#' { return true }
    if ch == '@' { return true }
    return false
}

// generateGoal - generates a goal (GoalStruct) from a string.
// For example, a string such as "can_swim($X), can_fly($X)"
// will become an And goal, with two complex statements,
// can_swim and can_fly, as subgoals.
//
// Param:  string of tokens
// Return: goal
//
func generateGoal(str string) (Goal, error) {
    tokens, err := Tokenize(str)
    if err != nil {
        panic(err.Error)
    }
    baseToken := groupTokens(tokens, 0)
    baseToken = groupAndTokens(baseToken)
    baseToken = groupOrTokens(baseToken)
    return tokenTreeToGoal(baseToken)
}


// Tokenize - Divides the given string into a series of tokens.
//
// Note: Parentheses can be part of a complex term: likes(Charles, Gina)
// or used to group terms: (father($_, $X); mother($_, $X))
//
// Params: string to parse
// Return: tokens
//         error
//
func Tokenize(str string) ([]TokenStruct, error) {

    tokens := []TokenStruct{}

    stkParenth := IntStack{} // Keeps track of parentheses.

    s := strings.TrimSpace(str)

    if len(s) == 0 {
        err := fmt.Errorf("Tokenize() - String is empty.")
        return tokens, err
    }

    runes := []rune(s)
    startIndex := 0
    length := len(runes)

    // Find a separator (comma, semicolon), if there is one.
    previous := '#' // random

    for i := startIndex; i < length; i++ {

        // Get top of stack.
        top, _ := stkParenth.Peek()

        ch := runes[i]
        if ch == '"' { // Ignore characters between quotes.
            for j := i + 1; j < length; j++ {
                ch = runes[j]
                if ch == '"' {
                    i = j
                    break
                }
            }
        } else if ch == '(' {
            // Is the previous character valid in a functor?
            if LetterNumberHyphen(previous) {
                stkParenth.Push(COMPLEX)
            } else {
                stkParenth.Push(GROUP)
                tokens = append(tokens, TokenLeaf("("))
                startIndex = i + 1
            }
        } else if ch == ')' {
            if top == NONE {
                err := fmt.Errorf("Tokenize() - Unmatched parenthesis: %v", s)
                return tokens, err
            }
            top, _ = stkParenth.Pop()
            if top == GROUP {
                subgoal := s[startIndex: i]
                tokens = append(tokens, TokenLeaf(subgoal))
                tokens = append(tokens, TokenLeaf(")"))
            } else if top != COMPLEX {
                err := fmt.Errorf("Tokenize() - Unmatched parenthesis: %v", s)
                return tokens, err
            }
        } else if ch == '[' {
            stkParenth.Push(LINKEDLIST)
        } else if ch == ']' {
            if top == NONE {
                err := fmt.Errorf("Tokenize() - Unmatched bracket: %v", s)
                return tokens, err
            }
            top, _ = stkParenth.Pop()
            if top != LINKEDLIST {
                err := fmt.Errorf("Tokenize() - Unmatched bracket: %v", s)
                return tokens, err
            }
        } else {
            // If not inside complex term or linked list...
            if top != COMPLEX && top != LINKEDLIST {
                if invalidBetweenTerms(ch) {
                    err := fmt.Errorf("Tokenize() - Invalid character: %v", s)
                    return tokens, err
                }
                if ch == ',' {   // AND
                    subgoal := s[startIndex: i]
                    tokens = append(tokens, TokenLeaf(subgoal))
                    tokens = append(tokens, TokenLeaf(","))
                    startIndex = i + 1
                } else if ch == ';' {   // OR
                    subgoal := s[startIndex: i]
                    tokens = append(tokens, TokenLeaf(subgoal))
                    tokens = append(tokens, TokenLeaf(";"))
                    startIndex = i + 1
                }
            }
        } // else
        previous = ch
    } // for

    if len(stkParenth) > 0 {
        err := fmt.Errorf("Tokenize() - Invalid character: %v", s)
        return tokens, err
    }

    if length - startIndex > 0 {
        subgoal := s[startIndex: length]
        tokens = append(tokens, TokenLeaf(subgoal))
    }

    return tokens, nil

} // Tokenize


// groupTokens - collects tokens within parentheses into groups.
// Converts a flat array of tokens into a tree of tokens.
//
// For example, this:   SUBGOAL SUBGOAL ( SUBGOAL  SUBGOAL )
// becomes:
//          GROUP
//            |
// SUBGOAL SUBGOAL GROUP
//                   |
//            SUBGOAL SUBGOAL
//
// There is a precedence order in subgoals. From highest to lowest.
//
//    groups (...)  -> GROUP
//    conjunction , -> AND
//    disjunction ; -> OR
//
// Params: flat array of tokens
//         starting index
// Return: base of token tree
//
func groupTokens(tokens []TokenStruct, index int) TokenStruct {

    newTokens := []TokenStruct{}
    size := len(tokens)

    for index < size {

        token := tokens[index]
        theType := token.theType

        if theType == LPAREN {
            index++
            // Make a GROUP token.
            t := groupTokens(tokens, index)
            newTokens = append(newTokens, t)
            // Skip past tokens already processed.
            index += t.numberOfChildren() + 1  // +1 for right parenthesis
        } else if theType == RPAREN {
            // Add all remaining tokens to the list.
            return TokenBranch(GROUP, newTokens)
        } else {
            newTokens = append(newTokens, token)
        }
        index++

    } // for

    return TokenBranch(GROUP, newTokens)

} // groupTokens


// groupAndTokens - groups tokens which are separated by commas.
//
// Param:  base of token tree
// Return: base of token tree
//
func groupAndTokens(token TokenStruct) TokenStruct {

    children    := token.children
    newChildren := []TokenStruct{}
    andList     := []TokenStruct{}

    for _, token := range children {

        theType := token.theType

        if theType == SUBGOAL {
            andList = append(andList, token)
        } else if theType == COMMA {
            // Nothing to do.
        } else if theType == SEMICOLON {
            // Must be end of comma separated list.
            size := len(andList)
            if size == 1 {
                newChildren = append(newChildren, andList[0])
            } else {
                newChildren = append(newChildren, TokenBranch(AND, andList))
            }
            newChildren = append(newChildren, token)
            andList = []TokenStruct{}
        } else if theType == GROUP {
            t := groupAndTokens(token)
            t  = groupOrTokens(token)
            andList = append(andList, t)
        }

    } // for

    size := len(andList)
    if size == 1 {
        newChildren = append(newChildren, andList[0])
    } else if size > 1 {
        newChildren = append(newChildren, TokenBranch(AND, andList))
    }

    token.children = newChildren
    return token

} // groupAndTokens


// groupOrTokens - groups tokens which are separated by semicolons.
//
// Param:  base of token tree
// Return: base of token tree
//
func groupOrTokens(token TokenStruct) TokenStruct {

    children    := token.children
    newChildren := []TokenStruct{}
    orList     := []TokenStruct{}

    for _, token := range children {

        theType := token.theType

        if theType == SUBGOAL || theType == AND || theType == GROUP {
            orList = append(orList, token)
        } else if theType == SEMICOLON {
            // Nothing to do.
        }

    } // for

    size := len(orList)
    if size == 1 {
        newChildren = append(newChildren, orList[0])
    } else if size > 1 {
        newChildren = append(newChildren, TokenBranch(OR, orList))
    }

    token.children = newChildren
    return token

} // groupOrTokens


// tokenTreeToGoal - produces a goal from the given token tree.
// Params: base of token tree
// Return: goal
//         error
func tokenTreeToGoal(token TokenStruct) (Goal, error) {

    var children []TokenStruct
    var operands []Goal

    if token.theType == SUBGOAL {
        g, err := ParseSubgoal(token.token)
        return g, err
    }

    if token.theType == AND {
        operands = []Goal{}
        children = token.children
        var err error
        var g Goal
        for _, child := range children {
            if child.theType == SUBGOAL {
               g, err = ParseSubgoal(child.token)
               operands = append(operands, g)
            } else if child.theType == GROUP {
               g, err = tokenTreeToGoal(child)
               operands = append(operands, g)
            }
        }
        return And(operands...), err
    }

    if token.theType == OR {
        operands = []Goal{}
        children = token.children
        var err error
        var g Goal
        for _, child := range children {
            if child.theType == SUBGOAL {
               g, err = ParseSubgoal(child.token)
               operands = append(operands, g)
            } else if child.theType == GROUP {
               g, err = tokenTreeToGoal(child)
               operands = append(operands, g)
            }
        }
        return Or(operands...), err
    }

    if token.theType == GROUP {
        if token.numberOfChildren() != 1 {
            panic("generateGoal - Group should have 1 child token.")
        }
        childToken := token.children[0]
        return tokenTreeToGoal(childToken)
    }

    return nil, fmt.Errorf("tokenTreeToGoal - Unknown token.")

}  // tokenTreeToGoal()


// showTokens - Displays a flat list of tokens for debugging purposes.
// Params: slice of tokens
func showTokens(tokens []TokenStruct) {
    first := true
    for _, token := range tokens {
        if !first { fmt.Print(" ") }
        first = false
        theType := token.theType
        if theType == SUBGOAL {
            fmt.Print(token.token)
        } else {
            fmt.Print(token.name)
        }
    } // for
    fmt.Print("\n")
} // showTokens
