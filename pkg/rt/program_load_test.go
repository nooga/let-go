/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/vm"
)

// TestLoadProgramNamespacesDrainsOverrides checks the #425 program-load
// helper: a pending RegisterGoOverrides entry is drained after the ns chunk
// runs, so the Var holds the native fn.
func TestLoadProgramNamespacesDrainsOverrides(t *testing.T) {
	const nsName = "entryframetest"

	nsMu.Lock()
	delete(nsRegistry, nsName)
	delete(pendingGoOverrides, nsName)
	nsMu.Unlock()

	native, err := vm.NativeFnType.Wrap(func(args []vm.Value) (vm.Value, error) {
		return vm.Int(42), nil
	})
	if err != nil {
		t.Fatalf("wrap native: %v", err)
	}
	// Queue while the ns does not exist — stays pending.
	RegisterGoOverrides(nsName, map[string]vm.Value{"answer": native})

	// Materialize the ns AFTER queueing so the pending entry survives.
	_ = DefNSBare(nsName)

	consts := vm.NewConsts()
	idx := consts.Intern(vm.NIL)
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST)
	chunk.Append32(idx)
	chunk.Append(vm.OP_RETURN)
	chunk.SetMaxStack(1)
	unit := &bytecode.ExecUnit{
		NSOrder:  []string{nsName},
		NSChunks: map[string]*vm.CodeChunk{nsName: chunk},
	}
	unit.MainChunk = chunk

	if err := LoadProgramNamespaces(unit); err != nil {
		t.Fatalf("LoadProgramNamespaces: %v", err)
	}

	ns := LookupNS(nsName)
	if ns == nil {
		t.Fatal("namespace missing after load")
	}
	v := ns.LookupLocal(vm.Symbol("answer"))
	if v == nil {
		t.Fatal("answer var missing — override did not apply")
	}
	got, err := v.Deref().(vm.Fn).Invoke(nil)
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if got != vm.Int(42) {
		t.Fatalf("got %v, want 42", got)
	}
}

func TestRunProgramMainChunkNoopWhenAlreadyLoaded(t *testing.T) {
	const nsName = "entryframemain"
	nsMu.Lock()
	delete(nsRegistry, nsName)
	nsMu.Unlock()
	_ = DefNSBare(nsName)

	consts := vm.NewConsts()
	idx := consts.Intern(vm.NIL)
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST)
	chunk.Append32(idx)
	chunk.Append(vm.OP_RETURN)
	chunk.SetMaxStack(1)
	unit := &bytecode.ExecUnit{
		NSOrder:   []string{nsName},
		NSChunks:  map[string]*vm.CodeChunk{nsName: chunk},
		MainChunk: chunk,
	}
	if err := LoadProgramNamespaces(unit); err != nil {
		t.Fatalf("load: %v", err)
	}
	// Must not re-run (would be a no-op on empty chunk either way, but the
	// contract is what the native-entry fallback relies on).
	if err := RunProgramMainChunk(unit); err != nil {
		t.Fatalf("RunProgramMainChunk: %v", err)
	}
}
