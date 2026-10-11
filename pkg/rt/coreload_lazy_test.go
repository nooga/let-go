//go:build !bootstrap

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// TestLoadCoreDefersUntilBytecodeLoaderRequires covers the runtime-only boot
// (LoadCore + bytecodeNSLoader, cmd/lg-runtime): a namespace left needs-load
// stays undecoded after boot, and its first require decodes and replays it
// through the bytecode loader.
func TestLoadCoreDefersUntilBytecodeLoaderRequires(t *testing.T) {
	const lazy = "walk"
	prevLoader, prevCore := GetNSLoader(), precompiledCore
	defer func() {
		SetNSLoader(prevLoader)
		precompiledCore = prevCore
	}()

	if err := LoadCore(); err != nil {
		t.Fatalf("LoadCore: %v", err)
	}
	if !precompiledCore.DeferredNS(lazy) {
		t.Fatalf("%s was decoded at boot", lazy)
	}
	if !NSNeedsLoad(lazy) {
		t.Fatalf("%s is not marked needs-load", lazy)
	}
	// LoadCore leaves hybrids to the loader too, so they are deferred as well.
	if !precompiledCore.DeferredNS("string") {
		t.Fatal("hybrid string was decoded at boot although LoadCore leaves it to the loader")
	}

	UseBytecodeNSLoader()
	ns, err := RequireNS(lazy)
	if err != nil {
		t.Fatalf("RequireNS(%s): %v", lazy, err)
	}
	if precompiledCore.DeferredNS(lazy) {
		t.Fatalf("%s still deferred after require", lazy)
	}
	v := ns.LookupLocal(vm.Symbol("postwalk"))
	if v == nil || !v.IsBound() {
		t.Fatalf("%s/postwalk unbound after require", lazy)
	}
	// (walk/postwalk identity [1 2]) round-trips: the materialized functions run.
	out, err := v.Invoke([]vm.Value{LookupCoreVar("identity").Deref(), vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2)})})
	if err != nil {
		t.Fatalf("postwalk: %v", err)
	}
	if got := out.String(); got != "[1 2]" {
		t.Fatalf("postwalk identity [1 2] = %s", got)
	}
}

// TestLoadCoreBundleDecodedCodeNeverLoadsDeferredSlot is the runtime half of
// the encoder's self-containment proof: after a deferring decode, no
// instruction of any decoded chunk names a const-pool slot that is still
// deferred. The encoder proved it from the code it emitted; this checks it
// against the bundle a process actually boots from.
func TestLoadCoreBundleDecodedCodeNeverLoadsDeferredSlot(t *testing.T) {
	unit, err := LoadCoreBundle(CoreLoadOptions{EagerHybrids: true})
	if err != nil {
		t.Fatal(err)
	}
	deferred := 0
	for _, name := range unit.NSOrder {
		if unit.DeferredNS(name) {
			deferred++
		}
	}
	if deferred == 0 {
		t.Fatal("nothing deferred: the bundle carries no namespace ranges?")
	}
	pool := unit.Consts
	seen := map[*vm.CodeChunk]bool{}
	var check func(name string, c *vm.CodeChunk)
	check = func(name string, c *vm.CodeChunk) {
		if c == nil || seen[c] {
			return
		}
		seen[c] = true
		code := c.Code()
		for i := 0; i < len(code); i += vm.OpcodeStride(code[i]) {
			if !vm.ReadsConst(code[i]) || i+1 >= len(code) {
				continue
			}
			idx := int(code[i+1])
			if pool.IsDeferred(idx) {
				t.Fatalf("decoded chunk of %s loads deferred const %d at ip %d", name, idx, i)
			}
			if f, ok := pool.AllValues()[idx].(*vm.Func); ok {
				check(name, f.Chunk())
			}
		}
	}
	for name, ch := range unit.NSChunks {
		check(name, ch)
	}
	for _, v := range pool.AllValues() {
		if f, ok := v.(*vm.Func); ok {
			check(f.FuncName(), f.Chunk())
		}
	}
	t.Logf("checked %d decoded chunks with %d namespaces deferred", len(seen), deferred)
}
