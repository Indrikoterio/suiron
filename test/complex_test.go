package main

// Test creation and unification of complex terms.
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "fmt"
    "testing"
)

// TestComplex = tests creation and unification of complex terms.
func TestComplex(t *testing.T) {

    fmt.Println("TestComplex")
    owns  := Atom("owns")
    john  := Atom("John")
    house := Atom("house")
    car := Atom("car")
    c1 := Complex{owns, john, house}  // owns(John, house)
    c2 := Complex{owns, john, house}  // owns(John, house)
    c3 := Complex{owns, john, car}    // owns(John, car)

    v, ok  := LogicVar("$X")
    if ok != nil { t.Error("LogicVar() - Invalid variable: $X") }
    c4 := Complex{owns, john, v}

    ss := SubstitutionSet{}

    newSS, success := c1.Unify(c2, ss)
    if !success { t.Error("c1 should unify with c2") }

    newSS, success = c1.Unify(c3, newSS)
    if success { t.Error("c1 should not unify with c3") }
    if len(newSS) > 0 { t.Error("Should not change substitution set.") }

    newSS, success = c1.Unify(c4, newSS)
    if !success { t.Error("c1 should unify with c4") }
    if len(newSS) != 1 { t.Error("Substitution set should contain one substitution.") }

    str := c3.String()
    if str != "owns(John, car)" {
        t.Error("String should be: owns(John, car). " + str)
    }

    key := c3.Key()
    if key != "owns/2" { t.Error("Key must be owns/2. " + key) }

} // TestComplex


// TestParseComplex - tests to confirm that the function ParseComplex()
// parses a complex term, and particularly its argument string, correctly.
// The creation of Atoms, Integers, Floats and Variables is confirmed.
// Backslashes are used to escape commas.
func TestParseComplex(t *testing.T) {

    testName := "TestParseComplex 1"
    fmt.Println(testName)

    // Make complex terms from strings. Test parsing of arguments.
    // Must use double backslash to escape a comma (term 5), because
    // backslash is also an escape character to the go compiler.
    c, _ := ParseComplex("dingo(Arthur, 414, 7.59, \"7.59\", This term\\, has a comma.)")

    fun := c.GetFunctor()
    if fun != "dingo" { t.Error("Invalid functor.") }
    if fun.TermType() != ATOM { t.Error("Invalid functor type.") }

    term1 := c.GetTerm(1)
    if term1.(Atom) != "Arthur" { t.Error("Invalid term. " + term1.(Atom)) }
    if term1.TermType() != ATOM { t.Error("Invalid term type. " + term1.(Atom)) }

    term2 := c.GetTerm(2)
    if term2.(Integer) != 414 { t.Error("Invalid term. 414") }
    if term2.TermType() != INTEGER { t.Error("Invalid term type. 414") }

    term3 := c.GetTerm(3)
    if term3.(Float) != 7.59 { t.Error("Invalid term. 7.59") }
    if term3.TermType() != FLOAT { t.Error("Invalid term type. 7.59") }

    // Any term enclosed by quotes is an Atom.
    term4 := c.GetTerm(4)
    if term4.(Atom) != "7.59" { t.Error("Invalid term. \"7.59\"") }
    if term4.TermType() != ATOM { t.Error("Invalid term type. \"7.59\"") }

    // Use a backslash to escape characters. In this case, a comma: \,
    term5 := c.GetTerm(5)
    if term5.(Atom) != "This term, has a comma." { t.Error("Invalid term. " + term5.(Atom)) }
    if term5.TermType() != ATOM { t.Error("Invalid term type. " + term5.(Atom)) }

    c2, _ := ParseComplex("double_quote(\\\")")
    quoteTerm := c2.GetTerm(1)
    if quoteTerm.(Atom) != "\"" { t.Error("Invalid quote-escape. " + quoteTerm.(Atom)) }

    testName = "TestParseComplex 2"
    fmt.Println(testName) //----------------------------------------

    c3, _ := ParseComplex("test($X, 3.14159, [\"a,b,c\", 3.14159, e | $Y])")
    functor := c3.GetFunctor()

    // Check functor.
    expected := "test"
    actual   := functor.String()
    if expected != actual {
        t.Error("\n" + testName + " - Expected: " + expected +
                "\n                          Was: " + actual)
    }

    // Check first term.
    term1 = c3.GetTerm(1)
    tt := term1.TermType()
    expected = "$X"
    actual = term1.String()
    if tt != VARIABLE {
        t.Error("\n" + testName + " " + actual + " - invalid type.")
        return
    }
    if expected != actual {
        t.Error("\n" + testName + " - Expected: " + expected +
                "\n                          Was: " + actual)
    }

    // Check second term.
    term2 = c3.GetTerm(2)
    tt = term2.TermType()
    if tt != FLOAT {
        t.Error("\n" + testName + " " + term2.String() + " - invalid type.")
    }
    pi := 3.14159
    strPi := fmt.Sprintf("%v", pi)
    fl := float64(term2.(Float))
    if pi != fl {
        t.Error("\n" + testName + " - Expected: " + strPi +
                "\n                          Was: " + term2.String())
    }

    // Check third term.
    term3 = c3.GetTerm(3)
    tt = term3.TermType()
    if tt != LINKEDLIST {
        t.Error("\n" + testName + " " + term3.String() + " - invalid type.")
        return
    }

    ll := term3.(LinkedListStruct)
    expected2 := 4
    actual2   := ll.GetCount()
    if expected2 != actual2 {
        msg := fmt.Sprintf("Expected list length: %v" +
                           "\n                 Was: %v", expected2, actual2)
        t.Error("\nTestParseComplex 2\n" + msg)
    }

    // Analyze the terms of the parsed linked list.
    ptr := &ll
    expectedTerms := [4]string{"a,b,c", "3.141590", "e", "$Y"}
    expectedTypes := [4]int{ATOM, FLOAT, ATOM, VARIABLE}
    i := 0
    for ptr != nil {
        actualTerm := ptr.GetTerm()
        if actualTerm == nil { break }
        if expectedTerms[i] != actualTerm.String() {
            msg := fmt.Sprintf("\nList term, expected: %v" +
                           "\n                Was: %v", expectedTerms[i], actualTerm)
            t.Error(msg)
        }
        if expectedTypes[i] != actualTerm.TermType() {
            msg := fmt.Sprintf("\nTerm type, expected: %v" +
                   "\n                Was: %v", expectedTypes[i], actualTerm.TermType())
            t.Error(msg)
        }
        ptr = ptr.GetNext()
        i++
    }

    // Check quotation mark errors. ---------------------------
    testName = "TestParseComplex 3"
    fmt.Println(testName)

    c4, e := ParseComplex("func(\"a, b, c\", d, e)")
    if e != nil {
        t.Error("func(\"a, b, c\", d, e) should not generate an error.")
    } else {
        expected := "a, b, c"
        actual := c4.GetTerm(1).String()
        if actual != expected {
            t.Error("\nFirst term should be: " + expected +
                    "\n                 Was: " + actual)
        }
    }

    _, e = ParseComplex("func(\"a, b, c, d, e)")
    checkParseComplexErrors(t, "Unmatched quotes: \"a, b, c, d, e", e)

    _, e = ParseComplex("func(a, b\"\", c, d, e)")
    checkParseComplexErrors(t, "Text before opening quote: b\"\"", e)

} // TestParseComplex


