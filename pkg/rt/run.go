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

// RunExecUnit replays a decoded unit: namespace chunks in dependency order
// first, then the main chunk.
func RunExecUnit(unit *bytecode.ExecUnit) error {
	for _, name := range unit.NSOrder {
		chunk := unit.NSChunks[name]
		if chunk == nil || chunk == unit.MainChunk {
			continue
		}
		if err := runChunk(chunk); err != nil {
			return fmt.Errorf("loading namespace %s: %w", name, err)
		}
	}
	return runChunk(unit.MainChunk)
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
