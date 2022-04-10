package suiron

// Defines the logical 'cut' operator (!).
//
// The cut operator stops backtracking. If the goals after a cut
// fails, the inference engine will not backtrack past the cut.
//
// Cleve Lendon

import (
    //"fmt"
)

type CutOp Operator

// Cut - creates a logical Cut operator.
func Cut(operands ...Goal) CutOp {
    return CutOp(operands)
}

// GetSolver - gets solution node for Cut operator.
func (c CutOp) GetSolver(kb KnowledgeBase,
                         parentSolution SubstitutionSet,
                         parentNode SolutionNode) SolutionNode {
    node := MakeCutSolutionNode(c, kb, parentSolution, parentNode)
    return node
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables
// Refer to comments in expression.go.
func (c CutOp) RecreateVariables(vars map[string]VariableStruct) Expression {
    return c
}

// ReplaceVariables
// Refer to comments in expression.go.
func (c CutOp) ReplaceVariables(ss SubstitutionSet) Expression {
    return c
}

// String - Creates a string representation of this operator.
func (c CutOp) String() string { return "!" }
