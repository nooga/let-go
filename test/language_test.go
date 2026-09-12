/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
)

var (
	consts        *vm.Consts
	coreNS        *vm.Namespace
	cleanBindings vm.BindingSnapshot
	harnessOnce   sync.Once

	// scratchSeq names the throwaway baseline namespaces the harness tests
	// load polluting fixtures into, one fresh name per use.
	scratchSeq atomic.Int64
)

// ensureHarness builds the one shared runtime the corpus harness runs every
// file in: the constant pool, the `require` resolver, the clean
// dynamic-binding baseline, and the clojure.core namespace each file starts
// from. TestRunner and the harness unit tests both call it, so whichever runs
// first sets the runtime up and the other reuses it.
func ensureHarness() {
	harnessOnce.Do(func() {
		consts = vm.NewConsts()
		// Set up a loader so rt.NS can autoload namespaces from files during tests.
		loaderCtx := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
		// Search paths for `require`: current dir for in-tree test helpers
		// (test/test.lg etc.), pkg/rt/gogen so tests can exercise the gogen
		// macro layer, and scripts/ for the shared script-side helpers.
		paths := []string{".", "../pkg/rt/gogen", "../scripts"}
		if extra := os.Getenv("LG_SOURCE_PATHS"); extra != "" {
			for _, p := range filepath.SplitList(extra) {
				if p != "" {
					paths = append(paths, p)
				}
			}
		}
		rt.SetNSLoader(resolver.NewNSResolver(loaderCtx, paths))

		// Per-file isolation baseline: a snapshot of the (clean) dynamic-binding
		// state taken before any test file runs. Each file is executed within this
		// scope and *ns* is reset below, so a file that leaves a dynamic var dirty
		// — e.g. *ns* left pointing at a scratch namespace by an in-ns or a
		// throwing (binding [*ns* ...] ...) body — cannot corrupt the unqualified
		// symbol resolution of files that run after it in this shared runtime.
		cleanBindings = vm.SnapshotBindings()
		coreNS = rt.NS(rt.NameCoreNS)
	})
}

// runFileTests implements spec 12.2: load the file, then run the namespace
// current at the end of the file through the public clojure.test API. It
// compiles no source strings and holds no test state. Returns (success,
// message-for-failure).
//
// The namespace to run is taken from the loader's own *ns* after the file
// finishes: the compiler exposes no "saw an ns form" flag, so the documented
// approximation is that a file which switched namespaces left *ns* pointing
// at something other than the baseline it was loaded into (and other than
// `test` itself). A file that never switched is load-only — nothing in it can
// be reached by run-tests, so a deftest there is the bug the harness reports.
func runFileTests(path string) (bool, string) {
	ensureHarness()
	return runFileTestsIn(path, coreNS)
}

// runFileTestsIn is runFileTests with the baseline namespace made explicit.
// The corpus walk passes coreNS, the shared clojure.core every .lg file starts
// in. A test that deliberately loads a polluting fixture — one that refers
// vars or interns names into whatever namespace it lands in — passes a
// throwaway namespace instead, so the pollution dies with the fixture rather
// than leaking into the runtime the rest of the package shares.
func runFileTestsIn(path string, baseline *vm.Namespace) (bool, string) {
	rt.CurrentNS.SetRoot(baseline)
	var ok bool
	var msg string
	_, _ = vm.RunWithBindings(cleanBindings, func() (vm.Value, error) {
		if loadErr := runFile(path, baseline); loadErr != nil {
			msg = "load: " + loadErr.Error()
			return vm.NIL, nil
		}
		finalNS, _ := rt.CurrentNS.Deref().(*vm.Namespace)
		if finalNS == nil {
			msg = "no current namespace after loading " + path
			return vm.NIL, nil
		}
		if finalNS == baseline || finalNS.Name() == "test" {
			// Load-only file. A deftest here would never run, so that is
			// the failure, not the missing ns form.
			if hasTestVars(finalNS) {
				msg = "deftest outside a namespace: " + path
				return vm.NIL, nil
			}
			ok = true
			return vm.NIL, nil
		}
		summary, err := rt.InvokeValue(rt.LookupVar("test", "run-tests").Deref(), []vm.Value{finalNS})
		if err != nil {
			msg = "run-tests: " + err.Error()
			return vm.NIL, nil
		}
		succ, err := rt.InvokeValue(rt.LookupVar("test", "successful?").Deref(), []vm.Value{summary})
		if err != nil {
			msg = "successful?: " + err.Error()
			return vm.NIL, nil
		}
		ok = succ == vm.TRUE
		if !ok {
			msg = "summary: " + summary.String()
		}
		return vm.NIL, nil
	})
	return ok, msg
}

