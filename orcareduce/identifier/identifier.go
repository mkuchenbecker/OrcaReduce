package identifier

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/mkuchenbecker/orcareduce/orcareduce"
)

type identifier struct {
	id     string
	scope  orcareduce.Scope
	depth  int
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
	for i := len(lineage) - 1; i >= 0; i-- {
		buf.WriteString(fmt.Sprintf("%s(%s)", lineage[i].Scope(), lineage[i].Value()))
		buf.WriteString(".")
	}
	return strings.Trim(buf.String(), ".")
}

func (this *identifier) Value() string {
	return this.id
}

func (this *identifier) Parent() orcareduce.ID {
	return this.parent
}

func (this *identifier) NewChild() orcareduce.ID {
	out := New()
	out.depth = this.depth + 1
	out.parent = this
	return out
}

func (this *identifier) Depth() int {
	return this.depth
}

func (this *identifier) Scope() orcareduce.Scope {
	return this.scope
}

func (this *identifier) WithScope(scope orcareduce.Scope) orcareduce.ID {
	this.scope = scope
	return this
}

func New() *identifier {
	sha := sha256.New()
	val := sha.Sum([]byte(time.Now().String())) // Re-hash the parent ID.
	return &identifier{
		id: string(val),
	}
}
