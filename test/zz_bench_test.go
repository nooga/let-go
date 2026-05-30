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

// BenchmarkClojureTestSuite measures the end-to-end wall time of
// running the full clojure-test-suite (a.k.a. "jank") through let-go:
// compile each .cljc file under clojure-test-suite/test/clojure and
// invoke its registered (deftest …) forms. This is the load-bearing
// "real Lisp at scale" benchmark for the ratchet — micro-benchmarks
// in pkg/vm catch primitive-level regressions; this one catches
// end-to-end regressions in the compile + run path.
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
	b.ReportMetric(float64(counters.files.Load())/n, "files")
	skips := counters.skipCompile.Load() + counters.skipPanic.Load() +
		counters.skipTimeout.Load() + counters.skipMem.Load()
	b.ReportMetric(float64(skips)/n, "skips")
}

// benchCounters is a lock-free assertion tally aggregated across all
// iterations of BenchmarkClojureTestSuite. atomic.Int64s keep the
// per-file goroutines from contending on a mutex.
type benchCounters struct {
	files       atomic.Int64
	pass        atomic.Int64
	fail        atomic.Int64
	skipCompile atomic.Int64
	skipPanic   atomic.Int64
	skipTimeout atomic.Int64
	skipMem     atomic.Int64
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
func runCompatTestBench(c *vm.Consts, filename string, counters *benchCounters) {
	type result struct {
		err       error
		isPanic   bool
		passCount int
		failCount int
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

		testNS := rt.NS("test")
		_, _, err := compiler.NewCompiler(c, testNS).CompileMultiple(
			strings.NewReader("(clear-registered-tests!)"),
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

		countersVar := testNS.Lookup("*report-counters*").(*vm.Var)
		_, _, _ = compiler.NewCompiler(c, testNS).CompileMultiple(
			strings.NewReader("(run-tests)"),
		)
		pc, fc := getCounters(countersVar.Deref())
		ch <- result{passCount: pc, failCount: fc}
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
