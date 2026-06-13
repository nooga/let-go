/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"sync"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestPmapvPreservesOrder(t *testing.T) {
	sq, _ := vm.NativeFnType.Wrap(func(a []vm.Value) (vm.Value, error) {
		n := int(a[0].(vm.Int))
		return vm.Int(n * n), nil
	})
	coll := vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3), vm.Int(4), vm.Int(5)})
	r, err := parallelMapV(vm.RootExecContext, []vm.Value{sq, coll})
	if err != nil {
		t.Fatalf("pmapv: %v", err)
	}
	got := r.(vm.ArrayVector)
	want := []int{1, 4, 9, 16, 25}
	if len(got) != len(want) {
		t.Fatalf("expected %d results, got %d", len(want), len(got))
	}
	for i, w := range want {
		if int(got[i].(vm.Int)) != w {
			t.Fatalf("at %d: want %d, got %v", i, w, got[i])
		}
	}
}

// TestPmapvSharesCallerBindings: the workers run synchronously while the
// caller's dynamic binding is live, so they read it via Var.Deref without
// any per-task snapshot.
func TestPmapvSharesCallerBindings(t *testing.T) {
	v := vm.NewVar(nil, "test", "*scale*")
	v.SetRoot(vm.Int(1))
	v.PushBinding(vm.Int(10))
	defer v.PopBinding()

	mul, _ := vm.NativeFnType.Wrap(func(a []vm.Value) (vm.Value, error) {
		return vm.Int(int(a[0].(vm.Int)) * int(v.Deref().(vm.Int))), nil
	})
	r, err := parallelMapV(vm.RootExecContext, []vm.Value{mul, vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)})})
	if err != nil {
		t.Fatalf("pmapv: %v", err)
	}
	got := r.(vm.ArrayVector)
	want := []int{10, 20, 30}
	for i, w := range want {
		if int(got[i].(vm.Int)) != w {
			t.Fatalf("worker did not see caller binding at %d: want %d, got %v", i, w, got[i])
		}
	}
}

func TestPmapvEmpty(t *testing.T) {
	id, _ := vm.NativeFnType.Wrap(func(a []vm.Value) (vm.Value, error) { return a[0], nil })
	r, err := parallelMapV(vm.RootExecContext, []vm.Value{id, vm.NIL})
	if err != nil {
		t.Fatalf("pmapv nil: %v", err)
	}
	if got := r.(vm.ArrayVector); len(got) != 0 {
		t.Fatalf("expected empty result for nil coll, got %d", len(got))
	}
}

// TestNextIDConcurrentUnique pins the atomic gensym counter: concurrent
// callers must each get a distinct id and there must be no data race
// (run under -race). The old non-atomic gensymID++ both raced and could
// hand out duplicates.
func TestNextIDConcurrentUnique(t *testing.T) {
	const n = 2000
	out := make([]int, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			out[i] = nextID()
		}(i)
	}
	wg.Wait()
	seen := make(map[int]struct{}, n)
	for _, x := range out {
		if _, dup := seen[x]; dup {
			t.Fatalf("duplicate gensym id %d under concurrency", x)
		}
		seen[x] = struct{}{}
	}
}
