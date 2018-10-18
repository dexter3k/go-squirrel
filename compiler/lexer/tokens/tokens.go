package tokens

import ()

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
