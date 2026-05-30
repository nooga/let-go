/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// buildLG builds the lg binary once into a temp dir and returns its path.
func buildLG(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "lg")
	out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("build lg: %v\n%s", err, out)
	}
	return bin
}

// TestWithScopeTeardownReleasesParkedTake: a go-block parked on (<! ch) inside
// with-scope is released when the block exits (scope cancel), runs to its end,
// and the program returns promptly — proving the take was cancelled, not fed.
func TestWithScopeTeardownReleasesParkedTake(t *testing.T) {
	bin := buildLG(t)
	src := `(def ch (chan)) ` +
		`(def d (atom :parked)) ` +
		`(with-scope [s] (go (<! ch) (reset! d :ran))) ` +
		`(println @d)`
	start := time.Now()
	out, err := exec.Command(bin, "-e", src).CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "ran") {
		t.Fatalf("expected go-block to run after cancel, got: %q", out)
	}
	if d := time.Since(start); d > 9*time.Second {
		t.Fatalf("with-scope teardown took %v — take was not cancelled", d)
	}
}

// TestWithScopeInterruptsSleep: a (sleep 10000) inside a scoped go-block is
// interrupted by teardown, so the whole run finishes well under 10s (bounded
// by *scope-drain-timeout-ms*, default 5000).
func TestWithScopeInterruptsSleep(t *testing.T) {
	bin := buildLG(t)
	src := `(with-scope [s] (go (sleep 10000))) (println :done)`
	start := time.Now()
	out, err := exec.Command(bin, "-e", src).CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if d := time.Since(start); d > 9*time.Second {
		t.Fatalf("sleep was not interrupted by teardown: %v\n%s", d, out)
	}
	if !strings.Contains(string(out), "done") {
		t.Fatalf("expected :done, got: %q", out)
	}
}
