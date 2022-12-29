package main

// Test the 'append' predicate. The append predicate is used to join
// terms into a list. For example:
//
// $X = raspberry, append(cherry, [strawberry, blueberry], $X, $Out).
//
// The last term of append() is an output term. For the above, $Out
// should unify with: [cherry, strawberry, blueberry, raspberry]
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestAppend(t *testing.T) {

    fmt.Println("TestAppend")

    X, _   := LogicVar("$X")
    Y, _   := LogicVar("$Y")
    Out, _ := LogicVar("$Out")

    test_append := Atom("test_append")
    orange  := Atom("orange")
    red     := Atom("red")
    green   := Atom("green")
    blue    := Atom("blue")
    purple  := Atom("purple")
    colours := MakeLinkedList(false, green, blue, purple)

    // test_append($Out) :- $X = red, $Y = colours, append($X, orange, $Y, $Out).

    ss := SubstitutionSet{}

    head := Complex{test_append, Out}
    u1   := Unify(X, red)
    u2   := Unify(Y, colours)
    ap   := Append(red, orange, colours, Out)
    body := And(u1, u2, ap)
    r1   := Rule(head, body)

    kb := KnowledgeBase{}
    kb.Add(r1)
    query := MakeQuery(test_append, Out)

    results, failure := SolveAll(query, kb, ss)
    if failure != "" {
        t.Error("TestAppend - " + failure)
    }

    if len(results) < 1 {
        t.Error("TestAppend - no results.")
        return
    }

    result := results[0].String()
    expected := "test_append([red, orange, green, blue, purple])"
    if result != expected {
        t.Error("\nTestAppend - Expected: " + expected +
                "\n                  Was: " + result)
    }

} // TestAppend
