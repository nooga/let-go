/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * KeySource is the host seam for blocking key input — the source behind
 * term/read-key and term/key-pending?. It is the input peer of the *out*
 * HostWriter and the *emit* HostEmitter: term ops consult a host-bound
 * capability at the *keys* var instead of reaching for os.Stdin / the
 * SharedArrayBuffer globals directly.
 *
 * Bound at the *keys* root per platform — a stdin+SIGWINCH source on native
 * (term.go), a SharedArrayBuffer source in the WASM bundle
 * (keysource_js_wasm.go) — and overridable per-Run by api.WithKeySource so
 * an embedder or test can feed synthetic keys.
 *
 * Scope: ReadKey + KeyPending only. Terminal geometry (size / set-size) stays
 * a separate term op, and wake — unblocking a parked ReadKey without a real
 * key — is deferred (see the note in pkg/rt/wasm/lg-host.js); it needs a
 * SAB-level protocol this seam leaves room for but doesn't introduce.
 */

package rt

import "github.com/nooga/let-go/pkg/vm"

// KeySource is a source of keypresses for term/read-key and key-pending?.
type KeySource interface {
	// ReadKey blocks until a key is available, returning the key's bytes as
	// a string. Returns "" for end-of-input (the read-key nil contract).
	ReadKey() (string, error)
	// KeyPending reports whether a key is buffered and ready, without
	// consuming it — the non-blocking peer of ReadKey.
	KeyPending() bool
}

// nopKeySource is the placeholder root: no input, nothing pending. Each
// platform's term install replaces it with a real source; it only governs
// stub platforms (plan9) and the pre-install window.
type nopKeySource struct{}

func (nopKeySource) ReadKey() (string, error) { return "", nil }
func (nopKeySource) KeyPending() bool         { return false }

// resolveKeySourceVar unwraps the current dynamic binding of varName (e.g.
// "*keys*") to a KeySource, mirroring resolveIOHandleVar / resolveEmitterVar.
// Returns nil if the var isn't installed or doesn't unwrap to a KeySource.
func resolveKeySourceVar(ec *vm.ExecContext, varName string) KeySource {
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
	if ks, ok := b.Unbox().(KeySource); ok {
		return ks
	}
	return nil
}

// boundKeySource returns the KeySource currently bound at *keys*, falling
// back to nopKeySource so the term ops always have something to call.
func boundKeySource(ec *vm.ExecContext) KeySource {
	if ks := resolveKeySourceVar(ec, "*keys*"); ks != nil {
		return ks
	}
	return nopKeySource{}
}
