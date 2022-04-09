package main

// Read in facts and rules from a text file (kings.txt)
// and executes a query. Who is Skule's grandfather?
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

// Reads in facts and rules from a text file.
func TestReadRules(t *testing.T) {

    fmt.Println("TestReadRules")
    fileName := "kings.txt"
    badFileName  := "kings!.txt"

    kb := KnowledgeBase{}

    err := LoadKBFromFile(kb, badFileName)
    expected := "open kings!.txt: no such file or directory"
    if err.Error() != expected {
        t.Error("\nTestReadRules - Should produce error:\n" + expected)
        return
    }

    err = LoadKBFromFile(kb, "badrule1.txt")
    expected = "Error - unmatched bracket: (\n" +
               "Check start of file."
    if err.Error() != expected {
        t.Error("\nTestReadRules - Should produce error:\n" + expected)
        return
    }

    err = LoadKBFromFile(kb, "badrule2.txt")
    expected = "Error - unmatched bracket: )\n" +
               "Error occurs after: parent(Godwin, Tostig)."
    if err.Error() != expected {
        t.Error("\nTestReadRules - Should produce error:\n" + expected)
        return
    }

    err = LoadKBFromFile(kb, "badrule3.txt")
    expected = "Check line 3: par"
    if err.Error() != expected {
        t.Error("\nTestReadRules - Should produce error:\n" + expected)
        return
    }

    err = LoadKBFromFile(kb, fileName)
    if err != nil { 
        t.Error("\nTestReadRules:\n", err.Error())
        return
    }

    goal, _ := ParseComplex("grandfather($X, Skule)")
    solution, failure := Solve(goal, kb, SubstitutionSet{})

    if len(failure) > 0 {
        t.Error("\nTestReadRules: No solution.\n", failure)
        return
    }

    grandfather := solution.GetTerm(1).String()
    expected = "Godwin"
    if grandfather != expected {
        t.Error("\nTestReadRules - Solution should be Godwin." +
                "\n                Was: ", grandfather)
        return
    }

}  // ReadRules
