/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"fmt"
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

func TestClassifyParityResultsSeparatesStrictFailuresFromOutputDivergences(t *testing.T) {
	results := []ParityResult{
		{Fixture: "bytecode-strict-fail.lg", BCOutput: "<<runtime-error>>", BCOK: false, AOTOutput: "ok", AOTOK: true},
		{Fixture: "gogen-strict-fail.lg", BCOutput: "ok", BCOK: true, AOTOutput: "<<runtime-error>>", AOTOK: false},
		{Fixture: "output-divergence.lg", BCOutput: "bytecode", BCOK: true, AOTOutput: "gogen", AOTOK: true},
		{Fixture: "agreement.lg", BCOutput: "same", BCOK: true, AOTOutput: "same", AOTOK: true},
	}

	strictFailures, diverged := classifyParityResults(results)
	if len(strictFailures) != 2 || !strictFailures["bytecode-strict-fail.lg"] || !strictFailures["gogen-strict-fail.lg"] {
		t.Fatalf("strict failures = %v; want both execution failures", strictFailures)
	}
	if len(diverged) != 1 || !diverged["output-divergence.lg"] {
		t.Fatalf("output divergences = %v; want only successful-but-different output", diverged)
	}
	for name := range strictFailures {
		if diverged[name] {
			t.Fatalf("strict failure %s was also classified as allowlistable output divergence", name)
		}
	}
}

func TestStrictFailureRemainsFatalWhenAllowlisted(t *testing.T) {
	results := []ParityResult{{
		Fixture: "strict-fail.lg", BCOutput: "<<runtime-error>>", BCOK: false,
		AOTOutput: "<<runtime-error>>", AOTOK: false,
	}}
	allowlisted := map[string]bool{"strict-fail.lg": true}

	errs := validateParityResults(results, allowlisted, allowlisted)
	foundStrictFailure := false
	for _, err := range errs {
		if strings.Contains(err, "STRICT-AOT failure") {
			foundStrictFailure = true
		}
	}
	if !foundStrictFailure {
		t.Fatalf("allowlisted strict failure produced errors %v; want unconditional STRICT-AOT failure", errs)
	}
	if len(errs) != 1 {
		t.Fatalf("strict failure was also processed as an allowlist divergence: %v", errs)
	}
}

func TestParityRederiveWritesBothAllowlists(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "test"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(parityRederiveEnv, "1")

	if !rederiveParityXfails(t, root, map[string]bool{"z.lg": true, "a.lg": true}) {
		t.Fatal("rederive mode was not activated")
	}
	want := "# Shrink-only allowlist of Tier-1 fixtures that currently DIVERGE\n" +
		"# between bytecode and -tags gogen_ir engines.\n" +
		"# Managed by TestParityGatePhase1: re-seed with `LETGO_PARITY_REDERIVE=1 go test -short=false -run TestParityGatePhase1 ./test/e2e`.\n" +
		"# A new divergence fails CI; a listed fixture that now agrees also fails (shrink-only).\n" +
		"# When committing an allowlist update, include the specific issue/fix in the reason comment.\n\n" +
		"a.lg\n" +
		"z.lg\n"
	for _, name := range []string{"parity-xfail-bytecode.txt", "parity-xfail-gogen.txt"} {
		path := filepath.Join(root, "test", name)
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Fatalf("%s bytes = %q; want byte-exact sorted output %q", name, got, want)
		}
	}
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
	Fixture   string
	BCOutput  string
	BCOK      bool
	AOTOutput string
	AOTOK     bool
}

const parityRederiveEnv = "LETGO_PARITY_REDERIVE"

func classifyParityResults(results []ParityResult) (strictFailures, diverged map[string]bool) {
	strictFailures = map[string]bool{}
	diverged = map[string]bool{}
	for _, result := range results {
		if !result.BCOK || !result.AOTOK {
			strictFailures[result.Fixture] = true
			continue
		}
		if result.BCOutput != result.AOTOutput {
			diverged[result.Fixture] = true
		}
	}
	return strictFailures, diverged
}

