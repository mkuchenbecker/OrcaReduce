package director

import (
	"math"
	"time"

	"github.com/mkuchenbecker/orcareduce/orcareduce"
	"github.com/mkuchenbecker/orcareduce/orcareduce/exceptions"
	"github.com/pkg/errors"
)

type config struct {
	maxAttempts int
}

func (c *config) Backoff(attempt int) time.Duration {
	return time.Duration(int64(math.Max(float64(attempt), 2))) * time.Second
}

func (c *config) Attempts() int {
	return c.maxAttempts
}

type director struct {
	actor   orcareduce.Reactor
	cfg     config
	handler exceptions.Handler
	id      orcareduce.ID
	clock   Clock
}

func (d *director) Direct() (err error) {
	defer func() {
		err = errors.Wrap(err, "unable to direct")
	}()
	actor := d.actor
	err = actor.Preconditions()
	if err != nil {
		return exceptions.PreconditionError(err.Error())
	}
	err = actor.Act()
	if err != nil {
		return err
	}
	err = actor.SignalSuccess()
	if err != nil {
		return err
	}
	actor.Notify()

	return nil
}

func (d *director) Run() (err error) {
	defer d.handler.HandlePanic(&err)
	runtime := NewRuntime(d.id, d.clock)
	defer func() {
		err := runtime.Save()
		_ = d.handler.HandleError(err)
	}()

	for i := 0; i < d.cfg.Attempts(); i++ {
		_, endFunc := runtime.StartRun()
		err := d.Direct()
		endFunc(err, "")
		if d.handler.HandleError(err) != nil {
			return err
		}
		d.cfg.Backoff(i)
	}
	return nil
}

func (d *director) RunAsync(err chan error) {
	err <- d.Run()
}
