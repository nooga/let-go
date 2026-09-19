/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Tests for issue #920: a scope's blocking natives (sleep, channel take/put)
// raise the Cancelled condition through their existing error slot on
// cancellation, instead of quietly returning nil. See also
// scope_cancelled_test.go, bound_fn_scope_test.go and scope_test.go for the
// pre-existing scope e2e tests updated for the same change.
package e2e

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// runLG runs the lg binary on src and returns (stdout, stderr) separately —
// unlike CombinedOutput-based helpers elsewhere in this package, tests here
// need to assert on stderr being EXACTLY empty (silent absorption), which a
// merged stream can't do.
func runLG(t *testing.T, bin string, src string) (stdout, stderr string, err error) {
	t.Helper()
	cmd := exec.Command(bin, "-e", src)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

// TestCancelledSleepRaisesCancelled: an interrupted (sleep ...) raises the
// Cancelled condition rather than returning nil.
func TestCancelledSleepRaisesCancelled(t *testing.T) {
	bin := buildLG(t)
	src := `(def result (promise)) ` +
		`(with-scope [s] ` +
		`  (go (deliver result ` +
		`        (try (sleep 10000) :no-cancellation ` +
		`             (catch Cancelled e (if (nil? e) :nil-caught :cancelled)))))` +
		`  (sleep 20)) ` +
		`(println (deref result 2000 :timeout))`
	start := time.Now()
	out, errOut, err := runLG(t, bin, src)
	if err != nil {
		t.Fatalf("run: %v\nstdout=%q stderr=%q", err, out, errOut)
	}
	if !strings.Contains(out, ":cancelled") {
		t.Fatalf("expected a cancelled sleep to raise Cancelled, got stdout=%q stderr=%q", out, errOut)
	}
	if d := time.Since(start); d > 9*time.Second {
		t.Fatalf("took %v — sleep was not interrupted", d)
	}
}

// TestCancelledTakeRaisesCancelled: an interrupted (<! ch) raises the
// Cancelled condition rather than returning nil.
func TestCancelledTakeRaisesCancelled(t *testing.T) {
	bin := buildLG(t)
	src := `(def ch (chan)) ` +
		`(def result (promise)) ` +
		`(with-scope [s] ` +
		`  (go (deliver result ` +
		`        (try (<! ch) :no-cancellation ` +
		`             (catch Cancelled e :cancelled))))` +
		`  (sleep 20)) ` +
		`(println (deref result 2000 :timeout))`
	start := time.Now()
	out, errOut, err := runLG(t, bin, src)
	if err != nil {
		t.Fatalf("run: %v\nstdout=%q stderr=%q", err, out, errOut)
	}
	if !strings.Contains(out, ":cancelled") {
		t.Fatalf("expected a cancelled take to raise Cancelled, got stdout=%q stderr=%q", out, errOut)
	}
	if d := time.Since(start); d > 9*time.Second {
		t.Fatalf("took %v — take was not interrupted", d)
	}
}

// TestCatchThrowableMissesCancelled: the idiomatic worker retry wrapper
// (catch Throwable e (log e) (recur)) must NOT catch cancellation — that is
// the entire point of #920 (a worker that swallows Throwable would otherwise
// never terminate on scope teardown). The worker increments an atom each
// pass through a (sleep 20); if Throwable caught Cancelled the loop would
// run all 100 iterations despite the scope closing after ~50ms.
func TestCatchThrowableMissesCancelled(t *testing.T) {
	bin := buildLG(t)
	src := `(def n (atom 0)) ` +
		`(def loops (atom 0)) ` +
		`(with-scope [s] ` +
		`  (go (loop [i 0] ` +
		`        (when (< i 100) ` +
		`          (try ` +
		`            (swap! n inc) ` +
		`            (sleep 20) ` +
		`            (swap! loops inc) ` +
		`            (catch Throwable e (swap! loops inc))) ` +
		`          (recur (inc i))))) ` +
		`  (sleep 50)) ` +
		`(println [@n @loops])`
	start := time.Now()
	out, errOut, err := runLG(t, bin, src)
	if err != nil {
		t.Fatalf("run: %v\nstdout=%q stderr=%q", err, out, errOut)
	}
	// n counts passes that started (before sleep); loops counts passes that
	// completed (whether normally or via the Throwable catch). If Throwable
	// caught Cancelled, loops would track n up to 100. Since it must NOT
	// catch it, exactly one iteration is started-but-never-completed: the one
	// whose sleep was interrupted by teardown.
	firstLine := strings.SplitN(strings.TrimSpace(out), "\n", 2)[0]
	fields := strings.Trim(firstLine, "[]")
	parts := strings.Fields(fields)
	if len(parts) != 2 {
		t.Fatalf("expected two numbers, got stdout=%q stderr=%q", out, errOut)
	}
	nVal, errN := strconv.Atoi(parts[0])
	loopsVal, errL := strconv.Atoi(parts[1])
	if errN != nil || errL != nil {
		t.Fatalf("expected numeric output, got stdout=%q", out)
	}
	if nVal != loopsVal+1 {
		t.Fatalf("expected exactly one in-flight iteration killed by cancellation (n=loops+1), got n=%d loops=%d — Throwable likely caught Cancelled", nVal, loopsVal)
	}
	if nVal >= 100 {
		t.Fatalf("loop ran to completion (n=%d) — cancellation did not stop it", nVal)
	}
	if d := time.Since(start); d > 9*time.Second {
		t.Fatalf("took %v — teardown did not interrupt the loop promptly", d)
	}
}

// TestCatchCancelledCatchesIt: an explicit (catch Cancelled e ...) DOES see
// the condition — the only spelling that does.
func TestCatchCancelledCatchesIt(t *testing.T) {
	bin := buildLG(t)
	src := `(def result (promise)) ` +
		`(with-scope [s] ` +
		`  (go (deliver result (try (sleep 10000) :ran (catch Cancelled e :caught))))` +
		`  (sleep 20)) ` +
		`(println (deref result 2000 :timeout))`
	out, errOut, err := runLG(t, bin, src)
	if err != nil {
		t.Fatalf("run: %v\nstdout=%q stderr=%q", err, out, errOut)
	}
	if !strings.Contains(out, ":caught") {
		t.Fatalf("expected (catch Cancelled e ...) to catch the condition, got stdout=%q stderr=%q", out, errOut)
	}
}

// TestFinallyRunsOnCancellationPath: finally still runs when a blocking op
// is interrupted by scope teardown, even with no catch present at all.
func TestFinallyRunsOnCancellationPath(t *testing.T) {
	bin := buildLG(t)
	src := `(def ran-finally (atom false)) ` +
		`(with-scope [s] ` +
		`  (go (try (sleep 10000) (finally (reset! ran-finally true)))) ` +
		`  (sleep 20)) ` +
		`(println @ran-finally)`
	out, errOut, err := runLG(t, bin, src)
	if err != nil {
		t.Fatalf("run: %v\nstdout=%q stderr=%q", err, out, errOut)
	}
	if !strings.Contains(out, "true") {
		t.Fatalf("expected finally to run on the cancellation path, got stdout=%q stderr=%q", out, errOut)
	}
}

// TestWithScopeStopsLoopEarly is the headline behaviour change from the
// issue: a loop bounded by iteration count, blocking on sleep each pass,
// stops early when its enclosing scope tears down instead of running to
// completion at full speed (the exact repro from the issue body).
func TestWithScopeStopsLoopEarly(t *testing.T) {
	bin := buildLG(t)
	src := `(def n (atom 0)) ` +
		`(with-scope [s] ` +
		`  (go (loop [i 0] (when (< i 100) (swap! n inc) (sleep 20) (recur (inc i))))) ` +
		`  (sleep 50)) ` +
		`(println @n)`
	start := time.Now()
	out, errOut, err := runLG(t, bin, src)
	if err != nil {
		t.Fatalf("run: %v\nstdout=%q stderr=%q", err, out, errOut)
	}
	firstLine := strings.SplitN(strings.TrimSpace(out), "\n", 2)[0]
	n, convErr := strconv.Atoi(firstLine)
	if convErr != nil {
		t.Fatalf("expected a number, got stdout=%q stderr=%q", out, errOut)
	}
	// ~50ms of teardown budget at 20ms/iteration is a handful of iterations,
	// nowhere near the uncancelled 100.
	if n >= 20 {
		t.Fatalf("expected the loop to stop early (well under 100), got n=%d", n)
	}
	if d := time.Since(start); d > 9*time.Second {
		t.Fatalf("took %v — loop was not actually interrupted", d)
	}
}

// TestOrdinaryTeardownIsSilent: with-scope tearing down several parked
// workers must print nothing to stderr — the absorption at the go* boundary
// (issue #920 part 3) swallows the Cancelled condition instead of routing
// it to *err* like an ordinary error.
func TestOrdinaryTeardownIsSilent(t *testing.T) {
	bin := buildLG(t)
	src := `(def ch (chan)) ` +
		`(with-scope [s] ` +
		`  (dotimes [_ 3] (go (<! ch))) ` +
		`  (go (sleep 10000)) ` +
		`  (sleep 30)) ` +
		`(println :done)`
	out, errOut, err := runLG(t, bin, src)
	if err != nil {
		t.Fatalf("run: %v\nstdout=%q stderr=%q", err, out, errOut)
	}
	if !strings.Contains(out, "done") {
		t.Fatalf("expected :done, got stdout=%q", out)
	}
	if strings.TrimSpace(errOut) != "" {
		t.Fatalf("expected silent teardown, got stderr=%q", errOut)
	}
}

// TestOrdinaryGoBlockErrorStillReachesErr: an ordinary (non-cancellation)
// error thrown inside a (go ...) block must still route to *err* (the
// process's stderr by default), unchanged by the cancellation-absorption fix
// in go* (issue #920 part 3 must not touch ordinary error handling).
func TestOrdinaryGoBlockErrorStillReachesErr(t *testing.T) {
	bin := buildLG(t)
	src := `(go (throw (ex-info "boom" {}))) (sleep 50) (println :done)`
	out, errOut, err := runLG(t, bin, src)
	if err != nil {
		t.Fatalf("run: %v\nstdout=%q stderr=%q", err, out, errOut)
	}
	if !strings.Contains(out, "done") {
		t.Fatalf("expected :done, got stdout=%q", out)
	}
	if !strings.Contains(errOut, "boom") {
		t.Fatalf("expected the ordinary go-block error to reach *err*/stderr, got stderr=%q", errOut)
	}
}
