package main

// Tests the Fail predicate.
//
// Fail causes the rule to fail. It is used to force backtracking.
//
// count(1).
// count(2).
// count(3).
// test :- count($X), print($X), fail.
// test :- nl, fail.
//
// This rule should print 1, 2 and 3 on separate lines.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestFail(t *testing.T) {

    fmt.Print("TestFail - ")

    count := Atom("count")
    test  := Atom("test")
    X, _  := LogicVar("$X")

    // Set up the knowledge base.
    kb := KnowledgeBase{}
    ss := SubstitutionSet{}

    f1 := Fact(Complex{count, Integer(1)})
    f2 := Fact(Complex{count, Integer(2)})
    f3 := Fact(Complex{count, Integer(3)})

    kb.Add(f1, f2, f3)

    head := Complex{test}
    body := And(
                Complex{count, X}, Print(Atom("%s "), X), Fail(),
            )
    r1 := Rule(head, body)
    kb.Add(r1)

    body2 := And( NL(), Fail() )
    r2 := Rule(head, body2)
    kb.Add(r2)

    goal := Complex{test}
    _, failure := SolveAll(goal, kb, ss)
    if failure == "" {
        t.Error("TestFail - " + failure)
    }

} // TestJoin
