/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
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
		// (test/test.lg etc.), pkg/rt/gogen so tests can exercise the
		// gogen macro layer, and the repo's scripts/ root so test files can
		// require the tooling namespaces that live there (spec-evidence.lg,
		// spec-evidence/*.lg, quality.*). scripts/ is UNCONDITIONAL: those tests
		// are part of the default suite, and requiring an env var to make
		// them compile meant a plain `go test ./test/` — which is what CI
		// and the pre-push hook run — failed with "unable to load namespace
		// spec-evidence". `go test ./test/` runs with cwd set to this
		// package's directory (test/), one level below the repo root, so
		// each relative entry is resolved relative to "..".
		nsPath := []string{".", "../pkg/rt/gogen", "../scripts"}
		// LG_SOURCE_PATHS still adds FURTHER namespace roots on top, for
		// one-off runs against an out-of-tree corpus. This harness runs
		// in-process rather than shelling out to ./lg, so the resolver never
		// otherwise reads the env var. An absolute entry is used as-is.
		for _, p := range resolver.ParseSearchPaths(os.Getenv("LG_SOURCE_PATHS")) {
			if !filepath.IsAbs(p) {
				p = filepath.Join("..", p)
			}
			nsPath = append(nsPath, p)
		}
		rt.SetNSLoader(resolver.NewNSResolver(loaderCtx, nsPath))

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

// testMetaValue returns the var's :test metadata value, or vm.NIL when it has
// none. A Meta() that does not implement ValueAt — including the nil interface
// a var with no metadata returns — is deliberately read as "no metadata", and
// so as "not a test".
func testMetaValue(v *vm.Var) vm.Value {
	m, ok := v.Meta().(interface{ ValueAt(vm.Value) vm.Value })
	if !ok {
		return vm.NIL
	}
	return m.ValueAt(vm.Keyword("test"))
}

// boundRoot reads a var's root without tripping over an unbound one, so the
// snapshot can record a test-ns-hook var that exists but holds no value yet.
func boundRoot(v *vm.Var) vm.Value {
	if !v.IsBound() {
		return vm.NIL
	}
	return v.Root()
}

// sameValue is the snapshot's identity test. vm.Value is an interface, and a
// few implementations (vectors, maps) are non-comparable Go types that would
// panic under ==, so an uncomparable value is conservatively read as "changed"
// rather than crashing the harness. Every :test value deftest produces is an
// Fn pointer, so this only guards against hand-attached metadata.
func sameValue(a, b vm.Value) bool {
	ta, tb := reflect.TypeOf(a), reflect.TypeOf(b)
	if ta != tb {
		return false
	}
	if ta == nil {
		return true
	}
	if !ta.Comparable() {
		return false
	}
	return a == b
}

// testSnapshot implements spec 12.2's test_snapshot: every var, in every
// namespace, that carries :test metadata mapped to that :test value, plus
// every namespace's test-ns-hook var mapped to its own value. A namespace
// whose only entry point is test-ns-hook interns no :test var, so without the
// second half testedNamespaces would never see it touched (spec 9.3).
func testSnapshot() map[*vm.Var]vm.Value {
	out := map[*vm.Var]vm.Value{}
	for _, ns := range rt.AllNSes() {
		for sym, v := range ns.AllVars() {
			if t := testMetaValue(v); t != vm.NIL {
				out[v] = t
				continue
			}
			if sym == vm.Symbol("test-ns-hook") {
				out[v] = boundRoot(v)
			}
		}
	}
	return out
}

// testedNamespaces implements spec 12.2's tested_namespaces: the namespaces
// holding at least one var whose :test value, or whose test-ns-hook value, is
// absent from `before` or not identical to what was recorded there — i.e. the
// namespaces in which this load defined or redefined a deftest or a
// test-ns-hook. Ordered by namespace name, so the run order is deterministic.
func testedNamespaces(before map[*vm.Var]vm.Value) []*vm.Namespace {
	var names []string
	byName := map[string]*vm.Namespace{}
	for name, ns := range rt.AllNSes() {
		touched := false
		for sym, v := range ns.AllVars() {
			cur := testMetaValue(v)
			if cur == vm.NIL && sym == vm.Symbol("test-ns-hook") {
				cur = boundRoot(v)
			}
			if cur == vm.NIL {
				continue
			}
			if prev, seen := before[v]; !seen || !sameValue(prev, cur) {
				touched = true
				break
			}
		}
		if touched {
			names = append(names, name)
			byName[name] = ns
		}
	}
	sort.Strings(names)
	out := make([]*vm.Namespace, 0, len(names))
	for _, n := range names {
		out = append(out, byName[n])
	}
	return out
}

// fileKind is spec 12.2's FileKind.
type fileKind int

const (
	kindTests fileKind = iota
	kindLoadOnly
	kindStrayDeftest
)

// classifyLoaded implements spec 12.2's classify_loaded. It reads only which
// namespaces gained a :test or test-ns-hook var since `before`, never how the
// file reached them, so a file using `ns` and one using only `in-ns` to reach
// the same namespace classify identically.
func classifyLoaded(initialNS *vm.Namespace, before map[*vm.Var]vm.Value) (fileKind, []*vm.Namespace) {
	touched := testedNamespaces(before)
	for _, ns := range touched {
		if ns == initialNS {
			// Harness policy, not clojure.test: the starting namespace is
			// shared by every file, so a deftest interned there cannot be
			// attributed to one. This fires even when the file also went on
			// to touch a namespace of its own.
			return kindStrayDeftest, nil
		}
	}
	if len(touched) == 0 {
		return kindLoadOnly, nil
	}
	return kindTests, touched
}

// runFileTests implements spec 12.2: snapshot the runtime's test vars, load
// the file, then run every namespace the load defined or redefined a deftest
// (or a test-ns-hook) in through the public clojure.test API, in one
// run-tests call. It compiles no source strings and holds no test state.
// Returns (success, message-for-failure).
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
		initialNS, _ := rt.CurrentNS.Deref().(*vm.Namespace)
		before := testSnapshot()
		if loadErr := runFile(path, baseline); loadErr != nil {
			msg = "load: " + loadErr.Error()
			return vm.NIL, nil
		}
		kind, namespaces := classifyLoaded(initialNS, before)
		switch kind {
		case kindStrayDeftest:
			msg = "deftest outside a namespace: " + path
			return vm.NIL, nil
		case kindLoadOnly:
			// Nothing gained a test; any assertions ran at load time.
			ok = true
			return vm.NIL, nil
		}
		args := make([]vm.Value, 0, len(namespaces))
		for _, ns := range namespaces {
			args = append(args, ns)
		}
		summary, err := rt.InvokeValue(rt.LookupVar("test", "run-tests").Deref(), args)
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
			// fixtures/ holds inputs read by tests via slurp (and, for the
			// quality corpus, files whose compilation order would matter);
			// they are data, not deftests.
			if info.Name() == "compat" || info.Name() == "clojure-test-suite" || info.Name() == "benches" || info.Name() == "gogen" || info.Name() == "gold-aot" || info.Name() == "native-entry" || info.Name() == "tools" || info.Name() == "fixtures" {
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
