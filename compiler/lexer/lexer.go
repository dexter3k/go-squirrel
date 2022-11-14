package lexer

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/dexter3k/go-squirrel/compiler/lexer/tokens"
)

var (
	ErrUnknownToken     = fmt.Errorf("Unknown token")
	ErrHexOverflow      = fmt.Errorf("Hexadecimal number is too big")
	ErrOctOverflow      = fmt.Errorf("Octal number is too big")
	ErrDecOverflow      = fmt.Errorf("Decimal number is too big")
	ErrFloatFormat      = fmt.Errorf("Malformed floating point value")
	ErrUnfinishedString = fmt.Errorf("String is left unterminated")
	ErrBadEscape        = fmt.Errorf("Unknown escape character")
	ErrBadCharacter     = fmt.Errorf("Unknown character")
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
	"__FILE__":    tokens.File,
	"__LINE__":    tokens.Line,
}

type TokenInfo struct {
	Token   tokens.Token
	String  string
	Integer uint64
	Float   float64

	Line   uint
	Column uint
}

type Lexer interface {
	Lex() (TokenInfo, error)
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
	l.next()

	return l
}

func (l *lexer) Lex() (TokenInfo, error) {
	for l.currentChar != 0 {
		switch l.currentChar {
		case '\t', '\r', ' ':
			l.next()
		case '\n':
			l.next()
			return TokenInfo{Token: tokens.Token('\n')}, nil
		case '#':
			l.skipLineComment()
		case '/':
			l.next()
			switch l.currentChar {
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
			l.next()
			if l.currentChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.Equal}, nil
			}
			return TokenInfo{Token: tokens.Token('=')}, nil
		case '<':
			l.next()
			switch l.currentChar {
			case '=':
				l.next()
				if l.currentChar == '>' {
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
			l.next()
			if l.currentChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.GreaterEqual}, nil
			} else if l.currentChar == '>' {
				l.next()
				if l.currentChar == '>' {
					l.next()
					return TokenInfo{Token: tokens.UnsignedShiftRight}, nil
				}
				return TokenInfo{Token: tokens.ShiftRight}, nil
			}
			return TokenInfo{Token: tokens.Token('<')}, nil
		case '!':
			l.next()
			if l.currentChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.NotEqual}, nil
			}
			return TokenInfo{Token: tokens.Token('!')}, nil
		case '@':
			l.next()
			if l.currentChar == '"' {
				l.next()
				return l.readString('"', true)
			}
			return TokenInfo{Token: tokens.Token('@')}, nil
		case '"':
			l.next()
			return l.readString('"', false)
		case '\'':
			l.next()
			return l.readString('\'', false)
		case '{', '}', '(', ')', '[', ']':
			fallthrough
		case ';', ',', '?', '^', '~':
			ch := l.currentChar
			l.next()
			return TokenInfo{Token: tokens.Token(ch)}, nil
		case '.':
			l.next()
			if l.currentChar != '.' {
				return TokenInfo{Token: tokens.Token('.')}, nil
			}
			l.next()
			if l.currentChar != '.' {
				return TokenInfo{String: ".."}, ErrUnknownToken
			}
			l.next()
			return TokenInfo{Token: tokens.VarParams}, nil
		case '&':
			l.next()
			if l.currentChar == '&' {
				l.next()
				return TokenInfo{Token: tokens.And}, nil
			}
			return TokenInfo{Token: tokens.Token('&')}, nil
		case '|':
			l.next()
			if l.currentChar == '|' {
				l.next()
				return TokenInfo{Token: tokens.Or}, nil
			}
			return TokenInfo{Token: tokens.Token('|')}, nil
		case ':':
			l.next()
			if l.currentChar == ':' {
				l.next()
				return TokenInfo{Token: tokens.DoubleColon}, nil
			}
			return TokenInfo{Token: tokens.Token(':')}, nil
		case '%':
			l.next()
			if l.currentChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.ModuloEqual}, nil
			}
			return TokenInfo{Token: tokens.Token('%')}, nil
		case '+':
			l.next()
			if l.currentChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.PlusEqual}, nil
			}
			if l.currentChar == '+' {
				l.next()
				return TokenInfo{Token: tokens.Increase}, nil
			}
			return TokenInfo{Token: tokens.Token('+')}, nil
		case '-':
			l.next()
			if l.currentChar == '=' {
				l.next()
				return TokenInfo{Token: tokens.MinusEqual}, nil
			}
			if l.currentChar == '-' {
				l.next()
				return TokenInfo{Token: tokens.Decrease}, nil
			}
			return TokenInfo{Token: tokens.Token('-')}, nil
		default:
			if isDigit(l.currentChar) {
				return l.readNumber()
			} else if isAlpha(l.currentChar) || l.currentChar == '_' {
				return l.readIdentifier()
			}
			
			return TokenInfo{
				Token:  tokens.Undefined,
				String: string(l.currentChar),
			}, ErrBadCharacter
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
	// l.next() // slash is already skipped
	l.next() // *
	for (l.currentChar != '*' || l.nextChar != '/') && l.currentChar != 0 {
		l.next()
	}
	l.next() // *
	l.next() // /
}

