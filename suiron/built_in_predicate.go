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
// This method assists the method RecreateVariables(). If the argument is
// a Variable, this method will substitute it with a new unique variable.
// If not, it just returns the argument as is.
//
// Params: unifiable argument
//         list of previously recreated variables
// Return: out argument
func recreateOneVar(argument Unifiable, vars map[Variable]Variable) Unifiable {
    tt := argument.TermType()
    if tt == VARIABLE {
        arg, _ := argument.(Variable)  // cast it
        return arg.RecreateVariables(vars).(Unifiable)
    } else if tt == LINKEDLIST {
        arg, _ := argument.(LinkedListStruct)  // cast it
        return arg.RecreateVariables(vars).(Unifiable)
    } else if tt == FUNCTION {
        arg, _ := argument.(Function)  // cast it
        return arg.RecreateVariables(vars).(Unifiable)
    }
    return argument
}

// RecreateVariables - In Prolog, and in this inference engine, the scope of
// a logic variable is the rule or goal in which it is defined. When the
// algorithm tries to solve a goal, it calls this method to ensure that the
// variables are unique. See comments in expression.go.
func (bips BuiltInPredicateStruct) RecreateVariables(
                       vars map[Variable]Variable) BuiltInPredicateStruct {
    newArguments := []Unifiable{}
    for i := 0; i < len(bips.Arguments); i++ {
        v := recreateOneVar(bips.Arguments[i], vars)
        newArguments = append(newArguments, v)
    }
    newBIP := BuiltInPredicateStruct{
               Name: bips.Name,
               Arguments: newArguments,
           }
    return newBIP
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
            theVar := arg.(Variable)
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
