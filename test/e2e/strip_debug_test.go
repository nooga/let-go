/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripDebugCompanion(t *testing.T) {
	lg := buildLG(t)
	dir := t.TempDir()
	source := filepath.Join(dir, "app.lg")
	lgb := filepath.Join(dir, "app.lgb")
	program := `(defn explode [] (+ 1 "not-a-number"))
(when-not *compiling-aot* (explode))`
	if err := os.WriteFile(source, []byte(program), 0644); err != nil {
		t.Fatal(err)
	}

	if out, err := exec.Command(lg, "-strip", "-c", lgb, source).CombinedOutput(); err != nil {
		t.Fatalf("compile stripped LGB: %v\n%s", err, out)
	}
	debugPath := lgb + ".debug"
	if info, err := os.Stat(debugPath); err != nil || info.Size() == 0 {
		t.Fatalf("debug companion was not created: info=%v err=%v", info, err)
	}

	run := func(env ...string) (int, string) {
		t.Helper()
		cmd := exec.Command(lg, lgb)
		cmd.Env = append(os.Environ(), env...)
		out, err := cmd.CombinedOutput()
		if err == nil {
			return 0, string(out)
		}
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("run stripped LGB: %v\n%s", err, out)
		}
		return exitErr.ExitCode(), string(out)
	}

	code, out := run()
	if code == 0 || !strings.Contains(out, "app.lg:") {
		t.Fatalf("companion did not restore source-located traceback (exit %d):\n%s", code, out)
	}

	code, out = run("LG_DEBUG_FILE=")
	if code == 0 || !strings.Contains(out, "<unknown>") || strings.Contains(out, "app.lg:") {
		t.Fatalf("disabling companion did not produce an unlocated traceback (exit %d):\n%s", code, out)
	}

	otherSource := filepath.Join(dir, "other.lg")
	otherLGB := filepath.Join(dir, "other.lgb")
	if err := os.WriteFile(otherSource, []byte(strings.Replace(program, "not-a-number", "different", 1)), 0644); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command(lg, "-strip", "-c", otherLGB, otherSource).CombinedOutput(); err != nil {
		t.Fatalf("compile second stripped LGB: %v\n%s", err, out)
	}
	code, out = run("LG_DEBUG_FILE=" + otherLGB + ".debug")
	if code == 0 || !strings.Contains(out, "SHA-256 mismatch") {
		t.Fatalf("mismatched companion was not rejected (exit %d):\n%s", code, out)
	}

	customDebug := filepath.Join(dir, "symbols", "app.debug")
	if err := os.Mkdir(filepath.Dir(customDebug), 0755); err != nil {
		t.Fatal(err)
	}
	customLGB := filepath.Join(dir, "custom.lgb")
	if out, err := exec.Command(lg, "-strip", "-debug-output", customDebug, "-c", customLGB, source).CombinedOutput(); err != nil {
		t.Fatalf("compile with custom debug output: %v\n%s", err, out)
	}
	if info, err := os.Stat(customDebug); err != nil || info.Size() == 0 {
		t.Fatalf("custom debug companion was not created: info=%v err=%v", info, err)
	}
	if _, err := os.Stat(customLGB + ".debug"); !os.IsNotExist(err) {
		t.Fatalf("unexpected default companion beside custom output: %v", err)
	}

	app := filepath.Join(dir, "app")
	if out, err := exec.Command(lg, "-strip", "-b", app, source).CombinedOutput(); err != nil {
		t.Fatalf("build stripped standalone bundle: %v\n%s", err, out)
	}
	if info, err := os.Stat(app + ".debug"); err != nil || info.Size() == 0 {
		t.Fatalf("bundle debug companion was not created: info=%v err=%v", info, err)
	}
	cmd := exec.Command(app)
	bundleOut, bundleErr := cmd.CombinedOutput()
	var bundleExit *exec.ExitError
	if !errors.As(bundleErr, &bundleExit) || !strings.Contains(string(bundleOut), "app.lg:") {
		t.Fatalf("standalone bundle did not load its companion: err=%v\n%s", bundleErr, bundleOut)
	}
}
