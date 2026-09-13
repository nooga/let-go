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

// TestBoundFnKeepsInvokingScope: bound-fn* conveys dynamic bindings but must
// not detach the call from the invoking execution's scope. A sleep inside a
// bound fn called from a scoped future is interrupted when that scope closes
// (sleep returns early on cancellation), and the captured binding is still
// visible inside the call. A bound fn detached from the scope sleeps its full
// ten seconds, so the deref times out instead.
func TestBoundFnKeepsInvokingScope(t *testing.T) {
	bin := buildLG(t)
	src := `(def ^:dynamic *x* :root) ` +
		`(def f (binding [*x* :caller] (bound-fn* (fn [] (sleep 10000) [*x* :woke])))) ` +
		`(def parent (scope-open)) ` +
		`(def seen (promise)) ` +
		`(future (deliver seen (f))) ` +
		`(sleep 50) ` +
		`(scope-close! parent 5000) ` +
		`(println (deref seen 1000 :timeout))`
	start := time.Now()
	out, err := exec.Command(bin, "-e", src).CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "[:caller :woke]") {
		t.Fatalf("expected [:caller :woke], got: %q", out)
	}
	if d := time.Since(start); d > 9*time.Second {
		t.Fatalf("bound fn escaped its invoking scope: %v", d)
	}
}
