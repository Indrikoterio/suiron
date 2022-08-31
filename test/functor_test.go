package main

// Test the functor() predicate.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestFunctor(t *testing.T) {

    fmt.Println("TestFunctor")

    kb := KnowledgeBase{}

    // Create logic variables.
    X, _  := LogicVar("$X")
    Y, _  := LogicVar("$Y")

    get       := Atom("get")
    animal, _ := ParseComplex("mouse(mammal, rodent)")

    // Make 'get' rule.
    // get($Y) :- functor(mouse(mammal, rodent), $X), $X = $Y.
    head := Complex{get, Y}
    body := And(Functor(animal, X), Unify(X, Y))
    r1 := Rule(head, body)
    kb.Add(r1)

    // Does ParseRule correctly part the functor predicate?
    r2, _ := ParseRule("get($Y) :- $X = cat(mammal, carnivore), functor($X, $Y).")
    kb.Add(r2)

    goal := MakeGoal(get, X)  // get($X)

    // Check the solutions for get($X).
    expected := [2]string{"mouse", "cat"}

    // Get the root solution node.
    root := goal.GetSolver(kb, SubstitutionSet{}, nil)

    for i := 0; i < 2; i++ {
        solution, found := root.NextSolution()
        if !found {
            fmt.Println("TestFunctor - expected two solutions")
            return
        }
        result := goal.ReplaceVariables(solution).(Complex)
        str := result.GetTerm(1).String()
        if str != expected[i] {
            t.Error("\nTestFunctor - Expected: " + expected[i] +
                    "\n                   Was: " + str)
            return
        }
    } // for

    // Check to make sure we can get the arity also.
    // check_arity($X, $Y) := functor(diamonds(forever, a girl's...), $X, $Y).

    mineral, _ := ParseComplex("diamonds(forever, a girl's best friend)")
    check_arity := Atom("check_arity")
    head  = Complex{check_arity, X, Y}
    body2 := Functor(mineral, X, Y)
    r3 := Rule(head, body2)
    kb.Add(r3)

    goal = MakeGoal(check_arity, X, Y)

    // Get the root solution node.
    root = goal.GetSolver(kb, SubstitutionSet{}, nil)

    solution, found := root.NextSolution()
    if !found {
        fmt.Println("TestFunctor - no solution.")
        return
    }

    result := goal.ReplaceVariables(solution).(Complex)
    functor := result.GetTerm(1).String()
    arity   := result.GetTerm(2).String()

    if functor != "diamonds" {
        t.Error("\nTestFunctor - Expected: diamonds" +
                "\n                   Was: " + functor)
        return
    }

    if arity != "2" {
        t.Error("\nTestFunctor - Expected: 2" +
                "\n                   Was: " + arity)
        return
    }

}  // TestFunctor
