/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package wasm

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GoBuildArgs builds the `go build` argument list for compiling the wasm binary
// at wasmPath, optionally with the given comma/space-separated build tags.
func GoBuildArgs(wasmPath string, buildTags string) []string {
	args := []string{"build"}
	if buildTags != "" {
		args = append(args, "-tags", buildTags)
	}
	args = append(args, "-o", wasmPath, ".")
	return args
}

// HasBuildTag reports whether tag appears in rawTags (a `,`/space/tab/newline
// separated build-tag list).
func HasBuildTag(rawTags, tag string) bool {
	for _, field := range strings.FieldsFunc(rawTags, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n'
	}) {
		if field == tag {
			return true
		}
	}
	return false
}

// WriteGogenIRWireup writes tmpDir/main_gogen_ir.go — a gogen_ir-tagged file that
// blank-imports every lowered-Go package under srcDir/pkg/rt/core_go_lowered so
// the wasm build links the native lowerings.
func WriteGogenIRWireup(tmpDir, srcDir string) error {
	pkgs, err := ListGogenIRPackages(srcDir)
	if err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("// Code generated for wasm gogen_ir wireup. DO NOT EDIT.\n\n")
	b.WriteString("//go:build gogen_ir\n\n")
	b.WriteString("package main\n")
	if len(pkgs) > 0 {
		b.WriteString("\nimport (\n")
		for _, pkg := range pkgs {
			fmt.Fprintf(&b, "\t_ %q\n", "github.com/nooga/let-go/pkg/rt/core_go_lowered/"+pkg)
		}
		b.WriteString(")\n")
	}
	return os.WriteFile(filepath.Join(tmpDir, "main_gogen_ir.go"), []byte(b.String()), 0644)
}

// ListGogenIRPackages returns the slash-separated relative paths of every
// lowered-Go package (a directory containing at least one .go file) under
// srcDir/pkg/rt/core_go_lowered, sorted.
func ListGogenIRPackages(srcDir string) ([]string, error) {
	root := filepath.Join(srcDir, "pkg", "rt", "core_go_lowered")
	var pkgs []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if path == root {
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		hasGo := false
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
				hasGo = true
				break
			}
		}
		if !hasGo {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		pkgs = append(pkgs, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(pkgs)
	return pkgs, nil
}
