/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/genmanifest"
)

// Differential self-AOT execution harness.
//
// let-go has two execution engines for the same source: the default bytecode
// VM, and — under `-tags gogen_ir` — the native IR pipeline, where the stdlib
// + IR passes are pre-lowered to Go (pkg/rt/core_go_lowered/). A program that
// produces a different result under the two engines is a lowering divergence.
//
// This test builds let-go twice (bytecode + gogen_ir), runs every
// test/gold-aot/*.lg fixture under both, and compares the last non-empty stdout
// line. The bytecode run is the reference (a fixture whose bytecode run fails
// is a broken fixture, not a divergence).
//
// The native tree is not yet fully green (see the gogen_ir pkg/ir suite and the
// repair plan), so divergences are tracked in a SHRINK-ONLY allowlist,
// test/gogen_aot_xfail.txt:
//
//   - a diverging fixture listed in the allowlist is tolerated (xfail);
//   - a NEW divergence not in the allowlist FAILS the test (execution regression);
//   - an allowlisted fixture that now AGREES also FAILS ("remove from xfail") —
//     the ratchet that forces the list to shrink as lowering bugs are fixed.
//
// Re-seed the allowlist from the current divergence set with
// LETGO_AOT_REDERIVE=1 (use after an intentional, reviewed change).
//
// Gated behind testing.Short(): plain `go test -short ./...` skips the
// double-build. CI runs it explicitly via `make gogen-diff`.

// goldAOTDir and aotXfailFile are resolved relative to the repo root at
// runtime (see TestGogenAOTDiff); this test package runs with its own
// directory as CWD, not the repo root.
const aotRedriveEnv = "LETGO_AOT_REDERIVE"

// buildLGTags builds the lg binary with the given build tags ("" = the default
// bytecode engine; "gogen_ir" = the native IR pipeline). Returns the path.
func buildLGTags(t *testing.T, tags string) string {
	t.Helper()
	name := "lg-bc"
	args := []string{"build", "-o", ""}
	if tags != "" {
		name = "lg-" + tags
		args = []string{"build", "-tags", tags, "-o", ""}
	}
	bin := filepath.Join(t.TempDir(), name)
	args[len(args)-1] = bin
	args = append(args, ".")
	cmd := exec.Command("go", args...)
	cmd.Dir = repoRoot(t) // build the root main package, not this test dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build lg (tags=%q): %v\n%s", tags, err, out)
	}
	return bin
}

// runFixture runs script under bin and returns (normalized full stdout, ok). ok
// is false when the process exits non-zero, in which case the token is a sentinel
// so an erroring engine always differs from a succeeding one. Comparing full
// output (not only the last line) catches mid-output divergences (P2).
func runFixture(bin, script string) (token string, ok bool) {
	cmd := exec.Command(bin, script)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "<<runtime-error>>", false
	}
	return normalizeOutput(out.String()), true
}

// normalizeOutput returns full stdout with per-line trailing whitespace and
// trailing blank lines stripped, so the cross-engine diff catches mid-output
// divergences rather than only the final line (P2).
func normalizeOutput(s string) string {
	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		lines[i] = strings.TrimRight(ln, " \t\r")
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}

func TestGogenAOTDiff(t *testing.T) {
	if testing.Short() {
		t.Skip("differential self-AOT harness builds let-go twice; run via `make gogen-diff`")
	}

	root := repoRoot(t)

	// The -tags gogen_ir build below imports the gitignored generated tree.
	// Validate it against its completeness sentinel first, so a torn or
	// missing tree fails with an actionable message instead of a wall of
	// "no required module provides package .../core_go_lowered/..." errors.
	if err := genmanifest.CheckTreeManifest(filepath.Join(root, "pkg/rt/core_go_lowered")); err != nil {
		t.Fatalf("generated lowered tree is not valid: %v\nrun `make generate` (or `make lowered`) and retry", err)
	}

	bc := buildLGTags(t, "")
	aot := buildLGTags(t, "gogen_ir")
	goldAOTDir := filepath.Join(root, "test/gold-aot")
	aotXfailFile := filepath.Join(root, "test/gogen_aot_xfail.txt")

	fixtures, err := filepath.Glob(filepath.Join(goldAOTDir, "*.lg"))
	if err != nil {
		t.Fatalf("glob %s: %v", goldAOTDir, err)
	}
	if len(fixtures) == 0 {
		t.Fatalf("no fixtures in %s", goldAOTDir)
	}
	sort.Strings(fixtures)

	diverged := map[string]bool{}
	for _, f := range fixtures {
		name := filepath.Base(f)
		bcOut, bcOk := runFixture(bc, f)
		if !bcOk {
			t.Errorf("%s: bytecode (reference) run failed — fix the fixture", name)
			continue
		}
		aotOut, aotOk := runFixture(aot, f)
		if !aotOk || bcOut != aotOut {
			diverged[name] = true
			t.Logf("DIVERGE %s: bc=%q aot=%q (aot_ok=%v)", name, bcOut, aotOut, aotOk)
		}
	}

	if os.Getenv(aotRedriveEnv) == "1" {
		writeXfail(t, aotXfailFile, diverged)
		t.Logf("re-seeded %s with %d diverging fixture(s)", aotXfailFile, len(diverged))
		return
	}

	allow := readXfail(t, aotXfailFile)

	// New divergence not on the allowlist → regression.
	for name := range diverged {
		if !allow[name] {
			t.Errorf("NEW gogen_ir execution divergence: %s (not in %s).\n"+
				"  A native-lowering regression: the fixture runs differently under -tags gogen_ir.\n"+
				"  Fix the lowering, or — if the divergence is known/triaged — re-seed the allowlist with\n"+
				"  `%s=1 go test -run TestGogenAOTDiff .` and commit the updated %s.",
				name, aotXfailFile, aotRedriveEnv, aotXfailFile)
		}
	}
	// Allowlisted fixture that now agrees → ratchet violation (must shrink).
	for name := range allow {
		if !diverged[name] {
			t.Errorf("%s is in %s but now AGREES under gogen_ir — remove it from the allowlist (shrink-only ratchet).",
				name, aotXfailFile)
		}
	}
}

// readXfail loads the allowlist: one fixture base-name per line, blank lines and
// '#' comments ignored. A missing file is an empty allowlist.
func readXfail(t *testing.T, path string) map[string]bool {
	t.Helper()
	m := map[string]bool{}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return m
		}
		t.Fatalf("read %s: %v", path, err)
	}
	for _, line := range strings.Split(string(b), "\n") {
		s := strings.TrimSpace(line)
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		m[s] = true
	}
	return m
}

// writeXfail rewrites the allowlist from the current divergence set (sorted),
// preserving a fixed header. Used only under LETGO_AOT_REDERIVE=1.
func writeXfail(t *testing.T, path string, diverged map[string]bool) {
	t.Helper()
	names := make([]string, 0, len(diverged))
	for n := range diverged {
		names = append(names, n)
	}
	sort.Strings(names)
	var b strings.Builder
	b.WriteString("# Shrink-only allowlist of test/gold-aot/*.lg fixtures that currently\n")
	b.WriteString("# DIVERGE between the bytecode and -tags gogen_ir engines.\n")
	b.WriteString("# Managed by TestGogenAOTDiff: re-seed with `LETGO_AOT_REDERIVE=1 go test -run TestGogenAOTDiff .`.\n")
	b.WriteString("# A new divergence not listed here fails CI; a listed fixture that starts\n")
	b.WriteString("# agreeing also fails (remove it). The list must only ever shrink.\n")
	for _, n := range names {
		b.WriteString(n)
		b.WriteString("\n")
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
