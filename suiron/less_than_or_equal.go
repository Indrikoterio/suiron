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
// If one of the numbers is an Integer and the other is a Float, the
// Integer will be converted to a Float for the comparison.
// (It remains an Integer.)
//
// Cleve Lendon

type LessThanOrEqualStruct BuiltInPredicateStruct

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
    infix, index := identifyInfix(runes)
    if infix != LESS_THAN_OR_EQUAL { return LessThanOrEqualStruct{}, false }
    term1, term2 := getLeftAndRight(runes, index, 2)
    return LessThanOrEqual(term1, term2), true
} // ParseLessThanOrEqual


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
func (lte LessThanOrEqualStruct) RecreateVariables(vars VarMap) Expression {
    bip := BuiltInPredicateStruct(lte).RecreateVariables(vars)
    return Expression(LessThanOrEqualStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (lte LessThanOrEqualStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(lte).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation of this comparison.
// For example: $X <= 8.
// Returns: string representation
func (lte LessThanOrEqualStruct) String() string {
    return comparisonString(lte.Arguments, " <= ")
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
    term1, type1, term2, type2 := getTermsToCompare(goal.Arguments, sn.ParentSolution)

    if type1 == ATOM || type2 == ATOM {
        result := compareAtoms(term1, type1, term2, type2)
        // If less than or equal.
        if result != 1 { return sn.ParentSolution, true }
        return sn.ParentSolution, false
    }

    if type1 == INTEGER && type2 == INTEGER {
        number1 := term1.(Integer)
        number2 := term2.(Integer)
        if number1 <= number2 { return sn.ParentSolution, true }
    } else {
        number1, number2 := twoFloats(term1, type1, term2, type2)
        if number1 <= number2 { return sn.ParentSolution, true }
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
