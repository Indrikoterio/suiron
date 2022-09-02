package suiron

// Join
//
// This built-in function joins strings (Atoms) to form a new string.
// It is used to join words and punctuation.
//
// Words are separated by spaces, but punctuation is attached directly
// to the previous word. For example:
//
//   $D1 = coffee, $D2 = "," , $D3 = tea, $D4 = or, $D5 = juice, $D6 = "?",
//   $X = join($D1, $D2, $D3, $D4, $D5, $D6).
//
// $X is bound to the Atom "coffee, tea or juice?".
//
// Note: All terms must be grounded. If not, the function fails.
//
// Cleve Lendon

import (
    "strings"
    //"fmt"
)

type JoinStruct BuiltInPredicateStruct

// Join - creates a JoinStruct, which holds the function's
// name and arguments. Join requires at least 2 arguments.
// Params: arguments (Unifiable)
// Return: JoinStruct
func Join(arguments ...Unifiable) JoinStruct {
    if len(arguments) < 2 {
        panic("Join - This function requires at least 2 arguments.")
    }
    return JoinStruct {
        Name: "join",
        Arguments: arguments,
    }
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (js JoinStruct) RecreateVariables(vars VarMap) Expression {
    bip := BuiltInPredicateStruct(js).RecreateVariables(vars)
    return Expression(JoinStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (js JoinStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(js).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  function_name(arg1, arg2, arg3)
func (js JoinStruct) String() string {
    return BuiltInPredicateStruct(js).String()
}

// joinWordsAndPunctuation - joins strings (Atoms) of words and
// punctuation together to form a new string.
// Params:
//     list of arguments
//     substitution set
// Return:
//     new Atom
//     success/failure flage
func joinWordsAndPunctuation(arguments []Unifiable, ss SubstitutionSet) (Atom, bool) {

    var sb strings.Builder

    count := 0
    for _, term := range arguments {

        at, ok := ss.CastAtom(term)
        // Should I convert numbers here?
        if !ok { return at, false }

        str := at.String()

        if count > 0 {
            if len(str) == 1 &&
               (str == "," || str == "." ||
                str == "?" || str == "!") {
                sb.WriteString(str)
            } else {
                sb.WriteString(" ")
                sb.WriteString(str)
            }
        } else {
            sb.WriteString(str)
        }
        count++
    }  // for

    return Atom(sb.String()), true

} // joinWordsAndPunctuation


// Unify - unifies the result of a function with another term (usually a variable).
// Params:
//    other unifiable term
//    substitution set
// Returns:
//    updated substitution set
//    success/failure flag
func (js JoinStruct) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {
    result, ok := joinWordsAndPunctuation(js.Arguments, ss)
    if !ok { return ss, false }
    return result.Unify(other, ss)
}

// TermType - returns a constant which identifies this type.
// This function satisfies the Unifiable interface.
func (f JoinStruct) TermType() int { return FUNCTION }
