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

// TestWithScopeTeardownReleasesParkedTake: a go-block parked on (<! ch) inside
// with-scope is released when the block exits (scope cancel), runs to its end,
// and the program returns promptly — proving the take was cancelled, not fed.
//
// Post-#920, a cancelled <! raises the Cancelled condition instead of quietly
// returning nil (that silent-nil behavior was the bug #920 fixed — it made a
// cancelled take indistinguishable from one that legitimately read nil). The
// (catch Cancelled e nil) below is what makes this test's proof explicit
// rather than incidental: reaching (reset! d :ran) now demonstrates the take
// was cancelled BECAUSE the catch fired, not merely because <! returned
// something falsy.
func TestWithScopeTeardownReleasesParkedTake(t *testing.T) {
	bin := buildLG(t)
	src := `(def ch (chan)) ` +
		`(def d (atom :parked)) ` +
		`(with-scope [s] (go (try (<! ch) (catch Cancelled e nil)) (reset! d :ran))) ` +
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
