/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/ir"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// compileViaDefault returns the *vm.CodeChunk emitted by the standard
// bytecode-only compile path (the baseline).
func compileViaDefault(b *testing.B, src string, name string) *vm.CodeChunk {
	b.Helper()
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	c := compiler.NewCompiler(consts, ns)
	c.SetSource("bench-default")
	if _, _, err := c.CompileMultiple(strings.NewReader(src)); err != nil {
		b.Fatalf("default compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol(name)).(*vm.Var).Deref()
	return v.(*vm.Func).Chunk()
}

// compileViaLispPipeline routes the form through:
//
//	ir.build/build-fn  →  ir.passes.pipeline/optimize-fn  →  ir/lower
//
// All in Lisp. Returns the *vm.CodeChunk.
func compileViaLispPipeline(b *testing.B, src string) *vm.CodeChunk {
	b.Helper()
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("bench-lisp-pipeline")
	expr := fmt.Sprintf(`(ir.passes.pipeline/compile-form (quote %s))`, src)
	_, result, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		b.Fatalf("lisp pipeline compile: %v", err)
	}
	boxed, ok := result.(*vm.Boxed)
	if !ok {
		b.Fatalf("expected *vm.Boxed, got %T (%v)", result, result)
	}
	chunk, ok := boxed.Unbox().(*vm.CodeChunk)
	if !ok {
		b.Fatalf("expected boxed *vm.CodeChunk, got %T", boxed.Unbox())
	}
	return chunk
}

// pipelineCorpus pairs a Lisp source with the metadata needed to
// execute the resulting chunk.
var pipelineCorpus = []struct {
	name string
	src  string
	args []vm.Value
}{
	// --- toy cases (kept from initial corpus) ---
	{
		name: "const-arith",
		src:  `(defn const-arith [] (+ 1 (* 2 3)))`,
		args: nil,
	},
	{
		name: "id",
		src:  `(defn id [x] x)`,
		args: []vm.Value{vm.Int(7)},
	},
	// --- shapes derived from real xsofy code ---
	// Tests cond + comparison + literal returns. Direct from xsofy/util.lg.
	{
		name: "sign",
		src:  `(defn sign [n] (cond (> n 0) 1 (< n 0) -1 :else 0))`,
		args: []vm.Value{vm.Int(-7)},
	},
	// Nested calls to core fns. xsofy/combat.lg's hit-chance shape.
	{
		name: "clamp",
		src:  `(defn clamp [lo hi v] (max lo (min hi v)))`,
		args: []vm.Value{vm.Int(0), vm.Int(100), vm.Int(150)},
	},
	// Duplicate ref → CSE candidate.
	{
		name: "square",
		src:  `(defn square [n] (* n n))`,
		args: []vm.Value{vm.Int(9)},
	},
	// Repeated subexpression + Const → CSE then constfold.
	{
		name: "poly2",
		src:  `(defn poly2 [x] (+ (* x x) (* x x) 5))`,
		args: []vm.Value{vm.Int(3)},
	},
	// Loop with invariant → LICM candidate.
	{
		name: "loop-sum",
		src: `(defn loop-sum [n]
		         (loop [i 0 acc 0]
		           (if (>= i n) acc (recur (inc i) (+ acc i)))))`,
		args: []vm.Value{vm.Int(10)},
	},
	// Let chain → tests basic let lowering through the pipeline.
	{
		name: "nested-let",
		src: `(defn nested-let [x y]
		         (let [a (+ x 1) b (+ y 2) c (+ a b)] (* c c)))`,
		args: []vm.Value{vm.Int(3), vm.Int(4)},
	},
	// Unused binding → DCE candidate.
	{
		name: "dead-load",
		src: `(defn dead-load [n]
		         (let [unused (+ 1 2) x (* n 5)] x))`,
		args: []vm.Value{vm.Int(4)},
	},
	// Algebraic identities — (* x 1), (- x 0), (+ x 0) — constfold try-identity!
	{
		name: "algebraic-identities",
		src: `(defn algebraic-identities [x]
		         (+ (* x 1) (- x 0) (+ x 0)))`,
		args: []vm.Value{vm.Int(7)},
	},
	// --- xsofy-derived: real game logic ---
	// Sign function from xsofy/util.lg — exercises cond + comparisons
	{
		name: "xsofy-sign",
		src:  `(defn xsofy-sign [n] (cond (> n 0) 1 (< n 0) -1 :else 0))`,
		args: []vm.Value{vm.Int(-7)},
	},
	// Strength modifier from xsofy/combat.lg — branched arithmetic
	{
		name: "strength-modifier",
		src: `(defn strength-modifier [strength str-req]
		         (let [diff (- strength str-req)]
		           (if (>= diff 0) (* diff 1) (* diff 2))))`,
		args: []vm.Value{vm.Int(10), vm.Int(5)},
	},
	// Manhattan distance with destructure-free signature — exercises let chains
	{
		name: "manhattan-flat",
		src: `(defn manhattan-flat [x1 y1 x2 y2]
		         (let [dx (- x2 x1) dy (- y2 y1)]
		           (+ (if (< dx 0) (- 0 dx) dx)
		              (if (< dy 0) (- 0 dy) dy))))`,
		args: []vm.Value{vm.Int(1), vm.Int(2), vm.Int(4), vm.Int(6)},
	},
	// Sum-to-n with explicit loop+recur — exercises loop-carried block-args
	{
		name: "sum-to-n",
		src: `(defn sum-to-n [n]
		         (loop [i 0 acc 0]
		           (if (>= i n) acc (recur (+ i 1) (+ acc i)))))`,
		args: []vm.Value{vm.Int(100)},
	},
	// let* is the special-form shape that the let macro expands to.
	// build-list previously only matched 'let.
	{
		name: "let-star",
		src:  `(defn let-star [n] (let* [x (+ n 1)] (* x 2)))`,
		args: []vm.Value{vm.Int(3)},
	},
	// Vector literal as expression — defn returns a constructed vector.
	{
		name: "vec-of",
		src:  `(defn vec-of [a b] [a b 99])`,
		args: []vm.Value{vm.Int(1), vm.Int(2)},
	},
	// Map literal as expression — defn returns a constructed map.
	{
		name: "map-of",
		src:  `(defn map-of [k v] {:k k :v v})`,
		args: []vm.Value{vm.Int(1), vm.Int(2)},
	},
	// Vector destructuring in let - xsofy/util.lg distance/manhattan shape.
	{
		name: "let-vector-destructure",
		src: `(defn let-vector-destructure [p1 p2]
		         (let [[x1 y1] p1
		               [x2 y2] p2]
		           (+ x1 y1 x2 y2)))`,
		args: []vm.Value{
			vm.ArrayVector{vm.Int(1), vm.Int(2)},
			vm.ArrayVector{vm.Int(3), vm.Int(4)},
		},
	},
	// Map :keys destructuring in let - xsofy/ai.lg world shape.
	{
		name: "let-map-keys-destructure",
		src: `(defn let-map-keys-destructure [world]
		         (let [{:keys [terrain width height]} world]
		           (+ terrain width height)))`,
		args: []vm.Value{vm.EmptyPersistentMap.
			Assoc(vm.Keyword("terrain"), vm.Int(1)).(*vm.PersistentMap).
			Assoc(vm.Keyword("width"), vm.Int(2)).(*vm.PersistentMap).
			Assoc(vm.Keyword("height"), vm.Int(3)).(*vm.PersistentMap)},
	},
	// Direct namespace-qualified symbol without an alias in the caller ns.
	{
		name: "direct-ns-math-sqrt",
		src:  `(defn direct-ns-math-sqrt [x] (math/sqrt x))`,
		args: []vm.Value{vm.Float(9)},
	},
}

// BenchmarkCompileAndRun compares the two paths on the corpus.
// Each iteration: compile (via one path), then execute, measure both
// stages combined.
func BenchmarkCompileAndRun(b *testing.B) {
	ensureLoader()

	for _, entry := range pipelineCorpus {
		b.Run(entry.name+"/default", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				chunk := compileViaDefault(b, entry.src, entry.name)
				frame := vm.NewFrame(chunk, entry.args)
				_, err := frame.Run()
				vm.ReleaseFrame(frame)
				if err != nil {
					b.Fatalf("run: %v", err)
				}
			}
		})
		b.Run(entry.name+"/lisp-pipeline", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				chunk := compileViaLispPipeline(b, entry.src)
				frame := vm.NewFrame(chunk, entry.args)
				_, err := frame.Run()
				vm.ReleaseFrame(frame)
				if err != nil {
					b.Fatalf("run: %v", err)
				}
			}
		})
	}
}