// TestComplexPanic - tests creation of complex terms which should
// or should not cause a panic.
/*
func TestComplexPanic(t *testing.T) {
    fmt.Println("TestComplexPanic")
    shouldPanic(t, ParseComplex, "verb(go, , base)")
    shouldNotPanic(t, ParseComplex, "dollar($)")
    shouldNotPanic(t, ParseComplex, "dollar(\\\"$\\\")")
}
*/

// shouldPanic - ensures that invalid argument strings cause a panic.
//     t - test argument
//     panickyFunction - function will should generate a panic
//     args - arguments which should cause a panic
func shouldPanic(t *testing.T, panickyFunction func(string) Complex, args string) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("Arguments should cause panic: " + args)
        }
    }()
    panickyFunction(args)
}

// shouldNotPanic - ensures that valid argument strings do not cause a panic.
//     t - test argument
//     panickyFunction - function will should not generate a panic
//     args - arguments which should not cause a panic
func shouldNotPanic(t *testing.T, myFunction func(string) Complex, args string) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("Arguments should not cause panic: " + args)
            fmt.Printf("    %v\n", r)
        }
    }()
    myFunction(args)
}

// checkParseComplexErrors - Check error messages generated by ParseComplex.
// Params:
//     testing
//     expected error message
//     error
func checkParseComplexErrors(t *testing.T, expected string, e error) {
    if e == nil {
        t.Error("\nShould show error message: " + expected)
        return
    }
    actual := e.Error()
    if actual != expected {
        t.Error("\nError message should be: " + expected +
                "\n                    Was: " + actual)
    }
} // checkParseComplexErrors
