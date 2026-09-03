/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// emitFailLine is the line external orchestrators scrape. gloat invokes
// scripts/lg-compile by path and greps this out of stdout to decide which
// functions stayed on the trampoline, so its shape is a compatibility
// contract, not a log message.
var emitFailLine = regexp.MustCompile(`^EMIT-FAIL (\S+) pkg=(\S+) returned (.+)$`)

// TestLgCompileEmitFailContract pins the two halves of the #596 stdout/exit
// contract that external callers depend on:
//
//   - the "EMIT-FAIL <path> pkg=<pkg> returned <type>" line, verbatim;
//   - EMIT-FAIL is NOT an exit failure. A fn that does not lower stays on the
//     trampoline and the program still runs, so partial lowering is a normal
//     outcome the orchestrator judges for itself. Only an entry-arity error
//     exits non-zero. legmacs' `set -e` build depends on this.
//
// No .lg fixture drives this: as of 2026-09-02 no construct probed (try/catch,
// atom, reify, deftype, dynamic set!, lazy-seq, macro-expanded body, multimethod
// dispatch) fails to lower any more, and a fixture picked for today's gaps would
// silently stop testing the path the moment lowering covered it. Both halves
// therefore inject a synthetic non-lowered result instead — the stdout half into
// write-packages!, where the line is printed, and the exit-status half into the
// real shim, by stubbing lg.compiler/main-from-args and evaluating
// scripts/lg-compile against the stub.
//
// The exit-status half MUST run the shim rather than assert on a helper. An
// earlier version of this test checked write-packages! and then ran the shim on
// a SUCCESSFUL compile, which never exercised a nonzero :emit-fails through the
// shim at all: reintroducing "exit 1 when :emit-fails is nonzero" left it
// passing. This version fails against that regression, which is the only reason
// to trust it.
func TestLgCompileEmitFailContract(t *testing.T) {
	bin := buildLG(t)
	root := repoRoot(t)

	runLG := func(t *testing.T, args ...string) (string, int) {
		t.Helper()
		cmd := exec.Command(bin, args...)
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		code := 0
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else if err != nil {
			t.Fatalf("running lg %v: %v\n%s", args, err, out)
		}
		return string(out), code
	}

	t.Run("EMIT-FAIL line shape is verbatim", func(t *testing.T) {
		dir := t.TempDir()
		probe := filepath.Join(dir, "probe.lg")
		// A result whose :source is not a string is exactly what the pipeline
		// hands back for a package it could not lower.
		src := `(require 'lg.compiler)
(let [specs   [{:path "/x/broken.lg" :pkg "broken" :reldir "broken" :lf ['(defn f [] 1)]}]
      results [{:fns [] :source nil}]
      fails   (lg.compiler/write-packages! "` + dir + `/out" specs results)]
  (println (str "fails=" fails)))
`
		if err := os.WriteFile(probe, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
		out, code := runLG(t, probe)
		if code != 0 {
			t.Fatalf("probe exited %d; output:\n%s", code, out)
		}

		var got string
		for _, line := range strings.Split(out, "\n") {
			if strings.HasPrefix(line, "EMIT-FAIL ") {
				got = line
				break
			}
		}
		if got == "" {
			t.Fatalf("no EMIT-FAIL line in output:\n%s", out)
		}
		m := emitFailLine.FindStringSubmatch(got)
		if m == nil {
			t.Fatalf("EMIT-FAIL line %q does not match the scraped shape %q", got, emitFailLine)
		}
		if m[1] != "/x/broken.lg" || m[2] != "broken" {
			t.Errorf("EMIT-FAIL path/pkg = %q/%q, want /x/broken.lg/broken", m[1], m[2])
		}
		// The count is returned to the caller, which is what lets a Go caller
		// decide differently from the shim.
		if !strings.Contains(out, "fails=1") {
			t.Errorf("write-packages! should return the EMIT-FAIL count; output:\n%s", out)
		}
	})

	// The regression guard: the real shim, driven with a result carrying
	// EMIT-FAILs and nothing else, must still exit 0 and report success.
	t.Run("EMIT-FAIL through the real shim exits 0", func(t *testing.T) {
		dir := t.TempDir()
		probe := filepath.Join(dir, "shim-inject.lg")
		src := `(require 'lg.compiler)
(in-ns 'lg.compiler)
(def main-from-args (fn [argv] {:emit-fails 3 :lowered-funcs 0 :entry nil}))
(in-ns 'user)
(eval (read-string (str "(do " (slurp "scripts/lg-compile") "\n)")))
`
		if err := os.WriteFile(probe, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
		out, code := runLG(t, probe)
		if code != 0 {
			t.Fatalf("shim exited %d with 3 EMIT-FAILs and no :error; want 0. output:\n%s", code, out)
		}
		if !strings.Contains(out, "done.") {
			t.Errorf("shim should report success through EMIT-FAILs; output:\n%s", out)
		}
	})

	// And the fatal case, so "exits 0" above is a real signal rather than a
	// shim that cannot fail.
	t.Run("an :error result exits non-zero through the real shim", func(t *testing.T) {
		dir := t.TempDir()
		probe := filepath.Join(dir, "shim-error.lg")
		src := `(require 'lg.compiler)
(in-ns 'lg.compiler)
(def main-from-args (fn [argv] {:error "synthetic failure" :emit-fails 0}))
(in-ns 'user)
(eval (read-string (str "(do " (slurp "scripts/lg-compile") "\n)")))
`
		if err := os.WriteFile(probe, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
		out, code := runLG(t, probe)
		if code == 0 {
			t.Fatalf("shim exited 0 on an :error result; want non-zero. output:\n%s", out)
		}
		if !strings.Contains(out, "synthetic failure") {
			t.Errorf("want the error text on stdout; got:\n%s", out)
		}
	})

	t.Run("end to end: a real compile exits 0, a bad entry arity exits 1", func(t *testing.T) {
		dir := t.TempDir()
		write := func(name, body string) string {
			p := filepath.Join(dir, name)
			if err := os.WriteFile(p, []byte(body), 0644); err != nil {
				t.Fatal(err)
			}
			return p
		}

		ok := write("ok.lg", "(ns app)\n(defn -main [] 0)\n")
		out, code := runLG(t, "scripts/lg-compile", filepath.Join(dir, "out-ok"), "tmpmod", ok)
		if code != 0 {
			t.Fatalf("a lowering run must exit 0; got %d:\n%s", code, out)
		}

		// The one case that does exit non-zero, so "exit 0" above is a real
		// signal rather than a shim that can never fail.
		bad := write("bad.lg", "(ns app)\n(defn -main [a b] (+ a b))\n")
		out, code = runLG(t, "scripts/lg-compile", "--entry-frame", filepath.Join(dir, "out-bad"), "tmpmod", bad)
		if code == 0 {
			t.Fatalf("an unsupported entry arity must exit non-zero; output:\n%s", out)
		}
		if !strings.Contains(out, "unsupported arity") {
			t.Errorf("want the entry-arity error on stdout; got:\n%s", out)
		}
	})
}
