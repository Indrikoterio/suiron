package suiron

// Complex (term)
//
// In Suiron, as in Prolog, a complex term consists of a functor
// (which must be an atom), and a list of arguments, which are atoms,
// numbers, or variables.
//
// Format:   functor(argument1, argument2, ...)
// Examples: owns(john, house),
//           owns($X, house)
//
// Note: 'Complex terms' are also called 'compound terms'.
//
// Also note: In Suiron, unlike in Prolog, a variable is defined
// with a dollar sign and a letter, eg.: owner($X, house)
// This allows for upper case atoms: owns(John, house).
//
// A complex term is implemented as a slice of unifiable terms.
// It can be instantiated as:
//     Complex{owns, john, house}
// Notice that the above has curley brackets {}, not round ones ().
// The first argument ('owns'), which is the functor, must be an
// Atom, or strange things might happen.
// Complex terms can also be created with the function ParseComplex():
//     ParseComplex("owns(John, house)")
//
// Cleve Lendon

import (
    "strings"
    "fmt"
)

type Complex []Unifiable

// ParseComplex - parses a string to produce a complex term.
// Example of usage:
//     c := ParseComplex("symptom(covid, fever)")
// Important: Backslash is used to escape characters, such as the comma.
// For example:
//     c := ParseComplex("punctuation(comma, \\,)")
// The backslash is doubled, because the Go compiler also interprets
// the backslash.
func ParseComplex(str string) (Complex, error) {

    s := strings.TrimSpace(str)
    length := len(s)

    if length == 0 {
        e := complexError("Length of string is 0", s)
        return Complex{}, e
    }

    if length > 1000 {
        e := complexError("String is too long", s)
        return Complex{}, e
    }

    first :=  s[0:1]
    if first == "$" || first == "(" {
        e := complexError("First character is invalid", s)
        return Complex{}, e
    }

    // Get indices.
    left, right, err := indicesOfParentheses([]rune(s))
    if err != nil { return Complex{}, err }

    if left == -1 { // If left is -1, right must also be -1.
        return parseFunctorTerms(s, "")
    }

    functor := strings.TrimSpace(s[0: left])
    args    := strings.TrimSpace(s[left + 1: right])
    return parseFunctorTerms(functor, args)

} // ParseComplex


// parseFunctorTerms - produces a complex term from two string arguments,
// the functor and a list of terms. For example:
//
//     c, err := parseFunctorTerms("father", "Anakin, Luke")
// produces (c =)
//     father(Anakin, Luke)
//
// Params: functor (string)
//         list of terms (string)
// Return: complex term
//         error
//
func parseFunctorTerms(functor string, terms string) (Complex, error) {
    f := Atom(functor)
    if terms == "" { return Complex{f}, nil }
    t, err := parseArguments(terms)
    if err != nil { return Complex{f}, err }
    unifiables := append([]Unifiable{f}, t...)
    return Complex(unifiables), nil
} // parseFunctorTerms

// complexError - creates an error for Complex terms.
// msg - error message
// str - string which caused the error
func complexError(msg string, str string) error {
    return fmt.Errorf("ParseComplex - %v: >%v<\n", msg, str)
}

// Arity - Returns the arity of a complex term.
// address(Tokyo, Shinjuku, Takadanobaba) has an arity of 3.
func (c Complex) Arity() int { return len(c) - 1 }

// Key - Creates a key (functor/arity).
// Eg. loves(Chandler, Monica) --> loves/2
func (c Complex) Key() string {
    return fmt.Sprintf("%v/%d", c[0].(Unifiable), len(c) - 1)
}

// GetFunctor - The functor is the first term: [functor, term1, term2, term3]
func (c Complex) GetFunctor() Atom { return c[0].(Atom) }

// GetTerm - Returns the indexed term. Term 0 is the functor.
// No error checking.
func (c Complex) GetTerm(index int) Unifiable { return c[index].(Unifiable) }

// TermType - Returns an integer constant which identifies the type.
func (c Complex) TermType() int { return COMPLEX }

// Unify - Unifies this complex term with another unifiable expression.
// Please refer to unifiable.go.
func (c Complex) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {

    otherType := other.TermType()
    if otherType == COMPLEX {
        cOther := other.(Complex)
        lenOther := len(cOther)
        if len(c) != lenOther { return ss, false }
        newSS := ss
        var ok bool
        for i := 0; i < lenOther; i++ {
            termA := c[i]
            termB := cOther[i]
            if termA.TermType() == ANONYMOUS { continue }
            if termB.TermType() == ANONYMOUS { continue }
            newSS, ok = termA.Unify(termB, newSS)
            if !ok { return ss, false }
        }
        return newSS, true
    }

    if otherType == VARIABLE { return other.Unify(c, ss) }
    if otherType == ANONYMOUS { return ss, true }

    return ss, false
} // Unify

// GetSolver - returns a solution node for Complex terms.
func (c Complex) GetSolver(kb KnowledgeBase,
                           parentSolution SubstitutionSet,
                           parentNode SolutionNode) SolutionNode {
    return MakeComplexSolutionNode(c, kb, parentSolution, parentNode)
}


//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - In Prolog, and in this inference engine, the scope of
// a logic variable is the rule or goal in which it is defined. When the
// algorithm tries to solve a goal, it calls this method to ensure that the
// variables are unique. See comments in expression.go.
func (c Complex) RecreateVariables(vars VarMap) Expression {
    newTerms := []Unifiable{}
    for i := 0; i < len(c); i++ {
        term := []Unifiable(c)[i]
        newTerms = append(newTerms, term.RecreateVariables(vars).(Unifiable))
        // Only variables are affected.
    }
    return Complex(newTerms)
}

// ReplaceVariables - replaces a bound variable with its binding.
// This method is used for displaying final results.
// Refer to comments in expression.go.
func (c Complex) ReplaceVariables(ss SubstitutionSet) Expression {
    newTerms := []Unifiable{}
    for i := 0; i < len(c); i++ {
        newTerms = append(newTerms, c[i].ReplaceVariables(ss).(Unifiable))
    }
    return Complex(newTerms)
}

// String - returns a string representation of this complex term.
// For example:  "owns(John, house)"
func (c Complex) String() string {
    length := len(c)
    functor := c[0].String()
    if length == 1 { return functor }
    var sb strings.Builder
    sb.WriteString(functor)
    sb.WriteString("(")
    for i := 1; i < length; i++ {
        if i != 1 { sb.WriteString(", ") }
        sb.WriteString(c[i].String())
    }
    sb.WriteString(")")
    return sb.String()
}
