/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * Emitter is the host seam for (js/emit ...): a sink for structured events
 * the guest dispatches to its host. The dual of the *out* HostWriter, and
 * simpler — fire-and-forget, no blocking, no wake, no return value.
 *
 * Bound at the *emit* root by the host: HostEmitter in the WASM bundle (see
 * hostemitter_js_wasm.go), a FuncEmitter via api.WithEmit for Go embedders.
 * Defaults to nopEmitter so (js/emit ...) is harmless when nothing is
 * listening, the same way *out* falls back when unbound. Resolution rides
 * the *emit* dynamic var, so (binding [*emit* ...] ...) works and per-Run
 * api bindings stay isolated between LetGo instances.
 */

package rt

import "github.com/nooga/let-go/pkg/vm"

// Emitter is a sink for guest-dispatched host events.
type Emitter interface {
	Emit(name, dataJSON string)
}

// FuncEmitter adapts a plain func into an Emitter, for api.WithEmit.
type FuncEmitter func(name, dataJSON string)

// Emit implements Emitter.
func (f FuncEmitter) Emit(name, dataJSON string) { f(name, dataJSON) }

// nopEmitter drops events. Root binding of *emit* until a host installs one.
type nopEmitter struct{}

func (nopEmitter) Emit(string, string) {}

// resolveEmitterVar unwraps the current dynamic binding of varName (e.g.
// "*emit*") to an Emitter, mirroring resolveIOHandleVar. Returns nil if the
// var isn't installed yet or its binding doesn't unwrap to an Emitter.
func resolveEmitterVar(ec *vm.ExecContext, varName string) Emitter {
	ns := lookupNSCached(NameCoreNS)
	if ns == nil {
		return nil
	}
	v := ns.LookupLocal(vm.Symbol(varName))
	if v == nil {
		return nil
	}
	b, ok := ec.Deref(v).(*vm.Boxed)
	if !ok {
		return nil
	}
	if e, ok := b.Unbox().(Emitter); ok {
		return e
	}
	return nil
}

// EmitVia dispatches through the current *emit* binding. No-op when unbound
// (early boot, or a host that never installed an emitter).
func EmitVia(ec *vm.ExecContext, name, dataJSON string) {
	if e := resolveEmitterVar(ec, "*emit*"); e != nil {
		e.Emit(name, dataJSON)
	}
}
