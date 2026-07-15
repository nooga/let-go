package vm

import (
	"sync"
	"testing"
)

func TestRootBindingsHelpers(t *testing.T) {
	v := NewVar(nil, "t", "x")
	v.SetRoot(Int(1))

	if _, ok := rootDerefHead(v); ok {
		t.Fatal("unbound: rootDerefHead ok=true, want false")
	}
	if rootHasBinding(v) {
		t.Fatal("unbound: rootHasBinding true")
	}

	rootPush(v, Int(2))
	if val, ok := rootDerefHead(v); !ok || val != Int(2) {
		t.Fatalf("bound head = %v,%v want 2,true", val, ok)
	}
	rootPush(v, Int(3)) // nested shadow
	if val, _ := rootDerefHead(v); val != Int(3) {
		t.Fatalf("nested head = %v want 3", val)
	}
	if !rootSetCurrent(v, Int(30)) {
		t.Fatal("setCurrent on bound returned false")
	}
	if val, _ := rootDerefHead(v); val != Int(30) {
		t.Fatalf("after set! head = %v want 30", val)
	}
	rootPop(v)
	if val, _ := rootDerefHead(v); val != Int(2) {
		t.Fatalf("after pop head = %v want 2 (restored)", val)
	}

	// snapshot captures live bindings top-last
	snap := rootSnapshot()
	if got := snap[v]; len(got) != 1 || got[0] != Int(2) {
		t.Fatalf("snapshot[v] = %v want [2]", got)
	}

	rootPop(v)
	if rootHasBinding(v) {
		t.Fatal("after full unbind: still bound")
	}
	if _, ok := rootSnapshot()[v]; ok {
		t.Fatal("after full unbind: still in snapshot/registry")
	}
	if rootSetCurrent(v, Int(9)) {
		t.Fatal("setCurrent on unbound returned true")
	}
}

func TestRootInstallReplacesState(t *testing.T) {
	a := NewVar(nil, "t", "a")
	a.SetRoot(Int(0))
	b := NewVar(nil, "t", "b")
	b.SetRoot(Int(0))
	rootPush(a, Int(1)) // a bound before install

	rootInstall(BindingSnapshot{b: {Int(7), Int(8)}}) // install: b bound (2 deep), a cleared
	if rootHasBinding(a) {
		t.Fatal("install: a should be cleared")
	}
	if val, _ := rootDerefHead(b); val != Int(8) {
		t.Fatalf("install: b head = %v want 8 (top)", val)
	}
	if _, ok := rootSnapshot()[a]; ok {
		t.Fatal("install: a still in registry")
	}
	if got := rootSnapshot()[b]; len(got) != 2 {
		t.Fatalf("install: b depth = %d want 2", len(got))
	}
	rootInstall(BindingSnapshot{}) // clear all
	if rootHasBinding(b) {
		t.Fatal("clear: b still bound")
	}
}

func TestRootPathWiredEndToEnd(t *testing.T) {
	v := NewVar(nil, "t", "e")
	v.SetRoot(Int(1))
	v.SetDynamic()

	if v.Deref() != Int(1) {
		t.Fatal("unbound deref != root")
	}
	v.PushBinding(Int(2))
	if v.Deref() != Int(2) {
		t.Fatal("bound deref != 2")
	}
	// child seeded from root must inherit the root binding
	child := RootExecContext.Child()
	if child.deref(v) != Int(2) {
		t.Fatalf("child deref = %v want 2 (inherited)", child.deref(v))
	}
	// child's own binding does not leak to root
	child.pushBinding(v, Int(3))
	if v.Deref() != Int(2) {
		t.Fatal("root deref changed by child binding")
	}
	if child.deref(v) != Int(3) {
		t.Fatal("child deref != own binding")
	}
	v.PopBinding()
	if v.Deref() != Int(1) {
		t.Fatal("after unbind deref != root")
	}
	// RunWithBindings brackets a root value
	_, _ = RunWithBindings(BindingSnapshot{v: {Int(42)}}, func() (Value, error) {
		if v.Deref() != Int(42) {
			t.Error("inside RunWithBindings deref != 42")
		}
		return NIL, nil
	})
	if v.Deref() != Int(1) {
		t.Fatal("after RunWithBindings deref != root")
	}
}

func TestRootBindingsConcurrent(t *testing.T) {
	const goroutines, iters = 16, 2000
	v := NewVar(nil, "t", "c")
	v.SetRoot(Int(-1))
	v.SetDynamic()
	valid := map[Value]bool{Int(-1): true, Int(7): true}

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				if id%3 == 0 {
					if got := v.Deref(); !valid[got] {
						t.Errorf("deref returned %v", got)
						return
					}
				} else if id%3 == 1 {
					v.PushBinding(Int(7))
					v.PopBinding()
				} else {
					_ = SnapshotBindings() // registry read racing binds
				}
			}
		}(g)
	}
	wg.Wait()
	if rootHasBinding(v) {
		t.Errorf("after balanced push/pop, still bound")
	}
	if _, ok := SnapshotBindings()[v]; ok {
		t.Errorf("after balanced push/pop, still in registry")
	}
}

func TestRootBindingsDistinctConcurrent(t *testing.T) {
	const goroutines, iters = 16, 2000
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			v := NewVar(nil, "t", "d")
			v.SetRoot(Int(int64(id)))
			v.SetDynamic()
			for i := 0; i < iters; i++ {
				v.PushBinding(Int(int64(id + 1000)))
				if v.Deref() != Int(int64(id+1000)) {
					t.Errorf("g%d bound deref wrong", id)
				}
				v.PopBinding()
				if v.Deref() != Int(int64(id)) {
					t.Errorf("g%d unbound deref wrong", id)
				}
			}
		}(g)
	}
	wg.Wait()
}
