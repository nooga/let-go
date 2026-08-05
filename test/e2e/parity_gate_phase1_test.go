/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/genmanifest"
)

// Phase-1 Backend Parity Gate.
//
// Runs Tier-1 fixtures (genuine AOT, pass strict mode, both backends agree)
// under both bytecode and -tags gogen_ir. Reports divergences against per-backend
// shrink-only allowlists (test/parity-xfail-{bytecode,gogen}.txt).
//
// Three-tier strict model:
// - Tier 1: Both backends pass *ir-compile-strict* and agree → counts toward parity
// - Tier 2: Bytecode fails strict (would trampoline) → excluded from parity coverage
// - Tier 3: Bytecode ≠ gogen_ir output → goes to allowlist (known miscompile)
//
// Shrink-only ratchet:
// - NEW divergence not on allowlist → FAIL (regression detected)
// - Allowlisted divergence now agrees → FAIL (remove from allowlist; improvement must be explicit)
//
// Gated: `go test -run TestParityGatePhase1 ./test/e2e` (or via `make parity-gate-phase1`)

// Tier-1 fixtures (both backends pass strict mode and agree).
// Updated by TestStrictModeAudit; manually curated for Phase-1.
var tier1Fixtures = []string{
	"arith.lg",
	"closure.lg",
	"destructure.lg",
	"quot.lg",
	"seq.lg",
	"setbang.lg",
	"setvar_load_stability.lg",
	"setvar_load_straightline.lg",
	"tryfinally.lg",
	"var.lg",
}

// Tier-2 fixtures (fail strict mode; excluded from parity coverage).
// These fixtures are NOT run in the gate — only Tier-1 counts toward parity.
var tier2Fixtures = []string{
	"binding.lg",
	"closure_capture.lg",
	"ir_pipeline.lg",
	"math_float.lg",
	"pipeline_dump_ir.lg",
	"typed_cross_fn.lg",
}

type ParityResult struct {
	Fixture string
	Backend string // "bytecode" or "gogen_ir"
	Output  string
	Ok      bool
}

func TestParityGatePhase1(t *testing.T) {
	if testing.Short() {
		t.Skip("parity gate builds let-go twice; run via `make parity-gate-phase1`")
	}

	root := repoRoot(t)

	if err := genmanifest.CheckTreeManifest(filepath.Join(root, "pkg/rt/core_go_lowered")); err != nil {
		t.Fatalf("generated lowered tree is not valid: %v\nrun `make generate` and retry", err)
	}

	bc := buildLGTags(t, "")
	aot := buildLGTags(t, "gogen_ir")
	goldAOTDir := filepath.Join(root, "test/gold-aot")
	wrapDir := t.TempDir()

	// Create strict-mode wrappers for all tier-1 fixtures
	wrappedFixtures := make(map[string]string)
	for _, name := range tier1Fixtures {
		path := filepath.Join(goldAOTDir, name)
		wrappedPath := filepath.Join(wrapDir, name)
		if err := createStrictWrapper(wrappedPath, path); err != nil {
			t.Fatalf("create wrapper for %s: %v", name, err)
		}
		wrappedFixtures[name] = wrappedPath
	}

	// Run tier-1 fixtures on both backends
	diverged := map[string]bool{}
	for _, name := range tier1Fixtures {
		wrappedPath := wrappedFixtures[name]
		bcOut, bcOk := runFixture(bc, wrappedPath)
		aotOut, aotOk := runFixture(aot, wrappedPath)

		// Both should succeed and agree
		if !bcOk || !aotOk || bcOut != aotOut {
			diverged[name] = true
			t.Logf("DIVERGE %s: bc_ok=%v bc_out=%q aot_ok=%v aot_out=%q",
				name, bcOk, bcOut, aotOk, aotOut)
		} else {
			t.Logf("AGREE %s: output=%q", name, bcOut)
		}
	}

	// Load per-backend allowlists
	bcXfail := readParityXfail(t, filepath.Join(root, "test/parity-xfail-bytecode.txt"))
	aotXfail := readParityXfail(t, filepath.Join(root, "test/parity-xfail-gogen.txt"))

	// Check for regressions: new divergence not on allowlist
	for name := range diverged {
		// Check both allowlists (could be in either backend's xfail)
		bcAllowed := bcXfail[name]
		aotAllowed := aotXfail[name]
		if !bcAllowed && !aotAllowed {
			t.Errorf("NEW backend parity divergence: %s (not in allowlists).\n"+
				"  A lowering divergence: the fixture runs differently under bytecode vs -tags gogen_ir.\n"+
				"  Fix the lowering, or — if triaged — add to test/parity-xfail-{bytecode,gogen}.txt\n"+
				"  and commit with a detailed reason comment.",
				name)
		}
	}

	// Check for stale allowlist entries: allowlisted fixture now agrees
	for name := range bcXfail {
		if !diverged[name] {
			t.Errorf("%s is in test/parity-xfail-bytecode.txt but now AGREES — remove it (shrink-only ratchet).",
				name)
		}
	}
	for name := range aotXfail {
		if !diverged[name] {
			t.Errorf("%s is in test/parity-xfail-gogen.txt but now AGREES — remove it (shrink-only ratchet).",
				name)
		}
	}

	t.Logf("\n=== PHASE-1 GATE SUMMARY ===")
	t.Logf("Tier-1 fixtures tested: %d", len(tier1Fixtures))
	t.Logf("Divergences detected:   %d", len(diverged))
	t.Logf("Bytecode allowlist:     %d entries", len(bcXfail))
	t.Logf("Gogen_ir allowlist:     %d entries", len(aotXfail))
}

// readParityXfail loads the per-backend allowlist: one fixture base-name per line,
// blank lines and '#' comments ignored.
func readParityXfail(t *testing.T, path string) map[string]bool {
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

// writeParityXfail rewrites the per-backend allowlist (used only with LETGO_PARITY_REDERIVE=1).
func writeParityXfail(t *testing.T, path string, diverged map[string]bool) {
	t.Helper()
	names := make([]string, 0, len(diverged))
	for n := range diverged {
		names = append(names, n)
	}
	sort.Strings(names)
	var b strings.Builder
	b.WriteString("# Shrink-only allowlist of Tier-1 fixtures that currently DIVERGE\n")
	b.WriteString("# between bytecode and -tags gogen_ir engines.\n")
	b.WriteString("# Managed by TestParityGatePhase1: re-seed with `LETGO_PARITY_REDERIVE=1 go test -run TestParityGatePhase1 ./test/e2e`.\n")
	b.WriteString("# A new divergence fails CI; a listed fixture that now agrees also fails (shrink-only).\n")
	b.WriteString("# When committing an allowlist update, include the specific issue/fix in the reason comment.\n\n")
	for _, n := range names {
		b.WriteString(n)
		b.WriteString("\n")
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
