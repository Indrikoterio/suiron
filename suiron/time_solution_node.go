package suiron

// Solution node for the Time operator.
//
// Cleve Lendon

import (
    "time"
    "fmt"
)

type TimeSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
    goalToTime Complex
    gttSolutionNode SolutionNode
}

func MakeTimeSolutionNode(goal Goal, kb KnowledgeBase,
                          parentSolution SubstitutionSet,
                          parentNode SolutionNode) SolutionNode {

    goalToTime := (goal.(TimeStruct)).Arguments[0].(Complex)

    node := TimeSolutionNodeStruct{
                SolutionNodeStruct: SolutionNodeStruct {
                                    Goal: goal,
                                    KnowledgeBase: kb,
                                    ParentSolution: parentSolution,
                                    ParentNode: parentNode },
                goalToTime: goalToTime,
                gttSolutionNode: goalToTime.GetSolver(
                                 kb, parentSolution, parentNode),
                moreSolutions: true,
            }
    return &node
}

// NextSolution
// If the search fails, the boolean value is false.
func (n *TimeSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {

    if n.NoBackTracking { return nil, false }
    if !n.moreSolutions { return nil, false }
    n.moreSolutions = false  // Only one solution.

    for _, term := range n.goalToTime {
        if term.TermType() == VARIABLE {
            v := term.(VariableStruct)
            if !n.ParentSolution.IsGroundVariable(v) {
                fmt.Printf("Time: Variable %v is not grounded.\n", term.String())
            }
        }
    } // for

    theStartTime := time.Now()
    solution, found  := n.gttSolutionNode.NextSolution()

    // time in milliseconds since the start of the query.
    elapsed := time.Since(theStartTime).Nanoseconds() / 1_000_000
    fmt.Printf("Elapsed time for %v: %d milliseconds\n", n.goalToTime, elapsed)

    return solution, found
}

// HasNextRule - returns true if the knowledge base contains untried
// rules for this node's goal. False otherwise.
func (n *TimeSolutionNodeStruct) HasNextRule() bool {
    if n.NoBackTracking { return false }
    return n.ruleNumber < n.count
}

// NextRule - fetches the next rule from the database, according to ruleNumber.
// The method HasNextRule must be called to ensure that a rule can be fetched
// from the knowledge base. If GetRule is called with invalid parameters, the
// knowledge base will panic.
func (n *TimeSolutionNodeStruct) NextRule() RuleStruct {
    rule := n.KnowledgeBase.GetRule(n.Goal, n.ruleNumber)
    n.ruleNumber++
    return rule
}

// SetNoBackTracking - rset the NoBackTracking flag.
// This flag is used to implement Cuts.
func (n *TimeSolutionNodeStruct) SetNoBackTracking() {
    n.NoBackTracking = true
}

// GetParentNode
func (n *TimeSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}
