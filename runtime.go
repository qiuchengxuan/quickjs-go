package quickjs

//#include "ffi.h"
import "C"
import (
	"runtime"
	"sync/atomic"

	_ "github.com/qiuchengxuan/go-quickjs/libquickjs"
)

type Runtime struct {
	raw  *C.JSRuntime
	free atomic.Bool
}

func (r *Runtime) Free() {
	if !r.free.Swap(true) {
		C.JS_FreeRuntime(r.raw)
	}
}

func NewRuntime() *Runtime {
	retval := &Runtime{raw: C.JS_NewRuntime()}
	if !globalConfig.ManualFree {
		runtime.SetFinalizer(retval, func(r *Runtime) { r.Free() })
	}
	return retval
}
