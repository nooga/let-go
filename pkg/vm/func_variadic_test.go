package vm

import "testing"

// A multi-arity fn's variadic variant must accept calls with exactly its
// fixed-arg count: (defn g ([a] 1) ([a b & c] c)) → (g 1 2) selects the
// variadic variant with an empty rest, like Clojure. variantFor used to
// require count >= arity INCLUDING the rest slot, rejecting exactly-min.
func TestVariadicMinArity(t *testing.T) {
	fixed1 := MakeFunc(1, false, nil)
	rest := MakeFunc(3, true, nil) // ([a b & c]) — 2 fixed + rest slot

	ma, err := MakeMultiArity([]Value{fixed1, rest})
	if err != nil {
		t.Fatal(err)
	}

	v, err := ma.variantFor(ma, 2)
	if err != nil {
		t.Fatalf("2-arg call must select the variadic variant, got error: %v", err)
	}
	if v != Fn(rest) {
		t.Fatalf("2-arg call selected %v, want the variadic variant", v)
	}

	for _, n := range []int{3, 5} {
		v, err = ma.variantFor(ma, n)
		if err != nil {
			t.Fatalf("%d-arg call: %v", n, err)
		}
		if v != Fn(rest) {
			t.Fatalf("%d-arg call selected %v, want the variadic variant", n, v)
		}
	}

	v, err = ma.variantFor(ma, 1)
	if err != nil {
		t.Fatalf("1-arg call: %v", err)
	}
	if v != Fn(fixed1) {
		t.Fatalf("1-arg call selected %v, want the fixed 1-arity variant", v)
	}

	if _, err = ma.variantFor(ma, 0); err == nil {
		t.Fatal("0-arg call below every arity must error")
	}
}

// A Closure-wrapped variadic Func rest arm follows the same rule.
func TestVariadicMinArityClosureRest(t *testing.T) {
	fixed1 := MakeFunc(1, false, nil)
	rest := &Closure{fn: MakeFunc(3, true, nil)}

	ma, err := MakeMultiArity([]Value{fixed1, rest})
	if err != nil {
		t.Fatal(err)
	}
	v, err := ma.variantFor(ma, 2)
	if err != nil {
		t.Fatalf("2-arg call must select the closure rest arm, got error: %v", err)
	}
	if v != Fn(rest) {
		t.Fatalf("2-arg call selected %v, want the closure rest arm", v)
	}
	if _, err = ma.variantFor(ma, 0); err == nil {
		t.Fatal("0-arg call below every arity must error")
	}
}

// A NativeFn rest arm (NewArityNativeFn) declares its variadic MINIMUM
// directly — its Arity() is the min call width, with no rest slot counted —
// so exactly-min stays accepted and below-min stays rejected.
func TestVariadicMinArityNativeRest(t *testing.T) {
	fixed1 := MakeFunc(1, false, nil)
	nrest := NewArityNativeFn("nrest", 2, true, func(ec *ExecContext, args []Value) (Value, error) {
		return NIL, nil
	})

	ma, err := MakeMultiArity([]Value{fixed1, nrest})
	if err != nil {
		t.Fatal(err)
	}
	v, err := ma.variantFor(ma, 2)
	if err != nil {
		t.Fatalf("2-arg call must select the native rest arm, got error: %v", err)
	}
	if v != Fn(nrest) {
		t.Fatalf("2-arg call selected %v, want the native rest arm", v)
	}
	v, err = ma.variantFor(ma, 1)
	if err != nil {
		t.Fatalf("1-arg call: %v", err)
	}
	if v != Fn(fixed1) {
		t.Fatalf("1-arg call selected %v, want the fixed variant", v)
	}
}
