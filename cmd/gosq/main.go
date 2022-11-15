package main

import (
    "flag"
    "fmt"

    "github.com/dexter3k/go-squirrel/sqvm"
)

var (
    showVersionInfo = flag.Bool("v", false, "display version info")
)

func init() {
    flag.Usage = func() {
        fmt.Printf("Usage of gosq: [OPTIONS] [program.nut [ARGS]]\n")
        flag.PrintDefaults()
    }
    flag.Parse()
}

func main() {
    if *showVersionInfo {
        fmt.Printf("go-squirrel wip\n")
        return
    }

    sq := sqvm.Open(1024)
    defer sq.Close()

    sq.SetPrintFunc(
        func(vm *sqvm.VM, format string, args ...any){
        },
        func(vm *sqvm.VM, format string, args ...any){
        },
    )

    sq.PushRootTable()

    // Register libraries

    // Register error handlers

    if len(flag.Args()) == 0 {
        fmt.Printf("Interpreter mode is not yet implemented")
        return
    } else {
        // Execute the file...
    }
}
