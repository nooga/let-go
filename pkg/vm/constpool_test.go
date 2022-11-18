package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstpool(t *testing.T) {
	c := NewConsts()
	v := []Value{Int(9), String("zoom"), Char('a'), Symbol("poop"), EmptyList, NIL, TRUE, FALSE}

	idx := []int{}
	for i := range v {
		idx = append(idx, c.Intern(v[i]))
	}

	assert.Equal(t, len(v), len(idx))
	assert.Equal(t, len(v), c.count())

	for i := range v {
		assert.Equal(t, v[i], c.get(idx[i]))
	}

	j := c.Intern(Int(9))
	assert.Equal(t, 0, j)
	assert.Equal(t, Int(9), c.get(j))
}
