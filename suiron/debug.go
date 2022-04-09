package suiron

// debug - has methods for printing debug messages.
//
// Cleve Lendon

import (
    "fmt"
    "reflect"
)

// DBG - prints a debug message.
// This function takes a variable number of arguments.
// In order to remove debug messages after testing,
// it is useful to include a string of x's, For example:
//     DBG("xxxxxx function arguments ", args)
// When testing is done, debugging messages can be grepped:
//     > grep xxx *.go
func DBG(args ...interface{}) {
    for _, arg := range args {
        fmt.Printf("%v ", arg)
    }
    fmt.Println("")
}

// DBKB - prints all rules and facts in the knowledge base.
func DBKB(kb KnowledgeBase) {
    fmt.Println(kb.FormatKB())
}

// displayLinkedListNode - displays one linked list node,
// indented according to level.
func displayLinkedListNode(level string, listPtr *LinkedListStruct) {
    var nilOrNot string
    fmt.Printf("%v %v\n", level, reflect.TypeOf(*listPtr))
    fmt.Printf("%v term: %v   %v\n", level, listPtr.term, reflect.TypeOf(listPtr.term))
    if listPtr.next == nil {
        nilOrNot = "nil"
    } else {
        nilOrNot = "not nil"
    }
    fmt.Printf("%v next: %v\n", level, nilOrNot)
    fmt.Printf("%v count: %d\n", level, listPtr.count)
    fmt.Printf("%v tail var: %v\n", level, listPtr.tailVar)
}

// DBLL - displays a linked list for debugging purposes.
func DBLL(message string, list LinkedListStruct) {
    level := ""
    fmt.Printf("%v\n", message)
    ptr := &list
    for ptr != nil {
        displayLinkedListNode(level, ptr)
        level = level + "    "
        ptr = ptr.next
    }
}
