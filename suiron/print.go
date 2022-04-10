package suiron

// Print - prints out a list of terms. These terms can be any unifiable,
// such as Atoms, Integers, Variables, etc. Variables are replaced with
// their ground terms. In a Suiron program, 'print' can be included in
// a rule as follows:
//
//   rule1 :- ..., print($X, b, c), ...
//
// Assuming that the variable $X is bound to 'a', the above would print out:
//
//   a, b, c
//
// In a Go program, the above is equivalent to:
//
//   And(..., Print(x, b, c), ...)
//
// (Assuming that x is a variable bound to 'a'.)
//
// If the first term is an Atom (string) which contains format specifiers
// (%s), it is treated as a format string. For example,
//
//   $Name = John, $Age = 23, print(%s is %s years old., $Name, $Age).
//
// will print out,
//
//   John is 23 years old.
//
// Commas which do not separate arguments, but are intended to be printed,
// can be escaped with a backslash, for example:
//
//   print(%s\, my friend\, is $s years old.\n, $Name, $Age)
//
// will print out,
//
//   John, my friend, is 23 years old.
//
// Alternatively, instead of using backslashes, a string can be enclosed
// within double quotes:
//
//   print("%s, my friend, is $s years old.\n", $Name, $Age)
//
// Cleve Lendon

import (
    "strings"
    "fmt"
)

type PrintStruct BuiltInPredicateStruct

const FORMAT_SPECIFIER = "%s"

// Print - creates a print predicate.
// This function is variadic.
func Print(arguments ...Unifiable) PrintStruct {
    return PrintStruct {
        Name: "print",
        Arguments: arguments,
    }
}

// GetSolver - gets solution node for Print predicate.
func (ps PrintStruct) GetSolver(kb KnowledgeBase,
                                parentSolution SubstitutionSet,
                                parentNode SolutionNode) SolutionNode {
    return makePrintSolutionNode(ps, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (ps PrintStruct) RecreateVariables(
                               vars map[VariableStruct]VariableStruct) Expression {
    bip := BuiltInPredicateStruct(ps).RecreateVariables(vars)
    return Expression(PrintStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (ps PrintStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(ps).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  predicate_name(arg1, arg2, arg3)
func (ps PrintStruct) String() string {
    return BuiltInPredicateStruct(ps).String()
}

//----------------------------------------------------------------
// Solution Node functions.
//    makePrintSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

type PrintSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

func makePrintSolutionNode(goal Goal, kb KnowledgeBase,
                           parentSolution SubstitutionSet,
                           parentNode SolutionNode) SolutionNode {

    node := PrintSolutionNodeStruct{
                SolutionNodeStruct: SolutionNodeStruct {
                                    Goal: goal,
                                    KnowledgeBase: kb,
                                    ParentSolution: parentSolution,
                                    ParentNode: parentNode },
                moreSolutions: true,
            }
    return &node
}

// NextSolution - calls a function to evaluate the current goal,
// based on its arguments and the substitution set.
// Returns:
//    updated substitution set
//    success/failure flag
// This function satisfies the SolutionNode interface.
func (sn *PrintSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if sn.NoBackTracking { return nil, false }
    if !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.
    goal := sn.Goal.(PrintStruct)
    return printTerms(goal.Arguments, sn.ParentSolution)
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *PrintSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *PrintSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}

//---------------------------------------------------------

// indexAt - gets the index of a substring within a string,
// starting from the given index.
// Params: string
//         substring to find
//         index of start of search
func indexAt(str string, substr string, start int) int {
    index := strings.Index(str[start:], substr)
    if index > -1 { index += start }
    return index
}

// isFormatString - returns true if the given string is a format string.
// (That is, the string contains a format specifier.)
func isFormatString(str string) bool {
    if strings.Index(str, FORMAT_SPECIFIER) < 0 { return false }
    return true
}

// splitFormatString
//
// A format string looks like this:
//   "Hello %s, my name is %s."
// This function will divide a string into substrings:
//   "Hello ", "%s", ", my name is ", "%s", "."
// Params: original string
// Return: array of substrings
func splitFormatString(str string) []string {
    sections := []string{}
    start := 0
    length := len(str)
    for start < length {
        index := indexAt(str, FORMAT_SPECIFIER, start)
        if index < 0 {
            sections = append(sections, str[start:])
            start = length
        } else {
            sections = append(sections, str[start: index])
            start = index
            index += 2
            sections = append(sections, str[start: index])
            start = index
        }
    }
    return sections
}  // splitFormatString

// getGround - if the term is a Variable, return its ground term.
// Else, return the term unchanged.
//
// Params:  term
//          substitution set
// Return:  ground term
func getGround(term Unifiable, ss SubstitutionSet) Unifiable {
    tt := term.TermType()
    if tt == VARIABLE {
        if ss.IsGroundVariable(term.(VariableStruct)) {
            gr, _ := ss.GetGroundTerm(term.(VariableStruct))
            return gr
        }
    }
    return term
} // getGround


// printTerms - prints out the ground terms of all arguments.
// If the first term has a format specifier, treat it as a format string.
func printTerms(arguments []Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

    if len(arguments) == 0 { return ss, false }

    // Get first argument.
    term0 := arguments[0]
    term0 = getGround(term0, ss)
    term0Str := term0.String()
    if isFormatString(term0Str) {
        formatSubstrings := splitFormatString(term0Str)
        count := 1
        for _, format := range formatSubstrings {
            if format == FORMAT_SPECIFIER {
                if count < len(arguments) {
                    t := arguments[count]
                    t = getGround(t, ss)
                    fmt.Print(t.String())
                    count++
                } else {
                    fmt.Print(format)
                }
            } else {
                fmt.Print(format)
            }
        } // for
    } else { // Not a format string.
        fmt.Print(term0Str)
        first := true
        for _, term := range arguments {
            if first {
                first = false
            } else {
                term = getGround(term, ss)
                fmt.Printf(", %v", term.String())
            }
        }
    }
    return ss, true  // Can't fail.

} // printTerms
