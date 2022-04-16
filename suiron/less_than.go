package suiron

// Less Than - compares integers and floating point numbers.
//
//    <($X, 18)
// or
//    $X < 18
//
// In the example above, the goal will succeed if the Variable $X is
// bound to an Integer or a Float which is less than 18.
// Otherwise, the goal will fail.
//
// If either side of the comparison is ungrounded, the goal will fail.
//
// If one of the numbers is an Integer and the other is a Float, the
// Integer will be converted to a Float for the comparison.
// (It remains an Integer.)
//
// Cleve Lendon

type LessThanStruct BuiltInPredicateStruct

// LessThan - creates a comparison predicate, LessThanStruct,
// which holds the predicate's name and arguments. LessThan requires
// 2 arguments, which can be Integers or Floats.
// Params: arguments (Unifiable)
// Return: LessThanStruct
func LessThan(arguments ...Unifiable) LessThanStruct {
    return LessThanStruct {
        Name: "less_than",
        Arguments: arguments,
    }
}

// ParseLessThan - creates a LessThanStruct from a string. If the
// string does not contain "<", the function returns with the
// success flag set to false.
// If there is an error in parsing one of the terms, the function
// causes a panic.
// Params:
//     string, eg.: $X < 18
// Return:
//     less-than predicate
//     success/failure flag
func ParseLessThan(str string) (LessThanStruct, bool) {
    term1, term2, ok := parseComparison(str, "<")
    if !ok { return LessThanStruct{}, false }
    return LessThan(term1, term2), true
} // ParseLessThan


// GetSolver - gets a solution node for this predicate.
// This function satisfies the Goal interface.
func (s LessThanStruct) GetSolver(kb KnowledgeBase,
                               parentSolution SubstitutionSet,
                               parentNode SolutionNode) SolutionNode {
    return makeLessThanSolutionNode(s, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (lt LessThanStruct) RecreateVariables(vars VarMap) Expression {
    bip := BuiltInPredicateStruct(lt).RecreateVariables(vars)
    return Expression(LessThanStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (lt LessThanStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(lt).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation of this comparison.
// For example: $X < 8.
// Returns: string representation
func (lt LessThanStruct) String() string {
    return comparisonString(lt.Arguments, " < ")
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeLessThanSolutionNode()
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

type LessThanSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

// makeLessThanSolutionNode - creates a solution node for this predicate.
func makeLessThanSolutionNode(goal Goal, kb KnowledgeBase,
                                     parentSolution SubstitutionSet,
                                     parentNode SolutionNode) SolutionNode {

    node := LessThanSolutionNodeStruct{
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
func (sn *LessThanSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {

    if sn.NoBackTracking || !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.

    goal  := sn.Goal.(LessThanStruct)
    term1, type1, term2, type2 := getTermsToCompare(goal.Arguments, sn.ParentSolution)

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
func (sn *LessThanSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
// This function satisfies the SolutionNode interface.
func (sn *LessThanSolutionNodeStruct) GetParentNode() SolutionNode {
    return sn.ParentNode
}
