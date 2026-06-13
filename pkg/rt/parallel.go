/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/nooga/let-go/pkg/vm"
)

// parallelMapV is `pmapv` — an EAGER, ordered parallel map that shares
// the caller's dynamic-var bindings rather than snapshotting them.
//
// `pmap` (the Lisp one) is built on `future`, and every future pays the
// global bindings mutex twice plus an O(active-bindings) copy in
// SnapshotBindings + RunWithBindings — because a future is async and may
// run after the caller's `binding` scope has unwound, so it must capture
// and reinstate that scope. That spawn cost serializes parallel work.
//
// pmapv is *synchronous*: it blocks until every element is done, so the
// caller's `binding` scope is still live throughout. The worker
// goroutines therefore just read the current global dynamic bindings via
// the (lock-free) Var.Deref — no snapshot, no per-task bindings mutex.
//
// Caveat: because workers share the caller's binding stack, a worker fn
// that PUSHES its own dynamic binding (a `binding` form) on a var would
// race the shared global stack. pmapv is for read-only-binding workloads
// (e.g. lowering a namespace's defns, which only read *ns*/*target*).
func parallelMapV(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
	if len(vs) != 2 {
		return vm.NIL, fmt.Errorf("pmapv expects 2 args (fn coll)")
	}
	fn, ok := vs[0].(vm.Fn)
	if !ok {
		return vm.NIL, fmt.Errorf("pmapv expected Fn as first arg")
	}
	// Convey the caller's context to the workers (Clojure pmap conveys).
	// Reads of the shared binding stack are safe; the push-a-binding caveat
	// below is unchanged.
	fn = ec.Bind(fn)
	if vs[1] == vm.NIL {
		return vm.NewArrayVector(nil), nil
	}
	seq, ok := vs[1].(vm.Sequable)
	if !ok {
		return vm.NIL, fmt.Errorf("pmapv expected a seqable second arg")
	}

	var items []vm.Value
	for s := seq.Seq(); s != nil; s = s.Next() {
		items = append(items, s.First())
	}
	n := len(items)
	if n == 0 {
		return vm.NewArrayVector(nil), nil
	}

	results := make([]vm.Value, n)
	errs := make([]error, n)

	workers := runtime.NumCPU()
	if workers > n {
		workers = n
	}

	var next int64 = -1
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				i := int(atomic.AddInt64(&next, 1))
				if i >= n {
					return
				}
				r, err := fn.Invoke([]vm.Value{items[i]})
				results[i] = r
				errs[i] = err
			}
		}()
	}
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return vm.NIL, err
		}
	}
	return vm.NewArrayVector(results), nil
}
