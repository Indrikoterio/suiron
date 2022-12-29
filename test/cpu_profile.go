package main

// Profile cpu usage.
// Based on TestTime, which executes a qsort algorithm.
//
// Cleve Lendon

import (
    . "github.com/indrikoterio/suiron/suiron"
    "fmt"
    "os"
    "runtime/pprof"
)

func main() {

    cpufile, err := os.Create("cpu.pprof")
    if err != nil { panic(err) }
    err = pprof.StartCPUProfile(cpufile)
    if err != nil { panic(err) }
    defer cpufile.Close()
    defer pprof.StopCPUProfile()

    fmt.Println("Profile CPU usage.")

    kb := KnowledgeBase{}
    err = LoadKBFromFile(kb, "qsort.txt")
    if err != nil { fmt.Println(err.Error()) }

    goal, _ := ParseQuery("measure")

    solution, _ := Solve(goal, kb, SubstitutionSet{})
    fmt.Println(solution)

} // TestTime
