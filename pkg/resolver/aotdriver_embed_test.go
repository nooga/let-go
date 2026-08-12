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

// TestAotDriverResolvesFromEmbeddedSource locks in nooga/let-go#596: the AOT
// compile driver ships inside the binary (pkg/rt/aotdriver_src.go) as an
// auxiliary embedded source, so it resolves with NO external source path and no
// let-go checkout on disk. Before this, driving lowering meant executing
// scripts/lg-compile from a checkout — a released binary had no copy of it.
//
// The empty search-path slice is the point: only embedded resolution can
// satisfy the require. Sibling of TestGogenResolvesFromEmbeddedSource, and
// gated the same way — the embed is excluded from the self-hosting bootstrap
// build, so this shipped-binary behavior is asserted for non-bootstrap only.
func TestAotDriverResolvesFromEmbeddedSource(t *testing.T) {
	if _, ok := rt.EmbeddedSource("lg.aotdriver"); !ok {
		t.Fatal("rt.EmbeddedSource(\"lg.aotdriver\") not found — the driver must be embedded via pkg/rt/aotdriver_src.go (#596)")
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
		"(require 'lg.aotdriver) (lg.aotdriver/basename-stem \"a/b/fib.lg\")"))
	if err != nil {
		t.Fatalf("resolving lg.aotdriver from embedded source failed: %v", err)
	}
	got, ok := val.(vm.String)
	if !ok || string(got) != "fib" {
		t.Fatalf("basename-stem returned %#v, want vm.String(\"fib\")", val)
	}
}
