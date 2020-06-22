package chaos

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStaticLatency(t *testing.T) {
	t.Parallel()
	latency := 100 * time.Millisecond
	static := NewStaticLatency(latency).(*staticLatency)

	static.sleep = func(duration time.Duration) {
		assert.Equal(t, latency, duration)
	}

	static.Latency()
}

func TestDynamicLatency(t *testing.T) {
	t.Parallel()
	rand.Seed(1)
	maxLatency := 100 * time.Millisecond
	dynamic := NewDynamicLatency(maxLatency).(*dynamicLatency)

	dynamic.sleep = func(duration time.Duration) {
		assert.True(t, maxLatency >= duration)
	}

	dynamic.Latency()
}
