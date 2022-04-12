package main

// parse_demo
//
// This is a demo program which parses simple English sentences, and checks
// for grammatical errors. It is not intended to be complete or practical.
//
// In order to understand the comments below, it is necessary to have
// a basic understanding of logic programming. Here are some references:
//
// http://athena.ecs.csus.edu/~mei/logicp/prolog.html
// https://courses.cs.washington.edu/courses/cse341/12au/prolog/basics.html
//
// The program starts at main. First it...
//   - creates an empty knowledge base
//   - loads part-of-speech data into a map (createPoSMap)
//   - creates some atoms, logic variables, and rules
//   - stores these rules into the knowledge base
//   - loads additional facts and rules from a file (LoadKBFromFile)
//   - read in a text file, sentences.txt
//   - splits the text into sentences (splitIntoSentences)
//
// Next, the program calls sentenceToFacts(). This function does several
// things. It calls sentenceToWords() to divide each sentence into words
// and punctuation. For example:
//
//    "They envy us."
//
// becomes...
//
//    ["They", "envy", "us", "."]
//
// Next it creates a linked list by calling MakeLinkedList():
//
//    [They, envy, us, .]
//
// Note: In Prolog, all words which begin with a capital letter are
// variables. In Suiron, variables begin with a dollar sign, eg. $X.
// A capitalized word, such as 'They', is an atom.
//
// sentenceToFacts() calls the function makeFacts(). This function
// makes facts which associate each word with a grammatical fact.
// For example:
//
//    word(we, pronoun(we , subject, first, plural))
//
// Note: Many words can have more than one part of speech. The word
// 'envy', for example, can be a noun or a verb. In order to parse
// English sentences, the program needs facts which identify all
// possible parts of speech:
//
//     word(envy, noun(envy, singular)).
//     word(envy, verb(envy, present, base)).
//
// Finally, the program calls the method Solve(), which tries to find
// a solution for the goal 'parse'.
//
// The arguments of Solve() are:
//     the goal - parse([They, envy, us, .], $X)
//     knowledge base
//     an empty substitution set
//
// During analysis, the rule words_to_pos/2 is applied to convert
// the input word list, created by sentenceToFacts(), into a list
// of terms which identify part of speech.
//
//   words_to_pos([$H1 | $T1], [$H2 | $T2]) :-
//                          word($H1, $H2), words_to_pos($T1, $T2).
//   words_to_pos([], []).
//
// The sentence "They envy us." will become:
//
// [pronoun(They, subject, third, plural), verb(envy, present, base),
//          pronoun(us, object, first, plural), period(.)]
//
// The inference rule 'sentence' identifies (unifies with) various
// types of sentence, such as:
//
//   subject pronoun, verb
//   subject noun, verb
//   subject pronoun, verb, object
//   subject noun, verb, object
//
// There are rules to check subject/verb agreement of these sentences:
//
//    check_pron_verb
//    check_noun_verb
//
// When a mismatch is found (*He envy), these rules print out an error
// message:
//
// 'He' and 'envy' do not agree.
//
// Cleve (Klivo) Lendon
//

import (
    . "github.com/indrikoterio/suiron/suiron"
    "strings"
    "bufio"
    "fmt"
    "os"
)

