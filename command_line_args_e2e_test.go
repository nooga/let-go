/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCommandLineArgs verifies core/*command-line-args* holds exactly the
// user's args — the positionals after the script — verbatim, and nil when
// there are none. It covers script mode (with and without a leading flag) and
// a bundled binary, which exercise lg's two set-points. os/args is left as the
// full process argv. (buildLG is defined in scope_e2e_test.go, same package.)
func TestCommandLineArgs(t *testing.T) {
	bin := buildLG(t)

	// writeScript writes body to app.lg in a fresh temp dir and returns its path.
	writeScript := func(t *testing.T, body string) string {
		t.Helper()
		p := filepath.Join(t.TempDir(), "app.lg")
		if err := os.WriteFile(p, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
		return p
	}

	// run executes lg with args and returns trimmed combined output.
	run := func(t *testing.T, args ...string) string {
		t.Helper()
		out, err := exec.Command(bin, args...).CombinedOutput()
		if err != nil {
			t.Fatalf("run lg %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}

	t.Run("script with args", func(t *testing.T) {
		app := writeScript(t, `(prn *command-line-args*)`)
		if got := run(t, app, "a", "b"); got != `("a" "b")` {
			t.Fatalf("got %q, want (\"a\" \"b\")", got)
		}
	})

	t.Run("flag before script does not shift args", func(t *testing.T) {
		app := writeScript(t, `(prn *command-line-args*)`)
		// "." is listed, so no source-paths transition warning is emitted —
		// combined output stays clean.
		if got := run(t, "-source-paths", ".", app, "a", "b"); got != `("a" "b")` {
			t.Fatalf("got %q, want (\"a\" \"b\")", got)
		}
	})

	t.Run("no args is nil", func(t *testing.T) {
		app := writeScript(t, `(prn *command-line-args*)`)
		if got := run(t, app); got != "nil" {
			t.Fatalf("got %q, want nil", got)
		}
	})

	t.Run("literal -- is preserved", func(t *testing.T) {
		app := writeScript(t, `(prn *command-line-args*)`)
		want := `("git" "checkout" "--" "file")`
		if got := run(t, app, "git", "checkout", "--", "file"); got != want {
			t.Fatalf("got %q, want %s", got, want)
		}
	})

	t.Run("bundled binary", func(t *testing.T) {
		app := writeScript(t, `(prn *command-line-args*)`)
		outBin := filepath.Join(t.TempDir(), "app")
		// The -b build runs the top-level form at AOT (printing nil); we only
		// care that it succeeds. The user args land when the bundle is run.
		if out, err := exec.Command(bin, "-b", outBin, app).CombinedOutput(); err != nil {
			t.Fatalf("bundle: %v\n%s", err, out)
		}
		out, err := exec.Command(outBin, "a", "b").CombinedOutput()
		if err != nil {
			t.Fatalf("run bundle: %v\n%s", err, out)
		}
		if got := strings.TrimSpace(string(out)); got != `("a" "b")` {
			t.Fatalf("got %q, want (\"a\" \"b\")", got)
		}
	})

	// os/args must stay the full process argv — the back-compat promise that
	// lets *command-line-args* be additive.
	t.Run("os/args still full argv", func(t *testing.T) {
		app := writeScript(t, `(prn *command-line-args*) (prn os/args)`)
		lines := strings.Split(run(t, app, "a", "b"), "\n")
		if len(lines) != 2 {
			t.Fatalf("want 2 output lines, got %d: %q", len(lines), lines)
		}
		if lines[0] != `("a" "b")` {
			t.Fatalf("*command-line-args*: got %q, want (\"a\" \"b\")", lines[0])
		}
		// os/args is the full argv as a vector, ending in the user args and
		// still carrying the program name and script path.
		if !strings.HasPrefix(lines[1], "[") || !strings.HasSuffix(lines[1], `"a" "b"]`) {
			t.Fatalf("os/args: got %q, want full argv vector ending in \"a\" \"b\"]", lines[1])
		}
		if !strings.Contains(lines[1], "app.lg") {
			t.Fatalf("os/args: got %q, want it to contain the script path", lines[1])
		}
	})
}
