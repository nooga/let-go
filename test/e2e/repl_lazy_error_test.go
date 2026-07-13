/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestREPLRecoversFromLazyErrorDuringResultRendering(t *testing.T) {
	bin := buildLG(t)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader("(map (fn [x] (throw \"boom while printing\")) [1])\n(+ 1 2)\n")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		t.Fatalf("REPL exited after lazy result error: %v\n%s", err, output.String())
	}
	got := output.String()
	if strings.Contains(got, "panic:") || strings.Contains(got, "goroutine ") {
		t.Fatalf("REPL exposed a raw Go panic:\n%s", got)
	}
	if !strings.Contains(got, "boom while printing") {
		t.Fatalf("REPL did not render the user error:\n%s", got)
	}
	if !strings.Contains(got, "\n3\n") {
		t.Fatalf("REPL did not evaluate the next form:\n%s", got)
	}
}

func TestExpressionRenderingConvertsLazyThrowToUserError(t *testing.T) {
	bin := buildLG(t)
	out, err := exec.Command(bin, "-e", "(map (fn [x] (throw \"boom in expression\")) [1])").CombinedOutput()
	if err == nil {
		t.Fatalf("expression runner returned zero after lazy result error:\n%s", out)
	}
	if _, ok := err.(*exec.ExitError); !ok {
		t.Fatalf("expression runner failed outside normal process exit: %v\n%s", err, out)
	}
	got := string(out)
	if strings.Contains(got, "panic:") || strings.Contains(got, "goroutine ") {
		t.Fatalf("expression runner exposed a raw Go panic:\n%s", got)
	}
	if !strings.Contains(got, "boom in expression") {
		t.Fatalf("expression runner did not render the user error:\n%s", got)
	}
}