// The demo program starts here.
func main() {

    // The knowledge base stores rules and facts.
    kb := KnowledgeBase{}

    // Load part of speech data from a text file.
    pos, err := createPoSMap("part_of_speech.txt")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // -------------------------------
    parse  := Atom("parse")
    words_to_pos := Atom("words_to_pos")
    word         := Atom("word")

    // Define variables.
    H1, _ := LogicVar("$H1")
    H2, _ := LogicVar("$H2")
    T1, _ := LogicVar("$T1")
    T2, _ := LogicVar("$T2")
    X, _   := LogicVar("$X")

    /*
     words_to_pos/2 is a rule to convert a list of words into a list
     of parts of speech. For example, the atom 'the' is converted to
     the Complex term 'article(the, definite)':

         words_to_pos([$H1 | $T1], [$H2 | $T2]) :- word($H1, $H2),
                                                   words_to_pos($T1, $T2).
         words_to_pos([], []).

     */

    head := Complex{words_to_pos, MakeLinkedList(true, H1, T1),
                                  MakeLinkedList(true, H2, T2)}
    body := And(Complex{word, H1, H2}, Complex{words_to_pos, T1, T2})
    rule := Rule(head, body)

    // Note: The Atom, MakeVar and Rule definitions above can be
    // replaced by a single line:
    // rule = ParseRule("words_to_pos([$H1 | $T1], [$H2 | $T2]) :- " +
    //                  "word($H1, $H2), words_to_pos($T1, $T2)")
    // ParseRule will parse the given string to produce a RuleStruct.
    // In Prolog, variables begin with a capital letter and atoms
    // begin with a lower case letter. Suiron is a little different.
    // The parser requires a dollar sign to identify variables.
    // An atom can begin with an upper case or lower case letter.

    kb.Add(rule)  // Add the rule to our knowledge base.

    rule = Fact(Complex{words_to_pos, LinkedListStruct{}, LinkedListStruct{}})
    // Alternative (simpler) rule definition:
    //rule, _ = ParseRule("words_to_pos([], [])")

    kb.Add(rule)

    // Rules for noun phrases.
    rule, _ = ParseRule("make_np([adjective($Adj, $_), " +
                        "noun($Noun, $Plur) | $T], [$NP | $Out]) :- " +
                        "!, $NP = np([$Adj, $Noun], $Plur), make_np($T, $Out)")
    kb.Add(rule)
    rule, _ = ParseRule("make_np([$H | $T], [$H | $T2]) :- make_np($T, $T2)")
    kb.Add(rule)
    rule, _ = ParseRule("make_np([], [])")
    kb.Add(rule)

    // Read facts and rules from file.
    fn := "demo_grammar.txt"
    err = LoadKBFromFile(kb, fn)
    if err != nil { 
        fmt.Println(err.Error())
        return
    }

    text, err := readFile("sentences.txt")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    sentences := splitIntoSentences(text)
    for _, sentence := range sentences {

        fmt.Print(sentence)

        // Delete previous 'word' facts. Don't want them to accumulate.
        delete(kb, "word/2")

        inList := sentenceToFacts(sentence, kb, pos)
        //DBKB(kb)

        goal := MakeGoal(parse, inList, X)

        _, failure := Solve(goal, kb, SubstitutionSet{})
        if len(failure) != 0 { fmt.Println(failure) }
        fmt.Print("\n")
    }
} // main


// readFile - reads a file into a string.
// Params: file name
// Return: file contents as string
func readFile(fileName string) (string, error) {

    file, err := os.Open(fileName)
    if err != nil { return "", err }
    defer file.Close()

    var sb strings.Builder
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        aLine := scanner.Text()
        sb.WriteString(aLine)
        sb.WriteString(" ")
    }
    return sb.String(), nil

} // readFile


// isPunc - returns true if the character is a punctuation mark,
// possibly marking the end of a sentence.
//
// Params: character to test
// Return: true/false
func isPunc(c rune) bool {
    if c == '!' || c == '?' || c == '.' { return true }
    return false
} // isPunc

// endOfWord - returns true if the current character is a space,
// or is at the end of a line.
//
// Params: character to test
// Return: true/false
func endOfWord(c rune) bool {
    if c == ' ' || c == '\n' { return true }
    return false
}  // endOfWord


// splitIntoSentences - splits a string of text into sentences, by searching
// for punctuation. The punctuation must be followed by a space.
// (The period in '3.14' doesn't mark the end of a sentence.)
//
// Params: input string
// Return: list of sentences
//
func splitIntoSentences(str string) []string {

    sentences := []string{}
    var sentence string
    previousIndex := 0
    previous3 := []rune{'a', 'a', 'a'}

    for i, c := range str {

        if endOfWord(c) && isPunc(previous3[2]) {
            if previous3[2] == '.' {
               // Check for H.G. Wells or H. G. Wells
               if previous3[0] != '.' && previous3[0] != ' ' {
                   sentence = strings.TrimSpace(str[previousIndex: i])
                   sentences = append(sentences, sentence)
                   previousIndex = i
               }
            } else {
                sentence = strings.TrimSpace(str[previousIndex: i])
                sentences = append(sentences, sentence)
                previousIndex = i
            }
        }
        previous3[0] = previous3[1]
        previous3[1] = previous3[2]
        previous3[2] = c

    } // for

    length := len(str)

    s := strings.TrimSpace(str[previousIndex: length])
    if len(s) > 0 {
        sentences = append(sentences, s)
    }

    return sentences

} // splitIntoSentences


// sentenceToFacts - divides a sentence it into words, and creates
// facts which are written to the knowledge base.
//
// Params: sentence
//         knowledge base
//         part of speech map
// Return: word list (linked list)
//
func sentenceToFacts(sentence string, kb KnowledgeBase,
                     pos map[string][]string) LinkedListStruct {

    words := sentenceToWords(sentence)

    terms := []Unifiable{}
    for _, word := range words {
        terms = append(terms, Atom(word))
    }

    wordList := MakeLinkedList(false, terms...)

    // Make word facts, such as: word(envy, noun(envy, singular)).
    facts := makeFacts(words, pos)
    for _, fact := range facts {
        kb.Add(fact)
    }

    return wordList

} // sentenceToFacts
