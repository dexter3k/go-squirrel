package lexer

import (
	"bufio"
	"io"

	"github.com/SMemsky/go-squirrel/compiler/lexer/tokens"
)

var keywords = map[string]tokens.Token{
	"while":       tokens.While,
	"do":          tokens.Do,
	"if":          tokens.If,
	"else":        tokens.Else,
	"break":       tokens.Break,
	"continue":    tokens.Continue,
	"return":      tokens.Return,
	"null":        tokens.Null,
	"function":    tokens.Function,
	"local":       tokens.Local,
	"for":         tokens.For,
	"foreach":     tokens.Foreach,
	"in":          tokens.In,
	"typeof":      tokens.Typeof,
	"base":        tokens.Base,
	"delete":      tokens.Delete,
	"try":         tokens.Try,
	"catch":       tokens.Catch,
	"throw":       tokens.Throw,
	"clone":       tokens.Clone,
	"yield":       tokens.Yield,
	"resume":      tokens.Resume,
	"switch":      tokens.Switch,
	"case":        tokens.Case,
	"default":     tokens.Default,
	"this":        tokens.This,
	"class":       tokens.Class,
	"extends":     tokens.Extends,
	"constructor": tokens.Constructor,
	"instanceof":  tokens.InstanceOf,
	"true":        tokens.True,
	"false":       tokens.False,
	"static":      tokens.Static,
	"enum":        tokens.Enum,
	"const":       tokens.Const,
	"__FILE__":    tokens.Line,
	"__LINE__":    tokens.File,
}

type TokenInfo struct {
	Token  tokens.Token
	String string

	Line   uint
	Column uint
}

type lexer struct {
	source *bufio.Reader

	currentChar rune
	nextChar    rune
}

func NewLexer(rr io.Reader) *lexer {
	l := &lexer{
		source: bufio.NewReader(rr),
	}

	l.next()

	return l
}

func (l *lexer) Lex() TokenInfo {
	l.next()
	for ; l.currentChar != 0; l.next() {
		switch l.currentChar {
		case '\t', '\r', ' ', '\n':
			break
		case '#':
			l.skipLineComment()
		case '/':
			switch l.nextChar {
			case '/':
				l.skipLineComment()
			case '*':
				l.skipBlockComment()
			case '=':
				l.next()
				return TokenInfo{Token: tokens.DivideEqual}
			case '>':
				l.next()
				return TokenInfo{Token: tokens.AttributeClose}
			default:
				return TokenInfo{Token: tokens.Token('/')}
			}
		case '=':
			if l.nextChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.Equal}
			}
			return TokenInfo{Token: tokens.Token('=')}
		case '<':
			switch l.nextChar {
			case '=':
				l.next()
				if l.nextChar == '>' {
					l.next()
					return TokenInfo{Token: tokens.ThreeWayCompare}
				}
				return TokenInfo{Token: tokens.LessEqual}
			case '-':
				l.next()
				return TokenInfo{Token: tokens.NewSlot}
			case '<':
				l.next()
				return TokenInfo{Token: tokens.ShiftLeft}
			case '/':
				l.next()
				return TokenInfo{Token: tokens.AttributeOpen}
			}
			return TokenInfo{Token: tokens.Token('<')}
		}
	}

	return TokenInfo{}
}

func (l *lexer) next() {
	l.currentChar = l.nextChar
	var err error
	if l.nextChar, _, err = l.source.ReadRune(); err != nil {
		l.nextChar = 0
	}
}

func (l *lexer) skipLineComment() {
	for l.currentChar != '\n' && l.currentChar != 0 {
		l.next()
	}
}

func (l *lexer) skipBlockComment() {
	l.next()
	for {
		l.next()
		switch l.currentChar {
		case 0:
			return
		case '*':
			if l.nextChar == '/' {
				l.next()
				return
			}
		}
	}
}
