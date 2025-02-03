package quickjs

//#include "ffi.h"
import "C"
import (
	"math/big"
)

type Value struct {
	context *Context
	raw     C.JSValue
}

func (v Value) String() string {
	cString := C.JS_ToCString(v.context.raw, v.raw)
	retval := C.GoString(cString)
	C.JS_FreeCString(v.context.raw, cString)
	return retval
}

func (v Value) getProperty(name string) C.JSValue {
	return C.JS_GetPropertyStr(v.context.raw, v.raw, strPtr(name+"\x00"))
}

func (v Value) GetProperty(name string) Value {
	return Value{v.context, v.getProperty(name)}
}

func (v Value) setProperty(name string, value C.JSValue) {
	C.JS_SetPropertyStr(v.context.raw, v.raw, strPtr(name+"\x00"), value)
}

func (v Value) Type() Type {
	switch {
	case C.JS_IsNull(v.raw) == 1:
		return TypeNull
	case C.JS_IsUndefined(v.raw) == 1:
		return TypeUndefined
	case C.JS_IsBool(v.raw) == 1:
		return TypeBool
	case C.JS_IsNumber(v.raw) == 1:
		return TypeNumber
	case C.JS_IsBigInt(v.context.raw, v.raw) == 1:
		return TypeBigInt
	case C.JS_IsString(v.raw) == 1:
		return TypeString
	case C.JS_IsSymbol(v.raw) == 1:
		return TypeSymbol
	case C.JS_IsObject(v.raw) == 1:
		return TypeObject
	default:
		return TypeNonPrimitive
	}
}

func (v Value) free() { C.JS_FreeValue(v.context.raw, v.raw) }

func (v Value) ToPrimitive() any {
	switch v.Type() {
	case TypeNull, TypeUndefined:
		return nil
	case TypeBool:
		return C.JS_ToBool(v.context.raw, v.raw) == 1
	case TypeNumber:
		if C.JS_IsInt(v.raw) == 1 {
			var retval C.int32_t
			C.JS_ToInt32(v.context.raw, &retval, v.raw)
			return int(retval)
		}
		var retval C.double
		C.JS_ToFloat64(v.context.raw, &retval, v.raw)
		return float64(retval)
	case TypeBigInt:
		var value C.int64_t
		assert0(C.JS_ToBigInt64(v.context.raw, &value, v.raw))
		if value != 0 { // not overflow
			return big.NewInt(int64(value))
		}
		var retval big.Int
		retval.SetString(v.String(), 10)
		return retval
	case TypeString:
		return v.String()
	case TypeObject:
		return Object{v}.ToPrimitive()
	default:
		return NonPrimitive{}
	}
}

// []uint8 not supported due to conflict with []byte
func (v Value) SetProperty(name string, value any) {
	v.setProperty(name, v.context.toJsValue(value))
}
