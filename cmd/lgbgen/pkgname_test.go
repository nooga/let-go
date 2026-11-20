//go:build bootstrap

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The same table holds for gogen/ns->go-pkg and gogen/ns->go-reldir, which
// name the packages lg compile writes; test/e2e covers that side.
func TestNamespaceToGoPackageAndDir(t *testing.T) {
	for _, tc := range []struct{ ns, pkg, dir string }{
		{"edn", "edn", "edn"},
		{"ir.lower-go", "lower_go", "ir/lower_go"},
		{"ir.passes.pipeline", "pipeline", "ir/passes/pipeline"},
		{"app.main", "main_", "app/main_"},
		{"main", "main_", "main_"},
		{"app.go", "go_", "app/go_"},
		{"app.type", "type_", "app/type_"},
		{"app.init", "init_", "app/init_"},
		{"main.util", "util", "main/util"},
		{"go.run-all", "run_all", "go/run_all"},
	} {
		if got := nsToGoPkgName(tc.ns); got != tc.pkg {
			t.Errorf("nsToGoPkgName(%q) = %q, want %q", tc.ns, got, tc.pkg)
		}
		if got := nsToGoRelDir(tc.ns); got != filepath.FromSlash(tc.dir) {
			t.Errorf("nsToGoRelDir(%q) = %q, want %q", tc.ns, got, tc.dir)
		}
	}
}

// lgbgen names the packages of the lowered core and gogen.lg names the ones
// lg compile writes, each from its own copy of the unusable names. A name in
// only one of them would give the same namespace two package names.
func TestUnusableGoPkgNamesMatchGogen(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "pkg", "rt", "gogen", "gogen.lg"))
	if err != nil {
		t.Fatal(err)
	}
	set := regexp.MustCompile(`(?s)\(def \^:private unusable-go-pkg-names.*?#\{(.*?)\}`).FindSubmatch(src)
	if set == nil {
		t.Fatal("unusable-go-pkg-names not found in gogen.lg")
	}
	var lg []string
	for _, m := range regexp.MustCompile(`"([^"]+)"`).FindAllSubmatch(set[1], -1) {
		lg = append(lg, string(m[1]))
	}
	var gen []string
	for name := range unusableGoPkgNames {
		gen = append(gen, name)
	}
	sort.Strings(lg)
	sort.Strings(gen)
	if a, b := strings.Join(lg, " "), strings.Join(gen, " "); a != b {
		t.Errorf("unusable package names differ\n gogen.lg: %s\n lgbgen:   %s", a, b)
	}
}
