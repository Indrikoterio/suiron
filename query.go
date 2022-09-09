package main

// Query - allows a user to query a Suiron knowledge base.
// Usage:
//
// > go run query.go test/kings.txt
//
// The command above will load facts and rules from kings.txt,
// and then prompt for a query with:
//
// Successfully loaded facts and rules from test/kings.txt
// ?- 
//
// To find the grandchildren of Godwin, the user would type in
// the following query:
//
// ?- grandfather(Godwin, $X)
//
// ($X is a variable.)
//
// The Suiron inference engine will output one result after each
// press of 'enter'. When solutions are exhausted, the inference
// engine will print out 'No'.
//
// grandfather(Godwin, Harold)
// grandfather(Godwin, Skule)
// No
// ?-
//
// Type <enter> to end the program.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "strings"
    "errors"
    "bufio"
    "fmt"
    "os"
)

func main() {

    kb := KnowledgeBase{}

    // Is there a file name?
    if len(os.Args) > 1 {
        fileName := os.Args[1]
        if _, err := os.Stat(fileName); err == nil {
            // Read in facts and rules.
            err = LoadKBFromFile(kb, fileName)
            if err != nil { 
                fmt.Println(err.Error())
                return
            }
            fmt.Printf("Successfully loaded facts and rules from %v\n", fileName)
        } else if errors.Is(err, os.ErrNotExist) {
            fmt.Printf("The file %v does not exist.\n", fileName)
            return
        }
    } else {
        fmt.Println("This is Suiron, an inference engine written in Go by Cleve Lendon.")
        fmt.Println("To load knowledge> go run query rules.txt")
    }

    reader := bufio.NewReader(os.Stdin)
    previous := ""   // Previous query.

    for {

        fmt.Print("?- ")  // Prompt for query.

        q, _ := reader.ReadString('\n')
        query := strings.TrimSpace(q)
        if len(query) == 0 { break }

        if query == "." {
            query = previous
        } else {
            previous = query
        }

        goal, err := ParseGoal(query)
        if err != nil {
            fmt.Println(err.Error())
            continue
        }

        // Get the root solution node.
        root := goal.GetSolver(kb, SubstitutionSet{}, nil)

        for {
            solution, found := root.NextSolution()
            if !found {
                fmt.Println("No")
                break
            }
            result := FormatSolution(goal, solution)
            fmt.Print(result)
            _, _ = reader.ReadString('\n')
        }
    } // for

} // main
