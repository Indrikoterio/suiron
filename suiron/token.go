package suiron

// Tokens are used to parse Suiron's goals. Each token represents
// a node in a token tree.
//
// A token leaf can be: SUBGOAL, COMMA, SEMICOLON, LPAREN, RPAREN.
// If a token is a branch node, its type will be: GROUP, AND, OR.
//
// For example, for this goal: (mother($X, $Y); father($X, $Y))
// the token types would be: LPAREN SUBGOAL SEMICOLON SUBGOAL RPAREN.
//
// There is a precedence to subgoals. From highest to lowest.
//
//    groups (...)  -> GROUP
//    conjunction , -> AND
//    disjunction ; -> OR
//
// Cleve Lendon

import (
    "strings"
)


// TokenStruct is a node in a token tree.
type TokenStruct struct {
    theType  int
    token    string
    children []TokenStruct
    name     string
}

// TokenLeaf - produces a leaf node (TokenStruct) for the given string.
// Valid leaf node types are: COMMA, SEMICOLON, LPAREN, RPAREN, SUBGOAL.
//
// Param:  symbol or subgoal (string)
// Return: token (TokenStruct)
func TokenLeaf(str string) TokenStruct {
    s := strings.TrimSpace(str)
    if s == "," {
        return TokenStruct{ theType: COMMA, token: s, name: "COMMA" }
    }
    if s == ";" {
        return TokenStruct{ theType: SEMICOLON, token: s, name: "SEMICOLON" }
    }
    if s == "(" {
        return TokenStruct{ theType: LPAREN, token: s, name: "LPAREN" }
    }
    if s == ")" {
        return TokenStruct{ theType: RPAREN, token: s, name: "RPAREN" }
    }
    return TokenStruct{ theType: SUBGOAL, token: s, name: "SUBGOAL" }
} // TokenLeaf


// TokenBranch - produces a branch node (TokenStruct) with the given
// child nodes. Valid branch node types are: GROUP, AND, OR.
// Param:  type of node
//         child nodes
// Return: a branch (parent) node
func TokenBranch(theType int, children []TokenStruct) TokenStruct {
    name := "??"
    if theType == GROUP {
        name = "GROUP"
    } else if theType == AND {
        name = "AND"
    } else if theType == OR {
        name = "OR"
    }
    return TokenStruct{ theType: theType, children: children, name: name }
} // TokenBranch


// numberOfChildren - returns the number of children, if any.
// Otherwise returns 1.
func (ts TokenStruct) numberOfChildren() int {
    return len(ts.children)
}

// String - produces a string representation of the node
// for debugging. Eg.
//
//    SUBGOAL > sister(Janelle, Amanda)
//
// Return: printable string
//
func (ts TokenStruct) String() string {
    var sb strings.Builder
    sb.WriteString(ts.name)
    if ts.theType == AND || ts.theType == OR {
        sb.WriteString(" > ")
        for _, child := range ts.children {
            sb.WriteString(child.name)
            sb.WriteString(" ")
        }
    } else if ts.theType == SUBGOAL {
        sb.WriteString(" > ")
        sb.WriteString(ts.token)
    }
    return sb.String()
} // String
