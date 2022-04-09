package suiron

// Variable - defines a logic variable for this inference engine.
// Prolog variables are represented by strings which begin with a capital letter,
// eg. X, Y, Noun. In this inference engine, a Variable begins with a dollar sign
// and a letter: $X, $Y, $Noun. This was done so that constants which begin with
// a capital letter do not have to be put inside quote marks. 
// (ie.: Harold, instead of "Harold".)
//
// Cleve Lendon

import (
    "strings"
    "unicode"
    "fmt"
)

type Variable struct {
    name   string
    id     int
}

// String - return this term as a string.
func (v Variable) String() string {
    if v.id == 0 { return v.name }
    return fmt.Sprint(v.name, "_", v.id)
}

var variableId int

// LogicVar - Factory function to create a logic Variable from a string.
// The variable must begin with a dollar sign and a letter. Eg. $X
// If it does not, a error is produced.
func LogicVar(str string) (Variable, error) {
    name := strings.TrimSpace(str)
    r := []rune(name)
    if len(r) < 2 {
        err := makeVariableError("variable must start with $ and letter", name)
        return Variable{ name: name, id: 0 }, err
    }
    first  := r[0]
    second := r[1]
    if first != '$' {
        err := makeVariableError("variable must start with $", name)
        return Variable{ name: name, id: 0 }, err
    }
    if !unicode.IsLetter(second) {
        err := makeVariableError("second character must be a letter", name)
        return Variable{ name: name, id: 0 }, err
    }
    return Variable{ name: name, id: 0 }, nil
}

// makeVariableError - creates an error for LogicVar().
// msg - error message
// str - string which caused the error
func makeVariableError(msg string, str string) error {
    return fmt.Errorf("LogicVar() - %v: >%v<\n", msg, str)
}

// CopyVar - creates a copy of the given Variable.
// The new variable has the same name as the copied variable,
// but the ID is different.
// There is no need to check validity, because the original
// has already been validated.
func CopyVar(v Variable) Variable {
    variableId += 1
    return Variable{ name: v.name, id: variableId }
}

// TermType - Returns an integer constant which identifies the type.
func (v Variable) TermType() int { return VARIABLE }

// Unify - Unifies this variable with another unifiable expression
// (if this variable is not already bound). Please refer to unifiable.go.
func (v Variable) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

    otherType := other.TermType()
    if otherType == VARIABLE {
        // A variable unifies with itself.
        if v.id == other.(Variable).id {
            if v.id == 0 {
                if v.name == other.(Variable).name { return ss, true }
            } else {
                return ss, true
            }
        }
    }

    // The unify method of a function evaluates the function, so
    // if the other expression is a function, call its unify method.
    if otherType == FUNCTION { return other.Unify(v, ss) }

    u, found := ss[v]
    if found {
        return u.Unify(other, ss)
    }

    newSS := ss.Copy()
    newSS[v] = other
    return newSS, true
} // Unify


// RecreateVariables - In Prolog, and in this inference engine, the scope of
// a logic variable is the rule or goal in which it is defined. When the
// algorithm tries to solve a goal, it calls this method to ensure that the
// variables are unique. See comments in expression.go.
// Note: This method creates variables from previously validated variables,
// so there is no need to validate the variable name by calling LogicVar().
func (v Variable) RecreateVariables(vars map[Variable]Variable) Expression {
    var newVar Variable
    var ok bool
    if newVar, ok = vars[v]; !ok {
        // Name has already been validated. No need to call LogicVar().
        variableId += 1
        newVar = Variable{ name: v.name, id: variableId }
        vars[v] = newVar
    }
    return Expression(newVar)
}

// ReplaceVariables - replaces a bound variable with its binding.
// This method is used for displaying final results.
// Refer to comments in expression.go.
func (v Variable) ReplaceVariables(ss SubstitutionSet) Expression {
    u, found := ss[v]
    if found {
        // Recursive
        return u.ReplaceVariables(ss)
    } else {
        return Expression(v)
    }
} // ReplaceVariables()
