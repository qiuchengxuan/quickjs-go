package quickjs

//#include "ffi.h"
import "C"
import (
	"math/big"
	"strconv"
)

const (
	tagBigInt    = -10
	tagSymbol    = -8
	tagString    = -7
	tagObject    = -1
	tagInt       = 0
	tagBool      = 1
	tagNull      = 2
	tagUndefined = 3
	tagFloat64   = 7
)

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
	TypeNotNative
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

func (v Value) Type() Type {
	switch C.JS_ValueTag(v.raw) {
	case tagNull:
		return TypeNull
	case tagUndefined:
		return TypeUndefined
	case tagBool:
		return TypeBool
	case tagInt, tagFloat64:
		return TypeNumber
	case tagBigInt:
		return TypeBigInt
	case tagString:
		return TypeString
	case tagSymbol:
		return TypeSymbol
	case tagObject:
		return TypeObject
	default:
		return TypeNotNative
	}
}

// Be aware that number could be int or double
// nil return indicates value is null or undefined or not primitive
func (v Value) ToPrimitive() any {
	switch v.Type() {
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
	case TypeString:
		return v.String()
	default:
		return nil
	}
}

func (v Value) JSONify() string {
	jsValue := C.JS_JSONStringify(v.context.raw, v.raw, null, null)
	output := Value{v.context, jsValue}.ToNative().(string)
	C.JS_FreeValue(v.context.raw, jsValue)
	return output
}

// Be aware that number could be int or double,
// BigInt could be int or big.Int,
// Plain object will be converted to map[string]any or []any
// Map will be converted to map[any]any
func (v Value) ToNative() any {
	switch v.Type() {
	case TypeNull:
		return nil
	case TypeUndefined:
		return Undefined
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
		if value > 0 { // 0 is possibly overflow, and negative value might be unsigned
			return int(value)
		}
		strVal := v.String()
		if len(strVal) < 19 {
			retval, _ := strconv.Atoi(strVal)
			return retval
		}
		var retval big.Int
		retval.SetString(strVal, 10)
		return retval
	case TypeString:
		return v.String()
	case TypeObject:
		return v.Object().ToNative()
	default:
		return NotNative{}
	}
}

func (v Value) free() {
	C.JS_FreeValue(v.context.raw, v.raw)
}
