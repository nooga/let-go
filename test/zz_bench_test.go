package test

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// BenchmarkClojureTestSuite measures the EXECUTION wall time of
// running the full clojure-test-suite (a.k.a. "jank") through let-go:
// for each .cljc file under clojure-test-suite/test/clojure, invoke its
// registered (deftest …) forms on an already-bootstrapped runtime. This
// is the load-bearing "real Lisp at scale" benchmark for the ratchet —
// micro-benchmarks in pkg/vm catch primitive-level regressions; this one
// catches end-to-end regressions in how compiled code RUNS.
//
// Compilation (IR generation) of each file is deliberately EXCLUDED from
// the measured ns/op: the headline number tracks execution speed only.
// Compile cost is the job of BenchmarkIRCompile (pkg/ir), and is also
// surfaced here as the informational "compile_ms" metric. Keeping the two
// concerns separate is the whole point — under -tags gogen_ir the suite
// runs core as native Go, so the ns/op delta vs the bytecode variant is a
// pure runtime comparison, not muddied by compile time.
//
// Three execution modes, selected by the LG_SUITE_IR env var crossed with
// the build tag (bench-ratchet emits one job per mode):
//
//	bytecode     LG_SUITE_IR unset, untagged   — each test fn body compiles
//	                                              via the direct Go path to
//	                                              bytecode. No IR.
//	ir_bytecode  LG_SUITE_IR=1, untagged       — *ir-compile* on: every body
//	                                              routes through the IR
//	                                              optimizer; the passes
//	                                              themselves run as bytecode.
//	aot_native   LG_SUITE_IR=1, -tags gogen_ir — *ir-compile* on AND the IR
//	                                              passes dispatch to native Go
//	                                              (the lowered core tree).
//
// All three measure EXECUTION wall time (ns/op); the extra compile cost the
// IR modes pay is surfaced separately as compile_ms. *ir-compile* is
// known-partial over arbitrary code (core.lg:890) — forms it can't lower
// fall back / skip, so ir_bytecode/aot_native legitimately show higher skip
// counts than bytecode. That divergence is signal: the trend tracks IR
// pipeline coverage as much as speed.
//
// One iteration runs the WHOLE suite (~hundreds of files). With
// -benchtime 1s Go runs a single iteration and reports its wall time
// as ns/op; -count > 1 gives variance bands.
//
// Submodule absent → b.Skip (so CI without the submodule still passes).
// Per-file 5s timeout (same as TestClojureTestSuite) protects against
// runaway files. Per-file panics/compile-errors are counted, not fatal —
// what we measure is the total time spent.
//
// Custom metrics reported alongside ns/op. These are PER-SUITE-RUN
// counts averaged over b.N iterations — `pass=5621` means each suite
// run produces 5621 passing assertions. We deliberately drop the
// `/op` suffix because in this bench one "op" is one full suite
// pass, and "5621 pass/op" reads as a ratio when it's really an
// absolute per-iteration count.
//
//	pass    — assertion passes per suite run
//	fail    — assertion failures per suite run
//	files   — files that produced a test result per suite run
//	skips   — files skipped (compile error, panic, timeout, mem-cap)
//
// Why these matter: wall time alone can't catch a correctness
// regression that flips passes into skips/errors with the same overall
// runtime. The ratchet currently checks ns/op + ratio_to_anchor; the
// custom metrics surface in the raw `go test -bench` output and in
// streaming .jsonl records when a future bench-ratchet learns to parse
// them. Right now they're documentation + dump-for-grep.
//
// Benchmark pair:
//   - BenchmarkClojureTestSuite              → execution-only cost
//   - BenchmarkClojureTestSuiteCompileAndRun → total request cost
//
// The latter answers the "IR crossover" question: once the IR pipeline itself
// runs natively, is routing forms through *ir-compile* no worse than the
// direct bytecode compiler on realistic jank-sized code?
func BenchmarkClojureTestSuiteCompileAndRun(b *testing.B) {
	compiler.SetMatchCljConditional(true)
	defer compiler.SetMatchCljConditional(false)

	suiteRoot := "clojure-test-suite/test/clojure"
	if _, err := os.Stat(filepath.Join(suiteRoot, "core_test")); os.IsNotExist(err) {
		b.Skip("clojure-test-suite submodule not initialized (run: git submodule update --init)")
	}

	var files []string
	for _, dir := range []string{"core_test", "string_test"} {
		matches, err := filepath.Glob(filepath.Join(suiteRoot, dir, "*.cljc"))
		if err != nil {
			b.Fatal(err)
		}
		files = append(files, matches...)
	}
	if len(files) == 0 {
		b.Fatal("no .cljc files found in ", suiteRoot)
	}

	origStdout := os.Stdout
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		b.Fatal("open /dev/null:", err)
	}
	os.Stdout = devnull
	defer func() {
		os.Stdout = origStdout
		devnull.Close()
	}()

	var counters benchCounters

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runSuiteOnce(b, files, &counters)
	}
	b.StopTimer()

	n := float64(b.N)
	if n < 1 {
		n = 1
	}
	b.ReportMetric(float64(counters.pass.Load())/n, "pass")
	b.ReportMetric(float64(counters.fail.Load())/n, "fail")
	b.ReportMetric(float64(counters.errs.Load())/n, "error")
	b.ReportMetric(float64(counters.tests.Load())/n, "tests")
	b.ReportMetric(float64(counters.files.Load())/n, "files")
	skips := counters.skipCompile.Load() + counters.skipPanic.Load() +
		counters.skipTimeout.Load() + counters.skipMem.Load()
	b.ReportMetric(float64(skips)/n, "skips")
	const msPerNs = 1.0 / 1e6
	b.ReportMetric(float64(counters.compileNanos.Load())*msPerNs/n, "compile_ms")
	b.ReportMetric(float64(counters.runNanos.Load())*msPerNs/n, "run_ms")
}

