/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// fnNSLoader adapts a func to NSLoader for tests.
type fnNSLoader func(name string) *vm.Namespace

func (f fnNSLoader) Load(name string) *vm.Namespace { return f(name) }

// TestBytecodeNSLoaderReplayFailure covers the decode-then-failed-replay
// shape: the namespace exists as a placeholder (pre-registered during VarRef
// decoding) and is marked needs-load, and then its bundled chunk fails to
// replay. The failure must surface through RequireNS — the placeholder must
// not be cached as a successfully loaded namespace — and must keep surfacing
// on a retry rather than silently "succeeding" the second time.
func TestBytecodeNSLoaderReplayFailure(t *testing.T) {
	const name = "bytecode-nsload-fixture"

	prevLoader := GetNSLoader()
	prevChunks := precompiledCoreNS
	defer func() {
		SetNSLoader(prevLoader)
		precompiledCoreNS = prevChunks
		RemoveNS(name)
	}()

	// A chunk whose replay fails: a single garbage opcode.
	bad := vm.NewCodeChunk(vm.NewConsts())
	bad.Append(1 << 30)
	precompiledCoreNS = map[string]*vm.CodeChunk{name: bad}

	// Simulate bundle decode: placeholder namespace + needs-load marker.
	DefNSBare(name)
	MarkNSNeedsLoad(name)
	UseBytecodeNSLoader()

	if ns, err := RequireNS(name); err == nil {
		t.Fatalf("want error from RequireNS after failed replay, got ns %v", ns)
	}
	// The failure must not have consumed the marker: requiring again retries
	// the load (and fails loudly again) instead of returning a half-loaded ns.
	if ns, err := RequireNS(name); err == nil {
		t.Fatalf("want error from RequireNS on retry, got ns %v", ns)
	} else if !strings.Contains(err.Error(), name) {
		t.Fatalf("want error to name the namespace, got: %v", err)
	}
}

// TestRequireNSRetryAfterLoaderFailure pins the registry-level contract the
// bytecode loader relies on: a loader that fails (restoring the needs-load
// marker before returning nil) makes RequireNS report the failure, and a
// subsequent require reinvokes the loader so a now-working loader succeeds.
func TestRequireNSRetryAfterLoaderFailure(t *testing.T) {
	const name = "nsload-retry-fixture"

	prevLoader := GetNSLoader()
	defer func() {
		SetNSLoader(prevLoader)
		RemoveNS(name)
	}()

	DefNSBare(name)
	MarkNSNeedsLoad(name)

	calls := 0
	SetNSLoader(fnNSLoader(func(n string) *vm.Namespace {
		if n != name {
			return nil
		}
		calls++
		MarkNSNeedsLoad(n) // failed-load contract: restore the marker
		return nil
	}))

	if ns, err := RequireNS(name); err == nil {
		t.Fatalf("want error from RequireNS after loader failure, got ns %v", ns)
	}
	if calls != 1 {
		t.Fatalf("want 1 loader call, got %d", calls)
	}

	loaded := vm.NewNamespace(name)
	SetNSLoader(fnNSLoader(func(n string) *vm.Namespace {
		if n != name {
			return nil
		}
		calls++
		return loaded
	}))

	ns, err := RequireNS(name)
	if err != nil {
		t.Fatalf("want retry to succeed, got error: %v", err)
	}
	if ns != loaded {
		t.Fatalf("want the loader's namespace from the retry, got %v", ns)
	}
	if calls != 2 {
		t.Fatalf("want 2 loader calls total, got %d", calls)
	}
}

// TestRequireNSNativeNamespaceLoaderNoop pins the other half of the contract:
// a pre-registered NATIVE namespace marked needs-load (the gogen pattern —
// native fns installed at init, marker set so a loader may overlay more) must
// still require successfully when the loader has nothing to add and returns
// nil WITHOUT restoring the marker. The pre-existing namespace is the result.
func TestRequireNSNativeNamespaceLoaderNoop(t *testing.T) {
	const name = "nsload-native-fixture"

	prevLoader := GetNSLoader()
	defer func() {
		SetNSLoader(prevLoader)
		RemoveNS(name)
	}()

	native := DefNSBare(name)
	MarkNSNeedsLoad(name)

	SetNSLoader(fnNSLoader(func(n string) *vm.Namespace {
		return nil // nothing to load for any namespace; no marker restored
	}))

	ns, err := RequireNS(name)
	if err != nil {
		t.Fatalf("want native namespace to satisfy require, got error: %v", err)
	}
	if ns != native {
		t.Fatalf("want the pre-registered native namespace, got %v", ns)
	}
}
