package quickjs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	context := NewRuntime().NewContext()
	value, err := context.Eval(`new Map([["key", "value"], ["int", 1]])`)
	assert.NoError(t, err)
	expected := map[string]any{"key": "value", "int": 1}
	assert.Equal(t, expected, value.ToPrimitive())
}
