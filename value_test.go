package quickjs

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetProperty(t *testing.T) {
	context := NewRuntime().NewContext()
	object := context.GlobalObject()
	for _, value := range []any{1, "1", 0.1, true} {
		object.SetProperty("value", value)
		assert.Equal(t, value, object.GetProperty("value").ToPrimitive())
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
		assert.Equal(t, expected, object.GetProperty("value").ToPrimitive())
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
		assert.Equal(t, expected, object.GetProperty("value").ToPrimitive())
	}
	values = []any{
		[]int8{1}, []int16{1}, []uint16{1}, []int32{1}, []uint32{1},
		[]float32{1}, []float64{1},
	}
	for _, value := range values {
		object.SetProperty("value", value)
		assert.Equal(t, value, object.GetProperty("value").ToPrimitive())
	}
}

func TestSetFunction(t *testing.T) {
	context := NewRuntime().NewContext()
	object := context.GlobalObject()
	counter := 0
	naiveFunc := func(args ...any) any {
		counter = len(args)
		return args[0].(int) + args[1].(int)
	}
	object.SetProperty("test", naiveFunc)
	retval, err := context.Eval("test(1, 2)")
	assert.NoError(t, err)
	assert.Equal(t, 2, counter)
	assert.Equal(t, 3, retval.ToPrimitive())
}
