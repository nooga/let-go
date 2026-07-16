/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

// Frame-execution tracing, gated by the dynamically-scoped *lg-trace* var.
//
// Tracing replaces the old `trace` special form. When *lg-trace* resolves
// truthy in a frame's execution context, that frame (and every callee, since
// the value propagates down the shared dynamic binding stack) prints its
// per-instruction stack + opcode dump — the same output the old special form
// produced, but scoped by `(binding [*lg-trace* true] ...)` instead of wrapping
// a single expression.
//
// The hot-path cost is kept to a single atomic Load per frame entry (the same
// shape as ProfilingEnabled): TraceArmed is false until *lg-trace* is first
// given a truthy value, so the more expensive per-frame Deref of TraceVar is
// only ever performed once tracing has actually been used in the process.

package vm

import "sync/atomic"

// TraceArmed is the coarse gate read once per frame entry in Frame.Run. It is
// flipped true the first time *lg-trace* is bound or set to a truthy value
// (see pushBinding and OP_SET_VAR). Once armed it stays armed for the process
// lifetime; a frame only actually traces when the precise per-frame Deref of
// TraceVar is also truthy, so disarming is unnecessary (an unbound/false
// *lg-trace* simply produces no trace).
var TraceArmed atomic.Bool

// TraceVarNS and TraceVarName identify the tracing control var, defined in
// the core namespace as `(def ^:dynamic *lg-trace* false)`. Both are matched
// so a user-namespace var that happens to be named *lg-trace* cannot hijack
// the gate.
const (
	TraceVarNS   = "core"
	TraceVarName = "*lg-trace*"
)

// TraceVar is the interned core `*lg-trace*` var, letting the VM — which cannot
// import pkg/rt — resolve the var's current dynamic value against a frame's
// ExecContext. It self-wires the first time *lg-trace* is bound or set truthy
// (armTraceIfTruthy), so no rt-side core-load hook is needed; nil until then,
// in which case tracing is inert.
var TraceVar atomic.Pointer[Var]

// armTraceIfTruthy flips the coarse gate on when the trace var receives a
// truthy value, capturing the var pointer the first time so later frame-entry
// checks are a plain pointer deref. Called from the (rare) dynamic-binding and
// set! paths; it costs a truthiness check plus — only until TraceVar is wired —
// a name compare, and nothing on the hot call path.
func armTraceIfTruthy(v *Var, val Value) {
	if v == nil || !IsTruthy(val) {
		return
	}
	current := TraceVar.Load()
	if v == current {
		TraceArmed.Store(true)
		return
	}
	if current == nil && v.NS() == TraceVarNS && v.VarName() == TraceVarName {
		// Use CompareAndSwap to safely set TraceVar the first time.
		if TraceVar.CompareAndSwap(nil, v) {
			TraceArmed.Store(true)
		}
	}
}
