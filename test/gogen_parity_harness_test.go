/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestGogenParitySelftest(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	run := exec.CommandContext(ctx, "bash", "scripts/gogen-parity.sh", "--selftest")
	run.Dir = root
	if out, err := run.CombinedOutput(); err != nil {
		t.Fatalf("gogen parity harness: %v\n%s", err, out)
	}
}
