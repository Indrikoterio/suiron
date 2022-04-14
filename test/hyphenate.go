package main

// Hyphenate - is a built-in predicate function which joins two words with
// a hyphen. Its purpose is to test built-in predicate functionality.
//
// Hyphenate is instantiated from bip_test.go, as follows:
//
//    c2 := Hyphenate(In, H, T, InErr, Err2)
//
// The arguments above are logic variables (type Variable). If this predicate
// were written in a Suiron source file, it would appear as follows:
//
// ..., hyphenate($In, $H, $T, $InErr, $Err2), ...
//
// Hyphenate takes the first two words from an input word list, and joins
// them with a hyphen. The new hyphenated word is bound to the second argument
// (head word). The remainder of the word list is bound to the third argument
// (tail of list). For example, if the input  were:
//
//    $In = [one, two, three, four]
//
// The output would be:
//
//    $H = one-two
//    $T = [three, four]
//
// The predicate also creates an error message, which is bound to the last
// argument. For the following input error list:
//
//    [first error]
//
// The output ($Err2) will be:
//
//    [another error, first error]
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    //"fmt"
)

type HyphenateStruct BuiltInPredicateStruct

// Hyphenate - creates a Hyphenate predicate.
func Hyphenate(arguments ...Unifiable) HyphenateStruct {
    if len(arguments) != 5 {
        panic("Hyphenate - This predicate requires 5 arguments.")
    }
    return HyphenateStruct {
        Name: "hyphenate",
        Arguments: arguments,
    }
}

// GetSolver - gets a solution node for this predicate.
// This function satisfies the Goal interface.
func (s HyphenateStruct) GetSolver(kb KnowledgeBase,
                                   parentSolution SubstitutionSet,
                                   parentNode SolutionNode) SolutionNode {
    return makeHyphenateSolutionNode(s, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (jws HyphenateStruct) RecreateVariables(vars VarMap) Expression {
    bip := BuiltInPredicateStruct(jws).RecreateVariables(vars)
    return Expression(HyphenateStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (jws HyphenateStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(jws).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  predicate_name(arg1, arg2, arg3)
func (jws HyphenateStruct) String() string {
    return BuiltInPredicateStruct(jws).String()
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeHyphenateSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

// A solution node holds the current state of the search for a solution.
// It contains the current goal, the number of the last rule fetched
// from the knowledge base, and a substitution set (which represents the
// solution so far).
// Built-in predicates produce only one solution for a given set of
// arguments. The boolean flag 'moreSolutions' is set to false after
// the first solution is returned.

type HyphenateSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

// makeHyphenateSolutionNode - creates a solution node for this predicate.
func makeHyphenateSolutionNode(goal Goal, kb KnowledgeBase,
                               parentSolution SubstitutionSet,
                               parentNode SolutionNode) SolutionNode {

    node := HyphenateSolutionNodeStruct{
                SolutionNodeStruct: MakeSolutionNode(goal, kb,
                                        parentSolution, parentNode),
                moreSolutions: true,
            }
    return &node
}

// NextSolution - calls a function (in this case, joinWithHyphen)
// to evaluate the current goal, based on its arguments and the
// substitution set.
// Returns:
//    updated substitution set
//    success/failure flag
// This function satisfies the SolutionNode interface.
func (sn *HyphenateSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if sn.NoBackTracking || !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.
    goal := sn.Goal.(HyphenateStruct)
    return joinWithHyphen(goal.Arguments, sn.ParentSolution)
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *HyphenateSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *HyphenateSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}

// joinWithHyphen - joins the first two words in a word list with
// a hyphen, and generates an error message. For example,
//
//    $In = [one, two, three, four]
//
// becomes,
//
//    $H = one-two
//    $T = [three, four]
//
// The function also takes an error list and adds a new error message, eg:
//
//    $InErr  = [first error]
//    $OutErr = [another error, first error]
//
// The 5 arguments are:
//
//    word list     - in
//    head word     - out
//    tail of list  - out
//    error list    - in
//    error list    - out
//
// Params:
//      list of 5 arguments
//      substitution set (= solution so far)
// Return:
//      updated substitution set
//      success/failure flag
//
func joinWithHyphen(arguments []Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

    // First argument must be a linked list.
    inList, ok := ss.CastLinkedList(arguments[0])
    if !ok { return ss, false }

    // The fourth argument must be a linked list.
    inErrors, ok := ss.CastLinkedList(arguments[3])
    if !ok { return ss, false }

    err := Atom("another error")
    // Add an error message to the error list
    newErrorList := MakeLinkedList(false, err, inErrors)

    // Flatten gets the first two items of a list and the
    // rest of the list. Thus, Flatten(2, ss) returns a slice
    // of 3 items: term1, term2, list of remaining terms
    terms, ok := inList.Flatten(2, ss)
    if !ok { return ss, false }

    term1, _ := ss.CastAtom(terms[0])
    term2, _ := ss.CastAtom(terms[1])

    // Join the two terms.
    str := term1.String() + "-" + term2.String()
    newHead := Atom(str)
    newTail := terms[2].(LinkedListStruct)

    // Unify output terms.
    ss, ok = newErrorList.Unify(arguments[4], ss)
    if !ok { return ss, false }
    ss, ok = newHead.Unify(arguments[1], ss)
    if !ok { return ss, false }
    ss, ok = newTail.Unify(arguments[2], ss)
    if !ok { return ss, false }
    return ss, true

} // joinWithHyphen
