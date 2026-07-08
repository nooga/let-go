/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TestSpikeIndexRPN performs parity verification and optional measurement
// of the four spike kernels vs the stack VM.
//
// SCENARIO-0032: LG_PROFILE=1 go test ./pkg/ir -run TestSpikeIndexRPN -v
//
// Kernels:
//
//	(a) typed numeric: (defn ksum [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc)))
//	(b) untyped boxed: same as (a), executed via spikeRun
//	(c) container: (defn kvec [v] (loop [i 0 acc 0] (if (< i (count v)) (recur (+ i 1) (+ acc (nth v i))) acc)))
//	(d) go-native: (defn kmax [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (max acc i)) acc)))
func TestSpikeIndexRPN(t *testing.T) {
	ensureLoader()

	// Parity section (always runs)
	t.Run("parity", func(t *testing.T) {
		testParityKernels(t)
	})

	// Measurement section (gated on LG_PROFILE)
	t.Run("measure", func(t *testing.T) {
		if os.Getenv("LG_PROFILE") == "" {
			t.Skip("set LG_PROFILE=1 to run measurements")
		}
		testMeasureKernels(t)
	})
}

// testParityKernels verifies that all 4 kernels produce identical results
// across spikeRunTyped, spikeRun, and stack VM.
func testParityKernels(t *testing.T) {
	testCases := []struct {
		name      string
		source    string
		stackName string
		args      []vm.Value
		label     string
	}{
		{
			name:      "kernel_a_typed_n10",
			source:    `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn ksum [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))))))`,
			stackName: "ksum",
			args:      []vm.Value{vm.Int(10)},
			label:     "kernel (a) typed, n=10",
		},
		{
			name:      "kernel_a_typed_n1000",
			source:    `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn ksum [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))))))`,
			stackName: "ksum",
			args:      []vm.Value{vm.Int(1000)},
			label:     "kernel (a) typed, n=1000",
		},
		{
			name:      "kernel_b_boxed_n10",
			source:    `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn m [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))))))`,
			stackName: "m",
			args:      []vm.Value{vm.Int(10)},
			label:     "kernel (b) boxed, n=10",
		},
		{
			name:      "kernel_b_boxed_n1000",
			source:    `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn m [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))))))`,
			stackName: "m",
			args:      []vm.Value{vm.Int(1000)},
			label:     "kernel (b) boxed, n=1000",
		},
		{
			name:      "kernel_c_container_v8",
			source:    `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn kvec [v] (loop [i 0 acc 0] (if (< i (count v)) (recur (+ i 1) (+ acc (nth v i))) acc))))))`,
			stackName: "kvec",
			args:      []vm.Value{vm.NewPersistentVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3), vm.Int(4), vm.Int(5), vm.Int(6), vm.Int(7), vm.Int(8)})},
			label:     "kernel (c) container, vec size 8",
		},
		{
			name:      "kernel_d_native_n10",
			source:    `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn kmax [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (max acc i)) acc))))))`,
			stackName: "kmax",
			args:      []vm.Value{vm.Int(10)},
			label:     "kernel (d) go-native, n=10",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			// Compile spike IR
			consts := vm.NewConsts()
			c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
			c.SetSource("test-spike-parity")
			_, irVal, err := c.CompileMultiple(strings.NewReader(tc.source))
			if err != nil {
				t.Fatalf("compile spike IR: %v", err)
			}

			spikeFn, decodeErr := decodeOptimizedIR(irVal)
			if decodeErr != nil {
				t.Fatalf("decode spike IR: %v", decodeErr)
			}

			// Compile stack VM
			ns := rt.NS(rt.NameCoreNS)
			stackC := compiler.NewCompiler(vm.NewConsts(), ns)
			stackC.SetSource("test-stack-parity")

			// Reconstruct the stack-side source from the quoted form
			var stackSource string
			if tc.stackName == "ksum" || tc.stackName == "m" {
				stackSource = fmt.Sprintf(`(defn %s [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc)))`, tc.stackName)
			} else if tc.stackName == "kvec" {
				stackSource = `(defn kvec [v] (loop [i 0 acc 0] (if (< i (count v)) (recur (+ i 1) (+ acc (nth v i))) acc)))`
			} else if tc.stackName == "kmax" {
				stackSource = `(defn kmax [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (max acc i)) acc)))`
			}

			if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(stackSource)); stackErr != nil {
				t.Fatalf("compile stack VM: %v", stackErr)
			}

			stackVar := ns.Lookup(vm.Symbol(tc.stackName))
			if stackVar == nil {
				t.Fatalf("stack VM: %s not found", tc.stackName)
			}
			stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
			stackChunk := stackFn.Chunk()

			// Run spike typed
			statsTyped := &spikeStats{}
			spikeTypedResult, spikeTypedErr := spikeRunTyped(spikeFn, tc.args, statsTyped)
			if spikeTypedErr != nil {
				t.Fatalf("spikeRunTyped error: %v", spikeTypedErr)
			}

			// Run spike boxed
			spikeBoxedResult, spikeBoxedErr := spikeRun(spikeFn, tc.args)
			if spikeBoxedErr != nil {
				t.Fatalf("spikeRun error: %v", spikeBoxedErr)
			}

			// Run stack VM
			stackFrame := vm.NewFrame(stackChunk, tc.args)
			stackResult, stackErr := stackFrame.Run()
			if stackErr != nil {
				t.Fatalf("stack VM error: %v", stackErr)
			}

			// Verify all three paths match
			if !valuesEqual(spikeTypedResult, stackResult) {
				t.Errorf("parity FAIL: spikeTyped=%v != stack=%v", spikeTypedResult, stackResult)
			}
			if !valuesEqual(spikeBoxedResult, stackResult) {
				t.Errorf("parity FAIL: spikeBoxed=%v != stack=%v", spikeBoxedResult, stackResult)
			}
			if !valuesEqual(spikeTypedResult, spikeBoxedResult) {
				t.Errorf("parity FAIL: spikeTyped=%v != spikeBoxed=%v", spikeTypedResult, spikeBoxedResult)
			}

			t.Logf("parity OK: result=%v", stackResult)
		})
	}
}

