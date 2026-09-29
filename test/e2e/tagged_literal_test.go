/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os/exec"
	"testing"
)

// tagged-literal? must print its answer and nothing else: no load-time warning
// on stderr. A warning would only show in a fresh process.
func TestTaggedLiteralPredicatePrintsNoWarning(t *testing.T) {
	out, err := exec.Command(buildLG(t), "-e", "(tagged-literal? 1)").CombinedOutput()
	if err != nil {
		t.Fatalf("lg -e: %v\n%s", err, out)
	}
	if string(out) != "false\n" {
		t.Fatalf("want exactly false, got:\n%s", out)
	}
}
