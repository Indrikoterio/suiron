package suiron

// Solution node for the Cut operator. (!)
//
// Cleve Lendon

import (
    //"fmt"
)

type CutSolutionNodeStruct SolutionNodeStruct

func MakeCutSolutionNode(c CutOp, kb KnowledgeBase,
                         parentSolution SubstitutionSet,
                         parentNode SolutionNode) SolutionNode {

    node := CutSolutionNodeStruct{
                Goal: c,
                KnowledgeBase: kb,
                ParentSolution: parentSolution,
                ParentNode: parentNode,
            }
    return &node
}

// NexSolution
func (n *CutSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {

    if n.NoBackTracking { return n.ParentSolution, false }
    n.NoBackTracking = true

    // Set NoBackTracking on all ancestors.
    parent := n.ParentNode
    for parent != nil {
        parent.SetNoBackTracking()
        parent = parent.GetParentNode()
    }
    return n.ParentSolution, true
}

// SetNoBackTracking - set the NoBackTracking flag.
// This flag is used to implement Cuts.
func (n *CutSolutionNodeStruct) SetNoBackTracking() {
    n.NoBackTracking = true
}

// GetParentNode
func (n *CutSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}
