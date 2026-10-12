package vm

import (
	"sync"
	"testing"
)

// The first evaluations of a capture-free multi-arity literal can run on
// several goroutines at once; all of them must get the one shared fn, or
// identical? is false between them.
func TestMultiArityLiteralConcurrentFirstUseIsOneObject(t *testing.T) {
	const goroutines = 64
	chunk := NewCodeChunk(NewConsts())
	for trial := range 50 {
		fns := []Value{MakeFunc(1, false, chunk), MakeFunc(2, false, chunk)}
		got := make([]*MultiArityFn, goroutines)
		var start, done sync.WaitGroup
		start.Add(1)
		for i := range got {
			done.Go(func() {
				start.Wait()
				ma, err := chunk.multiArityAt(fns)
				if err != nil {
					t.Error(err)
					return
				}
				got[i] = ma
			})
		}
		start.Done()
		done.Wait()
		for i, ma := range got {
			if ma != got[0] {
				t.Fatalf("trial %d: goroutine %d got %p, goroutine 0 got %p", trial, i, ma, got[0])
			}
		}
	}
}
