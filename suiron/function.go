package suiron

// Function - defines the interface for built-in functions.
// Built-in functions are Expressions. They must satisfy:
//
//    RecreateVariables
//    ReplaceVariables
//    String
//
// Cleve Lendon

type Function interface {
    RecreateVariables(map[Variable]Variable) Expression
    ReplaceVariables(SubstitutionSet) Expression
    String() string
}
