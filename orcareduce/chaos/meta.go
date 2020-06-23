package chaos

import (
	"fmt"
	"time"

	"github.com/mkuchenbecker/orcareduce/orcareduce"
)

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
				NewStaticLatency(20 * time.Millisecond),
				NewDynamicLatency(20 * time.Millisecond),
				NewRandomLatency(NewDynamicLatency(100*time.Millisecond), 0.1),
				NewRandomLatency(NewDynamicLatency(1000*time.Millisecond), 0.01),
				NewRandomLatency(NewDynamicLatency(5000*time.Millisecond), 0.001),
				NewRandomLatency(NewDynamicLatency(30*time.Second), 0.0001),
			},
		},
		errors: &metaErrors{
			errors: []orcareduce.ErrorInjector{
				NewRandomStaticErrorInjector(fmt.Errorf("[chaos] encountered an error"), 0.1),
			},
		},
	}
}
