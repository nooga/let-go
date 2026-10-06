/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestIRStressLowersEachDefnInItsInNsNamespace: the ir-stress harness lowers a
// defn in the namespace a top-level (in-ns 'x) switched to, as loading the file
// does. It used to lower every defn in the file's single (ns ...) namespace, so
// a defn written after (in-ns 'json) could not resolve json's own read-json and
// was counted as a lowering failure that real lowering does not have.
func TestIRStressLowersEachDefnInItsInNsNamespace(t *testing.T) {
	bin := buildLG(t)
	root := repoRoot(t)

	runStress := func(t *testing.T, dir, file string) string {
		t.Helper()
		cmd := exec.Command(bin, "scripts/ir-stress.lg", "lower-go", dir, file)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "LG_STRESS_PASSES=1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("ir-stress %s: %v\n%s", file, err, out)
		}
		return string(out)
	}

	t.Run("in-ns into a Go-installed and a fresh namespace", func(t *testing.T) {
		dir := t.TempDir()
		src := `(in-ns 'json)
(defn stress-probe-load [s] (read-json s))
(in-ns 'stress.probe.fresh)
(defn helper [] 1)
(defn uses-helper [] (helper))
`
		if err := os.WriteFile(filepath.Join(dir, "probe.lg"), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
		out := runStress(t, dir, "probe.lg")
		if !strings.Contains(out, "Total fixtures: 3") || !strings.Contains(out, "Failed: 0") {
			t.Fatalf("expected all 3 fixtures to lower, got:\n%s", out)
		}
	})

	t.Run("the corpus file that switches namespaces", func(t *testing.T) {
		// Two defns plus the two deftest bodies of the file's last segment,
		// whose (clojure.core/require '[test :refer :all]) must run in that
		// segment's namespace for deftest to expand at all.
		out := runStress(t, "./", "test/in_ns_auto_refer_test.lg")
		if !strings.Contains(out, "Total fixtures: 4") || !strings.Contains(out, "Failed: 0") {
			t.Fatalf("expected all 4 fixtures to lower, got:\n%s", out)
		}
	})
}

// TestIRStressGateAcceptsZeroFailureBaseline: once every corpus form lowers, the
// baseline is {:failed 0 :buckets {}}. The gate must accept it, pass a clean
// run, and report any failure as a regression; it used to reject an empty
// :buckets map as a baseline with nothing to ratchet against.
func TestIRStressGateAcceptsZeroFailureBaseline(t *testing.T) {
	bin := buildLG(t)
	root := repoRoot(t)
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.edn")
	if err := os.WriteFile(baseline, []byte("{:failed 0 :total 1 :passed 1 :buckets {}}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gate := func(t *testing.T, src string) (string, int) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, "probe.lg"), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(bin, "scripts/ir-stress.lg", "lower-go", dir, "probe.lg")
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "LG_STRESS_PASSES=1", "LG_STRESS_BASELINE="+baseline)
		out, err := cmd.CombinedOutput()
		code := 0
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else if err != nil {
			t.Fatalf("ir-stress: %v\n%s", err, out)
		}
		return string(out), code
	}

	t.Run("a clean run passes", func(t *testing.T) {
		out, code := gate(t, "(ns stress.zero.clean)\n(defn ok [] 1)\n")
		if code != 0 {
			t.Fatalf("expected exit 0, got %d:\n%s", code, out)
		}
	})
	t.Run("a new failure is a regression, not a gate error", func(t *testing.T) {
		out, code := gate(t, "(ns stress.zero.dirty)\n(defn bad [] (no-such-fn))\n")
		if code != 1 || !strings.Contains(out, "REGRESSED") {
			t.Fatalf("expected exit 1 with REGRESSED, got %d:\n%s", code, out)
		}
	})
}
