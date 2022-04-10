package main

// Capitalize
//
// This Suiron function capitalizes a word or name.
//
// It accepts one argument, an Atom, or a Variable bound
// to an Atom, and capitalizes it. (tokyo becomes Tokyo.)
//
// Note: Unlike a predicate, which returns a substitution
// set, a function returns a Uunifiable term. Therefore,
// it must be used with unification. Eg.
//
//    ..., $CapName = capitalize($Name),...
//
// Cleve Lendon

import (
    "strings"
    . "github.com/indrikoterio/suiron/suiron"
)

// Built-in functions and built-in predicates use the same struct.
type CapitalizeStruct BuiltInPredicateStruct

// Capitalize - creates the struct which defines this built-in
// function. Checks input arguments.
func Capitalize(arguments ...Unifiable) CapitalizeStruct {
    if len(arguments) != 1 {
        panic("Capitalize - requires 1 argument.")
    }
    return CapitalizeStruct {
        Name: "capitalize",
        Arguments: arguments,
    }
}

//----------------------------------------------------------------
// bifCapitalize - Capitalizes the first letter of the given word.
//
// Params:
//     list of arguments
//     substitution set
// Returns:
//    new unifiable
//    success/failure flag
//
func bifCapitalize(arguments []Unifiable, ss SubstitutionSet) (Unifiable, bool) {
    term, ok := ss.CastAtom(arguments[0])
    if !ok { return Atom("?"), false }
    str := term.String()
    str = strings.Title(strings.ToLower(str))
    return Atom(str), true
} // bifCapitalize

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (s CapitalizeStruct) RecreateVariables(
                               vars map[string]VariableStruct) Expression {
    bif := BuiltInPredicateStruct(s).RecreateVariables(vars)
    return Expression(CapitalizeStruct(*bif))
}

// ReplaceVariables - Refer to comments in expression.go.
func (s CapitalizeStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(s).ReplaceVariables(ss)
}  // ReplaceVariables

// String - creates a string representation.
// Returns:  function_name(arg1, arg2, arg3)
func (s CapitalizeStruct) String() string {
    return BuiltInPredicateStruct(s).String()
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
func (s CapitalizeStruct) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {
    result, ok := bifCapitalize(s.Arguments, ss)
    if !ok { return ss, false }
    return result.Unify(other, ss)
}

// TermType - returns a constant which identifies this type.
func (s CapitalizeStruct) TermType() int { return FUNCTION }
