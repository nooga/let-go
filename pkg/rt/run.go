/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

// Compiler-free run lifecycle shared by every entry point that executes
// precompiled bytecode: the full lg binary (.lgb scripts and appended-payload
// bundles) and the runtime-only cmd/lg-runtime. Lives in rt so a CLI fix lands
// once instead of per entry point, and so nothing here can pull the compiler
// into a runtime-only dependency graph.

package rt

import (
	"bytes"
	"fmt"
	"os"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/vm"
)

// LGBVarResolver resolves a decoded var reference against the live registry:
// the var's home namespace is materialized bare (no loader trigger) and the
// name is interned as a stub if not yet defined. The single resolver behind
// every .lgb decode site.
//
// It records var-ref hit/miss decode stats (LG_DECODE_TAG_STATS). The Note
// calls self-gate on bytecode.decodeStatsEnabled, so they're a no-op unless a
// caller has enabled stats around the decode; folding them in here — rather
// than a decode-site-specific resolver — makes the counter available at every
// canonical decode site, not just core boot where it started (#356).
func LGBVarResolver(nsName, name string) *vm.Var {
	n := DefNSBare(nsName)
	if v := n.LookupLocal(vm.Symbol(name)); v != nil {
		bytecode.NoteDecodeVarRefHit()
		return v
	}
	bytecode.NoteDecodeVarRefMiss(true)
	return n.DefStub(name)
}

// DecodeExecUnit decodes a .lgb payload (plain or bundle format), resolving
// var references with LGBVarResolver.
func DecodeExecUnit(data []byte) (*bytecode.ExecUnit, error) {
	return bytecode.DecodeToExecUnit(bytes.NewReader(data), LGBVarResolver)
}

// DecodeExecUnitWithDebug decodes a payload and attaches a verified external
// debug companion when the payload was produced by bytecode.SplitDebug.
func DecodeExecUnitWithDebug(data, debugData []byte) (*bytecode.ExecUnit, error) {
	if len(debugData) == 0 {
		return DecodeExecUnit(data)
	}
	return bytecode.DecodeToExecUnitBytesWithDebug(data, debugData, LGBVarResolver)
}

// DecodeExecUnitWithDebugFile loads and decodes a companion only when the LGB
// header says debug information was split. Ordinary artifacts therefore pay
// no filesystem lookup and ignore stale sidecars.
func DecodeExecUnitWithDebugFile(data []byte, artifactPath string) (*bytecode.ExecUnit, error) {
	if !bytecode.HasSplitDebug(data) {
		return DecodeExecUnit(data)
	}
	debugData, err := LoadDebugCompanion(artifactPath)
	if err != nil {
		return nil, err
	}
	return DecodeExecUnitWithDebug(data, debugData)
}

// LoadDebugCompanion reads debug information for artifactPath. LG_DEBUG_FILE,
// when set, selects an explicit file (or disables loading when empty); otherwise
// the conventional <artifact>.debug sibling is optional.
func LoadDebugCompanion(artifactPath string) ([]byte, error) {
	path, explicit := os.LookupEnv("LG_DEBUG_FILE")
	if explicit && path == "" {
		return nil, nil
	}
	if !explicit {
		path = artifactPath + bytecode.DebugCompanionSuffix
	}
	data, err := os.ReadFile(path)
	if err == nil {
		return data, nil
	}
	if !explicit && os.IsNotExist(err) {
		return nil, nil
	}
	return nil, fmt.Errorf("reading debug companion %s: %w", path, err)
}

// RunExecUnit replays a decoded unit: namespace chunks in dependency order
// first, then the main chunk (when it is not already one of those namespaces).
func RunExecUnit(unit *bytecode.ExecUnit) error {
	if err := LoadProgramNamespaces(unit); err != nil {
		return err
	}
	return RunProgramMainChunk(unit)
}

