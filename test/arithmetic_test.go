package main

// TestArithmetic
//
// Test built-in arithmetic functions: Add, Subtract, Multiply, Divide.
//
// f(x, y) = ((x + y) - 6) * 3.4 / 3.4
//
// f(3, 7)  = 4
// f(3, -7) = -10
//
// The rule is:
//
// calculate($X, $Y, $Out) :- $A = add($X, $Y), $B = substract($A, 6),
//                            $C = multiply($B, 3.4), $Out = divide($C, 3.4).
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestArithmetic(t *testing.T) {

    fmt.Println("TestArithmetic")

    i2 := Integer(2)
    i3 := Integer(3)
    i5 := Integer(5)
    //f5 := Float(5.0)
    pi := Float(3.14159)

    kb := KnowledgeBase{}

    //------------------------------------
    // First test. Make a rule.
    // test1($X) :- $X = add(2, 3, 5).

    test1 := Atom("test1")
    X, _  := LogicVar("$X")
    head  := Complex{test1, X}
    body  := Unify(X, Add(i2, i3, i5))
    r     := Rule(head, body)
    kb.Add(r)
    
    goal := MakeGoal(test1, X)

    solution, failure := Solve(goal, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestArithmetic 1 - " + failure)
        return
    }

    expected := Integer(10)
    actual := solution.GetTerm(1).(Integer)
    if actual != expected {
        str := fmt.Sprintf("\nTestArithmetic 1 - expected: %d" +
                           "\n                        Was: %d", expected, actual)
        t.Error(str)
    }

    //------------------------------------
    // Second test.
    // test2($X) :- $X = add(2, 3.14159).

    test2 := Atom("test2")
    head = Complex{test2, X}
    body = Unify(X, Add(i2, pi))
    r    = Rule(head, body)
    kb.Add(r)
    
    goal = MakeGoal(test2, X)

    solution, failure = Solve(goal, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestArithmetic 2 - " + failure)
        return
    }

    expected2 := Float(5.14159)
    actual2 := solution.GetTerm(1).(Float)
    if actual2 != expected2 {
        str := fmt.Sprintf("\nTestArithmetic 2 - expected: %f" +
                           "\n                        Was: %f", expected2, actual2)
        t.Error(str)
    }

    //------------------------------------
    // Third test. - test parsing.
    // test3($X) :- $X = add(7.922, 3).

    r, _ = ParseRule("test3($X) :- $X = add(7.922, 3).")
    kb.Add(r)
    
    goal, _ = ParseGoal("test3($X)")

    solution, failure = Solve(goal, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestArithmetic 3 - " + failure)
        return
    }

    expected3 := Float(10.922)
    actual3 := solution.GetTerm(1).(Float)
    if actual3 != expected3 {
        str := fmt.Sprintf("\nTestArithmetic 3 - expected: %f" +
                           "\n                        Was: %f", expected3, actual3)
        t.Error(str)
    }

    //------------------------------------
    // Fourth test. - subtraction.
    // test4($X) :- $X = subtract(5, 3, 2).

    test4 := Atom("test4")
    head  = Complex{test4, X}
    body  = Unify(X, Subtract(i5, i3, i2))
    r     = Rule(head, body)
    kb.Add(r)
    
    goal = MakeGoal(test4, X)

    solution, failure = Solve(goal, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestArithmetic 4 - " + failure)
        return
    }

    expected = Integer(0)
    actual = solution.GetTerm(1).(Integer)
    if actual != expected {
        str := fmt.Sprintf("\nTestArithmetic 4 - expected: %d" +
                           "\n                        Was: %d", expected, actual)
        t.Error(str)
    }

    //------------------------------------
    // Fifth test. - subtraction.
    // test5($X) :- $X = subtract(7.5, 2).

    r, _ = ParseRule("test5($X) :- $X = subtract(5.68, 3).")
    kb.Add(r)
    goal, _ = ParseGoal("test5($X)")

    solution, failure = Solve(goal, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestArithmetic 5 - " + failure)
        return
    }

    expected5 := Float(2.68)
    actual5 := solution.GetTerm(1).(Float)
    diff := expected5 - actual5

    if diff > 0.0000000000000005 {
       str := fmt.Sprintf("\nTestArithmetic 5 - expected: %f" +
                          "\n                        Was: %f", expected5, actual5)
       t.Error(str)
    }

} // TestArithmetic
