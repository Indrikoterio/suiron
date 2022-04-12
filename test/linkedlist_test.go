package main

// Tests the And and Or operators of the inference engine.
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "testing"
    "fmt"
)

// Tests the 'and' and 'or' operators.
func TestLinkedList(t *testing.T) {

    fmt.Println("TestLinkedList - MakeLinkedList")

    doctor        := Atom("doctor")
    carpenter     := Atom("carpenter")
    sales_manager := Atom("sales manager")
    scientist     := Atom("scientist")

    vars := map[string]VariableStruct{}

    jobs1 := MakeLinkedList(false, doctor, carpenter, sales_manager)
    jobs2 := MakeLinkedList(false, scientist, jobs1)

    theCount := jobs2.GetCount()
    if theCount != 4 {
        msg := fmt.Sprintf("List should have count of 4. Was: %v\n", theCount)
        t.Error(msg)
    }

    if EmptyList().String() != "[]" {
        t.Error("Empty list should display as []")
    }

    expected := "[scientist, doctor, carpenter, sales manager]"
    actual   := jobs2.String()
    if actual != expected {
        t.Error("\nList should display as: " + expected +
                "\n                   Was: " + actual)
    }

    v1, _ := LogicVar("$X")
    list1 := MakeLinkedList(true, scientist, doctor, v1)

    expected = "[scientist, doctor | $X]"
    actual   = list1.String()
    if actual != expected {
        t.Error("\nList should display as: " + expected +
                "\n                   Was: " + actual)
    }

    v2, _ := LogicVar("$Y")
    list2 := MakeLinkedList(true, v1, doctor, v2)

    expected = "[$X, doctor | $Y]"
    actual   = list2.String()
    if actual != expected {
        t.Error("\nList should display as: " + expected +
                "\n                   Was: " + actual)
    }

    fmt.Println("TestLinkedList - ParseLinkedList")

    list3, _ := ParseLinkedList("[]")
    theCount = list3.GetCount()
    if theCount != 0 {
        msg := fmt.Sprintf("\nEmpty list should have count of 0. Was: %v", theCount)
        t.Error(msg)
    }

    _, err2 := ParseLinkedList("[|]")
    if err2 == nil {
        t.Error("Should produce error: Missing argument")
        return
    }

    expected = "ParseLinkedList() - Missing argument: [|]"
    checkErrorMessage(t, expected, err2.Error())

    _, err3 := ParseLinkedList("[,]")
    if err3 == nil {
        t.Error("Should produce error: Missing argument")
        return
    }
    expected = "ParseLinkedList() - Missing argument: [,]"
    checkErrorMessage(t, expected, err3.Error())

    _, err4 := ParseLinkedList("[a,]")
    if err4 == nil {
        t.Error("Should produce error: Missing argument")
        return
    }
    expected = "ParseLinkedList() - Missing argument: [a,]"
    checkErrorMessage(t, expected, err4.Error())

    _, err5 := ParseLinkedList("[a, b")
    if err5 == nil {
        t.Error("Should produce error: Missing closing bracket")
        return
    }
    expected = "ParseLinkedList() - Missing closing bracket: [a, b"
    checkErrorMessage(t, expected, err5.Error())

    //-------------------------------------------------------
    // Check mismatched quotes.
    _, err6 := ParseLinkedList("[\"a, b, c, d, e]")
    if err6 == nil {
        t.Error("Should produce error: Unmatched quotes: \"a")
        return
    }
    expected = "Unmatched quotes: \"a"
    checkErrorMessage(t, expected, err6.Error())

    _, err7 := ParseLinkedList("[a, b\"\", c, d, e]")
    if err7 == nil {
        t.Error("Should produce error: Text before opening quote: b\"\"")
        return
    }
    expected = "Text before opening quote: b\"\""
    checkErrorMessage(t, expected, err7.Error())

    expected = "[lawyer, teacher, programmer, janitor]"
    jobs3, _ := ParseLinkedList(expected)
    actual = jobs3.String()
    if actual != expected {
        msg := fmt.Sprintf("\nTest LinkedList, expected: %v\n" +
               "                      Was: %v\n", expected, actual)
        t.Error(msg)
    }

    expected = "[lawyer, teacher, programmer | $X]"
    jobs4, _ := ParseLinkedList(expected)
    actual = jobs4.String()
    if actual != expected {
        msg := fmt.Sprintf("\nTest LinkedList, expected: %v" +
               "\n                      Was: %v\n", expected, actual)
        t.Error(msg)
    }

    //-----------------------------------------------------
    fmt.Println("TestLinkedList - Flatten")

    ss := SubstitutionSet{}

    // ----- Flatten test 1. -----
    flattened, ok := jobs4.Flatten(1, ss)
    flattenErrors(t, ok, "1", flattened, "lawyer", "[teacher, programmer | $X]")

    // ----- Flatten test 2. -----
    flattened, ok = jobs4.Flatten(2, ss)
    flattenErrors(t, ok, "2", flattened, "lawyer", "[programmer | $X]")

    // ----- Flatten test 3. -----
    flattened, ok = jobs4.Flatten(3, ss)
    flattenErrors(t, ok, "3", flattened, "lawyer", "[$X]")

    // ----- Flatten test 4. -----
    flattened, ok = jobs4.Flatten(4, ss)
    flattenErrors(t, ok, "4", flattened, "lawyer", "[]")

    fmt.Println("TestLinkedList - Unify")

    // Empty lists should unify.
    empty1 := LinkedListStruct{}
    empty2 := LinkedListStruct{}
    _, ok = empty1.Unify(empty2, ss)
    if !ok {
        t.Error("Unify - empty lists should unify. [] = []")
    }

    jobs5 := MakeLinkedList(false, doctor, carpenter, sales_manager)

    _, ok = jobs1.Unify(jobs5, ss)
    if !ok {
        t.Error("Unify 2 LinkedLists. jobs1 and jobs5 should unify.")
    }

    v1 = v1.RecreateVariables(vars).(VariableStruct)
    jobs6 := MakeLinkedList(true, doctor, carpenter, v1)
    newSS, _ := jobs5.Unify(jobs6, ss)
    binding := newSS[v1.ID()].String()
    expected = "[sales manager]"
    if binding != expected {
        t.Error("Unify - $X should unify with " + expected)
    }

    newSS, _ = jobs5.Unify(v1, ss)
    binding = newSS[v1.ID()].String()
    expected = "[doctor, carpenter, sales manager]"
    if binding != expected {
        t.Error("Unify - $X should unify with " + expected)
    }
    //fmt.Printf("---- %v\n", newSS)

    //-----------------------------------------------------
    fmt.Println("TestLinkedList - Count")

    // test_count($Out) :- $R = [doctor, carpenter, sales manager],
    //                     $S = [driver | $R],
    //                     count($S, $Out). 

    driver := Atom("driver")
    Out, _   := LogicVar("$Out")
    R, _     := LogicVar("$R")
    S, _     := LogicVar("$S")
    jobs7 := MakeLinkedList(true, driver, R)

    test_count := Atom("test_count")
    head := Complex{test_count, Out}
    body := And(
                Unify(R, jobs5),
                Unify(S, jobs7),
                Count(S, Out),
            )

    // Make rule, add to knowledge base.
    r1 := Rule(head, body)
    kb := KnowledgeBase{}
    kb.Add(r1)

    goal := MakeGoal(test_count, Out)
    ss = SubstitutionSet{}

    solution, failure := Solve(goal, kb, ss)
    if len(failure) != 0 {
        t.Error("TestLinkedList - Count - " + failure)
        return
    }

    count := solution.GetTerm(1)
    intCount := count.(Integer)

    if intCount != 4 {
        msg := fmt.Sprintf("\nLinkedList - Count - expected: %v" +
               "\n                          Was: %v\n", 4, intCount)
        t.Error(msg)
    }

}  // TestLinkedList


