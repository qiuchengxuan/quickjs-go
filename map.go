package quickjs

//#include "ffi.h"
import "C"

type Map struct{ Object }

func (m Map) Size() int {
	return m.GetProperty("size").ToPrimitive().(int)
}

func (m Map) foreach(fn func(key, value C.JSValueConst)) {
	wrapped := func(_ *C.JSContext, _ C.JSValueConst, args []C.JSValueConst) C.JSValueConst {
		fn(args[1], args[0])
		return C.JS_Null()
	}
	value := m.context.addCallback(&wrapped)
	if err := m.context.checkException(m.call("forEach", value)); err != nil {
		panic(err)
	}
	C.JS_FreeValue(m.context.raw, value)
}

func (m Map) ToPrimitive() map[string]any {
	retval := make(map[string]any, m.Size())
	m.foreach(func(jsKey, value C.JSValueConst) {
		key := Value{m.context, jsKey}.String()
		retval[key] = Value{m.context, value}.ToPrimitive()
	})
	return retval
}
