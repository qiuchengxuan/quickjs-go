package quickjs

//#include "ffi.h"
import "C"
import "time"

type Date struct{ Object }

func (d Date) ToPrimitive() time.Time {
	retval, _ := time.Parse("2006-01-02T15:04:051Z", d.String())
	return retval
}
