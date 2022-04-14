package suiron

// SubstitutionSet - is a map of bindings for logic variables.
//
// As the inference engine searches for a solution, it adds variable bindings
// to the substitution set. Because the substitution set contains all variable
// bindings of the search (so far), it can be thought of as the partial or
// complete solution.
//
// Cleve Lendon
//

import (
    "errors"
    "strings"
    "sort"
    "fmt"
)

type SubstitutionSet []*Unifiable

// IsBound() - A logic variable is bound if there exists an entry
// for it in the substitution set.
// Params: logic variable
// Return: true/false
func (ss SubstitutionSet) IsBound(v VariableStruct) bool {
    if v.id >= len(ss) { return false }
    uni := ss[v.id]
    return uni != nil
}


// GetBinding() - Returns the binding of a logic variable.
// If there is no binding, return an error.
// Params: logic variable
// Return: bound term
//         error
func (ss SubstitutionSet) GetBinding(v VariableStruct) (*Unifiable, error) {
    if v.id >= len(ss) {
        return nil, errors.New("Not bound: " + v.String())
    }
    term := ss[v.id]
    if term == nil {
        return nil, errors.New("Not bound: " + v.String())
    }
    return term, nil
}


// IsGroundVariable - A variable is 'ground' if it is ultimately
// bound to something other than a variable.
// Params: logic variable
// Return: true/false
func (ss SubstitutionSet) IsGroundVariable(v VariableStruct) bool {
    for {
        if v.id >= len(ss) { return false }
        u := ss[v.id]
        if u != nil {
            if (*u).TermType() != VARIABLE { return true }
            v = (*u).(VariableStruct)
        } else { return false }
    }
    return false
}

// GetGroundTerm - if the given term is a ground term, return it.
// If it's a variable, try to get its ground term. If the variable
// is bound to a ground term, return the term and set the success
// flag to true. Otherwise, return the variable and set the success
// flag to false.
// Params: term
// Return: ground term
//         success/failure flag
func (ss SubstitutionSet) GetGroundTerm(u Unifiable) (Unifiable, bool) {
    var u2 *Unifiable
    if u.TermType() != VARIABLE { return u, true }
    for {
        id := (u.(VariableStruct)).id
        if id >= len(ss) { return u, false }
        u2 = ss[id]
        if u2 != nil {
            if (*u2).TermType() != VARIABLE { return *u2, true }
        } else { return u, false }
        u = *u2
    }
    return *u2, false
} // GetGroundTerm


// String - creates a string representation of the substitution set,
// for debugging purposes.
func (ss SubstitutionSet) String() string {
    var sb strings.Builder
    sb.WriteString("\n----- Bindings -----\n")
    keys      := make([]int, 0, len(ss))
    keyValues := map[int]string{}
    for k := range ss {
        keys = append(keys, k)
        keyValues[k] = (*ss[k]).String()
    }
    sort.Ints(keys)
    for _, k := range keys {
        str := fmt.Sprintf("    %d: %v\n", k, keyValues[k])
        sb.WriteString(str)
    }
    sb.WriteString("--------------------\n")
    return sb.String()
}

// CastComplex - if the given Unifiable term is a Complex term,
// cast it as Complex and return it. If the given term is a Variable,
// get the ground term. If the ground term is a Complex term, cast
// it and return it. Otherwise fail.
// Params: Unifiable term
// Return: Complex term
//         success/failure flag
func (ss SubstitutionSet) CastComplex(term Unifiable) (Complex, bool) {
    tt := term.TermType()
    if tt == COMPLEX {
        comp, _ := term.(Complex)
        return comp, true
    }
    if tt == VARIABLE {
        varTerm, _ := term.(VariableStruct)
        if outTerm, ok := ss.GetGroundTerm(varTerm); ok {
            if outTerm.TermType() == COMPLEX {
                comp := outTerm.(Complex)
                return comp, true
            }
        }
    }
    return Complex{}, false
} // CastComplex()


// CastLinkedList - if the given Unifiable term is a linked list, cast it
// as a LinkedList and return it. If it is a Variable, get the ground term.
// If that term is a linked list, cast it as a LinkedList and return it.
// Otherwise fail.
//
// Params: Unifiable term
// Return: linked list
//         success/failure flag
//
func (ss SubstitutionSet) CastLinkedList(term Unifiable) (LinkedListStruct, bool) {
    tt := term.TermType()
    if tt == LINKEDLIST {
        list, _ := term.(LinkedListStruct)
        return list, true
    }
    if tt == VARIABLE {
        varTerm, _ := term.(VariableStruct)
        if outTerm, ok := ss.GetGroundTerm(varTerm); ok {
            if outTerm.TermType() == LINKEDLIST {
                list := outTerm.(LinkedListStruct)
                return list, true
            }
        }
    }
    return LinkedListStruct{}, false
} // CastLinkedList()


// CastAtom - if the given Unifiable term is an Atom, cast it
// as an Atom and return it. If it is a Variable, get the ground
// term. If that term is an Atom, cast it as an Atom and return it.
// Otherwise fail.
//
// Params: Unifiable term
// Return: Atom
//         success/failure flag
//
func (ss SubstitutionSet) CastAtom(term Unifiable) (Atom, bool) {
    tt := term.TermType()
    if tt == ATOM {
        at, _ := term.(Atom)
        return at, true
    }
    if tt == VARIABLE {
        varTerm, _ := term.(VariableStruct)
        if outTerm, ok := ss.GetGroundTerm(varTerm); ok {
            if outTerm.TermType() == ATOM {
                at := outTerm.(Atom)
                return at, true
            }
        }
    }
    return Atom(""), false
} // CastAtom()
