package compiler

import (
	"io"
	"log"

	"github.com/dexter3k/go-squirrel/sqvm"
	"github.com/dexter3k/go-squirrel/compiler/lexer"
	"github.com/dexter3k/go-squirrel/compiler/lexer/tokens"
)

func Compile(vm *sqvm.VM, filename string, r io.Reader) error {
	return nil
}

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
	c.commaExpression()
	// pop result of unused expression
}

func (c *compiler) commaExpression() {
	c.expression()
	for c.token == ',' {
		c.lex()
		c.expression()
		// pop result of last expression
	}
	// result of first expression always stays on stack and gets returned
}

func (c *compiler) expression() {
	c.prefixedExpression()
}

func (c *compiler) prefixedExpression() {
	c.factor()
	for {
		switch c.token {
		case '(':
			c.lex()
			c.functionCallArgs()
		default:
			return
		}
	}
}

func (c *compiler) factor() {
	switch c.token {
	case tokens.Identifier:
		c.lex()
	case tokens.StringLiteral:
		c.lex()
	case tokens.Integer:
		c.lex()
	case tokens.Float:
		c.lex()
	default:
		panic(ErrExpectExpression)
	}
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
