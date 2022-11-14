package main

import (
    "github.com/dexter3k/go-squirrel/sqvm"
)

func main() {
    sq := sqvm.Open(1024)
    defer sq.Close()

    sq.SetPrintFunc(
        func(vm *sqvm.VM, format string, args ...any){
        },
        func(vm *sqvm.VM, format string, args ...any){
        },
    )

    sq.PushRootTable()
}
