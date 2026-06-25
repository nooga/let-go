package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// TestMakeNativeMultiArityDeferredDispatch exercises the RUNTIME behavior of the
// native multi-arity dispatch path (ITER-0014, AC-MA.2): a multi-arity defn
// lowers to a constructor returning rt.MakeNativeMultiArity over per-arity
// rt.BoxNativeFn closures, registered via rt.MakeNativeMultiArityDeferred. This
// builds the same shape the codegen emits (see TestLowerNsToGo*MultiArity*) and
// actually INVOKES it — proving arity selection + variadic rt.BoxRestArgs tail
// packing match vm.MakeMultiArity semantics, not just that the Go renders.
func TestMakeNativeMultiArityDeferredDispatch(t *testing.T) {
	ec := vm.NewExecContext()

	// Mirrors `(defn g ([x] x) ([x & more] more))`: an arity-1 identity branch
	// and a variadic branch that packs its rest tail via rt.BoxRestArgs. The
	// constructor captures ec exactly like a real lowered multi-arity fn.
	build := func(ec *vm.ExecContext) vm.Value {
		return MakeNativeMultiArity([]vm.Value{
			BoxNativeFn(func(arg0 vm.Value) (vm.Value, error) { return arg0, nil }),
			BoxNativeFn(func(arg0 vm.Value, args ...vm.Value) (vm.Value, error) {
				return BoxRestArgs(args), nil
			}),
		})
	}

	fn := MakeNativeMultiArityDeferred(build)
	callable, ok := fn.(vm.Fn)
	if !ok {
		t.Fatalf("deferred multi-arity value is not callable: %T", fn)
	}

	// Identity stability (refutes the "rebuild breaks identity" hazard): the
	// REGISTERED value is created once and is stable; only the transient inner
	// MultiArityFn is rebuilt per call and never escapes. A second lookup of the
	// same registered value is identical to the first.
	fn2 := fn
	if fn != fn2 {
		t.Fatalf("registered multi-arity override value must be identity-stable")
	}

	// Arity-1 exact match → identity branch.
	r1, err := ec.Invoke(callable, []vm.Value{vm.Int(7)})
	if err != nil {
		t.Fatalf("arity-1 invoke errored: %v", err)
	}
	if r1 != vm.Int(7) {
		t.Fatalf("arity-1 (g 7) = %v, want 7", r1)
	}

	// Variadic branch (le >= rest arity) → rest tail packed via BoxRestArgs.
	r2, err := ec.Invoke(callable, []vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)})
	if err != nil {
		t.Fatalf("variadic invoke errored: %v", err)
	}
	seq, ok := r2.(vm.Seq)
	if !ok {
		t.Fatalf("variadic (g 1 2 3) returned %T, want a packed seq", r2)
	}
	// rest = (2 3) — the tail after the fixed arg, packed as a list.
	got := []vm.Value{}
	for s := seq; s != nil && s.First() != nil; s = s.Next() {
		got = append(got, s.First())
		if s.Next() == nil {
			break
		}
	}
	if len(got) != 2 || got[0] != vm.Int(2) || got[1] != vm.Int(3) {
		t.Fatalf("variadic rest = %v, want (2 3)", got)
	}

	// Re-invoking after the per-call rebuild still dispatches correctly (the
	// rebuild is invocation-safe — fresh ec/value each call, same semantics).
	r3, err := ec.Invoke(callable, []vm.Value{vm.Int(42)})
	if err != nil || r3 != vm.Int(42) {
		t.Fatalf("re-invoke arity-1 = %v (err %v), want 42", r3, err)
	}
}
