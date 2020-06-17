package state

import (
	"sync"

	"github.com/mkuchenbecker/orcareduce/orcareduce"
	"github.com/pkg/errors"
)

type State string

var (
	PENDING State = "PENDING"
	SUCCESS State = "SUCCESS"
	FAILURE State = "FAILURE"
)

type document struct {
	actor   orcareduce.Reactor
	state   State
	attempt int
}

type inMemorySink struct {
	data map[orcareduce.ID]document
	mux  sync.RWMutex
}

func (sink *inMemorySink) Success(r orcareduce.Reactor) (err error) {
	defer func() {
		err = errors.Wrap(err, "unable to record success")
	}()
	sink.mux.Lock()
	defer sink.mux.Unlock()
	val, ok := sink.data[r.ID()]
	if !ok {
		val = document{actor: r}
	}
	if val.state == FAILURE {
		return errors.New("record already marked as failure")
	}
	val.state = SUCCESS
	sink.data[r.ID()] = val
	return nil
}

func (sink *inMemorySink) Failure(r orcareduce.Reactor) {
	sink.mux.Lock()
	defer sink.mux.Unlock()
	val, ok := sink.data[r.ID()]
	if !ok {
		val = document{actor: r}
	}
	if val.state == SUCCESS {
		return
	}
	val.state = FAILURE
	sink.data[r.ID()] = val
	return
}

func (sink *inMemorySink) Attempt(r orcareduce.Reactor) error {
	sink.mux.Lock()
	defer sink.mux.Unlock()
	val, ok := sink.data[r.ID()]
	if !ok {
		val = document{actor: r}
	}
	if val.state != PENDING {
		return errors.New("only make an attempt on pending requests")
	}
	val.attempt++
	sink.data[r.ID()] = val
	return nil
}

type DatastoreActor struct {
	actor orcareduce.Reactor
	sink  orcareduce.ActorDataSink
}

func (d *DatastoreActor) SignalSuccess() error {
	return d.sink.Success(d)
}

func (d *DatastoreActor) Preconditions() error {
	return d.actor.Preconditions()
}

func (d *DatastoreActor) Act() error {
	err := d.sink.Attempt(d)
	if err != nil {
		return err
	}
	return d.Act()
}

func (d *DatastoreActor) Notify() {
	d.actor.Notify()
}

func (d *DatastoreActor) ID() orcareduce.ID {
	return d.ID()
}
