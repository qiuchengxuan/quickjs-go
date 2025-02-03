package quickjs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteCode(t *testing.T) {
	context := NewRuntime().NewContext()
	byteCode, err := context.Compile("1 + 1")
	assert.NoError(t, err)
	retval, err := context.EvalBinary(byteCode)
	assert.NoError(t, err)
	assert.Equal(t, 2, retval.ToPrimitive())
}
