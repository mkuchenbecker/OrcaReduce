package chaos

import (
	"math/rand"

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

func NewRandomStaticErrorInjector(err error, rate float64) orcareduce.ErrorInjector {
	return &randomStaticErrors{err: err, errorRate: rate}
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

func NewMetaErrorInjector(errors []orcareduce.ErrorInjector) orcareduce.ErrorInjector {
	return &metaErrors{errors: errors}
}
