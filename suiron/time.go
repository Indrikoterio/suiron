package suiron

// time - This predicate measures the execution time of a goal.
// The following example would print out the execution time of
// the qsort predicate in milliseconds.
//
//   ..., time(qsort()),...
//
// Time accepts only one goal, which must be a complex term.
//
// Cleve Lendon

import (
    //"fmt"
)

type TimeStruct BuiltInPredicateStruct

// Time - creates a TimeStruct, which holds the predicate's
// name and arguments. It accepts only one argument, which
// must be a complex term.
// Params: 1 complex term
// Return: TimeStruct
func Time(arguments ...Unifiable) TimeStruct {
    if len(arguments) != 1 {
        panic("Time - This predicate requires 1 argument.")
    }
    arg := arguments[0]
    tt := arg.TermType()
    if tt != COMPLEX {
        panic("Time - The argument must be a complex term.")
    }
    return TimeStruct {
        Name: "time",
        Arguments: arguments,
    }
}

// GetSolver - gets a solution node for the Time predicate.
func (ts TimeStruct) GetSolver(kb KnowledgeBase,
                               parentSolution SubstitutionSet,
                               parentNode SolutionNode) SolutionNode {

    return MakeTimeSolutionNode(ts, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (ts TimeStruct) RecreateVariables(
                               vars map[VariableStruct]VariableStruct) Expression {
    bip := BuiltInPredicateStruct(ts).RecreateVariables(vars)
    return Expression(TimeStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (ts TimeStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(ts).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  predicate_name(arg1)
func (ts TimeStruct) String() string {
    return BuiltInPredicateStruct(ts).String()
}
