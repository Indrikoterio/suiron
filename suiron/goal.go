package suiron

// Goal - is the base interface of all goal objects. Complex terms and
// operators such as And, Or, Unify etc. implement this interface.
//
// The method GetSolver() provides a solution node.
//
// Cleve Lendon

type Goal interface {

    Expression

    // GetSolver - gets a solution node for the current goal.
    GetSolver(kb KnowledgeBase, parentSolution SubstitutionSet,
              parentNode SolutionNode) SolutionNode
}
