package main

// Tests the 'join' function.
//
// 'Join' is a built-in function which joins a list of words and
// and punctuation to form a single string (= Atom). In the following
// source example,
//
//   $D1 = coffee, $D2 = "," , $D3 = tea, $D4 = or, $D5 = juice, $D6 = "?",
//   $X = join($D1, $D2, $D3, $D4, $D5, $D6).
//
// $X is bound to the Atom "coffee, tea or juice?".
//
// A built-in function is different from a built-in predicate, in that
// a built-in function returns a value (Atom, Integer or Float), which
// must be unified with something in order to be useful. All the
// arguments of a function must be constants or grounded variables.
// If not, the function fails.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestJoin(t *testing.T) {

    fmt.Println("TestJoin")

    D1, _  := LogicVar("$D1")
    D2, _  := LogicVar("$D2")
    D3, _  := LogicVar("$D3")
    D4, _  := LogicVar("$D4")
    D5, _  := LogicVar("$D5")
    D6, _  := LogicVar("$D6")
    Out, _ := LogicVar("$Out")

    would_you_like := Atom("Would you like...")

    coffee  := Atom("coffee")
    comma   := Atom(",")
    tea     := Atom("tea")
    or      := Atom("or")
    juice   := Atom("juice")
    question_mark := Atom("?")

    u1 := Unify(D1, coffee)
    u2 := Unify(D2, comma)
    u3 := Unify(D3, tea)
    u4 := Unify(D4, or)
    u5 := Unify(D5, juice)
    u6 := Unify(D6, question_mark)
    u7 := Unify(Out, Join(D1, D2, D3, D4, D5, D6))

    // Make rule.
    head := Complex{would_you_like, Out}
    body := And(u1, u2, u3, u4, u5, u6, u7)
    r1 := Rule(head, body)

    // Set up the knowledge base.
    kb := KnowledgeBase{}
    kb.Add(r1)

    ss := SubstitutionSet{}

    X, _ := LogicVar("$X")
    goal := Complex{would_you_like, X}
    Solve(goal, kb, ss)

    results, failure := Solve(goal, kb, ss)
    if failure != "" {
        t.Error("TestJoin - " + failure)
        return
    }

    expected := "coffee, tea or juice?"
    actual := results.GetTerm(1).String()

    if actual != expected {
        t.Error("\nTestJoin - Expected: " + expected +
                "\n                Was: " + actual)
    }

} // TestJoin
