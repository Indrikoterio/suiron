package suiron

// Include
//
// The built-in predicate 'Include' filters terms from an input list, according
// to a filter term. Its arguments are: filter, input list, output list. Eg.
//
// ...$InList = [male(Sheldon), female(Penny), female(Bernadette), male(Leonard)]
// ...include(male($_), $InList, $OutList),...
//
// Items in the input list which are unifiable with the filter will be written
// to the output list.
//
// The output list above will contain only males, [male(Sheldon), male(Leonard)]
//
// Cleve Lendon

type IncludeStruct BuiltInPredicateStruct

// Include - creates an IncludeStruct, which holds the name and arguments.
// Include requires 3 arguments.
// Params: arguments (Unifiable)
// Return: IncludeStruct
func Include(arguments ...Unifiable) IncludeStruct {
    if len(arguments) != 3 {
        panic("Include - Requires 3 arguments.")
    }
    return IncludeStruct {
        Name: "include",
        Arguments: arguments,
    }
}

// GetSolver - gets solution node for the Include predicate.
func (is IncludeStruct) GetSolver(kb KnowledgeBase,
                                  parentSolution SubstitutionSet,
                                  parentNode SolutionNode) SolutionNode {

    return makeIncludeSolutionNode(is, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (is IncludeStruct) RecreateVariables(vars VarMap) Expression {
    bip := BuiltInPredicateStruct(is).RecreateVariables(vars)
    return Expression(IncludeStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (is IncludeStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(is).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  predicate_name(arg1, arg2, arg3)
func (is IncludeStruct) String() string {
    return BuiltInPredicateStruct(is).String()
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeIncludeSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

type IncludeSolutionNodeStruct struct {
     SolutionNodeStruct
     moreSolutions bool
}

func makeIncludeSolutionNode(goal Goal, kb KnowledgeBase,
                             parentSolution SubstitutionSet,
                             parentNode SolutionNode) SolutionNode {

    node := IncludeSolutionNodeStruct{
                SolutionNodeStruct: SolutionNodeStruct {
                                    Goal: goal,
                                    KnowledgeBase: kb,
                                    ParentSolution: parentSolution,
                                    ParentNode: parentNode },
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
func (sn *IncludeSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if sn.NoBackTracking { return nil, false }
    if !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.
    goal := sn.Goal.(IncludeStruct)
    return filterList(goal.Arguments, sn.ParentSolution)
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *IncludeSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *IncludeSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}

// pass - determines whether the given term passes the filter.
//
// Params: term to test (Unifiable)
//         filter term  (Unifiable)
//         substitution set
// Return: true == pass, false == discard
//
func pass(term Unifiable, filter Unifiable, ss SubstitutionSet) bool {
      _, ok := filter.Unify(term, ss)
      if ok { return true }
      return false
}

// filterList - scans the given input list, and tries to unify
// each item with the filter goal. If unification succeeds, the
// item is included in the output list. The output list will be
// bound to the third argument.
// Params:
//      unifiable arguments
//      substitution set
// Return:
//      new substitution set
//      success/failure flag
func filterList(arguments []Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

    var outTerms []Unifiable

    nArgs := len(arguments)
    if nArgs != 3 { return ss, false }

    filterTerm := arguments[0]
    inputList, ok := ss.CastLinkedList(arguments[1])
    if !ok { return ss, false }

    // Iterate through the input list.
    ptr := &inputList
    for ptr != nil {
        if ptr.tailVar {
            v := ptr.term.(VariableStruct)
            groundTerm, ok := ss.GetGroundTerm(v)
            if !ok { return ss, false }
            termType := groundTerm.TermType()
            if termType == LINKEDLIST {
                ll := groundTerm.(LinkedListStruct)
                ptr = &ll
                continue
            } else {
                if pass(groundTerm, filterTerm, ss) {
                    outTerms = append(outTerms, groundTerm)
                }
            }
        } else {
            if ptr.term != nil {
                if pass(ptr.term, filterTerm, ss) {
                    outTerms = append(outTerms, ptr.term)
                }
            }
        }
        ptr = ptr.next
    }

    outList := MakeLinkedList(false, outTerms...)
    return arguments[2].Unify(outList, ss)

} // filterList
