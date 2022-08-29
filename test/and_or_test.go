package main

// Tests the And and Or operators of the inference engine.
// Cleve Lendon 2022

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

// Tests the 'and' and 'or' operators.
func TestAndOr(t *testing.T) {

    fmt.Println("TestAndOr")
    kb := KnowledgeBase{}

    c1, _ := ParseComplex("father(George, Frank)")
    c2, _ := ParseComplex("father(George, Sam)")
    c3, _ := ParseComplex("mother(Gina, Frank)")
    c4, _ := ParseComplex("mother(Gina, Sam)")
    c5, _ := ParseComplex("mother(Maria, Marcus)")
    c6, _ := ParseComplex("father(Frank, Marcus)")

    f1 := Fact(c1)
    f2 := Fact(c2)
    f3 := Fact(c3)
    f4 := Fact(c4)
    f5 := Fact(c5)
    f6 := Fact(c6)

    kb.Add(f1, f2, f3, f4, f5, f6)

    parent, _ := ParseComplex("parent($X, $Y)")
    father, _ := ParseComplex("father($X, $Y)")
    mother, _ := ParseComplex("mother($X, $Y)")
    or := Or(father, mother)

    r1 := Rule(parent, or)

    relative, _    := ParseComplex("relative($X, $Y)")
    grandfather, _ := ParseComplex("grandfather($X, $Y)")
    grandmother, _ := ParseComplex("grandmother($X, $Y)")
    or2 := Or(grandfather, father, grandmother, mother)

    r2 := Rule(relative, or2)
    //DBG("This is a rule: ", r2)

    father2, _ := ParseComplex("father($X, $Z)")
    parent2, _ := ParseComplex("parent($Z, $Y)")
    and := And(father2, parent2)

    r3 := Rule(grandfather, and)

    kb.Add(r1, r2, r3)
    //DBKB(kb)

    goal, _ := ParseGoal("relative($X, Marcus)")

    results, failure := SolveAll(goal, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestAndOr - " + failure)
        return
    }

    n := len(results)
    if n != 3 {
        msg := fmt.Sprintf("SolveAll - expected 3 results. Got %d.", n)
        t.Error(msg)
        return
    }

    // Check the solutions of relative($X, Marcus).
    expected := [3]string{"relative(George, Marcus)",
                          "relative(Frank, Marcus)",
                          "relative(Maria, Marcus)"}

    for i, r := range results {
        s := r.String()
        if s != expected[i] {
            t.Error("\nSolveAll - expected: " + expected[i] +
                    "\n                was: " + s)
        }
    }

}  // TestAndOr
