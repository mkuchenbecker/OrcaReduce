package identifier

import (
	"fmt"
	"testing"

	"github.com/mkuchenbecker/orcareduce/orcareduce"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentifierLineage(t *testing.T) {
	id := New().WithScope("parent_scope")
	child := id.NewChild().WithScope("child_scope")
	grandchild := child.NewChild().WithScope("grandchild_scope")

	lineage := grandchild.Lineage()
	require.Equal(t, 3, len(lineage))
	assert.Equal(t, id, lineage[2])
	assert.Equal(t, child, lineage[1])
	assert.Equal(t, grandchild, lineage[0])
}

func TestIdentifierValueStringEquivalence(t *testing.T) {
	id := identifier{
		id:    "bar",
		scope: orcareduce.Scope("foo"),
	}
	assert.Equal(t, "foo(bar)", id.String())
}

func TestIdentifierNewChild(t *testing.T) {
	id := New().WithScope("parent_scope")
	child := id.NewChild().WithScope("child_scope")
	grandchild := child.NewChild().WithScope("grandchild_scope")
	assert.Equal(t, 0, id.Depth())
	assert.Equal(t, 1, child.Depth())
	assert.Equal(t, 2, grandchild.Depth())

	assert.Equal(t, id, child.Parent())
	assert.Equal(t, child, grandchild.Parent())

	assert.Equal(t, fmt.Sprintf("parent_scope(%s).child_scope(%s).grandchild_scope(%s)",
		id.Value(),
		child.Value(),
		grandchild.Value(),
	),
		grandchild.String(),
	)
}
