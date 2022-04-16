package suiron

// Miscellaneous
// Cleve Lendon

const (

    NONE = iota

//----------term types----------
    ATOM
    INTEGER
    FLOAT
    VARIABLE
    COMPLEX     // prof(...)
    LINKEDLIST  // [...]
    ANONYMOUS   // anonymous (don't care) variable, $_
    FUNCTION

//--------parsing tokens--------
    SUBGOAL
    COMMA
    SEMICOLON
    LPAREN
    RPAREN
    GROUP
    AND
    OR

//-----------INFIXES-----------
    UNIFY          // =    Unify does unification.
    EQUAL          // ==   No unification. Simply compares.
    GREATER_THAN
    LESS_THAN
    GREATER_THAN_OR_EQUAL
    LESS_THAN_OR_EQUAL
)

var suironConstString = [...]string{ "NONE", "ATOM", "INTEGER",
    "FLOAT", "VARIABLE", "COMPLEX", "LINKEDLIST", "ANONYMOUS",
    "FUNCTION", "SUBGOAL", "COMMA", "SEMICOLON", "LPAREN", "RPAREN",
    "GROUP", "AND", "OR", "UNIFY", "EQUAL", "GREATER_THAN",
    "LESS_THAN", "GREATER_THAN_OR_EQUAL", "LESS_THAN_OR_EQUAL" }

func srConstToString(c int) string {
    if c < 0 || c >= len(suironConstString) { return "" }
    return suironConstString[c]
}
