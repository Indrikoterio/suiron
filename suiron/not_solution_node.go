package suiron

// Solution node for the Not operator.
// Eg.:  not($X = noun)
//
// Cleve Lendon

import (
    //"fmt"
)

type NotSolutionNodeStruct struct {
    SolutionNodeStruct
    operandSolutionNode SolutionNode
}

func makeNotSolutionNode(n NotOp, kb KnowledgeBase,
                         parentSolution SubstitutionSet,
                         parentNode SolutionNode) SolutionNode {

    operand := n[0]  // There must be 1 operand.
    osn := operand.GetSolver(kb, parentSolution, parentNode)

    node := NotSolutionNodeStruct{
                SolutionNodeStruct: SolutionNodeStruct {
                                    Goal: n,
                                    KnowledgeBase: kb,
                                    ParentSolution: parentSolution,
                                    ParentNode: parentNode },

                // There is only one operand.
                operandSolutionNode: osn,
            }
    return &node
}

// NexSolution - calls NextSolution() on the operand (which is a Goal).
// If there is a solution, the function will set the success flag to false.
// If there is no solution, the function will set the success flag to true.
// ('Not' means 'not unifiable'.) 
// Returns:  substitution set
//           success/failure flag
func (n *NotSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {

    if n.NoBackTracking { return nil, false }
    if n.ParentSolution == nil { return nil, false }

    _, found := n.operandSolutionNode.NextSolution()
    if found {
        return nil, false
    } else {
        solution := n.ParentSolution
        n.ParentSolution = nil
        return solution, true
    }
}

// HasNextRule - returns true if the knowledge base contains untried
// rules for this node's goal. False otherwise.
func (n *NotSolutionNodeStruct) HasNextRule() bool {
    if n.NoBackTracking { return false }
    return n.ruleNumber < n.count
}

// NextRule - fetches the next rule from the database, according to ruleNumber.
// The method HasNextRule must be called to ensure that a rule can be fetched
// from the knowledge base. If GetRule is called with invalid parameters, the
// knowledge base will panic.
func (n *NotSolutionNodeStruct) NextRule() RuleStruct {
    rule := n.KnowledgeBase.GetRule(n.Goal, n.ruleNumber)
    n.ruleNumber++
    return rule
}

// SetNoBackTracking - rset the NoBackTracking flag.
// This flag is used to implement Cuts.
func (n *NotSolutionNodeStruct) SetNoBackTracking() {
    n.NoBackTracking = true
}

// GetParentNode
func (n *NotSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}
