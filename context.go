package quickjs

//#include "ffi.h"
import "C"
import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

type Context struct {
	runtime     *Runtime
	raw         *C.JSContext
	global      C.JSValue
	evalRet     C.JSValue
	goValues    map[uintptr]any
	objectKinds map[C.JSValue]ObjectKind
	free        atomic.Bool
}

func (c *Context) addGoObject(value any) C.JSValue {
	jsObject := C.JS_NewObjectClass(c.raw, C.int(c.runtime.goObject))
	c.goValues[(uintptr)(C.JS_ValuePtr(jsObject))] = value
	data := goObjectData{value, c}
	dataPtr := C.malloc(C.size_t(unsafe.Sizeof(data)))
	*(*goObjectData)(dataPtr) = data
	C.JS_SetOpaque(jsObject, dataPtr)
	return jsObject
}

func (c *Context) addNaiveFunc(fn NaiveFunc) C.JSValue {
	callback := func(_ *C.JSContext, _ C.JSValueConst, args []C.JSValueConst) C.JSValueConst {
		goArgs := make([]any, len(args))
		for i, arg := range args {
			goArgs[i] = Value{c, arg}.ToNative()
		}
		retval, err := fn(goArgs...)
		if err != nil {
			return c.ThrowInternalError("%s", err)
		}
		return c.toJsValue(retval)
	}
	return c.addGoObject(callback)
}

func (c *Context) GlobalObject() Object {
	return Value{c, c.global}.Object()
}

func (c *Context) Compile(code string) (ByteCode, error) {
	codePtr := strPtr(code + "\x00")
	filename := "<input>\x00"
	flags := C.int(C.JS_EVAL_TYPE_GLOBAL | C.JS_EVAL_FLAG_COMPILE_ONLY)
	if C.JS_DetectModule(codePtr, strlen(code)) != 0 {
		flags |= C.JS_EVAL_TYPE_MODULE
	}
	jsValue := C.JS_Eval(c.raw, codePtr, strlen(code), strPtr(filename), flags)
	if err := c.checkException(jsValue); err != nil {
		return nil, err
	}
	var size C.size_t
	pointer := C.JS_WriteObject(c.raw, &size, jsValue, C.JS_WRITE_OBJ_BYTECODE)
	C.JS_FreeValue(c.raw, jsValue)
	if int(size) <= 0 {
		return nil, c.getException()
	}
	byteCode := C.GoBytes(unsafe.Pointer(pointer), C.int(size))
	C.js_free(c.raw, unsafe.Pointer(pointer))
	return byteCode, nil
}

func (c *Context) eval(code string) (C.JSValue, error) {
	codePtr := strPtr(code + "\x00")
	filename := "<input>\x00"
	flags := C.int(C.JS_EVAL_TYPE_GLOBAL)
	if C.JS_DetectModule(codePtr, strlen(code)) != 0 {
		flags |= C.JS_EVAL_TYPE_MODULE
	}
	jsValue := C.JS_Eval(c.raw, codePtr, strlen(code), strPtr(filename), flags)
	if err := c.checkException(jsValue); err != nil {
		return null, err
	}
	return jsValue, nil
}

// Return value must be consumed immediately before next Eval or EvalBinary
func (c *Context) Eval(code string) (Value, error) {
	C.JS_FreeValue(c.raw, c.evalRet)
	value, err := c.eval(code)
	c.evalRet = value
	return Value{c, value}, err
}

// Return value must be consumed immediately before next Eval or EvalBinary
func (c *Context) EvalBinary(byteCode ByteCode) (Value, error) {
	flags := C.int(C.JS_READ_OBJ_BYTECODE)
	object := C.JS_ReadObject(c.raw, bytesPtr(byteCode), C.size_t(len(byteCode)), flags)
	retval := c.assert(C.JS_EvalFunction(c.raw, c.assert(object)))
	C.JS_FreeValue(c.raw, c.evalRet)
	c.evalRet = retval
	return Value{c, retval}, nil
}

// Free context manually
func (c *Context) Free() {
	if c.free.Swap(true) {
		return
	}
	C.JS_FreeValue(c.raw, c.global)
	C.JS_FreeValue(c.raw, c.evalRet)
	C.JS_FreeContext(c.raw)
	c.runtime.Free()
}

type ContextGuard struct{ context *Context }

// Manipulate Context with os thread locked
func (g ContextGuard) With(fn func(*Context)) {
	// Reason unknown, without locking os thread will cause quickjs throw strange exception
	runtime.LockOSThread()
	fn(g.context)
	runtime.UnlockOSThread()
}

// NOTE: unsafe
func (g ContextGuard) Unwrap() *Context { return g.context }

func (g ContextGuard) Free() { g.context.Free() }

func (r *Runtime) NewContext() ContextGuard {
	r.refCount.Add(1)
	C.js_std_init_handlers(r.raw)

	jsContext := C.JS_NewContext(r.raw)
	C.JS_AddIntrinsicBigFloat(jsContext)
	C.JS_AddIntrinsicBigDecimal(jsContext)
	C.JS_AddIntrinsicOperators(jsContext)
	C.JS_EnableBignumExt(jsContext, C.int(1))

	object := C.JS_GetGlobalObject(jsContext)
	goValues := make(map[uintptr]any)
	context := &Context{runtime: r, raw: jsContext, global: object, goValues: goValues}
	proto := C.JS_NewObject(jsContext)
	C.JS_SetClassProto(jsContext, r.goObject, proto)
	objectKinds := make(map[C.JSValue]ObjectKind, KindDate+1)
	for i, name := range builtinKinds {
		jsValue, _ := context.GlobalObject().GetProperty(name)
		objectKinds[jsValue.raw] = ObjectKind(i + 1)
	}
	context.objectKinds = objectKinds
	if !globalConfig.ManualFree {
		runtime.SetFinalizer(context, func(c *Context) { c.Free() })
	}
	return ContextGuard{context}
}
