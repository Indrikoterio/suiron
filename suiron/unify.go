package suiron

// Unify - defines the unification predicate, which attempts to unify
// (bind) two terms. If the terms can be unified, the predicate succeeds.
// In a Suiron source file, unification can be defined as follows:
//
//    unify($X, pronoun)
// or
//    $X = pronoun
//
// In the examples above, unification will succeed if $X is unbound.
// Otherwise not.
//
// Note: Sometimes this is referred to as the unification operator, but
// it's actually a predicate.
//
// Cleve Lendon

import (
    //"fmt"
)

type UnifyStruct BuiltInPredicateStruct

// Unify - creates a unification predicate (UnifyStruct).
func Unify(arguments ...Unifiable) UnifyStruct {
    if len(arguments) != 2 {
        panic("Unify - This predicate requires 2 arguments.")
    }
    return UnifyStruct {
        Name: "unify",
        Arguments: arguments,
    }
}

// ParseUnify - creates a logical Unify operator from a string.
// If the string does not contain "=", the function returns with
// the success flag set to false.
// If there is an error in parsing one of the terms, the function
// causes a panic.
// Params:
//     string, eg.: $X = verb
// Return:
//     unify predicate
//     success/failure flag
func ParseUnify(str string) (UnifyStruct, bool) {
    runes := []rune(str)
    index := specialIndexOf(runes, []rune{'='})
    if index == -1 { return UnifyStruct{}, false }  // Not a Unify.
    arg1 := runes[0: index]
    arg2 := runes[index + 1:]
    term1, err := parseTerm(string(arg1))
    if err != nil { panic(err.Error()) }
    term2, err := parseTerm(string(arg2))
    if err != nil { panic(err.Error()) }
    return Unify(term1, term2), true
}


// GetSolver - gets a solution node for this predicate.
// This function satisfies the Goal interface.
func (s UnifyStruct) GetSolver(kb KnowledgeBase,
                               parentSolution SubstitutionSet,
                               parentNode SolutionNode) SolutionNode {
    return makeUnifySolutionNode(s, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (us UnifyStruct) RecreateVariables(vars VarMap) Expression {
    bip := BuiltInPredicateStruct(us).RecreateVariables(vars)
    return Expression(UnifyStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (us UnifyStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(us).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  "arg1 = arg2"
func (us UnifyStruct) String() string {
    term1 := us.Arguments[0].String()
    term2 := us.Arguments[1].String()
    return term1 + " = " + term2
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeUnifySolutionNode()
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

type UnifySolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}


// makeUnifySolutionNode - creates a solution node for this predicate.
func makeUnifySolutionNode(goal Goal, kb KnowledgeBase,
                           parentSolution SubstitutionSet,
                           parentNode SolutionNode) SolutionNode {

    node := UnifySolutionNodeStruct{
                SolutionNodeStruct: MakeSolutionNode(goal, kb,
                                                     parentSolution,
                                                     parentNode),
                moreSolutions: true,
            }

    return &node
}

// NextSolution - calls Unify to attempt to unify two terms.
// Returns:
//    updated substitution set
//    success/failure flag
// This function satisfies the SolutionNode interface.
func (sn *UnifySolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if sn.NoBackTracking || !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.
    goal  := sn.Goal.(UnifyStruct)
    term1 := goal.Arguments[0]
    term2 := goal.Arguments[1]
    return term1.Unify(term2, sn.ParentSolution)
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *UnifySolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *UnifySolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}
