/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// ITER-0011: value-dispatched (and otherwise non-type) multimethods fall back to
// the runtime vm.MultiFn. They are now EXPLICITLY registered in
// *defmulti-dispatchers* tagged :value/:unknown (visible to the ITER-0012 coverage
// detector) and trampoline BY TAG — the (= :type …) guard in defmulti-call-stmts
// rejects them, so NO native `switch v.(type)` is emitted for them.
//
// Two seams, both anti-hollow:
//   - classify-dispatch is unit-tested directly (pins the :type/:value/:unknown tag
//     a registry entry will carry — not just the downstream emission).
//   - a discrimination lowering test puts a type- AND a value-dispatched multifn in
//     ONE namespace and asserts the type one devirtualises (`_mm_…` + `.(type)`)
//     while the value one does NOT (no `_mm_vdarea` arm) — proving the tag actually
//     gates emission, rather than "everything trampolines" passing vacuously.

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestClassifyDispatchTags(t *testing.T) {
	ensureLoader()

	cases := []struct {
		expr string
		want vm.Keyword
	}{
		// (fn [x] (type x)) → :type (native-eligible)
		{`(ir.lower-go/classify-dispatch (quote (fn [x] (type x))))`, vm.Keyword("type")},
		// named type-dispatch fn is still the type shape
		{`(ir.lower-go/classify-dispatch (quote (fn td [x] (type x))))`, vm.Keyword("type")},
		// keyword dispatch (e.g. (defmulti area :shape)) → analyzable value dispatch
		{`(ir.lower-go/classify-dispatch :shape)`, vm.Keyword("value")},
		// non-type (fn …) form → value dispatch
		{`(ir.lower-go/classify-dispatch (quote (fn [x] (:kind x))))`, vm.Keyword("value")},
		// bare symbol naming a fn defined elsewhere → not statically analyzable
		{`(ir.lower-go/classify-dispatch (quote my-dispatch-fn))`, vm.Keyword("unknown")},
	}
	for _, c := range cases {
		got := runLispExpr(t, c.expr)
		if got != c.want {
			t.Fatalf("classify-dispatch %s = %v, want %v", c.expr, got, c.want)
		}
	}
}

func TestValueDispatchMultimethodTrampolines(t *testing.T) {
	ensureLoader()

	// Intern both multimethods so the callers' builds resolve them.
	runLispExpr(t, `(do
	  (defmulti vdtype (fn [x] (type x)) :default :default)
	  (defmethod vdtype (type [])  [x] :vec)
	  (defmethod vdtype :default    [x] :other)
	  (defmulti vdarea :shape :default :default)
	  (defmethod vdarea :circle [s] :c)
	  (defmethod vdarea :default [s] :o))`)

	v := runLispExpr(t, `(ir.passes.pipeline/lower-ns-to-go "vdpkg" (quote core)
	  (quote [(defmulti vdtype (fn [x] (type x)) :default :default)
	          (defmethod vdtype (type [])  [x] :vec)
	          (defmethod vdtype :default    [x] :other)
	          (defmulti vdarea :shape :default :default)
	          (defmethod vdarea :circle [s] :c)
	          (defmethod vdarea :default [s] :o)
	          (defn call-vdtype [x] (vdtype x))
	          (defn call-vdarea [s] (vdarea s))]))`)

	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go string, got %T: %v", v, v)
	}
	rendered := string(s)

	// The TYPE-dispatched multifn devirtualises: native type-switch + a method arm.
	for _, want := range []string{".(type)", "case vm.ArrayVector:", "_mm_vdtype_0(ec"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("type-dispatch multifn should devirtualise, missing %q:\n%s", want, rendered)
		}
	}
	// The VALUE-dispatched multifn does NOT devirtualise: no lowered method fn, no
	// native arm for it. Its call site trampolines to the runtime vm.MultiFn.
	if strings.Contains(rendered, "_mm_vdarea") {
		t.Fatalf("value-dispatch multifn must NOT emit a native method arm (_mm_vdarea):\n%s", rendered)
	}
}
