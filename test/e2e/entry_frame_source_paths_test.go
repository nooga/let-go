/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// An entry-frame binary searches for source only where LG_SOURCE_PATHS says,
// as lg does when the variable is set (#198): it never adds the directory it
// runs in, so a stray .lg file there cannot satisfy a run-time require.
func TestEntryFrameRequireSearchesOnlyLGSourcePaths(t *testing.T) {
	if testing.Short() {
		t.Skip("builds an entry-frame binary; run without -short")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	root := repoRoot(t)
	bin := buildLG(t)
	out := t.TempDir()
	app := filepath.Join(out, "app.lg")
	// The require runs at load time, in the bytecode main chunk, so the
	// result depends only on the frame's search path.
	writeFile(t, app, `(ns cwdapp)

(def found
  (try (require 'cwdlib) ((resolve 'cwdlib/hello))
       (catch e :not-found)))

(defn -main [] (println found))
`)
	if o, err := runCmd(ctx, t, bin, root, []string{"scripts/lg-compile", "--entry-frame", out, "cwdapp", app}); err != nil {
		t.Fatalf("lg-compile --entry-frame: %v\n%s", err, o)
	}
	if o, err := runCmd(ctx, t, bin, root, []string{"-c", filepath.Join(out, "program.lgb"), "-entry-frame-entry", "cwdapp/-main", app}); err != nil {
		t.Fatalf("lg -c: %v\n%s", err, o)
	}
	writeFile(t, filepath.Join(out, "go.mod"), "module cwdapp\n\ngo "+rootGoDirective(t, root)+
		"\n\nrequire github.com/nooga/let-go v0.0.0\n\nreplace github.com/nooga/let-go => "+root+"\n")
	if o, err := runCmd(ctx, t, "go", out, []string{"mod", "tidy"}); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, o)
	}
	exe := filepath.Join(out, "app")
	if o, err := runCmd(ctx, t, "go", out, []string{"build", "-o", exe, "."}); err != nil {
		t.Fatalf("go build: %v\n%s", err, o)
	}

	cwd := t.TempDir()
	writeFile(t, filepath.Join(cwd, "cwdlib.lg"), "(ns cwdlib)\n(defn hello [] :loaded)\n")
	run := func(env ...string) string {
		t.Helper()
		cmd := exec.CommandContext(ctx, exe)
		cmd.Dir = cwd
		cmd.Env = append(append(os.Environ(), "LG_SOURCE_PATHS="), env...)
		got, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("binary: %v\n%s", err, got)
		}
		return strings.TrimSpace(string(got))
	}
	unset := func() string {
		t.Helper()
		cmd := exec.CommandContext(ctx, exe)
		cmd.Dir = cwd
		for _, kv := range os.Environ() {
			if !strings.HasPrefix(kv, "LG_SOURCE_PATHS=") {
				cmd.Env = append(cmd.Env, kv)
			}
		}
		got, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("binary: %v\n%s", err, got)
		}
		return strings.TrimSpace(string(got))
	}
	if got := unset(); got != ":not-found" {
		t.Errorf("LG_SOURCE_PATHS unset: got %q, want :not-found (the run directory is not searched)", got)
	}
	if got := run("LG_SOURCE_PATHS=" + cwd); got != ":loaded" {
		t.Errorf("LG_SOURCE_PATHS=%s: got %q, want :loaded", cwd, got)
	}
}
