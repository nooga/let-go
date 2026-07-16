/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// evalTrace compiles and runs src in a fresh compiler over the shared core
// namespace, independent of TestRunner's package-level consts (which may not
// be initialized yet depending on test order).
func evalTrace(t *testing.T, src string) {
	t.Helper()
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	if _, _, err := ctx.CompileMultiple(strings.NewReader(src)); err != nil {
		t.Fatalf("eval %q: %v", src, err)
	}
}

// resetTrace restores the process-global trace gate so this test cannot leak
// an armed gate (and its per-frame TraceVar deref cost) into other tests.
func resetTrace() {
	vm.TraceArmed.Store(false)
	vm.TraceVar.Store(nil)
}

// TestLgTraceArming verifies the *lg-trace* gate wiring the review asked for:
// a falsy binding must not arm, a user-namespace *lg-trace* must not hijack
// the gate, and `(binding [*lg-trace* true] ...)` — the documented entry
// point — must arm via pushBinding.
func TestLgTraceArming(t *testing.T) {
	defer resetTrace()
	resetTrace()

	evalTrace(t, `(binding [*lg-trace* false] 1)`)
	if vm.TraceArmed.Load() {
		t.Fatal("falsy (binding [*lg-trace* false] ...) must not arm tracing")
	}

	// A same-named var in a user namespace must not wire the gate: the match
	// is ns-qualified to core/*lg-trace*.
	evalTrace(t, `
		(ns test.lg-trace-hijack)
		(def ^:dynamic *lg-trace* false)
		(binding [*lg-trace* true] 1)
		(in-ns 'core)`)
	if vm.TraceArmed.Load() || vm.TraceVar.Load() != nil {
		t.Fatal("user-namespace *lg-trace* binding must not arm core tracing")
	}

	// The PR's own documented usage: dynamic binding arms the gate.
	evalTrace(t, `(binding [*lg-trace* true] (+ 1 2))`)
	if !vm.TraceArmed.Load() {
		t.Fatal("(binding [*lg-trace* true] ...) did not arm tracing")
	}
	tv := vm.TraceVar.Load()
	if tv == nil {
		t.Fatal("TraceVar was not wired by the truthy binding")
	}
	if tv.NS() != vm.TraceVarNS || tv.VarName() != vm.TraceVarName {
		t.Fatalf("TraceVar wired to wrong var: %s/%s", tv.NS(), tv.VarName())
	}

	// set! at the root (the other arming path) still works after reset.
	resetTrace()
	evalTrace(t, `(set! *lg-trace* true)`)
	if !vm.TraceArmed.Load() {
		t.Fatal("(set! *lg-trace* true) did not arm tracing")
	}
	evalTrace(t, `(set! *lg-trace* false)`)

	// alter-var-root (which with-redefs is built on) is the third mutation
	// path and must arm as well.
	resetTrace()
	evalTrace(t, `(alter-var-root #'*lg-trace* (constantly true))`)
	if !vm.TraceArmed.Load() {
		t.Fatal("(alter-var-root #'*lg-trace* (constantly true)) did not arm tracing")
	}
	evalTrace(t, `(alter-var-root #'*lg-trace* (constantly false))`)
}

// TestLgTraceShadowingRegression verifies that a user-defined trace function
// is not shadowed by the VM's trace implementation (the regression this PR
// fixes). A mismatch THROWS, which surfaces through evalTrace's error check —
// the assertion is load-bearing, not a printed diagnostic.
func TestLgTraceShadowingRegression(t *testing.T) {
	defer resetTrace()
	resetTrace()

	evalTrace(t, `
		(defn trace [s]
			(with-out-str
				(println s)))
		(let [output (trace "Hello")]
			(when-not (= output "Hello\n")
				(throw (str "user-defined trace was shadowed; got: " (pr-str output)))))
	`)
}

// captureStdout runs body with os.Stdout redirected to a pipe and returns
// everything written. The pipe is drained concurrently so trace output larger
// than the kernel pipe buffer cannot deadlock the writer, and nothing is
// truncated (a fixed-size single Read would be).
func captureStdout(t *testing.T, body func() error) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	oldStdout := os.Stdout
	os.Stdout = w
	done := make(chan string, 1)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()
	bodyErr := body()
	os.Stdout = oldStdout
	w.Close()
	out := <-done
	r.Close()
	if bodyErr != nil {
		t.Fatalf("eval failed: %v", bodyErr)
	}
	return out
}

// TestLgTraceOutputAndPropagation verifies that *lg-trace* binding produces
// trace output, that it propagates into CALLED functions (not just the frame
// the binding wraps), and that after the binding exits the still-armed gate
// stays silent because the per-frame *lg-trace* deref is falsy.
func TestLgTraceOutputAndPropagation(t *testing.T) {
	defer resetTrace()
	resetTrace()

	// inner is called with a sentinel argument only IT receives, so its
	// frame-entry line ("run[42077]") is distinguishable from outer's
	// ("run[5]") — outer alone tracing must not pass this test. The call must
	// be NON-tail: a tail call reuses the caller's frame (tracing continues,
	// but no new frame-entry banner is printed), so a tail-positioned inner
	// would emit no run[42077] line even with propagation working.
	output := captureStdout(t, func() error {
		ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
		_, _, err := ctx.CompileMultiple(strings.NewReader(`
			(defn inner [x]
				(+ x 1))
			(defn outer [x]
				(+ 1 (inner 42077)))
			(binding [*lg-trace* true]
				(outer 5))
		`))
		return err
	})

	if !strings.Contains(output, "run[5]") {
		t.Fatalf("trace output missing outer's frame entry run[5]:\n%s", output)
	}
	if !strings.Contains(output, "run[42077]") {
		t.Fatalf("binding did not propagate into the callee: no run[42077] frame entry for inner:\n%s", output)
	}
	if !strings.Contains(output, "LOAD_ARG") {
		t.Fatalf("trace output missing LOAD_ARG opcode:\n%s", output)
	}

	// Post-binding silence: do NOT reset the armed state. TraceArmed stays
	// true and TraceVar stays wired — exactly the state after a real binding
	// scope exits — so what must suppress tracing here is the per-frame
	// *lg-trace* VALUE deref resolving falsy, not the coarse gate.
	if !vm.TraceArmed.Load() || vm.TraceVar.Load() == nil {
		t.Fatal("precondition: trace gate should remain armed after the binding exits")
	}
	output2 := captureStdout(t, func() error {
		ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
		_, _, err := ctx.CompileMultiple(strings.NewReader(`
			(defn simple [x]
				(+ x 1))
			(simple 10)
		`))
		return err
	})

	if strings.Contains(output2, "run[") || strings.Contains(output2, "LOAD_ARG") {
		t.Fatalf("armed-but-unbound tracing must stay silent; got:\n%s", output2)
	}
}