// KernelMeasurement describes one kernel × size combination.
type KernelMeasurement struct {
	Name             string
	Size             int
	SizeLabel        string
	DecodeMs         float64
	StackOps         uint64
	SpikeOps         int
	NsOpStack        float64
	NsOpSpikeBoxed   float64
	NsOpSpikeTyped   float64
	AllocsStack      float64
	AllocsSpikeBoxed float64
	AllocsSpikeTyped float64
	Speedup          float64 // typed vs stack
}

// testMeasureKernels runs the measurement suite.
func testMeasureKernels(t *testing.T) {
	vm.ResetProfile()
	vm.ProfilingEnabled.Store(false)

	measurements := []KernelMeasurement{}

	// Kernel (a): typed numeric, n ∈ {10, 100, 1000, 10000}
	for _, n := range []int{10, 100, 1000, 10000} {
		m := measureKernel(t, &kernelSpec{
			name:      fmt.Sprintf("kernel_a_%d", n),
			label:     fmt.Sprintf("(a) typed n=%d", n),
			sizeLabel: fmt.Sprintf("%d", n),
			source:    `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn ksum [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))))))`,
			stackName: "ksum",
			args: func() []vm.Value {
				return []vm.Value{vm.Int(int64(n))}
			},
			profiledSource: fmt.Sprintf(`(defn ksum_p [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc)))
(ksum_p %d)`, n),
		})
		measurements = append(measurements, m)
	}

	// Kernel (b): untyped boxed, n ∈ {10, 100, 1000, 10000}
	for _, n := range []int{10, 100, 1000, 10000} {
		m := measureKernel(t, &kernelSpec{
			name:      fmt.Sprintf("kernel_b_%d", n),
			label:     fmt.Sprintf("(b) boxed n=%d", n),
			sizeLabel: fmt.Sprintf("%d", n),
			source:    `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn m [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc))))))`,
			stackName: "m",
			args: func() []vm.Value {
				return []vm.Value{vm.Int(int64(n))}
			},
			profiledSource: fmt.Sprintf(`(defn m_p [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc)))
(m_p %d)`, n),
		})
		measurements = append(measurements, m)
	}

	// Kernel (c): container, vector sizes ∈ {8, 64, 512}
	for _, sz := range []int{8, 64, 512} {
		args := make([]vm.Value, sz)
		for i := 0; i < sz; i++ {
			args[i] = vm.Int(int64(i))
		}
		m := measureKernel(t, &kernelSpec{
			name:      fmt.Sprintf("kernel_c_%d", sz),
			label:     fmt.Sprintf("(c) container sz=%d", sz),
			sizeLabel: fmt.Sprintf("%d", sz),
			source:    `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn kvec [v] (loop [i 0 acc 0] (if (< i (count v)) (recur (+ i 1) (+ acc (nth v i))) acc))))))`,
			stackName: "kvec",
			args: func() []vm.Value {
				vec := make([]vm.Value, sz)
				for i := 0; i < sz; i++ {
					vec[i] = vm.Int(int64(i))
				}
				return []vm.Value{vm.NewPersistentVector(vec)}
			},
			profiledSource: fmt.Sprintf(`(defn kvec_p [v] (loop [i 0 acc 0] (if (< i (count v)) (recur (+ i 1) (+ acc (nth v i))) acc)))
(kvec_p (vec %s))`, vecLiteral(sz)),
		})
		measurements = append(measurements, m)
	}

	// Kernel (d): go-native, n ∈ {10, 100, 1000}
	for _, n := range []int{10, 100, 1000} {
		m := measureKernel(t, &kernelSpec{
			name:      fmt.Sprintf("kernel_d_%d", n),
			label:     fmt.Sprintf("(d) native n=%d", n),
			sizeLabel: fmt.Sprintf("%d", n),
			source:    `(ir.passes.pipeline/optimize-fn (ir.build/build-fn (quote (defn kmax [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (max acc i)) acc))))))`,
			stackName: "kmax",
			args: func() []vm.Value {
				return []vm.Value{vm.Int(int64(n))}
			},
			profiledSource: fmt.Sprintf(`(defn kmax_p [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (max acc i)) acc)))
(kmax_p %d)`, n),
		})
		measurements = append(measurements, m)
	}

	// Print measurement table
	t.Logf("\n=== Measurement Table ===")
	t.Logf("%10s | %6s | %10s | %10s | %12s | %12s | %12s | %8s | %8s %8s %8s",
		"Kernel", "Size", "DecodeMs", "StackOps", "ns/op stack", "ns/op spike-b", "ns/op spike-t", "Speedup", "Allocs-S", "Allocs-B", "Allocs-T")
	t.Logf("%s", strings.Repeat("-", 130))

	for _, m := range measurements {
		t.Logf("%10s | %6s | %10.3f | %10d | %12.1f | %12.1f | %12.1f | %7.2fx | %8.2f %8.2f %8.2f",
			m.Name, m.SizeLabel, m.DecodeMs, m.StackOps, m.NsOpStack, m.NsOpSpikeBoxed, m.NsOpSpikeTyped, m.Speedup,
			m.AllocsStack, m.AllocsSpikeBoxed, m.AllocsSpikeTyped)
	}

	// Go/no-go verdict
	t.Logf("\n=== Go/No-Go Verdict ===")
	kernelVerdicts := map[string][]float64{
		"kernel_a": {},
		"kernel_b": {},
		"kernel_c": {},
		"kernel_d": {},
	}
	for _, m := range measurements {
		for prefix := range kernelVerdicts {
			if strings.HasPrefix(m.Name, prefix) {
				kernelVerdicts[prefix] = append(kernelVerdicts[prefix], m.Speedup)
				break
			}
		}
	}

	for _, kernel := range []string{"kernel_a", "kernel_b", "kernel_c", "kernel_d"} {
		speedups := kernelVerdicts[kernel]
		if len(speedups) == 0 {
			continue
		}
		geomean := geometricMean(speedups)
		verdict := "YES"
		if geomean <= 1.0 {
			verdict = "NO"
		}
		t.Logf("%s: geom-mean speedup %.2fx → %s", kernel, geomean, verdict)
	}
}

