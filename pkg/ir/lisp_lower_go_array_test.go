/*
 * Copyright (c) 2026 Matt Parrett
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// #1055: aget/aset/alength on a statically typed array lower to the typed
// accessors on *vm.TypedArray, so the index and the element stay native and
// the arithmetic around them does not fall back onto rt.*Value helpers.

func lowerArrayFn(t *testing.T, src string) string {
	t.Helper()
	fn := buildLispIR(t, src)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")
	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v (reason=%v)", got, result.ValueAt(vm.Keyword("reason")))
	}
	return bindAndRenderGoDecl(t, result)
}

func mustContain(t *testing.T, rendered string, wants ...string) {
	t.Helper()
	for _, w := range wants {
		if !strings.Contains(rendered, w) {
			t.Fatalf("expected %q in the lowered Go, got:\n%s", w, rendered)
		}
	}
}

func mustNotContain(t *testing.T, rendered string, donts ...string) {
	t.Helper()
	for _, d := range donts {
		if strings.Contains(rendered, d) {
			t.Fatalf("did not expect %q in the lowered Go, got:\n%s", d, rendered)
		}
	}
}

// A ^doubles param is a *vm.TypedArray, aget on it is AtFloat64 with a native
// int64 index, and the accumulator stays float64: no boxed element, no
// AddValue. The trampoline (which does box the index) stays as the else branch
// of the root guard.
func TestLowerGoDoublesHintAgetLowersToTypedAccessor(t *testing.T) {
	ensureLoader()
	rendered := lowerArrayFn(t, `(defn sum-arr [^doubles a n]
	  (loop [i 0 s 0.0]
	    (if (>= i n) s (recur (inc i) (+ s (aget a i))))))`)
	mustContain(t, rendered, "arg0 *vm.TypedArray", ".AtFloat64(i)", "rt.NativePrimsIntact()", `"aget"`, "s float64")
	mustNotContain(t, rendered, "AddValue", "s vm.Value")
}

// aset with an int-typed value into a double-array widens the value to float64,
// as TypedArray.Set coerces an Int, and lowers to SetFloat64.
func TestLowerGoDoublesHintAsetWidensIntValue(t *testing.T) {
	ensureLoader()
	// if/do rather than when: a when whose body is a side-effecting call plus
	// recur fails in legalize on main too, independently of arrays.
	rendered := lowerArrayFn(t, `(defn fill [^doubles a n]
	  (loop [i 0]
	    (if (< i n)
	      (do (aset a i (* 2 i)) (recur (inc i)))
	      a)))`)
	mustContain(t, rendered, ".SetFloat64(i, float64(")
	mustNotContain(t, rendered, "CoreAsetf")
}

// The array kind flows from the constructor: a double-array built in the same
// function types the local *vm.TypedArray, and aget/aset/alength on it take
// the typed path without any hint.
func TestLowerGoDoubleArrayConstructorTypesLocal(t *testing.T) {
	ensureLoader()
	rendered := lowerArrayFn(t, `(defn mk [n]
	  (let [a (double-array n)]
	    (aset a 0 1.5)
	    (+ (aget a 0) (alength a))))`)
	mustContain(t, rendered, "*vm.TypedArray", ".SetFloat64(", ".AtFloat64(", ".Len()")
	mustNotContain(t, rendered, "CoreAgetf", "CoreAsetf", "CoreAlengthf")
}

// ^longs is an int-array: AtInt64/SetInt64, element int64. The index must be
// a native int (here a ^long param); a vm.Value index keeps the boxed path.
func TestLowerGoLongsHintUsesInt64Accessors(t *testing.T) {
	ensureLoader()
	rendered := lowerArrayFn(t, `(defn bump [^longs a ^long i]
	  (aset a i (inc (aget a i))))`)
	mustContain(t, rendered, "arg0 *vm.TypedArray", ".AtInt64(", ".SetInt64(")
	mustNotContain(t, rendered, "CoreAgetf", "CoreAsetf")
}

// A caller that holds an unproven vm.Value and passes it to a lowered sibling
// with a *vm.TypedArray param gets the same guarded call a scalar param gets:
// a comma-ok assert, the direct call when it holds, the trampoline when the
// runtime value is not an array. The sibling is a hand-seeded registry entry,
// as in TestCrossNsLoweredDirectCallEmit.
func TestLowerGoArrayParamCallIsGuarded(t *testing.T) {
	ensureLoader()
	runLispExpr(t, `(do (create-ns (quote arrns)) (intern (quote arrns) (quote first-elt)))`)

	fn := buildLispIR(t, `(defn caller [x] (arrns/first-elt x))`)
	optimizeLispIR(t, fn)
	passVarCounter++
	varName := fmt.Sprintf("*arr-guard-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, fn)

	v := runLispExpr(t, fmt.Sprintf(`
	  (let [reg {(ir.lower-go/registry-key (quote arrns) "first-elt" 1)
	             {:go-name "FirstElt" :arity 1 :needs-error? true
	              :param-specs ["*vm.TypedArray"] :param-types [[:array :float]]
	              :result-spec "float64"
	              :native? false
	              :go-pkg "github.com/nooga/let-go/pkg/rt/core_go_lowered/arrns"}}]
	    (binding [ir.lower-go/*lowered-registry* reg
	              ir.lower-go/*native-imports-used* (atom #{})]
	      (ir.lower-go/lower %s :strict)))`, varName))
	m, ok := v.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("expected lower to return a map, got %T", v)
	}
	rendered := bindAndRenderGoDecl(t, m)
	mustContain(t, rendered, ".(*vm.TypedArray)", ".Kind() == vm.ArrayFloat", "arrns.FirstElt(ec, tg",
		`rt.LookupVar("arrns", "first-elt")`)
}

// A proven array of ANOTHER kind (a long-array into a ^doubles sibling) never
// takes the direct call: bytecode ignores the hint and aget works on any
// array, so the lowered caller bails to the trampoline rather than hand an
// int-array to AtFloat64. A proven non-array (a string) bails the same way;
// before this, the dtype pointer branch asserted on the concrete literal and
// emitted Go that does not compile.
func TestLowerGoArrayKindMismatchBailsToTrampoline(t *testing.T) {
	ensureLoader()
	runLispExpr(t, `(do (create-ns (quote arrns2)) (intern (quote arrns2) (quote first-elt)))`)

	for _, src := range []string{
		`(defn caller [] (arrns2/first-elt (long-array 3)))`,
		`(defn caller [] (arrns2/first-elt "nope"))`,
	} {
		fn := buildLispIR(t, src)
		optimizeLispIR(t, fn)
		passVarCounter++
		varName := fmt.Sprintf("*arr-mismatch-fn-%d*", passVarCounter)
		rt.NS(rt.NameCoreNS).Def(varName, fn)
		v := runLispExpr(t, fmt.Sprintf(`
		  (let [reg {(ir.lower-go/registry-key (quote arrns2) "first-elt" 1)
		             {:go-name "FirstElt" :arity 1 :needs-error? true
		              :param-specs ["*vm.TypedArray"] :param-types [[:array :float]]
		              :result-spec "float64"
		              :native? false
		              :go-pkg "github.com/nooga/let-go/pkg/rt/core_go_lowered/arrns2"}}]
		    (binding [ir.lower-go/*lowered-registry* reg
		              ir.lower-go/*native-imports-used* (atom #{})]
		      (ir.lower-go/lower %s :strict)))`, varName))
		m, ok := v.(*vm.PersistentMap)
		if !ok {
			t.Fatalf("expected lower to return a map, got %T", v)
		}
		rendered := bindAndRenderGoDecl(t, m)
		// (The long-array result still asserts *vm.TypedArray on its way into
		// its own typed local; what must be absent is the direct call.)
		mustContain(t, rendered, `rt.LookupVar("arrns2", "first-elt")`)
		mustNotContain(t, rendered, "arrns2.FirstElt")
	}
}

// An unhinted array param stays vm.Value and aget on it keeps today's path
// (the var dispatch here; the native primitive when the registry is loaded);
// the intrinsic only fires on a declared kind.
func TestLowerGoUnhintedArrayParamStaysBoxed(t *testing.T) {
	ensureLoader()
	rendered := lowerArrayFn(t, `(defn get0 [a] (aget a 0))`)
	mustContain(t, rendered, "arg0 vm.Value", `rt.LookupVar("clojure.core", "aget")`)
	mustNotContain(t, rendered, "AtFloat64", "AtValue")
}

// Function-level recur bypasses the ordinary direct-call kind guard. A
// changing or unproven kind must keep the function on bytecode, rather than
// feed the next iteration's typed accessor or emit an invalid assignment.
func TestLowerGoArrayRecurMismatchFallsBack(t *testing.T) {
	ensureLoader()
	for _, operand := range []string{`(long-array 3)`, `"x"`, `x`} {
		t.Run(operand, func(t *testing.T) {
			fn := buildLispIR(t, fmt.Sprintf(`(defn array-recur [^doubles a n x]
     (if (< n 1) (aget a 0) (recur %s (dec n) x)))`, operand))
			optimizeLispIR(t, fn)
			result := lowerGo(t, fn, ":bridge")
			if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("fallback") {
				t.Fatalf("expected bytecode fallback, got %v", result)
			}
			if decl := result.ValueAt(vm.Keyword("decl")); decl != vm.NIL {
				t.Fatalf("fallback emitted a declaration: %v", decl)
			}
		})
	}
}

func TestLowerGoArrayRecurSameKindStaysNative(t *testing.T) {
	ensureLoader()
	for _, operand := range []string{`a`, `(double-array 3)`} {
		t.Run(operand, func(t *testing.T) {
			rendered := lowerArrayFn(t, fmt.Sprintf(`(defn array-recur [^doubles a n]
     (if (< n 1) (aget a 0) (recur %s (dec n))))`, operand))
			mustContain(t, rendered, "a0 *vm.TypedArray", "a0.AtFloat64(0)")
		})
	}
}

// A rejected sibling must disappear from the direct-call registry too, so
// callers execute its original bytecode body and ignore the array hint.
func TestLowerGoArrayRecurFallbackCallerUsesTrampoline(t *testing.T) {
	ensureLoader()
	rendered := runLispExpr(t, `(do
   (create-ns (quote arrayrecur))
   (intern (quote arrayrecur) (quote changing))
   (intern (quote arrayrecur) (quote caller))
   (ir.passes.pipeline/lower-ns-to-go "arrayrecur" (quote arrayrecur)
    [(quote (defn changing [^doubles a n]
       (if (< n 1) (aget a 0) (recur (long-array 3) (dec n)))))
     (quote (defn caller [] (changing (double-array 3) 1)))]))`)
	source, ok := rendered.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go, got %T", rendered)
	}
	mustContain(t, string(source), `rt.CachedVarFn(&__v_arrayrecur_changing, "arrayrecur", "changing")`)
	mustNotContain(t, string(source), "func Changing(", "Changing(ec")
	got := runLispExpr(t, `(do
   (defn review-array-recur [^doubles a n]
     (if (< n 1) (aget a 0) (recur (long-array 3) (dec n))))
   (review-array-recur (double-array 3) 1))`)
	if got != vm.Int(0) {
		t.Fatalf("bytecode fallback = %v, want 0", got)
	}
}
