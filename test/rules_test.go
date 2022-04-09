package main

// Test creation of logic rules.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestRules(t *testing.T) {

    fmt.Println("TestRules")
    // grandparent($X, $Y) :- parent($X, $Z), parent($Z, $Y).
    grandparent := Atom("grandparent")
    parent      := Atom("parent")
    X, _ := LogicVar("$X")
    Y, _ := LogicVar("$Y")
    Z, _ := LogicVar("$Z")

    c1 := Complex{grandparent, X, Y}
    c2 := Complex{parent, X, Z}
    c3 := Complex{parent, Z, Y}

    andOp := And(c2, c3)
    r1 := Rule(c1, andOp)

    expected := "grandparent($X, $Y) :- parent($X, $Z), parent($Z, $Y)."
    actual := r1.String()
    if actual != expected {
        t.Error("\nRule should be: " + expected + "\nWas:            " + actual)
    }
}
