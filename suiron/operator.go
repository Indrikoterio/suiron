package suiron

// Operator - Base type for And, Or, Not logic operators, etc.
//
// Cleve Lendon

import (
    //"strings"
    //"fmt"
)

// An operator consists of a list of operands (= Goals).
type Operator []Goal

// copy - makes a copy of this operator.
func (op Operator) copy() Operator {
    operands := []Goal{}
    for _, g := range op { operands = append(operands, g) }
    return Operator(operands)
}

// getOperand - gets operand by index.
func (op Operator) getOperand(i int) Goal {
    return op[i]
}

// getHeadOperand - gets the first operand of the operand list.
func (op Operator) getHeadOperand() Goal {
    return op[0]
}

// getTailOperands - gets the tail of the operand list.
// (All operands except head.)
func (op Operator) getTailOperands() Operator {
    operands := []Goal{}
    for n, g := range op {
        if n == 0 { continue }
        operands = append(operands, g)
    }
    return Operator(operands)
}

// RecreateVariables - The scope of a logic variable is the rule in which
// it is defined. This method satisfies the Expression and Goal interfaces.
// Refer to comments in expression.go.
func RecreateVariablesForOperators(op []Goal, vars VarMap) []Goal {
    newGoals := []Goal{}
    for i := 0; i < len(op); i++ {
        goal := op[i]
        newGoals = append(newGoals, goal.RecreateVariables(vars).(Goal))
    }
    return newGoals
}

// ReplaceVariablesForOperators()
// Operators must implement ReplaceVariables(), in order to satisfy the
// Expression and Goal interfaces. Refer to comments in expression.go.
func ReplaceVariablesForOperators(op []Goal, ss SubstitutionSet) []Goal {
    newGoals := []Goal{}
    for i := 0; i < len(op); i++ {
        goal := op[i]
        newGoals = append(newGoals, goal.ReplaceVariables(ss).(Goal))
    }
    return newGoals
}
