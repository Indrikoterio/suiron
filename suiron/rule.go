package suiron

// Rule - Defines a rule or fact.
// Rules have the format: head :- body.
// Facts are the same as rules without a body. For example:
//   grandfather(G, C) :- father(G, A), parent(A, C).   <-- rule
//   father(john, kaitlyn).                             <-- fact
//
// Cleve Lendon

import (
    "strings"
    "fmt"
)

// head :- body.
type RuleStruct struct {
    head Complex
    body Goal
}

// Rule - Factory function to create a Rule.
func Rule(head Complex, body Goal) RuleStruct {
    return RuleStruct{ head: head, body: body }
}

// Fact - Factory function to create a Fact.
func Fact(head Complex) RuleStruct {
    return RuleStruct{ head: head, body: nil }
}

// ParseRule - create a fact or rule from a string representation.
// Examples of usage:
//    c := ParseRule("male(Harold).")
//    c := ParseRule("father($X, $Y) :- parent($X, $Y), male($X).")
func ParseRule(str string) (RuleStruct, error) {

    s := strings.TrimSpace(str)
    length := len(s)
    if length < 4 {
        err := fmt.Errorf("ParseRule() - Invalid string. >%v<\n", s)
        return RuleStruct{}, err
    }

    // Remove final period.
    ch := s[length - 1:]
    if ch == "." {
        s = s[0: length - 1]
        length = length - 1
    }
    index := strings.Index(s, ":-")

    if (index > -1) {

        strHead := s[0: index]
        strBody := s[index + 2:]

        // Make sure there is not a second ':-'.
        index2 := strings.Index(strBody, ":-")
        if index2 >= 0 {
            err := fmt.Errorf("ParseRule() - Invalid rule.\n%v", s)
            return RuleStruct{}, err
        }

        head, err := ParseSubgoal(strHead)
        if err != nil {
            return RuleStruct{}, err
        }
        
        body := generateGoal(strBody)
        return RuleStruct{head.(Complex), body}, nil

    } else { // Must be a fact (no body).
        head, err := ParseComplex(s)
        rule := RuleStruct{ head: head, body: nil }
        return rule, err
    }

} // ParseRule


// Key - generates a key from the head term.
// Eg. loves(Chandler, Monica) --> loves/2
func (r RuleStruct) Key() string {
    return r.head.Key()
}

// GetHead - returns the head of this rule, which is Complex type.
func (r RuleStruct) GetHead() Complex { return r.head }

// GetBody - returns the body of this rule, which is Goal type.
func (r RuleStruct) GetBody() Goal { return r.body }

// RecreateVariables - In Prolog, and in this inference engine, the scope of
// a logic variable is the rule or goal in which it is defined. When the
// algorithm tries to solve a goal, it calls this method to ensure that the
// variables are unique. See comments in expression.go.
func (r RuleStruct) RecreateVariables(vars map[string]VariableStruct) Expression {
    newHead := r.head.RecreateVariables(vars).(Complex)
    var newBody Goal = nil
    if r.body != nil {
        newBody = r.body.RecreateVariables(vars).(Goal)
    }
    return RuleStruct{ head: newHead, body: newBody }
} // RecreateVariables

// ReplaceVariables - replaces a bound variable with its binding.
// This method is used for displaying final results.
// Refer to comments in expression.go.
func (r RuleStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    return r.body.ReplaceVariables(ss)
} // ReplaceVariables

// String - displays a rule in its source form, eg. father(Anakin, Luke).
func (r RuleStruct) String() string {
    var sb strings.Builder
    sb.WriteString(r.head.String())
    if r.body != nil {
        sb.WriteString(" :- ")
        sb.WriteString(r.body.String())
    }
    sb.WriteString(".")
    return sb.String()
} // String
