/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// #554: a fn whose body lowers to a fully typed signature — the canonical case
// being a loop that infers an int result, `func SumLoop(ec) int` — used to
// register no override at all, because override-uniform-value? demanded
// vm.Value in AND out. The lowered code was emitted and then never reached: the
// var kept its bytecode fn, so every dynamic call ran the VM.
//
// Measured on benchmark/gloat/loop-recur.clj, where the loop sits inside -main:
// the AOT binary ran at 60.0ms against the VM's 60.5ms (no benefit at all),
// while the same program with a boxed return ran 6.9ms — 8.7x, from lowering
// that was already happening.
func TestTypedReturnRegistersBoxedOverride(t *testing.T) {
	ensureLoader()

	v := runLispExpr(t, `(do (create-ns (quote tret))
      (ir.passes.pipeline/lower-ns-to-go "tret" (quote tret)
        [(quote (defn sum-loop [] (loop [i 0 acc 0] (if (< i 1000000) (recur (inc i) (+ acc i)) acc))))]))`)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go source string, got %T", v)
	}
	src := string(s)

	// Precondition: the body really did lower to a typed int return. If this
	// stops holding the test is no longer exercising #554.
	if !regexp.MustCompile(`func SumLoop\(ec \*vm\.ExecContext\) int`).MatchString(src) {
		t.Fatalf("expected a typed `func SumLoop(ec *vm.ExecContext) int`:\n%s", src)
	}

	// The fix: it registers, and the wrapper boxes the typed result.
	if !strings.Contains(src, "RegisterGoOverrides") {
		t.Fatalf("typed-return fn emitted no init()/RegisterGoOverrides — it is unreachable as an override:\n%s", src)
	}
	if !regexp.MustCompile(`"sum-loop":\s*__gogen_wrap`).MatchString(src) {
		t.Fatalf("sum-loop is NOT in the override map:\n%s", src)
	}
	if !strings.Contains(src, "return vm.Int(SumLoop(ec)), nil") {
		t.Fatalf("expected the wrapper to box the typed return as vm.Int:\n%s", src)
	}
}

// The return side is lifted; the PARAM side deliberately is not. Unboxing an
// arg needs a proof that the runtime value really is that primitive — the
// direct-call path gets that at the call site (coerce-arg-to-type), a dynamic
// override boundary does not. A typed PARAM must therefore still stay off the
// override path rather than registering a wrapper that would mis-assert.
func TestTypedParamStillSkipsOverride(t *testing.T) {
	ensureLoader()

	v := runLispExpr(t, `(do (create-ns (quote tparam))
      (ir.passes.pipeline/lower-ns-to-go "tparam" (quote tparam)
        [(quote (defn twice [^long n] (* 2 n)))]))`)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go source string, got %T", v)
	}
	src := string(s)

	// Precondition: a typed PARAM really was emitted.
	if !regexp.MustCompile(`func Twice\(ec \*vm\.ExecContext, arg0 int\)`).MatchString(src) {
		t.Fatalf("expected a typed param `func Twice(ec *vm.ExecContext, arg0 int)`:\n%s", src)
	}
	if regexp.MustCompile(`"twice":\s*__gogen_wrap`).MatchString(src) {
		t.Fatalf("typed-PARAM fn registered an override; the wrapper cannot soundly unbox args[i]:\n%s", src)
	}
}

// Multi-arity sibling of TestTypedReturnRegistersBoxedOverride. Both arities of
// `choose` return bool, so before #554's multi-arity arm the whole fn was
// disqualified from override registration (multi-arity-override-groups required
// every arity to be uniformly vm.Value) and emitted no init() at all — the two
// native helpers were unreachable from a dynamic call.
func TestTypedReturnRegistersMultiArityOverride(t *testing.T) {
	ensureLoader()

	v := runLispExpr(t, `(do (create-ns (quote tmulti))
      (ir.passes.pipeline/lower-ns-to-go "tmulti" (quote tmulti)
        [(quote (defn choose ([x] (= x 1)) ([x y] (= x y))))]))`)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go source string, got %T", v)
	}
	src := string(s)

	// Precondition: both arities lowered with typed bool returns.
	for _, want := range []string{"func Choose_1(", "func Choose_2("} {
		if !strings.Contains(src, want) {
			t.Fatalf("expected %s in the lowered output:\n%s", want, src)
		}
	}
	if !regexp.MustCompile(`func Choose_1\([^)]*\) bool`).MatchString(src) {
		t.Fatalf("expected Choose_1 to return a typed bool:\n%s", src)
	}

	// The fix: one combined adapter registered under the bare name, with each
	// typed arm boxed independently.
	if !regexp.MustCompile(`"choose":\s*__gogen_wrap`).MatchString(src) {
		t.Fatalf("multi-arity typed-return fn is NOT registered (#554):\n%s", src)
	}
	for _, want := range []string{"vm.Boolean(Choose_1(ec, args[0]))", "vm.Boolean(Choose_2(ec, args[0], args[1]))"} {
		if !strings.Contains(src, want) {
			t.Fatalf("expected each arm boxed — missing %q:\n%s", want, src)
		}
	}
}
