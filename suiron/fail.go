package suiron

// Defines the fail operator. This logic operator always fails.
//
// Cleve Lendon

import (
    //"fmt"
)

type FailOp []Goal

// Fail - creates a logical Fail operator.
func Fail(operands ...Goal) FailOp {
    return FailOp(operands)
}

// GetSolver - gets solution node for Fail operator.
func (f FailOp) GetSolver(kb KnowledgeBase,
                         parentSolution SubstitutionSet,
                         parentNode SolutionNode) SolutionNode {
    node := MakeFailSolutionNode(f, kb, parentSolution, parentNode)
    return node
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables
// Refer to comments in expression.go.
func (c FailOp) RecreateVariables(vars map[Variable]Variable) Expression {
    return c
}

// ReplaceVariables
// Refer to comments in expression.go.
func (c FailOp) ReplaceVariables(ss SubstitutionSet) Expression {
    return c
}

// String - Creates a string representation of this operator.
func (c FailOp) String() string { return "fail" }
