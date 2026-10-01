package vm

import (
	"sync"
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

func TestConstpoolDistinguishesEmptyListFromNilList(t *testing.T) {
	c := NewConsts()
	nilList := NewList([]Value{NIL})

	emptyIdx := c.Intern(EmptyList)
	nilListIdx := c.Intern(nilList)

	assert.NotEqual(t, emptyIdx, nilListIdx)
	assert.Equal(t, EmptyList, c.get(emptyIdx))
	assert.Equal(t, nilList, c.get(nilListIdx))
}

// A collection constant is shared by every goroutine that runs its chunk, so
// interning it must fill its hash cache: run with -race.
func TestConstpoolPrehashesCollectionConstants(t *testing.T) {
	c := NewConsts()
	v := NewList([]Value{Int(1), Int(2)})
	c.Intern(v)

	start := make(chan struct{})
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			HashValue(v)
		}()
	}
	close(start)
	wg.Wait()
}
