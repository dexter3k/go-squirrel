package main

import (
	// "fmt"
	"strings"

	"github.com/dexter3k/go-squirrel/compiler"
)

const programSource = `

print("Hello, world!")

`

func main() {
	comp := compiler.NewCompiler(strings.NewReader(programSource))
	_, err := comp.Compile()
	if err != nil {
		panic(err)
	}
}
