/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	userRead  os.FileMode = 1 << 8
	userWrite os.FileMode = 1 << 7
	groupRead os.FileMode = 1 << 5
	otherRead os.FileMode = 1 << 2

	testFilePerm = userRead | userWrite | groupRead | otherRead
)

// TestReadLineStdin verifies (read-line) with no args reads the current *in*,
// like Clojure, while the 1-arg handle form keeps working.
func TestReadLineStdin(t *testing.T) {
	bin := buildLG(t)

	run := func(t *testing.T, stdin string, body string) string {
		t.Helper()
		p := filepath.Join(t.TempDir(), "app.lg")
		if err := os.WriteFile(p, []byte(body), testFilePerm); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(bin, p)
		cmd.Stdin = strings.NewReader(stdin)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("run lg: %v\n%s", err, out)
		}
		return strings.TrimSpace(string(out))
	}

	t.Run("no args reads stdin", func(t *testing.T) {
		if got := run(t, "hello\n", `(println (read-line))`); got != "hello" {
			t.Fatalf("got %q, want hello", got)
		}
	})

	t.Run("nil at EOF", func(t *testing.T) {
		if got := run(t, "", `(prn (read-line))`); got != "nil" {
			t.Fatalf("got %q, want nil", got)
		}
	})

	t.Run("explicit handle still works", func(t *testing.T) {
		if got := run(t, "hi\n", `(println (read-line *in*))`); got != "hi" {
			t.Fatalf("got %q, want hi", got)
		}
	})

	t.Run("no args respects rebound *in*", func(t *testing.T) {
		input := filepath.Join(t.TempDir(), "input.txt")
		if err := os.WriteFile(input, []byte("from-file\n"), testFilePerm); err != nil {
			t.Fatal(err)
		}
		body := fmt.Sprintf(`(let [f (open %q :read)]
  (binding [*in* f]
    (println (read-line))))`, input)
		if got := run(t, "from-stdin\n", body); got != "from-file" {
			t.Fatalf("got %q, want from-file", got)
		}
	})
}
