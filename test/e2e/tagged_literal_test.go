/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os/exec"
	"testing"
)

// The first tagged-literal? call loads let-go.tagged-literal, which defines its
// own tagged-literal?; loading it must not warn that the name shadows
// clojure.core's. The warning goes to stderr, so it needs a fresh process.
func TestFirstTaggedLiteralCallPrintsNoShadowWarning(t *testing.T) {
	out, err := exec.Command(buildLG(t), "-e", "(tagged-literal? 1)").CombinedOutput()
	if err != nil {
		t.Fatalf("lg -e: %v\n%s", err, out)
	}
	if string(out) != "false\n" {
		t.Fatalf("want exactly false, got:\n%s", out)
	}
}
