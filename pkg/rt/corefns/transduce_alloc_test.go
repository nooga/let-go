/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package corefns_test

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// evalV evaluates a let-go source string and returns the result. Uses the same
// compiler infrastructure as the language test suite.
func evalV(t *testing.T, src string) vm.Value {
	t.Helper()
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	if ns == nil {
		t.Fatal("core namespace not found")
	}
	ctx := compiler.NewCompiler(consts, ns)
	ctx.SetSource("<transduce-alloc-test>")
	_, result, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("evalV(%q) compile error: %v", src, err)
	}
	return result
}

// TestTransduceNoIntermediateSeqAlloc asserts that the fused transducer form
// allocates strictly less than the naive lazy-seq form. The whole point of
// transducers is to eliminate intermediate lazy-seq allocations (LazySeq,
// ArrayVectorSeq nodes). This test proves the optimization works.
func TestTransduceNoIntermediateSeqAlloc(t *testing.T) {
	// Set up the namespace loader so require works if needed
	consts := vm.NewConsts()
	loaderCtx := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	rt.SetNSLoader(resolver.NewNSResolver(loaderCtx, []string{"."}))

	// Two forms that should produce the same result:
	// naive: (vec (map f [1..8]))    — allocates intermediate lazy-seq
	// fused: (into [] (map f) [1..8]) — deforested; no intermediate seq
	naive := `(vec (map (fn [x] (* 2 x)) [1 2 3 4 5 6 7 8]))`
	fused := `(into [] (map (fn [x] (* 2 x))) [1 2 3 4 5 6 7 8])`

	// Warm up: run each once to trigger JIT, class loading, etc.
	_ = evalV(t, naive)
	_ = evalV(t, fused)

	// Measure allocations per run over 100 iterations (high N for stable estimate)
	naiveAllocs := testing.AllocsPerRun(100, func() { _ = evalV(t, naive) })
	fusedAllocs := testing.AllocsPerRun(100, func() { _ = evalV(t, fused) })

	if fusedAllocs >= naiveAllocs {
		t.Fatalf("expected fused transduce to allocate less than naive lazy-seq: fused=%.0f naive=%.0f",
			fusedAllocs, naiveAllocs)
	}
	t.Logf("fused=%.0f naive=%.0f allocs/run (%.0f%% reduction)",
		fusedAllocs, naiveAllocs, (1-fusedAllocs/naiveAllocs)*100)
}
