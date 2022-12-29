package main

// Tests the 'time' predicate, which measures execution time
// of a goal. The test loads a qsort algorithm from the file
// qsort.txt, and runs it.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestTime(t *testing.T) {

    fmt.Println("TestTime:")

    kb := KnowledgeBase{}
    err := LoadKBFromFile(kb, "qsort.txt")
    if err != nil { 
        t.Error("\nTestTime:\n", err.Error())
        return
    }

    query, _ := ParseQuery("measure")

    _, failure := Solve(query, kb, SubstitutionSet{})
    if len(failure) > 0 {
        t.Errorf("\nTestTime: %v\n", failure)
        return
    }

} // TestTime
