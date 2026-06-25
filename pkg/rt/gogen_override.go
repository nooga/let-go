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
	if defs := pendingGoOverrides[ns.Name()]; defs != nil {
		for name, fn := range defs {
			ns.Def(name, fn)
		}
		delete(pendingGoOverrides, ns.Name())
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
