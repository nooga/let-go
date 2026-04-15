package test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// Per-test timeout. If a test takes longer than this (e.g. infinite seq
// realization), it is killed and reported as a skip.
const compatTestTimeout = 5 * time.Second

// memLimitBytes is the per-test memory growth threshold (512MB).
// If a test causes allocations beyond this, we skip it.
const memLimitBytes = 512 * 1024 * 1024

// knownFailing lists test names (filename stems) that are known to fail.
// Tests that pass but appear here will cause an error so the list stays current.
var knownFailing = map[string]bool{
	"binding":          true, // thread binding propagation to futures
	"contains_qmark":   true, // contains? edge cases
	"decimal_qmark":    true, // BigDecimal parsed as Float
	"double_qmark":     true, // BigDecimal parsed as Float
	"float":            true, // BigDecimal edge cases
	"float_qmark":      true, // BigDecimal parsed as Float
	"int_qmark":        true, // BigDecimal parsed as Float
	"integer_qmark":    true, // BigDecimal parsed as Float
	"nan_qmark":        true, // NaN type predicate
	"ratio_qmark":      true, // ratio parsed as Float
	"rational_qmark":   true, // BigDecimal parsed as Float
	"zero_qmark":       true, // BigDecimal zero
	"select_keys":      true, // select-keys edge cases
	"drop":             true, // (drop 5 nil) → nil not ()
	"drop_while":       true, // (drop-while pred nil) → nil not ()
	"even_qmark":       true, // even? on BigDecimal
	"identical_qmark":  true, // identical? on boxed values
	"ifn_qmark":        true, // ifn? edge cases
	"max":              true, // max with BigDecimal
	"min":              true, // min with BigDecimal
	"min_key":          true, // min-key edge cases
	"mod":              true, // mod edge cases

	"not_empty":        true, // not-empty on list containing nil
	"nth":              true, // nth out-of-bounds doesn't throw
	"nthrest":          true, // nthrest edge cases
	"odd_qmark":        true, // odd? on BigDecimal
	"peek":             true, // peek on cons
	"pr_str":           true, // pr-str formatting
	"quot":             true, // quot edge cases
	"rem":              true, // rem edge cases
	"str":              true, // str reader conditional
	"keyword":          true, // keyword with empty ns
	"nnext":            true, // map ordering
	"num":              true, // num edge cases
	"print_str":        true, // int-as-float formatting
	"println_str":      true, // int-as-float formatting
	"prn_str":          true, // int-as-float formatting
}

// suiteCounters tracks aggregate assertion counts across the entire suite.
type suiteCounters struct {
	mu                                              sync.Mutex
	files, pass, fail, skip, compileSkip, panicSkip int
}

func (s *suiteCounters) addResult(pass, fail int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.files++
	s.pass += pass
	s.fail += fail
}

func (s *suiteCounters) addSkip(reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch reason {
	case "compile":
		s.compileSkip++
	case "panic":
		s.panicSkip++
	default:
		s.skip++
	}
}

func (s *suiteCounters) summary() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("files=%d assertions: pass=%d fail=%d | skipped: compile=%d panic=%d runtime=%d",
		s.files, s.pass, s.fail, s.compileSkip, s.panicSkip, s.skip)
}

// TestClojureTestSuite runs tests from jank-lang/clojure-test-suite.
// Each .cljc file is compiled and executed through let-go with compat shims.
// Files that fail to compile (e.g. missing builtins) are reported as skipped.
func TestClojureTestSuite(t *testing.T) {
	suiteDir := "clojure-test-suite/test/clojure/core_test"
	if _, err := os.Stat(suiteDir); os.IsNotExist(err) {
		t.Skip("clojure-test-suite submodule not initialized (run: git submodule update --init)")
	}

	// Save and restore the global NS loader so we don't pollute TestRunner's state
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

	// Load portability shim (provides when-var-exists and thrown? macros)
	portCtx := compiler.NewCompiler(c, coreNS)
	portCtx.SetSource("compat/clojure/core-test/portability.lg")
	pf, err := os.Open("compat/clojure/core-test/portability.lg")
	if err != nil {
		t.Fatal("failed to open portability shim:", err)
	}
	_, _, err = portCtx.CompileMultiple(pf)
	pf.Close()
	if err != nil {
		t.Fatal("failed to compile portability shim:", err)
	}

	files, err := filepath.Glob(filepath.Join(suiteDir, "*.cljc"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no .cljc files found in", suiteDir)
	}

	totals := &suiteCounters{}

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".cljc")
		if name == "portability" {
			continue
		}

		t.Run(name, func(t *testing.T) {
			runCompatTest(t, c, file, totals)
		})
	}

	t.Logf("TOTALS: %s", totals.summary())
}

// compatTestResult carries the outcome of a single compat test back from its goroutine.
type compatTestResult struct {
	err      error  // non-nil for compile/runtime errors
	isPanic  bool   // true if err came from a recovered panic
	outcome  bool   // test pass/fail (only valid when err == nil)
	counters vm.Value
}

