/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
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

// TestCompiledProgramRequiresEmbeddedNamespace guards #954: a compiled program
// that requires a core namespace loaded from embedded source, rather than from
// the core bundle at boot, must get that namespace's definitions. zip is a
// standalone embedded namespace; ir.build is part of the ir.* pipeline #356
// moved out of the bundle.
func TestCompiledProgramRequiresEmbeddedNamespace(t *testing.T) {
	lg := buildLG(t)
	dir := t.TempDir()
	source := filepath.Join(dir, "app.lg")
	program := `(ns app (:require [zip :as z] [ir.build :as b]))
(println "zip:" (z/node (z/vector-zip [1 2])))
(println "ir.build:" (fn? b/build-fn))`
	if err := os.WriteFile(source, []byte(program), 0o644); err != nil {
		t.Fatal(err)
	}
	want := []string{"zip: [1 2]", "ir.build: true"}

	check := func(t *testing.T, cmd *exec.Cmd) {
		t.Helper()
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("run: %v\n%s", err, out)
		}
		for _, line := range want {
			if !strings.Contains(string(out), line) {
				t.Fatalf("output lacks %q:\n%s", line, out)
			}
		}
	}

	t.Run("lgb", func(t *testing.T) {
		lgb := filepath.Join(dir, "app.lgb")
		if out, err := exec.Command(lg, "-c", lgb, source).CombinedOutput(); err != nil {
			t.Fatalf("compile: %v\n%s", err, out)
		}
		check(t, exec.Command(lg, lgb))
	})
	t.Run("standalone", func(t *testing.T) {
		bin := filepath.Join(dir, "app")
		if out, err := exec.Command(lg, "-b", bin, source).CombinedOutput(); err != nil {
			t.Fatalf("bundle: %v\n%s", err, out)
		}
		check(t, exec.Command(bin))
	})
}
