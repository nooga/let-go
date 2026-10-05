/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// precompiledChunkLoader is a bundle-only namespace loader over
// PrecompiledNSChunk, the seam the resolver loads bundled namespaces through.
type precompiledChunkLoader struct{ t *testing.T }

func (l precompiledChunkLoader) Load(name string) *vm.Namespace {
	chunk, err := PrecompiledNSChunk(name)
	if err != nil {
		l.t.Errorf("PrecompiledNSChunk(%s): %v", name, err)
		return nil
	}
	if chunk == nil {
		return nil
	}
	saved := rt.CurrentNS.Deref()
	defer rt.CurrentNS.SetRoot(saved)
	f := vm.NewFrame(chunk, nil)
	_, err = f.RunProtected()
	vm.ReleaseFrame(f)
	if err != nil {
		l.t.Errorf("replaying %s: %v", name, err)
		return nil
	}
	return rt.LookupNS(name)
}

// countDecode runs fn with decode statistics enabled and returns what it
// decoded.
func countDecode(t *testing.T, fn func()) bytecode.DecodeStats {
	t.Helper()
	bytecode.ResetDecodeStats()
	bytecode.SetDecodeStatsEnabled(true)
	defer bytecode.SetDecodeStatsEnabled(false)
	fn()
	return bytecode.SnapshotDecodeStats()
}

// TestInitFromLGBDefersLazyNamespaceChunks pins that startup decodes only the
// chunks it runs: a namespace LoadCoreBundle leaves needs-load (edn here) has
// its code decoded when it is first required, not by every process at boot.
// Otherwise code added to a lazily loaded namespace would grow the startup
// gate (BenchmarkInitFromLGB) although no process runs it before a require.
func TestInitFromLGBDefersLazyNamespaceChunks(t *testing.T) {
	if len(rt.CoreCompiledLGB) == 0 {
		t.Skip("no precompiled core_compiled.lgb")
	}
	const lazy = "edn"

	// Everything the bundle holds, as a plain eager decode sees it.
	full := countDecode(t, func() {
		if _, err := bytecode.DecodeToExecUnitBytes(rt.CoreCompiledLGB, rt.LGBVarResolver); err != nil {
			t.Fatal(err)
		}
	})
	if full.Chunks == 0 {
		t.Fatal("eager decode reported no chunks; decode stats are not wired")
	}

	// A fresh boot must decode strictly less than the bundle: at least the
	// lazily loaded namespace's chunks stay undecoded.
	rt.ClearNSNeedsLoad(lazy)
	startup := countDecode(t, func() { initFromLGB(t) })
	t.Logf("bundle chunks=%d words=%d; startup chunks=%d words=%d", full.Chunks, full.CodeWords, startup.Chunks, startup.CodeWords)
	if !rt.NSNeedsLoad(lazy) {
		t.Fatalf("%s is not marked needs-load after init", lazy)
	}
	if startup.Chunks >= full.Chunks || startup.CodeWords >= full.CodeWords {
		t.Fatalf("startup decoded the whole bundle (%d/%d chunks, %d/%d code words): a lazily loaded namespace's chunks are decoded at boot", startup.Chunks, full.Chunks, startup.CodeWords, full.CodeWords)
	}

	// Requiring the namespace materializes its chunks, exactly once. The
	// compiler package installs no loader of its own; this one replays the
	// precompiled chunk the way resolver.execPrecompiled does.
	prevLoader := rt.GetNSLoader()
	rt.SetNSLoader(precompiledChunkLoader{t})
	defer rt.SetNSLoader(prevLoader)
	first := countDecode(t, func() {
		if ns := rt.NS(lazy); ns == nil {
			t.Fatalf("require %s failed", lazy)
		}
	})
	if first.Chunks == 0 {
		t.Fatalf("requiring %s decoded no chunks; it was decoded at boot after all", lazy)
	}
	if !lazyDefBound(lazy, "pretty") {
		t.Fatalf("%s/pretty is unbound after require", lazy)
	}
	again := countDecode(t, func() { rt.NS(lazy) })
	if again.Chunks != 0 {
		t.Fatalf("a second require of %s decoded %d chunks again", lazy, again.Chunks)
	}
	if startup.Chunks+first.Chunks > full.Chunks {
		t.Fatalf("startup (%d) + %s (%d) decoded more chunks than the bundle holds (%d)", startup.Chunks, lazy, first.Chunks, full.Chunks)
	}
}
