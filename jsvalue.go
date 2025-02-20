package quickjs

//#include "ffi.h"
import "C"
import (
	"bytes"
	"encoding/json"
	"math"
	"math/big"
	"reflect"
	"unsafe"
)

var null = C.JS_Null()

func (c *Context) toJsValue(value any) C.JSValue {
	switch value := value.(type) {
	case bool:
		intValue := 0
		if value {
			intValue = 1
		}
		return C.JS_NewBool(c.raw, C.int(intValue))
	case int8:
		return C.JS_NewInt32(c.raw, C.int32_t(value))
	case int16:
		return C.JS_NewInt32(c.raw, C.int32_t(value))
	case int32:
		return C.JS_NewInt32(c.raw, C.int32_t(value))
	case int64:
		return C.JS_NewInt64(c.raw, C.int64_t(value))
	case int:
		return C.JS_NewInt64(c.raw, C.int64_t(value))
	case uint8:
		return C.JS_NewInt32(c.raw, C.int32_t(value))
	case uint16:
		return C.JS_NewInt32(c.raw, C.int32_t(value))
	case uint32:
		return C.JS_NewInt64(c.raw, C.int64_t(value))
	case uint64:
		if value <= math.MaxInt64 {
			return C.JS_NewInt64(c.raw, C.int64_t(value))
		}
		return C.JS_NewFloat64(c.raw, C.double(value))
	case uint:
		if value <= math.MaxInt64 {
			return C.JS_NewInt64(c.raw, C.int64_t(value))
		}
		return C.JS_NewFloat64(c.raw, C.double(value))
	case float32:
		return C.JS_NewFloat64(c.raw, C.double(value))
	case float64:
		return C.JS_NewFloat64(c.raw, C.double(value))
	case *big.Int:
		return c.toJsValue(*value)
	case big.Int:
		arg := C.JS_NewString(c.raw, strPtr(value.String()+"\x00"))
		bigInt, _ := c.GlobalObject().GetProperty("BigInt")
		retval := bigInt.Object().call(null, 1, &arg)
		C.JS_FreeValue(c.raw, arg)
		return retval
	case string:
		newStr := value + "\x00"
		return C.JS_NewString(c.raw, strPtr(newStr))
	case []int8:
		return newTypedArray(c, value, typedArrayInt8)
	case []uint8: // also []byte
		return newTypedArray(c, value, typedArrayUint8)
	case []int16:
		return newTypedArray(c, value, typedArrayInt16)
	case []uint16:
		return newTypedArray(c, value, typedArrayUint16)
	case []int32:
		return newTypedArray(c, value, typedArrayInt32)
	case []uint32:
		return newTypedArray(c, value, typedArrayUint32)
	case []float32:
		return newTypedArray(c, value, typedArrayFloat32)
	case []float64:
		return newTypedArray(c, value, typedArrayFloat64)
	case []any:
		array := Value{c, C.JS_NewArray(c.raw)}.Object().Array()
		for i, item := range value {
			array.Set(i, item)
		}
		return array.raw
	case map[string]any:
		object := Value{c, C.JS_NewObject(c.raw)}.Object()
		for key, value := range value {
			object.setProperty(key, c.toJsValue(value))
		}
		return object.raw
	case NaiveFunc:
		return c.addNaiveFunc(value)
	case json.Marshaler:
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(value); err != nil {
			return null
		}
		buf.WriteByte(0)
		data := buf.Bytes()
		dataPtr := (*C.char)(unsafe.Pointer(&data[0]))
		return C.JS_ParseJSON(c.raw, dataPtr, sliceSize(data)-1, nil)
	default:
		if value == Undefined {
			return C.JS_Undefined()
		}
		if value == nil {
			return null
		}
		valueOf := reflect.ValueOf(value)
		switch valueOf.Kind() {
		case reflect.Map:
			class, _ := c.GlobalObject().GetProperty("Map")
			items := make([]any, 0, valueOf.Len())
			iter := valueOf.MapRange()
			for iter.Next() {
				item := []any{iter.Key().Interface(), iter.Value().Interface()}
				items = append(items, item)
			}
			jsValue := c.toJsValue(items)
			jsMap := C.JS_CallConstructor(c.raw, class.raw, 1, &jsValue)
			C.JS_FreeValue(c.raw, jsValue)
			return jsMap
		case reflect.Array, reflect.Slice:
			array := Value{c, C.JS_NewArray(c.raw)}.Object().Array()
			for i := 0; i < valueOf.Len(); i++ {
				array.Set(i, valueOf.Index(i).Interface())
			}
			return array.raw
		default:
			return C.JS_Undefined()
		}
	}
}
