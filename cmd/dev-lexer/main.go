package main

import (
	"fmt"
	"strings"

	"github.com/SMemsky/go-squirrel/compiler/lexer"
)

const programSource = `

// /= */

.
.........-><---
%=
%
// ..

foobar

07776x
0xfffffffffffffffffffff

`

func main() {
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
		fmt.Println(uint(token.Token), token.String)
	}
}
