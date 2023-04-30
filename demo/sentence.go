package main

// sentence - functions which divide an English sentence into
// a list of words and punctuation.
//
// Cleve (Klivo) Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "strings"
    //"fmt"
)

const MAX_WORDS_IN_SENTENCE int = 120

// isAnApostrophe - tests whether a character is an apostrophe.
//
// Params: character (rune)
// Return: true/false
func isAnApostrophe(ch rune) bool {
    if ch == '\'' || ch == '\u02bc' { return true }
    return false
}

// isPunctuation - determines whether a character is punctuation.
// EXCEPT if the character is a period (.).
// A period could be part of an abbreviation or number (eg. 37.49).
//
// Params: character (rune)
// Return: true/false
func isPunctuation(ch rune) bool {
    if ch == '.' { return false }
    if ch >= '!' && ch <= '/' { return true }
    if ch >= ':' && ch <= '@' { return true }
    if ch == '\u2013'         { return true }  // en-dash
    if isQuoteMark(ch) >= 0   { return true }
    return false
}

// isQuoteMark - checks whether the character is a quote mark ("'«).
//
// If yes, return the index of the quote mark.
// If no, return -1.
//
// Params:  character (rune)
// Return:  index of quote (or -1)
func isQuoteMark(ch rune) int {
    for i := 0; i < len(leftQuotes); i++ {
        if ch == leftQuotes[i] { return i }
        if ch == rightQuotes[i] { return i }
    }
    return -1
} // isQuoteMark

var leftQuotes  = []rune{'\'', '"', '\u00ab', '\u2018', '\u201c' }
var rightQuotes = []rune{'\'', '"', '\u00bb', '\u2019', '\u201d' }


// endOfSentence - determines whether a period is at the end of a sentence.
// (If it is at the end, it must be punctuation.)
//
// Params: sentence (runes)
//         index
// Return: true/false
//
func endOfSentence(sentence []rune, index int) bool {
    length := len(sentence)
    if index >= length - 1 { return true }
    for index < length {
        ch := sentence[index]
        index++
        if LetterNumberHyphen(ch) { return false }
    }
    return true
} // endOfSentence


// getWords - divides a sentence into a list of words and punctuation.
//
// Params: sentence string
// Return: words and punctuation (slice of strings)
//
func getWords(sentence string) []string {

    words := []string{}
    numberOfWords := 0

    r := []rune(sentence)
    length := len(r)

    startIndex := 0
    var lastIndex int

    for startIndex < length && numberOfWords < MAX_WORDS_IN_SENTENCE {

        character := ' '

        // Skip spaces, etc.
        for startIndex < length {
            character = r[startIndex]
            if character > ' ' { break }
            startIndex++
        }
        if startIndex >= length { break }

        // A period at the end of a sentence is punctuation.
        // A period in the middle is probably part of an abbreviation
        // or number, eg.: 7.3
        if character == '.' && endOfSentence(r, startIndex) {
            words = append(words, ".")
            startIndex++
        } else if (isPunctuation(character)) {
            words = append(words, string(character))
            startIndex++
        } else if (LetterNumberHyphen(character)) {

            for lastIndex = startIndex + 1; lastIndex < length; lastIndex++ {
                character = r[lastIndex]
                if character == '.' {
                    if endOfSentence(r, lastIndex) { break }
                    // There might be an apostrophe within the word: don't, we've
                } else if isAnApostrophe(character) {
                    if lastIndex < length - 1 {
                        ch2 := r[lastIndex + 1]
                        if !LetterNumberHyphen(ch2) { break }
                    }
                } else {
                    if !LetterNumberHyphen(character) { break }
                }
            } // for

            word := r[startIndex: lastIndex]
            words = append(words, string(word))

            numberOfWords++

            startIndex = lastIndex

        } else {  // unknown character.
            startIndex++
        }
    } // for

    return words

} // getWords

// sentenceToWords - divides a sentence into words.
//
// Params: original sentence
// Return: list of words
//
func sentenceToWords(sentence string) []string {
    // Clean up the string. New line becomes a space.
    s := strings.ReplaceAll(sentence, "\n", " ")
    // Divide string into words and punctuation.
    return getWords(s)
}
