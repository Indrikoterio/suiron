package main

// Test the 'built-in predicate' functionality.
//
// Suiron has built-in predicates, such as 'append', 'nl' and 'print'.
// It also provides a mechanism to allow programmers to write their own
// 'built-in' predicates. (That is, predicates written in Go, which can
// be called as part of Suiron's rule declaration language. The file
// bip_template.go can be used as a template for this purpose.)
//
// The predicate to be tested here is hyphenate(...). Please refer to
// comments in hyphenate.go for a description of how this predicate works. 
//
// In order to test this functionality, the following rules will be written
// to the knowledgebase:
//
//   join_all($In, $Out, $InErr, $OutErr) :- hyphenate($In, $H, $T, $InErr, $Err2),
//                                              join_all([$H | $T], $Out, $Err2, $OutErr).
//   join_all([$H], $H, $X, $X).
//
//   bip_test($Out, $OutErr) :- join_all([sister, in, law], $Out, [first error], $OutErr).
//
// Suiron will solve for $Out and $OutErr.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

// Set up facts and rules.
func TestBuiltInPredicate(t *testing.T) {

    fmt.Println("TestBuiltInPredicate")
    kb := KnowledgeBase{}

    // Create logic variables.
    X, _  := LogicVar("$X")
    Y, _  := LogicVar("$Y")
    H, _  := LogicVar("$H")
    T, _  := LogicVar("$T")

    In, _     := LogicVar("$In")
    InErr, _  := LogicVar("$InErr")
    Out, _    := LogicVar("$Out")
    OutErr, _ := LogicVar("$OutErr")
    Err2, _   := LogicVar("$Err2")

    bip_test  := Atom("bip_test")
    join_all  := Atom("join_all")

    c1 := Complex{join_all, In, Out, InErr, OutErr}
    c2 := Hyphenate(In, H, T, InErr, Err2)
    ll := MakeLinkedList(true, H, T)
    c3 := Complex{join_all, ll, Out, Err2, OutErr}
    body := And(c2, c3)
    r1 := Rule(c1, body)

    l2 := MakeLinkedList(false, H)
    c4 := Complex{join_all, l2, H, X, X}
    r2 := Fact(c4)

    // bip_test($Out, $OutErr) :- join_all([sister, in, law], $Out, [first error], $OutErr).
    c7 := Complex{bip_test, Out, OutErr}
    l3 := MakeLinkedList(false, Atom("sister"), Atom("in"), Atom("law"))
    l4 := MakeLinkedList(false, Atom("first error"))
    c8 := Complex{join_all, l3, Out, l4, OutErr}
    r3 := Rule(c7, c8)

    kb.Add(r1, r2, r3)  // Add rules to knowledgebase.

    // Show the knowledgebase.
    //DBKB(kb)

    goal := Complex{bip_test, X, Y}
    solution, failure := Solve(goal, kb, SubstitutionSet{})

    if len(failure) != 0 {
        t.Error("TestBuiltInPredicate - " + failure)
        return
    }

    // Check the solutions of bip_test($Out, $OutError).
    expected := "sister-in-law [another error, another error, first error]"

    out    := solution.GetTerm(1)
    outErr := solution.GetTerm(2)
    actual := fmt.Sprintf("%v %v", out, outErr)

    if actual != expected {
        t.Error("\nTestBuiltInPredicate - Expected: " + expected +
                "\n                            Was: " + actual)
        return
    }

}  // TestBuiltInPredicate