func BenchmarkClojureTestSuite(b *testing.B) {
	compiler.SetMatchCljConditional(true)
	defer compiler.SetMatchCljConditional(false)

	suiteRoot := "clojure-test-suite/test/clojure"
	if _, err := os.Stat(filepath.Join(suiteRoot, "core_test")); os.IsNotExist(err) {
		b.Skip("clojure-test-suite submodule not initialized (run: git submodule update --init)")
	}

	// Discover files outside the timer.
	var files []string
	for _, dir := range []string{"core_test", "string_test"} {
		matches, err := filepath.Glob(filepath.Join(suiteRoot, dir, "*.cljc"))
		if err != nil {
			b.Fatal(err)
		}
		files = append(files, matches...)
	}
	if len(files) == 0 {
		b.Fatal("no .cljc files found in ", suiteRoot)
	}

	// Redirect stdout to /dev/null for the duration of the benchmark.
	// The Lisp test framework's `(run-tests)` prints PASS/FAIL lines
	// per assertion — hundreds of thousands of lines — which both slows
	// the bench down and, more critically, corrupts go test's benchmark
	// output line ("BenchmarkX-N\t…\tN ns/op") by inserting our prints
	// between the bench prefix and its measurement, breaking the parser
	// in cmd/bench-ratchet.
	origStdout := os.Stdout
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		b.Fatal("open /dev/null:", err)
	}
	os.Stdout = devnull
	defer func() {
		os.Stdout = origStdout
		devnull.Close()
	}()

	var counters benchCounters

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runSuiteOnce(b, files, &counters)
	}
	b.StopTimer()

	// Report per-iteration custom metrics. Counters are summed across
	// b.N iterations; divide so each metric reads as "per outer iter."
	n := float64(b.N)
	if n < 1 {
		n = 1
	}
	b.ReportMetric(float64(counters.pass.Load())/n, "pass")
	b.ReportMetric(float64(counters.fail.Load())/n, "fail")
	// :error is the IR modes' dominant failure outcome — an assertion that
	// threw and was caught by run-tests' runner. Counting it (not just pass +
	// fail) is what keeps the headline honest: pass alone collapses under IR
	// because the missing assertions ERRORED, they didn't vanish.
	b.ReportMetric(float64(counters.errs.Load())/n, "error")
	b.ReportMetric(float64(counters.tests.Load())/n, "tests")
	b.ReportMetric(float64(counters.files.Load())/n, "files")
	skips := counters.skipCompile.Load() + counters.skipPanic.Load() +
		counters.skipTimeout.Load() + counters.skipMem.Load()
	b.ReportMetric(float64(skips)/n, "skips")
	// Override the headline ns/op with EXECUTION-only wall time (the
	// worker's own clock), so the timeline tracks how fast compiled code
	// runs — not how long it took to compile. Compile cost is surfaced
	// separately as compile_ms and is the dedicated job of BenchmarkIRCompile.
	// ReportMetric with the "ns/op" unit replaces the built-in timer value.
	b.ReportMetric(float64(counters.runNanos.Load())/n, "ns/op")
	const msPerNs = 1.0 / 1e6
	b.ReportMetric(float64(counters.compileNanos.Load())*msPerNs/n, "compile_ms")
}

