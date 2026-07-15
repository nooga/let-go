package vm

import (
	"sync"
	"testing"
)

func TestRootBindDepthGate(t *testing.T) {
	v := NewVar(nil, "t", "x")
	v.SetRoot(Int(1))

	// unbound (also covers non-dynamic): reads root, gate skips the stack
	if got := v.Deref(); got != Int(1) {
		t.Fatalf("unbound deref = %v, want 1", got)
	}
	if rootHasBinding(v) {
		t.Fatalf("unbound: rootHasBinding true, want false")
	}

	// bound: deref sees the binding
	v.PushBinding(Int(2))
	if !rootHasBinding(v) {
		t.Fatalf("bound: rootHasBinding false, want true")
	}
	if got := v.Deref(); got != Int(2) {
		t.Fatalf("bound deref = %v, want 2", got)
	}

	// nested shadow
	v.PushBinding(Int(3))
	if got := v.Deref(); got != Int(3) {
		t.Fatalf("nested deref = %v, want 3", got)
	}
	v.PopBinding()
	if got := v.Deref(); got != Int(2) {
		t.Fatalf("after inner pop deref = %v, want 2 (restored)", got)
	}

	// fully unbound: back to root, counter zero (the PreviouslyBound case)
	v.PopBinding()
	if rootHasBinding(v) {
		t.Fatalf("after unbind: rootHasBinding true, want false")
	}
	if got := v.Deref(); got != Int(1) {
		t.Fatalf("previously-bound deref = %v, want 1 (root)", got)
	}
}

func TestRunWithBindingsCounterConsistent(t *testing.T) {
	v := NewVar(nil, "t", "y")
	v.SetRoot(Int(10))
	v.SetDynamic()

	snap := BindingSnapshot{v: {Int(99)}}
	_, _ = RunWithBindings(snap, func() (Value, error) {
		if got := v.Deref(); got != Int(99) {
			t.Errorf("inside RunWithBindings deref = %v, want 99", got)
		}
		if !rootHasBinding(v) {
			t.Errorf("inside RunWithBindings: rootHasBinding false, want true")
		}
		return NIL, nil
	})
	if got := v.Deref(); got != Int(10) {
		t.Errorf("after RunWithBindings deref = %v, want 10 (root)", got)
	}
	if rootHasBinding(v) {
		t.Errorf("after RunWithBindings: rootHasBinding true, want false")
	}
}

func TestVarDerefChildIsolation(t *testing.T) {
	v := NewVar(nil, "t", "z")
	v.SetRoot(Int(0))
	v.SetDynamic()
	v.PushBinding(Int(1)) // root binding
	defer v.PopBinding()

	child := RootExecContext.Child() // seeded from root snapshot
	if got := child.deref(v); got != Int(1) {
		t.Fatalf("child sees root binding = %v, want 1", got)
	}
	child.pushBinding(v, Int(2)) // child-local
	if got := child.deref(v); got != Int(2) {
		t.Fatalf("child own binding = %v, want 2", got)
	}
	// child binding must NOT touch the root binding state or the root value
	if !rootHasBinding(v) {
		t.Fatalf("child push lost root binding, want true")
	}
	if got := v.Deref(); got != Int(1) {
		t.Fatalf("root deref after child push = %v, want 1 (unaffected)", got)
	}
}

func TestVarDerefConcurrentBindDeref(t *testing.T) {
	const goroutines, iters = 16, 2000
	v := NewVar(nil, "t", "c")
	v.SetRoot(Int(-1))
	v.SetDynamic()

	valid := map[Value]bool{Int(-1): true, Int(7): true} // root or the pushed value
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				if id%2 == 0 {
					if got := v.Deref(); !valid[got] {
						t.Errorf("deref returned unexpected %v", got)
						return
					}
				} else {
					v.PushBinding(Int(7))
					_ = v.Deref()
					v.PopBinding()
				}
			}
		}(g)
	}
	wg.Wait()
	if rootHasBinding(v) {
		t.Errorf("after balanced push/pop: rootHasBinding true, want false")
	}
}

func TestVarDerefDistinctConcurrent(t *testing.T) {
	const goroutines, iters = 16, 2000
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			v := NewVar(nil, "t", "d") // distinct var per goroutine
			v.SetRoot(Int(int64(id)))
			v.SetDynamic()
			for i := 0; i < iters; i++ {
				v.PushBinding(Int(int64(id + 1000)))
				if got := v.Deref(); got != Int(int64(id+1000)) {
					t.Errorf("g%d deref = %v, want %d", id, got, id+1000)
				}
				v.PopBinding()
				if got := v.Deref(); got != Int(int64(id)) {
					t.Errorf("g%d unbound deref = %v, want %d", id, got, id)
				}
			}
		}(g)
	}
	wg.Wait()
}
