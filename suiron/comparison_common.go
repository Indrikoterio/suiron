package suiron

// comparison_common - This file contains functions which are common
// to all mathematical comparison functions (<= >= < > == etc.).
//
//     parseComparison()
//     getTermsToCompare()
//     twoFloats()
//     comparisonString()
//
// Cleve Lendon

import (
    "strings"
    "fmt"
)

const errNotGround     = "Cannot compare. Variable %v is not grounded."
const errNotNumber     = "Cannot compare. Not a number: %v"
const errCannotCompare = "Cannot compare. Invalid term type: %v %T"

// getTermsToCompare - gets two terms from the argument array and
// returns their ground terms and types. If a term is not grounded,
// there will be chaos, pandemonium, and a panic.
// Params: array of unifiable terms
//         substitution set
// Return: grounded term1,
//         type of term1,
//         grounded term2
//         type of term2
func getTermsToCompare(terms []Unifiable, ss SubstitutionSet) (Unifiable, int, Unifiable, int) {
    term1 := terms[0]
    term2 := terms[1]
    ground1, ok := ss.GetGroundTerm(term1)
    if !ok {
        msg := fmt.Sprintf(errNotGround, ground1)
        panic(msg)
    }
    ground2, ok := ss.GetGroundTerm(term2)
    if !ok {
        msg := fmt.Sprintf(errNotGround, ground2)
        panic(msg)
    }
    return ground1, ground1.TermType(), ground2, ground2.TermType()
} // getTermsToCompare

// twoFloats
// In order to compare two numbers, they should be the same type,
// both Integers or both Floats. If one number is an Integer, and
// the other is a Float, the Integer should be converted to a Float.
// If one of the terms is not a number, this function will have
// a panic attack.
// Params: term1
//         type of term 1
//         term2
//         type of term 2
// Return: term1 as Float
//         term2 as Float
func twoFloats(term1 Unifiable, type1 int,
               term2 Unifiable, type2 int) (Float, Float) {

    var f1, f2 Float

    if type1 == FLOAT {
        f1 = term1.(Float)
    } else {
        if type1 == INTEGER {
            f1 = Float(term1.(Integer))
        } else {
            msg := fmt.Sprintf(errNotNumber, term1)
            panic(msg)
        }
    }

    if type2 == FLOAT {
        f2 = term2.(Float)
    } else {
        if type2 == INTEGER {
            f2 = Float(term2.(Integer))
        } else {
            msg := fmt.Sprintf(errNotNumber, term2)
            panic(msg)
        }
    }

    return f1, f2

} // twoFloats


// compareAtoms - does a string compare on Atoms. Returns -1 for less
// than, 0 for equal, and 1 for greater than. If one of the terms is an
// Integer or a Float, it must be converted to an Atom for the comparison.
// If one of the terms is not an Atom, Integer or Float, the function
// will cause a panic.
// Params: term1
//         type of term 1
//         term2
//         type of term 2
// Return: result (-1, 0, 1)
func compareAtoms(term1 Unifiable, type1 int,
                  term2 Unifiable, type2 int) int {

    var a1, a2 Atom

    if type1 == ATOM {
        a1 = term1.(Atom)
    } else {
        if type1 == INTEGER {
            a1 = Atom(fmt.Sprintf("%d", term1.(Integer)))
        } else if type1 == FLOAT {
            a1 = Atom(fmt.Sprintf("%f", term1.(Float)))
        } else {
            msg := fmt.Sprintf(errCannotCompare, term1, term1)
            panic(msg)
        }
    }

    if type2 == ATOM {
        a2 = term2.(Atom)
    } else {
        if type2 == INTEGER {
            a2 = Atom(fmt.Sprintf("%d", term2.(Integer)))
        } else if type2 == FLOAT {
            a2 = Atom(fmt.Sprintf("%f", term2.(Float)))
        } else {
            msg := fmt.Sprintf(errCannotCompare, term2, term2)
            panic(msg)
        }
    }

    return strings.Compare(string(a1), string(a2))

} // compareAtoms


// comparisonString - creates a string representation for comparisons,
// for example, "$X <= 5".
// Params: slice of arguments
//         operator (eg. " <= ", " >= ")
// Return: string representation
func comparisonString(terms []Unifiable, operator string) string {
    term1 := terms[0].String()
    term2 := terms[1].String()
    return term1 + operator + term2
} // comparisonString
