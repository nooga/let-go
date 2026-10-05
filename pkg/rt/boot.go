/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/vm"
)

// precompiledCore is the decoded core bundle, kept so a required namespace can
// be replayed on demand without the compiler — decoding its chunks first if
// the boot deferred them. Populated by LoadCore, consumed by the bytecode-only
// NSLoader.
var precompiledCore *bytecode.ExecUnit

// LoadCore boots the runtime from the embedded core .lgb: decode it, replay
// core + the lg baseline namespaces, then register the remaining namespace
// chunks for on-demand loading via bytecodeNSLoader. It touches no compiler
// code, so it can boot a build where pkg/compiler is not linked (see
// cmd/lg-runtime).
//
// Hybrid namespaces (native fns + a bundled lg chunk, e.g. async) are left for
// the on-demand loader — EagerHybrids is false — safe today because async's
// bundled content is macro-only and macros are fully expanded before a .lgb
// exists. LoadCore's whole reason for existing over BootCore is that it
// installs a namespace loader (UseBytecodeNSLoader), so it doesn't need the
// eager pass; a future hybrid bundling runtime fns would flip this to true.
func LoadCore() error {
	unit, err := LoadCoreBundle(CoreLoadOptions{EagerHybrids: false})
	if err != nil {
		return err
	}
	// Retain the unit so bytecodeNSLoader can replay a namespace on
	// (require ...).
	precompiledCore = unit
	return nil
}

// bytecodeNSLoader loads namespaces only from the precompiled bundle: on
// (require 'ns) it replays the stored .lgb chunk. It never reads or compiles
// source — the runtime-only counterpart to the resolver+compiler NSLoader, and
// the enforcement point of the "cannot execute arbitrary source" guarantee.
type bytecodeNSLoader struct{}

func (bytecodeNSLoader) Load(name string) *vm.Namespace {
	if precompiledCore == nil {
		return nil
	}
	chunk, err := precompiledCore.NSChunk(name)
	if err == nil && chunk == nil {
		return nil // not in the bundle; source loading is disabled
	}
	if err == nil {
		err = runChunk(chunk)
	}
	if err != nil {
		// NSLoader has no error channel: report, restore the needs-load marker,
		// and return nil. The namespace was pre-registered as a placeholder
		// during bytecode decoding, so without the restored marker the registry
		// would find the placeholder and treat the load as a success — caching a
		// half-initialized namespace. The marker is what makes RequireNS report
		// this failure and a later require retry the load. (The registry itself
		// can't restore it: a loader returning nil with a pre-existing namespace
		// also happens on legitimate native+lazy namespaces like gogen, where
		// the pre-existing namespace IS the result.)
		fmt.Fprintf(os.Stderr, "error: failed to load precompiled namespace %s: %s\n", name, err)
		MarkNSNeedsLoad(name)
		return nil
	}
	return NS(name)
}

// UseBytecodeNSLoader installs the bytecode-only loader. A runtime-only entry
// point calls this after LoadCore so that dynamic require resolves against the
// bundle and rejects source loading.
func UseBytecodeNSLoader() { SetNSLoader(bytecodeNSLoader{}) }