// benchCounters is a lock-free assertion tally aggregated across all
// iterations of BenchmarkClojureTestSuite. atomic.Int64s keep the
// per-file goroutines from contending on a mutex.
type benchCounters struct {
	files       atomic.Int64
	pass        atomic.Int64
	fail        atomic.Int64
	errs        atomic.Int64 // run-tests :error count (assertion threw, caught by runner)
	tests       atomic.Int64 // run-tests :test count (deftest fns actually invoked)
	skipCompile atomic.Int64
	skipPanic   atomic.Int64
	skipTimeout atomic.Int64
	skipMem     atomic.Int64
	// Wall-time split (nanoseconds), summed across files. compileNanos is
	// excluded from the bench timer; runNanos mirrors what the timer counts.
	compileNanos atomic.Int64
	runNanos     atomic.Int64
}

// suiteIRMode reports whether LG_SUITE_IR selects an IR-compile variant
// (ir_bytecode / aot_native). Unset or "0"/"false" → the plain bytecode mode.
func suiteIRMode() bool {
	v := os.Getenv("LG_SUITE_IR")
	return v == "1" || v == "true"
}

func runSuiteOnce(b *testing.B, files []string, counters *benchCounters) {
	origLoader := rt.GetNSLoader()
	defer rt.SetNSLoader(origLoader)

	c := vm.NewConsts()
	coreNS := rt.NS(rt.NameCoreNS)
	loaderCtx := compiler.NewCompiler(c, coreNS)
	rt.SetNSLoader(resolver.NewNSResolver(loaderCtx, []string{
		"compat",
		"clojure-test-suite/test",
		".",
	}))

	portCtx := compiler.NewCompiler(c, coreNS)
	portCtx.SetSource("compat/clojure/core-test/portability.lg")
	pf, err := os.Open("compat/clojure/core-test/portability.lg")
	if err != nil {
		b.Fatal("open portability shim:", err)
	}
	_, _, err = portCtx.CompileMultiple(pf)
	pf.Close()
	if err != nil {
		b.Fatal("compile portability shim:", err)
	}

	// Mode selection (see BenchmarkClojureTestSuite doc): when LG_SUITE_IR is
	// set, route every test-file fn body through the IR-optimizing pipeline by
	// turning on *ir-compile*. Untagged → the passes run as bytecode
	// (ir_bytecode); under -tags gogen_ir → they dispatch to native Go
	// (aot_native). Unset → direct bytecode compile. This is setup, not part
	// of the measured compile/run split, so it stays out of the counters.
	if suiteIRMode() {
		irCtx := compiler.NewCompiler(c, coreNS)
		if _, _, err := irCtx.CompileMultiple(strings.NewReader(
			"(require 'ir.passes.pipeline)\n(set! *ir-compile* true)",
		)); err != nil {
			b.Fatal("enable *ir-compile*:", err)
		}
		defer func() {
			reset := compiler.NewCompiler(c, coreNS)
			_, _, _ = reset.CompileMultiple(strings.NewReader("(set! *ir-compile* false)"))
		}()
	}

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".cljc")
		if name == "portability" {
			continue
		}
		runCompatTestBench(c, file, counters)
	}
}

