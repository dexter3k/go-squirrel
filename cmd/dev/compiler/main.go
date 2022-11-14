package main

import (
	"fmt"
	"strings"

	"github.com/dexter3k/go-squirrel/compiler"
	"github.com/dexter3k/go-squirrel/compiler/lexer"
)

const programSource = `

print("Hello, world!"), print("Comma!")

`

func main() {
	compiler := compiler.NewCompiler(strings.NewReader(programSource))
	if compiler != nil {
		fmt.Println("Compiler ok")
	}
	compiler.Compile()

	lex := lexer.NewLexer(strings.NewReader(programSource))
	for {
		token, err := lex.Lex()
		if err != nil {
			fmt.Println(err, token.String)
			break
		}
		if token.Token == 0 {
			break
		}
		fmt.Println(uint(token.Token), token.String, token.Integer, token.Float)
	}
}
