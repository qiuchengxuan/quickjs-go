package quickjs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	context := NewRuntime().NewContext()
	value, err := context.Eval("new Date(8.64e15)")
	assert.NoError(t, err)
	expected := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, value.ToPrimitive())
}
