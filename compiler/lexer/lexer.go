package lexer

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/SMemsky/go-squirrel/compiler/lexer/tokens"
)

var (
	ErrUnknownToken = fmt.Errorf("Unknown token")
	ErrHexOverflow  = fmt.Errorf("Hexadecimal number is too big")
	ErrOctOverflow  = fmt.Errorf("Octal number is too big")
	ErrDecOverflow  = fmt.Errorf("Decimal number is too big")
	ErrFloatFormat  = fmt.Errorf("Malformed floating point value")
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
	Token   tokens.Token
	String  string
	Integer uint64
	Float   float64

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
			if isDigit(l.currentChar) {
				return l.readNumber()
			}
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

// Integer, Hex, Octal, Float, Scientific
// Hex format:
// 0xFFFFFF, 0Xfffffff, 0xffffff, 0XFFFFFF
// Octal format:
// 07777777
// Scientific format:
// . e E + -
// TODO: This function is a mess
func (l *lexer) readNumber() (TokenInfo, error) {
	var builder strings.Builder

	if l.currentChar == '0' && (l.nextChar == 'x' || l.nextChar == 'X') {
		l.next()

		for i := 0; isHex(l.nextChar); i++ {
			if i == 2*8 { // max u64 hex value takes 16 chars
				return TokenInfo{
					Token:  tokens.Integer,
					String: builder.String(),
				}, ErrHexOverflow
			}
			l.next()
			builder.WriteRune(l.currentChar)
		}
		value, err := strconv.ParseUint(builder.String(), 16, 64)
		if err != nil {
			return TokenInfo{
				Token:  tokens.Integer,
				String: builder.String(),
			}, ErrHexOverflow
		}
		return TokenInfo{
			Token:   tokens.Integer,
			String:  builder.String(),
			Integer: value,
		}, nil
	} else if l.currentChar == '0' && isOctal(l.nextChar) {
		for i := 0; isOctal(l.nextChar); i++ {
			if i == 22 { // max u64 oct value takes 22 chars
				return TokenInfo{
					Token:  tokens.Integer,
					String: builder.String(),
				}, ErrOctOverflow
			}
			l.next()
			builder.WriteRune(l.currentChar)
		}
		value, err := strconv.ParseUint(builder.String(), 8, 64)
		if err != nil {
			return TokenInfo{
				Token:  tokens.Integer,
				String: builder.String(),
			}, ErrOctOverflow
		}
		return TokenInfo{
			Token:   tokens.Integer,
			String:  builder.String(),
			Integer: value,
		}, nil
	} else {
		// Read int/float first and then parse integer exponent if present
		builder.WriteRune(l.currentChar)

		var hasDot bool
		for isDigit(l.nextChar) || l.nextChar == '.' {
			// Detect double point
			if l.nextChar == '.' {
				if hasDot {
					return TokenInfo{
						Token:  tokens.Integer,
						String: builder.String(),
					}, ErrFloatFormat
				} else {
					hasDot = true
				}
			}

			l.next()
			builder.WriteRune(l.currentChar)
		}

		// Optionally read exponent
		if isExponent(l.nextChar) {
			l.next()

			sign := l.nextChar == '-'
			if l.nextChar == '+' || l.nextChar == '-' {
				l.next()
			}

			var expBuilder strings.Builder
			for isDigit(l.nextChar) {
				l.next()
				expBuilder.WriteRune(l.currentChar)
			}
			if l.nextChar == '.' || isExponent(l.nextChar) {
				return TokenInfo{
					Token: tokens.Float,
				}, ErrFloatFormat
			}

			value, err := strconv.ParseFloat(builder.String(), 10)
			if err != nil {
				return TokenInfo{
					Token:  tokens.Float,
					String: builder.String(),
				}, ErrFloatFormat
			}

			exponent, err := strconv.ParseUint(expBuilder.String(), 10, 64)
			if err != nil {
				return TokenInfo{
					Token: tokens.Float,
				}, ErrFloatFormat
			}
			if sign {
				exponent = -exponent
			}

			return TokenInfo{
				Token:  tokens.Float,
				String: builder.String(),
				Float:  value * math.Pow10(int(exponent)),
			}, nil
		}

		if hasDot {
			value, err := strconv.ParseFloat(builder.String(), 10)
			if err != nil {
				return TokenInfo{
					Token:  tokens.Float,
					String: builder.String(),
				}, ErrFloatFormat
			}
			return TokenInfo{
				Token:  tokens.Float,
				String: builder.String(),
				Float:  value,
			}, nil
		}

		value, err := strconv.ParseUint(builder.String(), 10, 64)
		if err != nil {
			return TokenInfo{
				Token:  tokens.Integer,
				String: builder.String(),
			}, ErrDecOverflow
		}
		return TokenInfo{
			Token:   tokens.Integer,
			String:  builder.String(),
			Integer: value,
		}, nil
	}

	return TokenInfo{}, nil
}
