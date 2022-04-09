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
)

var suironConstString = [...]string{ "NONE", "ATOM", "INTEGER",
    "FLOAT", "VARIABLE", "COMPLEX", "LINKEDLIST", "ANONYMOUS",
    "FUNCTION", "SUBGOAL", "COMMA", "SEMICOLON", "LPAREN",
    "RPAREN", "GROUP", "AND", "OR" }

func srConstToString(c int) string {
    if c < 0 || c >= len(suironConstString) { return "" }
    return suironConstString[c]
}
