package main

// Testing the Not operator.
//
// parent(Sarah, Daniel).
// parent(Richard, Daniel).
// female(Sarah).
// mother($X, $Y) :- female($X), parent($X, $Y).
// father($X, $Y) :- parent($X, $Y), not(female($X)).
//
// ?- father($X, Daniel)
//
// -----------------------------------
// A second and third test.
//
// friend(Sheldon).
// friend(Leonard).
// friend(Penny).
// invite($X) :- friend($X), not($X = Sheldon).
// invite2($X) :- friend($X), not($X = Leonard).
//
// ?- invite($X)
// ?- invite2($X)
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestNot(t *testing.T) {

    fmt.Println("TestNot")

    // For first test.
    // ?- father($X, Daniel).

    f1, _ := ParseRule("parent(Sarah, Daniel)")
    f2, _ := ParseRule("parent(Richard, Daniel)")
    f3, _ := ParseRule("female(Sarah)")

    r1, _ := ParseRule("mother($X, $Y) :- female($X), parent($X, $Y).")
    r2, _ := ParseRule("father($X, $Y) :- parent($X, $Y), not(female($X)).")

    // For second test.
    f4, _  := ParseRule("friend(Sheldon)")
    f5, _  := ParseRule("friend(Leonard)")
    f6, _  := ParseRule("friend(Penny)")

    r3, _ := ParseRule("invite($X)  :- friend($X), not($X = Sheldon).")
    r4, _ := ParseRule("invite2($X) :- friend($X), not($X = Leonard).")

    kb := KnowledgeBase{}
    kb.Add(f1, f2, f3, f4, f5, f6, r1, r2, r3, r4)

    // ?- father($X, Daniel)
    goal, _ := ParseGoal("father($X, Daniel)")
    result, failure := Solve(goal, kb, SubstitutionSet{})
    if failure != "" { t.Error("TestNot - " + failure); return }

    expected := "Richard"
    actual := result[1].String()

    if actual != expected {
        t.Error("\nTestNot - Expected: " + expected +
                "\n               Was: " + actual)
    }

    // Second test.
    // ?- invite($X)

    goal, _ = ParseGoal("invite($X)")
    solutions, failure := SolveAll(goal, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestNot - " + failure)
        return
    }

    expected2 := [2]string{ "Leonard", "Penny" }

    for i, r := range solutions {
        s := r.GetTerm(1).String()
        if s != expected2[i] {
            t.Error("\nTestNot - expected: " + expected2[i] +
                    "\n               Was: " + s)
        }
    }

    // Third test.
    // ?- invite2($X)

    goal, _ = ParseGoal("invite2($X)")
    solutions, failure = SolveAll(goal, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestNot - " + failure)
        return
    }

    expected3 := [2]string{ "Sheldon", "Penny" }

    for i, r := range solutions {
        s := r.GetTerm(1).String()
        if s != expected3[i] {
            t.Error("\nTestNot - expected: " + expected2[i] +
                    "\n               Was: " + s)
        }
    }

} // TestNot