// LoadProgramNamespaces replays every namespace chunk in unit.NSOrder, draining
// ApplyGoOverrides after each so gogen_ir NativeFn overrides land on the Vars
// the bytecode just installed. This is the program-half of an AOT native-entry
// frame (#425): BootCore already owns core boot+reapply; the generated main.go
// owns only this load, then either calls the lowered entry directly or falls
// back to RunProgramMainChunk.
//
// Unlike the historical RunExecUnit loop, this does NOT skip MainChunk —
// for a single-ns program the ns chunk IS MainChunk, and the native entry
// still needs those Vars (non-lowered callees, override targets) installed
// before the direct Go call. RunProgramMainChunk then no-ops when MainChunk
// was already covered here.
func LoadProgramNamespaces(unit *bytecode.ExecUnit) error {
	if unit == nil {
		return fmt.Errorf("LoadProgramNamespaces: nil ExecUnit")
	}
	for _, name := range unit.NSOrder {
		chunk := unit.NSChunks[name]
		if chunk == nil {
			continue
		}
		if err := runChunk(chunk); err != nil {
			return fmt.Errorf("loading namespace %s: %w", name, err)
		}
		if ns := LookupNS(name); ns != nil {
			ApplyGoOverrides(ns)
		}
	}
	return nil
}

// RunProgramMainChunk runs unit.MainChunk on the VM when it was not already
// replayed by LoadProgramNamespaces (i.e. when MainChunk is not one of the
// NSOrder namespaces). The AOT native-entry VM-fallback path (#425 Finding 2)
// uses this after LoadProgramNamespaces when the recognized main/-main form
// did not lower.
func RunProgramMainChunk(unit *bytecode.ExecUnit) error {
	if unit == nil {
		return fmt.Errorf("RunProgramMainChunk: nil ExecUnit")
	}
	if unit.MainChunk == nil {
		return nil
	}
	for _, name := range unit.NSOrder {
		if unit.NSChunks[name] == unit.MainChunk {
			// Already loaded (and override-drained) by LoadProgramNamespaces.
			return nil
		}
	}
	return runChunk(unit.MainChunk)
}

// InvokeProgramEntry looks up name ("-main" or "main") in the given
// namespace and Invokes it with args. Used by the #425 native-entry
// VM-fallback frame after LoadProgramNamespaces has installed defs — a bare
// MainChunk re-run would no-op for single-ns programs whose MainChunk is the
// ns chunk itself.
//
// namespace must be the source namespace that owns the selected entry. The
// compiler already chose that namespace; searching CurrentNS or the global
// registry would let a competing entry in another namespace steal the call.
func InvokeProgramEntry(ec *vm.ExecContext, namespace, name string, args []vm.Value) error {
	if ec == nil {
		return fmt.Errorf("InvokeProgramEntry: nil ExecContext")
	}
	if namespace == "" {
		return fmt.Errorf("InvokeProgramEntry: empty namespace")
	}
	ns := LookupNS(namespace)
	if ns == nil {
		return fmt.Errorf("InvokeProgramEntry: namespace %s not found", namespace)
	}
	v := ns.LookupLocal(vm.Symbol(name))
	if v == nil {
		return fmt.Errorf("InvokeProgramEntry: %s/%s not found", namespace, name)
	}
	fn, ok := v.Deref().(vm.Fn)
	if !ok {
		return fmt.Errorf("InvokeProgramEntry: %s/%s is not callable", namespace, name)
	}
	_, err := ec.Invoke(fn, args)
	return err
}

func runChunk(c *vm.CodeChunk) error {
	f := vm.NewFrame(c, nil)
	_, err := f.RunProtected()
	vm.ReleaseFrame(f)
	return err
}

// SetCommandLineArgs publishes the user's CLI args — the positionals after the
// script — to core/*command-line-args*: nil when there are none, else a seq of
// strings. The entry point is the only layer that knows authoritatively where
// the script ends and the user's args begin, so it computes them once and
// every consumer reads the var instead of slicing os/args by hand.
func SetCommandLineArgs(args []string) {
	var val vm.Value = vm.NIL
	if len(args) > 0 {
		vs := make([]vm.Value, len(args))
		for i, a := range args {
			vs[i] = vm.String(a)
		}
		val = vm.NewList(vs)
	}
	CoreNS.Lookup("*command-line-args*").(*vm.Var).SetRoot(val)
}

// InstallPersistentStorage binds core/*storage* to the default file-backed
// store for storeID. Failure disables storage with a warning rather than
// aborting startup.
func InstallPersistentStorage(storeID string) {
	store, err := NewDefaultFileStorage(storeID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: storage disabled: %v\n", err)
		return
	}
	if v := LookupCoreVar("*storage*"); v != nil {
		v.SetRoot(vm.NewBoxed(store))
	}
}
