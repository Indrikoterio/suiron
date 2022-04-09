package suiron

// Solution node for the Fail operator.
//
// Cleve Lendon

import (
    //"fmt"
)

type FailSolutionNodeStruct SolutionNodeStruct

func MakeFailSolutionNode(f FailOp, kb KnowledgeBase,
                         parentSolution SubstitutionSet,
                         parentNode SolutionNode) SolutionNode {

    node := FailSolutionNodeStruct{
                Goal: f,
                KnowledgeBase: kb,
                ParentSolution: parentSolution,
                ParentNode: parentNode,
            }
    return &node
}

// NexSolution
// The Fail operator always fails. Return false.
func (n *FailSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    return n.ParentSolution, false
}

// SetNoBackTracking - set the NoBackTracking flag.
func (n *FailSolutionNodeStruct) SetNoBackTracking() {
    n.NoBackTracking = true
}

// GetParentNode
func (n *FailSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}
