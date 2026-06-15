/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"os/exec"
	"strings"
	"testing"
)

// TestDeftypeSkeletonNativeLowering drives the native-lowering trampoline
// (scripts/gogen-trampoline.lg) over the deftype walking-skeleton fixture
// (test/gogen/deftype_skeleton.lg) and asserts it passes.
//
// The trampoline lowers the fixture's defprotocol/deftype/defn forms to a Go
// package, wires that package in under a build tag so its init() registers the
// Go-native overrides, then `require`s the fixture as a library ns — which makes
// the resolver drain the overrides so (run) dispatches to the GENERATED code. It
// checks (run) matches the bytecode VM AND that the entry var is actually a
// native fn under the tag (so the comparison isn't a vacuous bytecode-to-
// bytecode run). See the script header for the full mechanism.
//
// This is the committed, `go test ./...`-discovered gate for the feature; the
// trampoline script is the reusable, fixture-agnostic suite it delegates to.
// (buildLG is defined in scope_e2e_test.go, same package.)
func TestDeftypeSkeletonNativeLowering(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping native-lowering trampoline (builds go subprocesses) in -short mode")
	}

	// The trampoline is itself an lg script: run it on a freshly-built lg, and
	// pass that same lg as the bytecode-side binary it compares against.
	lg := buildLG(t)
	out, err := exec.Command(lg, "scripts/gogen-trampoline.lg", "--lg", lg).CombinedOutput()
	if err != nil {
		t.Fatalf("gogen-trampoline.lg failed: %v\n%s", err, out)
	}
	// Guard against a vacuous pass: the trampoline must have confirmed the
	// tagged run dispatched to the lowered code, not silently re-run bytecode.
	if !strings.Contains(string(out), "native-dispatch ok") {
		t.Fatalf("trampoline did not confirm native dispatch:\n%s", out)
	}
}
