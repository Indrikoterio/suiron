package suiron

// Or - Defines a logical Or operator.
//
// Cleve Lendon

import (
    "strings"
)

type OrOp Operator

// Or - creates a logical Or operator.
// Params:  operands
func Or(operands ...Goal) OrOp {
      return OrOp(operands)
}

// String - Creates a string for debugging purposes.
func (o OrOp) String() string {
    var sb strings.Builder
    for n, k := range o {
        if n != 0 { sb.WriteString("; ") }
        sb.WriteString(k.String())
    }
    return sb.String()
}

// GetSolver - gets solution node for Or operator.
func (o OrOp) GetSolver(kb KnowledgeBase,
                      parentSolution SubstitutionSet,
                      parentNode SolutionNode) SolutionNode {
    return MakeOrSolutionNode(o, kb, parentSolution, parentNode)
}


// RecreateVariables - Operators must implement RecreateVariables(),
// in order to satisfy the Expression and Goal interfaces.
// Refer to comments in expression.go.
func (o OrOp) RecreateVariables(vars map[string]VariableStruct) Expression {
    return OrOp(RecreateVariablesForOperators(o, vars))
}

// ReplaceVariables - Operators must implement ReplaceVariables(),
// in order to satisfy the Expression and Goal interfaces.
// Refer to comments in expression.go.
func (o OrOp) ReplaceVariables(ss SubstitutionSet) Expression {
    return OrOp(ReplaceVariablesForOperators(o, ss))
}
