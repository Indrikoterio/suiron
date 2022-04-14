package suiron

// Variable - defines a logic variable for this inference engine.
// Prolog variables are represented by strings which begin with a capital letter,
// eg. X, Y, Noun. In this inference engine, a Variable begins with a dollar sign
// and a letter: $X, $Y, $Noun. This was done so that constants which begin with
// a capital letter do not have to be put inside quote marks. 
// (ie.: Harold, instead of "Harold".)
//
// Cleve Lendon

import "unsafe"

//#include <string.h>
import "C"

import (
    "strings"
    "unicode"
    "fmt"
)

type VariableStruct struct {
    id     int
    name   string
}

// String - return this term as a string.
func (v VariableStruct) String() string {
    if v.id == 0 { return v.name }
    return fmt.Sprint(v.name, "_", v.id)
}

var variableId int

// LogicVar - Factory function to create a logic Variable from a string.
// The variable must begin with a dollar sign and a letter. Eg. $X
// If it does not, a error is produced.
func LogicVar(str string) (VariableStruct, error) {
    name := strings.TrimSpace(str)
    r := []rune(name)
    if len(r) < 2 {
        err := makeVariableError("variable must start with $ and letter", name)
        return VariableStruct{ name: name, id: 0 }, err
    }
    first  := r[0]
    second := r[1]
    if first != '$' {
        err := makeVariableError("variable must start with $", name)
        return VariableStruct{ name: name, id: 0 }, err
    }
    if !unicode.IsLetter(second) {
        err := makeVariableError("second character must be a letter", name)
        return VariableStruct{ name: name, id: 0 }, err
    }
    return VariableStruct{ name: name, id: 0 }, nil
}

// makeVariableError - creates an error for LogicVar().
// msg - error message
// str - string which caused the error
func makeVariableError(msg string, str string) error {
    return fmt.Errorf("LogicVar() - %v: >%v<\n", msg, str)
}

// CopyVar - creates a copy of the given VariableStruct.
// The new variable has the same name as the copied variable,
// but the ID is different.
// There is no need to check validity, because the original
// has already been validated.
func CopyVar(v VariableStruct) VariableStruct {
    variableId += 1
    return VariableStruct{ name: v.name, id: variableId }
}

// TermType - Returns an integer constant which identifies the type.
func (v VariableStruct) TermType() int { return VARIABLE }

// ID - Returns the ID number of the given variable.
func (v VariableStruct) ID() int { return v.id }

// Unify - Unifies this variable with another unifiable expression
// (if this variable is not already bound). Please refer to unifiable.go.
func (v VariableStruct) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

    // LogicVar() creates variables with an ID of 0, but variables
    // are recreated (given a new ID) when a rule is fetched from the
    // knowledge base.
    // If a variable has an ID of 0 here, something was done incorrectly.
    // The following statement avoids endless loops, which occur when the
    // substitution set has a variable with ID = 0 at location 0.
    if v.id == 0 { return ss, false }

    otherType := other.TermType()

    if otherType == VARIABLE {
        // A variable unifies with itself.
        if v.id == other.(VariableStruct).id { return ss, true }
    }

    // The unify method of a function evaluates the function, so
    // if the other expression is a function, call its unify method.
    if otherType == FUNCTION { return other.Unify(v, ss) }
 
    var u *Unifiable
    lengthSrc := len(ss)
    if v.id < lengthSrc && ss[v.id] != nil {
        u = ss[v.id]
        return (*u).Unify(other, ss)
    }

    lengthDst := lengthSrc
    if v.id >= lengthDst { lengthDst = v.id + 1 }
    newSS := make(SubstitutionSet, lengthDst)

    //copy(newSS, ss)
    //copy(unsafe.Slice((**Unifiable)(unsafe.Pointer(&newSS[0])), lengthSrc), ss)

    // As fast as possible.
    ptrSize := (int)(unsafe.Sizeof(&ss))
    if lengthSrc > 0 {
        C.memcpy(unsafe.Pointer(&newSS[0]),
                 unsafe.Pointer(&ss[0]),
                 C.size_t(lengthSrc * ptrSize))
    }

    newSS[v.id] = &other
    return newSS, true

} // Unify


// RecreateVariables - The scope of a logic variable is the rule or goal in
// which it is defined. When the algorithm tries to solve a goal, it calls
// this method to ensure that the variables are unique.
// See comments in expression.go.
// Note: This method creates variables from previously validated variables,
// so there is no need to validate the variable name by calling LogicVar().
// Params: map of previously recreated variables
// Return: new variable (as Expression)
func (v VariableStruct) RecreateVariables(vars VarMap) Expression {
    var newVar VariableStruct
    var ok bool
    strVar := v.String()
    if newVar, ok = vars[strVar]; !ok {
        // Name has already been validated. No need to call LogicVar().
        variableId++
        newVar = VariableStruct{ name: v.name, id: variableId }
        vars[strVar] = newVar
    }
    return Expression(newVar)
} // RecreateVariables()

// ReplaceVariables - replaces a bound variable with its binding.
// This method is used for displaying final results.
// Refer to comments in expression.go.
func (v VariableStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    var u *Unifiable
    if v.id < len(ss) && ss[v.id] != nil {
        u = ss[v.id]
        return (*u).ReplaceVariables(ss)
    } else {
        return Expression(v)
    }
} // ReplaceVariables()
