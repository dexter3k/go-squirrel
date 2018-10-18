package main

import (
	"fmt"
	"strings"

	"github.com/SMemsky/go-squirrel/compiler/lexer"
)

const programSource = `

// /= */

`

func main() {
	lex := lexer.NewLexer(strings.NewReader(programSource))
	for {
		token := lex.Lex()
		if token.Token == 0 {
			break
		}
		fmt.Println(uint(token.Token))
	}
}
