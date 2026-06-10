//go:build js && wasm

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * HostWriter routes runtime output to the WASM bundle's JS host directly, via
 * a global _lgOutput(string) callback, instead of writing to os.Stdout/Stderr
 * and relying on the bundle's JS-side fs.writeSync shim to intercept fd 1/2.
 * Installed as the root binding of the *out* / *err* vars (see
 * NewWriterHandle), it makes (println ...) and the error sites reach
 * LetGoHost.onOutput as the Go-side dual of api.WithStdout.
 *
 * Same hidden contract shape as js/emit's _lgEmit: a global function
 *   _lgOutput(s string)
 * defined per mode by the bundle bootstrap (worker version postMessages to the
 * main thread; main-thread version calls LetGoHost._output directly). Resolved
 * per Write, so boot order doesn't matter; if the bundle hasn't wired it
 * (running outside the official host), the write is dropped rather than
 * erroring — fire-and-forget, like emit.
 */

package rt

import "syscall/js"

// HostWriter is an io.Writer that forwards bytes to the JS _lgOutput global.
type HostWriter struct{}

// NewHostWriter returns a HostWriter. Pass it to NewWriterHandle to build an
// IOHandle suitable for the *out* / *err* vars.
func NewHostWriter() *HostWriter { return &HostWriter{} }

func (w *HostWriter) Write(p []byte) (int, error) {
	out := js.Global().Get("_lgOutput")
	if out.IsUndefined() {
		return len(p), nil
	}
	out.Invoke(string(p))
	return len(p), nil
}
