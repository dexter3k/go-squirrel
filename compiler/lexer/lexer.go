package lexer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/SMemsky/go-squirrel/compiler/lexer/tokens"
)

var (
	ErrUnknownToken = fmt.Errorf("Unknown token")
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

func (l *lexer) Lex() (TokenInfo, error) {
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
				return TokenInfo{Token: tokens.DivideEqual}, nil
			case '>':
				l.next()
				return TokenInfo{Token: tokens.AttributeClose}, nil
			default:
				return TokenInfo{Token: tokens.Token('/')}, nil
			}
		case '=':
			if l.nextChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.Equal}, nil
			}
			return TokenInfo{Token: tokens.Token('=')}, nil
		case '<':
			switch l.nextChar {
			case '=':
				l.next()
				if l.nextChar == '>' {
					l.next()
					return TokenInfo{Token: tokens.ThreeWayCompare}, nil
				}
				return TokenInfo{Token: tokens.LessEqual}, nil
			case '-':
				l.next()
				return TokenInfo{Token: tokens.NewSlot}, nil
			case '<':
				l.next()
				return TokenInfo{Token: tokens.ShiftLeft}, nil
			case '/':
				l.next()
				return TokenInfo{Token: tokens.AttributeOpen}, nil
			}
			return TokenInfo{Token: tokens.Token('<')}, nil
		case '>':
			if l.nextChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.GreaterEqual}, nil
			} else if l.nextChar == '>' {
				l.next()
				if l.nextChar == '>' {
					l.next()
					return TokenInfo{Token: tokens.UnsignedShiftRight}, nil
				}
				return TokenInfo{Token: tokens.ShiftRight}, nil
			}
			return TokenInfo{Token: tokens.Token('<')}, nil
		case '!':
			if l.nextChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.NotEqual}, nil
			}
			return TokenInfo{Token: tokens.Token('!')}, nil
		case '{', '}', '(', ')', '[', ']':
		case ';', ',', '?', '^', '~':
			return TokenInfo{Token: tokens.Token(l.currentChar)}, nil
		case '.':
			if l.nextChar != '.' {
				return TokenInfo{Token: tokens.Token('.')}, nil
			}
			l.next()
			if l.nextChar != '.' {
				return TokenInfo{String: ".."}, ErrUnknownToken
			}
			l.next()
			return TokenInfo{Token: tokens.VarParams}, nil
		case '&':
			if l.nextChar == '&' {
				l.next()
				return TokenInfo{Token: tokens.And}, nil
			}
			return TokenInfo{Token: tokens.Token('&')}, nil
		case '|':
			if l.nextChar == '|' {
				l.next()
				return TokenInfo{Token: tokens.Or}, nil
			}
			return TokenInfo{Token: tokens.Token('|')}, nil
		case ':':
			if l.nextChar == ':' {
				l.next()
				return TokenInfo{Token: tokens.DoubleColon}, nil
			}
			return TokenInfo{Token: tokens.Token(':')}, nil
		case '%':
			if l.nextChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.ModuloEqual}, nil
			}
			return TokenInfo{Token: tokens.Token('%')}, nil
		case '+':
			if l.nextChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.PlusEqual}, nil
			}
			if l.nextChar == '+' {
				l.next()
				return TokenInfo{Token: tokens.Increase}, nil
			}
			return TokenInfo{Token: tokens.Token('+')}, nil
		case '-':
			if l.nextChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.MinusEqual}, nil
			}
			if l.nextChar == '-' {
				l.next()
				return TokenInfo{Token: tokens.Decrease}, nil
			}
			return TokenInfo{Token: tokens.Token('-')}, nil
		default:
			return TokenInfo{Token: tokens.Token(l.currentChar)}, nil
		}
	}

	return TokenInfo{}, nil
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
