package main

// Test creation and unification of Variables.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestVariables(t *testing.T) {

    // vars - Keeps track of previously recreated variables.
    vars := VarMap{}

    fmt.Println("TestVariables")
    W, _ := LogicVar("$W")
    W = W.RecreateVariables(vars).(VariableStruct)
    X, _ := LogicVar("$X")
    X = X.RecreateVariables(vars).(VariableStruct)
    Y, _ := LogicVar("$Y")
    Y = Y.RecreateVariables(vars).(VariableStruct)
    Z, _ := LogicVar("$Z")
    Z = Z.RecreateVariables(vars).(VariableStruct)
    age := Integer(43)
    pi  := Float(3.14159)

    pronoun  := Atom("pronoun")
    me       := Atom("me")
    first    := Atom("first")
    sing     := Atom("singular")
    acc      := Atom("accusative")

    person, _    := LogicVar("$Person")
    plurality, _ := LogicVar("$Plurality")
    case_, _     := LogicVar("$Case")

    c1 := MakeQuery(pronoun, me, first, sing, acc)
    c2 := MakeQuery(pronoun, me, person, plurality, case_)

    newSS, ok := X.Unify(X, SubstitutionSet{})
    if !ok { t.Error("TestVariables - unification should succeed: $X = $X") }

    newSS, ok = X.Unify(Y, newSS)
    if !ok { t.Error("TestVariables - unification should succeed: $X = $Y") }

    newSS, ok = Y.Unify(pronoun, newSS)
    if !ok { t.Error("TestVariables - unification should succeed: $Y = pronoun") }

    newSS, ok = X.Unify(pronoun, newSS)
    if !ok { t.Error("TestVariables - unification should succeed: $X = pronoun") }

    newSS, ok = Z.Unify(age, newSS)
    if !ok { t.Error("TestVariables - unification should succeed: $Z = 43") }

    newSS, ok = W.Unify(pi, newSS)
    if !ok { t.Error("TestVariables - unification should succeed: $W = 3.14159") }

    newSS, ok = Z.Unify(W, newSS)
    if ok { t.Error("TestVariables - unification should not succeed: $Z = $W") }

    // Unify complex terms.
    newSS, ok = c1.Unify(c2, newSS)
    if !ok { t.Error("TestVariables - unification should succeed: c1 = c2") }

    expected := "pronoun(me, first, singular, accusative)"
    actual := c2.ReplaceVariables(newSS).(Complex).String()

    if actual != expected {
        t.Error("TestVariables - failed to unify complex terms.\n" +
                actual + "\n" + expected)
    }
}
