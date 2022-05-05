package suiron

// Not
//
// In Prolog, the Not operator means 'currently not provable'.
// Many Prolog implementations use \+ for this operator, because
// the word 'not' can be misleading.
//
// 'Not' takes only one operand. Eg. not(father(Jonathan, $_))
//
// Cleve Lendon

import(
    "fmt"
)

type NotOp Operator

// Not - creates an NotOp type, which holds the operator's operand.
// Params: operands (Goal)
// Return: NotOp
func Not(operands ...Goal) NotOp {
    if len(operands) != 1 { panic("Not - Accepts 1 operand.") }
    return NotOp(operands)
}

// GetSolver - gets solution node for the Not operator.
func (n NotOp) GetSolver(kb KnowledgeBase,
                         parentSolution SubstitutionSet,
                         parentNode SolutionNode) SolutionNode {

    return makeNotSolutionNode(n, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (n NotOp) RecreateVariables(vars VarMap) Expression {
    return NotOp(RecreateVariablesForOperators(n, vars))
}

// ReplaceVariables - Refer to comments in expression.go.
func (n NotOp) ReplaceVariables(ss SubstitutionSet) Expression {
    return NotOp(ReplaceVariablesForOperators(n, ss))
}

// String - Creates a string for debugging purposes.
func (n NotOp) String() string {
    return fmt.Sprintf("not(%v)", n[0])
}
