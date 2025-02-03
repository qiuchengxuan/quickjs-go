package quickjs

//#include "ffi.h"
import "C"

type Set struct{ Object }

func (s Set) Size() int {
	return s.GetProperty("size").ToPrimitive().(int)
}

func (s Set) foreach(fn func(value C.JSValueConst)) {
	wrapped := func(_ *C.JSContext, _ C.JSValueConst, args []C.JSValueConst) C.JSValueConst {
		fn(args[0])
		return C.JS_Null()
	}
	s.call("forEach", s.context.addCallback(&wrapped))
}

func (s Set) ToPrimitive() []any {
	retval := make([]any, 0, s.Size())
	s.foreach(func(value C.JSValueConst) {
		retval = append(retval, Value{s.context, value}.ToPrimitive())
	})
	return retval
}
