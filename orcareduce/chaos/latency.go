package chaos

import (
	"math/rand"
	"time"

	"github.com/mkuchenbecker/orcareduce/orcareduce"
)

type sleepFunction func(time.Duration)

type staticLatency struct {
	latency time.Duration
	sleep   sleepFunction
}

func (cfg *staticLatency) Latency() {
	time.Sleep(cfg.latency)
}

func NewStaticLatency(latency time.Duration) orcareduce.LatencyInjector {
	return &staticLatency{
		sleep:   time.Sleep,
		latency: latency,
	}
}

type dynamicLatency struct {
	maxLatency time.Duration
	sleep      sleepFunction
}

func (cfg *dynamicLatency) Latency() {
	val := rand.Int63n(int64(cfg.maxLatency))
	cfg.sleep(time.Duration(val))
}

func NewDynamicLatency(maxLatency time.Duration) orcareduce.LatencyInjector {
	return &dynamicLatency{
		sleep:      time.Sleep,
		maxLatency: maxLatency,
	}
}

type randomLatency struct {
	latency orcareduce.LatencyInjector
	percent float64
}

func (cfg *randomLatency) Latency() {
	if rand.Float64() < cfg.percent {
		cfg.latency.Latency()
	}
}

func NewRandomLatency(latency orcareduce.LatencyInjector, percent float64) orcareduce.LatencyInjector {
	return &randomLatency{latency: latency, percent: percent}
}

type metaLatency struct {
	latency []orcareduce.LatencyInjector
}

func (this *metaLatency) Latency() {
	for _, latency := range this.latency {
		latency.Latency()
	}
}

func NewMetaLatency(latency []orcareduce.LatencyInjector) orcareduce.LatencyInjector {
	return &metaLatency{latency: latency}
}
