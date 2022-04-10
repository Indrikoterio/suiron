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
    //"fmt"
)

type SubstitutionSet map[VariableStruct]Unifiable

// Copy() - Makes a copy of the substitution set.
func (ss SubstitutionSet) Copy() SubstitutionSet {
    newSS := SubstitutionSet{}
    for k, v := range ss {
        newSS[k] = v
    }
    return newSS
}

// IsBound() - A logic variable is bound if there exists an entry
// for it in the substitution set.
func (ss SubstitutionSet) IsBound(v VariableStruct) bool {
    _, found := ss[v]
    return found
}


// GetBinding() - Returns the binding of a logic variable.
// If there is no binding, return an error.
func (ss SubstitutionSet) GetBinding(v VariableStruct) (Unifiable, error) {
    unifiableTerm, found := ss[v]
    if found { return unifiableTerm, nil }
    return unifiableTerm, errors.New("Not bound: " + v.String())
}


// IsGroundVariable - A variable is 'ground' if it is ultimately
// bound to something other than a variable.
func (ss SubstitutionSet) IsGroundVariable(v VariableStruct) bool {
    for {
        if u, ok := ss[v]; ok {
            if u.TermType() != VARIABLE { return true }
            v = u.(VariableStruct)
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
    var u2 Unifiable
    var ok bool
    if u.TermType() != VARIABLE { return u, true }
    for {
        if u2, ok = ss[u.(VariableStruct)]; ok {
            if u2.TermType() != VARIABLE { return u2, true }
        } else { return u, false }
        u = u2
    }
    return u2, false
} // GetGroundTerm


// String - creates a string representation of the substitution set,
// for debugging purposes.
func (ss SubstitutionSet) String() string {
    var sb strings.Builder
    sb.WriteString("\n----- Bindings -----\n")
    keys      := make([]string, 0, len(ss))
    keyValues := map[string]string{}
    for k := range ss {
        strK := k.String()
        keys = append(keys, strK)
        keyValues[strK] = ss[k].String()
    }
    sort.Strings(keys)
    for _, k := range keys {
        sb.WriteString("    " + k + ": " + keyValues[k] + "\n")
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