// runCompatTestBench is the benchmark-friendly twin of runCompatTest.
// Differences:
//   - takes a *vm.Consts directly, no testing.T
//   - all failures (compile error, panic, timeout, mem-limit) are
//     recorded as atomic counters; the bench keeps going
//   - per-assertion counts are also recorded via the Lisp framework's
//     *report-counters* on success
//
// Behavior otherwise matches runCompatTest: 5s per-file timeout, 512MB
// per-file memory ceiling, panic recovery in a goroutine.
//
// Timing: the worker measures compile and execution wall time separately
// (compileNanos / runNanos). The benchmark's headline ns/op is overridden
// in BenchmarkClojureTestSuite to report runNanos only — compile time is
// surfaced as the informational "compile_ms" metric. We deliberately do
// NOT use b.StopTimer/StartTimer here: with a per-file watchdog goroutine
// the start/stop transitions race the result, so we drive the reported
// metric straight from the worker's own clock instead.
func runCompatTestBench(c *vm.Consts, filename string, counters *benchCounters) {
	type result struct {
		err       error
		isPanic   bool
		passCount int
		failCount int
		errCount  int
		testCount int
	}
	ch := make(chan result, 1)
	runtime.GC()
	baseAlloc := currentAlloc()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				_ = string(debug.Stack())
				ch <- result{err: anyToErr(r), isPanic: true}
			}
		}()

		compileStart := time.Now()
		testNS := rt.NS("test")
		// Reset per-file test state. Crucially this clears *once-fixtures* /
		// *each-fixtures* too: clear-registered-tests! only resets the test
		// registry, but run-tests wraps ALL execution in the once-fixtures, and
		// those leak across files (a fixture registered by file N would wrap
		// file N+1's tests). That leak is harmless when every fixture runs
		// correctly (bytecode), but under *ir-compile* a single misbehaving
		// fixture escapes run-tests and silently zeros every later file. Each
		// .cljc is an independent namespace, so its fixtures must not persist.
		_, _, err := compiler.NewCompiler(c, testNS).CompileMultiple(
			strings.NewReader("(clear-registered-tests!) (set! *once-fixtures* []) (set! *each-fixtures* [])"),
		)
		if err != nil {
			ch <- result{err: err}
			return
		}

		rt.DefNSBare(nsNameFromCompatPath(filename))

		coreNS := rt.NS(rt.NameCoreNS)
		ctx := compiler.NewCompiler(c, coreNS)
		ctx.SetSource(filename)
		f, err := os.Open(filename)
		if err != nil {
			ch <- result{err: err}
			return
		}
		err = compileProtected(ctx, f)
		f.Close()
		if err != nil {
			ch <- result{err: err, isPanic: strings.HasPrefix(err.Error(), "panic:")}
			return
		}
		counters.compileNanos.Add(time.Since(compileStart).Nanoseconds())

		// (run-tests) executes the already-compiled deftest fns. This is the
		// execution phase — the only thing the headline ns/op reflects.
		runStart := time.Now()
		countersVar := testNS.Lookup("*report-counters*").(*vm.Var)
		_, _, _ = compiler.NewCompiler(c, testNS).CompileMultiple(
			strings.NewReader("(run-tests)"),
		)
		pc, fc, ec, tc := getCountersFull(countersVar.Deref())
		counters.runNanos.Add(time.Since(runStart).Nanoseconds())
		ch <- result{passCount: pc, failCount: fc, errCount: ec, testCount: tc}
	}()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.After(compatTestTimeout)

	for {
		select {
		case res := <-ch:
			switch {
			case res.err != nil && res.isPanic:
				counters.skipPanic.Add(1)
			case res.err != nil:
				counters.skipCompile.Add(1)
			default:
				counters.files.Add(1)
				counters.pass.Add(int64(res.passCount))
				counters.fail.Add(int64(res.failCount))
				counters.errs.Add(int64(res.errCount))
				counters.tests.Add(int64(res.testCount))
			}
			return
		case <-deadline:
			counters.skipTimeout.Add(1)
			return
		case <-ticker.C:
			runtime.GC()
			if currentAlloc()-baseAlloc > memLimitBytes {
				runtime.GC()
				counters.skipMem.Add(1)
				return
			}
		}
	}
}

func anyToErr(v any) error {
	if err, ok := v.(error); ok {
		return err
	}
	return &simpleErr{msg: "panic: " + reprAny(v)}
}

type simpleErr struct{ msg string }

func (e *simpleErr) Error() string { return e.msg }

func reprAny(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case error:
		return x.Error()
	default:
		return ""
	}
}
