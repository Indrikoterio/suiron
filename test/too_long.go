package main

// TooLong - this predicate sleeps, to simulate a calculation which
// takes too long to execute. It is used to test the "Time out" feature
// of SolveAll(), in solutions.go. The test is in solutions_test.go.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "time"
    //"fmt"
)

type TooLongStruct BuiltInPredicateStruct

// TooLong - creates a TooLong predicate.
// This function is variadic.
func TooLong(arguments ...Unifiable) TooLongStruct {
    ts := TooLongStruct {
        Name: "TooLong",
        Arguments: arguments,
    }
    return ts
}

// String - creates a string representation.
func (ts TooLongStruct) String() string {
    return "too_long"
}

// getSolver - gets solution node for this predicate.
func (ts TooLongStruct) GetSolver(kb KnowledgeBase,
                                  parentSolution SubstitutionSet,
                                  parentNode SolutionNode) SolutionNode {
    node := makeTooLongSolutionNode(ts, kb, parentSolution, parentNode)
    return node
}

// RecreateVariables - Refer to comments in expression.go.
func (ts TooLongStruct) RecreateVariables(
                        vars map[string]VariableStruct) Expression {
    return ts
}

// ReplaceVariables - Refer to comments in expression.go.
func (ts TooLongStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return ts
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeTooLongSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

type TooLongSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

func makeTooLongSolutionNode(goal Goal, kb KnowledgeBase,
                             parentSolution SubstitutionSet,
                             parentNode SolutionNode) SolutionNode {

    node := TooLongSolutionNodeStruct{
                SolutionNodeStruct: MakeSolutionNode(goal, kb,
                                        parentSolution, parentNode),
                moreSolutions: true,
            }

    return &node
}

// NextSolution - calls longCalculation, which takes too long.
// Returns:
//    updated substitution set
//    success/failure flag
// This function satisfies the SolutionNode interface.
func (sn *TooLongSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if sn.NoBackTracking || !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.
    // Arguments are not needed here for this predicate.
    return longCalculation(sn.ParentSolution)
}

// SetNoBackTracking - set the NoBackTracking flag.
// This flag is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *TooLongSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *TooLongSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}

// longCalculation - Sleep for 10 second.
// Returns unchanged substitution set and true for success.
func longCalculation(ss SubstitutionSet) (SubstitutionSet, bool) {
    time.Sleep(10 * time.Second)
    return ss, true
} // longCalculation()