func (l *lexer) readString(delimiter rune, verbatim bool) (TokenInfo, error) {
	var builder strings.Builder
	for {
		for ; l.currentChar != delimiter; l.next() {
			switch l.currentChar {
			case 0:
				return TokenInfo{
					Token:  tokens.StringLiteral,
					String: builder.String(),
				}, ErrUnfinishedString
			case '\n':
				if !verbatim {
					return TokenInfo{
						Token:  tokens.StringLiteral,
						String: builder.String(),
					}, ErrUnfinishedString
				}
				builder.WriteRune('\n')
			case '\\':
				if verbatim {
					builder.WriteRune('\\')
					continue
				}
				l.next()
				switch l.currentChar {
				case 't':
					builder.WriteRune('\t')
				case 'a':
					builder.WriteRune('\a')
				case 'b':
					builder.WriteRune('\b')
				case 'n':
					builder.WriteRune('\n')
				case 'r':
					builder.WriteRune('\r')
				case 'v':
					builder.WriteRune('\v')
				case 'f':
					builder.WriteRune('\f')
				case '0':
					builder.WriteRune(0)
				case '\\':
					builder.WriteRune('\\')
				case '"':
					builder.WriteRune('"')
				case '\'':
					builder.WriteRune('\'')
				default:
					return TokenInfo{
						Token:  tokens.StringLiteral,
						String: builder.String(),
					}, ErrBadEscape
				}
			default:
				builder.WriteRune(l.currentChar)
			}
		}

		if verbatim && l.nextChar == '"' {
			builder.WriteRune('"')
			l.next()
		} else {
			break
		}
	}

	l.next() // skip delimiter

	if delimiter == '\'' {
		if builder.Len() != 1 {
			return TokenInfo{
				Token:  tokens.Integer,
				String: builder.String(),
			}, ErrBadCharacter
		}
		var first rune
		for _, c := range builder.String() {
			first = c
			break
		}
		return TokenInfo{
			Token:   tokens.Integer,
			String:  builder.String(),
			Integer: uint64(first),
		}, nil
	}

	return TokenInfo{Token: tokens.StringLiteral, String: builder.String()}, nil
}

func (l *lexer) readIdentifier() (TokenInfo, error) {
	var builder strings.Builder
	for isAlnum(l.currentChar) || l.currentChar == '_' {
		builder.WriteRune(l.currentChar)
		l.next()
	}

	if token, present := keywords[builder.String()]; present {
		return TokenInfo{Token: token, String: builder.String()}, nil
	}
	return TokenInfo{Token: tokens.Identifier, String: builder.String()}, nil
}

