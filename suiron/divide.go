package suiron

// Divide
//
// A built-in function to divide numbers. Eg.:
//
//   $X = divide(7, 3, 2),...   // (7 / 3 / 2)
//
// Divide always produces a floating point number.
//
// Cleve Lendon

import (
    "fmt"
)

type DivideStruct BuiltInPredicateStruct

// Divide - creates a DivideStruct, which holds the function's
// name and arguments. Divide requires at least 2 arguments.
// Params: arguments (Unifiable)
// Return: DivideStruct
func Divide(arguments ...Unifiable) DivideStruct {
    if len(arguments) < 2 {
        panic("Divide - requires at least 2 arguments.")
    }
    return DivideStruct {
        Name: "divide",
        Arguments: arguments,
    }
}

//----------------------------------------------------------------
// bifDivide - Divides all arguments together.
// All arguments must be bound.
//
// Params:
//     list of arguments
//     substitution set
// Returns:
//     new unifiable
//     success/failure flag
//
func bifDivide(arguments []Unifiable, ss SubstitutionSet) (Unifiable, bool) {

    ground := []Unifiable{}  // Array of ground terms.

    // Get ground terms.
    for _, arg := range arguments {
        c, ok := ss.GetGroundTerm(arg)
        if !ok {
            s := fmt.Sprintf("Divide - Argument is not ground: %v", arg)
            panic(s)
        }
        ground = append(ground, c)
    }

    var result Float
    arg := ground[0]
    tt  := arg.TermType()
    if tt == INTEGER {
        result = Float(arg.(Integer))
    } else {
        result = arg.(Float)
    }
    for n, arg := range ground {
        if n == 0 { continue }
        tt = arg.TermType()
        if tt == INTEGER {
            result /= Float(arg.(Integer))
        } else {
            result /= arg.(Float)
        }
    }
    return Unifiable(result), true

} // bifDivide


//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (ms DivideStruct) RecreateVariables(vars VarMap) Expression {
    bif := BuiltInPredicateStruct(ms).RecreateVariables(vars)
    return Expression(DivideStruct(*bif))
}

// ReplaceVariables - Refer to comments in expression.go.
func (ms DivideStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(ms).ReplaceVariables(ss)
}  // ReplaceVariables

// String - creates a string representation.
// Returns:  function_name(arg1, arg2, arg3)
func (ms DivideStruct) String() string {
    return BuiltInPredicateStruct(ms).String()
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
func (as DivideStruct) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {
    result, ok := bifDivide(as.Arguments, ss)
    if !ok { return ss, false }
    return result.Unify(other, ss)
}

// TermType - returns a constant which identifies this type.
func (as DivideStruct) TermType() int { return FUNCTION }
