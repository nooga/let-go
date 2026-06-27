/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// ITER-0010: native Go lowering of type-dispatched defmulti/defmethod. A call to
// a same-ns type-dispatched multimethod must devirtualise to a native
// `switch dispatchArg.(type) { case <GoType>: <method-fn>(…) … default: <runtime
// vm.MultiFn> }` — NOT a runtime hash-dispatch trampoline. This is the anti-stub
// proof: it asserts the rendered Go actually contains the type-switch with method
// arms and a runtime default. (test/defmulti_native_lowering_test.lg pins the
// matching runtime semantics under the bytecode VM.)

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestDefmultiNativeTypeSwitchLowering(t *testing.T) {
	ensureLoader()

	// Intern the multimethod so the caller's build resolves it (a real defmulti,
	// not just a form in the list).
	runLispExpr(t, `(do
	  (defmulti tdnative (fn [x] (type x)) :default :default)
	  (defmethod tdnative (type [])     [x] :vec)
	  (defmethod tdnative (type (list)) [x] :lst)
	  (defmethod tdnative :default      [x] :other))`)

	v := runLispExpr(t, `(ir.passes.pipeline/lower-ns-to-go "tdpkg" (quote core)
	  (quote [(defmulti tdnative (fn [x] (type x)) :default :default)
	          (defmethod tdnative (type [])     [x] :vec)
	          (defmethod tdnative (type (list)) [x] :lst)
	          (defmethod tdnative :default      [x] :other)
	          (defn call-tdnative [x] (tdnative x))]))`)

	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go string, got %T: %v", v, v)
	}
	rendered := string(s)

	// Native type-switch with both concrete arms calling the lowered method fns.
	for _, want := range []string{
		".(type)",
		"case vm.ArrayVector:",
		"case *vm.List:",
		"_mm_tdnative_0(ec",
		"_mm_tdnative_1(ec",
		"default:",
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("native multimethod lowering missing %q:\n--- rendered Go ---\n%s", want, rendered)
		}
	}
	// The default arm is the safety net: it must trampoline to the runtime
	// vm.MultiFn (cached-var IFn dispatch), never silently drop unmatched types.
	if !strings.Contains(rendered, "CachedVarFn") && !strings.Contains(rendered, "InvokeValueEC") {
		t.Fatalf("expected a runtime-MultiFn default arm (CachedVarFn/InvokeValueEC):\n%s", rendered)
	}
	// And it must NOT be the old pure-trampoline shape (no type-switch at all).
	if !strings.Contains(rendered, "switch") {
		t.Fatalf("expected a native switch, got a pure trampoline:\n%s", rendered)
	}
	// P1 guard (#326): the native arms are only valid while the multifn var
	// still holds its frozen native baseline. The switch must be gated by
	// rt.MultiFnNativeFrozen so a late defmethod falls back to runtime dispatch,
	// and the ns must register the multifn for freezing at load.
	if !strings.Contains(rendered, `rt.MultiFnNativeFrozen(&__v_clojure_core_tdnative, "clojure.core", "tdnative")`) {
		t.Fatalf("native type-switch must be guarded by rt.MultiFnNativeFrozen:\n%s", rendered)
	}
	if !strings.Contains(rendered, `rt.RegisterNativeMultiFns("core", []string{"tdnative"})`) {
		t.Fatalf("ns must register its native multimethods for load-time freezing:\n%s", rendered)
	}
}

// P1 regression: repeated defmethod definitions for the same dispatch value are
// last-write-wins at runtime (the later method replaces the earlier in the
// MultiFn map). The native type-switch must collapse them to a SINGLE
// `case <GoType>:` arm calling the latest method — two arms would (a) be invalid
// Go (duplicate case) and (b) let an older/stale-arity arm win over the runtime's
// latest method.
func TestDefmultiNativeDuplicateDispatchCollapsed(t *testing.T) {
	ensureLoader()

	runLispExpr(t, `(do
	  (defmulti tddup (fn [x] (type x)) :default :default)
	  (defmethod tddup (type []) [x] :first)
	  (defmethod tddup (type []) [x] :second)
	  (defmethod tddup :default  [x] :other))`)

	v := runLispExpr(t, `(ir.passes.pipeline/lower-ns-to-go "tdduppkg" (quote core)
	  (quote [(defmulti tddup (fn [x] (type x)) :default :default)
	          (defmethod tddup (type []) [x] :first)
	          (defmethod tddup (type []) [x] :second)
	          (defmethod tddup :default  [x] :other)
	          (defn call-tddup [x] (tddup x))]))`)

	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go string, got %T: %v", v, v)
	}
	rendered := string(s)

	// Exactly one ArrayVector arm — the duplicate must be collapsed.
	if n := strings.Count(rendered, "case vm.ArrayVector:"); n != 1 {
		t.Fatalf("expected exactly 1 `case vm.ArrayVector:` arm (last-wins collapse), got %d:\n%s", n, rendered)
	}
	// Last-wins: the surviving arm is the SECOND method (_mm_tddup_1), not the first.
	if !strings.Contains(rendered, "_mm_tddup_1(ec") {
		t.Fatalf("collapsed arm must call the LAST method (_mm_tddup_1):\n%s", rendered)
	}
	if strings.Contains(rendered, "_mm_tddup_0(ec") {
		t.Fatalf("stale first method (_mm_tddup_0) must not be dispatched after a later redefinition:\n%s", rendered)
	}
}
