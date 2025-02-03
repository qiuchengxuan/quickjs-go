package quickjs

//#include "ffi.h"
import "C"
import (
	"unsafe"

	"golang.org/x/exp/constraints"
)

type Array struct{ Value }

func (a Array) Len() int {
	return a.GetProperty("length").ToPrimitive().(int)
}

func (a Array) Get(index int) Value {
	value := C.JS_GetPropertyUint32(a.context.raw, a.raw, C.uint32_t(index))
	return Value{a.context, value}
}

func (a Array) ToPrimitive() []any {
	length := a.Len()
	retval := make([]any, 0, length)
	for i := 0; i < length; i++ {
		retval = append(retval, a.Get(i).ToPrimitive())
	}
	return retval
}

type ArrayBuffer struct{ Value }

func (b ArrayBuffer) Len() int {
	return b.GetProperty("byteLength").ToPrimitive().(int)
}

func (b ArrayBuffer) ToPrimitive() []byte {
	size := C.size_t(b.Len())
	out := C.JS_GetArrayBuffer(b.context.raw, &size, b.raw)
	return C.GoBytes(unsafe.Pointer(out), C.int(size))
}

type Number interface {
	constraints.Integer | constraints.Float
}

type TypedArray[T Number] struct{ Value }

func (a TypedArray[T]) Len() int {
	return a.GetProperty("length").ToPrimitive().(int)
}

func (a TypedArray[T]) ToPrimitive() []T {
	buf := C.JS_GetTypedArrayBuffer(a.context.raw, a.raw, nil, nil, nil)
	bytes := ArrayBuffer{Value{a.context, buf}}.ToPrimitive()
	var t T
	sizeOf := int(unsafe.Sizeof(t))
	return unsafe.Slice((*T)(unsafe.Pointer(&bytes[0])), len(bytes)/sizeOf)
}
