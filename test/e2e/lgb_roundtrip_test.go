/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestLgbRoundtrip runs scripts/lgb-roundtrip.sh against a freshly built lg.
// The selftest runs in every lane. The examples corpus compiles each program
// in all four serialization modes, which costs several seconds, so it runs
// only outside -short.
func TestLgbRoundtrip(t *testing.T) {
	lg := buildLG(t)
	run := func(t *testing.T, args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, "bash", append([]string{"scripts/lgb-roundtrip.sh"}, args...)...)
		cmd.Dir = repoRoot(t)
		cmd.Env = append(os.Environ(), "LG="+lg)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("lgb-roundtrip.sh %v: %v\n%s", args, err, out)
		}
	}

	t.Run("selftest", func(t *testing.T) { run(t, "--selftest") })
	t.Run("examples", func(t *testing.T) {
		if testing.Short() {
			t.Skip("compiles every example four ways; runs in the non-short e2e step")
		}
		run(t)
	})
}
