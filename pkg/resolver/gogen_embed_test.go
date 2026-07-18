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

// TestGogenResolvesFromEmbeddedSource locks in nooga/let-go#425 (self-contain
// gogen): shipped binaries embed gogen.lg (pkg/rt/gogen_src.go) and register it
// as an auxiliary embedded source, so (require 'gogen) resolves with NO external
// source path — EmbeddedSource returns it before any search path is consulted.
// Before #425 gogen was an external classpath dir (deps.edn), so compiling from
// the wrong cwd failed with `Can't resolve gogen/ns->go-pkg`. The empty
// search-path slice below is the point: only embedded resolution can satisfy the
// require.
//
// Gated //go:build !bootstrap: the embed is deliberately excluded from the
// self-hosting bootstrap build (gogen resolves via classpath there — see
// pkg/rt/gogen_src.go), so this shipped-binary behavior is only asserted for
// non-bootstrap builds.
func TestGogenResolvesFromEmbeddedSource(t *testing.T) {
	if _, ok := rt.EmbeddedSource("gogen"); !ok {
		t.Fatal("rt.EmbeddedSource(\"gogen\") not found — gogen must be embedded via pkg/rt/gogen_src.go (#425)")
	}

	prev := rt.GetNSLoader()
	defer rt.SetNSLoader(prev)

	consts := vm.NewConsts()
	ctx := compiler.NewCompiler(consts, rt.NS("user"))
	rt.SetNSLoader(NewNSResolver(ctx, []string{})) // no search paths: embed-only
	ctx.SetSource("<test>")

	// Reference gogen/ns->go-pkg — a macro-layer .lg defn, not a native fn — so
	// success proves the embedded *source* loaded, not merely the native ns
	// installed at boot. If gogen's source didn't resolve, the compile fails with
	// the pre-#425 "Can't resolve gogen/ns->go-pkg".
	_, _, err := ctx.CompileMultiple(strings.NewReader("(require 'gogen) (gogen/ns->go-pkg \"ir.lower-go\")"))
	if err != nil {
		t.Fatalf("resolving gogen from embedded source failed (regressed #425?): %v", err)
	}
}