func validateParityResults(results []ParityResult, bcXfail, aotXfail map[string]bool) []string {
	strictFailures, diverged := classifyParityResults(results)
	var errs []string
	for _, result := range results {
		if strictFailures[result.Fixture] {
			errs = append(errs, fmt.Sprintf(
				"STRICT-AOT failure for Tier-1 fixture %s: bytecode_ok=%v gogen_ir_ok=%v; strict failures cannot be allowlisted",
				result.Fixture, result.BCOK, result.AOTOK))
		}
	}
	for name := range diverged {
		if !bcXfail[name] && !aotXfail[name] {
			errs = append(errs, fmt.Sprintf("NEW backend parity output divergence: %s (not in allowlists)", name))
		}
	}
	for name := range bcXfail {
		if !strictFailures[name] && !diverged[name] {
			errs = append(errs, fmt.Sprintf("%s is in test/parity-xfail-bytecode.txt but now AGREES — remove it (shrink-only ratchet)", name))
		}
	}
	for name := range aotXfail {
		if !strictFailures[name] && !diverged[name] {
			errs = append(errs, fmt.Sprintf("%s is in test/parity-xfail-gogen.txt but now AGREES — remove it (shrink-only ratchet)", name))
		}
	}
	return errs
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

	// Run tier-1 fixtures on both backends.
	results := make([]ParityResult, 0, len(tier1Fixtures))
	for _, name := range tier1Fixtures {
		wrappedPath := wrappedFixtures[name]
		bcOut, bcOk := runFixture(bc, wrappedPath)
		aotOut, aotOk := runFixture(aot, wrappedPath)
		results = append(results, ParityResult{
			Fixture: name, BCOutput: bcOut, BCOK: bcOk, AOTOutput: aotOut, AOTOK: aotOk,
		})
	}

	strictFailures, diverged := classifyParityResults(results)
	for _, result := range results {
		switch {
		case strictFailures[result.Fixture]:
			t.Logf("STRICT-AOT FAILURE %s: bytecode_ok=%v gogen_ir_ok=%v",
				result.Fixture, result.BCOK, result.AOTOK)
		case diverged[result.Fixture]:
			t.Logf("OUTPUT DIVERGENCE %s: bytecode=%q gogen_ir=%q",
				result.Fixture, result.BCOutput, result.AOTOutput)
		default:
			t.Logf("AGREE %s: output=%q", result.Fixture, result.BCOutput)
		}
	}

	rederived := rederiveParityXfails(t, root, diverged)
	var bcXfail, aotXfail map[string]bool
	if rederived {
		t.Logf("re-seeded parity allowlists with %d output divergence(s)", len(diverged))
		bcXfail, aotXfail = diverged, diverged
	} else {
		bcXfail = readParityXfail(t, filepath.Join(root, "test/parity-xfail-bytecode.txt"))
		aotXfail = readParityXfail(t, filepath.Join(root, "test/parity-xfail-gogen.txt"))
	}

	for _, err := range validateParityResults(results, bcXfail, aotXfail) {
		t.Error(err)
	}

	t.Logf("\n=== PHASE-1 GATE SUMMARY ===")
	t.Logf("Tier-1 fixtures tested: %d", len(tier1Fixtures))
	t.Logf("Divergences detected:   %d", len(diverged))
	t.Logf("Strict-AOT failures:    %d", len(strictFailures))
	t.Logf("Bytecode allowlist:     %d entries", len(bcXfail))
	t.Logf("Gogen_ir allowlist:     %d entries", len(aotXfail))
}

func rederiveParityXfails(t *testing.T, root string, diverged map[string]bool) bool {
	t.Helper()
	if os.Getenv(parityRederiveEnv) != "1" {
		return false
	}
	writeParityXfail(t, filepath.Join(root, "test/parity-xfail-bytecode.txt"), diverged)
	writeParityXfail(t, filepath.Join(root, "test/parity-xfail-gogen.txt"), diverged)
	return true
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
	b.WriteString("# Managed by TestParityGatePhase1: re-seed with `LETGO_PARITY_REDERIVE=1 go test -short=false -run TestParityGatePhase1 ./test/e2e`.\n")
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
