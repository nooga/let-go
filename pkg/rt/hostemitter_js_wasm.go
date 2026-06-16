//go:build js && wasm

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * HostEmitter routes (js/emit ...) to the WASM bundle's JS host via the
 * global _lgEmit(name, dataJson) callback — the dual of HostWriter's
 * _lgOutput. Installed as the *emit* root by the generated bundle main (see
 * wasmMainTmpl in wasm.go), it makes (js/emit :stats {...}) reach
 * LetGoHost.onEmit without js/emit touching syscall/js directly.
 *
 * Same hidden contract shape as HostWriter: a global function
 *   _lgEmit(name string, dataJson string)
 * defined per mode by the bundle bootstrap (worker version postMessages to
 * the main thread; main-thread version dispatches the CustomEvent directly).
 * Resolved per Emit, so boot order doesn't matter; if the bundle hasn't
 * wired it, the event is dropped rather than erroring — fire-and-forget.
 */

package rt

import "syscall/js"

// HostEmitter forwards events to the JS _lgEmit global.
type HostEmitter struct{}

// NewHostEmitter returns a HostEmitter suitable for the *emit* root binding.
func NewHostEmitter() *HostEmitter { return &HostEmitter{} }

func (e *HostEmitter) Emit(name, dataJSON string) {
	emit := js.Global().Get("_lgEmit")
	if emit.IsUndefined() {
		return
	}
	emit.Invoke(name, dataJSON)
}
