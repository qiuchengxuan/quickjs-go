package quickjs

//#include "ffi.h"
import "C"
import (
	"math"
	"unsafe"
)

type callback = func(*C.JSContext, C.JSValueConst, []C.JSValueConst) C.JSValueConst

//export proxyCall
func proxyCall(raw *C.JSContext, this C.JSValueConst, argc C.int, argv *C.JSValueConst) C.JSValue {
	refs := unsafe.Slice(argv, argc)
	var double C.double
	C.JS_ToFloat64(raw, &double, this)
	callback := *(*callback)(unsafe.Pointer((uintptr)(math.Float64bits(float64(double)))))
	return callback(raw, this, refs)
}
