package suiron

// Functor
//
// Functor is a built-in predicate to get the functor and arity
// of a complex term. Eg.:
//
//     functor(boss(Zack, Stephen), $Func, $Arity)
//
// The first term must be the complex term to be tested.
//
// $Func will bind to 'boss' and $Arity will bind to the integer
// 2 (because there are two arguments, Zack and Stephen). Arity
// is optional:
//
//     functor(boss(Zack, Stephen), $Func)
//
// The following goal will succeed.
//
//     $X = boss(Zack, Stephen), functor($X, boss)
//
// The next goal will not succeed, because the arity is incorrect:
//
//     functor($X, boss, 3)
//
// If the second argument has an asterisk at the end, the match will
// test only the start of the string. For example, the following
// will succeed:
//
//     $X = noun_phrase(the blue sky), functor($X, noun*)
//
// TODO:
// Perhaps the functionality could be expanded to accept a regex
// string for the second argument.
//
// Cleve Lendon

import (
    "strings"
)

type FunctorStruct BuiltInPredicateStruct

// Functor - creates a FunctorStruct, which holds the predicate's
// name and arguments. Functor requires 2 or 3 arguments
// Params: arguments (Unifiable)
// Return: FunctorStruct
func Functor(arguments ...Unifiable) FunctorStruct {
    if len(arguments) < 2 {
        panic("Functor - Requires at least 2 arguments.")
    }
    if len(arguments) > 3 {
        panic("Functor - Too many arguments.")
    }
    return FunctorStruct {
        Name: "functor",
        Arguments: arguments,
    }
}

// GetSolver - gets solution node for Functor predicate.
func (as FunctorStruct) GetSolver(kb KnowledgeBase,
                                  parentSolution SubstitutionSet,
                                  parentNode SolutionNode) SolutionNode {

    return makeFunctorSolutionNode(as, kb, parentSolution, parentNode)
}

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - Refer to comments in expression.go.
func (fs FunctorStruct) RecreateVariables(vars VarMap) Expression {
    bip := BuiltInPredicateStruct(fs).RecreateVariables(vars)
    return Expression(FunctorStruct(*bip))
}

// ReplaceVariables - Refer to comments in expression.go.
func (fs FunctorStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return BuiltInPredicateStruct(fs).ReplaceVariables(ss)
}  // ReplaceVariables


// String - creates a string representation.
// Returns:  predicate_name(arg1, arg2, arg3)
func (fs FunctorStruct) String() string {
    return BuiltInPredicateStruct(fs).String()
}

//----------------------------------------------------------------
// Solution Node functions.
//    makeFunctorSolutionNode()
//    NextSolution()
//    SetNoBackTracking()
//----------------------------------------------------------------

type FunctorSolutionNodeStruct struct {
     SolutionNodeStruct
     moreSolutions bool
}

func makeFunctorSolutionNode(goal Goal, kb KnowledgeBase,
                             parentSolution SubstitutionSet,
                             parentNode SolutionNode) SolutionNode {

    node := FunctorSolutionNodeStruct{
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
func (sn *FunctorSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {
    if sn.NoBackTracking { return nil, false }
    if !sn.moreSolutions { return nil, false }
    sn.moreSolutions = false  // Only one solution.
    goal := sn.Goal.(FunctorStruct)
    return evaluate(goal.Arguments, sn.ParentSolution)
}

// SetNoBackTracking - set the NoBackTracking flag,
// which is used to implement Cuts.
// This function satisfies the SolutionNode interface.
func (sn *FunctorSolutionNodeStruct) SetNoBackTracking() {
    sn.NoBackTracking = true
}

// GetParentNode
func (n *FunctorSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}


// evaluate - determines the functor and arity of the first argument.
// Binds the functor to the second argument, and the arity to the
// third argument, if there is one. Returns the new substitution set
// and a success/failure flag.
// Params:
//      unifiable arguments
//      substitution set
// Return:
//      new substitution set
//      success/failure flag
func evaluate(arguments []Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

    nArgs := len(arguments)
    if nArgs < 2 || nArgs > 3 { return ss, false }

    first, ok := ss.CastComplex(arguments[0])
    if !ok { return ss, false }

    functor := first.GetFunctor()
    strFunc := string(functor)

    newSS := ss

    // Get second argument (functor)
    second := arguments[1]
    tt := second.TermType()
    if tt == ATOM {
        str := second.String()
        len := len(str)
        if str[len - 1] == '*' {
            str2 := str[0: len-1]
            if !strings.HasPrefix(strFunc, str2) { return ss, false }
        } else {
            if strFunc != str { return ss, false }
        }
    } else {
        newSS, ok = second.Unify(functor, ss)
        if !ok { return newSS, false}
    }

    if nArgs == 3 {
        third := arguments[2]
        arity := first.Arity()
        return third.Unify(Integer(arity), newSS)
    }

    return newSS, true

} // evaluate
