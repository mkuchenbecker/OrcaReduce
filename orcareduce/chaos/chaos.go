package chaos

import (
	"fmt"
	"math/rand"
	"time"
)

type randomStaticErrors struct {
	err       error
	errorRate float64
}

func (cfg *randomStaticErrors) Error() error {
	if rand.Float64() <= cfg.errorRate {
		return cfg.err
	}
	return nil
}

type staticLatency struct {
	latency time.Duration
}

type metaErrors struct {
	errors []ErrorInjector
}

func (cfg *metaErrors) Error() error {
	for _, injector := range cfg.errors {
		err := injector.Error()
		if err != nil {
			return err
		}
	}
	return nil
}

func (cfg *staticLatency) Latency() {
	time.Sleep(cfg.latency)
}

type dynamicLatency struct {
	maxLatency time.Duration
}

func (cfg *dynamicLatency) Latency() {
	val := rand.Int63n(int64(cfg.maxLatency))
	time.Sleep(time.Duration(val))
}

type randomLatency struct {
	latency LatencyInjector
	percent float64
}

func (cfg *randomLatency) Latency() {
	if rand.Float64() < cfg.percent {
		cfg.latency.Latency()
	}
}

type metaLatency struct {
	latency []LatencyInjector
}

func (this *metaLatency) Latency() {
	for _, latency := range this.latency {
		latency.Latency()
	}
}

type metaInjector struct {
	latency LatencyInjector
	errors  ErrorInjector
}

func (cfg *metaInjector) Latency() {
	cfg.latency.Latency()
}

func (cfg *metaInjector) Error() error {
	return cfg.errors.Error()
}

func NewDefault() Injector {
	return &metaInjector{
		latency: &metaLatency{
			latency: []LatencyInjector{
				&staticLatency{latency: 20 * time.Millisecond},
				&dynamicLatency{maxLatency: 20 * time.Millisecond},
				&randomLatency{percent: 0.1, latency: &dynamicLatency{maxLatency: 100 * time.Millisecond}},
				&randomLatency{percent: 0.01, latency: &dynamicLatency{maxLatency: 1000 * time.Millisecond}},
				&randomLatency{percent: 0.001, latency: &dynamicLatency{maxLatency: 5000 * time.Millisecond}},
				&randomLatency{percent: 0.0001, latency: &dynamicLatency{maxLatency: 30 * time.Second}},
			},
		},
		errors: &metaErrors{
			errors: []ErrorInjector{
				&randomStaticErrors{err: fmt.Errorf("[chaos] encountered an error"), errorRate: 0.1},
			},
		},
	}
}
