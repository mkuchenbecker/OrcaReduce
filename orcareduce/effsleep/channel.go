package effsleep

import (
	"context"
	"errors"
	"sync"

	"github.com/mkuchenbecker/orcareduce/orcareduce/chaos"
)

type ChannelInjector interface {
	Done()
	Sync()
}

type chanInjector chan bool

func NewInjector() ChannelInjector {
	out := chanInjector(make(chan bool))
	return &out
}

func (this *chanInjector) Done() {
	*this <- true
}

func (this *chanInjector) Sync() {
	<-*this
}

type InjectorManager interface {
	Subscribe(key string) ChannelInjector
	Done(key string) bool
	Sync(key string) bool
}

type injectorManager struct {
	injectors map[string]ChannelInjector
	mux       sync.RWMutex
}

func NewInjectorManager() InjectorManager {
	return &injectorManager{injectors: make(map[string]ChannelInjector)}
}

func (this *injectorManager) Subscribe(key string) ChannelInjector {
	injector := NewInjector()
	this.injectors[key] = injector
	return injector
}

func (this *injectorManager) Done(key string) (ok bool) {
	injector, ok := this.injectors[key]
	if !ok {
		return
	}
	injector.Done()
	return
}

func (this *injectorManager) Sync(key string) (ok bool) {
	injector, ok := this.injectors[key]
	if !ok {
		return
	}
	injector.Done()
	return
}

type IndeterminateRuntime struct {
	chaos chaos.Injector
}

const InjectorManagerContextKey = "INJECTION_MANAGER"

func Finished(ctx context.Context) error {
	manager := ctx.Value(InjectorManagerContextKey)
	if manager == nil {
		return errors.New("unable to find manager")
	}
	man, ok := manager.(ChannelInjector)
	if !ok {
		return errors.New("manager wrong type")
	}
	man.Done()
	return nil
}

func (this *IndeterminateRuntime) Run(doneFunc func(), i *int) {
	defer doneFunc()
	this.chaos.Latency()
	*i = *i + 1
}
func (this *IndeterminateRuntime) RunCtx(ctx context.Context, i *int) error {
	this.chaos.Latency()
	*i = *i + 1
	return Finished(ctx)
}
