package suiron

// Goal - is the base interface of all goal objects. Complex terms
// and operators such as And and Or etc. implement this interface.
//
// The method GetSolver() provides a solution node.
//
// Cleve Lendon

type Goal interface {

    Expression

    // GetSolver - gets a solution node for the current goal.
    GetSolver(kb KnowledgeBase, parentSolution SubstitutionSet,
              parentNode SolutionNode) SolutionNode
}

func MakeGoal(terms ...Unifiable) Complex {

    // The main bottleneck in Suiron is the time it takes
    // to copy the substitution set. The substitution set
    // is as large as the highest variable ID. Therefore
    // variableId should be set to 0 for every query.
    variableId = 0

    newTerms := makeLogicVariablesUnique(terms...)
    return Complex(newTerms)

} // MakeGoal

// ParseGoal - creates a goal (Complex term) from a text string.
// Ensures that all logic variables have a unique ID.
// Params: string representation of goal
// Return: complex term (slice of Unifiables)
//         error
func ParseGoal(str string) (Complex, error) {

    // The main bottleneck in Suiron is the time it takes
    // to copy the substitution set. The substitution set
    // is as large as the highest variable ID. Therefore
    // variableId should be set to 0 for every query.
    variableId = 0

    c, err := ParseComplex(str)
    if err != nil { return c, err }
    terms := []Unifiable(c) // get terms
    newTerms := makeLogicVariablesUnique(terms...)
    return Complex(newTerms), nil
} // MakeGoal

// makeLogicVariablesUnique - Long explanation.
// A substitution set keeps track of the bindings of logic variables.
// In order to avoid the overhead of hashing, the substitution set is
// indexed by the ID numbers of these variables. If two logic vars had
// the same ID, this would cause the search for a solution to fail.
// The function LogicVar() creates logic variables with a name and an
// ID number, which is always 0. This is OK, because whenever a rule
// is fetched from the knowledge base, it's variables are recreated,
// by calling RecreateVariables().
// However, goals are not fetched from the knowledge base. If a goal
// is created, it is necessary to ensure that any logic variables it
// contains do not have an index of 0.
func makeLogicVariablesUnique(terms ...Unifiable) []Unifiable {
    newTerms := []Unifiable{}
    vars := make(VarMap)
    for _, term := range terms {
        newTerm := term.RecreateVariables(vars).(Unifiable)
        newTerms = append(newTerms, newTerm)
    }
    return newTerms
} // makeLogicVariablesUnique
