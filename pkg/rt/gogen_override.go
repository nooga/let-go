/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
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
			installGoOverride(ns, name, fn)
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

// installedGoOverrides holds every override installed so far, by namespace,
// so that ApplyGoOverrides can put them back after the namespace loads again.
// The core bundle decodes every namespace it holds up front, as stub Vars, so
// an override registered at init lands on a stub; the namespace's chunk
// replays later, on demand, and its `def`s set each Var's root to the
// bytecode fn. ApplyGoOverrides runs after that replay and restores these.
var installedGoOverrides = map[string]map[string]vm.Value{}

// installGoOverride makes fn the value of ns/name and records it in
// installedGoOverrides. An existing Var keeps its identity and takes fn as
// its root: code compiled before the override (a bundle decodes its var
// references up front) holds that Var, and ns.Def would replace it with a new
// one the compiled code never reads. A name with no Var yet is defined as
// usual.
func installGoOverride(ns *vm.Namespace, name string, fn vm.Value) {
	installed := installedGoOverrides[ns.Name()]
	if installed == nil {
		installed = map[string]vm.Value{}
		installedGoOverrides[ns.Name()] = installed
	}
	installed[name] = fn
	v := ns.LookupLocal(vm.Symbol(name))
	if v == nil {
		ns.Def(name, fn)
		return
	}
	if nf, ok := fn.(*vm.NativeFn); ok {
		nf.SetName(name)
	}
	v.SetRoot(fn)
}

// ApplyGoOverrides installs any pending overrides for ns, then puts back every
// override already installed on it, replacing whatever the bytecode or source
// load that just finished produced. No-op for a namespace with no overrides.
func ApplyGoOverrides(ns *vm.Namespace) {
	if ns == nil {
		return
	}
	if defs := pendingGoOverrides[ns.Name()]; defs != nil {
		for name, fn := range defs {
			installGoOverride(ns, name, fn)
		}
		delete(pendingGoOverrides, ns.Name())
	}
	for name, fn := range installedGoOverrides[ns.Name()] {
		if v := ns.LookupLocal(vm.Symbol(name)); v == nil || v.Deref() != fn {
			installGoOverride(ns, name, fn)
		}
	}
	if names := pendingNativeMultiFns[ns.Name()]; names != nil {
		freezeNativeMultiFns(ns, names)
		delete(pendingNativeMultiFns, ns.Name())
	}
}

// Native-baked multimethods (gogen_ir): a lowered package emitting native
// type-switch dispatch arms registers its multifn names here so the runtime
// can freeze them as the native baseline once the namespace finishes loading
// — i.e. after bytecode replay has created the multifn var and applied every
// build-time defmethod. A later defmethod then replaces the var with an
// unfrozen MultiFn, and the generated guard falls back to runtime dispatch.
var pendingNativeMultiFns = map[string][]string{}

// RegisterNativeMultiFns queues multimethod names whose native dispatch arms
// must be frozen for ns. If ns has already finished loading, the multifns are
// frozen immediately (the immediate-apply path RegisterGoOverrides also uses);
// otherwise they wait for ApplyGoOverrides.
func RegisterNativeMultiFns(nsName string, names []string) {
	if len(names) == 0 {
		return
	}
	if ns := LookupNS(nsName); ns != nil {
		freezeNativeMultiFns(ns, names)
		return
	}
	pendingNativeMultiFns[nsName] = append(pendingNativeMultiFns[nsName], names...)
}

// freezeNativeMultiFns marks each named var's MultiFn value as the native
// baseline. Names that are absent or not bound to a MultiFn are skipped — the
// generated guard treats a non-MultiFn / unfrozen value as "use runtime
// dispatch", so a miss is safe, never incorrect.
func freezeNativeMultiFns(ns *vm.Namespace, names []string) {
	for _, name := range names {
		v := ns.LookupLocal(vm.Symbol(name))
		if v == nil {
			continue
		}
		if mm, ok := v.Deref().(*vm.MultiFn); ok {
			mm.FreezeNative()
		}
	}
}