// hasTestVars reports whether any var interned in ns carries :test metadata —
// i.e. whether the namespace holds anything run-tests would run.
func hasTestVars(ns *vm.Namespace) bool {
	for _, v := range ns.AllVars() {
		// A Meta() that does not implement ValueAt — including the nil
		// interface a var with no metadata returns — is deliberately read as
		// "no metadata", and so as "not a test".
		m, ok := v.Meta().(interface{ ValueAt(vm.Value) vm.Value })
		if ok && m.ValueAt(vm.Keyword("test")) != vm.NIL {
			return true
		}
	}
	return false
}

// runFile compiles one .lg file with ns as the namespace it starts in.
func runFile(filename string, ns *vm.Namespace) error {
	if ns == nil {
		fmt.Println("namespace not found")
		return nil
	}
	ctx := compiler.NewCompiler(consts, ns)
	ctx.SetSource(filename)
	return rt.WithFile(filename, func() error {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		_, _, err = ctx.CompileMultiple(f)
		errc := f.Close()
		if err != nil {
			return err
		}
		if errc != nil {
			return errc
		}
		return nil
	})
}

func TestRunner(t *testing.T) {
	ensureHarness()

	file, err := os.Open("./")
	assert.NoError(t, err)
	// removed unused names := file.Readdirnames(0)
	err = file.Close()
	assert.NoError(t, err)

	// Walk the corpus one file at a time rather than loading everything and
	// running once at the end: a per-file subtest keeps Go's reporting and
	// the per-file runtime isolation runFileTests gives each file.
	err = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip directories that contain non-test .lg files (e.g.
			// compat/ holds the corpus runner, which is invoked manually;
			// benches/ holds standalone benchmarks + bug-repro fixtures that
			// run for seconds and assert nothing — run by hand with ./lg;
			// gogen/ holds native-lowering harness fixtures driven by the
			// deftype_skeleton_lowering_e2e_test, not bytecode deftests).
			// native-entry/ holds AOT gate fixtures (programs with -main, not
			// deftests) driven by TestNativeEntryASTGate; tools/ holds
			// test-scoped generators invoked with arguments by e2e tests.
			if info.Name() == "compat" || info.Name() == "clojure-test-suite" || info.Name() == "benches" || info.Name() == "gogen" || info.Name() == "gold-aot" || info.Name() == "native-entry" || info.Name() == "tools" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".lg" {
			return nil
		}
		name := info.Name()
		t.Run(name, func(t *testing.T) {
			ok, msg := runFileTests(path)
			assert.True(t, ok, "some tests failed in "+name+": "+msg)
		})
		return nil
	})
	assert.NoError(t, err)
}

// TestHarnessRejectsDeftestOutsideNamespace pins the load-only rule of spec
// 12.2: a file with no `ns` form is loaded but never run through run-tests,
// so a deftest in one would silently never execute. That is the failure the
// harness reports — not the missing ns form.
func TestHarnessRejectsDeftestOutsideNamespace(t *testing.T) {
	ensureHarness()

	dir := t.TempDir()
	path := filepath.Join(dir, "stray.lg")
	// :refer :all is not cosmetic here: deftest's own template splices a
	// bare, unqualified `is` (see test.lg's note on this dialect's
	// syntax-quote), so (test/deftest ... (test/is ...)) does not compile.
	if err := os.WriteFile(path, []byte("(require '[test :refer :all])\n(deftest stray (is true))\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// A load-only file interns and refers into whatever namespace it is
	// loaded in, and this one refers every public of `test` — deftest, is,
	// testing, are, use-fixtures. Loading it into the shared clojure.core
	// would make those resolvable unqualified in every corpus file that runs
	// afterwards, silently masking a file that used them without requiring
	// them. Give it a throwaway baseline instead, so the whole blast radius
	// is a namespace nothing else ever looks at. The name is unique per run
	// so -count=N cannot reuse a dirtied one.
	scratch := rt.DefNSBare(fmt.Sprintf("test.harness-load-only-scratch-%d", scratchSeq.Add(1)))

	// Guard against a future reordering or refactor quietly reintroducing the
	// dependency on this test running after TestRunner. It is a before/after
	// diff, not an absolute check: corpus files of the same shape already
	// refer test's publics into clojure.core when TestRunner has run first,
	// so what must hold is that THIS fixture changes none of these mappings —
	// which is true in any order only while it loads somewhere else.
	core := rt.NS(rt.NameCoreNS)
	watched := []vm.Symbol{"stray", "deftest", "is", "testing", "use-fixtures"}
	before := make(map[vm.Symbol]vm.Value, len(watched))
	for _, sym := range watched {
		before[sym] = core.Lookup(sym)
	}

	ok, msg := runFileTestsIn(path, scratch)
	if ok || !strings.Contains(msg, "deftest outside a namespace") {
		t.Fatalf("expected load-only failure, got ok=%v msg=%q", ok, msg)
	}

	for _, sym := range watched {
		if got := core.Lookup(sym); got != before[sym] {
			t.Errorf("fixture leaked %q into the shared clojure.core baseline "+
				"(was %v, now %v); it must load into a throwaway namespace",
				sym, before[sym], got)
		}
	}
}
