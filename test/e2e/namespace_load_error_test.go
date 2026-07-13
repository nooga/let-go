/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestREPLFormatsNamespaceLoadCompileChainAndContinues(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "broken.lg")
	source := "(ns broken)\n(def broken-value\n  (fn []\n    (let [:tag 1] 1)))\n"
	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}

	bin := buildLG(t)
	cmd := exec.Command(bin, "-source-paths", dir)
	cmd.Stdin = strings.NewReader("(require 'broken)\n(+ 1 2)\n")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("REPL exited after namespace-load error: %v\n%s", err, output.String())
	}

	got := output.String()
	if strings.Count(got, "error:") != 1 {
		t.Fatalf("namespace-load error rendered more than once:\n%s", got)
	}
	for _, want := range []string{
		file + ":2:",
		file + ":3:",
		file + ":4:",
		"compiling def value",
		"compiling fn body",
		"let binding name must be a symbol",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("formatted namespace-load error missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "unable to load namespace") {
		t.Fatalf("REPL replaced the original error:\n%s", got)
	}
	if !strings.Contains(got, "\n3\n") {
		t.Fatalf("REPL did not evaluate the next form:\n%s", got)
	}
}