// Integer, Hex, Octal, Float, Scientific
// Hex format:
// 0xFFFFFF, 0Xfffffff, 0xffffff, 0XFFFFFF
// Octal format:
// 07777777
// Scientific format:
// . e E + -
// TODO: This function is a mess
func (l *lexer) readHexPart() (string, error) {
	var builder strings.Builder

	for i := 0; isHex(l.currentChar); i++ {
		if i == 2*8 { // max u64 hex value takes 16 chars
			return builder.String(), ErrHexOverflow
		}

		builder.WriteRune(l.currentChar)
		l.next()
	}

	return builder.String(), nil
}

func (l *lexer) readOctPart() (string, error) {
	var builder strings.Builder

	for i := 0; isOctal(l.currentChar); i++ {
		if i == 3*8 { // max u64 hex value takes 22 chars
			return builder.String(), ErrOctOverflow
		}

		builder.WriteRune(l.currentChar)
		l.next()
	}

	return builder.String(), nil
}

func (l *lexer) readNumber() (TokenInfo, error) {
	if l.currentChar == '0' && (l.nextChar == 'x' || l.nextChar == 'X') {
		l.next()
		l.next()

		str, err := l.readHexPart()
		if err != nil {
			return TokenInfo{
				Token:  tokens.Integer,
				String: str,
			}, err
		}

		val, err := strconv.ParseUint(str, 16, 64)
		if err != nil {
			panic(err) // TODO
		}

		return TokenInfo{
			Token:   tokens.Integer,
			String:  str,
			Integer: val,
		}, nil
	}

	if l.currentChar == '0' && isOctal(l.nextChar) {
		l.next()

		str, err := l.readOctPart()
		if err != nil {
			return TokenInfo{
				Token:  tokens.Integer,
				String: str,
			}, err
		}

		val, err := strconv.ParseUint(str, 8, 64)
		if err != nil {
			panic(err) // TODO
		}

		return TokenInfo{
			Token:   tokens.Integer,
			String:  str,
			Integer: val,
		}, nil
	}

	// Read int/float first and then parse integer exponent if present

	var builder strings.Builder
	builder.WriteRune(l.currentChar)
	l.next()

	var hasDot bool
	for isDigit(l.currentChar) || l.currentChar == '.' {
		// Detect double point
		if l.currentChar == '.' {
			if hasDot {
				return TokenInfo{
					Token:  tokens.Integer,
					String: builder.String(),
				}, ErrFloatFormat
			} else {
				hasDot = true
			}
		}

		builder.WriteRune(l.currentChar)
		l.next()
	}

	// If we have a dot, then parse float right away
	var float float64
	if hasDot {
		var err error
		float, err = strconv.ParseFloat(builder.String(), 10)
		if err != nil {
			return TokenInfo{
				Token:  tokens.Float,
				String: builder.String(),
			}, ErrFloatFormat
		}
	}

	// Optionally read exponent
	if isExponent(l.currentChar) {
		l.next()

		sign := l.currentChar == '-'
		if l.currentChar == '+' || l.currentChar == '-' {
			l.next()
		}

		var expBuilder strings.Builder
		for isDigit(l.currentChar) {
			expBuilder.WriteRune(l.currentChar)
			l.next()
		}
		if l.currentChar == '.' || isExponent(l.currentChar) {
			return TokenInfo{
				Token: tokens.Float,
			}, ErrFloatFormat
		}

		exp, err := strconv.ParseUint(expBuilder.String(), 10, 64)
		if err != nil {
			return TokenInfo{
				Token: tokens.Float,
			}, ErrFloatFormat
		}
		if sign {
			exp = -exp
		}

		return TokenInfo{
			Token:  tokens.Float,
			String: builder.String(),
			Float:  float * math.Pow10(int(int64(exp))),
		}, nil
	}

	if hasDot {
		return TokenInfo{
			Token:  tokens.Float,
			String: builder.String(),
			Float:  float,
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
