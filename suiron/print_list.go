package suiron

// PrintList - This built-in predicate prints out a linked list in
// a readable form. It's mainly for debugging purposes. For example:
//
// a_rule :- $X = [a, b, c], $List = [1, 2, 3 | $X], print_list($List).
//
// The rule above should print out: 1, 2, 3, a, b, c
//
// print_list() skips terms which are not ground.
//
// Cleve Lendon

import (
    "fmt"
)

type PrintListStruct BuiltInPredicateStruct

// PrintList - creates a print_list predicate.
func PrintList(arguments ...Unifiable) PrintListStruct {
    return PrintListStruct {
        Name: "print_list",
        Arguments: arguments,
    }
}

// GetSolver - gets solution node for PrintList predicate.
func (pls PrintListStruct) GetSolver(kb KnowledgeBase,
                                     parentSolution SubstitutionSet,
                                     parentNode SolutionNode) SolutionNode {
    return makePrintListSolutionNode(pls, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (pls PrintListStruct) RecreateVariables(vars VarMap) Expression {
    bip := BuiltInPredicateStruct(pls).RecreateVariables(vars)
    return Expression(PrintListStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (pls PrintListStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(pls).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  predicate_name(arg1, arg2, arg3)
func (pls PrintListStruct) String() string {
    return BuiltInPredicateStruct(pls).String()
}

//----------------------------------------------------------------
// Solution Node functions.
//    makePrintListSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

type PrintListSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

func makePrintListSolutionNode(goal Goal, kb KnowledgeBase,
                               parentSolution SubstitutionSet,
                               parentNode SolutionNode) SolutionNode {

    node := PrintListSolutionNodeStruct{
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
func (sn *PrintListSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {

    ss := sn.ParentSolution

    if sn.NoBackTracking { return ss, false }
    if !sn.moreSolutions { return ss, false }
    sn.moreSolutions = false  // Only one solution.

    goal := sn.Goal.(PrintListStruct)

    if len(goal.Arguments) == 0 { return ss, false }
    for _, term := range goal.Arguments {
        term, ok := ss.GetGroundTerm(term)
        if !ok { continue }
        tt := term.TermType()
        if tt == LINKEDLIST {
            showLinkedList(term.(LinkedListStruct), ss)
        }
    }
    return ss, true  // Can't fail.
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *PrintListSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *PrintListSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}

// showLinkedList - Prints out all terms in a list.
//
// Params: linked list
//         substitution set
//
func showLinkedList(theList LinkedListStruct, ss SubstitutionSet) {
    first := true
    ptr := &theList
    for ptr != nil {
        term := ptr.term
        if term == nil { break }
        gt, ok := ss.GetGroundTerm(term)
        if ok {
            if ptr.tailVar {
                tt := gt.TermType()
                if tt == LINKEDLIST {
                    ll := gt.(LinkedListStruct)
                    ptr = &ll
                    gt, ok = ss.GetGroundTerm(ptr.term)
                }
            }
            if ok {
                if !first { fmt.Print(", ") }
                first = false
                fmt.Print(gt)
            }
        }
        ptr = ptr.next
    }
    fmt.Print("\n")
} // showLinkedList
