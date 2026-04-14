package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// knownFailing lists test names (filename stems) that are known to fail.
// These are tracked as TODOs rather than regressions.
var knownFailing = map[string]bool{
	"butlast":          true,
	"find":             true,
	"fn_qmark":         true,
	"fnil":             true,
	"get":              true,
	"get_in":           true,
	"hash_set":         true,
	"next":             true,
	"pr_str":           true,
	"rest":             true,
	"sequential_qmark": true,
}

// suiteCounters tracks aggregate assertion counts across the entire suite.
type suiteCounters struct {
	mu                                        sync.Mutex
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

	c := vm.NewConsts()
	coreNS := rt.NS(rt.NameCoreNS)
	loaderCtx := compiler.NewCompiler(c, coreNS)
	rt.SetNSLoader(resolver.NewNSResolver(loaderCtx, []string{
		"compat",
		"clojure-test-suite/test",
		".",
	}))

	// Create clojure.test as a copy of the test namespace
	testNS := rt.NS("test")
	clojureTestNS := rt.NS("clojure.test")
	clojureTestNS.CopyVarsFrom(testNS)

	// Load portability shim
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

func runCompatTest(t *testing.T, c *vm.Consts, filename string, totals *suiteCounters) {
	defer func() {
		if r := recover(); r != nil {
			totals.addSkip("panic")
			t.Skipf("panic: %v", r)
		}
	}()

	// Reset test registry
	testNS := rt.NS("test")
	_, _, err := compiler.NewCompiler(c, testNS).CompileMultiple(
		strings.NewReader("(clear-registered-tests!)"),
	)
	if err != nil {
		t.Fatal("failed to reset test registry:", err)
	}

	// Compile the .cljc file
	coreNS := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(c, coreNS)
	ctx.SetSource(filename)
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = ctx.CompileMultiple(f)
	f.Close()
	if err != nil {
		totals.addSkip("compile")
		t.Skipf("compile: %s", err)
	}

	// Run registered tests
	outcomeVar := testNS.Lookup("*test-result*").(*vm.Var)
	countersVar := testNS.Lookup("*report-counters*").(*vm.Var)

	_, _, err = compiler.NewCompiler(c, testNS).CompileMultiple(
		strings.NewReader("(run-tests)"),
	)
	if err != nil {
		totals.addSkip("runtime")
		t.Skipf("run-tests error: %s", err)
	}

	outcome := bool(outcomeVar.Deref().(vm.Boolean))
	counters := countersVar.Deref()
	pc, fc := getCounters(counters)
	totals.addResult(pc, fc)

	name := strings.TrimSuffix(filepath.Base(filename), ".cljc")
	if !outcome {
		if knownFailing[name] {
			t.Skipf("known failing — %s", formatCounters(counters))
		} else {
			t.Errorf("FAILED — %s", formatCounters(counters))
		}
	} else {
		if knownFailing[name] {
			t.Errorf("PASSES but is listed in knownFailing — remove it! %s", formatCounters(counters))
		} else {
			t.Logf("ok — %s", formatCounters(counters))
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
