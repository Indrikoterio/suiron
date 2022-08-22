package suiron

// Multiply
//
// A built-in function to multiply numbers. Eg.:
//
//   $X = multiply(7, 3, 2),...   // (7 * 3 * 2)
//
// Cleve Lendon

import (
    "fmt"
)

type MultiplyStruct BuiltInPredicateStruct

// Multiply - creates a MultiplyStruct, which holds the function's
// name and arguments. Multiply requires at least 2 arguments.
// Params: arguments (Unifiable)
// Return: MultiplyStruct
func Multiply(arguments ...Unifiable) MultiplyStruct {
    if len(arguments) < 2 {
        panic("Multiply - requires at least 2 arguments.")
    }
    return MultiplyStruct {
        Name: "multiply",
        Arguments: arguments,
    }
}

//----------------------------------------------------------------
// bifMultiply - Multiplies all arguments together.
// All arguments must be bound.
//
// Params:
//     list of arguments
//     substitution set
// Returns:
//     new unifiable
//     success/failure flag
//
func bifMultiply(arguments []Unifiable, ss SubstitutionSet) (Unifiable, bool) {

    ground := []Unifiable{}  // Array of ground terms.
    hasFloat := false

    // Get ground terms.
    for _, arg := range arguments {
        c, ok := ss.GetGroundTerm(arg)
        if !ok {
            s := fmt.Sprintf("Multiply - Argument is not ground: %v", arg)
            panic(s)
        }
        ground = append(ground, c)
        tt := c.TermType()
        if tt == FLOAT { hasFloat = true }
    }

    if hasFloat {
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
                result *= Float(arg.(Integer))
            } else {
                result *= arg.(Float)
            }
        }
        return Unifiable(result), true
    } else {
        result := ground[0].(Integer)
        for n, arg := range ground {
            if n == 0 { continue }
            result *= arg.(Integer)
        }
        return Unifiable(result), true
    }

} // bifMultiply


//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (ms MultiplyStruct) RecreateVariables(vars VarMap) Expression {
    bif := BuiltInPredicateStruct(ms).RecreateVariables(vars)
    return Expression(MultiplyStruct(*bif))
}

// ReplaceVariables - Refer to comments in expression.go.
func (ms MultiplyStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(ms).ReplaceVariables(ss)
}  // ReplaceVariables

// String - creates a string representation.
// Returns:  function_name(arg1, arg2, arg3)
func (ms MultiplyStruct) String() string {
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
func (as MultiplyStruct) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {
    result, ok := bifMultiply(as.Arguments, ss)
    if !ok { return ss, false }
    return result.Unify(other, ss)
}

// TermType - returns a constant which identifies this type.
func (as MultiplyStruct) TermType() int { return FUNCTION }
