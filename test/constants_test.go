package main

// Test creation and unification of constants (Atoms, Integers, Floats).
// Cleve Lendon  2022

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestAtom(t *testing.T) {
    fmt.Println("TestAtom")
    a1 := Atom("This is an atom.")
    a2 := Atom("This is an atom.")
    a3 := Atom("Just another.")
    ss := SubstitutionSet{}
    newSS, ok := a1.Unify(a2, ss)
    if !ok { t.Error("Unification must succeed: a1 = a2") }
    newSS, ok = a1.Unify(a3, newSS)
    if ok { t.Error("Unification must fail: a1 != a3") }
    if len(newSS) > 0 { t.Error("Must not change substitution set.") }
}


func TestInteger(t *testing.T) {
    fmt.Println("TestInteger")
    i1 := Integer(45)
    i2 := Integer(45)
    i3 := Integer(46)
    ss := SubstitutionSet{}
    newSS, ok := i1.Unify(i2, ss)
    if !ok { t.Error("Unification must succeed: i1 = i2") }
    newSS, ok = i1.Unify(i3, newSS)
    if ok { t.Error("Unification must fail: i1 != i3") }
    if len(newSS) > 0 { t.Error("Must not change substitution set.") }
}

func TestFloat(t *testing.T) {
    fmt.Println("TestFloat")
    a1 := Atom("An atom")
    i1 := Integer(45)
    f1 := Float(45.0)
    f2 := Float(45.0)
    f3 := Float(45.0000000001)
    ss := SubstitutionSet{}
    newSS, ok := f1.Unify(f2, ss)
    if !ok { t.Error("Unification must succeed: f1 = f2") }
    newSS, ok = f1.Unify(f3, newSS)
    if ok { t.Error("Unification must fail: f1 != f3") }
    newSS, ok = f1.Unify(a1, newSS)
    if ok { t.Error("Unification must fail: f1 != a1") }
    newSS, ok = f1.Unify(i1, newSS)
    if ok { t.Error("Unification must fail: f1 != i1") }
    if len(newSS) > 0 { t.Error("Must not change substitution set.") }
}

