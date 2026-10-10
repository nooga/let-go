package vm

import (
	"sync"
	"testing"
)

// TestInternConcurrentFirstCallsShareOneVar pins Intern's atomic interning:
// goroutines racing to intern a name the namespace does not have yet must all
// get the Var that ends up in the registry. A lookup-then-Def split lets a
// later Def replace an earlier caller's Var, stranding code compiled against
// it. The race detector cannot see that logical race, so the test checks Var
// identity over many rounds.
func TestInternConcurrentFirstCallsShareOneVar(t *testing.T) {
	const rounds = 500
	const workers = 8

	for r := 0; r < rounds; r++ {
		ns := NewNamespace("intern.race")
		got := make([]*Var, workers)
		start := make(chan struct{})
		var wg sync.WaitGroup
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func(w int) {
				defer wg.Done()
				<-start
				got[w] = ns.Intern("x", Int(int64(w)))
			}(w)
		}
		close(start)
		wg.Wait()

		want := ns.LookupLocal("x")
		for w, va := range got {
			if va != want {
				t.Fatalf("round %d: worker %d got a Var that is not the registered one", r, w)
			}
		}
	}
}

// TestInternKeepsVarFlagsAndMeta pins that Intern only replaces the root. A
// source-level (def ...) re-applies definition metadata; Intern does not.
func TestInternKeepsVarFlagsAndMeta(t *testing.T) {
	ns := NewNamespace("intern.flags")
	va := ns.Def("x", Int(1))
	meta := NewPersistentMap([]Value{Keyword("custom"), TRUE})
	va.SetMeta(meta)
	va.SetDynamic()
	va.SetPrivate()

	if got := ns.Intern("x", Int(2)); got != va {
		t.Fatalf("Intern replaced the existing Var")
	}
	if va.Root() != Int(2) {
		t.Fatalf("root = %v, want 2", va.Root())
	}
	if !va.IsDynamic() || !va.IsPrivate() {
		t.Fatalf("flags lost: dynamic=%v private=%v", va.IsDynamic(), va.IsPrivate())
	}
	if va.Meta() != meta {
		t.Fatalf("meta = %v, want %v", va.Meta(), meta)
	}
}

// TestInternTracksGuardedRoot pins that Intern goes through SetRoot, so a
// guarded native-primitive root counts as deviated while overridden and as
// intact once restored.
func TestInternTracksGuardedRoot(t *testing.T) {
	ns := NewNamespace("intern.guard")
	canonical, err := NativeFnType.WrapNoErr(func([]Value) Value { return NIL })
	if err != nil {
		t.Fatal(err)
	}
	ns.Def("prim", canonical).GuardRoot()
	base := guardedRootDeviations.Load()

	ns.Intern("prim", Int(1))
	if got := guardedRootDeviations.Load() - base; got != 1 {
		t.Fatalf("deviations after override = %+d, want +1", got)
	}
	ns.Intern("prim", canonical)
	if got := guardedRootDeviations.Load() - base; got != 0 {
		t.Fatalf("deviations after restore = %+d, want 0", got)
	}
}

// TestInternShadowsReferredVarWithoutMutatingIt pins that Intern only reuses
// a Var the namespace owns. A name referred in from another namespace gets a
// new local Var; the referred Var keeps its root.
func TestInternShadowsReferredVarWithoutMutatingIt(t *testing.T) {
	lib := NewNamespace("intern.lib")
	libVar := lib.Def("x", Int(1))
	user := NewNamespace("intern.user")
	user.Refer(lib, "", true)

	local := user.Intern("x", Int(2))
	if local == libVar {
		t.Fatalf("Intern reused the referred Var")
	}
	if libVar.Root() != Int(1) {
		t.Fatalf("referred root = %v, want 1", libVar.Root())
	}
	if user.Lookup("x") != local {
		t.Fatalf("unqualified x does not resolve to the local Var")
	}
}
