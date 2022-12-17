package suiron

// Suiron supports singly linked lists (called LinkedList), which are
// represented as a list of items between square brackets, as in Prolog:
//
//    []             // empty list
//    [a, b, c, d]
//    [a, b | $Z]    // vertical bar separates head and tail
//
// LinkedLists are implemented as a chain of structs. Each link (node)
// of the list has a value (the term), and a link to the next element
// of the list. The last element of the list is an empty list.
//
// A vertical bar (or pipe) '|', is used to divide the list between
// head terms and the tail, which is everything left over. When two
// lists are unified, a tail variable binds with all the left-over
// tail items in the other list. For example, in the following code,
//
//    [a, b, c, d, e] = [$X, $Y | $Z]
//
// Variable $X binds to a.
// Variable $Y binds to b.
// Variable $Z binds to [c, d, e].
//
// There are two are functions to create a LinkedList from within Go
// source code: MakeLinkedList() or ParseLinkedList().
//
// MakeLinkedList() takes a boolean and a variable number of arguments.
// The boolean indicates that the last argument is a tail variable,
// when true. Examples:
//
//   list := MakeLinkedList(false, a, b, c)    // [a, b, c]
//   list := MakeLinkedList(true, a, b, c, X)  // [a, b, c | $X]
//
// ParseLinkedList() parses a string to produce a LinkedList:
//
//   list := ParseLinkedList("[a, b, c]")
//   list := ParseLinkedList("[a, b, c | $X]")  // X is a tail Variable
//
// Cleve Lendon

import (
    "strings"
    "fmt"
)

type LinkedListStruct struct {
    term    Unifiable
    next    *LinkedListStruct
    count   int  // number of terms in the list
    tailVar bool
}

var emptyList LinkedListStruct = LinkedListStruct{}

// About tailVar:
// It is necessary to distinguish between
//
//     [$A, $B, $X]
// and
//     [$A, $B | $Y]
//
// If the last item in a list is a variable, it can be an ordinary
// list term ($X), or a tail variable ($Y). A tail variable will
// unify with the tail of another list. A non-tail variable will
// not - it can only unify with one term from the other list.
// The tailVar flag is used to indicate that the last variable is
// a tail variable, as in the second case above.

// EmptyList - returns an empty LinkedList.
func EmptyList() LinkedListStruct {
    return emptyList
}

// MakeLinkedList - makes a singly-linked list, such as [a, b, c]
//
//    list := MakeLinkedList(false, a, b, c)
//
// The first parameter, vbar, is set true for lists which have
// a tail variable, such as [a, b, c | $Tail]
//
//    list := MakeLinkedList(true, a, b, c, $Tail)
//
// Params:  vbar - vertical bar flag
//          args - list of unifiable arguments
// Return:  linked list
//
func MakeLinkedList(vbar bool, args ...Unifiable) LinkedListStruct {
    var term Unifiable
    nArgs := len(args)
    if nArgs == 0 { return emptyList }
    tailPtr := &emptyList  // Point to empty list
    tailVar := vbar  // Last variable is a tail variable.
    num := 1
    lastIndex := nArgs - 1
    for i := lastIndex; i > 0; i-- {
        term = args[i]
        if i == lastIndex && term.TermType() == LINKEDLIST {
            t, _ := term.(LinkedListStruct) // I know it's a LinkedList.
            tailPtr = &t
            // If the last term is empty [], there
            // is no need to add it to the tail.
            if tailPtr.term == nil {
               tailPtr = nil
            } else {
               num = tailPtr.count + 1
            }
        } else {
            tailPtr = &LinkedListStruct{ term: term, next: tailPtr,
                                   count: num, tailVar: tailVar }
            num += 1
        }
        tailVar = false
    }
    return LinkedListStruct{ term: args[0], next: tailPtr,
                       count: num, tailVar: tailVar }
}

// equalEscape - compares the indexed character (rune) with the given
// character. If they are equal, the function will return true, except
// if the indexed character is proceeded by an backslash. Because...
// characters which are escaped by a backslash should not be interpreted
// by the parser; they need to be included as they are. Examples:
//
//    text := []rune("OK, sure.")
//    if equalEscape(text, 2, ',')  <-- returns true
//
//    text2 := []rune("OK\\, sure.")  <-- double backslash for Go
//    if equalEscape(text2, 3, ',')  <-- returns false
//
func equalEscape(runes []rune, index int, ch rune) bool {
    r := runes[index]
    if r == ch {
        if index > 0 {
            previous := runes[index - 1]
            if previous == '\\' { return false }
        }
        return true
    }
    return false
}