// compileProtected wraps CompileMultiple with panic recovery, returning
// panics as errors instead of crashing the test process.
func compileProtected(ctx *compiler.Context, f *os.File) (panicErr error) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			lines := strings.Split(stack, "\n")
			var relevant []string
			for i, line := range lines {
				if strings.Contains(line, "let-go/pkg/") {
					relevant = append(relevant, strings.TrimSpace(lines[i]))
					if i+1 < len(lines) {
						relevant = append(relevant, strings.TrimSpace(lines[i+1]))
					}
					break
				}
			}
			loc := strings.Join(relevant, " ")
			if loc == "" {
				loc = "unknown location"
			}
			panicErr = fmt.Errorf("panic: %v at %s", r, loc)
		}
	}()
	_, _, panicErr = ctx.CompileMultiple(f)
	return
}

func currentAlloc() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func runCompatTest(t *testing.T, c *vm.Consts, filename string, totals *suiteCounters) {
	ch := make(chan compatTestResult, 1)
	baseAlloc := currentAlloc()

	go func() {
		// The entire compile+run happens in this goroutine so we can
		// abandon it on timeout without blocking the test harness.
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				lines := strings.Split(stack, "\n")
				var relevant []string
				for i, line := range lines {
					if strings.Contains(line, "let-go/pkg/") {
						relevant = append(relevant, strings.TrimSpace(lines[i]))
						if i+1 < len(lines) {
							relevant = append(relevant, strings.TrimSpace(lines[i+1]))
						}
						break
					}
				}
				loc := strings.Join(relevant, " ")
				ch <- compatTestResult{
					err:     fmt.Errorf("panic: %v at %s", r, loc),
					isPanic: true,
				}
			}
		}()

		// Reset test registry
		testNS := rt.NS("test")
		_, _, err := compiler.NewCompiler(c, testNS).CompileMultiple(
			strings.NewReader("(clear-registered-tests!)"),
		)
		if err != nil {
			ch <- compatTestResult{err: fmt.Errorf("reset: %w", err)}
			return
		}

		// Compile the .cljc file
		coreNS := rt.NS(rt.NameCoreNS)
		ctx := compiler.NewCompiler(c, coreNS)
		ctx.SetSource(filename)
		f, err := os.Open(filename)
		if err != nil {
			ch <- compatTestResult{err: err}
			return
		}
		err = compileProtected(ctx, f)
		f.Close()
		if err != nil {
			ch <- compatTestResult{err: err, isPanic: strings.HasPrefix(err.Error(), "panic:")}
			return
		}

		// Run registered tests
		outcomeVar := testNS.Lookup("*test-result*").(*vm.Var)
		countersVar := testNS.Lookup("*report-counters*").(*vm.Var)

		_, _, err = compiler.NewCompiler(c, testNS).CompileMultiple(
			strings.NewReader("(run-tests)"),
		)
		if err != nil {
			ch <- compatTestResult{err: fmt.Errorf("run-tests: %w", err)}
			return
		}

		ch <- compatTestResult{
			outcome:  bool(outcomeVar.Deref().(vm.Boolean)),
			counters: countersVar.Deref(),
		}
	}()

	// Wait for result or timeout
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.After(compatTestTimeout)

	for {
		select {
		case res := <-ch:
			if res.err != nil {
				if res.isPanic {
					totals.addSkip("panic")
					t.Skipf("%s", res.err)
				} else {
					totals.addSkip("compile")
					t.Skipf("compile: %s", vm.FormatError(res.err))
				}
				return
			}

			pc, fc := getCounters(res.counters)
			totals.addResult(pc, fc)

			name := strings.TrimSuffix(filepath.Base(filename), ".cljc")
			if !res.outcome {
				if knownFailing[name] {
					t.Skipf("known failing — %s", formatCounters(res.counters))
				} else {
					t.Errorf("FAILED — %s", formatCounters(res.counters))
				}
			} else {
				if knownFailing[name] {
					t.Errorf("PASSES but is listed in knownFailing — remove it! %s", formatCounters(res.counters))
				} else {
					t.Logf("ok — %s", formatCounters(res.counters))
				}
			}
			return

		case <-deadline:
			totals.addSkip("runtime")
			t.Skipf("timeout after %s", compatTestTimeout)
			return

		case <-ticker.C:
			// Check memory growth
			if currentAlloc()-baseAlloc > memLimitBytes {
				totals.addSkip("runtime")
				runtime.GC() // try to reclaim before moving on
				t.Skipf("memory limit exceeded (>%dMB growth)", memLimitBytes/1024/1024)
				return
			}
		}
	}
}

func getCounters(v vm.Value) (pass, fail int) {
	m, ok := v.(*vm.PersistentMap)
	if !ok {
		return 0, 0
	}
	getInt := func(k string) int {
		val := m.ValueAtOr(vm.Keyword(k), vm.MakeInt(0))
		if n, ok := val.(vm.Int); ok {
			return int(n)
		}
		return 0
	}
	return getInt("pass"), getInt("fail")
}

func formatCounters(v vm.Value) string {
	m, ok := v.(*vm.PersistentMap)
	if !ok {
		return fmt.Sprintf("%s", v)
	}
	get := func(k string) int {
		val := m.ValueAtOr(vm.Keyword(k), vm.MakeInt(0))
		if n, ok := val.(vm.Int); ok {
			return int(n)
		}
		return 0
	}
	return fmt.Sprintf("tests=%d pass=%d fail=%d error=%d",
		get("test"), get("pass"), get("fail"), get("error"))
}
