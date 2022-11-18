package compiler

import (
	"fmt"
	"io"
	"log"

	"github.com/dexter3k/go-squirrel/sqvm"
	"github.com/dexter3k/go-squirrel/compiler/lexer"
	"github.com/dexter3k/go-squirrel/compiler/lexer/tokens"
)

// Compile the code from reader and push resulting closure onto vm stack.
func Compile(vm *sqvm.VM, filename string, r io.Reader) (*sqvm.FuncProto, error) {
	return NewCompiler(r).Compile()
}

type compiler struct {
	lexer lexer.Lexer

	f *state

	token     tokens.Token
	tokenInfo lexer.TokenInfo
	lastToken tokens.Token
}

func NewCompiler(rr io.Reader) *compiler {
	return &compiler{
		lexer: lexer.NewLexer(rr),
	}
}

func (c *compiler) Compile() (*sqvm.FuncProto, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(c.token)
			log.Println("Recovered from:", r)
		}
	}()

	// Create root function definition with args "this" and "vargv"
	c.f = newState()

	c.lex()

	for c.token != 0 {
		c.statement()

		if c.lastToken != '}' && c.lastToken != ';' {
			c.optionalSemicolon()
		}
	}

	// Return the built function prototype
	return c.f.makeFuncProto()
}

func (c *compiler) lex() {
	for true {
		c.lastToken = c.token

		var err error
		c.tokenInfo, err = c.lexer.Lex()
		c.token = c.tokenInfo.Token
		if err != nil {
			panic(err)
		}
		if c.token != '\n' {
			break
		}
	}
}

func (c *compiler) statement() {
	switch c.token {
	case ';':
		c.lex()
	default:
		c.commaExpression()
		c.f.popTarget()
	}
}

func (c *compiler) commaExpression() {
	c.expression()

	for c.token == ',' {
		c.lex()

		// Discard result for the last expression
		c.f.popTarget()

		c.expression()
	}
}

func (c *compiler) expression() {
	c.prefixedExpression()
}

func (c *compiler) prefixedExpression() error {
	if err := c.factor(); err != nil {
		return err
	}

	for {
		switch c.token {
		case '(':
			c.lex()
			c.functionCallArgs()
		default:
			return nil
		}
	}
}

func (c *compiler) factor() error {
	switch c.token {
	case tokens.Identifier:
		id := c.f.makeString(c.tokenInfo.String)
		fmt.Printf("factor: identifier %d->%q\n", id, c.tokenInfo.String)
		c.lex()
	case tokens.StringLiteral:
		fmt.Printf("factor: string\n")

		c.lex()
	case tokens.Integer:
		fmt.Printf("factor: integer\n")

		c.lex()
	case tokens.Float:
		fmt.Printf("factor: float\n")

		c.lex()
	default:
		return ErrExpectExpression
	}

	return nil
}

func (c *compiler) functionCallArgs() {
	for c.token != ')' {
		c.expression()
		if c.token == ',' {
			c.lex()
			if c.token == ')' {
				panic(ErrExpectArgument)
			}
		}
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
