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

// TestBoundFnKeepsInvokingScope: bound-fn* conveys dynamic bindings but must
// not detach the call from the invoking execution's scope. A sleep inside a
// bound fn called from a scoped future is interrupted when that scope closes,
// while the captured binding is still visible inside the call.
func TestBoundFnKeepsInvokingScope(t *testing.T) {
	bin := buildLG(t)
	src := `(def ^:dynamic *x* :root) ` +
		`(def f (binding [*x* :caller] (bound-fn* (fn [] (sleep 10000) [*x* (scope-cancelled?)])))) ` +
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
	if !strings.Contains(string(out), "[:caller true]") {
		t.Fatalf("expected [:caller true], got: %q", out)
	}
	if d := time.Since(start); d > 9*time.Second {
		t.Fatalf("bound fn escaped its invoking scope: %v", d)
	}
}

// TestEvalRunsInCallerContext: a form passed to eval runs under the caller's
// dynamic bindings, so with-out-str around eval captures what the evaluated
// code prints, including from inside a future.
func TestEvalRunsInCallerContext(t *testing.T) {
	bin := buildLG(t)
	src := `(println (pr-str @(future (with-out-str (eval (read-string "(println :captured)"))))))`
	out, err := exec.Command(bin, "-e", src).CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), `":captured\n"`) {
		t.Fatalf("expected captured output, got: %q", out)
	}
}
