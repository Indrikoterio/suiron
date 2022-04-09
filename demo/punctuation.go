package main

// punctuation - makes logic terms for punctuation symbols: ()?![] etc.
// Cleve (Klivo) Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
)

// makePunctuationTerm - makes a complex term for punctuation.
// Params: symbol (string)
// Return: complex term
//         success/failure flag, true = success
func makePunctuationTerm(symbol string) (Complex, bool) {

    if len(symbol) != 1 { return nil, false } // fail

    var c Complex

    if (symbol == ".") {
        c, _ = ParseComplex("period(.)")
    } else if (symbol == ",") {
        c, _ = ParseComplex("comma(\\,)")  // Must escape the comma with backslash.
    } else if (symbol == "?") {
        c, _ = ParseComplex("question_mark(?)")
    } else if (symbol == "!") {
        c, _ = ParseComplex("exclamation_mark(!)")
    } else if (symbol == ":") {
        c, _ = ParseComplex("colon(:)")
    } else if (symbol == ";") {
        c, _ = ParseComplex("semicolon(;)")
    } else if (symbol == "-") {
        c, _ = ParseComplex("dash(-)")
    } else if (symbol == "\"") {
        // The second argument is for comparisons.
        c, _ = ParseComplex("quote_mark(\", \")")
    } else if (symbol == "'") {
        c, _ = ParseComplex("quote_mark(', ')")
    } else if (symbol == "«") {
        c, _ = ParseComplex("quote_mark(«, «)")
    } else if (symbol == "»") {
        c, _ = ParseComplex("quote_mark(», «)")
    } else if (symbol == "‘") {
        c, _ = ParseComplex("quote_mark(‘, ‘)")
    } else if (symbol == "’") {
        c, _ = ParseComplex("quote_mark(’, ‘)")
    } else if (symbol == "“") {
        c, _ = ParseComplex("quote_mark(“, “)")
    } else if (symbol == "”") {
        c, _ = ParseComplex("quote_mark(”, “)")
    } else if (symbol == "(") {
        c, _ = ParseComplex("bracket((, ()")
    } else if (symbol == ")") {
        c, _ = ParseComplex("bracket(), ()")
    } else if (symbol == "[") {
        c, _ = ParseComplex("bracket([, [)")
    } else if (symbol == "]") {
        c, _ = ParseComplex("bracket(], [)")
    } else if (symbol == "<") {
        c, _ = ParseComplex("bracket(<, <)")
    } else if (symbol == ">") {
        c, _ = ParseComplex("bracket(>, <)")
    }

    if c != nil { return c, true }
    return nil, false

} // makePunctuationTerm
