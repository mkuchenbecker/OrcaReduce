package orcareduce

import (
	"bytes"
	"crypto/sha256"
	"time"
)

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
	NewScopedChild(scope Scope) ID
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

type identifier struct {
	id     string
	scope  Scope
	level  int
	parent *identifier
}

func (this *identifier) Lineage() []ID {
	out := make([]ID, 0)
	curr := this
	for {
		out = append(out, curr)
		if curr.parent == nil {
			break
		}
		curr = curr.parent
	}
	return out
}

func (this *identifier) String() string {
	lineage := this.Lineage()
	var buf bytes.Buffer
	for i := len(lineage) - 1; i >= 0; i++ {
		buf.WriteString(lineage[i].Value())
	}
	return buf.String()
}

func (this *identifier) Value() string {
	return this.id
}

func (this *identifier) Parent() ID {
	return this.parent
}

func (this *identifier) NewChild() ID {
	out := NewID()
	out.level = this.level + 1
	out.parent = this
	return out
}

func (this *identifier) NewScopedChild(scope Scope) ID {
	out := NewID()
	out.scope = scope
	out.level = this.level + 1
	out.parent = this
	return out
}

func NewID() *identifier {
	sha := sha256.New()
	val := sha.Sum([]byte(time.Now().String())) // Re-hash the parent ID.
	return &identifier{
		id: string(val),
	}
}
