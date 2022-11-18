package sqvm

type ObjectType int

const (
	TypeNull ObjectType = iota
	TypeInteger
	TypeFloat
	TypeString
	TypeTable
	TypeArray
	TypeUserData
	TypeClosure
	TypeNativeClosure
	TypeGenerator
	TypeUserPointer
	TypeBool
	TypeInstance
	TypeClass
	TypeWeakRef
)

type FuncProto struct {
}
