package quickjs

//#include "ffi.h"
import "C"
import (
	"errors"
	"math"
	"runtime"
	"sync/atomic"
	"unsafe"
)

type Context struct {
	runtime   *Runtime
	raw       *C.JSContext
	filename  C.int
	global    C.JSValue
	proxy     C.JSValue
	makeProxy C.JSValue
	evalRet   C.JSValue
	funcs     []callback
	free      atomic.Bool
}

func (c *Context) getException() error {
	exception := C.JS_GetException(c.raw)
	return errors.New(Value{c, exception}.String())
}

func (c *Context) checkException(value C.JSValue) error {
	if C.JS_IsException(value) == 1 {
		return c.getException()
	}
	return nil
}

func (c *Context) addCallback(callback *callback) C.JSValue {
	pointer := math.Float64frombits(uint64(uintptr(unsafe.Pointer(callback))))
	handler := C.JS_NewFloat64(c.raw, C.double(pointer))
	args := []C.JSValue{c.proxy, handler}
	return C.JS_Call(c.raw, c.makeProxy, C.JS_Null(), C.int(len(args)), &args[0])
}

func (c *Context) addNaiveFunc(fn NaiveFunc) C.JSValue {
	callback := func(_ *C.JSContext, _ C.JSValueConst, args []C.JSValueConst) C.JSValueConst {
		goArgs := make([]any, len(args))
		for i, arg := range args {
			goArgs[i] = Value{c, arg}.ToPrimitive()
		}
		return c.toJsValue(fn(goArgs...))
	}
	c.funcs = append(c.funcs, callback) // avoid GC
	return c.addCallback(&c.funcs[len(c.funcs)-1])
}

func (c *Context) GlobalObject() Value {
	return Value{c, c.global}
}

func (c *Context) Compile(code string) (ByteCode, error) {
	codePtr := strPtr(code + "\x00")
	filename := "<input>\x00"
	flags := C.int(C.JS_EVAL_TYPE_GLOBAL | C.JS_EVAL_FLAG_COMPILE_ONLY)
	if C.JS_DetectModule(codePtr, strlen(code)) != 0 {
		flags |= C.JS_EVAL_TYPE_MODULE
	}
	value := C.JS_Eval(c.raw, codePtr, strlen(code), strPtr(filename), flags)
	if err := c.checkException(value); err != nil {
		return nil, err
	}
	var size C.size_t
	pointer := C.JS_WriteObject(c.raw, &size, value, C.JS_WRITE_OBJ_BYTECODE)
	C.JS_FreeValue(c.raw, value)
	if int(size) <= 0 {
		return nil, c.getException()
	}
	byteCode := make(ByteCode, int(size))
	copy(byteCode, C.GoBytes(unsafe.Pointer(pointer), C.int(size)))
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
	value := C.JS_Eval(c.raw, codePtr, strlen(code), strPtr(filename), flags)
	if err := c.checkException(value); err != nil {
		return C.JS_Null(), err
	}
	return value, nil
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
	retval := C.JS_ReadObject(c.raw, bytesPtr(byteCode), C.size_t(len(byteCode)), flags)
	if err := c.checkException(retval); err != nil {
		return Value{c, C.JS_Null()}, err
	}
	retval = C.JS_EvalFunction(c.raw, retval)
	if err := c.checkException(retval); err != nil {
		return Value{c, C.JS_Null()}, err
	}
	C.JS_FreeValue(c.raw, c.evalRet)
	c.evalRet = retval
	return Value{c, retval}, nil
}

func (c *Context) Free() {
	if c.free.Swap(true) {
		return
	}
	C.JS_FreeValue(c.raw, c.global)
	C.JS_FreeValue(c.raw, c.proxy)
	C.JS_FreeValue(c.raw, c.evalRet)
	C.JS_FreeContext(c.raw)
}

func (r *Runtime) NewContext() *Context {
	C.js_std_init_handlers(r.raw)

	jsContext := C.JS_NewContext(r.raw)
	C.JS_AddIntrinsicBigFloat(jsContext)
	C.JS_AddIntrinsicBigDecimal(jsContext)
	C.JS_AddIntrinsicOperators(jsContext)
	C.JS_EnableBignumExt(jsContext, C.int(1))

	fn := (*C.JSCFunction)(unsafe.Pointer(C.proxyCall))
	object := C.JS_GetGlobalObject(jsContext)
	proxy := C.JS_NewCFunction(jsContext, fn, nil, C.int(0))
	context := &Context{runtime: r, raw: jsContext, global: object, proxy: proxy}
	if !globalConfig.ManualFree {
		runtime.SetFinalizer(context, func(c *Context) { c.Free() })
	}
	makeProxy := "(proxy, handler) => function() { return proxy.call(handler, ...arguments) }"
	context.makeProxy, _ = context.eval(makeProxy)
	return context
}
