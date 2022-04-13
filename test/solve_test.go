package main

// Tests methods which search for solutions.
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

// Test the solve() method, which finds one solution.
func TestSolve(t *testing.T) {

    fmt.Println("TestSolve")
    kb := KnowledgeBase{}
    ss := SubstitutionSet{}

    hobby     := Atom("hobby")
    chess     := Atom("chess")
    dance     := Atom("dance")
    gardening := Atom("gardening")

    tim    := Atom("Tim")
    robert := Atom("Robert")
    sarah  := Atom("Sarah")

    c1 := Complex{hobby, tim, dance}
    c2 := Complex{hobby, robert, gardening}
    c3 := Complex{hobby, sarah, chess}

    f1 := Fact(c1)
    f2 := Fact(c2)
    f3 := Fact(c3)

    kb.Add(f1, f2, f3)

    x, _ := LogicVar("$X")

    // Do not use Complex{} to create a goal, because variables
    // must have unique IDs. MakeGoal calls recreateVariables.
    // Bad> goal := Complex{hobby, tim, x}
    goal := MakeGoal(hobby, tim, x)  // Goal is: hobby(Tim, $X)

    expected := "hobby(Tim, dance)"
    actual, failure := Solve(goal, kb, ss)
    if len(failure) != 0 {
        t.Error("TestSolveAll - " + failure)
        return
    }

    s := actual.String()
    if s != expected {
        t.Error("\nTestSolve - Expected: " + expected +
                "\n                 Was: " + s)
    }

    y, _ := LogicVar("$Y")

    // Do not use Complex{} to create a goal, because variables
    // must have unique IDs. MakeGoal calls recreateVariables.
    // Bad> goal := goal = Complex{hobby, x, y}
    goal = MakeGoal(hobby, x, y)

    results, failure := SolveAll(goal, kb, ss)
    if len(failure) != 0 {
        t.Error("TestSolveAll - " + failure)
        return
    }

    exp1 := "hobby(Tim, dance)"
    exp2 := "hobby(Robert, gardening)"
    exp3 := "hobby(Sarah, chess)"

    if len(results) != 3 {
        t.Error("TestSolveAll - There should be 3 results.")
        return
    }

    s = results[0].String()
    if exp1 != s {
        t.Error("\nTestSolveAll - Expected: " + exp1 +
                "\n                    Was: " + s)
    }

    s = results[1].String()
    if exp2 != s {
        t.Error("\nTestSolveAll - Expected: " + exp2 +
                "\n                    Was: " + s)
    }

    s = results[2].String()
    if exp3 != s {
        t.Error("\nTestSolveAll - Expected: " + exp3 +
                "\n                    Was: " + s)
    }

    fmt.Println("TestTimeOut 1")

    c4 := Complex{Atom("Time out test.")}
    // The predicate too_long has a sleep timer
    // which should cause a timeout error.
    r1 := Rule(c4, TooLong())
    kb.Add(r1)

    // Here it's OK to use c4 as a goal, without calling
    // MakeGoal(), because c4 does not contain variables.
    goal = c4
    _, failure = Solve(goal, kb, ss)
    if len(failure) == 0 {
        t.Error("TestTimeOut - this test should time out.")
        return
    }

    fmt.Println("TestTimeOut 2")
    // Second timeout test. Escape from endless loop.
    // endless($X) := endless($X)

    endless := Atom("endless")
    cEndless := Complex{endless, x}  // Term is:  endless($X)
    r2 := Rule(cEndless, cEndless) // Rule is: endless($X) :- endless($X).
    kb.Add(r2)

    // Here it's OK to use Complex{} as a goal, without calling
    // MakeGoal(), because the goal does not contain variables.
    goal = Complex{endless, Atom("loop")} // Goal is: endless(loop)
    _, failure = Solve(goal, kb, SubstitutionSet{})
    //fmt.Printf("----------- %v\n", failure)
    if len(failure) == 0 {
        t.Error("TestTimeOut - this test should time out.")
        return
    }

}  // TestSolve
