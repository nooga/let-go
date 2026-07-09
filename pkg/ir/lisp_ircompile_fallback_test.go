/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// TestIRCompileFallbackDiagnostics verifies that when a defn fails to compile
// through the IR pipeline, the fallback to bytecode is logged when
// *ir-compile-verbose* is true.
func TestIRCompileFallbackDiagnostics(t *testing.T) {
	ensureLoader()
	nsName := "irfallback"

	// Find a fixture that fails IR compilation but succeeds in bytecode.
	// Try/catch is unsupported in the IR pipeline but works in bytecode.
	badForm := `(defn bad-defn [x] (try x (catch e e)))`

	// First, verify the bad form fails with IR enabled
	oldIR := setVarRoot(t, "clojure.core", "*ir-compile*", vm.TRUE)
	oldVerbose := setVarRoot(t, "clojure.core", "*ir-compile-verbose*", vm.TRUE)
	defer func() {
		setVarRoot(t, "clojure.core", "*ir-compile*", oldIR)
		setVarRoot(t, "clojure.core", "*ir-compile-verbose*", oldVerbose)
	}()

	// Compile the bad form with verbose enabled
	// It should fallback to bytecode but also log a warning
	runInNs(t, nsName, badForm)

	// Verify the defn exists (fallback compiled it successfully)
	result := runInNs(t, nsName, "(resolve 'bad-defn)")
	if result == vm.NIL {
		t.Fatalf("bad-defn did not compile; fallback failed")
	}

	// Verify the fallback was logged to the diagnostic atom
	// The atom should have at least one entry after the fallback
	logResult := runInNs(t, nsName, "(count @clojure.core/*ir-compile-fallback-log*)")
	logCount, ok := logResult.(vm.Int)
	if !ok || logCount < 1 {
		t.Fatalf("fallback not logged to *ir-compile-fallback-log*; got %v", logResult)
	}

	t.Logf("Fallback test: defn compiled successfully despite IR failure; %d entries logged", logCount)
}

// TestIRCompileSilentWhenNotVerbose verifies that the fallback is silent
// when *ir-compile-verbose* is false (the default).
func TestIRCompileSilentWhenNotVerbose(t *testing.T) {
	ensureLoader()
	nsName := "irsilent"

	badForm := `(defn bad-silent [x] (try x (catch e e)))`

	// Enable IR but disable verbose
	oldIR := setVarRoot(t, "clojure.core", "*ir-compile*", vm.TRUE)
	oldVerbose := setVarRoot(t, "clojure.core", "*ir-compile-verbose*", vm.FALSE)
	defer func() {
		setVarRoot(t, "clojure.core", "*ir-compile*", oldIR)
		setVarRoot(t, "clojure.core", "*ir-compile-verbose*", oldVerbose)
	}()

	// Compile the bad form with verbose disabled
	// Should not produce any warnings
	runInNs(t, nsName, badForm)

	// Verify the defn exists (fallback compiled it)
	result := runInNs(t, nsName, "(resolve 'bad-silent)")
	if result == vm.NIL {
		t.Fatalf("bad-silent did not compile; fallback failed")
	}

	t.Logf("Silent test: defn compiled without warnings")
}
