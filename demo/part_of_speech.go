package main

// part_of_speech - has functions which read in a list of words
// with part-of-speech tags, and creates a hash map keyed by word.
//
// In addition, there are methods to create Facts which can be
// analyzed by Suiron.
//
// The part-of-speech tags in part_of_speech.txt are from Penn
// State's Treebank tagset. There is a reference here:
//
//   https://sites.google.com/site/partofspeechhelp/home
//
// ABN pre-quantifier (half, all)
// AP post-determiner (many, several, next)
// AT article (a, the, no)
// BE be
// BED were
// BEDZ was
// BEG being
// BEM am
// BEN been
// BER are, art
// BBB is
// CC coordinating conjunction
// CD cardinal digit
// DT determiner
// EX existential there (like: “there is” … think of it like “there exists”)
// FW foreign word
// IN preposition/subordinating conjunction
// JJ adjective 'big'
// JJR adjective, comparative 'bigger'
// JJS adjective, superlative 'biggest'
// LS list marker 1)
// MD modal could, will
// NN noun, singular 'desk'
// NNS noun plural 'desks'
// NNP proper noun, singular 'Harrison'
// NNPS proper noun, plural 'Americans'
// OD ordinal numeral (first, 2nd)
// NPS proper noun, plural Vikings
// PDT predeterminer 'all the kids'
// PN nominal pronoun (everybody, nothing)
// PP$ possessive personal pronoun (my, our)
// PP$$ second (nominal) personal pronoun (mine, ours)
// PPO objective personal pronoun (me, him, it, them)
// PPS 3rd. singular nominative pronoun (he, she, it, one)
// PPSS other nominative personal pronoun (I, we, they, you)
// POS possessive ending parent's
// PRP personal pronoun I, he, she
// PRP$ possessive pronoun my, his, hers
// QL qualifier (very, fairly)
// QLP post-qualifier (enough, indeed)
// RB adverb very, silently,
// RBR adverb, comparative better
// RBS adverb, superlative best
// RP particle give up
// SYM symbol
// TO to go 'to' the store.
// UH interjection errrrrm
// VB verb, base form take
// VBD verb, past tense took
// VBG verb, gerund/present participle taking
// VBN verb, past participle taken
// VBP verb, sing. present, non-3d take
// VBZ verb, 3rd person sing. present takes
// WDT wh-determiner which
// WP wh-pronoun who, what
// WP$ possessive wh-pronoun whose
// WRB wh-abverb where, when
//
// Cleve (Klivo) Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "strings"
    "bufio"
    "fmt"
    "os"
)

const FILENAME = "part_of_speech.txt"

// Word functor. Use capitals to distinguish from the variable 'word'.
const WORD = Atom("word")

const noun = Atom("noun")
const verb = Atom("verb")
const pronoun = Atom("pronoun")
const adjective = Atom("adjective")
const participle = Atom("participle")
const preposition = Atom("preposition")
const unknown = Atom("unknown")

// tenses
const past    = Atom("past")
const present = Atom("present")

// voice
const active  = Atom("active")
const passive = Atom("passive")

// Person, for verbs.
const first_sing  = Atom("first_sing")  // I am
const second_sing = Atom("second_sing") // Thou art
const third_sing  = Atom("third_sing")  // it is
const base        = Atom("base")        // you see

// Person, for pronouns
const first  = Atom("first")  // I, me, we, us
const second = Atom("second") // you
const third  = Atom("third")  // he, him, she, her, it, they, them

// Plurality for nouns and pronouns
const singular = Atom("singular") // table, mouse
const plural   = Atom("plural")   // tables, mice
const both     = Atom("both")     // you

// For adjectives.
const positive    = Atom("positive")    // good
const comparative = Atom("comparative") // better
const superlative = Atom("superlative") // best

// For adverbs.
const adverb = Atom("adverb")  // happily

// For articles.
const article    = Atom("article")    // the, a, an
const definite   = Atom("definite")   // the
const indefinite = Atom("indefinite") // a, an

// For pronouns. (case)
const subject = Atom("subject")  // subject
const object  = Atom("object")   // object

// Punctuation.
const punctuation = Atom("punctuation")

// createPoSMap - reads in part-of-speech data from a file,
// and creates a map of PoS tags, indexed by a word string.
// Params: file name
// Return: PoS map
//         error
func createPoSMap(fileName string) (map[string][]string, error) {

    // Map: word / Part of Speech.
    wordPoS := map[string][]string{}

    file, err := os.Open(fileName)
    if err != nil { return nil, err }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        aLine := strings.TrimSpace(scanner.Text())
        index := strings.Index(aLine, " ")
        if index < 0 {
            wordPoS[aLine] = nil
        } else {
            word := aLine[0: index]
            // pos is an array of strings
            pos := strings.Split(aLine[index + 1:], " ")
            wordPoS[word] = pos
        }
    } // for
    return wordPoS, nil

} // createPoSMap


