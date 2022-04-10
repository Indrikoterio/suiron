package suiron

// KnowledgeBase - defines a dictionary of Prolog-like rules and facts.
// The dictionary (a map) is indexed by a key, which is created from the
// functor and arity. For example, for the fact mother(Carla, Caitlyn),
// the key would be "mother/2".
// Each key indexes a slice of Rules which have the same key.
//
// Cleve Lendon

import (
    "sort"
    "strings"
    "fmt"
)

// A knowledge base is a dictionary indexed by a key.
// Each indexed item is a slice of rules and/or facts.
type KnowledgeBase map[string][]RuleStruct

// Add - adds facts and rules to the knowledge base.
// Eg.  knowledgebase.Add(fact1, fact2, rule1, rule2)
func (kb KnowledgeBase) Add(rules ...RuleStruct) {
    for _, rule := range rules {
        key := rule.Key()
        sliceOfRules, found := kb[key]
        if !found {
            sliceOfRules = []RuleStruct{}
            kb[key] = append(sliceOfRules, rule)
        } else {
            kb[key] = append(sliceOfRules, rule)
        }
    }
} // Add


// GetRule - fetches a rule (or fact) from the knowledge base.
// Rules are indexed by functor/arity (eg. sister/2) and by index number.
// The variables of the retrieved rule must be made unique, by calling
// createVariables().
func (kb KnowledgeBase) GetRule(goal Goal, i int) RuleStruct {
    key := goal.(Complex).Key()
    list, ok := kb[key]
    if !ok {
        // Should never happen.
        panic("KnowledgeBase, GetRule - rule does not exist: " + key + "\n")
    }
    if i >= len(list) {
        // Should never happen.
        msg := fmt.Sprintf("KnowledgeBase, GetRule - index out of range: %v %d\n", key, i)
        panic(msg)
    }
    rule := list[i]
    rule2 := rule.RecreateVariables(make(map[string]VariableStruct))
    return rule2.(RuleStruct)
}

// FormatKB - formats the knowledge base facts and rules for display.
// This method is useful for diagnostics. The keys are sorted.
func (kb KnowledgeBase) FormatKB() string {
    var sb strings.Builder
    sb.WriteString("\n########## Contents of Knowledge Base ##########\n")
    keys := make([]string, 0, len(kb))
    for k := range kb { keys = append(keys, k) }
    sort.Strings(keys)
    for _, k := range keys {
        sb.WriteString(k + "\n")
        for i := 0; i < len(kb[k]); i++ {
            sb.WriteString("    " + kb[k][i].String() + "\n")
        }
    }
    return sb.String()
} // FormatKB


// getRuleCount - counts the number of rules for the given goal.
// When the execution time has been exceeded, this function will
// return 0. Zero indicates that all rules have been exhausted.
// Params:  goal
// Returns: count
func (kb KnowledgeBase) getRuleCount(goal Goal) int {

    // The following line will stop the search for a solution if
    // the execution time (300 msecs by default) has timed out.
    if suironHasTimedOut { return 0 }

    key := goal.(Complex).Key()
    listOfRules := kb[key]
    return len(listOfRules)

} // getRuleCount

