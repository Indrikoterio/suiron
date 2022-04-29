package suiron

// Solutions
//
// Methods which search the knowledge space for solutions.
//
// Cleve Lendon

import (
    "strings"
    "fmt"
)

// Solve - finds one solution for the given goal.
// The solution is returned as a string.
// A second string indicates the reason for failure, as follows:
//    "" (success)
//    "No" (no solution)
//    "Other reason"
// Note: This method only finds the first result.
// See SolveAll below.
// Params:  goal
//          knowledgebase
//          substitution set (previous bindings)
// Returns: solution
//          reason for failure
func Solve(goal Complex, kb KnowledgeBase, ss SubstitutionSet) (solution Complex, failure string) {

    defer func() {  // Catch panics.
        if r := recover(); r != nil {
            solution = goal
            failure = fmt.Sprintf("%v", r)
        }
    }()

    //SetMaxTimeMilliseconds(100)
    SetStartTime()

    timer := MakeTimer()  // For execution time-out.

    solutionChannel := make(chan Complex)

    // Get the root solution node.
    root := goal.GetSolver(kb, ss, nil)

    // Get the next solution.
    go func(out chan<- Complex) {
        newSS, found := root.NextSolution()
        if found {
            out <- goal.ReplaceVariables(newSS).(Complex)
        } else {
            out <- nil
        }
    }(solutionChannel)

    select {
    case <-timer.C:
        solution = nil
        failure = "Time out."
        suironHasTimedOut = true  // Stop searching for a solution.
    case solution = <-solutionChannel:
        timer.Stop()
        if solution == nil {
            failure = "No"
        } else {
            failure = ""
        }
    }

    return solution, failure

}  // Solve


// SolveAll - finds all solutions for the given goal.
// The solutions are returned as a list (slice) of strings.
// A second return value indicates the reason for failure,
// as follows:
//    "" (success)
//    "No" (no solution)
//    "Other reason"
// Params:  goal
//          knowledgebase
//          substitution set (previous bindings)
// Returns: solutions, failure
//
func SolveAll(goal Complex, kb KnowledgeBase, ss SubstitutionSet) (solutions []Complex, failure string) {

    var newSS  SubstitutionSet
    var found  bool

    defer func() {  // Catch panics.
        if r := recover(); r != nil {
            failure  = fmt.Sprintf("%v", r)
        }
    }()

    //SetMaxTimeMilliseconds(100)
    SetStartTime()
    timer := MakeTimer()  // For execution time-out.

    solutionChannel := make(chan []Complex)

    // Get the next solution.
    go func(out chan<- []Complex) {

        // Get the root solution node.
        root := goal.GetSolver(kb, ss, nil)

        // Get the next solution.
        newSS, found = root.NextSolution()

        for found {
            // Replace variables with their bound constants.
            result := goal.ReplaceVariables(newSS)
            solutions = append(solutions, result.(Complex))
            newSS, found = root.NextSolution()
        }
        out <- solutions

    }(solutionChannel)

    select {
    case <-timer.C:
        failure = "Time out."
        suironHasTimedOut = true  // Stop searching for a solution.
    case solutions = <-solutionChannel:
        timer.Stop()
        if len(solutions) == 0 {
            failure = "No"
        } else {
            failure = ""
        }
    }

    return solutions, failure

}  // SolveAll

// FormatSolution - formats a string to display the variable bindings
// of a solution. For example, if the query were: grandfather(Godwin, $X),
// then the function would return: $X = Harold
// Params:
//    query
//    bindings (substitution set)
// Return:
//    solution in string format
func FormatSolution(query Complex, bindings SubstitutionSet) string {
    if bindings != nil {
        var sb strings.Builder
        result := query.ReplaceVariables(bindings).(Complex)
        first := true  // first variable
        for n, term := range query {
            tt := term.TermType()
            if tt == VARIABLE {
                v := term.(VariableStruct)
                if !first { sb.WriteString(", ") }
                s := fmt.Sprintf("%v = %v", v.name, result[n])
                sb.WriteString(s)
                first = false
            }
        }
        return sb.String()
    }
    return "No"
} // FormatSolution