// displayPoS - displays the entire contents of
// the wordPoS map, for debugging purposes.
// Params: map of Parts of Speech, keyed by word
func displayPoS(wordPoS map[string][]string) {
    for k, v := range wordPoS {
        fmt.Printf("%v ", k)
        for _, pos := range v {
            fmt.Printf("%v  ", pos)
        }
    }
} // displayPoS


// lowerCaseExceptI - makes a word lower case,
// except if it's the pronoun I.
// Params: word
// Return: lower case word
func lowerCaseExceptI(word string) string {
    if word == "I" { return word }
    if strings.HasPrefix(word, "I'") { return word }
    return strings.ToLower(word)
}


// makePronounTerm - creates a pronoun term based on the given word
// and its tag. Eg. pronoun(they, subject, third, plural).
// Note: This function does not handle the pronoun 'you'. 'You' is
// dealt with separately, by makeYouFacts().
//
// Params: word
//         lower case word
//         part of speech tag
// Return: complex term
//         success/failure flag
//
func makePronounTerm(word, lower, tag string) (Complex, bool) {

    var term Complex

    if strings.HasPrefix(tag, "PPS") { // PPS or PPSS
        if lower == "we" {
            term = Complex{pronoun, Atom(word), subject, first, plural}
        } else if lower == "they" {
            term = Complex{pronoun, Atom(word), subject, third, plural}
        } else if lower == "I" {
            term = Complex{pronoun, Atom(word), subject, first, singular}
        } else {  // he she it
            term = Complex{pronoun, Atom(word), subject, third, singular}
        }
    } else if strings.HasPrefix(tag, "PPO") {
        if lower == "us" {
            term = Complex{pronoun, Atom(word), object, first, plural}
        } else if lower == "them" {
            term = Complex{pronoun, Atom(word), object, third, plural}
        } else if lower == "me" {
            term = Complex{pronoun, Atom(word), object, first, singular}
        } else {
            term = Complex{pronoun, Atom(word), object, third, singular}
        }
    }

    if term == nil { return term, false}
    return term, true

} // makePronounTerm


// makeYouFacts - creates facts for the pronoun 'you', for example:
//     word(you, pronoun(you, subject, second, singular)).
//
// Params: word
// Return: facts
func makeYouFacts(word string) []RuleStruct {

    facts := []RuleStruct{}
    you := Atom(word)

    pronouns := []Complex{
        Complex{pronoun, you, subject, second, singular},
        Complex{pronoun, you, object, second, singular},
        Complex{pronoun, you, subject, second, plural},
        Complex{pronoun, you, object, second, plural},
    }

    for _, term := range pronouns {
        newTerm := Complex{WORD, you, term}
        fact := Fact(newTerm)
        facts = append(facts, fact)
    }

    return facts

} // makeYouFacts


// makeVerbTerm - creates a verb term, eg. verb(listen, present, base).
//
// Params: word
//         part of speech tag
// Return: term
//         success/failure flag
func makeVerbTerm(word string, tag string) (Complex, bool) {

    var term Complex

    if tag == "VB" {
        term = Complex{verb, Atom(word), present, base}
    } else if tag == "VBZ" {
        term = Complex{verb, Atom(word), present, third_sing}
    } else if tag == "VBD" {
        term = Complex{verb, Atom(word), past, past}
    } else if tag == "VBG" {
        term = Complex{participle, Atom(word), active}
    } else if tag == "VBN" {
        term = Complex{participle, Atom(word), passive}
    }

    if term == nil { return term, false}
    return term, true

} // makeVerbTerm


// makeNounTerm - creates a noun term, eg. noun(speaker, singular).
//
// Params: word
//         part of speech tag
// Return: term
//         success/failure flag
func makeNounTerm(word string, tag string) (Complex, bool) {

    var term Complex

    if tag == "NN" {
        term = Complex{noun, Atom(word), singular}
    } else if tag == "NNS" {
        term = Complex{noun, Atom(word), plural}
    } else if tag == "NNP" {
        term = Complex{noun, Atom(word), singular}
    }

    if term == nil { return term, false}
    return term, true

} // makeNounTerm


// makeAdjectiveTerm - creates an adjective term, eg. adjective(happy).
//
// Params: word
//         part of speech tag
// Return: term
//         success/failure flag
func makeAdjectiveTerm(word string, tag string) (Complex, bool) {

    var term Complex

    if tag == "JJ" {
        term = Complex{adjective, Atom(word), positive}
    } else if tag == "JJR" {
        term = Complex{adjective, Atom(word), comparative}
    } else if tag == "JJS" {
        term = Complex{adjective, Atom(word), superlative}
    }

    if term == nil { return term, false}
    return term, true

} // makeAdjectiveTerm

