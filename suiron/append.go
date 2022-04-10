package suiron

// Append
//
// This built-in predicate appends terms to make a LinkedList.
// For example:
//
//   ..., $X = a, append($X, b, [c, d, e], [f, g], $Out), ...
//
// The last argument of the append predicate is the output argument.
// The variable $Out will be bound to [a, b, c, d, e, f, g].
//
// Input arguments can be Integers, Floats, Atoms, Variables,
// Complex terms, or LinkedLists.
//
// There must be at least 2 arguments.
//
// Cleve Lendon

import (
    //"fmt"
)

type AppendStruct BuiltInPredicateStruct

// Append - creates an AppendStruct, which holds the predicate's
// name and arguments. Append requires at least 2 arguments.
// Params: arguments (Unifiable)
// Return: AppendStruct
func Append(arguments ...Unifiable) AppendStruct {
    if len(arguments) < 2 {
        panic("Append - This predicate requires at least 2 arguments.")
    }
    return AppendStruct {
        Name: "append",
        Arguments: arguments,
    }
}

// GetSolver - gets solution node for Append predicate.
func (as AppendStruct) GetSolver(kb KnowledgeBase,
                                 parentSolution SubstitutionSet,
                                 parentNode SolutionNode) SolutionNode {

    return makeAppendSolutionNode(as, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (as AppendStruct) RecreateVariables(
                               vars map[string]VariableStruct) Expression {
    bip := BuiltInPredicateStruct(as).RecreateVariables(vars)
    return Expression(AppendStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (as AppendStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(as).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  predicate_name(arg1, arg2, arg3)
func (as AppendStruct) String() string {
    return BuiltInPredicateStruct(as).String()
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeAppendSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

type AppendSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

func makeAppendSolutionNode(goal Goal, kb KnowledgeBase,
                            parentSolution SubstitutionSet,
                            parentNode SolutionNode) SolutionNode {

    node := AppendSolutionNodeStruct{
                SolutionNodeStruct: SolutionNodeStruct {
                                    Goal: goal,
                                    KnowledgeBase: kb,
                                    ParentSolution: parentSolution,
                                    ParentNode: parentNode },
                moreSolutions: true,
            }
    return &node
}

// NextSolution - calls a function  to evaluate the current goal,
// based on its arguments and the substitution set.
// Returns:
//    updated substitution set
//    success/failure flag
// This function satisfies the SolutionNode interface.
func (sn *AppendSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if sn.NoBackTracking { return nil, false }
    if !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.
    goal := sn.Goal.(AppendStruct)
    return appendTerms(goal.Arguments, sn.ParentSolution)
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *AppendSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *AppendSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}


// appendTerms - appends n - 1 arguments together in a linked list
// and binds the result to last argument.
func appendTerms(arguments []Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

    length := len(arguments)

    if length < 2 { return ss, false }

    argList := []Unifiable{}

    for i := 0; i < length - 1; i++ {

        term := arguments[i]

        // Get ground term.
        tt := term.TermType()
        if tt == VARIABLE {
            if (ss.IsGroundVariable(term.(VariableStruct))) {
                term, _ = ss.GetGroundTerm(term.(VariableStruct))
            }
        }

        tt = term.TermType()
        if tt == ATOM || tt == INTEGER || tt == FLOAT {
            argList = append(argList, term)
        } else if tt == COMPLEX {
            argList = append(argList, term)
        } else if tt == LINKEDLIST {
            t := term.(LinkedListStruct)
            list := &t
            for {
               head := list.term
               if head == nil { break }
               argList = append(argList, head)
               list = list.next
               if list.term == nil { break }
            }
        }
    } // for

    outList := MakeLinkedList(false, argList...)
    lastTerm := arguments[length - 1]
    return lastTerm.Unify(outList, ss)

} // appendTerms
