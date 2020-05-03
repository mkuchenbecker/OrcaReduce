package prod

import (
	"github.com/mkuchenbecker/orcareduce/orcareduce"
	"github.com/mkuchenbecker/orcareduce/orcareduce/exceptions"
	"golang.org/pkg/errors"
)

type queueDirector struct {
	actors []orcareduce.Reactor
}

func (d * queueDirector) GetNextActor() (actor orcareduce.Reactor, err error) {
	if len(d.actors) == 0 {
		return nil, exceptions.ErrActorsDepleted
	}
	actor = d.actors[0]
	return actor,nil
}


func (d *queueDirector) Direct() (err error) {
	defer func() {
		err = errors.Wrap(err,"unable to direct")
	}
	actor, err := d.GetNextActor()
	if err != nil {
		return  err
	}
	err = actor.Preconditions()
	if err != nil {
		return err
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