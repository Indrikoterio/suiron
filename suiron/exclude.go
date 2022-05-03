package suiron

// Exclude
//
// The built-in predicate 'Exclude' filters terms from an input list, according
// to a filter term. Its arguments are: filter, input list, output list. Eg.
//
// ...
// $InList = [male(Sheldon), female(Penny), female(Bernadette), male(Leonard)],
// exclude(male($_), $InList, $OutList),
//
// Items in the input list which are unifiable with the filter will NOT be
// written to the output list.
//
// The output list above will contain only females,
//    [female(Penny), female(Bernadette)]
//
// Cleve Lendon

type ExcludeStruct BuiltInPredicateStruct

// Exclude - creates an ExcludeStruct, which holds the name and arguments.
// Exclude requires 3 arguments.
// Params: arguments (Unifiable)
// Return: ExcludeStruct
func Exclude(arguments ...Unifiable) ExcludeStruct {
    if len(arguments) != 3 {
        panic("Exclude - Requires 3 arguments.")
    }
    return ExcludeStruct {
        Name: "exclude",
        Arguments: arguments,
    }
}

// GetSolver - gets solution node for the Exclude predicate.
func (xs ExcludeStruct) GetSolver(kb KnowledgeBase,
                                  parentSolution SubstitutionSet,
                                  parentNode SolutionNode) SolutionNode {

    return makeExcludeSolutionNode(xs, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (xs ExcludeStruct) RecreateVariables(vars VarMap) Expression {
    bip := BuiltInPredicateStruct(xs).RecreateVariables(vars)
    return Expression(ExcludeStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (xs ExcludeStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(xs).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  predicate_name(arg1, arg2, arg3)
func (xs ExcludeStruct) String() string {
    return BuiltInPredicateStruct(xs).String()
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeExcludeSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

type ExcludeSolutionNodeStruct struct {
     SolutionNodeStruct
     moreSolutions bool
}

func makeExcludeSolutionNode(goal Goal, kb KnowledgeBase,
                             parentSolution SubstitutionSet,
                             parentNode SolutionNode) SolutionNode {

    node := ExcludeSolutionNodeStruct{
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
func (sn *ExcludeSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if sn.NoBackTracking { return nil, false }
    if !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.
    goal := sn.Goal.(ExcludeStruct)
    return excludeFromList(goal.Arguments, sn.ParentSolution)
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *ExcludeSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *ExcludeSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}

// passFilter - determines whether the given term passes the filter.
//
// Params: term to test (Unifiable)
//         filter term  (Unifiable)
//         substitution set
// Return: true == pass, false == discard
//
func passFilter(term Unifiable, filter Unifiable, ss SubstitutionSet) bool {
      _, ok := filter.Unify(term, ss)
      if !ok { return true }
      return false
}

// excludeFromList - scans the given input list, and tries to unify
// each item with the filter goal. If unification succeeds, the
// item is excluded from the output list. The output list will be
// bound to the third argument.
// Params:
//      unifiable arguments
//      substitution set
// Return:
//      new substitution set
//      success/failure flag
func excludeFromList(arguments []Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

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
                if passFilter(groundTerm, filterTerm, ss) {
                    outTerms = append(outTerms, groundTerm)
                }
            }
        } else {
            if ptr.term != nil {
                if passFilter(ptr.term, filterTerm, ss) {
                    outTerms = append(outTerms, ptr.term)
                }
            }
        }
        ptr = ptr.next
    }

    outList := MakeLinkedList(false, outTerms...)
    return arguments[2].Unify(outList, ss)

} // excludeFromList
