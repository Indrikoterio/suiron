package suiron

// Complex Solution Node
//
// Solution node for complex terms (= compound terms).
//
// This node has a child node, which is the next subgoal. The method
// NextSolution() will check to see if the child has a solution.
// If it does, the solution (substitution set) is returned, with
// the success flag set to true.
//
// Otherwise, NextSolution() fetches rules/facts from the knowledge
// base, and tries to unify the head of these rules and facts with
// the goal. If a matching fact is found, the solution is returned,
// with the success flag set to true.
//
// (Note, a fact is a rule without a body.)
//
// Otherwise, the body node of the rule becomes the child node, and
// the algorithm tries to find a solution (substitution set) for the
// child. It will return the child solution with the success flag
// set to true, or set to false for failure.
//
// Cleve Lendon

import (
    //"fmt"
)

type ComplexSolutionNodeStruct struct {
    SolutionNodeStruct
    child SolutionNode
}

func MakeComplexSolutionNode(g Complex, kb KnowledgeBase,
                             parentSolution SubstitutionSet,
                             parentNode SolutionNode) SolutionNode {

    node := ComplexSolutionNodeStruct{
                SolutionNodeStruct: SolutionNodeStruct{
                                    Goal: Goal(g),
                                    KnowledgeBase: kb,
                                    ParentSolution: parentSolution,
                                    ParentNode: nil },
                child: nil,
            }
    // Count the number of rules or facts which match the goal.
    node.count = kb.getRuleCount(Goal(g))

    // For debugging.
    //if node.count == 0 {
    //    fmt.Printf("Missing rule: " + g.Key())
    //}

    return &node
}

// NextSolution - initiates or continues the search for a solution.
// If the search succeeds, the method returns the updated substitution
// set, and sets the success flag to true.
// If the search fails, the success flag is set to false.
func (n *ComplexSolutionNodeStruct) NextSolution() (SubstitutionSet, bool) {

    if n.NoBackTracking { return nil, false }

    if n.child != nil {
        solution, found := n.child.NextSolution()
        if found { return solution, true }
    }

    n.child = nil

    for n.HasNextRule() == true {
        rule := n.NextRule()

        head := rule.GetHead()
        solution, success := head.Unify(n.Goal.(Unifiable), n.ParentSolution)

        if success {
            body := rule.body
            if body == nil { return solution, true }
            n.child = body.GetSolver(n.KnowledgeBase, solution, n);
            childSolution, ok := n.child.NextSolution()
            if ok { return childSolution, true }
        }
    }
    return nil, false

} // NextSolution()

// HasNextRule - returns true if the knowledge base contains untried
// rules for this node's goal. False otherwise.
func (n *ComplexSolutionNodeStruct) HasNextRule() bool {
    if n.NoBackTracking { return false }
    return n.ruleNumber < n.count
}

// NextRule - fetches the next rule from the database, according to ruleNumber.
// The method HasNextRule must called to ensure that a rule can be fetched
// from the knowledge base. If GetRule is called with invalid parameters, the
// knowledge base will panic.
func (n *ComplexSolutionNodeStruct) NextRule() RuleStruct {
    rule := n.KnowledgeBase.GetRule(n.Goal, n.ruleNumber)
    n.ruleNumber += 1
    return rule
}

// SetNoBackTracking - rset the NoBackTracking flag.
// This flag is used to implement Cuts.
func (n *ComplexSolutionNodeStruct) SetNoBackTracking() {
    n.NoBackTracking = true
}

// GetParentNode
func (n *ComplexSolutionNodeStruct) GetParentNode() SolutionNode {
    return n.ParentNode
}

// GetChild - returns the child node.
func (n *ComplexSolutionNodeStruct) GetChild() SolutionNode {
    return n.child
}
