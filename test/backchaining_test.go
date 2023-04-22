package main

// Test the backchaining functionality of the inference engine.
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

// Set up facts and rules.
func TestBackChaining(t *testing.T) {

    fmt.Println("TestBackChaining")
    kb := KnowledgeBase{}
    ss := SubstitutionSet{}

    // Create logic variables.
    X, _ := LogicVar("$X")
    Y, _ := LogicVar("$Y")
    Z, _ := LogicVar("$Z")

    parent   := Atom("parent")
    ancestor := Atom("ancestor")
    charles  := Atom("Charles")
    tony     := Atom("Tony")
    maria    := Atom("Maria")
    bill     := Atom("Bill")
    audrey   := Atom("Audrey")

    // Create a few facts.
    c1 := Complex{parent, bill, audrey} // parent(Bill, Audrey)
    c2 := Complex{parent, maria, bill}  // parent(Maria, Bill)
    c3 := Complex{parent, tony, maria}
    c4 := Complex{parent, charles, tony}

    f1 := Fact(c1)
    f2 := Fact(c2)
    f3 := Fact(c3)
    f4 := Fact(c4)

    // Register the above facts in the knowledgebase.
    kb.Add(f1, f2, f3, f4)

    head := Complex{ancestor, X, Y}
    c5   := Complex{parent, X, Y}
    c6   := Complex{parent, X, Z}
    c7   := Complex{ancestor, Z, Y}

    // ancestor($X, $Y) := parent($X, $Y).
    r1 := Rule(head, c5)

    // ancestor($X, $Y) := parent($X, $Z), ancestor($Z, $Y).
    body := And(c6, c7)
    r2 := Rule(head, body)

    // Register the above rules in the knowledgebase.
    kb.Add(r1, r2)

    query := MakeQuery(ancestor, charles, Y)

    // Check the solutions of ancestor(Charles, $Y).
    expected := [4]string{"ancestor(Charles, Tony)",
                          "ancestor(Charles, Maria)",
                          "ancestor(Charles, Bill)",
                          "ancestor(Charles, Audrey)"}


    solutions, failure := SolveAll(query, kb, ss)
    if len(failure) != 0 {
        t.Error("TestBackChaining - " + failure)
        return
    }

    for i, r := range solutions {
        s := r.String()
        if s != expected[i] {
            t.Error("\nBackChaining - expected: " + expected[i] +
                    "\n                    Was: " + s)
        }
    }

}  // TestBackChaining
