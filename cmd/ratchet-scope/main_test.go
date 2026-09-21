package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAffectedDecisions pins the decision table with the paths that actually
// motivated this tool. The two halves matter for different reasons: the RUN
// cases are the gate doing its job, and the SKIP cases are the pushes that
// used to be blocked by a stale bar they had nothing to do with.
func TestAffectedDecisions(t *testing.T) {
	// A stand-in file set rather than the real one, so these cases test the
	// decision logic and not the dependency graph. TestRealClosure covers the
	// graph itself. Values are the owning package, as `go list` reports it.
	deps := map[string]string{
		"pkg/vm/var.go":                  "pkg/vm",
		"pkg/vm/testdata/fixture.edn":    "pkg/vm",
		"pkg/ir/build.lg":                "pkg/ir",
		"pkg/rt/lang.go":                 "pkg/rt",
		"pkg/rt/core/core.lg":            "pkg/rt",
		"pkg/rt/core_compiled.lgb":       "pkg/rt",
		"pkg/rt/core_go_lowered/core.go": "pkg/rt/core_go_lowered",
	}

	for _, tc := range []struct {
		path string
		run  bool
		why  string
	}{
		// Measured code, directly.
		{"pkg/vm/var.go", true, "the VM itself"},
		{"pkg/ir/build.lg", true, "the IR itself"},
		// Measured code, transitively — the case the path-regex approach
		// could not express.
		{"pkg/rt/lang.go", true, "runtime, imported by the benchmarks"},
		// Not a .go file, but in a measured package: embedded or read data
		// changes behaviour just as source does.
		{"pkg/vm/testdata/fixture.edn", true, "an embedded asset go list reports"},
		// The case that motivated file-list matching: it sits in the `test`
		// package's directory but belongs to no build, and the suite benchmark
		// reads only the vendored clojure-test-suite corpus.
		{"test/quality_inputs_test.lg", false, "in a root's directory but in no build"},
		// Generated artifacts and their sources: loaded at runtime, imported
		// by nothing.
		{"pkg/rt/core/core.lg", true, "embedded by pkg/rt"},
		{"pkg/rt/core_compiled.lgb", true, "embedded by pkg/rt"},
		{"pkg/rt/core_go_lowered/core.go", true, "compiled under -tags gogen_ir"},
		{"pkg/rt/generated.sums", true, "generated-artifact manifest"},
		// Build inputs that change generated code beneath every package.
		{"go.mod", true, "module graph"},
		{"mise.toml", true, "Go toolchain version"},
		// Read by the suite benchmark at run time, so no package lists them.
		{"test/compat/clojure/core-test/portability.lg", true, "the suite portability shim"},
		{"test/compat/clojure/core_test/abs.cljc", true, "a namespace the suite loader can resolve"},
		{"test/clojure-test-suite", true, "the suite submodule gitlink"},
		{"test/clojure-test-suite/test/clojure/core_test/abs.cljc", true, "a file inside the suite corpus"},
		// Out of scope: nothing the measured binaries compile or load.
		{"docs/perf/ratchet.md", false, "documentation"},
		{"docs/perf/baseline.json", false, "the gate's reference, not measured code"},
		{"cmd/bench-ratchet/main.go", false, "the gate tool itself"},
		{"cmd/ratchet-scope/main.go", false, "this tool"},
		{".github/workflows/ci.yml", false, "CI config"},
		{"README.md", false, "documentation"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			reason := affected([]string{tc.path}, deps, present)
			if got := reason != ""; got != tc.run {
				t.Errorf("affected(%q) RUN = %v, want %v (%s); reason = %q",
					tc.path, got, tc.run, tc.why, reason)
			}
			if tc.run && reason == "" {
				t.Error("a RUN decision must carry a reason to print")
			}
		})
	}
}

// present stands in for a checkout in which every changed path still exists,
// so the decision-table tests exercise matching without the deletion rule.
func present(string) bool { return true }

// TestDeletedPathsUnderTheClosureRun pins the deletion rule. The closure is
// computed from the tree after the push, so a deleted file is never in it; the
// review repro was deleting the ratchet's own anchor benchmark, which read as
// SKIP.
func TestDeletedPathsUnderTheClosureRun(t *testing.T) {
	deps := map[string]string{
		"pkg/vm/vm.go":        "pkg/vm",
		"pkg/rt/lang.go":      "pkg/rt",
		"pkg/rt/core/core.lg": "pkg/rt",
	}
	deleted := func(string) bool { return false }
	for _, tc := range []struct {
		path string
		run  bool
		why  string
	}{
		{"pkg/vm/bench_ratchet_anchor_test.go", true, "the anchor benchmark, in a measured package"},
		{"pkg/rt/core/sub/gone.lg", true, "an embedded subtree removed whole"},
		{"docs/perf/old.md", false, "documentation"},
		{"cmd/bench-ratchet/gone.go", false, "the gate tool, outside the closure"},
		{"stray.txt", false, "a file at the repository root"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			reason := affected([]string{tc.path}, deps, deleted)
			if got := reason != ""; got != tc.run {
				t.Errorf("affected(deleted %q) RUN = %v, want %v (%s); reason = %q",
					tc.path, got, tc.run, tc.why, reason)
			}
		})
	}
	// The same paths, still present, fall back to exact matching: the rule
	// applies to deletions only, so an out-of-closure file beside measured
	// code does not start running the ratchet.
	if reason := affected([]string{"pkg/vm/notes.txt"}, deps, present); reason != "" {
		t.Errorf("present pkg/vm/notes.txt should not be affected; reason = %q", reason)
	}
}

