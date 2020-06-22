package chaos

import (
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mkuchenbecker/orcareduce/orcareduce"
	"github.com/mkuchenbecker/orcareduce/orcareduce/mock"
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

func TestRandomLatency(t *testing.T) {
	t.Parallel()
	rand.Seed(1)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockLatency := mock.NewMockLatencyInjector(mockCtrl)

	random := NewRandomLatency(mockLatency, 1)
	mockLatency.EXPECT().Latency().Times(1)
	random.Latency()

	random = NewRandomLatency(mockLatency, 0)
	random.Latency()
}

func TestMetaLatency(t *testing.T) {
	t.Parallel()
	rand.Seed(1)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockLatency0 := mock.NewMockLatencyInjector(mockCtrl)
	mockLatency0.EXPECT().Latency().Times(1)
	mockLatency1 := mock.NewMockLatencyInjector(mockCtrl)
	mockLatency1.EXPECT().Latency().Times(1)

	meta := NewMetaLatency([]orcareduce.LatencyInjector{mockLatency0, mockLatency1})
	meta.Latency()
}