// kernelSpec describes a kernel to measure.
type kernelSpec struct {
	name           string
	label          string
	sizeLabel      string
	source         string
	stackName      string
	args           func() []vm.Value
	profiledSource string
}

// measureKernel performs full measurement for one kernel × size.
// timeEngines measures ns/op for the three engines with INTERLEAVED
// rounds and per-engine medians. Two audit findings drive the shape:
// (1) single short runs are scheduler-noise-dominated → calibrate each
// engine's per-round count to >=25ms and take the median of 5 rounds;
// (2) sequential arms let an ambient load spike skew one arm's RATIO →
// interleave the engines' rounds so load hits all arms alike.
func timeEngines(bodies ...func()) []float64 {
	n := len(bodies)
	bs := make([]int, n)
	for e, body := range bodies {
		B := 64
		for {
			start := time.Now()
			for i := 0; i < B; i++ {
				body()
			}
			if time.Since(start) >= 25*time.Millisecond || B >= 1<<22 {
				break
			}
			B *= 4
		}
		bs[e] = B
	}
	const rounds = 5
	samples := make([][]float64, n)
	for r := 0; r < rounds; r++ {
		for e, body := range bodies {
			start := time.Now()
			for i := 0; i < bs[e]; i++ {
				body()
			}
			samples[e] = append(samples[e], float64(time.Since(start).Nanoseconds())/float64(bs[e]))
		}
	}
	out := make([]float64, n)
	for e := range samples {
		sort.Float64s(samples[e])
		out[e] = samples[e][rounds/2]
	}
	return out
}

