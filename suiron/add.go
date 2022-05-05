package suiron

// Add
//
// A built-in function to add two or more numbers together. Eg.:
//
//   $X = add(7, 3, 2),...
//
// Cleve Lendon

import (
    "fmt"
)

type AddStruct BuiltInPredicateStruct

// Add - creates an AddStruct, which holds the function's name and arguments.
// Add requires at least 2 arguments.
// Params: arguments (Unifiable)
// Return: AddStruct
func Add(arguments ...Unifiable) AddStruct {
    if len(arguments) < 2 {
        panic("Add - requires at least 2 arguments.")
    }
    return AddStruct {
        Name: "add",
        Arguments: arguments,
    }
}

//----------------------------------------------------------------
// bifAdd - Adds all arguments together. All arguments must be bound.
//
// Params:
//     list of arguments
//     substitution set
// Returns:
//    new unifiable
//    success/failure flag
//
func bifAdd(arguments []Unifiable, ss SubstitutionSet) (Unifiable, bool) {

    ground := []Unifiable{}  // Array of ground terms.
    hasFloat := false

    for _, arg := range arguments {
        c, ok := ss.GetGroundTerm(arg)
        if !ok {
            s := fmt.Sprintf("Add - Unbound argument: %v", arg)
            panic(s)
        }
        ground = append(ground, c)
        tt := c.TermType()
        if tt == FLOAT { hasFloat = true }
    }

    if hasFloat {
        sum := Float(0.0)
        for _, arg := range ground {
            tt := arg.TermType()
            if tt == INTEGER {
                i := arg.(Integer)
                sum += Float(i)
            } else {
                sum += arg.(Float)
            }
        }
        return Unifiable(sum), true
    } else {
        sum := Integer(0)
        for _, arg := range ground {
            sum += arg.(Integer)
        }
        return Unifiable(sum), true
    }
} // bifAdd


//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (as AddStruct) RecreateVariables(vars VarMap) Expression {
    bif := BuiltInPredicateStruct(as).RecreateVariables(vars)
    return Expression(AddStruct(*bif))
}

// ReplaceVariables - Refer to comments in expression.go.
func (as AddStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(as).ReplaceVariables(ss)
}  // ReplaceVariables

// String - creates a string representation.
// Returns:  function_name(arg1, arg2, arg3)
func (as AddStruct) String() string {
    return BuiltInPredicateStruct(as).String()
}

//----------------------------------------------------------------
// Unify() and TermType() satisfy the Unifiable interface.
//----------------------------------------------------------------

// Unify - unifies the result of a function with another term
// (usually a variable).
//
// Params:
//    other unifiable term
//    substitution set
// Returns:
//    updated substitution set
//    success/failure flag
func (as AddStruct) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {
    result, ok := bifAdd(as.Arguments, ss)
    if !ok { return ss, false }
    return result.Unify(other, ss)
}

// TermType - returns a constant which identifies this type.
func (as AddStruct) TermType() int { return FUNCTION }
