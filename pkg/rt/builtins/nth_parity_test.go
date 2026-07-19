/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package builtins

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// emptyLazy returns a *vm.LazySeq that realizes to empty (nil head).
func emptyLazy() vm.Value {
	f, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) { return vm.NIL, nil })
	return vm.NewLazySeq(f.(vm.Fn))
}

// TestNthMatchesCanonical guards #613's builtins.Nth/Nth3 (the AOT direct-call
// targets for nth) against VM/AOT semantic drift on the edge cases the reviewer
// flagged. If these diverge from clojure.core/nth, the direct-linking
// optimization silently changes results between the VM and AOT paths.
func TestNthMatchesCanonical(t *testing.T) {
	vec := vm.NewArrayVector([]vm.Value{vm.Int(1)})

	// (nth nil 0) => nil, NOT an error (nil is an empty coll for nth).
	if v, err := Nth(vm.NIL, vm.Int(0)); err != nil || v != vm.NIL {
		t.Errorf("Nth(nil, 0) = %v, %v; want nil with no error", v, err)
	}

	// (nth [1] :bad :missing) => THROW on the non-numeric index, NOT :missing —
	// not-found is only for out-of-bounds, not a bad index type.
	if _, err := Nth3(vec, vm.Keyword("bad"), vm.Keyword("missing")); err == nil {
		t.Error("Nth3([1], :bad, :missing) returned no error; want a throw on the non-numeric index")
	}
	// The 2-arg form already threw here; keep it pinned.
	if _, err := Nth(vec, vm.Keyword("bad")); err == nil {
		t.Error("Nth([1], :bad) returned no error; want a throw on the non-numeric index")
	}

	// An EMPTY lazy seq must be forced: (nth empty-lazy 0) is out of bounds, and
	// the 3-arg form yields the default — not a nil First() off an unforced seq.
	if _, err := Nth(emptyLazy(), vm.Int(0)); err == nil {
		t.Error("Nth(empty-lazy, 0) returned no error; want out-of-bounds")
	}
	if v, err := Nth3(emptyLazy(), vm.Int(0), vm.Keyword("missing")); err != nil || v != vm.Keyword("missing") {
		t.Errorf("Nth3(empty-lazy, 0, :missing) = %v, %v; want :missing", v, err)
	}
}
