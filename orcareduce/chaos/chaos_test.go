package chaos

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultChaos(t *testing.T) {
	rand.Seed(5) // Static seed so the test is deterministic.
	chaos := NewDefault()

	now := time.Now()
	chaos.Latency()
	assert.True(t, time.Since(now) >= time.Millisecond*20)

	errCount := 0
	for i := 0; i < 1000; i++ {
		if err := chaos.Error(); err != nil {
			errCount++
		}
	}
	assert.Equal(t, 105, errCount)
}
