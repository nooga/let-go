//go:build !bootstrap

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package resolver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// An out-of-tree package registers its lg layer with rt.RegisterEmbeddedSource
// (nooga/let-go#1063). The namespace resolves with no search path, and shadows
// a file of the same name on the search path, because embedded resolution comes
// first.
func TestRegisteredEmbeddedSourceResolves(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "test"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Each case gets its own namespace: a loaded namespace stays loaded, and a
	// second require of it would not consult the resolver again.
	for _, tc := range []struct {
		desc, ns string
		paths    []string
	}{
		{"no search path", "test.registered-lib", []string{}},
		{"shadowing a file on the search path", "test.registered-shadow", []string{dir}},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			rt.RegisterEmbeddedSource(tc.ns, "(ns "+tc.ns+`) (defn origin [] "embedded")`)
			file := strings.ReplaceAll(strings.TrimPrefix(tc.ns, "test."), "-", "_") + ".lg"
			if err := os.WriteFile(filepath.Join(dir, "test", file),
				[]byte("(ns "+tc.ns+`) (defn origin [] "search-path")`), 0o644); err != nil {
				t.Fatal(err)
			}

			prev := rt.GetNSLoader()
			defer rt.SetNSLoader(prev)

			consts := vm.NewConsts()
			ctx := compiler.NewCompiler(consts, rt.NS("user"))
			rt.SetNSLoader(NewNSResolver(ctx, tc.paths))
			ctx.SetSource("<test>")

			_, v, err := ctx.CompileMultiple(strings.NewReader(
				"(require '" + tc.ns + ") (" + tc.ns + "/origin)"))
			if err != nil {
				t.Fatalf("require: %v", err)
			}
			if got := v.String(); !strings.Contains(got, "embedded") {
				t.Fatalf("origin = %s, want the registered source's \"embedded\"", got)
			}
		})
	}
}
