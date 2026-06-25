/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"

	"github.com/nooga/let-go/pkg/vm"
)

// Go-native overrides for generated IR-stack namespaces.
//
// A Go-lowered package's init() registers its defns here. When the host
// resolver finishes loading a namespace (bytecode replay or source compile),
// it drains the pending overrides for that namespace, clobbering any vars
// the bytecode/source path produced with native Go implementations.
//
// Without any registrations (no gogen_ir build tag, no generated packages
// imported), the maps stay empty and the hook is a single map lookup.

var pendingGoOverrides = map[string]map[string]vm.Value{}

// RegisterGoOverrides queues a set of name → NativeFn bindings. If the
// target namespace already exists in the registry, the defs are applied
// immediately (the host has already finished bundle replay); otherwise
// they sit in pendingGoOverrides until ApplyGoOverrides drains them.
//
// Init-order hazard this guards against: pkg/compiler's init() replays
// the core bundle and calls ApplyGoOverrides(coreNS) BEFORE the blank-
// imported lowered packages in main get to run their own init(). Without
// the immediate-apply path here, every override registered after that
// point would sit in the queue forever — the queue would be drained
// once when empty, then filled and never re-drained. Symptom is the
// dispatch counters reporting zero and `(reduce + (range 1000))` going
// through bytecode dispatch even with -tags gogen_ir + all the blank
// imports wired up.
func RegisterGoOverrides(nsName string, defs map[string]vm.Value) {
	if len(defs) == 0 {
		return
	}
	if ns := LookupNS(nsName); ns != nil {
		for name, fn := range defs {
			ns.Def(name, fn)
		}
		return
	}
	existing := pendingGoOverrides[nsName]
	if existing == nil {
		existing = make(map[string]vm.Value, len(defs))
		pendingGoOverrides[nsName] = existing
	}
	for k, v := range defs {
		existing[k] = v
	}
}

// ApplyGoOverrides drains any pending overrides for ns and Defs them onto
// the namespace, replacing whatever bytecode or source replay produced.
// No-op if nothing is pending.
func ApplyGoOverrides(ns *vm.Namespace) {
	if ns == nil {
		return
	}
	defs := pendingGoOverrides[ns.Name()]
	if defs == nil {
		return
	}
	for name, fn := range defs {
		ns.Def(name, fn)
	}
	delete(pendingGoOverrides, ns.Name())
}

// GoVarInit is a top-level (def NAME value) lowered to native Go: Init builds
// the var's root value when given an ExecContext. Unlike a function override
// (a NativeFn that captures ec at call time), a def value — e.g. a parser
// combinator tree holding native closures — must be CONSTRUCTED at namespace
// load with an ec, so it is registered as an initializer drained after replay.
type GoVarInit struct {
	Name string
	Init func(*vm.ExecContext) (vm.Value, error)
}

// pendingGoVarInits holds, per namespace, the ordered var initializers a
// lowered package registered. Order is source order so a phase-1 pre-intern
// pass can make forward references resolve before any initializer runs.
var pendingGoVarInits = map[string][]GoVarInit{}

// RegisterGoVarInits queues a lowered package's top-level def initializers.
// They are drained by ApplyGoVarInits once the host finishes replaying the
// namespace's bytecode (which has an ExecContext in scope). Always queued —
// unlike RegisterGoOverrides there is no immediate-apply path, since running
// an initializer requires an ec the registration site does not have.
func RegisterGoVarInits(nsName string, inits []GoVarInit) {
	if len(inits) == 0 {
		return
	}
	pendingGoVarInits[nsName] = append(pendingGoVarInits[nsName], inits...)
}

// ApplyGoVarInits drains any pending var initializers for ns, running each
// with ec. Each initializer is a lowered top-level (def …): it interns its var
// and SetRoots the natively-built value itself, replacing whatever bytecode or
// source replay produced. Two phases: every target var is interned first (so a
// later def's value referencing an earlier def's var resolves), then each
// initializer runs in registration (source) order. An initializer that errors
// is reported and skipped, leaving the bytecode value in place. No-op if
// nothing is pending.
func ApplyGoVarInits(ns *vm.Namespace, ec *vm.ExecContext) {
	if ns == nil {
		return
	}
	inits := pendingGoVarInits[ns.Name()]
	if inits == nil {
		return
	}
	if ec == nil {
		ec = vm.RootExecContext
	}
	// Phase 1: ensure every target var exists so forward references resolve.
	for _, gi := range inits {
		if ns.LookupLocal(vm.Symbol(gi.Name)) == nil {
			ns.Def(gi.Name, vm.NIL)
		}
	}
	// Phase 2: run each initializer in source order — the lowered :def interns
	// and SetRoots the var as a side effect; the returned var value is ignored.
	for _, gi := range inits {
		if _, err := gi.Init(ec); err != nil {
			fmt.Fprintf(os.Stderr, "warning: go var-init %s/%s failed: %s\n", ns.Name(), gi.Name, err)
		}
	}
	delete(pendingGoVarInits, ns.Name())
}
