package main

// Tests for knowledge base.
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

// Create a knowledge base, add facts, and display them.
func TestKnowledgeBase(t *testing.T) {

    fmt.Println("TestKnowledgeBase")
    kb := KnowledgeBase{}

    age := Atom("age")
    loves := Atom("loves")
    //vehicle := Atom("vehicle")

    n1 := Atom("Ross")
    n2 := Atom("Rachel")
    n3 := Atom("Chandler")
    n4 := Atom("Monica")
    a1 := Integer(27)
    a2 := Integer(28)
    // ages
    c1 := Complex{age, n1, a1}
    c2 := Complex{age, n2, a2}
    c3 := Complex{age, n3, a2}
    c4 := Complex{age, n4, a1}
    // loves
    c5 := Complex{loves, n1, n2}
    c6 := Complex{loves, n2, n1}
    c7 := Complex{loves, n3, n4}
    c8 := Complex{loves, n4, n3}

    x, _ := LogicVar("$X")
    y, _ := LogicVar("$Y")
    c9 := Complex{loves, x, y}

    // Facts are rules.
    kb.Add(Fact(c1))
    kb.Add(Fact(c2))
    kb.Add(Fact(c3))
    kb.Add(Fact(c4))
    kb.Add(Fact(c5))
    kb.Add(Fact(c6))
    kb.Add(Fact(c7))
    kb.Add(Fact(c8))

    expected := `
########## Contents of Knowledge Base ##########
age/2
    age(Ross, 27).
    age(Rachel, 28).
    age(Chandler, 28).
    age(Monica, 27).
loves/2
    loves(Ross, Rachel).
    loves(Rachel, Ross).
    loves(Chandler, Monica).
    loves(Monica, Chandler).
`   //------------------------End of expected

    actual := kb.FormatKB()
    if actual != expected {
        t.Error("KnowledgeBase is different than expected.\n" + actual)
    }

    expected = "loves(Rachel, Ross)."
    rule1 := kb.GetRule(c9, 1)
    actual = rule1.String()

    if actual != expected {
        t.Error("\nGetRule, expected: " + expected +
                "\nWas:               " + actual)
    }
}  // TestKnowledgeBase
