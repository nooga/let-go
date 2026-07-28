/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Package primgen scans //lg:-annotated Go sources and emits the
// zz_primitives_generated.go registrar. It is deliberately a PURE source →
// source code generator: it depends only on go/ast + text templates and does
// NOT import the let-go runtime (pkg/compiler / pkg/rt).
//
// This independence is the whole point. The registrar it emits is what lets the
// NEXT-generation runtime boot; a generator that booted the current runtime
// would depend on the very artifacts it is about to produce (p0 depending on
// p1). During a large primitive migration — closures moved out of installLangNS
// but not yet in the generated registrar — that runtime cannot boot, so a
// runtime-coupled generator would deadlock the regeneration it exists to drive.
package primgen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

// KebabCase converts a Go identifier (CamelCase) to a let-go kebab name. Used
// both to default a primitive's lg-name during scanning and by lginterop's
// skeleton emitter, so it lives in the shared runtime-free package.
func KebabCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 {
			prev := rune(s[i-1])
			if unicode.IsUpper(r) {
				if unicode.IsLower(prev) {
					b.WriteByte('-')
				} else if i+1 < len(s) && unicode.IsLower(rune(s[i+1])) {
					b.WriteByte('-')
				}
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// Generate scans srcDir for //lg:native directives and writes the generated
// primitive registrar to outPath. goPkg is the Go import path of the scanned
// sources (used for cross-package symbol references); empty falls back to the
// per-spec GoPkg or the builtins default.
func Generate(srcDir, outPath, goPkg string) error {
	var allSpecs []primSpec

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("read directory %s: %w", srcDir, err)
	}

	for _, entry := range entries {
		// Skip test files and generated output: an //lg:native annotation in
		// a _test.go or zz_ file would emit an unconditional reference to a
		// symbol the normal build doesn't compile.
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") ||
			strings.HasSuffix(entry.Name(), "_test.go") ||
			strings.HasPrefix(entry.Name(), "zz_") {
			continue
		}

		fpath := filepath.Join(srcDir, entry.Name())
		src, err := os.ReadFile(fpath)
		if err != nil {
			return fmt.Errorf("read %s: %w", fpath, err)
		}

		// Files under a build constraint may reference symbols absent from
		// other targets; the registrar compiles unconditionally, so skip them.
		if hasBuildConstraint(src) {
			continue
		}

		specs, err := scanSource(fpath, src)
		if err != nil {
			return fmt.Errorf("parse %s: %w", fpath, err)
		}
		allSpecs = append(allSpecs, specs...)
	}

	if len(allSpecs) == 0 {
		fmt.Printf("lgprimgen: no //lg:native directives found in %s; generating empty stub\n", srcDir)
		// Fall through to generate an empty RegisterGeneratedPrimitives() stub
	}

	for i := range allSpecs {
		if goPkg != "" {
			allSpecs[i].GoPkg = goPkg
		} else if allSpecs[i].GoPkg == "" {
			allSpecs[i].GoPkg = "github.com/nooga/let-go/pkg/rt/builtins"
		}
	}

	output := emitFile(allSpecs)

	formatted, err := gofmtCode(output)
	if err != nil {
		return fmt.Errorf("gofmt: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(outPath), err)
	}
	if err := os.WriteFile(outPath, []byte(formatted), 0644); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}

	fmt.Printf("lgprimgen: generated primitives → %s (%d specs)\n", outPath, len(allSpecs))
	return nil
}

// ScanNativeNames returns the set of let-go names already claimed by //lg:native
// declarations in dir's Go sources, skipping excludeBase (a file basename, e.g.
// "lang.go"). The hoist codemod uses it to avoid hoisting a closure whose
// lg-name a pre-existing native decl already owns — doing so would emit a
// duplicate registration and fail to compile the generated registrar. Reuses
// the same scanner lgprimgen runs, so both agree on what a name is.
func ScanNativeNames(dir, excludeBase string) (map[string]bool, error) {
	names := map[string]bool{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") ||
			strings.HasSuffix(entry.Name(), "_test.go") ||
			strings.HasPrefix(entry.Name(), "zz_") ||
			entry.Name() == excludeBase {
			continue
		}
		fpath := filepath.Join(dir, entry.Name())
		src, err := os.ReadFile(fpath)
		if err != nil {
			return nil, err
		}
		if hasBuildConstraint(src) {
			continue
		}
		specs, err := scanSource(fpath, src)
		if err != nil {
			return nil, err
		}
		for _, s := range specs {
			names[s.LgName] = true
		}
	}
	return names, nil
}

// gofmtCode formats Go source code using the gofmt command. On failure it
// surfaces gofmt's stderr (which carries the <line>:<col>: syntax-error
// location) and dumps the unformatted source to a temp file, so a bad emitted
// registrar is diagnosable instead of a bare "exit status 2".
func gofmtCode(src string) (string, error) {
	cmd := exec.Command("gofmt")
	cmd.Stdin = strings.NewReader(src)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		dump := filepath.Join(os.TempDir(), "lgprimgen-unformatted.go")
		_ = os.WriteFile(dump, []byte(src), 0644)
		return "", fmt.Errorf("gofmt failed: %w\n%s(unformatted source written to %s)", err, stderr.String(), dump)
	}
	return string(output), nil
}
