package main

// Tests 'print_list' predicate.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestPrintList(t *testing.T) {

    //----------------------------------------------------
    fmt.Println("TestPrintList: - should print:\n1, 2, 3, a, b, c")

    kb := KnowledgeBase{}

    head, _  := ParseComplex("print_list_test")
    u1, _ := ParseUnify("$X = [a, b, c]")
    u2, _ := ParseUnify("$List = [1, 2, 3 | $X]")
    p, _  := ParseSubgoal("print_list($List)")
    body := And(u1, u2, p)

    r1 := Rule(head, body)
    kb.Add(r1)

    query, _ := ParseQuery("print_list_test")
    Solve(query, kb, SubstitutionSet{})

} // TestPrintList
