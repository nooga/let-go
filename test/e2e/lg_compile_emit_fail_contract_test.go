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
// The EMIT-FAIL half drives lg.compiler/write-packages! directly with a
// synthetic non-lowered result rather than a .lg fixture chosen for not
// lowering. That is deliberate: as of 2026-09-02 no construct probed (try/catch,
// atom, reify, deftype, dynamic set!, lazy-seq, macro-expanded body, multimethod
// dispatch) fails to lower any more, and a fixture picked for today's gaps would
// silently stop testing the path the moment lowering covered it. write-packages!
// is where the line is printed and where the count that does NOT reach the exit
// code is returned, so pinning it there tests the contract rather than the gap.
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

	t.Run("exit status tracks :error, not :emit-fails", func(t *testing.T) {
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
