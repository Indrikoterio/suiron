package main

// Test the built-in comparison predicates: > >= == <= <
// Eg.:
//    .., $X <= 23,...
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

func TestComparison(t *testing.T) {

    fmt.Println("TestComparison")

    X, _   := LogicVar("$X")
    Y, _   := LogicVar("$Y")
    Z, _   := LogicVar("$Z")

    passed   := Atom("passed")
    failed   := Atom("failed")
    Beth     := Atom("Beth")
    Albert   := Atom("Albert")
    Samantha := Atom("Samantha")
    Trevor   := Atom("Trevor")
    Joseph   := Atom("Joseph")

    test := Atom("test")
    test_greater_than          := Atom("test_greater_than")
    test_greater_than_or_equal := Atom("test_greater_than_or_equal")
    test_less_than             := Atom("test_less_than")
    test_less_than_or_equal    := Atom("test_less_than_or_equal")
    test_equal                 := Atom("test_equal")

    kb := KnowledgeBase{}
    ss := SubstitutionSet{}

    var head Complex

    head = Complex{test_greater_than, X, Y, Z}    
    body := And(GreaterThan(X, Y), Cut(), Unify(Z, passed))
    r1 := Rule(head, body)

    head = Complex{test_greater_than, Anon(), Anon(), Z}
    body2 := Unify(Z, failed)
    r2 := Rule(head, body2)

    head = Complex{test_less_than, X, Y, Z}    
    body3 := And(LessThan(X, Y), Cut(), Unify(Z, passed))
    r3 := Rule(head, body3)

    head = Complex{test_less_than, Anon(), Anon(), Z}
    body4 := Unify(Z, failed)
    r4 := Rule(head, body4)

    head = Complex{test_greater_than_or_equal, X, Y, Z}    
    body5 := And(GreaterThanOrEqual(X, Y), Cut(), Unify(Z, passed))
    r5 := Rule(head, body5)

    head = Complex{test_greater_than_or_equal, Anon(), Anon(), Z}
    body6 := Unify(Z, failed)
    r6 := Rule(head, body6)

    head = Complex{test_less_than_or_equal, X, Y, Z}    
    body7 := And(LessThanOrEqual(X, Y), Cut(), Unify(Z, passed))
    r7 := Rule(head, body7)

    head = Complex{test_less_than_or_equal, Anon(), Anon(), Z}
    body8 := Unify(Z, failed)
    r8 := Rule(head, body8)

    head = Complex{test_equal, X, Y, Z}    
    body9 := And(Equal(X, Y), Cut(), Unify(Z, passed))
    r9 := Rule(head, body9)

    head = Complex{test_equal, Anon(), Anon(), Z}
    body10 := Unify(Z, failed)
    r10 := Rule(head, body10)

    head = Complex{test, Z}

    body11 := Complex{test_greater_than, Integer(4), Integer(3), Z}
    r11 := Rule(head, body11)
    body12 := Complex{test_greater_than, Beth, Albert, Z}
    r12 := Rule(head, body12)
    body13 := Complex{test_greater_than, Integer(2), Integer(3), Z}
    r13 := Rule(head, body13)

    body14 := Complex{test_less_than, Float(1.6), Float(7.2), Z}
    r14 := Rule(head, body14)
    body15 := Complex{test_less_than, Samantha, Trevor, Z}
    r15 := Rule(head, body15)
    body16 := Complex{test_less_than, Float(4.222), Float(4.), Z}
    r16 := Rule(head, body16)

    body17 := Complex{test_greater_than_or_equal, Integer(4), Float(4.0), Z}
    r17 := Rule(head, body17)
    body18 := Complex{test_greater_than_or_equal, Joseph, Joseph, Z}
    r18 := Rule(head, body18)
    body19 := Complex{test_greater_than_or_equal, Float(3.9), Float(4.0), Z}
    r19 := Rule(head, body19)

    body20 := Complex{test_less_than_or_equal, Float(7.000), Integer(7), Z}
    r20 := Rule(head, body20)
    body21 := Complex{test_less_than_or_equal, Float(7.000), Float(7.1), Z}
    r21 := Rule(head, body21)
    body22 := Complex{test_less_than_or_equal, Float(0.0), Integer(-20), Z}
    r22 := Rule(head, body22)

    body23 := Complex{test_equal, Joseph, Joseph, Z}
    r23 := Rule(head, body23)
    body24 := Complex{test_equal, Joseph, Trevor, Z}
    r24 := Rule(head, body24)

    kb.Add(r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11, r12)
    kb.Add(r13, r14, r15, r16, r17, r18, r19, r20, r21, r22, r23, r24)

    goal, _ := ParseGoal("test($Z)")

    solutions, failure := SolveAll(goal, kb, ss)
    if len(failure) != 0 {
        t.Error("TestComparison - " + failure)
        return
    }

    expected := [14]string{
                   "passed", "passed", "failed",
                   "passed", "passed", "failed",
                   "passed", "passed", "failed",
                   "passed", "passed", "failed",
                   "passed", "failed",
                   }

    for i, r := range solutions {
        s := r.GetTerm(1).String()
        if s != expected[i] {
            t.Error("\nTestComparison - expected: " + expected[i] +
                    "\n                      Was: " + s)
        }
    }

} // TestComparison