// ParseLinkedList - parses a string to create a linked list.
// For example,
//     list, err := ParseLinkedList("[a, b, c | $X]")
// Produces an error if the string is invalid.
func ParseLinkedList(str string) (LinkedListStruct, error) {

    s := strings.TrimSpace(str)
    r := []rune(s)
    length := len(r)

    list := LinkedListStruct{}   // Make an empty list.

    if length < 2 {
        err := parseLinkedListError("String is too short", s)
        return list, err
    }

    first := r[0]
    if first != '[' {
        err := parseLinkedListError("Missing opening bracket", s)
        return list, err
    }
    last := r[length - 1]
    if last != ']' {
        err := parseLinkedListError("Missing closing bracket", s)
        return list, err
    }

    if length == 2 { return list, nil }

    arguments := r[1: length-1] // remove brackets
    argLength := len(arguments)

    vbar        := false
    endIndex    := argLength

    openQuote   := false
    numQuotes   := 0
    roundDepth  := 0   // depth of round parentheses (())
    squareDepth := 0   // depth of square brackets [[]]

    var i int
    for i = argLength - 1; i >= 0; i-- {
        if openQuote {
            if equalEscape(arguments, i, '"') {
                openQuote = false
                numQuotes += 1
            }
        } else {
            if equalEscape(arguments, i, ']') {
                squareDepth++
            } else if equalEscape(arguments, i, '[') {
                squareDepth--
            } else if equalEscape(arguments, i, ')') {
                 roundDepth++;
            } else if equalEscape(arguments, i, '(') {
                 roundDepth--
            } else if roundDepth == 0 && squareDepth == 0 {
                if equalEscape(arguments, i, '"') {
                    openQuote = true  // first quote
                    numQuotes += 1
                } else if equalEscape(arguments, i, ',') {
                    strTerm := string(arguments[i + 1: endIndex])
                    strTerm = strings.TrimSpace(strTerm)
                    if len(strTerm) == 0 {
                        err := parseLinkedListError("Missing argument", s)
                        return list, err
                    }
                    error := checkQuotes(strTerm, numQuotes)
                    if error != nil { return list, error }
                    term, err := parseTerm(strTerm)
                    if err == nil {
                        list = linkFront(term, false, list)
                        endIndex = i
                    } else { return emptyList, err }
                    numQuotes = 0
                } else if equalEscape(arguments, i, '|') {  // Must be a tail variable.
                    if vbar {
                        err := parseLinkedListError("Too many vertical bars.", s)
                        return list, err
                    }
                    strTerm := string(arguments[i + 1: endIndex])
                    strTerm = strings.TrimSpace(strTerm)
                    if len(strTerm) == 0 {
                        err := parseLinkedListError("Missing argument", s)
                        return list, err
                    }
                    term, err := LogicVar(strTerm)
                    if err != nil {
                        err := parseLinkedListError("Require variable after vertical bar", s)
                        return list, err
                    }
                    vbar = true
                    list = linkFront(term, true, list);  // tail variable
                    endIndex = i
                }
            }
        } // openQuote

        if i == 0 {
            strTerm := string(arguments[0: endIndex])
            strTerm = strings.TrimSpace(strTerm)
            if len(strTerm) == 0 {
                err := parseLinkedListError("Missing argument", strTerm)
                return list, err
            }
            error := checkQuotes(strTerm, numQuotes)
            if error != nil { return list, error }
            term, err := parseTerm(strTerm)
            if err != nil { return list, err }
            list = linkFront(term, false, list)
            endIndex = i
        }

    } // for

    return list, nil

} // ParseLinkedList()


// linkFront - adds a term to the front of the given linked list (LinkedList).
// Params:  term to add in
//          tail variable flag
//          linked list
// Return:  new linked list
// Note: The tail variable flag is true when the term is tail variable
func linkFront(term Unifiable, tailVar bool, list LinkedListStruct) LinkedListStruct {
    count := list.count + 1
    return LinkedListStruct{ term: term, next: &list, count: count, tailVar: tailVar }
}

// parseLinkedListError - creates an error for ParseLinkedList().
// msg - error message
// str - string which caused the error
func parseLinkedListError(msg string, str string) error {
    return fmt.Errorf("ParseLinkedList() - %v: %v", msg, str)
}

// Flatten - partially flattens this linked list.
// If the number of terms requested is two, this function will return
// a slice of the first and second terms, and the tail of the linked
// list. In other words, the list [a, b, c, d] becomes a slice containing
// a, b, and the linked list [c, d]. The function returns the resulting
// slice and a boolean to indicate success or failure.
func (ll LinkedListStruct) Flatten(numOfTerms int, ss SubstitutionSet) ([]Unifiable, bool) {
    ptr := &ll
    outList := []Unifiable{}
    if numOfTerms < 1 { return outList, false }
    for i := 0; i < numOfTerms; i++ {
        if ptr == nil { return outList, false } // fail
        term := ptr.term
        if ptr.tailVar {  // Is this node a tail variable?
            varTerm, _ := term.(VariableStruct)
            list, ok := ss.CastLinkedList(varTerm)
            if ok {
                ptr = &list
                term = ptr.term
            }
        }
        if term == nil {
            return outList, false
        }
        outList = append(outList, term)
        ptr = ptr.next
    }
    if ptr != nil {
        outList = append(outList, *ptr)
    }
    return outList, true
} // Flatten

