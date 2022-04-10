//package main
package suiron

// BIPTemplate
//
// This file is a template for writing built-in predicates (BIP) for Suiron.
//
// Search and replace the string 'BIPTemplate', everywhere it appears,
// with the name of your predicate. Write your predicate specific code
// in the function bipEvaluate(), and change its name to something
// meaningful. Adjust comments appropriately and rename this file.
//
// Cleve Lendon

import (
    // Uncomment this import if the predicate
    // is outside of the suiron package.
    //. "github.com/indrikoterio/suiron/suiron"
)

type BIPTemplateStruct BuiltInPredicateStruct

// BIPTemplate - creates the struct which defines this built-in predicate.
// Checks input arguments.
func BIPTemplate(arguments ...Unifiable) BIPTemplateStruct {
    if len(arguments) != 4 {
        panic("BIPTemplate - This predicate requires 4 arguments.")
    }
    return BIPTemplateStruct {
        Name: "BIPTemplate",
        Arguments: arguments,
    }
}

// GetSolver - gets a solution node for this predicate.
// This function satisfies the Goal interface.
func (s BIPTemplateStruct) GetSolver(kb KnowledgeBase,
                                    parentSolution SubstitutionSet,
                                    parentNode SolutionNode) SolutionNode {
    return makeBIPTemplateSolutionNode(s, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (s BIPTemplateStruct) RecreateVariables(
                               vars map[string]VariableStruct) Expression {
    bip := BuiltInPredicateStruct(s).RecreateVariables(vars)
    return Expression(BIPTemplateStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (s BIPTemplateStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(s).ReplaceVariables(ss)
}  // ReplaceVariables

// String - creates a string representation.
// Returns: predicate_name(arg1, arg2, arg3)
func (s BIPTemplateStruct) String() string {
    return BuiltInPredicateStruct(s).String()
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeBIPTemplateSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

// A solution node holds the current state of the search for a solution.
type BIPTemplateSolutionNodeStruct struct {
    SolutionNodeStruct
    moreSolutions bool
}

// makeBIPTemplateSolutionNode - creates a solution node for this predicate.
func makeBIPTemplateSolutionNode(goal Goal, kb KnowledgeBase,
                              parentSolution SubstitutionSet,
                              parentNode SolutionNode) SolutionNode {

    node := BIPTemplateSolutionNodeStruct{
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
func (sn *BIPTemplateSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if sn.NoBackTracking || !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.
    goal := sn.Goal.(BIPTemplateStruct)
    return bipEvaluate(goal.Arguments, sn.ParentSolution)
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *BIPTemplateSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *BIPTemplateSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}

//----------------------------------------------------------------
// bipEvaluate - does the logical work of this built-in predicate.
//
// Params:
//      list of arguments
//      substitution set (= solution so far)
// Return:
//      updated substitution set
//      success/failure flag
//
func bipEvaluate(arguments []Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

    // Do stuff here, for example:
    // ss, ok = someUnifiable.Unify(arguments[0], ss)
    // if !ok { return ss, false }
    return ss, true

} // bipEvaluate