// BenchmarkExecOnly compares ONLY execution of the produced chunks,
// not the compile time. This isolates "did the Lisp pipeline produce
// faster-running bytecode?" from "is the pipeline itself slow?"
func BenchmarkExecOnly(b *testing.B) {
	ensureLoader()

	for _, entry := range pipelineCorpus {
		defChunk := compileViaDefault(&testing.B{}, entry.src, entry.name)
		lispChunk := compileViaLispPipeline(&testing.B{}, entry.src)

		b.Run(entry.name+"/default", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				frame := vm.NewFrame(defChunk, entry.args)
				_, err := frame.Run()
				vm.ReleaseFrame(frame)
				if err != nil {
					b.Fatalf("run: %v", err)
				}
			}
		})
		b.Run(entry.name+"/lisp-pipeline", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				frame := vm.NewFrame(lispChunk, entry.args)
				_, err := frame.Run()
				vm.ReleaseFrame(frame)
				if err != nil {
					b.Fatalf("run: %v", err)
				}
			}
		})
	}
}

// Sanity check: both paths produce chunks that execute and return
// equivalent values (modulo Lisp-pipeline bugs we haven't surfaced).
// pendingPhase5 are corpus entries whose lisp-pipeline path produces
// bytecode with bad stack tracking — they panic at runtime in
// vm.(*Frame).nth or vm.(*Frame).push. Phase 5 of the plan
// (docs/superpowers/plans/2026-05-22-ir-pipeline-faithful-port.md)
// fixes Lower's BranchTarget args materialization and nested-if join-
// block handling; remove from this set as each becomes green.
// pendingPhase5 is empty — all corpus entries are exercised so failure
// modes (panics, value mismatches, validate errors) surface in test
// output. Entries still failing post-Phase 5 should be triaged via the
// test output rather than skipped here.
var pendingPhase5 = map[string]string{}

