//package main
package suiron

// BIFTemplate
//
// This file is a template for writing built-in functions (BIF) for Suiron.
//
// Search and replace the string 'BIFTemplate', everywhere it appears,
// with the name of your function. Write your function specific code
// in bifEvaluate(), and rename bifEvaluate to something meaningful.
// Adjust comments appropriately and rename this file.
//
// Cleve Lendon 2022

import (
    // Uncomment this import if the function
    // is outside of the suiron package.
    //. "github.com/indrikoterio/suiron/suiron"
)

// Built-in functions and built-in predicates use the same struct.
type BIFTemplateStruct BuiltInPredicateStruct

// BIFTemplate - creates the struct which defines this built-in function.
// Checks input arguments.
func BIFTemplate(arguments ...Unifiable) BIFTemplateStruct {
    if len(arguments) < 2 {
        panic("BIFTemplate - requires at least 2 arguments.")
    }
    return BIFTemplateStruct {
        Name: "BIFTemplate",
        Arguments: arguments,
    }
}

//----------------------------------------------------------------
// bifEvaluate - does the logical work of this built-in function.
// Params:
//     list of arguments
//     substitution set
// Returns:
//    new unifiable
//    success/failure flag
func bifEvaluate(arguments []Unifiable, ss SubstitutionSet) (Unifiable, bool) {
    // Do something!
    return Atom("Return the result."), true
} // bifEvaluate

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (s BIFTemplateStruct) RecreateVariables(vars VarMap) Expression {
    bif := BuiltInPredicateStruct(s).RecreateVariables(vars)
    return Expression(BIFTemplateStruct(*bif))
}

// ReplaceVariables - Refer to comments in expression.go.
func (s BIFTemplateStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(s).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  function_name(arg1, arg2, arg3)
func (s BIFTemplateStruct) String() string {
    return BuiltInPredicateStruct(s).String()
}

//----------------------------------------------------------------
// Unify() and TermType() satisfy the Unifiable interface.
//----------------------------------------------------------------

// Unify - unifies the result of a function with another term (usually a variable).
// Params:
//    other unifiable term
//    substitution set
// Returns:
//    updated substitution set
//    success/failure flag
func (s BIFTemplateStruct) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {
    result, ok := bifEvaluate(s.Arguments, ss)
    if !ok { return ss, false }
    return result.Unify(other, ss)
}

// TermType - returns a constant which identifies this type.
func (s BIFTemplateStruct) TermType() int { return FUNCTION }
