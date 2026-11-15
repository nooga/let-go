//go:build bootstrap

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"path/filepath"
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
