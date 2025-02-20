package quickjs

import "encoding/json"

type ByteCode []byte

type NotNative struct{}

type undefined struct{}

var Undefined *undefined = nil

type NaiveFunc = func(...any) (any, error)

type JSONValue interface {
	json.Marshaler
	json.Unmarshaler
}

type AsJSONValue[T any] struct{ value T }

func (c AsJSONValue[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

func (c AsJSONValue[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.value)
}
