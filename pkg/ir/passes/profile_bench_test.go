/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package passes

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/ir"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// compileFn compiles a `(defn name [args] body)` form and returns the
// resulting bytecode chunk along with the arity declared on the fn.
func compileFn(t testing.TB, src, name string) (*vm.CodeChunk, int) {
	t.Helper()
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	if _, _, err := ctx.CompileMultiple(strings.NewReader(src)); err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol(name))
	if v == nil {
		t.Fatalf("symbol %s not found", name)
	}
	fnVal := v.(*vm.Var).Deref()
	fn, ok := fnVal.(*vm.Func)
	if !ok {
		t.Fatalf("%s is not a Func: %T", name, fnVal)
	}
	return fn.Chunk(), fn.Arity()
}

// optimize runs Build → ConstFold → CSE → DCE → Lower and returns the
// optimized chunk.
func optimize(t testing.TB, chunk *vm.CodeChunk, name string, arity int) *vm.CodeChunk {
	t.Helper()
	irFn, err := ir.Build(chunk, name, arity, false)
	if err != nil {
		t.Fatalf("Build %s: %v", name, err)
	}
	ConstFold(irFn)
	CSE(irFn)
	DCE(irFn)
	out, err := ir.Lower(irFn)
	if err != nil {
		t.Fatalf("Lower %s: %v", name, err)
	}
	return out
}

// runWithProfile runs chunk N times with `args`, reporting the bytecode
// size, the result of the last invocation, and the per-opcode counts
// observed (which will be N times what a single run records).
func runWithProfile(t testing.TB, chunk *vm.CodeChunk, args []vm.Value, n int) (size int, result vm.Value, counts []vm.ProfileSample) {
	t.Helper()
	vm.ResetProfile()
	vm.ProfilingEnabled.Store(true)
	defer vm.ProfilingEnabled.Store(false)

	for i := 0; i < n; i++ {
		frame := vm.NewFrame(chunk, args)
		var err error
		result, err = frame.Run()
		vm.ReleaseFrame(frame)
		if err != nil {
			t.Fatalf("run: %v", err)
		}
	}
	size = len(chunk.Code())
	counts = vm.ProfileSnapshot()
	return
}

func sumOpcodeCounts(counts []vm.ProfileSample) uint64 {
	var total uint64
	for _, s := range counts {
		total += s.Count
	}
	return total
}

func formatProfile(label string, size int, counts []vm.ProfileSample, n int) string {
	var sb strings.Builder
	total := sumOpcodeCounts(counts)
	fmt.Fprintf(&sb, "%s: %d-word bytecode, %d total opcode dispatches across %d runs (%d per run)\n",
		label, size, total, n, total/uint64(n))
	for _, s := range counts {
		fmt.Fprintf(&sb, "  %-18s %10d  (%d per run)\n", s.Name, s.Count, s.Count/uint64(n))
	}
	return sb.String()
}

// TestProfile_SumtoBaselineVsOptimized characterizes the impact of the
// IR optimization pipeline on a tight loop. (sumto 10) runs 10 inner
// loop iterations; we run it 100 times to get a stable opcode profile.
//
// This is a measurement, not a regression test — there are no pass/fail
// thresholds. It exists so before/after numbers are reproducible.
func TestProfile_SumtoBaselineVsOptimized(t *testing.T) {
	src := `(defn sumto-prof [n] (loop [i 0 acc 0] (if (< i n) (recur (inc i) (+ acc i)) acc)))`
	chunk, arity := compileFn(t, src, "sumto-prof")
	const runs = 100
	args := []vm.Value{vm.Int(10)}

	baseSize, baseResult, baseCounts := runWithProfile(t, chunk, args, runs)
	optChunk := optimize(t, chunk, "sumto-prof", arity)
	optSize, optResult, optCounts := runWithProfile(t, optChunk, args, runs)

	if baseResult.String() != optResult.String() {
		t.Errorf("baseline (%s) and optimized (%s) results differ", baseResult, optResult)
	}
	if baseResult.String() != "45" {
		t.Errorf("expected (sumto 10) = 45, got %s", baseResult)
	}

	t.Logf("\n%s\n%s", formatProfile("BASELINE", baseSize, baseCounts, runs),
		formatProfile("OPTIMIZED", optSize, optCounts, runs))

	baseTotal := sumOpcodeCounts(baseCounts)
	optTotal := sumOpcodeCounts(optCounts)
	delta := int64(optTotal) - int64(baseTotal)
	t.Logf("\nopcode dispatches: baseline=%d  optimized=%d  delta=%+d (%.1f%%)",
		baseTotal, optTotal, delta, 100*float64(delta)/float64(baseTotal))
	t.Logf("bytecode size:     baseline=%d  optimized=%d  delta=%+d",
		baseSize, optSize, optSize-baseSize)
}

// TestProfile_MandelbrotHelpersBaselineVsOptimized: same for the three
// arithmetic helpers from examples/mandelbrot.lg (cadd, cmul, cmagsq).
// Each helper is called with a fixed set of args; we run the body 1000
// times to amplify CSE opportunities.
func TestProfile_MandelbrotHelpersBaselineVsOptimized(t *testing.T) {
	src := `
(def scale 1000)
(defn cadd [a b c d] (vector (+ a c) (+ b d)))
(defn cmul [a b c d]
  (vector (/ (- (* a c) (* b d)) scale)
          (/ (+ (* a d) (* b c)) scale)))
(defn cmagsq [a b]
  (/ (+ (* a a) (* b b)) scale))
`
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	if _, _, err := ctx.CompileMultiple(strings.NewReader(src)); err != nil {
		t.Fatalf("compile: %v", err)
	}

	cases := []struct {
		name string
		args []vm.Value
	}{
		{"cadd", []vm.Value{vm.Int(100), vm.Int(200), vm.Int(50), vm.Int(75)}},
		{"cmul", []vm.Value{vm.Int(1500), vm.Int(2500), vm.Int(800), vm.Int(1200)}},
		{"cmagsq", []vm.Value{vm.Int(1500), vm.Int(2500)}},
	}

	const runs = 1000
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := ns.Lookup(vm.Symbol(c.name))
			fnVal := v.(*vm.Var).Deref()
			fn := fnVal.(*vm.Func)
			chunk := fn.Chunk()
			arity := fn.Arity()

			baseSize, baseResult, baseCounts := runWithProfile(t, chunk, c.args, runs)
			optChunk := optimize(t, chunk, c.name, arity)
			optSize, optResult, optCounts := runWithProfile(t, optChunk, c.args, runs)

			if baseResult.String() != optResult.String() {
				t.Errorf("%s: baseline (%s) and optimized (%s) results differ", c.name, baseResult, optResult)
			}

			t.Logf("\n%s\n%s", formatProfile(c.name+" BASELINE", baseSize, baseCounts, runs),
				formatProfile(c.name+" OPTIMIZED", optSize, optCounts, runs))

			baseTotal := sumOpcodeCounts(baseCounts)
			optTotal := sumOpcodeCounts(optCounts)
			delta := int64(optTotal) - int64(baseTotal)
			pct := 100 * float64(delta) / float64(baseTotal)
			t.Logf("%s opcode dispatches: baseline=%d  optimized=%d  delta=%+d (%.1f%%)  bytecode-size: baseline=%d optimized=%d",
				c.name, baseTotal, optTotal, delta, pct, baseSize, optSize)
		})
	}
}
