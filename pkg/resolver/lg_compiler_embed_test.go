//go:build !bootstrap

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package resolver

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TestLgCompilerResolvesFromEmbeddedSource locks in the property a released
// binary depends on: the AOT compile driver resolves with NO external source
// path and no let-go checkout on disk. The empty search-path slice is the
// point, since only embedded resolution can satisfy the require. lg.compiler is
// kept out of the bytecode bundle (cmd/lgbgen.isBundleSkippedTool), so
// resolution from embedded SOURCE on demand is the path that matters.
func TestLgCompilerResolvesFromEmbeddedSource(t *testing.T) {
	if _, ok := rt.EmbeddedSource("lg.compiler"); !ok {
		t.Fatal("rt.EmbeddedSource(\"lg.compiler\") not found — the driver must ship as core source at pkg/rt/core/lg/compiler.lg (#596)")
	}

	prev := rt.GetNSLoader()
	defer rt.SetNSLoader(prev)

	consts := vm.NewConsts()
	ctx := compiler.NewCompiler(consts, rt.NS("user"))
	rt.SetNSLoader(NewNSResolver(ctx, []string{})) // no search paths: embed-only
	ctx.SetSource("<test>")

	// Call a pure driver fn rather than merely requiring the ns: success proves
	// the embedded source compiled AND that its transitive requires (gogen,
	// ir.passes.pipeline, ir.passes.entry-frame) also resolved embed-only, which
	// is the property that makes a checkout-free `lg compile` possible.
	// basename-stem is chosen because it needs no filesystem or lowering state.
	_, val, err := ctx.CompileMultiple(strings.NewReader(
		"(require 'lg.compiler) (lg.compiler/basename-stem \"a/b/fib.lg\")"))
	if err != nil {
		t.Fatalf("resolving lg.compiler from embedded source failed: %v", err)
	}
	got, ok := val.(vm.String)
	if !ok || string(got) != "fib" {
		t.Fatalf("basename-stem returned %#v, want vm.String(\"fib\")", val)
	}
}
