package suiron

// Greater Than or Equal - compares integers and floating point numbers.
//
//    >=($X, 18)
// or
//    $X >= 18
//
// In the example above, the goal will succeed if the Variable $X is
// bound to an Integer or a Float which is greater than or equal to 18.
// Otherwise, the goal will fail.
//
// If either side of the comparison is ungrounded, the goal will fail.
//
// If one of the numbers is an Integer and the other is a Float, the
// Integer will be converted to a Float for the comparison.
// (It remains an Integer.)
//
// Cleve Lendon

type GreaterThanOrEqualStruct BuiltInPredicateStruct

// GreaterThanOrEqual - creates a comparison predicate, GreaterThanOrEqualStruct,
// which holds the predicate's name and arguments. GreaterThanOrEqual requires
// 2 arguments, which can be Integers or Floats.
// Params: arguments (Unifiable)
// Return: GreaterThanOrEqualStruct
func GreaterThanOrEqual(arguments ...Unifiable) GreaterThanOrEqualStruct {
    return GreaterThanOrEqualStruct {
        Name: "greater_than_or_equal",
        Arguments: arguments,
    }
}

// ParseGreaterThanOrEqual - creates a GreaterThanOrEqualStruct from
// a string. If the string does not contain ">=", the function
// returns with the success flag set to false.
// If there is an error in parsing one of the terms, the function
// causes a panic.
// Params:
//     string, eg.: $X >= 18
// Return:
//     greater-than-or-equal predicate
//     success/failure flag
func ParseGreaterThanOrEqual(str string) (GreaterThanOrEqualStruct, bool) {
    runes := []rune(str)
    infix, index := identifyInfix(runes)
    if infix != GREATER_THAN_OR_EQUAL { return GreaterThanOrEqualStruct{}, false }
    term1, term2 := getLeftAndRight(runes, index, 2)
    return GreaterThanOrEqual(term1, term2), true
} // ParseGreaterThanOrEqual


// GetSolver - gets a solution node for this predicate.
// This function satisfies the Goal interface.
func (s GreaterThanOrEqualStruct) GetSolver(kb KnowledgeBase,
                               parentSolution SubstitutionSet,
                               parentNode SolutionNode) SolutionNode {
    return makeGreaterThanOrEqualSolutionNode(s, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (gte GreaterThanOrEqualStruct) RecreateVariables(vars VarMap) Expression {
    bip := BuiltInPredicateStruct(gte).RecreateVariables(vars)
    return Expression(GreaterThanOrEqualStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (gte GreaterThanOrEqualStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(gte).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation of this comparison.
// For example: $X >= 8.
// Returns: string representation
func (gte GreaterThanOrEqualStruct) String() string {
    return comparisonString(gte.Arguments, " >= ")
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeGreaterThanOrEqualSolutionNode()
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

type GreaterThanOrEqualSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

// makeGreaterThanOrEqualSolutionNode - creates a solution node for this predicate.
func makeGreaterThanOrEqualSolutionNode(goal Goal, kb KnowledgeBase,
                                     parentSolution SubstitutionSet,
                                     parentNode SolutionNode) SolutionNode {

    node := GreaterThanOrEqualSolutionNodeStruct{
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
func (sn *GreaterThanOrEqualSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {

    if sn.NoBackTracking || !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.

    goal  := sn.Goal.(GreaterThanOrEqualStruct)
    term1, type1, term2, type2 := getTermsToCompare(goal.Arguments, sn.ParentSolution)

    if type1 == ATOM || type2 == ATOM {
        result := compareAtoms(term1, type1, term2, type2)
        // If greater than or equal.
        if result != -1 { return sn.ParentSolution, true }
        return sn.ParentSolution, false
    }

    if type1 == INTEGER && type2 == INTEGER {
        number1 := term1.(Integer)
        number2 := term2.(Integer)
        if number1 >= number2 { return sn.ParentSolution, true }
    } else {
        number1, number2 := twoFloats(term1, type1, term2, type2)
        if number1 >= number2 { return sn.ParentSolution, true }
    }

    return sn.ParentSolution, false

} // NextSolution


// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *GreaterThanOrEqualSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
// This function satisfies the SolutionNode interface.
func (sn *GreaterThanOrEqualSolutionNodeStruct) GetParentNode() SolutionNode {
    return sn.ParentNode
}
