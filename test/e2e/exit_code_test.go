/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const (
	exitCodeUserRead  os.FileMode = 1 << 8
	exitCodeUserWrite os.FileMode = 1 << 7
	exitCodeGroupRead os.FileMode = 1 << 5
	exitCodeOtherRead os.FileMode = 1 << 2

	exitCodeTestFilePerm = exitCodeUserRead | exitCodeUserWrite | exitCodeGroupRead | exitCodeOtherRead
)

// TestErrorExitCode verifies lg exits nonzero when a script or -e expression
// fails with a runtime error, and zero when it succeeds. Previously runtime
// errors were printed but the process still exited 0, so shell pipelines and
// CI runs could not detect failure.
func TestErrorExitCode(t *testing.T) {
	bin := buildLG(t)

	writeScript := func(t *testing.T, body string) string {
		t.Helper()
		p := filepath.Join(t.TempDir(), "app.lg")
		if err := os.WriteFile(p, []byte(body), exitCodeTestFilePerm); err != nil {
			t.Fatal(err)
		}
		return p
	}

	// run executes lg and returns its exit code and combined output.
	run := func(t *testing.T, args ...string) (int, string) {
		t.Helper()
		out, err := exec.Command(bin, args...).CombinedOutput()
		if err != nil {
			var ee *exec.ExitError
			if !errors.As(err, &ee) {
				t.Fatalf("run lg %v: %v\n%s", args, err, out)
			}
			return ee.ExitCode(), string(out)
		}
		return 0, string(out)
	}

	t.Run("failing script exits nonzero", func(t *testing.T) {
		app := writeScript(t, `(undefined-fn-xyz)`)
		if code, out := run(t, app); code == 0 {
			t.Fatalf("want nonzero exit, got 0; output:\n%s", out)
		}
	})

	t.Run("throwing script exits nonzero", func(t *testing.T) {
		app := writeScript(t, `(throw (ex-info "boom" {}))`)
		if code, out := run(t, app); code == 0 {
			t.Fatalf("want nonzero exit, got 0; output:\n%s", out)
		}
	})

	t.Run("succeeding script exits zero", func(t *testing.T) {
		app := writeScript(t, `(println :ok)`)
		if code, out := run(t, app); code != 0 {
			t.Fatalf("want exit 0, got %d; output:\n%s", code, out)
		}
	})

	t.Run("failing -e exits nonzero", func(t *testing.T) {
		if code, out := run(t, "-e", `(undefined-fn-xyz)`); code == 0 {
			t.Fatalf("want nonzero exit, got 0; output:\n%s", out)
		}
	})

	t.Run("succeeding -e exits zero", func(t *testing.T) {
		if code, out := run(t, "-e", `(+ 1 2)`); code != 0 {
			t.Fatalf("want exit 0, got %d; output:\n%s", code, out)
		}
	})

	t.Run("caught error still exits zero", func(t *testing.T) {
		app := writeScript(t, `(println (try (throw (ex-info "boom" {})) (catch e :caught)))`)
		if code, out := run(t, app); code != 0 {
			t.Fatalf("want exit 0, got %d; output:\n%s", code, out)
		}
	})

	// The one carve-out in runMain: a failing script under -r drops into the
	// REPL, which recovers and exits clean. Feeding /dev/null gives the REPL an
	// immediate EOF so it returns instead of blocking.
	t.Run("failing script with -r exits zero", func(t *testing.T) {
		app := writeScript(t, `(undefined-fn-xyz)`)
		cmd := exec.Command(bin, "-r", app)
		cmd.Stdin = nil // nil Stdin reads from the null device -> EOF
		out, err := cmd.CombinedOutput()
		if err != nil {
			var ee *exec.ExitError
			if !errors.As(err, &ee) {
				t.Fatalf("run lg -r %s: %v\n%s", app, err, out)
			}
			t.Fatalf("want exit 0, got %d; output:\n%s", ee.ExitCode(), out)
		}
	})
}
