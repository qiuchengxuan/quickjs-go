package quickjs

import (
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainObject(t *testing.T) {
	NewRuntime().NewContext().With(func(context *Context) {
		retval, err := context.Eval("new Object([1, 2])")
		assert.NoError(t, err)
		expected := []any{1, 2}
		assert.Equal(t, expected, retval.ToNative())

		context.GlobalObject().SetProperty("plain", expected)
		retval, err = context.GlobalObject().GetProperty("plain")
		assert.NoError(t, err)
		assert.Equal(t, expected, retval.ToNative())
	})
	NewRuntime().NewContext().With(func(context *Context) {
		retval, err := context.Eval("new Object({a: 1, b: {c: 2}})")
		assert.NoError(t, err)
		expected := map[string]any{"a": 1, "b": map[string]any{"c": 2}}
		assert.Equal(t, expected, retval.ToNative())

		context.GlobalObject().SetProperty("plain2", expected)
		retval, err = context.GlobalObject().GetProperty("plain2")
		assert.NoError(t, err)
		assert.Equal(t, expected, retval.ToNative())
	})
}

func TestSetProperty(t *testing.T) {
	NewRuntime().NewContext().With(func(context *Context) {
		object := context.GlobalObject()
		for _, value := range []any{1, "1", 0.1, true} {
			object.SetProperty("value", value)
			property, err := object.GetProperty("value")
			assert.NoError(t, err)
			assert.Equal(t, value, property.ToNative())
		}
		values := []any{
			int8(1), int16(1), int32(1), int64(1),
			uint(1), uint8(1), uint16(1), uint32(1), uint64(1),
		}
		for _, value := range values {
			object.SetProperty("value", value)
			var expected int
			valueOf := reflect.ValueOf(value)
			if valueOf.CanInt() {
				expected = int(valueOf.Int())
			} else {
				expected = int(valueOf.Uint())
			}
			property, err := object.GetProperty("value")
			assert.NoError(t, err)
			assert.Equal(t, expected, property.ToNative())
		}
		values = []any{
			uint(math.MaxUint), uint32(math.MaxUint32),
			int64(math.MaxInt64), uint64(math.MaxUint64),
		}
		for _, value := range values {
			object.SetProperty("value", value)
			var expected float64
			valueOf := reflect.ValueOf(value)
			if valueOf.CanInt() {
				expected = float64(valueOf.Int())
			} else {
				expected = float64(valueOf.Uint())
			}
			property, err := object.GetProperty("value")
			assert.NoError(t, err)
			assert.Equal(t, expected, property.ToNative())
		}
		values = []any{
			[]int8{1}, []int16{1}, []uint16{1}, []int32{1}, []uint32{1},
			[]float32{1}, []float64{1},
		}
		for _, value := range values {
			object.SetProperty("value", value)
			property, err := object.GetProperty("value")
			assert.NoError(t, err)
			assert.Equal(t, value, property.ToNative())
		}
	})
}

func TestCallGoFunction(t *testing.T) {
	NewRuntime().NewContext().With(func(context *Context) {
		object := context.GlobalObject()
		counter := 0
		naiveFunc := func(args ...any) (any, error) {
			counter = len(args)
			return args[0].(int) + args[1].(int), nil
		}
		object.SetProperty("test", naiveFunc)
		retval, err := context.Eval("test(1, 2);")
		assert.NoError(t, err)
		assert.Equal(t, 2, counter)
		assert.Equal(t, 3, retval.ToNative())
	})
}

func TestFinalizer(t *testing.T) {
	runtime := NewRuntime()
	runtime.NewContext().With(func(context *Context) {
		value := context.addGoObject(nil)
		Value{context, value}.free()
		assert.Zero(t, 0, len(context.goValues))
	})
}

func BenchmarkGetKind(b *testing.B) {
	NewRuntime().NewContext().With(func(context *Context) {
		retval, err := context.Eval("new Date()")
		assert.NoError(b, err)
		assert.Equal(b, KindDate, retval.Object().Kind())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = retval.Object().Kind()
		}
	})
}

func BenchmarkObjectFromNative(b *testing.B) {
	NewRuntime().NewContext().With(func(context *Context) {
		expected := make(map[string]any)
		for i := 0; i < 16; i++ {
			expected[strconv.Itoa(i)] = i
		}
		global := context.GlobalObject()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			global.SetProperty("whatever", expected)
		}
		b.StopTimer()
	})
}

func BenchmarkCallGoFunction(b *testing.B) {
	NewRuntime().NewContext().With(func(context *Context) {
		naiveFunc := func(args ...any) any {
			return args[0].(int) + args[1].(int)
		}
		global := context.GlobalObject()
		global.SetProperty("test", naiveFunc)
		retval, err := context.Eval("test(1, 2)")
		assert.NoError(b, err)
		assert.Equal(b, 3, retval.ToNative())
		for i := 0; i < b.N; i++ {
			_, _ = context.Eval("test(1, 2)")
		}
	})
}
