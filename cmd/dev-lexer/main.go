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

100

07776x
0x100

0.2000001
100.0
100.

1.0e309
1e300
1e+10

foobar
while
true
false
__FILE__
__LINE__

@"
каждую пятницу
одно
и...

\n

"

"Foo\nbar"

'a'
'b'
'\t'

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
		fmt.Println(uint(token.Token), token.String, token.Integer, token.Float)
	}
}
