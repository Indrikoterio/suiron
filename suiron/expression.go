package suiron

// Expression - is the base interface of unifiable items (constants,
// variables, complex terms) and goals (complex terms and operators).
//
// Note: In this inference engine, variables begin with a dollar sign
// and a letter. Eg. $X, $Age
// Constants are atoms or numbers. Atoms can start with an upper case
// or a lower case letter. Eg. destiny, Utred.
// Numbers can be integers or floating point numbers: 15, 7.9
//
// Cleve Lendon

// VarMap defines a map which is used for RecreateVariables,
// to keep track of previously recreated variables.
type VarMap map[string]VariableStruct

type Expression interface {

    // RecreateVariables - creates unique variables every time the
    // inference engine fetches a rule from the knowledge base.
    //
    // The scope of a variable is the rule in which it is defined.
    // For example, in the knowledge base we have:
    //
    //    grandparent($X, $Y) :- parent($X, $Z), parent($Z, $Y).
    //    parent($X, $Y) :- father($X, $Y).
    //    parent($X, $Y) :- mother($X, $Y).
    //    mother(Martha, Jackie)
    //    ... other facts and rules
    //
    // To find a solution for the goal 'grandparent(Frank, $X)',
    // the inference engine will fetch the rule parent/2 from the
    // knowledge base. Each time the rule is fetched, the variables
    // $X and $Y must be unique, different from the $X and $Y which
    // were previously fetched.
    //
    // A variable is identified by its print name and a unique id
    // number. For example, the first time the parent/2 rule is
    // fetched, the variable $X might become '$X_22'. The second
    // time it might become '$X_23'.
    //
    RecreateVariables(newVars VarMap) Expression

    // ReplaceVariables() is called after a solution has been found.
    // It replaces logic variables with the constants which they are
    // bound to, in order to display results.
    //
    // For example, consider the goal 'grandfather(Frank, $X)'.
    // If the substitution set (i.e. solution), has $X bound to $Y
    // bound to Cindy, $X would be replaced by Cindy to give:
    //
    //     grandfather(Frank, Cindy).
    //
    ReplaceVariables(ss SubstitutionSet) Expression

    // String - returns a string representation of the expression.
    String() string
}
