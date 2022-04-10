package suiron

// This file defines constants (Atoms, Integers, Floats) for Suiron.
// Atoms are equivalent to strings. Integers and Floats are 64-bit.
// Cleve Lendon

import (
    "fmt"
)

// Atom is equivalent to a string.
// In this inference engine, an atom can start with an upper case
// or a lower case letter. (Unlike Prolog.)
type Atom string

// TermType - Returns an integer constant which identifies this type.
func (a Atom) TermType() int { return ATOM }

// Unify - unifies an Atom with another term. If both terms are Atoms,
// and equal, then Unify succeeds. If they are not equal, Unify fails.
// If one of the terms is an unbound Variable, then Unify binds the
// Variable to the Atom, records the binding in the substitution set,
// and returns with success. If the Variable is already bound to a
// different Atom, Unify will fail.
func (a Atom) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {
    otherType := other.TermType()
    if otherType == ATOM {
        if a == other.(Atom) { return ss, true }  // success
        return ss, false   // failure
    }
    if otherType == VARIABLE { return other.Unify(a, ss) }
    if otherType == ANONYMOUS { return ss, true } // success
    return ss, false // failure
}

// String - return this term as a string.
func (a Atom) String() string { return string(a) }

// RecreateVariables - creates unique variables every time the
// inference engine fetches a rule from the knowledge base.
// A constant is not a variable, so this function simply returns
// the constant. This function satisfies the Expression interface.
func (a Atom) RecreateVariables(m map[string]VariableStruct) Expression {
    return a;
}

// ReplaceVariables() is called after a solution has been found.
// It replaces logic variables with the constants which they are
// bound to, in order to display results.
// For constants, ReplaceVariables() simply returns the constant.
// This function satisfies the Expression interface.
func (a Atom) ReplaceVariables(ss SubstitutionSet) Expression {
    return a;
}

//--------------------------------------------
// Define a 64 bit integer.
type Integer int64

// TermType - returns an integer constant which identifies this type.
func (i Integer) TermType() int { return INTEGER }

// Unify - unifies an Integer with another term. If both terms are
// Integers, and equal, then Unify succeeds. If they are not equal,
// Unify fails. If one of the terms is an unbound Variable, then Unify
// binds the Variable to the Integer, records the binding in the
// substitution set, and returns with success. If the Variable is
// already bound to a different Integer, Unify will fail.
func (i Integer) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {
    otherType := other.TermType()
    if otherType == INTEGER {
        if i == other.(Integer) { return ss, true }  // success
        return ss, false   // failure
    }
    if otherType == VARIABLE { return other.Unify(i, ss) }
    if otherType == ANONYMOUS { return ss, true } // success
    return ss, false // failure
}

// String - return this term as a string.
func (i Integer) String() string { return fmt.Sprintf("%d", i) }

// RecreateVariables - creates unique variables every time the
// inference engine fetches a rule from the knowledge base.
// A constant is not a variable, so this function simply returns
// the constant. This function satisfies the Expression interface.
func (i Integer) RecreateVariables(m map[string]VariableStruct) Expression {
    return i;
}

// ReplaceVariables() is called after a solution has been found.
// It replaces logic variables with the constants which they are
// bound to, in order to display results.
// For constants, ReplaceVariables() simply returns the constant.
// This function satisfies the Expression interface.
func (i Integer) ReplaceVariables(ss SubstitutionSet) Expression {
    return i;
}


//--------------------------------------------
// Define a 64 bit floating point number.
type Float float64

// TermType - Returns an integer constant which identifies this type.
func (f Float) TermType() int { return FLOAT }

// Unify - unifies a floating point number with another term.
// If both terms are floating point numbers, and equal, then
// Unify succeeds. If they are not equal, Unify fails.
// If one of the terms is an unbound Variable, then Unify binds
// the Variable to the floating point number, records the binding
// in the substitution set, and returns with success.
// If the Variable is already bound to a different floating point
// number, Unify will fail.
func (f Float) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {
    otherType := other.TermType()
    if otherType == FLOAT {
        if f == other.(Float) { return ss, true }  // success
        return ss, false   // failure
    }
    if otherType == VARIABLE { return other.Unify(f, ss) }
    if otherType == ANONYMOUS { return ss, true } // success
    return ss, false // failure
}

// String - return this term as a string.
func (f Float) String() string { return fmt.Sprintf("%f", f) }

// RecreateVariables - creates unique variables every time the
// inference engine fetches a rule from the knowledge base.
// A constant is not a variable, so this function simply returns
// the constant. This function satisfies the Expression interface.
func (f Float) RecreateVariables(m map[string]VariableStruct) Expression {
    return f;
}

// ReplaceVariables() is called after a solution has been found.
// It replaces logic variables with the constants which they are
// bound to, in order to display results.
// For constants, ReplaceVariables() simply returns the constant.
// This function satisfies the Expression interface.
func (f Float) ReplaceVariables(ss SubstitutionSet) Expression {
    return f;
}