func measureKernel(t *testing.T, spec *kernelSpec) KernelMeasurement {
	result := KernelMeasurement{
		Name:      spec.name,
		SizeLabel: spec.sizeLabel,
	}

	// Compile spike IR (once, outside timing)
	decodeStart := time.Now()
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("measure-spike-decode")
	_, irVal, err := c.CompileMultiple(strings.NewReader(spec.source))
	if err != nil {
		t.Fatalf("%s: compile spike IR: %v", spec.label, err)
	}
	spikeFn, decodeErr := decodeOptimizedIR(irVal)
	if decodeErr != nil {
		t.Fatalf("%s: decode spike IR: %v", spec.label, decodeErr)
	}
	result.DecodeMs = time.Since(decodeStart).Seconds() * 1000

	// Compile stack VM (once, outside timing)
	ns := rt.NS(rt.NameCoreNS)
	stackC := compiler.NewCompiler(vm.NewConsts(), ns)
	stackC.SetSource("measure-stack-compile")
	var stackSource string
	if spec.stackName == "ksum" {
		stackSource = `(defn ksum [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc)))`
	} else if spec.stackName == "m" {
		stackSource = `(defn m [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc)))`
	} else if spec.stackName == "kvec" {
		stackSource = `(defn kvec [v] (loop [i 0 acc 0] (if (< i (count v)) (recur (+ i 1) (+ acc (nth v i))) acc)))`
	} else if spec.stackName == "kmax" {
		stackSource = `(defn kmax [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (max acc i)) acc)))`
	}
	if _, _, stackErr := stackC.CompileMultiple(strings.NewReader(stackSource)); stackErr != nil {
		t.Fatalf("%s: compile stack VM: %v", spec.label, stackErr)
	}

	stackVar := ns.Lookup(vm.Symbol(spec.stackName))
	if stackVar == nil {
		t.Fatalf("%s: stack VM: %s not found", spec.label, spec.stackName)
	}
	stackFn := stackVar.(*vm.Var).Deref().(*vm.Func)
	stackChunk := stackFn.Chunk()

	// Measure stack VM ops (single run with profiling)
	vm.ResetProfile()
	vm.ProfilingEnabled.Store(true)
	stackFrame := vm.NewFrame(stackChunk, spec.args())
	_, _ = stackFrame.Run()
	vm.ProfilingEnabled.Store(false)
	for _, s := range vm.ProfileSnapshot() {
		result.StackOps += s.Count
	}

	// Measure spike ops (single run with profiling)
	vm.ResetProfile()
	vm.ProfilingEnabled.Store(true)
	statsSpike := &spikeStats{}
	_, _ = spikeRunTyped(spikeFn, spec.args(), statsSpike)
	vm.ProfilingEnabled.Store(false)
	result.SpikeOps = statsSpike.unboxOps + statsSpike.boxedArithOps + statsSpike.boxOps + statsSpike.callOps

	// Repetition count for the allocs census below.
	B := 200

	// Wall-time: calibrated median-of-rounds per engine (single short
	// runs are dominated by scheduler noise — audit finding).
	args := spec.args()
	times := timeEngines(
		func() {
			stackFrame := vm.NewFrame(stackChunk, args)
			_, _ = stackFrame.Run()
			vm.ReleaseFrame(stackFrame)
		},
		func() {
			_, _ = spikeRun(spikeFn, args)
		},
		func() {
			stats := &spikeStats{}
			_, _ = spikeRunTyped(spikeFn, args, stats)
		},
	)
	result.NsOpStack, result.NsOpSpikeBoxed, result.NsOpSpikeTyped = times[0], times[1], times[2]

	// Speedup
	result.Speedup = result.NsOpStack / result.NsOpSpikeTyped

	// Allocs: use testing.AllocsPerRun for stack VM
	result.AllocsStack = float64(testing.AllocsPerRun(B, func() {
		stackFrame := vm.NewFrame(stackChunk, args)
		_, _ = stackFrame.Run()
	}))

	// Allocs: spike boxed
	result.AllocsSpikeBoxed = float64(testing.AllocsPerRun(B, func() {
		_, _ = spikeRun(spikeFn, args)
	}))

	// Allocs: spike typed
	result.AllocsSpikeTyped = float64(testing.AllocsPerRun(B, func() {
		stats := &spikeStats{}
		_, _ = spikeRunTyped(spikeFn, args, stats)
	}))

	return result
}

// geometricMean computes the geometric mean of a slice of floats.
func geometricMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	product := 1.0
	for _, v := range values {
		product *= v
	}
	return math.Pow(product, 1.0/float64(len(values)))
}

// vecLiteral generates a Lisp literal for a vector of integers 0..n-1.
func vecLiteral(n int) string {
	var buf strings.Builder
	buf.WriteString("[")
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteString(" ")
		}
		fmt.Fprintf(&buf, "%d", i)
	}
	buf.WriteString("]")
	return buf.String()
}
