package suiron

// BuiltInPredicate - This file contains the struct and
// functions which are common to all of Suiron's built-in
// predicates and functions.
//
// BuiltInPredicateStruct records the name of the function/
// predicate, and holds its arguments, a slice of Unifiables.
//
// RecreateVariables(), ReplaceVariables(), and String() are
// called by predicates and functions to satisfy the Expression
// interface. See expression.go for details.
//
// Cleve Lendon

import (
    "strings"
    //"fmt"
)

type BuiltInPredicateStruct struct {
    Name      string
    Arguments []Unifiable
}

// recreateOneVar - recreates one variable.
//
// This function assists RecreateVariables(). If the given term is a Variable,
// recreateOneVar() will return a new unique variable. If the term is a linked
// list, recreateOneVar() will call its RecreateVariables() function, to recreate
// all the elements of the linked list. Otherwise, it will return the term as is.
//
// Params: Unifiable term
//         list of previously recreated Variables
// Return: tew Unifiable term
func recreateOneVar(term Unifiable, vars VarMap) Unifiable {
    tt := term.TermType()
    if tt == VARIABLE {
        t, _ := term.(VariableStruct)  // cast it
        return t.RecreateVariables(vars).(Unifiable)
    } else if tt == COMPLEX {
        t, _ := term.(Complex)  // cast it
        return t.RecreateVariables(vars).(Unifiable)
    } else if tt == LINKEDLIST {
        t, _ := term.(LinkedListStruct)  // cast it
        return t.RecreateVariables(vars).(Unifiable)
    } else if tt == FUNCTION {
        t, _ := term.(Function)  // cast it
        return t.RecreateVariables(vars).(Unifiable)
    }
    return term
} // recreateOneVar()


// recreateVars - recreates all Variables in a slice of Unifiable terms.
// Params: slice of Unifiable terms
//         map of previously recreated Variables
// Return: slice of new terms
func recreateVars(terms []Unifiable, vars VarMap) []Unifiable {
    newTerms := []Unifiable{}
    for _, term := range terms {
        v := recreateOneVar(term, vars)
        newTerms = append(newTerms, v)
    }
    return newTerms
} // recreateVars()

// RecreateVariables - The scope of a logic variable is the rule in which
// it is defined. When the algorithm tries to solve a goal, it calls this
// function to ensure that the variables are unique.
// See comments in expression.go.
func (bips BuiltInPredicateStruct) RecreateVariables(vars VarMap) *BuiltInPredicateStruct {
    newArguments := recreateVars(bips.Arguments, vars)
    ptrBIP := new(BuiltInPredicateStruct)
    ptrBIP.Name = bips.Name
    ptrBIP.Arguments = newArguments
    return ptrBIP
}

// ReplaceVariables - replaces a bound variable with its binding.
// This method is used for displaying final results.
// Refer to comments in expression.go.
func (bips BuiltInPredicateStruct) ReplaceVariables(ss SubstitutionSet) Expression {

    for _, arg := range bips.Arguments {
        tt := arg.TermType()
        if tt == ATOM || tt == INTEGER || tt == FLOAT || tt == COMPLEX {
            return arg
        } else if tt == VARIABLE {
            theVar := arg.(VariableStruct)
            if ss.IsBound(theVar) {
                exp := theVar.ReplaceVariables(ss)
                return exp
            } else {
               return Atom("BIP error")
            } 
        } else {
            return Atom("BIP error 2")
        }
    }
    return nil

}  // ReplaceVariables


// String - creates a string for debugging purposes.
// Format:  PredName(arg1, arg2, arg3)
func (bips BuiltInPredicateStruct) String() string {
    var sb strings.Builder
    sb.WriteString(bips.Name)
    sb.WriteString("(")
    for n, k := range bips.Arguments {
        if n != 0 { sb.WriteString(", ") }
        sb.WriteString(k.String())
    }
    sb.WriteString(")")
    return sb.String()
}
