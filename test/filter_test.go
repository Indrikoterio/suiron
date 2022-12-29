package main

// Test the built-in predicates include() and exclude(). Eg.:
//
// $People = [male(Sheldon), male(Leonard), male(Raj), male(Howard),
//            female(Penny), female(Bernadette), female(Amy)]
// list_wimmin($W) :- include(female($_), $People, $W).
// list_nerds($N)  :- exclude(female($_), $People, $N).
// 
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestFilter(t *testing.T) {

    fmt.Println("TestFilter")

    W, _ := LogicVar("$W")
    N, _ := LogicVar("$N")

    c1, _ := ParseComplex("male(Sheldon)")
    c2, _ := ParseComplex("male(Leonard)")
    c3, _ := ParseComplex("male(Raj)")
    c4, _ := ParseComplex("male(Howard)")
    c5, _ := ParseComplex("female(Penny)")
    c6, _ := ParseComplex("female(Bernadette)")
    c7, _ := ParseComplex("female(Amy)")

    people := MakeLinkedList(false, c1, c2, c3, c4, c5, c6, c7)

    list_wimmin, _ := ParseComplex("list_wimmin($W)")
    list_nerds, _  := ParseComplex("list_nerds($N)")
    filter, _      := ParseComplex("female($_)")

    inc := Include(filter, people, W)
    ex  := Exclude(filter, people, N)

    kb := KnowledgeBase{}
    r1 := Rule(list_wimmin, inc)
    r2 := Rule(list_nerds, ex)
    kb.Add(r1, r2)

    query, _ := ParseQuery("list_wimmin($W)")
    result, failure := Solve(query, kb, SubstitutionSet{})
    if failure != "" {
        t.Error("TestFilter - " + failure)
    }

    actual := result[1].String()
    expected := "[female(Penny), female(Bernadette), female(Amy)]"
    if actual != expected {
        t.Error("\nTestFilter - Expected: " + expected +
                "\n                  Was: " + actual)
    }

    query, _ = ParseQuery("list_nerds($W)")
    result, failure = Solve(query, kb, SubstitutionSet{})
    if failure != "" {
        t.Error("TestFilter - " + failure)
    }

    actual = result[1].String()
    expected = "[male(Sheldon), male(Leonard), male(Raj), male(Howard)]"
    if actual != expected {
        t.Error("\nTestFilter - Expected: " + expected +
                "\n                  Was: " + actual)
    }

} // TestFilter
