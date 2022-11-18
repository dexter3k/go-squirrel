package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/dexter3k/go-squirrel/sqvm"
    "github.com/dexter3k/go-squirrel/compiler"
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

func runFile(vm *sqvm.VM, filename string, args []string) int {
    f, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Unable to open %q: %w\n", filename, err)
        return 1
    }
    defer f.Close()

    // Compiler pushes the resulting closure onto the vm stack
    if err := compiler.Compile(vm, filename, f); err != nil {
        fmt.Printf("Unable to compile file %q: %w\n", filename, err)
        return 1
    }

    // Push the args
    vm.PushRootTable() // root table as the local space of the script
    for _, arg := range args {
        vm.PushString(arg)
    }

    // Perform the call
    if err := vm.Call(1 + len(args), true, true); err != nil {
        fmt.Printf("Call failed with err %w\n", err)
        return 1
    }

    // Expect an integer return type
    retType := vm.GetType(-1)
    if retType == sqvm.TypeInteger {
        return int(vm.GetInteger(-1))
    }

    return 0
}

func main() {
    os.Exit(mainWithCode())
}

func mainWithCode() int {
    if *showVersionInfo {
        fmt.Printf("go-squirrel wip\n")
        return 0
    }

    vm := sqvm.Open(1024)
    defer vm.Close()

    vm.SetPrintFunc(
        func(vm *sqvm.VM, format string, args ...any){
        },
        func(vm *sqvm.VM, format string, args ...any){
        },
    )

    vm.PushRootTable()
    // Register libraries
    // Register error handlers

    if len(flag.Args()) == 0 {
        fmt.Printf("Interpreter mode is not yet implemented")
        return 1
    } else {
        return runFile(vm, flag.Args()[0], flag.Args()[1:])
    }
}
