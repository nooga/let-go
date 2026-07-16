/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"runtime"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// Re-evaluating source with regex literals must not grow the process-global
// constant pool: each Eval interns into a per-eval CHILD pool whose lifetime
// is its chunk's, so transient constants are collectible instead of rooting
// the global pool (one leaked entry per eval, ~650B retained, unbounded in a
// long-lived REPL/nREPL session).
func TestRegexLiteralReEvalDoesNotGrowConstPool(t *testing.T) {
	if _, err := Eval(`(re-find #"x" "xyz")`); err != nil {
		t.Fatalf("warmup eval: %v", err)
	}
	before := len(consts.AllValues())
	for i := 0; i < 1000; i++ {
		if _, err := Eval(`(re-find #"x" "xyz")`); err != nil {
			t.Fatalf("eval %d: %v", i, err)
		}
	}
	if growth := len(consts.AllValues()) - before; growth != 0 {
		t.Fatalf("global const pool grew by %d entries across 1000 evals — transient constants are rooting the global pool", growth)
	}
}

// A transient eval's regex constant must become garbage once nothing
// references its chunk: the child pool is retained by the chunk alone.
func TestTransientRegexConstantIsCollectible(t *testing.T) {
	collected := make(chan struct{})
	func() {
		v, err := Eval(`#"transient-collectible-pattern"`)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		re, ok := v.(*vm.Regex)
		if !ok {
			t.Fatalf("expected *vm.Regex, got %T", v)
		}
		runtime.SetFinalizer(re, func(*vm.Regex) { close(collected) })
	}()
	for i := 0; i < 20; i++ {
		// Recycled VM frames keep stale stack-slot residue until reused;
		// run a filler eval so pooled frames overwrite any slot still
		// pointing at the regex, then collect.
		if _, err := Eval(`(+ 1 2)`); err != nil {
			t.Fatalf("filler eval: %v", err)
		}
		runtime.GC()
		select {
		case <-collected:
			return
		default:
		}
	}
	t.Fatal("transient regex constant never became collectible — its const pool is still rooted")
}

// Distinct literal occurrences of the SAME pattern must remain distinct
// objects: Clojure regex equality is identity, (= #"x" #"x") is false.
func TestRegexLiteralIdentityEquality(t *testing.T) {
	v, err := Eval(`(= #"same pat" #"same pat")`)
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	if vm.IsTruthy(v) {
		t.Fatal("(= #\"same pat\" #\"same pat\") must be false — distinct literals were merged")
	}
	v, err = Eval(`(let [r #"same pat"] (= r r))`)
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	if !vm.IsTruthy(v) {
		t.Fatal("a regex must equal itself")
	}
}

// The remaining transient eval paths must not root the process-global pool
// either (review finding on #496): each compiles through a per-evaluation
// child pool via NewTransientCompiler. One probe per path; the assertion is
// the same as TestRegexLiteralReEvalDoesNotGrowConstPool's.
func assertNoGlobalPoolGrowth(t *testing.T, label string, evalOnce func(i int)) {
	t.Helper()
	evalOnce(-1) // warmup: first use may legitimately intern shared constants
	before := len(consts.AllValues())
	for i := 0; i < 100; i++ {
		evalOnce(i)
	}
	if growth := len(consts.AllValues()) - before; growth != 0 {
		t.Fatalf("%s: global const pool grew by %d entries across 100 evaluations — transient constants are rooting the global pool", label, growth)
	}
}

func TestLoadStringDoesNotGrowConstPool(t *testing.T) {
	assertNoGlobalPoolGrowth(t, "load-string", func(i int) {
		if _, err := Eval(`(load-string "(re-find #\"x\" \"xyz\")")`); err != nil {
			t.Fatalf("eval %d: %v", i, err)
		}
	})
}

func TestNativeEvalDoesNotGrowConstPool(t *testing.T) {
	assertNoGlobalPoolGrowth(t, "eval", func(i int) {
		if _, err := Eval(`(eval (read-string "(re-find #\"x\" \"xyz\")"))`); err != nil {
			t.Fatalf("eval %d: %v", i, err)
		}
	})
}

func TestEvalInNSDoesNotGrowConstPool(t *testing.T) {
	ns := rt.NS(rt.NameCoreNS)
	assertNoGlobalPoolGrowth(t, "EvalInNS (pods)", func(i int) {
		if _, err := evalInNSChild(`(re-find #"x" "xyz")`, ns); err != nil {
			t.Fatalf("eval %d: %v", i, err)
		}
	})
}

// Nested transient evals (eval inside eval) create SIBLING child pools of the
// global pool — each level's constants retained only by its own chunks. Pins
// that layering ambiguity can't quietly root the global pool either.
func TestNestedEvalDoesNotGrowConstPool(t *testing.T) {
	assertNoGlobalPoolGrowth(t, "nested eval", func(i int) {
		if _, err := Eval(`(eval (read-string "(eval (read-string \"(re-find #\\\"x\\\" \\\"xyz\\\")\"))"))`); err != nil {
			t.Fatalf("eval %d: %v", i, err)
		}
	})
}
