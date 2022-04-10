package suiron

// Less Than or Equal - compares integers and floating point numbers.
//
//    <=($X, 18)
// or
//    $X <= 18
//
// In the example above, the goal will succeed if the Variable $X is
// bound to an Integer or a Float which is less than or equal to 18.
// Otherwise, the goal will fail.
//
// If either side of the comparison is ungrounded, the goal will fail.
//
// If one of the numbers is an Integer and the other is a Float,
// the Integer will be converted to a Float for the comparison.
// (It remains an Integer.)
//
// Cleve Lendon

import (
    //"time"
    "fmt"
)

type LessThanOrEqualStruct BuiltInPredicateStruct

const lteRequire2  = "LessThanOrEqual - requires 2 numeric arguments."
const lteNotGround = "LessThanOrEqual - variable %v is not grounded."
const lteNotNumber = "LessThanOrEqual - not a number: %v"

// LessThanOrEqual - creates a comparison predicate, LessThanOrEqualStruct,
// which holds the predicate's name and arguments. LessThanOrEqual requires
// 2 arguments, which can be Integers or Floats.
// Params: arguments (Unifiable)
// Return: LessThanOrEqualStruct
func LessThanOrEqual(arguments ...Unifiable) LessThanOrEqualStruct {
    return LessThanOrEqualStruct {
        Name: "less_than_or_equal",
        Arguments: arguments,
    }
}

// ParseLessThanOrEqual - creates a LessThanOrEqualStruct from
// a string. If the string does not contain "<=", the function
// returns with the success flag set to false.
// If there is an error in parsing one of the terms, the function
// causes a panic.
// Params:
//     string, eg.: $X <= 18
// Return:
//     less-than-or-equal predicate
//     success/failure flag
func ParseLessThanOrEqual(str string) (LessThanOrEqualStruct, bool) {
    runes := []rune(str)
    index := specialIndexOf(runes, []rune{'<', '='})
    if index == -1 { return LessThanOrEqualStruct{}, false }  // Not Less Than or Equal.
    arg1 := runes[0: index]
    arg2 := runes[index + 2:]
    term1, err := parseTerm(string(arg1))
    if err != nil { panic(err.Error()) }
    term2, err := parseTerm(string(arg2))
    if err != nil { panic(err.Error()) }
    return LessThanOrEqual(term1, term2), true
}


// GetSolver - gets a solution node for this predicate.
// This function satisfies the Goal interface.
func (s LessThanOrEqualStruct) GetSolver(kb KnowledgeBase,
                               parentSolution SubstitutionSet,
                               parentNode SolutionNode) SolutionNode {
    return makeLessThanOrEqualSolutionNode(s, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (lte LessThanOrEqualStruct) RecreateVariables(
                               vars map[VariableStruct]VariableStruct) Expression {
    bip := BuiltInPredicateStruct(lte).RecreateVariables(vars)
    return Expression(LessThanOrEqualStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (lte LessThanOrEqualStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(lte).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  "arg1 = arg2"
func (lte LessThanOrEqualStruct) String() string {
    term1 := lte.Arguments[0].String()
    term2 := lte.Arguments[1].String()
    return term1 + " <= " + term2
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeLessThanOrEqualSolutionNode()
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

type LessThanOrEqualSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

// makeLessThanOrEqualSolutionNode - creates a solution node for this predicate.
func makeLessThanOrEqualSolutionNode(goal Goal, kb KnowledgeBase,
                                     parentSolution SubstitutionSet,
                                     parentNode SolutionNode) SolutionNode {

    node := LessThanOrEqualSolutionNodeStruct{
                SolutionNodeStruct: MakeSolutionNode(goal, kb,
                                                     parentSolution,
                                                     parentNode),
                moreSolutions: true,
            }
    return &node
}

// NextSolution - compares two numbers, float or integer.
// Returns:
//    updated substitution set
//    success/failure flag
// This function satisfies the SolutionNode interface.
func (sn *LessThanOrEqualSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {

    if sn.NoBackTracking || !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.

    goal  := sn.Goal.(LessThanOrEqualStruct)
    term1 := goal.Arguments[0]
    term2 := goal.Arguments[1]

    ground1, ok := sn.ParentSolution.GetGroundTerm(term1)
    if !ok {
        msg := fmt.Sprintf(lteNotGround, term1)
        panic(msg)
    }

    ground2, ok := sn.ParentSolution.GetGroundTerm(term2)
    if !ok {
        msg := fmt.Sprintf(lteNotGround, term2)
        panic(msg)
    }

    termType1 := ground1.TermType()
    termType2 := ground2.TermType()

    if termType1 != FLOAT && termType1 != INTEGER {
        msg := fmt.Sprintf(lteNotNumber, ground1)
        panic(msg)
    }

    if termType2 != FLOAT && termType2 != INTEGER {
        msg := fmt.Sprintf(lteNotNumber, ground2)
        panic(msg)
    }

    if termType1 == INTEGER && termType2 == INTEGER {
        number1 := ground1.(Integer)
        number2 := ground2.(Integer)
        if number1 <= number2 { return sn.ParentSolution, true }
        return sn.ParentSolution, false
    }

    if termType1 == FLOAT && termType2 == FLOAT {
        number1 := ground1.(Float)
        number2 := ground2.(Float)
        if number1 <= number2 { return sn.ParentSolution, true }
        return sn.ParentSolution, false
    }

    if termType1 == FLOAT && termType2 == INTEGER {
        number1 := float64(ground1.(Float))
        number2 := float64(ground2.(Integer))
        if number1 <= number2 { return sn.ParentSolution, true }
        return sn.ParentSolution, false
    }

    if termType1 == INTEGER && termType2 == FLOAT {
        number1 := float64(ground1.(Integer))
        number2 := float64(ground2.(Float))
        if number1 <= number2 { return sn.ParentSolution, true }
        return sn.ParentSolution, false
    }

    return sn.ParentSolution, false

} // NextSolution


// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *LessThanOrEqualSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
// This function satisfies the SolutionNode interface.
func (sn *LessThanOrEqualSolutionNodeStruct) GetParentNode() SolutionNode {
    return sn.ParentNode
}