func TestPipelineRoundTrip_ProducesExecutableChunks(t *testing.T) {
	ensureLoader()
	for _, entry := range pipelineCorpus {
		t.Run(entry.name, func(t *testing.T) {
			if reason, pending := pendingPhase5[entry.name]; pending {
				t.Skipf("pending Phase 5 (%s)", reason)
			}
			// Wrap with recover so unexpected panics surface as test
			// failures rather than killing the whole binary.
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panicked: %v", r)
				}
			}()
			defChunk := func() *vm.CodeChunk {
				consts := vm.NewConsts()
				ns := rt.NS(rt.NameCoreNS)
				c := compiler.NewCompiler(consts, ns)
				_, _, err := c.CompileMultiple(strings.NewReader(entry.src))
				if err != nil {
					t.Fatalf("default compile: %v", err)
				}
				v := ns.Lookup(vm.Symbol(entry.name)).(*vm.Var).Deref()
				return v.(*vm.Func).Chunk()
			}()
			lispChunk := func() *vm.CodeChunk {
				consts := vm.NewConsts()
				c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
				expr := fmt.Sprintf(`(ir.passes.pipeline/compile-form (quote %s))`, entry.src)
				_, result, err := c.CompileMultiple(strings.NewReader(expr))
				if err != nil {
					t.Fatalf("lisp compile: %v", err)
				}
				boxed, ok := result.(*vm.Boxed)
				if !ok {
					t.Fatalf("expected *vm.Boxed, got %T", result)
				}
				return boxed.Unbox().(*vm.CodeChunk)
			}()

			run := func(c *vm.CodeChunk) vm.Value {
				f := vm.NewFrame(c, entry.args)
				out, err := f.Run()
				vm.ReleaseFrame(f)
				if err != nil {
					t.Fatalf("run: %v", err)
				}
				return out
			}
			defResult := run(defChunk)
			lispResult := run(lispChunk)
			if !vm.ValueEquals(defResult, lispResult) {
				t.Errorf("results differ: default=%s lisp=%s", defResult, lispResult)
			}
		})
	}
	_ = ir.Op(0) // keep import alive even if unused above
}

func TestIRCompileDefnUsesCallerNamespaceAliases(t *testing.T) {
	ensureLoader()

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	_, _, err := c.CompileMultiple(strings.NewReader(`
		(require 'ir.passes.pipeline)
		(ns ir-band2.alias-test
		  (:require [math :as m]))
		(set! *ir-compile* true)
		(defn alias-sqrt [x] (m/sqrt x))
	`))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	v := rt.NS("ir-band2.alias-test").Lookup(vm.Symbol("alias-sqrt")).(*vm.Var).Deref()
	f := vm.NewFrame(v.(*vm.Func).Chunk(), []vm.Value{vm.Float(9)})
	out, err := f.Run()
	vm.ReleaseFrame(f)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if out.String() != "3" && out.String() != "3.0" {
		t.Fatalf("alias-sqrt result: got %s, want 3", out)
	}
}
