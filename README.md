quickjs
=======

Go bindings to QuickJS: a fast, small, and embeddable ES2020 JavaScript interpreter.

Example
-------

```go
package main

import (
	"github.com/qiuchengxuan/go-quickjs"
)

func main() {
	quickjs.NewRuntime().NewContext().With(context *Context) {
		value, _ := context.Eval("new Map()")
		_ = value.ToNative().(map[string]any)

		context := quickjs.NewRuntime().NewContext()
		value, _ = context.Eval(`let value = "value"`)
		_ = context.GlobalObject().GetProperty("value").ToNative() // Should be "value"

		byteCode, _ := context.Compile("1 + 1")
		value, _ = context.EvalBinary(byteCode)
		value.ToNative() // Should be 2

		context.GlobalObject().SetProperty("value", 0.1)
		value, _ = context.Eval(`value + 0.1`)
		_ = value.ToNative() // should be 0.2

		counter := 0
		context.GlobalObject().SetProperty("sum", func(args ...any) any {
		    return args[0].(int) + args[1].(int)
		})
		value, _ = context.Eval("test(1, 2)")
		_ = value.ToNative() // should be 3 and counter should be 2
    })
}
```

Set property to JS object
-------------------------

When setting properties to JS objects, values are converted as following:

| Go Value                  | JS Value    |
|---------------------------|-------------|
| nil                       | null        |
| Undefined                 | undefined   |
| bool                      | boolean     |
| (u)int(*)/float32/float64 | Number      |
| big.Int                   | bigint      |
| string                    | string      |
| []uint8                   | Uint8Array  |
| []uint16                  | Uint16Array |
| []uint32                  | Uint32Array |
| []int8                    | Int8Array   |
| []int16                   | Int16Array  |
| []int32                   | Int32Array  |
| []any or map[string]any   | object      |
| map[\*]\*                 | Map         |
| []\*                      | Array       |
| *                         | undefined   |

Convert to native value from JS
-------------------------------

Value returned by `Eval` or `GetProperty` can be further exported as
Go representation with `ToPrimitive` or `ToNative`.

Value converted as following:

| JS Value     | Go Value                |
|--------------|-------------------------|
| null         | nil                     |
| undefined    | Undefined               |
| boolean      | bool                    |
| Number       | int or float64          |
| bigint       | int or big.Int          |
| string       | string                  |
| object       | []any or map[string]any |
| Array        | []any                   |
| ArrayBuffer  | []byte                  |
| Uint8Array   | []uint8                 |
| Uint16Array  | []uint16                |
| Uint32Array  | []uint32                |
| Int8Array    | []int8                  |
| Int16Array   | []int16                 |
| Int32Array   | []int32                 |
| Float32Array | []float32               |
| Float64Array | []float64               |
| Map          | map[any]any             |
| Set          | []any                   |
| Date         | time.Time               |
| *            | NotNative               |
