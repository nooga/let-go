//go:build !bootstrap

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/vm"
)

// BootCore decodes and runs the embedded precompiled clojure.core bundle
// (CoreCompiledLGB) so a program that links only pkg/rt + pkg/vm — with no
// compiler — can resolve and invoke the full clojure.core surface plus the lg
// baseline namespaces. It returns a fresh ExecContext ready to Invoke.
//
// The motivating case is an AOT-lowered binary (scripts/lg-compile output): its
// generated funcs call core via ec.Invoke(rt.CachedVarFn(..., "clojure.core",
// name), ...), which needs those vars bound. Programs that touch only NATIVE
// core fns (str, +, vector, …) already work from a bare vm.NewExecContext(),
// because installers register those at package init; BootCore is what
// additionally binds the .lg-DEFINED core (and let-go.core / let-go.types).
//
// Call once at startup. This mirrors the runtime half of pkg/compiler/eval.go's
// loadPrecompiledBundle without linking the compiler; a future refactor should
// unify the two behind this function so there is a single core-load path.
func BootCore() (*vm.ExecContext, error) {
	if len(CoreCompiledLGB) == 0 {
		return nil, fmt.Errorf("BootCore: embedded core is empty (built -tags bootstrap?)")
	}

	// Namespaces present before decode are native-backed (installers ran at
	// package init); a bundle chunk for one of them is a hybrid ns whose chunk
	// must run eagerly (below), since qualified refs bypass the on-demand loader.
	preexisting := map[string]bool{}
	for name := range AllNSes() {
		preexisting[name] = true
	}

	// DefNSBare gives a var a home ns without triggering a loader; DefStub avoids
	// warn-on-shadow while the chunk Defs the real value later.
	resolve := func(nsName, name string) *vm.Var {
		n := DefNSBare(nsName)
		if v := n.LookupLocal(vm.Symbol(name)); v != nil {
			return v
		}
		return n.DefStub(name)
	}

	unit, err := bytecode.DecodeToExecUnitBytes(CoreCompiledLGB, resolve)
	if err != nil {
		return nil, fmt.Errorf("BootCore: decode core: %w", err)
	}

	// Replay core's def/defn/defmacro definitions.
	if err := runChunk(unit.MainChunk); err != nil {
		return nil, fmt.Errorf("BootCore: run core chunk: %w", err)
	}

	// The lg baseline namespaces (let-go.core, …) are auto-refer'd everywhere but
	// never explicitly required, so nothing else runs them — do it eagerly, after
	// core (whose defs they depend on).
	baseline := map[string]bool{}
	for _, name := range LgBaselineNSNames() {
		baseline[name] = true
		if unit.NSChunks != nil {
			if err := runChunk(unit.NSChunks[name]); err != nil {
				return nil, fmt.Errorf("BootCore: run baseline %s: %w", name, err)
			}
		}
	}

	// Hybrid namespaces (native fns + bundled lg source, e.g. async) are reachable
	// via qualified symbols without a require, so their chunks must run eagerly
	// too or the lg-defined vars stay nil stubs.
	if unit.NSChunks != nil {
		for name, ch := range unit.NSChunks {
			if name == "core" || baseline[name] || !preexisting[name] {
				continue
			}
			if err := runChunk(ch); err != nil {
				return nil, fmt.Errorf("BootCore: run hybrid %s: %w", name, err)
			}
		}
	}

	return vm.NewExecContext(), nil
}

func runChunk(ch *vm.CodeChunk) error {
	if ch == nil {
		return nil
	}
	f := vm.NewFrame(ch, nil)
	_, err := f.RunProtected()
	vm.ReleaseFrame(f)
	return err
}
