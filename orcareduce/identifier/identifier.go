package identifier

import (
	"bytes"
	"crypto/sha256"
	"time"

	"github.com/mkuchenbecker/orcareduce/orcareduce"
)

type identifier struct {
	id     string
	scope  orcareduce.Scope
	level  int
	parent *identifier
}

func (this *identifier) Lineage() []orcareduce.ID {
	out := make([]orcareduce.ID, 0)
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

func (this *identifier) Parent() orcareduce.ID {
	return this.parent
}

func (this *identifier) NewChild() orcareduce.ID {
	out := New()
	out.level = this.level + 1
	out.parent = this
	return out
}

func (this *identifier) NewScopedChild(scope orcareduce.Scope) orcareduce.ID {
	out := New()
	out.scope = scope
	out.level = this.level + 1
	out.parent = this
	return out
}

func New() *identifier {
	sha := sha256.New()
	val := sha.Sum([]byte(time.Now().String())) // Re-hash the parent ID.
	return &identifier{
		id: string(val),
	}
}
