# Suiron - A Golang Inference Engine

Suiron is an inference engine written in Go. The rule declaration syntax is very similar to Prolog, but there are some differences.

This brief README does not present a detailed explanation of how inference engines work, so a basic understanding of Prolog is required. Documentation will be expanded in time.

## Briefly

An inference engines analyzes facts and rules which are stored in a knowledge base. Suiron has a parser which loads these facts and rules from a text-format source file.

Below is an example of a fact, which means "June is the mother of Theodore":

```
mother(June, Theodore).
```

Here we see the main difference between Suiron and Prolog. In Prolog, lower case words are 'atoms' (that is, string constants) and upper case words are variables. In Suiron, atoms can be lower case or upper case. Thus 'mother', 'June' and 'Theodore' are all atoms. Suiron's atoms can even contain spaces.

```
mother(June, The Beaver).
```

Suiron's variables are defined by putting a dollar sign in front of the variable name, for example, $Child. A query to determine June's children would be written:

```
mother(June, $Child).
```

Please refer to [variable.go](suiron/variable.go).

The [anonymous](suiron/anonymous.go) variable must also begin with a dollar sign: $\_ . A simple underscore '\_' is treated as an atom. Below is an example a rule which contains an anonymous variable:

```
voter($P) :- $P = person($_, $Age), $Age >= 18.
```

<hr><br>

Facts and rules can also be created dynamically within a Go application program. The fact
mother(June, Theodore) could be created by calling the function ParseComplex().

```
    fact := ParseComplex("mother(June, Theodore).")
```

Please refer to [complex.go](suiron/complex.go).

The query mother(June, $Child) could be created in Go as follows:

```
mother   := Atom("mother")
June     := Atom("June")
child, _ := LogicVar("$Child")
query    := MakeGoal(mother, June, child)
```

Please refer to [variable.go](suiron/variable.go) and [goal.go](suiron/goal.go) for more details.

Suiron also supports integer and floating point numbers, which are implemented as int64 and float64. These are parsed by Go's strconv package:

```
    f, err := strconv.ParseFloat(str, 64)
    i, err := strconv.ParseInt(str, 10, 64)
```

If a Float and an Integer are compared, the Integer will be converted to a Float for the comparison.

Please refer to [constants.go](suiron/constants.go).

Of course, Suiron supports linked lists, which work the same way as Prolog lists. A linked list can be loaded from a file:

```
   ..., [a, b, c, d] = [$Head | $Tail], ...
```

or created dynamically:

```
    X := ParseLinkedList("[a, b, c, d]")
    Y := MakeLinkedList(true, $Head, $Tail)
```

Please refer to [linkedlist.go](suiron/linkedlist.go).

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

In the top folder is a program called [query.go](query.go), which loads facts and rules from a file, and allows the user to query the knowledge base. Query can be run in a terminal window as follows:

```
go run query.go test/kings.txt
```

The user will be prompted for a query with this prompt: ?-

The query below will print out all father/child relationships.

```
?- father($F, $C).
```

After typing enter, the program will print out solutions, one after each press of Enter, until there are no more solutions, as indicated by 'No'.

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

The program [parse_demo.go](demo/parse_demo.go) demonstrates how to set up a knowledge base and make queries. If you intend to incorporate Suiron into your own project, this is a good reference. There are detailed comments in the header.

To run parse_demo, move to the demo folder and execute the batch file 'run'.

```
 cd demo
 ./run
```

Suiron doesn't have a lot of built-in predicates, but it does have: [append.go](suiron/append.go), [functor.go](suiron/functor.go), [print.go](suiron/print.go), [new_line.go](suiron/new_line.go), [include.go](suiron/include.go), [exclude.go](suiron/exclude.go), greater_than (etc.)


...and some arithmetic functions: [add.go](suiron/add.go), [subtract.go](suiron/subtract.go), [multiply.go](suiron/multiply.go), [divide.go](suiron/divide.go)

Please refer to the test programs for examples of how to use these.

To run the tests, open a terminal window, go to the test folder, and execute 'run'.

```
 cd test
 ./run
```

Suiron allows you to write your own built-in predicates and functions. The files [bip_template.go](suiron/bip_template.go) and [bif_template.go](suiron/bif_template.go) can be used as templates. Please read the comments in the headers of these files.

The files [hyphenate.go](test/hyphenate.go) and [capitalize.go](test/capitalize.go) in the test directory can also be used for reference.

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
