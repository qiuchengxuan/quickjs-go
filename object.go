package quickjs

//#include "ffi.h"
import "C"

type ObjectKind uint8

const (
	KindArray ObjectKind = iota
	KindArrayBuffer
	KindInt8Array
	KindInt16Array
	KindInt32Array
	KindUint8Array
	KindUint16Array
	KindUint32Array
	KindFloat32Array
	KindFloat64Array
	KindMap
	KindSet
	KindDate
	KindUnknown
)

var builtinObjects = [13]string{
	"Array", "ArrayBuffer",
	"Int8Array", "Int16Array", "Int32Array",
	"Uint8Array", "Uint16Array", "Uint32Array",
	"Float32Array", "Float64Array",
	"Map", "Set", "Date",
}

type Object struct{ Value }

func (o Object) instanceOf(value Value) bool {
	return C.JS_IsInstanceOf(o.context.raw, o.raw, value.raw) == 1
}

func (o Object) Kind() ObjectKind {
	if C.JS_IsArray(o.context.raw, o.raw) == 1 {
		return KindArray
	}
	globals := Value{o.context, C.JS_GetGlobalObject(o.context.raw)}
	retval := KindUnknown
	for i, name := range builtinObjects[1:] {
		if value := globals.GetProperty(name); o.instanceOf(value) {
			retval = ObjectKind(i + 1)
			break
		}
	}
	C.JS_FreeValue(o.context.raw, globals.raw)
	return retval
}

func (o Object) ToPrimitive() any {
	switch o.Kind() {
	case KindArray:
		return Array(o).ToPrimitive()
	case KindArrayBuffer:
		return ArrayBuffer(o).ToPrimitive()
	case KindInt8Array:
		return TypedArray[int8](o).ToPrimitive()
	case KindInt16Array:
		return TypedArray[int16](o).ToPrimitive()
	case KindInt32Array:
		return TypedArray[int32](o).ToPrimitive()
	case KindUint8Array:
		return TypedArray[uint8](o).ToPrimitive()
	case KindUint16Array:
		return TypedArray[uint16](o).ToPrimitive()
	case KindUint32Array:
		return TypedArray[uint32](o).ToPrimitive()
	case KindFloat32Array:
		return TypedArray[float32](o).ToPrimitive()
	case KindFloat64Array:
		return TypedArray[float64](o).ToPrimitive()
	case KindMap:
		return Map{o}.ToPrimitive()
	case KindSet:
		return Set{o}.ToPrimitive()
	case KindDate:
		return Date(o).ToPrimitive()
	default:
		return NonPrimitive{}
	}
}

func (o Object) call(name string, args ...C.JSValue) C.JSValue {
	fn := o.GetProperty(name)
	var arg0 *C.JSValue
	if len(args) > 0 {
		arg0 = &args[0]
	}
	value := C.JS_Call(o.context.raw, fn.raw, o.raw, C.int(len(args)), arg0)
	fn.free()
	return value
}
