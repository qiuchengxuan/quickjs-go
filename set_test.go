package quickjs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	context := NewRuntime().NewContext()
	value, err := context.Eval(`new Set([1, 2, 1])`)
	assert.NoError(t, err)
	assert.Equal(t, []any{1, 2}, value.ToPrimitive())
}
