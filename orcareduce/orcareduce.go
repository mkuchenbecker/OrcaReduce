package orcareduce

// Reactor is a repeatable actor.
type Reactor interface {
	Preconditions() error
	Act() error
	SignalSuccess() error
	Notify()
	ID() ID
}

type Precondition interface {
	Check() error
}

// ID is a generalized interface for identifying reactors.
type ID interface {
	String() string
	Value() string
	Lineage() []ID
	Parent() ID
	NewChild() ID
	Depth() int
	Scope() Scope
	WithScope(Scope) ID
}

type Director interface {
	GetNextActor() error
	Act() error
	ID() ID
}

type ActorDataSink interface {
	Success(Reactor) error
	Failure(Reactor)
	Attempt(Reactor) error
}

type Scope string