// flattenErrors - creates error messages for testing Flatten().
// Params:  testing pointer
//          ok - true if flatten succeeded
//          test number
//          flattened linked list
//          expected first term
//          expected last term
func flattenErrors(t *testing.T,
                   ok bool,
                   testNum string,
                   flattened []Unifiable,
                   expectedFirst string, expectedLast string) {
    if !ok {
        t.Error("Flatten test " + testNum + " fails. Cannot flatten linked list.")
        return
    }
    n := len(flattened)
    str1 := flattened[0].String()
    if str1 != expectedFirst {
        t.Error("\nFlatten test " + testNum + " fails." +
                "\nFirst term should be: " + expectedFirst +
                "\n                 was: " + str1)
    }
    str2 := flattened[n - 1].String()
    if str2 != expectedLast {
        t.Error("\nFlatten test " + testNum + " fails. " +
            "\nLast term should be: " + expectedLast +
            "\n                was: " + str2)
    }

} // flattenErrors()


// checkErrorMessage - If the error message is different from expected,
// report an error.
// Params:
//     testing
//     expected message
//     actual message
func checkErrorMessage(t *testing.T, expected string, actual string) {
    if actual != expected {
        t.Error("\nError message should be: " + expected +
                "\n                    Was: " + actual)
    }
}
