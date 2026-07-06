/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"bytes"
	"fmt"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/vm"
)

// precompiledCoreNS holds the core bundle's per-namespace chunks so a required
// namespace can be replayed on demand without the compiler. Populated by
// LoadCore, consumed by the bytecode-only NSLoader.
var precompiledCoreNS map[string]*vm.CodeChunk

// LoadCore boots the runtime from the embedded core .lgb: decode it, replay the
// core main chunk (which defines all of core), then register the remaining
// namespace chunks for on-demand loading. It imports only bytecode + vm and
// touches no compiler code, so it can boot a build where pkg/compiler is not
// linked (see the runtime_only build). It mirrors the boot fast-path in
// pkg/compiler/eval.go (loadPrecompiledBundle), minus the compiler-side pieces
// (the global const pool for further compilation, and postCoreInit's eval
// builtins) that a runtime-only build deliberately omits.
func LoadCore() error {
	resolve := func(nsName, name string) *vm.Var {
		n := DefNSBare(nsName)
		if v := n.LookupLocal(vm.Symbol(name)); v != nil {
			return v
		}
		return n.DefStub(name)
	}
	unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(CoreCompiledLGB), resolve)
	if err != nil {
		return fmt.Errorf("decode core bundle: %w", err)
	}
	f := vm.NewFrame(unit.MainChunk, nil)
	_, err = f.RunProtected()
	vm.ReleaseFrame(f)
	if err != nil {
		return fmt.Errorf("replay core: %w", err)
	}
	if unit.NSChunks != nil {
		precompiledCoreNS = unit.NSChunks
		for name := range precompiledCoreNS {
			if name != NameCoreNS {
				MarkNSNeedsLoad(name)
			}
		}
	}
	return nil
}

// bytecodeNSLoader loads namespaces only from the precompiled bundle: on
// (require 'ns) it replays the stored .lgb chunk. It never reads or compiles
// source — the runtime-only counterpart to the resolver+compiler NSLoader, and
// the enforcement point of the "cannot execute arbitrary source" guarantee.
type bytecodeNSLoader struct{}

func (bytecodeNSLoader) Load(name string) *vm.Namespace {
	chunk := precompiledCoreNS[name]
	if chunk == nil {
		return nil // not in the bundle; source loading is disabled
	}
	f := vm.NewFrame(chunk, nil)
	_, _ = f.RunProtected()
	vm.ReleaseFrame(f)
	return NS(name)
}

// UseBytecodeNSLoader installs the bytecode-only loader. A runtime-only entry
// point calls this after LoadCore so that dynamic require resolves against the
// bundle and rejects source loading.
func UseBytecodeNSLoader() { SetNSLoader(bytecodeNSLoader{}) }
