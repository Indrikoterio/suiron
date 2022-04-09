package suiron

// NewLine - Prints out a new line character. The parser recognizes nl.
// For example:
//
//   greeting :- print("Are you OK, world?"), nl.
//
// In a Go program file, call the function NL() to output a new line.
// The above rule could be written as:
//
//     msg  := Atom("Are you OK, world?")
//     pr   := Print(msg)
//     body := And(pr, NL())
//     head := Complex{Atom("greeting")}
//     rule := Rule(head, body)
//
// Cleve Lendon

import ( "fmt" )

type NewLineStruct BuiltInPredicateStruct

// NL - creates a predicate which outputs a new-line code.
// No arguments are needed.
func NL() NewLineStruct {
    return NewLineStruct { Name: "nl" }
}

// GetSolver - gets solution node for new line predicates.
func (nls NewLineStruct) GetSolver(kb KnowledgeBase,
                                   parentSolution SubstitutionSet,
                                   parentNode SolutionNode) SolutionNode {
    return makeNewLineSolutionNode(nls, kb, parentSolution, parentNode)
}

// RecreateVariables - Refer to comments in expression.go.
func (nls NewLineStruct) RecreateVariables(
                            vars map[Variable]Variable) Expression {
    return nls
}

// ReplaceVariables - Refer to comments in expression.go.
func (nls NewLineStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return nls
}

// String - creates a string representation.
func (nls NewLineStruct) String() string { return "nl" }

//----------------------------------------------------------------
// Solution Node functions.
//    makeNewLineSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

type NewLineSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

// makeNewLineSolutionNode - creates a solution node for this predicate.
func makeNewLineSolutionNode(goal Goal, kb KnowledgeBase,
                             parentSolution SubstitutionSet,
                             parentNode SolutionNode) SolutionNode {

    node := NewLineSolutionNodeStruct{
                SolutionNodeStruct: SolutionNodeStruct {
                                    Goal: goal,
                                    KnowledgeBase: kb,
                                    ParentSolution: parentSolution,
                                    ParentNode: parentNode },
                moreSolutions: true,
            }
    return &node
}

// NextSolution - simply prints out a new line character.
// This function satisfies the SolutionNode interface.
func (nlsn *NewLineSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if nlsn.NoBackTracking { return nil, false }
    if !nlsn.moreSolutions { return nil, false }
    nlsn.moreSolutions = false  // Only one solution.
    fmt.Print("\n")
    // No changes to substitution set. Never fails.
    return nlsn.ParentSolution, true
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (nlsn *NewLineSolutionNodeStruct) SetNoBackTracking() {
    nlsn.NoBackTracking = true
}

// GetParentNode
func (n *NewLineSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}
