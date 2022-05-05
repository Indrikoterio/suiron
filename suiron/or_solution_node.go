package suiron

// Solution node for the Or operator.
//
// Cleve Lendon

import (
//    "fmt"
)

type OrSolutionNodeStruct struct {
    SolutionNodeStruct
    headSolutionNode SolutionNode
    tailSolutionNode SolutionNode
    operatorTail OrOp
    parentSolution SubstitutionSet
}

func makeOrSolutionNode(o OrOp, kb KnowledgeBase,
                        parentSolution SubstitutionSet,
                        parentNode SolutionNode) SolutionNode {

    op := Operator(o)
    headOp  := op.getHeadOperand()
    tailOps := OrOp(op.getTailOperands())

    node := OrSolutionNodeStruct{
                SolutionNodeStruct: SolutionNodeStruct{
                                    Goal: o,
                                    KnowledgeBase: kb,
                                    ParentSolution: parentSolution,
                                    ParentNode: parentNode },

                // The first operand could be anything,
                // perhaps a complex term.
                headSolutionNode: headOp.GetSolver(kb, parentSolution, parentNode),

                // The operator tail is the same as the original
                // 'Or' minus the first goal.
                operatorTail: tailOps,

                parentSolution: parentSolution,
            }
    
    return &node
}

// NexSolution - recursively calls NextSolution on all subgoals.
// If the search succeeds, the boolean return value is true, and
// the substitution set is updated.
// If the search fails, the boolean value is false.
func (o *OrSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {

    var solution SubstitutionSet
    var found bool

    if o.NoBackTracking { return nil, false }

    if o.tailSolutionNode != nil {
        return o.tailSolutionNode.NextSolution()
    }

    solution, found = o.headSolutionNode.NextSolution()

    if found || len(o.operatorTail) == 0 {
        return solution, found
    } else {
        // tailSolutionNode has to be a new OrSolutionNode.
        o.tailSolutionNode = o.operatorTail.
                                 GetSolver(o.KnowledgeBase,
                                           o.ParentSolution, o)
        return o.tailSolutionNode.NextSolution()
    }
}

// HasNextRule - returns true if the knowledge base contains untried
// rules for this node's goal. False otherwise.
func (o *OrSolutionNodeStruct) HasNextRule() bool {
    if o.NoBackTracking { return false }
    return o.ruleNumber < o.count
}

// NextRule - fetches the next rule from the database, according to ruleNumber.
// The method HasNextRule must be called to ensure that a rule can be fetched
// from the knowledge base. If GetRule is called with invalid parameters, the
// knowledge base will panic.
func (o *OrSolutionNodeStruct) NextRule() RuleStruct {
    rule := o.KnowledgeBase.GetRule(o.Goal, o.ruleNumber)
    o.ruleNumber++
    return rule
}

// SetNoBackTracking - set the NoBackTracking flag.
// This flag is used to implement Cuts.
func (o *OrSolutionNodeStruct) SetNoBackTracking() {
    o.NoBackTracking = true
}

// GetParentNode
func (n *OrSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}