// Unify - unifies this LinkedList with a Variable or another LinkedList.
// Two lists can unify if they have the same number of items, and each
// corresponding pair of items can unify. Or, if one of the lists ends in
// a tail Variable (eg. [a, b, | $X]), the tail Variable can unify with
// the remainder of the other list.
// The method returns an updated substitution set and a boolean flag which
// indicates success or failure. Please refer to unifiable.go.
func (ll LinkedListStruct) Unify(other Unifiable,
                           ss SubstitutionSet) (SubstitutionSet, bool) {

    var ok bool

    if other.TermType() == LINKEDLIST {

        o, _ := other.(LinkedListStruct) // I know that this is a LinkedList.
        newSS := ss

        // Empty lists unify. [] = []
        if ll.term == nil && o.term == nil {
            return ss, true
        }

        thisList  := &ll
        otherList := &o

        var thisTerm Unifiable
        var thisIsTailVar bool

        var otherTerm Unifiable
        var otherIsTailVar bool

        for thisList != nil && otherList != nil {

            thisTerm  = thisList.term
            otherTerm = otherList.term
            thisIsTailVar  = thisList.tailVar
            otherIsTailVar = otherList.tailVar

            if thisIsTailVar && otherIsTailVar {
                if otherTerm.TermType() == ANONYMOUS { return newSS, true }
                if thisTerm.TermType()  == ANONYMOUS { return newSS, true }
                return thisTerm.Unify(otherTerm, newSS)
            } else if thisIsTailVar {
                return thisTerm.Unify(*otherList, newSS)
            } else if otherIsTailVar {
                return otherTerm.Unify(*thisList, newSS)
            } else {
               if thisTerm == nil && otherTerm == nil {
                   return newSS, true
               }
               if thisTerm == nil || otherTerm == nil { return ss, false }
               newSS, ok = thisTerm.Unify(otherTerm, newSS)
               if !ok { return ss, false }
            }
            thisList  = thisList.next
            otherList = otherList.next
        } // for

        return ss, false

    } else if other.TermType() == VARIABLE {
        return other.Unify(ll, ss)
    }

    return ss, false  // failure

}  // Unify

// TermType - Returns an integer constant which identifies the type.
// This function satisfies the Unifiable interface.
func (ll LinkedListStruct) TermType() int { return LINKEDLIST }

// GetTerm - returns the top term of this list.
func (ll LinkedListStruct) GetTerm() Unifiable { return ll.term }

// GetNext - returns pointer to the rest of the list.
func (ll LinkedListStruct) GetNext() *LinkedListStruct { return ll.next }

// GetCount - returns the number of items in the list.
func (ll LinkedListStruct) GetCount() int { return ll.count }

//----------------------------------------------------------------
// RecreateVariables(), ReplaceVariables(), and String() satisfy
// the Expression interface.
//----------------------------------------------------------------

// RecreateVariables - The scope of a logic variable is the rule or goal in
// which it is defined. When the algorithm tries to solve a goal, it calls
// this method to ensure that the variables are unique.
// See comments in expression.go.
func (ll LinkedListStruct) RecreateVariables(vars VarMap) Expression {
    newTerms := []Unifiable{}
    thisList := &ll
    vbar  := thisList.tailVar
    term     := thisList.term
    if term == nil { return emptyList }
    for term != nil {
        newTerms = append(newTerms, term.RecreateVariables(vars).(Unifiable))
        vbar = thisList.tailVar
        thisList = thisList.next
        if (thisList == nil) { break }
        term = thisList.term
    }
    newLinkedList := MakeLinkedList(vbar, newTerms...)
    return newLinkedList
} // RecreateVariables()


// ReplaceVariables - replaces a bound variable with its binding.
// This method is used for displaying final results.
// Refer to comments in expression.go.
func (ll LinkedListStruct) ReplaceVariables(ss SubstitutionSet) Expression {
    newTerms := []Unifiable{}
    var newTerm Unifiable
    thisList := &ll
    term     := thisList.term
    if term == nil { return ll }
    for term != nil {
        newTerm = term.ReplaceVariables(ss).(Unifiable)
        tt := newTerm.TermType()
        if tt == LINKEDLIST {
            list := newTerm.(LinkedListStruct)
            ptr := &list
            head := ptr.term
            for head != nil {
                newTerms = append(newTerms, head)
                ptr = ptr.next
                if ptr == nil { break }
                head = ptr.term
            }
        } else {  // not a list
           newTerms = append(newTerms, newTerm)
        }
        thisList = thisList.next
        if thisList == nil { break }
        term = thisList.term
    }
    
    result := MakeLinkedList(false, newTerms...)
    return result

} // ReplaceVariables()

// String - returns a string representation of this list.
func (ll LinkedListStruct) String() string {
    ptr := &ll
    if ptr.term == nil { return "[]" }
    var sb strings.Builder
    sb.WriteString("[" + ptr.term.String())
    for ptr.next != nil {
        ptr = ptr.next
        if ptr.term == nil {
            break
        } else if ptr.tailVar {
            sb.WriteString(" | " + ptr.term.String())
        } else {
            sb.WriteString(", " + ptr.term.String())
        }
    }
    sb.WriteString("]")
    return sb.String()
}

/*
    Scan the list from head to tail,
    Curse recursion, force a fail.
    Hold your chin, hypothesize.
    Predicate logic never lies.
*/
