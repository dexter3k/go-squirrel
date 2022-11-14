package main

import (
	"fmt"
	"strings"

	"github.com/dexter3k/go-squirrel/compiler/lexer"
)

// const programSource = `

// // /= */

// .
// .........-><---
// %=
// %
// // ..

// foobar

// 100

// 07776x
// 0x100

// 0.2000001
// 100.0
// 100.

// 1.0e309
// 1e300
// 1e+10

// foobar
// while
// true
// false
// __FILE__
// __LINE__

// @"
// каждую пятницу
// одно
// и...

// \n

// "

// "Foo\nbar"

// 'a'
// 'b'
// '\t'

// `

const programSource = `
@"сомтхинг"

"Foo\nbar"

'a'
'b'
'\t'


</ flippy = 10, second = [1, 2, 3] /> // attrs

`

func main() {
	lex := lexer.NewLexer(strings.NewReader(programSource))
	for {
		token, err := lex.Lex()
		if err != nil {
			fmt.Printf("%s, %q\n", err, token.String)
			break
		}

		if token.Token == 0 {
			break
		}

		fmt.Printf("type=%3d, s=%q, i=%d, f=%f\n", uint(token.Token), token.String, token.Integer, token.Float)
	}
}
