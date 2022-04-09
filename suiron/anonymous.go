package suiron

// This file defines the anonymous variable ($_). The anonymous
// variable unifies with any term. It is used when the term which
// it unifies with is irrelevant.
//
// Cleve Lendon

type Anonymous struct {}

var anon Anonymous = Anonymous{}

// Anon - returns the anonymous variable ($_).
func Anon() Anonymous { return anon }

// String - return this term as a string.
func (a Anonymous) String() string { return "$_" }

// TermType - Returns an integer constant which identifies the type.
func (a Anonymous) TermType() int { return ANONYMOUS }

// Unify - unifies a variable with another unifiable expression.
// For the anonymous variable ($_), no substitution necessary.
func (a Anonymous) Unify(other Unifiable, ss SubstitutionSet) (SubstitutionSet, bool) {
    return ss, true
} // Unify


// RecreateVariables - creates unique variables every time the
// inference engine fetches a rule from the knowledge base. For the
// anonymous variable, RecreateVariables() simply returns itself.
// This function satisfies the Expression interface.
func (a Anonymous) RecreateVariables(vars map[Variable]Variable) Expression {
    return a
}


// ReplaceVariables - replaces a bound variable with its binding.
// For the anonymous variable, ReplaceVariables() simply returns itself.
// This function satisfies the Expression interface.
func (a Anonymous) ReplaceVariables(ss SubstitutionSet) Expression {
    return a
} // ReplaceVariables()
