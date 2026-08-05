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
	"testing"

	"github.com/nooga/let-go/pkg/genmanifest"
)

// Strict-mode audit harness.
//
// Classifies each test/gold-aot/*.lg fixture under *ir-compile-strict* binding:
// - Tier 1: Compiles and runs successfully under strict mode (genuine AOT)
// - Tier 2: Compiles default, fails strict (silent trampoline fallback; excluded from coverage)
// - Tier 3: Diverges bytecode ≠ gogen output (goes to per-backend allowlist)
//
// Each fixture is wrapped in a strict-mode binding:
//   (binding [*ir-compile* true *ir-compile-strict* true]
//     (eval (read-string (str "(do " (slurp "fixture.lg") ")"))))
//
// The wrapper is created in a temp directory before running.
// A throw from strict mode is classified as Tier 2 (cannot currently AOT-lower).
// Divergence between backends is classified as Tier 3 (miscompile).
// Agreement under strict is classified as Tier 1 (genuine AOT).

type AuditResult struct {
	Fixture   string
	Tier      int    // 1, 2, or 3
	Reason    string // explanation
	BCOutput  string // bytecode output (or error token)
	AOTOutput string // gogen_ir output (or error token)
}

func TestStrictModeAudit(t *testing.T) {
	if testing.Short() {
		t.Skip("strict-mode audit builds let-go twice; run via `make strict-audit`")
	}

	root := repoRoot(t)

	if err := genmanifest.CheckTreeManifest(filepath.Join(root, "pkg/rt/core_go_lowered")); err != nil {
		t.Fatalf("generated lowered tree is not valid: %v\nrun `make generate` and retry", err)
	}

	bc := buildLGTags(t, "")
	aot := buildLGTags(t, "gogen_ir")
	goldAOTDir := filepath.Join(root, "test/gold-aot")

	fixtures, err := filepath.Glob(filepath.Join(goldAOTDir, "*.lg"))
	if err != nil {
		t.Fatalf("glob %s: %v", goldAOTDir, err)
	}
	if len(fixtures) == 0 {
		t.Fatalf("no fixtures in %s", goldAOTDir)
	}
	sort.Strings(fixtures)

	// Create strict-mode wrappers in a temp directory
	wrapDir := t.TempDir()
	wrappedFixtures := make(map[string]string) // name -> wrapped path
	for _, f := range fixtures {
		name := filepath.Base(f)
		wrappedPath := filepath.Join(wrapDir, name)
		if err := createStrictWrapper(wrappedPath, f); err != nil {
			t.Fatalf("create strict wrapper for %s: %v", name, err)
		}
		wrappedFixtures[name] = wrappedPath
	}

	// Run audit and collect results
	var results []AuditResult
	for _, f := range fixtures {
		name := filepath.Base(f)
		wrappedPath := wrappedFixtures[name]
		result := auditFixture(t, bc, aot, name, wrappedPath)
		results = append(results, result)
	}

	// Print audit summary
	tier1 := 0
	tier2 := 0
	tier3 := 0
	for _, r := range results {
		switch r.Tier {
		case 1:
			tier1++
		case 2:
			tier2++
		case 3:
			tier3++
		}
		status := fmt.Sprintf("Tier %d", r.Tier)
		t.Logf("%-30s %s: %s", r.Fixture, status, r.Reason)
	}

	t.Logf("\n=== AUDIT SUMMARY ===")
	t.Logf("Tier 1 (genuine AOT):      %2d fixtures", tier1)
	t.Logf("Tier 2 (strict-fail tramp): %2d fixtures", tier2)
	t.Logf("Tier 3 (divergent):        %2d fixtures", tier3)
	t.Logf("TOTAL:                     %2d fixtures", len(results))

	// Print detailed evidence for each fixture
	t.Logf("\n=== DETAILED EVIDENCE ===")
	for _, r := range results {
		t.Logf("\n%s (Tier %d):", r.Fixture, r.Tier)
		t.Logf("  Reason: %s", r.Reason)
		if r.Tier == 3 {
			t.Logf("  Bytecode: %q", r.BCOutput)
			t.Logf("  Gogen_IR: %q", r.AOTOutput)
		}
	}
}

func createStrictWrapper(dest, source string) error {
	// Read the source fixture
	content, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("read %s: %v", source, err)
	}

	// Create wrapper that evaluates the fixture content within strict mode binding
	wrapper := fmt.Sprintf(`(binding [*ir-compile* true *ir-compile-strict* true]
  (eval (read-string
    (str "(do " %q ")"))))
`, string(content))

	if err := os.WriteFile(dest, []byte(wrapper), 0o644); err != nil {
		return fmt.Errorf("write %s: %v", dest, err)
	}
	return nil
}

func auditFixture(t *testing.T, bc, aot, name, wrappedPath string) AuditResult {
	t.Helper()

	bcOut, bcOk := runFixture(bc, wrappedPath)
	aotOut, aotOk := runFixture(aot, wrappedPath)

	// Classify based on results
	result := AuditResult{
		Fixture:   name,
		BCOutput:  bcOut,
		AOTOutput: aotOut,
	}

	// Bytecode strict-mode failure: Tier 2 (trampoline fallback)
	if !bcOk {
		result.Tier = 2
		result.Reason = "bytecode fails strict mode (would trampoline); excluded from coverage"
		return result
	}

	// Divergence (one side fails or outputs differ): Tier 3
	if !aotOk || bcOut != aotOut {
		result.Tier = 3
		if !aotOk {
			result.Reason = "gogen_ir fails strict mode; bytecode succeeds (miscompile)"
		} else {
			result.Reason = fmt.Sprintf("bytecode ≠ gogen_ir: %q ≠ %q", bcOut, aotOut)
		}
		return result
	}

	// Both succeed and agree: Tier 1 (genuine AOT)
	result.Tier = 1
	result.Reason = "both backends pass strict mode and agree"
	return result
}
