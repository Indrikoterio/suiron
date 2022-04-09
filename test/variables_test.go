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

    fmt.Println("TestVariables")
    W, _ := LogicVar("$W")
    X, _ := LogicVar("$X")
    Y, _ := LogicVar("$Y")
    Z, _ := LogicVar("$Z")
    age := Integer(43)
    pi  := Float(3.14159)
    ss  := SubstitutionSet{}

    pronoun  := Atom("pronoun")
    me       := Atom("me")
    first    := Atom("first")
    sing     := Atom("singular")
    acc      := Atom("accusative")

    person, _    := LogicVar("$Person")
    plurality, _ := LogicVar("$Plurality")
    case_, _     := LogicVar("$Case")

    c1 := Complex{pronoun, me, first, sing, acc}
    c2 := Complex{pronoun, me, person, plurality, case_}

    newSS, ok := X.Unify(X, ss)
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

    expected := `
----- Bindings -----
    $Case: accusative
    $Person: first
    $Plurality: singular
    $W: 3.141590
    $X: $Y
    $Y: pronoun
    $Z: 43
--------------------
`
    actual := newSS.String()

    if actual != expected {
        t.Error("TestVariables, error in substitution set.\n" + actual + expected)
    }

}