// makeArticleTerm - creates terms for articles, eg. article(the, definite).
//
// Params: word
// Return: term
//         success/failure flag
func makeArticleTerm(word string) (Complex, bool) {

    var term Complex

    wordLower := strings.ToLower(word)
    if wordLower == "the" {
        term = Complex{article, Atom(word), definite}
    } else {
        term = Complex{article, Atom(word), indefinite}
    }

    if term == nil { return term, false}
    return term, true

} // makeArticleTerm

// makeAdverbTerm - creates adverb terms, eg. adverb(happily).
//
// Params: word
// Return: term
//         success/failure flag
func makeAdverbTerm(word string) (Complex, bool) {
    term := Complex{adverb, Atom(word)}
    return term, true
} // makeAdverbTerm

// makePrepositionTerm - creates preposition terms, eg. preposition(from).
//
// Params: word
// Return: term
//         success/failure flag
func makePrepositionTerm(word string) (Complex, bool) {
    term := Complex{preposition, Atom(word)}
    return term, true
} // makePrepositionTerm


// makeUnknownTerm - creates terms for words with unknown part of speech.
//
// Params: word
// Return: term
//         success/failure flag
func makeUnknownTerm(word string) (Complex, bool) {
    term := Complex{unknown, Atom(word)}
    return term, true
} // makeUnknownTerm


// makeTerm - creates a complex term object for an English word.
// The second parameter is a part of speech tag, such as NNS or VBD.
// Tags are listed at the top of this file.
//
// Params: word
//         lower case word
//         part of speech tag
// Return: complex term
//         bool
func makeTerm(word, lower, tag string) (Complex, bool) {
    if strings.HasPrefix(tag, "VB") {
        return makeVerbTerm(word, tag)
    } else if strings.HasPrefix(tag, "NN") {
        return makeNounTerm(word, tag)
    } else if strings.HasPrefix(tag, "PP") {
        return makePronounTerm(word, lower, tag)
    } else if strings.HasPrefix(tag, "JJ") {
        return makeAdjectiveTerm(word, tag)
    } else if strings.HasPrefix(tag, "AT") {
        return makeArticleTerm(word)
    } else if strings.HasPrefix(tag, "IN") {
        return makePrepositionTerm(word)
    } else if strings.HasPrefix(tag, "RB") {
        return makeAdverbTerm(word)
    }
    return Complex{}, false
} // makeTerm


// wordToFacts - takes a word string and produces facts for the
// knowledge base.
//
// For some words, the part of speech is unambiguous. For example,
// 'the' can only be a definite article:
//
//      article(the, definite)
//
// Other words can have more than one part of speech. The word
// 'envy', for example, might be a noun or a verb.
//
//      noun(envy, singular)
//      verb(envy, present, base)
//
// For 'envy', a parsing algorithm must be able to test both
// possibilities. Therefore, the inference engine will need two
// facts for the knowledge base:
//
//      word(envy, noun(envy, singular)).
//      word(envy, verb(envy, present, base)).
//
// Params: word (string)
//         part of speech data (map)
// Return: facts
//
func wordToFacts(word string, pos map[string][]string) []RuleStruct {

    lower := lowerCaseExceptI(word)

    // Handle pronoun 'you', which is very ambiguous.
    if lower == "you" { return makeYouFacts(word) }

    length := len(word)
    if length == 1 { // Maybe this is punctuation.
        term, ok := makePunctuationTerm(word)
        if ok {
            wordTerm := Complex{WORD, Atom(word), term}
            return []RuleStruct{ Fact(wordTerm) }
        }
    }

    facts := []RuleStruct{}

    posData := pos[word]
    if posData == nil {
        posData = pos[lower]
    }

    if posData != nil && len(posData) > 0 {
        for _, pos := range posData {
            term, ok := makeTerm(word, lower, pos);
            if ok {
                wordTerm := Complex{WORD, Atom(word), term}
                fact := Fact(wordTerm)
                facts = append(facts, fact)
            }
        }
    }

    if len(facts) < 1 {
        term := Complex{unknown, Atom(word)}
        wordTerm := Complex{WORD, Atom(word), term}
        fact := Fact(wordTerm)
        facts = append(facts, fact)
    }

    return facts

} // wordToFacts


// makeFacts - takes a list of words, and creates a list of
// facts which can be analyzed by the inference engine.
// The word 'envy', for example, should produce two facts.
//
//      word(envy, noun(envy, singular)).
//      word(envy, verb(envy, present, base)).
//
// Note: A Fact is the same as a Rule without a body.
//
// Params: list of words
//         map of part of speech data
// Return: list of facts
//
func makeFacts(words []string, pos map[string][]string) []RuleStruct {
    facts := []RuleStruct{}
    for _, word := range words {
        wordFacts := wordToFacts(word, pos)
        for _, wordFact := range wordFacts {
            facts = append(facts, wordFact)
        }
    }
    return facts
} // makeFacts
