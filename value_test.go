package quickjs

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToNative(t *testing.T) {
	NewRuntime().NewContext().With(func(context *Context) {
		jsValue, err := context.GlobalObject().GetProperty("nonExist")
		assert.NoError(t, err)
		assert.Equal(t, Undefined, jsValue.ToNative())

		jsValue, err = context.Eval("null")
		assert.NoError(t, err)
		assert.Equal(t, nil, jsValue.ToNative())

		jsValue, err = context.Eval("true")
		assert.NoError(t, err)
		assert.Equal(t, true, jsValue.ToNative())

		jsValue, err = context.Eval("1")
		assert.NoError(t, err)
		assert.Equal(t, 1, jsValue.ToNative())

		jsValue, err = context.Eval("1.1")
		assert.NoError(t, err)
		assert.Equal(t, 1.1, jsValue.ToNative())

		jsValue, err = context.Eval(strconv.Itoa(math.MaxInt32))
		assert.NoError(t, err)
		assert.Equal(t, math.MaxInt32, jsValue.ToNative())

		jsValue, err = context.Eval(strconv.Itoa(math.MaxInt64))
		assert.NoError(t, err)
		assert.Equal(t, float64(math.MaxInt64), jsValue.ToNative())

		jsValue, err = context.Eval("0n")
		assert.NoError(t, err)
		assert.Equal(t, 0, jsValue.ToNative())

		jsValue, err = context.Eval("1n")
		assert.NoError(t, err)
		assert.Equal(t, 1, jsValue.ToNative())

		jsValue, err = context.Eval("-1n")
		assert.NoError(t, err)
		assert.Equal(t, -1, jsValue.ToNative())

		jsValue, err = context.Eval(strconv.Itoa(math.MaxInt64) + "n")
		assert.NoError(t, err)
		assert.Equal(t, math.MaxInt64, jsValue.ToNative())

		jsValue, err = context.Eval(fmt.Sprintf("%d", uint64(math.MaxUint64)) + "n")
		assert.NoError(t, err)
		var bigInt big.Int
		bigInt.SetUint64(math.MaxUint64)
		assert.Equal(t, bigInt, jsValue.ToNative())

		jsValue, err = context.Eval(`"string"`)
		assert.NoError(t, err)
		assert.Equal(t, "string", jsValue.ToNative())
	})
}

func TestFromNative(t *testing.T) {
	NewRuntime().NewContext().With(func(context *Context) {
		global := context.GlobalObject()
		global.SetProperty("testValue", nil)
		retval, err := context.Eval("testValue == null")
		assert.NoError(t, err)
		assert.True(t, retval.ToNative().(bool))

		global.SetProperty("testValue", true)
		property, err := global.GetProperty("testValue")
		assert.NoError(t, err)
		assert.True(t, property.ToNative().(bool))

		integers := []any{
			int(1), int8(1), int16(1), int32(1), int64(1),
			uint(1), uint8(1), uint16(1), uint32(1), uint64(1),
		}
		for _, value := range integers {
			global.SetProperty("testValue", value)
			property, err := global.GetProperty("testValue")
			assert.NoError(t, err)
			assert.Equal(t, 1, property.ToNative())
		}

		bigNumber := big.NewInt(math.MaxInt64)
		bigNumber.Add(bigNumber, bigNumber)
		global.SetProperty("bigNumber", bigNumber)
		property, err = global.GetProperty("bigNumber")
		assert.NoError(t, err)
		assert.Equal(t, *bigNumber, property.ToNative())

		numbers := []any{
			math.MaxUint32, int64(math.MaxUint32),
			uint(math.MaxUint32), uint32(math.MaxUint32), uint64(math.MaxUint32),
			float64(math.MaxUint32),
		}
		for _, value := range numbers {
			global.SetProperty("testValue", value)
			property, err := global.GetProperty("testValue")
			assert.NoError(t, err)
			assert.Equal(t, float64(math.MaxUint32), property.ToNative())
		}

		global.SetProperty("testValue", float32(1.1))
		property, err = global.GetProperty("testValue")
		assert.NoError(t, err)
		assert.Equal(t, float64(float32(1.1)), property.ToNative())
	})
}

func TestJSON(t *testing.T) {
	NewRuntime().NewContext().With(func(context *Context) {
		value, err := context.Eval("new Object({a: 1, b: 2})")
		assert.NoError(t, err)
		expected := `{"a":1,"b":2}`
		assert.Equal(t, expected, value.JSONify())

		type ut struct {
			A int `json:"a"`
			B int `json:"b"`
		}
		var actual ut
		assert.NoError(t, value.Object().UnmarshalJSON(&actual))
		assert.Equal(t, ut{1, 2}, actual)

		global := context.GlobalObject()
		global.SetProperty("jsonValue", json.RawMessage([]byte(expected)))
		value, err = global.GetProperty("jsonValue")
		assert.NoError(t, err)
		assert.Equal(t, expected, value.JSONify())
	})
}

func TestUndefined(t *testing.T) {
	NewRuntime().NewContext().With(func(context *Context) {
		context.GlobalObject().SetProperty("whatever", Undefined)
		value, _ := context.Eval("whatever")
		assert.NotEqual(t, nil, value.ToNative())
		assert.Equal(t, Undefined, value.ToNative())
	})
}

func BenchmarkGetType(b *testing.B) {
	NewRuntime().NewContext().With(func(context *Context) {
		jsValue, err := context.Eval("new Object()")
		assert.NoError(b, err)
		assert.Equal(b, TypeObject, jsValue.Type())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = jsValue.Type()
		}
	})
}
