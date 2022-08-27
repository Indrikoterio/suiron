package suiron

// Solution node for the And operator.
//
// Cleve Lendon

import (
    //"fmt"
)

type AndSolutionNodeStruct struct {
    SolutionNodeStruct
    headSolutionNode SolutionNode
    tailSolutionNode SolutionNode
    operatorTail AndOp
}

func makeAndSolutionNode(a AndOp, kb KnowledgeBase,
                         parentSolution SubstitutionSet,
                         parentNode SolutionNode) SolutionNode {

    op := Operator(a)
    headOp  := op.getHeadOperand()
    tailOps := AndOp(op.getTailOperands())

    hsn := headOp.GetSolver(kb, parentSolution, parentNode)

    node := AndSolutionNodeStruct{
                SolutionNodeStruct: SolutionNodeStruct {
                                    Goal: a,
                                    KnowledgeBase: kb,
                                    ParentSolution: parentSolution,
                                    ParentNode: parentNode },

                // The first operand could be anything,
                // perhaps a complex term.
                headSolutionNode: hsn,

                // The operator tail is the same as the original
                // 'And' minus the first goal.
                operatorTail: tailOps,
            }
    return &node
}

// NextSolution - recursively calls NextSolution on all subgoals.
// If the search succeeds, the boolean return value is true, and
// the substitution set is updated.
// If the search fails, the boolean value is false.
func (n *AndSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {

    var solution SubstitutionSet
    var found bool

    if n.NoBackTracking { return nil, false }

    if n.tailSolutionNode != nil {
        solution, found = n.tailSolutionNode.NextSolution()
        if found { return solution, true }
    }

    solution, found = n.headSolutionNode.NextSolution()
    for found {
        if len(n.operatorTail) == 0 {
            return solution, true
        } else {
            // tailSolutionNode has to be a new AndSolutionNode.
            n.tailSolutionNode = n.operatorTail.
                                   GetSolver(n.KnowledgeBase, solution, n)
            tailSolution, found := n.tailSolutionNode.NextSolution()
            if found { return tailSolution, true }
        }
        solution, found = n.headSolutionNode.NextSolution()
    }
    return nil, false
}

// HasNextRule - returns true if the knowledge base contains untried
// rules for this node's goal. False otherwise.
func (n *AndSolutionNodeStruct) HasNextRule() bool {
    if n.NoBackTracking { return false }
    return n.ruleNumber < n.count
}

// NextRule - fetches the next rule from the database, according to ruleNumber.
// The method HasNextRule must be called to ensure that a rule can be fetched
// from the knowledge base. If GetRule is called with invalid parameters, the
// knowledge base will panic.
func (n *AndSolutionNodeStruct) NextRule() RuleStruct {
    rule := n.KnowledgeBase.GetRule(n.Goal, n.ruleNumber)
    n.ruleNumber++
    return rule
}

// SetNoBackTracking - set the NoBackTracking flag.
// This flag is used to implement Cuts.
func (n *AndSolutionNodeStruct) SetNoBackTracking() {
    n.NoBackTracking = true
}

// GetParentNode
func (n *AndSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}
