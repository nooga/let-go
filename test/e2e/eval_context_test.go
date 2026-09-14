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
)

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
