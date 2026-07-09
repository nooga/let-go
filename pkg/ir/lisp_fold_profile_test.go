/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// ITER-0035 (static arm): fold-over-rest unroll helps-vs-hinders profiler.
// Measures the STATIC cost of unrolling on the gogen-Go target across an N-sweep:
// caller code size and element-dispatch count, unroll ON vs OFF. No binary
// builds — this is the code-growth curve that tunes *max-unroll* and frames the
// runtime arm. Run: go test ./pkg/ir -run TestProfileUnrollStatic -v
func TestProfileUnrollStatic(t *testing.T) {
	if os.Getenv("LG_PROFILE") == "" {
		t.Skip("set LG_PROFILE=1 to run the fold-unroll static profiler (a measurement tool, not a gate)")
	}
	ensureLoader()
	ns := []int{1, 2, 3, 4, 6, 8, 12, 16, 24, 32, 48, 64}
	t.Logf("fold-over-rest unroll — STATIC gogen-Go cost sweep (caller = Parse)")
	t.Logf("%4s | %-26s | %-26s | %8s", "N", "OFF (loop call)", "ON (unrolled)", "ON/OFF")
	t.Logf("%4s | %8s %8s %6s | %8s %8s %6s | %8s", "", "bytes", "dispatch", "ifs", "bytes", "dispatch", "ifs", "size×")
	var prevOnBytes int
	for i, n := range ns {
		off := sliceFunc(lowerFoldNsMode(t, n, false), "Parse")
		on := sliceFunc(lowerFoldNsMode(t, n, true), "Parse")
		offB, onB := len(off), len(on)
		offD := strings.Count(off, "InvokeValueEC")
		onD := strings.Count(on, "InvokeValueEC")
		offI := strings.Count(off, "IsTruthy")
		onI := strings.Count(on, "IsTruthy")
		ratio := float64(onB) / float64(offB)
		perElem := ""
		if i > 0 {
			perElem = fmt.Sprintf("  +%d B/elem", (onB-prevOnBytes)/(n-ns[i-1]))
		}
		prevOnBytes = onB
		t.Logf("%4d | %8d %8d %6d | %8d %8d %6d | %7.1f%s", n, offB, offD, offI, onB, onD, onI, ratio, perElem)
	}
	t.Logf("OFF caller = a single dispatch to the combinator (O(1) code, O(N) runtime loop).")
	t.Logf("ON  caller = N-way short-circuit chain (O(N) code, N element dispatches, no loop overhead).")
}

// ── Bytecode arm: does unroll cut executed VM ops (the loop conditionals)? ──

// runInNs compiles+runs prog in namespace nsName (definitions persist across calls).
func runInNs(t *testing.T, nsName, prog string) vm.Value {
	t.Helper()
	ns := rt.NS(nsName)
	c := compiler.NewCompiler(vm.NewConsts(), ns)
	c.SetSource("bcprof")
	_, res, err := c.CompileMultiple(strings.NewReader(prog))
	if err != nil {
		t.Fatalf("run %q: %v", prog, err)
	}
	return res
}

// unrolledAnySrc mirrors specialize-fold's output: the unrolled any(N) chain.
func unrolledAnySrc(name string, n int) string {
	params := "[input"
	for i := 0; i < n; i++ {
		params += fmt.Sprintf(" e%d", i)
	}
	params += "]"
	body := "false"
	for i := n - 1; i >= 0; i-- {
		body = fmt.Sprintf("(if (e%d input) true %s)", i, body)
	}
	return fmt.Sprintf("(defn %s %s %s)", name, params, body)
}

func profiledOps(t *testing.T, nsName, prog string) uint64 {
	t.Helper()
	vm.ResetProfile()
	vm.ProfilingEnabled.Store(true)
	runInNs(t, nsName, prog)
	vm.ProfilingEnabled.Store(false)
	var total uint64
	for _, s := range vm.ProfileSnapshot() {
		total += s.Count
	}
	return total
}

func TestProfileUnrollBytecode(t *testing.T) {
	if os.Getenv("LG_PROFILE") == "" {
		t.Skip("set LG_PROFILE=1 to run the bytecode-VM unroll profiler")
	}
	ensureLoader()
	const B = 200000
	t.Logf("fold-over-rest unroll — BYTECODE-VM cost (all-false rules → full N-scan)")
	t.Logf("%4s | %10s %10s %7s | %12s %12s %7s", "N", "ops/loop", "ops/unroll", "ops×", "ns/op loop", "ns/op unrl", "time×")
	for _, n := range []int{2, 4, 8, 16, 32} {
		nsN := fmt.Sprintf("bcprof%d", n)
		// definitions: N all-false rules, the loop combinator, the unrolled fn, input.
		var defs strings.Builder
		args := ""
		for i := 0; i < n; i++ {
			fmt.Fprintf(&defs, "(defn r%d [x] false)\n", i)
			args += fmt.Sprintf(" r%d", i)
		}
		defs.WriteString(`(defn loopc [input & fs] (loop [v fs] (if (empty? v) false (if ((first v) input) true (recur (rest v))))))` + "\n")
		defs.WriteString(unrolledAnySrc("unrollc", n) + "\n")
		runInNs(t, nsN, defs.String())

		driverLoop := fmt.Sprintf("(loopc 42%s)", args)
		driverUnroll := fmt.Sprintf("(unrollc 42%s)", args)

		// op counts (single invocation each).
		opsLoop := profiledOps(t, nsN, driverLoop)
		opsUnroll := profiledOps(t, nsN, driverUnroll)

		// wall-time: B invocations via dotimes (loop overhead cancels in the ratio).
		timeIt := func(driver string) time.Duration {
			prog := fmt.Sprintf("(dotimes [_ %d] %s)", B, driver)
			start := time.Now()
			runInNs(t, nsN, prog)
			return time.Since(start)
		}
		tLoop := timeIt(driverLoop)
		tUnroll := timeIt(driverUnroll)
		nsLoop := float64(tLoop.Nanoseconds()) / B
		nsUnroll := float64(tUnroll.Nanoseconds()) / B
		t.Logf("%4d | %10d %10d %6.2fx | %12.0f %12.0f %6.2fx", n,
			opsLoop, opsUnroll, float64(opsLoop)/float64(opsUnroll),
			nsLoop, nsUnroll, nsLoop/nsUnroll)
	}
	t.Logf("ops = executed VM opcodes for ONE call over N rules; ns/op = per-call wall-time.")
}
