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
// calculate($X, $Y, $Out) :- $A = add($X, $Y), $B = subtract($A, 6),
//                            $C = multiply($B, 3.4), $Out = divide($C, 3.4).
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "math"
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
    
    query := MakeQuery(test1, X)

    solution, failure := Solve(query, kb, SubstitutionSet{})
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
    
    query = MakeQuery(test2, X)

    solution, failure = Solve(query, kb, SubstitutionSet{})
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
    
    query, _ = ParseQuery("test3($X)")

    solution, failure = Solve(query, kb, SubstitutionSet{})
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
    
    query = MakeQuery(test4, X)

    solution, failure = Solve(query, kb, SubstitutionSet{})
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
    query, _ = ParseQuery("test5($X)")

    solution, failure = Solve(query, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestArithmetic 5 - " + failure)
        return
    }

    expected5 := Float(2.68)
    actual5 := solution.GetTerm(1).(Float)
    diff := expected5 - actual5

    if math.Abs(float64(diff)) > 0.0000000000000005 {
       str := fmt.Sprintf("\nTestArithmetic 5 - expected: %f" +
                          "\n                        Was: %f", expected5, actual5)
       t.Error(str)
    }

    //------------------------------------
    // Sixth test. - multiplication.
    // test6($X) :- $X = multiply(4, 2).

    r, _ = ParseRule("test6($X) :- $X = multiply(4, 2).")
    kb.Add(r)
    query, _ = ParseQuery("test6($X)")

    solution, failure = Solve(query, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestArithmetic 6 - " + failure)
        return
    }

    expected6 := Integer(8)
    actual6 := solution.GetTerm(1).(Integer)

    if actual6 != expected6 {
       str := fmt.Sprintf("\nTestArithmetic 6 - expected: %d" +
                          "\n                        Was: %d", expected6, actual6)
       t.Error(str)
    }

    //------------------------------------
    // Seventh test. - multiplication.
    // test7($X) :- $X = multiply(3.14159, 2).

    r, _ = ParseRule("test7($X) :- $X = multiply(3.14159, 2).")
    kb.Add(r)
    query, _ = ParseQuery("test7($X)")

    solution, failure = Solve(query, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestArithmetic 7 - " + failure)
        return
    }

    expected7 := Float(3.14159 * 2)
    actual7 := solution.GetTerm(1).(Float)

    diff = expected7 - actual7

    if math.Abs(float64(diff)) > 0.0000000000000005 {
       str := fmt.Sprintf("\nTestArithmetic 7 - expected: %f" +
                          "\n                        Was: %f", expected7, actual7)
       t.Error(str)
    }

    //------------------------------------
    // Eighth test. - divide.
    // test8($X) :- $X = divide(4, 2).

    r, _ = ParseRule("test8($X) :- $X = divide(4, 2).")
    kb.Add(r)
    query, _ = ParseQuery("test8($X)")

    solution, failure = Solve(query, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestArithmetic 8 - " + failure)
        return
    }

    expected8 := Float(2)
    actual8 := solution.GetTerm(1).(Float)

    diff = expected8 - actual8

    if math.Abs(float64(diff)) > 0.0000000000000005 {
       str := fmt.Sprintf("\nTestArithmetic 8 - expected: %f" +
                          "\n                        Was: %f", expected8, actual8)
       t.Error(str)
    }

    //------------------------------------
    // Ninth test. - a formula.
    // test9($X) :- $X = divide(4, 2).
    //
    // f(x, y) = ((x + y) - 6) * 3.4 / 3.4
    //
    // f(3, 7)  = 4
    // f(3, -7) = -10
    //
    // The rule is:
    //
    // calculate($X, $Y, $Out) :- $A = add($X, $Y),
    //                            $B = subtract($A, 6),
    //                            $C = multiply($B, 3.4),
    //                            $Out = divide($C, 3.4).

    Y, _   := LogicVar("$Y")
    A, _   := LogicVar("$A")
    B, _   := LogicVar("$B")
    C, _   := LogicVar("$C")
    Out, _ := LogicVar("$Out")

    head, _ = ParseComplex("calculate($X, $Y, $Out)")
    u1 := Unify(A, Add(X, Y))
    u2 := Unify(B, Subtract(A, Integer(6)))
    u3 := Unify(C, Multiply(B, Float(3.4)))
    u4 := Unify(Out, Divide(C, Float(3.4)))

    r = Rule(head, And(u1, u2, u3, u4))
    kb.Add(r)

    calc, _ := ParseQuery("calculate(3.0, 7.0, $Out)")

    solution, failure = Solve(calc, kb, SubstitutionSet{})
    if len(failure) != 0 {
        t.Error("TestArithmetic 9 - " + failure)
        return
    }

    expected9 := Float(4)
    actual9 := solution.GetTerm(3).(Float)

    diff = expected9 - actual9

    if math.Abs(float64(diff)) > 0.0000000000000005 {
       str := fmt.Sprintf("\nTestArithmetic 9 - expected: %f" +
                          "\n                        Was: %f", expected9, actual9)
       t.Error(str)
    }

} // TestArithmetic
