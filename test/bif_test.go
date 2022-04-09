package main

// Test the built-in function feature.
//
// Suiron provides a way for programmers to write their own
// 'built-in' functions. (That is, a way to write functions for
// Suiron's rule declaration language, in Go. bif_template.go
// can be used as a template for this purpose.)
//
// The function to be tested here is capitalize(...). Please
// refer to comments in capitalize.go for a description of how
// this function works. 
//
// The following rule will be written to the knowledgebase:
//
//   test($In, $Out) :- capitalize($In) = $Out.
//
// The goal to be tested is:
//
//   test(london, $X).
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

// Set up facts and rules.
func TestBuiltInFunction(t *testing.T) {

    fmt.Println("TestBuiltInFunction")
    kb := KnowledgeBase{}

    // Create logic variables.
    X, _  := LogicVar("$X")

    In, _     := LogicVar("$In")
    Out, _    := LogicVar("$Out")

    test  := Atom("test")

    c1 := Complex{test, In, Out}
    c2 := Unify(Capitalize(In), Out)
    r1 := Rule(c1, c2)

    kb.Add(r1)  // Add rule to knowledgebase.

    goal := Complex{test, Atom("london"), X}
    solution, failure := Solve(goal, kb, SubstitutionSet{})

    if len(failure) != 0 {
        t.Error("TestBuiltInFunction - " + failure)
        return
    }

    expected := "London"
    actual := solution.GetTerm(2).String()

    if actual != expected {
        t.Error("\nTestBuiltInFunction - Expected: " + expected +
                "\n                           Was: " + actual)
        return
    }

}  // TestBuiltInFunction
