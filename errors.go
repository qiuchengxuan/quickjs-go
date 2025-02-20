package quickjs

//#include "ffi.h"
import "C"
import (
	"fmt"
	"strings"
)

func (c *Context) getException() error {
	value := Value{c, C.JS_GetException(c.raw)}
	cause := value.String()
	stack, _ := value.Object().GetProperty("stack")
	if stack.Type() == TypeUndefined {
		return &Error{Cause: cause}
	}
	err := &Error{Cause: cause, Stack: stack.String()}
	C.JS_FreeValue(c.raw, value.raw)
	return err
}

func (c *Context) checkException(value C.JSValue) error {
	if C.JS_IsException(value) == 1 {
		return c.getException()
	}
	return nil
}

func (c *Context) assert(value C.JSValue) C.JSValue {
	if err := c.checkException(value); err != nil {
		panic(err)
	}
	return value
}

func (c *Context) ThrowInternalError(format string, args ...any) C.JSValue {
	var buf strings.Builder
	fmt.Fprintf(&buf, format, args...)
	buf.WriteByte(0)
	return C.ThrowInternalError(c.raw, strPtr(buf.String()))
}
