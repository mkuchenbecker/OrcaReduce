package exceptions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreconditionError(t *testing.T) {
	err := PreconditionError("error")
	assert.Equal(t, "preconditions not met: error", err.Error())
}
