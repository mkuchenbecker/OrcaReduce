package effsleep

import (
	"context"
	"testing"

	"github.com/mkuchenbecker/orcareduce/orcareduce/chaos"
	"github.com/stretchr/testify/assert"
)

func TestInjector(t *testing.T) {
	t.Parallel()

	syncFunctionName := "IndeterminateRuntime.Run"
	manager := NewInjectorManager()
	injector := manager.Subscribe(syncFunctionName)
	m := map[string]ChannelInjector{"0":}
	ctx := context.WithValue(context.Background(), InjectorManagerContextKey, injector)

	runtime := IndeterminateRuntime{chaos: chaos.NewDefault()}

	i := 0
	go func() {
		for i := 0; i < 10; i++ {
			assert.NoError(t, runtime.RunCtx(ctx, &i))
		}
	}()
	injector.Sync()
	assert.Equal(t, 1, i)

}

type SyncFunc func() error

type SafeGoRunner struct {
}



func RunAsync(