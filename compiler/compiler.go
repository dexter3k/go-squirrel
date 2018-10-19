package compiler

import (
	"io"
	"log"

	"github.com/SMemsky/go-squirrel/compiler/lexer"
	"github.com/SMemsky/go-squirrel/compiler/lexer/tokens"
)

type compiler struct {
	lexer lexer.Lexer

	token     tokens.Token
	tokenInfo lexer.TokenInfo
	lastToken tokens.Token
}

func NewCompiler(rr io.Reader) *compiler {
	return &compiler{
		lexer: lexer.NewLexer(rr),
	}
}

func (c *compiler) Compile() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(c.token)
			log.Println("Recovered from:", r)
		}
	}()
	c.lex()
	for c.token != 0 {
		c.statement()
		if c.lastToken != '}' && c.lastToken != ';' {
			c.optionalSemicolon()
		}
	}
}

func (c *compiler) lex() {
	repeat := true
	for repeat {
		c.lastToken = c.token
		var err error
		c.tokenInfo, err = c.lexer.Lex()
		c.token = c.tokenInfo.Token
		if err != nil {
			panic(err)
		}
		repeat = c.token == '\n'
	}
}

func (c *compiler) statement() {
	c.prefixedExpression()
}

func (c *compiler) prefixedExpression() {
	c.factor()
	if c.token == '(' {
		c.lex()
		c.functionCallArgs()
	}
}

func (c *compiler) factor() {
	if c.token == tokens.Identifier {
		c.lex()
	}
}

func (c *compiler) functionCallArgs() {
	for c.token != ')' {
		// c.expression()
		// if c.token == ',' {
		//     c.lex()
		//     if c.token == ')' {
		//         panic(ErrExpectArgument)
		//     }
		// }
		c.lex()
	}
	c.lex()
}

func (c *compiler) optionalSemicolon() {
	if c.token == ';' {
		c.lex()
		return
	}
	if !c.isEndOfStatement() {
		panic(ErrExpectStatementEnd)
	}
}

func (c *compiler) isEndOfStatement() bool {
	return c.token == 0 || c.lastToken == '\n' || c.token == '}' || c.token == ';'
}
