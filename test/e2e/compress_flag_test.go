/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestCompressFlagRequiresCompileOrBundle(t *testing.T) {
	lg := buildLG(t)
	cmd := exec.Command(lg, "-z")
	out, err := cmd.CombinedOutput()
	var ee *exec.ExitError
	if !errors.As(err, &ee) || ee.ExitCode() != 2 {
		t.Fatalf("want exit 2 for bare -z, got %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "-z requires -c or -b") {
		t.Fatalf("want -z requires message, got:\n%s", out)
	}
}
