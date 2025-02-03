package quickjs

type ByteCode []byte

type NonPrimitive struct{}

type NaiveFunc = func(...any) any

type Type uint8

const (
	TypeNull Type = iota
	TypeUndefined
	TypeBool
	TypeNumber
	TypeBigInt
	TypeString
	TypeSymbol
	TypeObject
	TypeNonPrimitive
)
