/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import "sync"

// Dynamic, inspectable registry of Go packages that expose top-level
// functions callable directly from lowered Go output. The IR pipeline
// (lower-go) consults this registry at Pass-1 time and seeds
// *lowered-registry* so call-assign-stmts emits a direct Go call
// instead of an rt.InvokeValue + rt.LookupVar trampoline.
//
// Self-registration shape:
//
//	package corefns
//	import "github.com/nooga/let-go/pkg/rt"
//	func init() {
//	    rt.RegisterNativeModule(&rt.NativeModule{
//	        GoPkg:     "github.com/nooga/let-go/pkg/rt/corefns",
//	        Namespace: "clojure.core",
//	        Fns: map[string]rt.NativeDirectFn{
//	            "seq": {GoIdent: "Seq", Arity: 1,
//	                    ParamSpecs: []string{"vm.Value"},
//	                    ResultSpec: "vm.Value", NeedsError: true},
//	        },
//	    })
//	}
//
// Each contributing package self-registers — no centrally-maintained
// list. New modules added to the binary appear in the registry without
// touching IR-pipeline code. The IR pipeline reads the registry once
// at codegen time; runtime is read-only after init.

// NativeDirectFn describes one direct-callable Go function. The shape
// matches the per-arity entry that *lowered-registry* holds for lowered
// Lisp defns, so the IR pipeline can treat both uniformly.
//
// The current MVP scopes direct calls to vm.Value-only signatures with
// a (vm.Value, error) return. Variadic and typed-primitive params fall
// through to the existing rt.InvokeValue trampoline at the call site;
// they're representable here so future extensions can opt in.
type NativeDirectFn struct {
	GoIdent    string   // exported Go identifier, e.g. "Seq"
	Arity      int      // fixed arity; -1 for variadic
	Variadic   bool     // true → final ParamSpec is the rest slice element type
	ParamSpecs []string // Go type strings, e.g. []string{"vm.Value"}
	ResultSpec string   // Go type of single result, e.g. "vm.Value"
	NeedsError bool     // result tuple is (ResultSpec, error)
}

// NativeModule is one Go package's contribution to one let-go
// namespace. A package can register N functions; each is direct-callable
// from lowered Go.
type NativeModule struct {
	GoPkg     string                    // import path, e.g. "github.com/nooga/let-go/pkg/rt/corefns"
	Namespace string                    // target let-go namespace, e.g. "clojure.core"
	Fns       map[string]NativeDirectFn // let-go name → descriptor (e.g. "seq" → {GoIdent:"Seq", ...})
}

var (
	nativeMu      sync.RWMutex
	nativeMods    = map[string][]*NativeModule{}            // namespace → modules contributing to it
	nativeFnIndex = map[string]map[string]*NativeDirectFn{} // ns → lisp-name → descriptor (flat lookup)
)

// RegisterNativeModule self-registers a Go package as a contributor to
// a let-go namespace. Safe to call from init(); merges into the global
// registry without overwriting prior contributions to the same namespace.
//
// If the same (ns, name) is registered twice, the latest wins — this
// matches the var rebinding semantics on the let-go side and lets a
// downstream package selectively override a builtin.
func RegisterNativeModule(m *NativeModule) {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	nativeMods[m.Namespace] = append(nativeMods[m.Namespace], m)
	if nativeFnIndex[m.Namespace] == nil {
		nativeFnIndex[m.Namespace] = map[string]*NativeDirectFn{}
	}
	for name, fn := range m.Fns {
		fn := fn
		nativeFnIndex[m.Namespace][name] = &fn
	}
}

// LookupNativeDirect returns the descriptor for (ns, name) or nil if no
// module has registered it. The returned pointer is stable for the
// lifetime of the process.
func LookupNativeDirect(ns, name string) *NativeDirectFn {
	nativeMu.RLock()
	defer nativeMu.RUnlock()
	if m, ok := nativeFnIndex[ns]; ok {
		return m[name]
	}
	return nil
}

// LookupNativeModule returns the modules contributing to ns. Useful for
// inspection / debugging / codegen tools that need to walk by package.
func LookupNativeModule(ns string) []*NativeModule {
	nativeMu.RLock()
	defer nativeMu.RUnlock()
	mods := nativeMods[ns]
	if len(mods) == 0 {
		return nil
	}
	out := make([]*NativeModule, len(mods))
	copy(out, mods)
	return out
}

// AllNativeModules returns a snapshot of the entire registry. Useful for
// tools that need the full view — the IR pipeline at codegen time, a
// REPL `(rt/native-modules)` inspector, parity diff scripts.
func AllNativeModules() map[string][]*NativeModule {
	nativeMu.RLock()
	defer nativeMu.RUnlock()
	out := make(map[string][]*NativeModule, len(nativeMods))
	for ns, mods := range nativeMods {
		cp := make([]*NativeModule, len(mods))
		copy(cp, mods)
		out[ns] = cp
	}
	return out
}
