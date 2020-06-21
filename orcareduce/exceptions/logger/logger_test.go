package exceptions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger_Errorf(t *testing.T) {
	expectedFormat := "format %s"
	arg := "arg"

	logger := logger{
		errorf: func(format string, args ...interface{}) {
			assert.Equal(t, expectedFormat, format)
			require.Equal(t, 1, len(args))
			assert.Equal(t, "arg", args[0].(string))
		},
	}
	logger.Errorf(expectedFormat, arg)
	DefaultLogger().Errorf(expectedFormat, arg)
}

func TestLogger_Infof(t *testing.T) {
	expectedFormat := "format %s"
	arg := "arg"

	logger := logger{
		infof: func(format string, args ...interface{}) {
			assert.Equal(t, expectedFormat, format)
			require.Equal(t, 1, len(args))
			assert.Equal(t, "arg", args[0].(string))
		},
	}
	logger.Infof(expectedFormat, arg)
	DefaultLogger().Infof(expectedFormat, arg)
}

func TestLogger_Error(t *testing.T) {
	expectedError := fmt.Errorf("error")
	expectedFormat := "%s: %+v"
	expectedMessage := "msg"

	logger := logger{
		errorf: func(format string, args ...interface{}) {
			assert.Equal(t, expectedFormat, format)
			require.Equal(t, 2, len(args))
			assert.Equal(t, expectedMessage, args[0].(string))
			assert.Equal(t, expectedError, args[1].(error))
		},
	}
	logger.Error(expectedError, expectedMessage)
	DefaultLogger().Error(expectedError, expectedMessage)
}
