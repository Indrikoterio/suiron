package suiron

// SolutionNode - represents a node in a 'proof tree'.
//
// Complex terms and operators (And, Or, Unify, etc.) implement a
// method called GetSolver(), which returns a SolutionNode specific
// to that complex term or operator.
//
// The method NextSolution() starts the search for a solution.
// When a solution is found, the search stops. Each node preserves
// its state (goal, ruleNumber, etc.). Calling NextSolution() again
// will continue the search from where it left off.
//
// Cleve Lendon

type SolutionNodeStruct struct {

    KnowledgeBase KnowledgeBase
    ParentSolution SubstitutionSet
    ParentNode SolutionNode
    NoBackTracking bool

    Goal Goal  // goal being solved
    ruleNumber int
    count int  // counts number of rules/facts
}

// MakeSolutionNode - makes a solution node with the given arguments:
// Params:
//     goal
//     knowledgebase
//     parent solution (substitution set)
//     solution node of parent
// Return:
//     a solution node struct
func MakeSolutionNode(goal Goal, kb KnowledgeBase,
                      parentSolution SubstitutionSet,
                      parentNode SolutionNode) SolutionNodeStruct {

    node := SolutionNodeStruct {
                Goal: goal,
                KnowledgeBase: kb,
                ParentSolution: parentSolution,
                ParentNode: parentNode }
    return node
}


type SolutionNode interface {
    NextSolution() (SubstitutionSet, bool)
    SetNoBackTracking()
    GetParentNode() SolutionNode
}
