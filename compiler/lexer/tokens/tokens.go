package tokens

import (
	"fmt"
)

type Token uint

const Undefined Token = 0

const (
	Identifier Token = iota + 256 // skip printable characters
	StringLiteral
	Integer
	Float
	Base
	Delete
	Equal
	NotEqual
	LessEqual
	GreaterEqual
	Switch
	Arrow
	And
	Or
	If
	Else
	While
	Break
	For
	Do
	Null
	Foreach
	In
	NewSlot
	Modulo
	Local
	Clone
	Function
	Return
	Typeof
	Negate // unary minus
	PlusEqual
	MinusEqual
	Continue
	Yield
	Try
	Catch
	Throw
	ShiftLeft
	ShiftRight
	Resume
	DoubleColon
	Case
	Default
	This
	Increase
	Decrease
	ThreeWayCompare
	UnsignedShiftRight
	Class
	Extends
	Constructor
	InstanceOf
	VarParams
	Line
	File
	True
	False
	MultiplyEqual
	DivideEqual
	ModuloEqual
	AttributeOpen
	AttributeClose
	Static
	Enum
	Const
	RawCall
)

func (t Token) String() string {
	if t == Undefined {
		return "tokens.Undefined"
	}

	if t < Identifier {
		return fmt.Sprintf("tokens.Printable(%q)", rune(t))
	}

	if t > RawCall {
		return fmt.Sprintf("tokens.Token(%d)", uint(t))
	}

	return []string{
		"Identifier", "StringLiteral", "Integer", "Float", "Base", "Delete", "Equal",
		"NotEqual", "LessEqual", "GreaterEqual", "Switch", "Arrow", "And", "Or", "If",
		"Else", "While", "Break", "For", "Do", "Null", "Foreach", "In", "NewSlot",
		"Modulo", "Local", "Clone", "Function", "Return", "Typeof", "Negate",
		"PlusEqual", "MinusEqual", "Continue", "Yield", "Try", "Catch", "Throw",
		"ShiftLeft", "ShiftRight", "Resume", "DoubleColon", "Case", "Default", "This",
		"Increase", "Decrease", "ThreeWayCompare", "UnsignedShiftRight", "Class",
		"Extends", "Constructor", "InstanceOf", "VarParams", "Line", "File", "True",
		"False", "MultiplyEqual", "DivideEqual", "ModuloEqual", "AttributeOpen",
		"AttributeClose", "Static", "Enum", "Const", "RawCall",
	}[t - Identifier]
}