// TestAffectedReportsFirstMatchInMixedSet pins that one in-scope path carries
// the whole change set, which is the conservative direction: a push that
// touches the VM and fifty docs files still runs.
func TestAffectedReportsFirstMatchInMixedSet(t *testing.T) {
	deps := map[string]string{"pkg/vm/vm.go": "pkg/vm"}
	paths := []string{"README.md", "docs/a.md", "pkg/vm/vm.go", "docs/b.md"}
	reason := affected(paths, deps, present)
	if reason == "" {
		t.Fatal("a change set containing pkg/vm/vm.go must be affected")
	}
	if !strings.Contains(reason, "pkg/vm/vm.go") {
		t.Errorf("reason = %q, want it to name the in-scope path", reason)
	}
}

// TestNormalizeHandlesPathSpellings covers the spellings a VCS can emit for
// the same file. A missed normalisation reads as "not affected", which is the
// unsafe direction.
func TestNormalizeHandlesPathSpellings(t *testing.T) {
	deps := map[string]string{"pkg/vm/vm.go": "pkg/vm"}
	for _, p := range []string{"pkg/vm/vm.go", "./pkg/vm/vm.go", " pkg/vm/vm.go ", "pkg/vm/vm.go\r"} {
		if affected([]string{p}, deps, present) == "" {
			t.Errorf("affected(%q) = not affected, want affected", p)
		}
	}
}

// TestExactFileMatchHasNoPrefixBleed pins that matching is by exact file, not
// by prefix. The old tree-matching rule had to defend against "testdata/" being
// swallowed by the "test" package; an exact lookup cannot make that mistake,
// and this test is what proves the property is now structural.
func TestExactFileMatchHasNoPrefixBleed(t *testing.T) {
	deps := map[string]string{"test/zz_bench_test.go": "test", "pkg/vm/vm.go": "pkg/vm"}
	for path, want := range map[string]bool{
		"test/zz_bench_test.go": true,
		"pkg/vm/vm.go":          true,
		"test/x.lg":             false,
		"testdata/x.edn":        false,
		"test-helpers/a.go":     false,
		"pkg/vmx/a.go":          false,
		"pkg/vm/vm.go.bak":      false,
	} {
		if got := matchedDepFile(path, deps) != ""; got != want {
			t.Errorf("matchedDepFile(%q) matched = %v, want %v", path, got, want)
		}
	}
}

// TestReadPathsIgnoresBlankLines guards the stdin contract: a trailing newline
// or a blank line must not read as a path.
func TestReadPathsIgnoresBlankLines(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "paths")
	if err := os.WriteFile(f, []byte("a.go\n\n  \nb.go\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := readPaths(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "a.go" || got[1] != "b.go" {
		t.Errorf("readPaths = %q, want [a.go b.go]", got)
	}
}

// TestRealClosureCoversTheRuntime runs the actual dependency expansion. It is
// the guard against the roots list going stale: if a benchmark moves to a
// package not listed in benchmarkRoots, or the runtime stops being reachable
// from the roots, the closure silently narrows and the gate stops firing on
// changes it should catch. Asserting a known-transitive package is in, and a
// known-unrelated one is out, catches both directions.
func TestRealClosureCoversTheRuntime(t *testing.T) {
	if testing.Short() {
		t.Skip("runs go list over the module")
	}
	files, err := dependencyFiles()
	if err != nil {
		t.Fatalf("dependencyFiles: %v", err)
	}
	// Real files that must be in scope. pkg/rt/lang.go is the transitive case
	// (not a benchmark root); the two pkg/rt artifacts are the embed case,
	// which is what lets the trigger list stay short.
	for _, want := range []string{
		"pkg/vm/vm.go",
		"pkg/rt/lang.go",
		"pkg/rt/core/core.lg",
		"pkg/rt/core_compiled.lgb",
	} {
		if files[want] == "" {
			t.Errorf("file set is missing %q — benchmarkRoots or the field list may be stale", want)
		}
	}
	// The gate tooling is not compiled into the benchmarks. If it ever were,
	// editing the ratchet would again require a benchmark run — the behaviour
	// this tool exists to remove.
	for _, notWant := range []string{
		"cmd/bench-ratchet/main.go",
		"cmd/ratchet-scope/main.go",
		"docs/perf/baseline.json",
	} {
		if pkg := files[notWant]; pkg != "" {
			t.Errorf("file set unexpectedly contains %q (owned by %s)", notWant, pkg)
		}
	}
	// The regression that motivated file-list matching: a .lg test lives in the
	// `test` package's directory but is part of no build, and the suite
	// benchmark reads only the vendored corpus.
	for _, notWant := range []string{
		"test/quality_inputs_test.lg",
		"test/seq_test.lg",
	} {
		if pkg := files[notWant]; pkg != "" {
			t.Errorf("%q should not be in any build (reported under %s); tree matching has crept back", notWant, pkg)
		}
	}
}
