package suiron

// Count
//
// This predicate counts the number of items in a linked list.
// The first argument is the list to check. The second is the result.
// For example:
//
//     ..., count([a, b, c], $X),...
//
// The variable $X will be bound to 3.
//
// If the last item in the list is a tail variable, which is bound
// to another list, the predicate 'count' will count the number of
// items in both lists.
//
// Cleve Lendon

type CountStruct BuiltInPredicateStruct

// Count - creates the struct which defines this built-in predicate.
// Checks input arguments.
func Count(arguments ...Unifiable) CountStruct {
    if len(arguments) != 2 {
        panic("Count - This predicate requires 2 arguments.")
    }
    return CountStruct {
        Name: "count",
        Arguments: arguments,
    }
}

// GetSolver - gets a solution node for this predicate.
// This function satisfies the Goal interface.
func (s CountStruct) GetSolver(kb KnowledgeBase,
                                    parentSolution SubstitutionSet,
                                    parentNode SolutionNode) SolutionNode {
    return makeCountSolutionNode(s, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (s CountStruct) RecreateVariables(
                               vars map[string]VariableStruct) Expression {
    bip := BuiltInPredicateStruct(s).RecreateVariables(vars)
    return Expression(CountStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (s CountStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(s).ReplaceVariables(ss)
}  // ReplaceVariables

// String - creates a string representation.
// Returns: predicate_name(arg1, arg2, arg3)
func (s CountStruct) String() string {
    return BuiltInPredicateStruct(s).String()
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeCountSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

// A solution node holds the current state of the search for a solution.
type CountSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

// makeCountSolutionNode - creates a solution node for this predicate.
func makeCountSolutionNode(goal Goal, kb KnowledgeBase,
                              parentSolution SubstitutionSet,
                              parentNode SolutionNode) SolutionNode {

    node := CountSolutionNodeStruct{
                SolutionNodeStruct: MakeSolutionNode(goal, kb,
                                        parentSolution, parentNode),
                moreSolutions: true,
            }
    return &node
}

// NextSolution - calls a function to evaluate the current goal,
// based on its arguments and the substitution set.
// Returns:
//    updated substitution set
//    success/failure flag
// This function satisfies the SolutionNode interface.
func (sn *CountSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if sn.NoBackTracking || !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.
    goal := sn.Goal.(CountStruct)
    return countLL(goal.Arguments, sn.ParentSolution)
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *CountSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *CountSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}

//----------------------------------------------------------------
// countLL - counts the number of items in a linked list.
// If the last item is a tail variable, eg. [a, b | $T], and
// that variable is bound to another list, count the terms in
// the second list also.
//
// Params:
//      list of arguments
//      substitution set (= solution so far)
// Return:
//      updated substitution set
//      success/failure flag
//
func countLL(arguments []Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

    linkedList, ok := ss.CastLinkedList(arguments[0])
    if !ok { return ss, ok }

    ptr := &linkedList
    count := 0
    for ptr != nil {
        if ptr.tailVar {
            v := ptr.term.(VariableStruct)
            groundTerm, ok := ss.GetGroundTerm(v)
            if !ok { return ss, false }
            termType := groundTerm.TermType()
            if termType == LINKEDLIST {
                ll := groundTerm.(LinkedListStruct)
                ptr = &ll
            } else {
                count++
            }
        } else {
            count++
        }
        ptr = ptr.next
    }

    intCount := Integer(count)
    return arguments[1].Unify(intCount, ss)

} // countLL
