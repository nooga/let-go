/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestScopeCancelledPredicate: scope-cancelled? is false for a live scope,
// true for an explicit scope after it is closed, and true from inside a
// future whose parent scope was closed — sleep returns early and silently on
// cancellation, so this predicate is the only way Lisp code can observe it.
func TestScopeCancelledPredicate(t *testing.T) {
	bin := buildLG(t)
	src := `(def parent (scope-open)) ` +
		`(def before (scope-cancelled? parent)) ` +
		`(def seen (promise)) ` +
		`(def child-before (promise)) ` +
		`(future (deliver child-before (scope-cancelled?)) (sleep 10000) (deliver seen (scope-cancelled?))) ` +
		`(deref child-before) ` +
		`(scope-close! parent 5000) ` +
		`(println [before @child-before (scope-cancelled? parent) (deref seen 1000 :timeout)])`
	start := time.Now()
	out, err := exec.Command(bin, "-e", src).CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "[false false true true]") {
		t.Fatalf("expected [false false true true], got: %q", out)
	}
	if d := time.Since(start); d > 9*time.Second {
		t.Fatalf("sleep was not interrupted by teardown: %v", d)
	}
}
