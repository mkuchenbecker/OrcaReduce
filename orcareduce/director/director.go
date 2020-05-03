package prod

import (
	"github.com/mkuchenbecker/orcareduce/orcareduce"
	"github.com/mkuchenbecker/orcareduce/orcareduce/exceptions"
	"github.com/pkg/errors"
)



type config struct {
	retries int
	backoffFactor int
}

func (c *config) Backoff(retry int) time.Duration {
	return math.Max(1,2) * time.Second
}

func (c *config) Retries() int {

}


type director struct {
	actor orcareduce.Reactor
	cfg config
}

func (d *director) Direct() (err error) {
	defer func() {
		err = errors.Wrap(err,"unable to direct")
	}()
	actor, err := d.GetNextActor()
	if err != nil {
		return  err
	}
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


func (d *director) Run() {
	for i:= 0; i < d.config.Retries(); i++ {


	}
}

