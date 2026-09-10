/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNamespaceShadowWarningHarness(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	lg := filepath.Join(t.TempDir(), "lg")
	build := exec.CommandContext(ctx, "go", "build", "-o", lg, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build lg: %v\n%s", err, out)
	}

	run := exec.CommandContext(ctx, "bash", "test/namespace_shadow_warning_test/run.sh")
	run.Dir = root
	run.Env = replaceEnv(os.Environ(), "LG", lg)
	if out, err := run.CombinedOutput(); err != nil {
		t.Fatalf("shadow-warning harness: %v\n%s", err, out)
	}
}

func replaceEnv(env []string, key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	for _, entry := range env {
		if !strings.HasPrefix(entry, prefix) {
			out = append(out, entry)
		}
	}
	return append(out, prefix+value)
}
