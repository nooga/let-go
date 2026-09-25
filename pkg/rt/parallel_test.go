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

// boundFn builds a bound-fn* wrapper over fn, capturing whatever v is
// currently bound to on ec (v must already have an active push via
// v.PushBinding for anything to be captured).
func boundFn(t *testing.T, ec *vm.ExecContext, fn vm.Fn) vm.Fn {
	t.Helper()
	bfVar := NS(NameCoreNS).Lookup(vm.Symbol("bound-fn*"))
	if bfVar == nil {
		t.Fatal("core/bound-fn* not found")
	}
	bfStar, ok := bfVar.(*vm.Var).Deref().(vm.Fn)
	if !ok {
		t.Fatal("core/bound-fn* is not an Fn")
	}
	wrapped, err := ec.Invoke(bfStar, []vm.Value{fn})
	if err != nil {
		t.Fatalf("bound-fn*: %v", err)
	}
	f, ok := wrapped.(vm.Fn)
	if !ok {
		t.Fatal("bound-fn* did not return an Fn")
	}
	return f
}

// TestPmapvBoundFnWorkersDoNotRaceSharedContext pins the pmapv/bound-fn*
// interaction: pmapv deliberately shares ONE ec across every worker
// goroutine (ec.Bind, unlike future*/go* which each get a private child
// context). bound-fn* used to push/pop its captured bindings directly on
// that shared ec; pops are LIFO per var, not per goroutine, so one
// worker's pop could remove another worker's still-live frame. Each
// worker here captures a DISTINCT value for *scale* at bound-fn creation
// time, so cross-contamination between workers is directly observable
// (not just "same value happened to survive"). Run under -race: the fix
// (push/pop on a private ec.Child(), not the shared callEc) must also be
// data-race free.
func TestPmapvBoundFnWorkersDoNotRaceSharedContext(t *testing.T) {
	v := vm.NewVar(nil, "test", "*scale*")
	v.SetRoot(vm.Int(-1))

	const n = 64
	fns := make([]vm.Fn, n)
	for i := 0; i < n; i++ {
		v.PushBinding(vm.Int(i))
		// Ctx-aware: v.Deref() (no args) always reads the ROOT binding only
		// (pkg/vm/var.go) and would never see a child ExecContext's pushed
		// binding at all, silently passing regardless of the bug under test.
		// Real .lg code reads a dynamic var via the ec-aware path (the
		// bytecode VM's LOAD_VAR resolves through ec.Deref); reader must go
		// through that same path to actually exercise it.
		reader := vm.NewCtxNativeFn("", func(callEc *vm.ExecContext, _ []vm.Value) (vm.Value, error) {
			return callEc.Deref(v), nil
		})
		fns[i] = boundFn(t, vm.RootExecContext, reader)
		v.PopBinding()
	}

	invoker, _ := vm.NativeFnType.Wrap(func(a []vm.Value) (vm.Value, error) {
		i := int(a[0].(vm.Int))
		return vm.RootExecContext.Invoke(fns[i], nil)
	})
	coll := make([]vm.Value, n)
	for i := 0; i < n; i++ {
		coll[i] = vm.Int(i)
	}
	r, err := parallelMapV(vm.RootExecContext, []vm.Value{invoker, vm.NewArrayVector(coll)})
	if err != nil {
		t.Fatalf("pmapv: %v", err)
	}
	got := r.(vm.ArrayVector)
	for i := 0; i < n; i++ {
		if int(got[i].(vm.Int)) != i {
			t.Fatalf("worker %d saw captured *scale*=%v, want %d (shared-context corruption)", i, got[i], i)
		}
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
