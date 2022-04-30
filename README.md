# Suiron - A Golang Inference Engine.

Suiron is an inference engine written in Go. The rule declaration syntax is very similar to Prolog, but there are some differences.

This brief README does not present a detailed explanation of how inference engines work, so a basic understanding of Prolog is required. Documentation will be expanded in time.

## Briefly

Suiron analyzes facts and rules which are recorded in a knowledge base. These facts and rules can be loaded from a text file, or created dynamically within a Go application program.

In the code sample below, the fact 'mother(June, Theodore).', meaning 'June is the mother of Theodore', is defined in a Go source program by calling the function ParseComplex().

```
    fact := ParseComplex("mother(June, Theodore).")
```

In a text-format file, the above fact would be written:

```
mother(June, Theodore).
```

Please refer to [suiron/complex.go](suiron/complex.go).

In Prolog, words which begin with a lower case letter (eg. mother) are atoms, and words which begin with an upper case letter are variables. In Suiron, atoms can be upper case or lower case. Thus 'June' and 'Theodore' are atoms. Suiron's atoms are implemented as string constants. They can even contain spaces.

```
mother(June, The Beaver).
```

Suiron also supports integer and floating point numbers, which are implemented as int64 and float64. They are parsed by Go's strconv package:

```
    f, err := strconv.ParseFloat(str, 64)
    i, err := strconv.ParseInt(str, 10, 64)
```

Please refer to [suiron/constants.go](suiron/constants.go).

Suiron's variables are defined by putting a dollar sign in front of the variable name, for example, $Child. In the Go code sample below, 'child' is a logic variable.

```
mother   := Atom("mother")
June     := Atom("June")
child, _ := LogicVar("$Child")
goal     := MakeGoal(mother, June, child)
```

Please refer to [suiron/variable.go](suiron/variable.go) and [suiron/goal.go](suiron/goal.go).

LogicVar() returns two values, a Suiron logic variable and a parsing error, if there is an error. In a source file of facts and rules, the goal above would be written:

```
mother(June, $Child).
```

The <i>anonymous variable</i> must also begin with a dollar sign: $\_ . A simple underscore '\_' is treated as an atom. Below is an example a rule which contains an anonymous variable:

```
voter($P) :- $P = person($_, $Age), $Age >= 18.
```

Of course, Suiron supports linked lists, which work the same way as Prolog lists. A linked list can be defined dynamically:

```
    list := MakeLinkedList(true, a, b, c, Tail)
```

or loaded from a text file:

```
   ..., list = [a, b, c | $Tail], ...
```

Please refer to [suiron/linkedlist.go](suiron/variable.go).

## Requirements

Suiron was developed and tested with Go version 1.17.

[https://go.dev/](https://go.dev/)

## Cloning

To clone the repository, run the following command in a terminal window:

```
 git clone git@github.com:Indrikoterio/suiron.git
```

The repository has three folders:

```
 suiron/suiron
 suiron/test
 suiron/demo
```

The code for the inference engine itself is in the subfolder /suiron.

The subfolder /test contains Go programs which test the basic functionality of Suiron.

The subfolder /demo contains a simple demo program which parses English sentences.

## Usage

In the top folder is a program called query.go, which loads facts and rules from a file, and allows the user to query the knowledge base. Query can be run in a terminal window as follows:

```
go run query.go test/kings.txt
?- father($F, $C).
$F = Godwin, $C = Harold II
$F = Godwin, $C = Tostig
$F = Godwin, $C = Edith
$F = Tostig, $C = Skule
$F = Harold II, $C = Harold
No
?-
```

To use Suiron in your own project, copy the subfolder 'suiron' to your project folder. You will have to include:

```
import (
    . "github.com/indrikoterio/suiron/suiron"
)
```

in your source file.

The program [demo/parse_demo.go](demo/parse_demo.go) demonstrates how to set up a knowledge base and make queries. If you intend to incorporate Suiron into your own project, this is a good reference. There are detailed comments in the header.

To run parse_demo, move to the demo folder and execute the batch file 'run'.

```
 cd demo
 ./run
```

Suiron doesn't have a lot of built-in predicates, but it does have:

```
append, functor, print, nl, include, exclude, greater_than (etc.)
```

...and some arithmetic functions:

```
add, subtract, multply, divide
```

Please refer to the test programs for examples of how to use these.

To run the tests, open a terminal window, go to the test folder, and execute 'run'.

```
 cd test
 ./run
```

If you need write your own built-in predicates and functions (in Go), refer to [test/hyphenate.go](test/hyphenate.go) and [test/capitalize.go](test/capitalize.go) to see how this is done.

## Developer

Suiron was developed by Cleve (Klivo) Lendon.

## Contact

To contact the developer, send email to indriko@yahoo.com . Comments, suggestions and criticism are welcomed.

## History

First release, April 2022.

## Reference

The code structure of this inference engine is inspired by the Predicate Calculus Problem Solver presented in chapters 23 and 24 of 'AI Algorithms...' by Luger and Stubblefield. I highly recommend this book.

```
AI Algorithms, Data Structures, and Idioms in Prolog, Lisp, and Java
George F. Luger, William A. Stubblefield, ©2009 | Pearson Education, Inc. 
ISBN-13: 978-0-13-607047-4
ISBN-10: 0-13-607047-7
```

## License

The source code for Suiron is licensed under the MIT license, which you can find in [LICENSE](LICENSE).
