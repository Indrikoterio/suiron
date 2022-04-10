package suiron

// And - Defines a logical And operator (Conjunction).
//
// Cleve Lendon

import (
    "strings"
)

type AndOp Operator

// And - creates a logical And operator (Conjunction).
// Params:  operands
func And(operands ...Goal) AndOp {
    return AndOp(operands)
}


// GetSolver - gets solution node for And operator.
func (a AndOp) GetSolver(kb KnowledgeBase,
                       parentSolution SubstitutionSet,
                       parentNode SolutionNode) SolutionNode {
    node := MakeAndSolutionNode(a, kb, parentSolution, parentNode)
    return node
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Operators must implement RecreateVariables(),
// in order to satisfy the Expression and Goal interfaces.
// Refer to comments in expression.go.
func (a AndOp) RecreateVariables(vars map[string]VariableStruct) Expression {
    return AndOp(RecreateVariablesForOperators(a, vars))
}

// ReplaceVariables - Operators must implement ReplaceVariables(),
// in order to satisfy the Expression and Goal interfaces.
// Refer to comments in expression.go.
func (a AndOp) ReplaceVariables(ss SubstitutionSet) Expression {
    return AndOp(ReplaceVariablesForOperators(a, ss))
}

// String - Creates a string for debugging purposes.
func (a AndOp) String() string {
    var sb strings.Builder
    for n, k := range a {
        if n != 0 { sb.WriteString(", ") }
        sb.WriteString(k.String())
    }
    return sb.String()
}
