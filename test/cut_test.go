package main

// Test the cut operator.
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestCut(t *testing.T) {

    fmt.Println("TestCut")
    kb := KnowledgeBase{}
    ss := SubstitutionSet{}

    /* Test this rule:
        cut_rule :- !, a = b.  % This fails.
        cut_rule :- print(*** This should NOT print. ***).
        cut_rule(OK).
        test($X) :- cut_rule, $X = Bad.
        test($X) :- cut_rule($X).
       Note: cut_rule/0 and cut_rule/1 are two different rules.
    */

    // Set up facts and rules.

    c1 := Complex{Atom("cut_rule")}
    a1 := And(Cut(), Unify(Atom("a"), Atom("b")))
    a2 := Print(Atom("*** This should NOT print. ***"))
    r1 := Rule(c1, a1)  // cut_rule :- !, a = b.
    r2 := Rule(c1, a2)  // cut_rule :- print(*** This should NOT print. ***).
    c2 := Complex{Atom("cut_rule"), Atom("OK")}
    f1 := Fact(c2)  // cut_rule(OK).

    X, _ := LogicVar("$X")
    c3   := Complex{Atom("test"), X}
    c4   := Complex{Atom("cut_rule")}
    a3   := And(c4, Unify(X, Atom("Bad")))
    r3   := Rule(c3, a3)   // test($X) :- cut_rule, $X = Bad.
    c5   := Complex{Atom("cut_rule"), X}
    r4   := Rule(c3, c5)   // test($X) :- cut_rule($X).
    kb.Add(r1, r2, f1, r3, r4)

    //DBKB(kb)
    goal := MakeGoal(Atom("test"), X)

    solutions, failure := SolveAll(goal, kb, ss)

    if failure != "" {
        t.Error("TestCut - " + failure)
        return
    }

    length := len(solutions)
    if length != 1 {
        t.Error("TestCut - expected 1 result. There were ", length, ".")
        return
    }

    result := solutions[0].GetTerm(1).String()
    expected := "OK"
    if result != expected {
        t.Error("\nTestCut - Expected: " + expected +
                "\n               Was: " + result)
    }

    /*
     * handicapped(John).
     * handicapped(Mary).
     * has_small_children(Mary).
     * is_elderly(Diane)
     * is_elderly(John)
     * priority_seating($Name, $YN) :- handicapped($Name), $YN = Yes, !.
     * priority_seating($Name, $YN) :- has_small_children($Name), $YN = Yes, !.
     * priority_seating($Name, $YN) :- is_elderly($Name), $YN = Yes, !.
     * priority_seating($Name, No).
     */

    handicapped        := Atom("handicapped")
    has_small_children := Atom("has_small_children")
    is_elderly         := Atom("is_elderly")
    priority_seating   := Atom("priority_seating")
    yes                := Atom("Yes")
    no                 := Atom("No")

    John  := Atom("John")
    Mary  := Atom("Mary")
    Diane := Atom("Diane")

    Name, _ := LogicVar("$Name")
    YN, _   := LogicVar("$YN")

    d1 := Complex{handicapped, John}
    d2 := Complex{handicapped, Mary}
    d3 := Complex{has_small_children, Mary}
    d4 := Complex{is_elderly, Diane}
    d5 := Complex{is_elderly, John}

    fact1 := Fact(d1)
    fact2 := Fact(d2)
    fact3 := Fact(d3)
    fact4 := Fact(d4)
    fact5 := Fact(d5)

    kb.Add(fact1, fact2, fact3, fact4, fact5)

    h1 := Complex{priority_seating, Name, YN}
    b1 := And(Complex{handicapped, Name}, Unify(YN, yes), Cut())
    rule1 := Rule(h1, b1)

    h2 := Complex{priority_seating, Name, YN}
    b2 := And(Complex{has_small_children, Name}, Unify(YN, yes), Cut())
    rule2 := Rule(h2, b2)

    h3 := Complex{priority_seating, Name, YN}
    b3 := And(Complex{is_elderly, Name}, Unify(YN, yes), Cut())
    rule3 := Rule(h3, b3)

    h4 := Complex{priority_seating, Name, no}
    fact6 := Fact(h4)

    kb.Add(rule1, rule2, rule3, fact6)
    goal = MakeGoal(priority_seating, John, X)

    solutions, failure = SolveAll(goal, kb, ss)

    if failure != "" {
        t.Error("TestCut - " + failure)
        return
    }

    result2 := solutions[0].GetTerm(2).String()

    expected2 := "Yes"
    if result2 != expected2 {
        t.Error("\nTestCut - Expected: " + expected2 +
                "\n               Was: " + result2)
    }

}  // TestBackChaining
