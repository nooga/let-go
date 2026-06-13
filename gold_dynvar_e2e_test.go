/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Gold differential tests: the SAME .cljc script run under let-go must produce
// byte-for-byte the same last-output line as Clojure. The verified Clojure 1.12
// / babashka reference output lives in a committed test/gold/<name>.out file.
//
// By DEFAULT these tests are clojure-free: each runs its script under let-go and
// compares the result to the committed .out file. No JVM, no network, no
// `clojure` dependency — fast and deterministic in CI.
//
// Re-deriving the .out goldens from a live `clojure` is opt-in via
// LETGO_GOLD_REDERIVE=1 (clojure must be on PATH). In that mode each test runs
// the script under real Clojure, OVERWRITES its .out file, and then runs the
// usual let-go assertion against the refreshed gold — so one
// `LETGO_GOLD_REDERIVE=1 go test` both refreshes the goldens and verifies let-go
// still matches. Trigger it by hand, or via the CI `gold-differential` job that
// fires on workflow_dispatch or when the resolved clojure version changes.
//
// Scenario coverage of each script:
//   - dynvar_threads.cljc:   dynamic-var (namespace variable) binding
//     propagation across nested threads — future conveyance, isolation from
//     later rebinds, nested-future inheritance, concurrent non-interference,
//     bound-fn capture, thread-local set! mutation.
//   - ns_threads.cljc:       *ns* conveying into futures, isolating across
//     concurrent futures, inheriting through nested futures, in-ns mutation.
//   - dynvar_callbacks.cljc: dynamic bindings conveying through EAGER
//     callback-invoking builtins (mapv, filterv, reduce, run!, sort-by, swap!,
//     transducer application, with-out-str) when the binding lives in a child
//     ExecContext (inside a future). A context-free Fn.Invoke at any of those
//     sites resolves the var against the root context and breaks every scenario.
//
// The scripts deliberately stay inside the subset where let-go and Clojure
// agree. let-go is intentionally MORE PERMISSIVE than Clojure in two spots the
// scripts avoid: (a) set! with no active binding (Clojure throws; let-go mutates
// the root — load-bearing for test.lg), and (b) set! on a conveyed binding in a
// future that opened no binding of its own (Clojure throws; let-go permits it).

const (
	dynvarThreadsScript   = "test/gold/dynvar_threads.cljc"
	nsThreadsScript       = "test/gold/ns_threads.cljc"
	dynvarCallbacksScript = "test/gold/dynvar_callbacks.cljc"
)

// rederiveGoldEnv, when set to "1", makes the gold tests re-derive their .out
// files from a live `clojure` before asserting let-go against them.
const rederiveGoldEnv = "LETGO_GOLD_REDERIVE"

func TestGoldDynvarThreadsMatchesClojure(t *testing.T)   { checkGold(t, dynvarThreadsScript) }
func TestGoldNsThreadsMatchesClojure(t *testing.T)       { checkGold(t, nsThreadsScript) }
func TestGoldDynvarCallbacksMatchesClojure(t *testing.T) { checkGold(t, dynvarCallbacksScript) }

// goldPath maps a .cljc script to its committed expected-output file:
// test/gold/foo.cljc -> test/gold/foo.out.
func goldPath(script string) string {
	return strings.TrimSuffix(script, ".cljc") + ".out"
}

// checkGold runs script under let-go and asserts its last non-empty output line
// equals the committed gold file. With LETGO_GOLD_REDERIVE=1 it first re-derives
// that gold file from real Clojure (which must then be on PATH).
func checkGold(t *testing.T, script string) {
	t.Helper()
	gold := goldPath(script)

	if os.Getenv(rederiveGoldEnv) == "1" {
		clj, err := exec.LookPath("clojure")
		if err != nil {
			t.Fatalf("%s=1 but no `clojure` on PATH: %v", rederiveGoldEnv, err)
		}
		ref := lastNonEmptyLine(stdoutOf(t, clj, "-M", script))
		if ref == "" {
			t.Fatalf("clojure produced no output for %s", script)
		}
		if err := os.WriteFile(gold, []byte(ref+"\n"), 0o644); err != nil {
			t.Fatalf("write gold %s: %v", gold, err)
		}
		t.Logf("re-derived %s from clojure: %q", gold, ref)
	}

	want, err := os.ReadFile(gold)
	if err != nil {
		t.Fatalf("read gold %s: %v\n(run `%s=1 go test` with clojure on PATH to (re)generate it)", gold, err, rederiveGoldEnv)
	}
	wantLine := strings.TrimSpace(string(want))

	lg := buildLG(t)
	got := lastNonEmptyLine(stdoutOf(t, lg, script))
	if got != wantLine {
		t.Fatalf("let-go output = %q, want %q (from %s)", got, wantLine, gold)
	}
}

// stdoutOf runs cmd and returns its stdout, failing the test on a non-zero exit.
// stderr is captured only to surface in failure messages.
func stdoutOf(t *testing.T, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s %v failed: %v\nstderr:\n%s", name, args, err, errb.String())
	}
	return out.String()
}

func lastNonEmptyLine(s string) string {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if t := strings.TrimSpace(lines[i]); t != "" {
			return t
		}
	}
	return ""
}
