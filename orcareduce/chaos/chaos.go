package chaos

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mkuchenbecker/orcareduce/orcareduce"
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

type metaErrors struct {
	errors []orcareduce.ErrorInjector
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

type metaInjector struct {
	latency orcareduce.LatencyInjector
	errors  orcareduce.ErrorInjector
}

func (cfg *metaInjector) Latency() {
	cfg.latency.Latency()
}

func (cfg *metaInjector) Error() error {
	return cfg.errors.Error()
}

func NewDefault() orcareduce.Injector {
	return &metaInjector{
		latency: &metaLatency{
			latency: []orcareduce.LatencyInjector{
				&staticLatency{latency: 20 * time.Millisecond},
				&dynamicLatency{maxLatency: 20 * time.Millisecond},
				&randomLatency{percent: 0.1, latency: &dynamicLatency{maxLatency: 100 * time.Millisecond}},
				&randomLatency{percent: 0.01, latency: &dynamicLatency{maxLatency: 1000 * time.Millisecond}},
				&randomLatency{percent: 0.001, latency: &dynamicLatency{maxLatency: 5000 * time.Millisecond}},
				&randomLatency{percent: 0.0001, latency: &dynamicLatency{maxLatency: 30 * time.Second}},
			},
		},
		errors: &metaErrors{
			errors: []orcareduce.ErrorInjector{
				&randomStaticErrors{err: fmt.Errorf("[chaos] encountered an error"), errorRate: 0.1},
			},
		},
	}
}
