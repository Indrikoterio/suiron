package main

// Tests 'print' predicate.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestPrint(t *testing.T) {

    fmt.Println("TestPrint 1:")
    fmt.Println("Persian, king, [Cyrus, Cambysis, Darius]")

    persian := Atom("Persian")
    king    := Atom("king")
    K1 := Atom("Cyrus")
    K2 := Atom("Cambysis")
    K3 := Atom("Darius")
    print_test := Atom("print_test")

    list := MakeLinkedList(false, K1, K2, K3)

    X, _ := LogicVar("$X")
    Y, _ := LogicVar("$Y")

    ss := SubstitutionSet{}

    // print_test :- $X = king, $Y = [Cyrus, Cambysis, Darius], print(persian, $X, $Y), nl.
    head := Complex{print_test}
    c2   := Unify(X, king)
    c3   := Unify(Y, list)
    c4   := Print(persian, X, Y)
    body := And(c2, c3, c4, NL())

    // Set up the knowledge base.
    kb := KnowledgeBase{}
    r1 := Rule(head, body)

    kb.Add(r1)

    goal := MakeGoal(print_test)
    Solve(goal, kb, ss)

    fmt.Println("TestPrint 2:")
    fmt.Println("Hello World, my name is Cleve.")

    world := Atom("World")
    cleve := Atom("Cleve")

    format_string := Atom("Hello %s, my name is %s.\n")
    print_test2 := Atom("print_test2")
    c5 := Complex{print_test2}
    c6 := Print(format_string, world, cleve)
    r2 := Rule(c5, c6)
    kb.Add(r2)

    goal = MakeGoal(print_test2)
    Solve(goal, kb, ss)

} // TestPrint
