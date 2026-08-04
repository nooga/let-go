/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"reflect"
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

// mkLogChunk builds a chunk that invokes logFn with a single label arg, so a
// replay leaves an observable, ordered trace of which chunks ran.
func mkLogChunk(logFn vm.Value, label vm.Value) *vm.CodeChunk {
	consts := vm.NewConsts()
	fnIdx := consts.Intern(logFn)
	argIdx := consts.Intern(label)
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST)
	chunk.Append32(fnIdx)
	chunk.Append(vm.OP_LOAD_CONST)
	chunk.Append32(argIdx)
	chunk.Append(vm.OP_INVOKE)
	chunk.Append32(1)
	chunk.Append(vm.OP_RETURN)
	chunk.SetMaxStack(2)
	return chunk
}

func valStrings(vs []vm.Value) []string {
	out := make([]string, len(vs))
	for i, v := range vs {
		out[i] = string(v.(vm.String))
	}
	return out
}

// TestInvokeProgramEntryUsesSelectedNamespace pins the #425 VM-fallback
// contract: the compiler-selected namespace owns the entry lookup. A competing
// -main in another namespace must not steal the call via CurrentNS or a
// registry sweep.
func TestInvokeProgramEntryUsesSelectedNamespace(t *testing.T) {
	const (
		nsWant  = "entryns.want"
		nsOther = "entryns.other"
	)
	for _, name := range []string{nsWant, nsOther} {
		nsMu.Lock()
		delete(nsRegistry, name)
		nsMu.Unlock()
		_ = DefNSBare(name)
	}

	var got string
	wantFn, err := vm.NativeFnType.Wrap(func(args []vm.Value) (vm.Value, error) {
		got = "want"
		return vm.NIL, nil
	})
	if err != nil {
		t.Fatalf("wrap want: %v", err)
	}
	otherFn, err := vm.NativeFnType.Wrap(func(args []vm.Value) (vm.Value, error) {
		got = "other"
		return vm.NIL, nil
	})
	if err != nil {
		t.Fatalf("wrap other: %v", err)
	}
	LookupNS(nsWant).Def("-main", wantFn)
	LookupNS(nsOther).Def("-main", otherFn)

	savedNS := CurrentNS.Deref()
	defer CurrentNS.SetRoot(savedNS)
	// Point CurrentNS at the competing namespace — the old global search would
	// prefer this. The selected-namespace path must ignore it.
	CurrentNS.SetRoot(LookupNS(nsOther))

	ec := vm.NewExecContext()
	got = ""
	if err := InvokeProgramEntry(ec, nsWant, "-main", nil); err != nil {
		t.Fatalf("InvokeProgramEntry: %v", err)
	}
	if got != "want" {
		t.Fatalf("invoked %q, want selected namespace's -main", got)
	}

	if err := InvokeProgramEntry(ec, "", "-main", nil); err == nil {
		t.Fatal("empty namespace: expected error")
	}
	if err := InvokeProgramEntry(ec, "entryns.missing", "-main", nil); err == nil {
		t.Fatal("missing namespace: expected error")
	}
}

// TestRunExecUnitReplayOrderAndMainOnce pins the observable contract of
// RunExecUnit after the #425 refactor onto LoadProgramNamespaces +
// RunProgramMainChunk: every NSOrder chunk runs in order, and the main chunk
// runs exactly once — whether it is a distinct chunk (runs last) or aliases a
// namespace chunk (runs in place, not a second time).
func TestRunExecUnitReplayOrderAndMainOnce(t *testing.T) {
	var log []vm.Value
	logFn, err := vm.NativeFnType.Wrap(func(args []vm.Value) (vm.Value, error) {
		log = append(log, args[0])
		return vm.NIL, nil
	})
	if err != nil {
		t.Fatalf("wrap native: %v", err)
	}

	nsA := mkLogChunk(logFn, vm.String("A"))
	nsB := mkLogChunk(logFn, vm.String("B"))

	t.Run("separate main chunk runs last", func(t *testing.T) {
		log = nil
		unit := &bytecode.ExecUnit{
			NSOrder:   []string{"nsA", "nsB"},
			NSChunks:  map[string]*vm.CodeChunk{"nsA": nsA, "nsB": nsB},
			MainChunk: mkLogChunk(logFn, vm.String("main")),
		}
		if err := RunExecUnit(unit); err != nil {
			t.Fatalf("RunExecUnit: %v", err)
		}
		if got, want := valStrings(log), []string{"A", "B", "main"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("main aliasing a namespace chunk runs exactly once", func(t *testing.T) {
		log = nil
		unit := &bytecode.ExecUnit{
			NSOrder:   []string{"nsA", "nsB"},
			NSChunks:  map[string]*vm.CodeChunk{"nsA": nsA, "nsB": nsB},
			MainChunk: nsB, // MainChunk IS the last ns chunk
		}
		if err := RunExecUnit(unit); err != nil {
			t.Fatalf("RunExecUnit: %v", err)
		}
		// nsB must not run a second time as "the main chunk".
		if got, want := valStrings(log), []string{"A", "B"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v (main chunk double-ran)", got, want)
		}
	})
}
