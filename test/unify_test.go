package main

// Tests 'unify' predicate. Eg. $X = pronoun
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestUnify(t *testing.T) {

    fmt.Println("TestUnify 1") //-------------

    // First test is:
    // test($X) :- $X = pronoun.
    // Query is test($X)

    X, _    := LogicVar("$X")
    pronoun := Atom("pronoun")
    test    := Atom("test")
    head    := Complex{test, X}
    body    := Unify(X, pronoun)

    r1 := Rule(head, body)

    // Set up the knowledge base.
    kb := KnowledgeBase{}
    kb.Add(r1)

    query := MakeQuery(test, X)
    solution, failure := Solve(query, kb, SubstitutionSet{})

    if len(failure) != 0 {
        t.Error("TestUnify - Failure: " + failure)
        return
    }

    expected := "test(pronoun)"
    actual   := solution.String()
    if expected != actual {
        t.Error("\nTestUnify - Expected: " + expected +
                "\n                 Was: " + actual)
    }

    // Second test is:
    // test2($A, $B, $C) := [eagle, parrot, raven, sparrow] = [$A, $B | $C].
    // Query is test2($A, $B, $C)

    A, _ := LogicVar("$A")
    B, _ := LogicVar("$B")
    C, _ := LogicVar("$C")

    eagle   := Atom("eagle")
    parrot  := Atom("parrot")
    raven   := Atom("raven")
    sparrow := Atom("sparrow")
    birds   := MakeLinkedList(false, eagle, parrot, raven, sparrow)
    list    := MakeLinkedList(true, A, B, C)

    test2   := Atom("test2")
    head2   := Complex{test2, A, B, C}
    body2   := Unify(birds, list)

    r2 := Rule(head2, body2)
    kb.Add(r2)

    query = MakeQuery(test2, A, B, C)
    solution, failure = Solve(query, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("\nTestUnify - Solve failed.")
        return
    }

    expected = "test2(eagle, parrot, [raven, sparrow])"
    actual   = solution.String()

    if expected != actual {
        t.Error("\nTestUnify - Expected: " + expected +
                "\n                 Was: " + actual)
    }

    // Test the parsing functionality for 
    fmt.Println("TestUnify 2")  //-------------

    /*
        unify_test($X, $Y, $Z) :- lawyer = lawyer,
                                  job(programmer, $Z) = job($Y, janitor),
                                  $W = $X, job($W).
    */

    lawyer := Atom("lawyer")

    c1, _ := ParseComplex("job(lawyer)")
    c2, _ := ParseComplex("job(teacher)")
    c3, _ := ParseComplex("job(programmer)")
    c4, _ := ParseComplex("job(janitor)")

    f1 := Fact(c1)
    f2 := Fact(c2)
    f3 := Fact(c3)
    f4 := Fact(c4)
    kb.Add(f1, f2, f3, f4)

    u1 := Unify(lawyer, lawyer)
    u2, _ := ParseUnify("job(programmer, $Z) = job($Y, janitor)")
    u3, _ := ParseUnify("$W = $X")

    head, _ = ParseComplex("unify_test($X, $Y, $Z)")
    c, _ := ParseComplex("job($W)")
    body3 := And(u1, u2, u3, c)
    r1 = Rule(head, body3)
    kb.Add(r1)

    query, _ = ParseQuery("unify_test($X, $Y, $Z)")
    solutions, failure := SolveAll(query, kb, SubstitutionSet{})

    // Expected solutions of unify_test($X, $Y, $Z).
    expected2 := [4]string{"unify_test(lawyer, programmer, janitor)",
                           "unify_test(teacher, programmer, janitor)",
                           "unify_test(programmer, programmer, janitor)",
                           "unify_test(janitor, programmer, janitor)"}

    if len(solutions) != 4 {
        t.Error("TestUnify - Parse - Expecting 4 solutions.")
        return
    }

    for i := 0; i < 4; i++ {
        solution := solutions[i]
        actual := solution.String()
        exp := expected2[i]
        if actual != exp {
            t.Error("\nTestUnify - Expected: " + exp +
                    "\n                 Was: " + actual)
            return
        }
    }

    // Test the parsing functionality for 
    fmt.Println("TestUnify 3")   //-------------

    /*
      second_test($Y) :- $X = up, $Y = down, $X = $Y.
      This query must fail.
     */

    u1, _ = ParseUnify("$X = up")
    u2, _ = ParseUnify("$Y = down")
    u3, _ = ParseUnify("$X = $Y")
    head, _ = ParseComplex("second_test($Y)")
    body4 := And(u1, u2, u3)
    r2 = Rule(head, body4)
    kb.Add(r2)

    query, _ = ParseQuery("second_test($Y)")
    _, failure = SolveAll(query, kb, SubstitutionSet{})

    if failure != "No" {
        t.Error("TestUnify - Query must fail.")
    }

} // TestUnify
