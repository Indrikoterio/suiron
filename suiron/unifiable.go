package suiron

// Unifiable - is the base interface for unifiable expressions.
// (Constants, Variables, Complex terms)
//
// The method 'unify' unifies (binds) predicate calculus expressions
// (Constants, Variables, Complex/Compound terms). A substitution set
// records the bindings made while searching for a solution.
//
// For example, if the source program has: $X = water, then the unify
// method for the variable $X will add { $X : 'water' } to the
// substitution set.
//
// The method TermType() returns an integer which identifies the type
// of term. This was defined because reflection is slow.
//
// Cleve Lendon

type Unifiable interface {

    Expression

    // Unify - unifies two terms, if possible. If a variable is
    // unified (bound) to another term, the binding is recorded
    // in the substitution set.
    // 
    // Parameters:
    //     other unifiable term
    //     substitution set
    // Returns:
    //     new substitution set
    //     flag, true = success, false = failure
    Unify(u Unifiable, ss SubstitutionSet) (SubstitutionSet, bool)

    // TermType - returns an integer which defines the type of expression.
    // Eg. ATOM, INTEGER, etc.
    TermType() int
}
