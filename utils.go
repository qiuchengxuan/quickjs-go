package quickjs

//#include "ffi.h"
import "C"
import "unsafe"

//go:inline
func strPtr(text string) *C.char {
	return (*C.char)(unsafe.Pointer(unsafe.StringData(text)))
}

//go:inline
func strlen(text string) C.size_t { return C.size_t(len(text)) }

//go:inline
func bytesPtr(bytes []byte) *C.uint8_t {
	return (*C.uint8_t)(unsafe.Pointer(&bytes[0]))
}

//go:inline
func slicePtr[T any](slice []T) *C.uint8_t {
	return (*C.uint8_t)(unsafe.Pointer(&slice[0]))
}

//go:inline
func sliceSize[T any](slice []T) C.size_t {
	var t T
	return C.size_t(len(slice)) * C.size_t(unsafe.Sizeof(t))
}

func assert0(value C.int) {
	if value != 0 {
		panic("Assert fail")
	}
}
