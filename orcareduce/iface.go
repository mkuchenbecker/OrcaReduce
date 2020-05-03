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
	Bytes() []byte
	Parent() ID
}

type Director interface {
	GetNextActor() error
	Act() error
}
