/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/nooga/let-go/pkg/vm"
)

// regPrimDebug, when LG_REGPRIM_DEBUG is set, logs every generated-primitive
// registration (name, requested ns, resolved canonical ns, whether it landed
// in the canonical ns). It is the per-primitive trace used to root-cause a
// batch of newly hoisted //lg:native primitives without bisecting the whole
// set: grep the run for `landed-canonical=false`.
var regPrimDebug = os.Getenv("LG_REGPRIM_DEBUG") != ""

// Generated-primitive binding lifecycle.
//
// RegisterGeneratedPrimitives (zz_primitives_generated.go) Defs each native
// adapter into its target namespace during rt init. A namespace whose source
// loads AFTER init — e.g. clojure.string on first (require ...) — re-runs its
// bootstrap Lisp definitions (string.lg's `(defn upper-case ...)`), which
// overwrites the native adapter with the bootstrap closure. Each generated
// Def is therefore recorded here, and reapplyGeneratedPrimitives re-Defs the
// adapters right after a namespace's on-demand load completes, so the native
// implementation is what user code observes. Explicit user redefinition
// (def / with-redefs / alter-var-root) still wins: it happens after the load,
// and the reapply fires only on the load itself.
var (
	genPrimMu       sync.RWMutex
	genPrimBindings = map[string]map[string]vm.Value{} // canonical ns name -> var name -> adapter
)

// NativePrimsIntact reports that no native-primitive var root is currently
// overridden (with-redefs / alter-var-root / re-def). Lowered Go call sites
// emitted for native-module callees consult this before taking the baked
// direct call; on false they fall back to the var-dispatch trampoline so the
// override is observed across lowered function boundaries.
func NativePrimsIntact() bool {
	if vm.GuardedRootsIntact() {
		return true
	}
	// Off the fast path (a root deviated): under LG_GUARD_DEBUG, name the
	// offending native primitive(s) — the diagnostic that pinpoints which
	// hoisted primitive is being shadowed instead of a bare intact=false.
	if guardDebug {
		fmt.Fprintf(os.Stderr, "[GUARD] deviated: %v\n", deviatedGuardedVars())
	}
	return false
}

var guardDebug = os.Getenv("LG_GUARD_DEBUG") != ""

// deviatedGuardedVars lists every interned var whose guarded root has deviated
// from its canonical native root, "<ns>/<name>". Diagnostic only.
func deviatedGuardedVars() []string {
	nsMu.RLock()
	names := make([]string, 0, len(nsRegistry))
	for k := range nsRegistry {
		names = append(names, k)
	}
	nsMu.RUnlock()
	var dev []string
	for _, nsName := range names {
		ns := LookupNS(nsName)
		if ns == nil {
			continue
		}
		for sym, v := range ns.AllVars() {
			if v.GuardDeviated() {
				dev = append(dev, nsName+"/"+string(sym))
			}
		}
	}
	sort.Strings(dev)
	return dev
}

// guardModuleVars marks the CURRENT root of every var a native module
// direct-calls (corefns seq/first/…, builtins vector/cons/…) as canonical,
// wiring them into the deviation counter behind NativePrimsIntact. Vars that
// don't exist yet (adapters Def'd right after registration) are guarded by
// defGeneratedPrimitive instead. Called from RegisterNativeModule.
func guardModuleVars(m *NativeModule) {
	ns := LookupNS(m.Namespace)
	if ns == nil {
		return
	}
	for key, d := range m.Fns {
		name := d.LgName
		if name == "" {
			name = key
			if i := strings.IndexByte(name, '@'); i >= 0 {
				name = name[:i]
			}
		}
		if v := ns.LookupLocal(vm.Symbol(name)); v != nil {
			v.GuardRoot()
		}
	}
}

// setPrimitiveRoot binds a native adapter as a var's root, mutating the
// EXISTING var in place when the name is already interned. Namespace.Def
// creates a fresh Var each call, which would split var identity: code
// compiled before this binding holds the old Var object and stops observing
// later redefinitions through the language path (which interns via
// LookupOrAdd and SetRoots the same object).
func setPrimitiveRoot(ns *vm.Namespace, name string, v vm.Value) *vm.Var {
	if existing := ns.LookupLocal(vm.Symbol(name)); existing != nil {
		return existing.SetRoot(v)
	}
	return ns.Def(name, v)
}

// defGeneratedPrimitive binds a generated native adapter into ns and records
// the binding so it can be reapplied after the namespace's source loads.
// Called only from the generated RegisterGeneratedPrimitives.
func defGeneratedPrimitive(ns *vm.Namespace, nsName, name string, v vm.Value) {
	// GuardRoot marks the adapter as the var's canonical root: lowered
	// direct-call sites stay on the native fast path only while the root is
	// untouched (vm.GuardedRootsIntact), so with-redefs/alter-var-root in a
	// caller's dynamic extent is observed even across lowered function
	// boundaries.
	setPrimitiveRoot(ns, name, v).GuardRoot()
	canonical := resolveNSAlias(nsName)
	if regPrimDebug {
		landed := false
		if cns := LookupNS(canonical); cns != nil {
			landed = cns.LookupLocal(vm.Symbol(name)) != nil
		}
		fmt.Fprintf(os.Stderr, "[REGPRIM] %s/%s requested-ns=%s canonical=%s landed-canonical=%v\n",
			nsName, name, nsName, canonical, landed)
	}
	genPrimMu.Lock()
	m := genPrimBindings[canonical]
	if m == nil {
		m = map[string]vm.Value{}
		genPrimBindings[canonical] = m
	}
	m[name] = v
	genPrimMu.Unlock()
}

// AuditGeneratedPrimitives verifies that every recorded generated primitive is
// actually bound in its canonical namespace, and returns one human-readable
// line per primitive that is NOT — including where the name IS bound instead,
// if anywhere. It never panics and never stops at the first problem: it is the
// non-terminating diagnostic run around the bootstrap core.lg compile, so a
// batch of newly hoisted //lg:native primitives surfaces ALL its registration
// failures at once (`+ registered under "clojure.core", not canonical "core"`)
// rather than the compile dying on the first `Can't resolve X`. Empty slice =
// every generated primitive resolves in its canonical namespace.
func AuditGeneratedPrimitives() []string {
	type binding struct {
		ns, name string
	}
	genPrimMu.RLock()
	bindings := make([]binding, 0, len(genPrimBindings))
	for nsName, m := range genPrimBindings {
		for name := range m {
			bindings = append(bindings, binding{nsName, name})
		}
	}
	genPrimMu.RUnlock()

	var problems []string
	for _, b := range bindings {
		ns := LookupNS(b.ns)
		if ns == nil {
			problems = append(problems, fmt.Sprintf("%s/%s: canonical namespace does not exist", b.ns, b.name))
			continue
		}
		if v := ns.LookupLocal(vm.Symbol(b.name)); v == nil || !v.IsBound() {
			problems = append(problems, fmt.Sprintf("%s/%s: not bound in canonical ns%s", b.ns, b.name, whereNameBound(b.name, resolveNSAlias(b.ns))))
		}
	}
	sort.Strings(problems)
	return problems
}

// whereNameBound scans the namespace registry for any namespace OTHER than
// exclude that has name bound, returning a " (found bound in: …)" suffix so a
// misregistered primitive points straight at the namespace it landed in (the
// exact signature of the alias-resolution bug: bound in "clojure.core", absent
// from canonical "core"). Empty string when the name is bound nowhere else.
func whereNameBound(name, exclude string) string {
	nsMu.RLock()
	defer nsMu.RUnlock()
	var found []string
	for k, ns := range nsRegistry {
		if k == exclude {
			continue
		}
		if v := ns.LookupLocal(vm.Symbol(name)); v != nil && v.IsBound() {
			found = append(found, k)
		}
	}
	if len(found) == 0 {
		return ""
	}
	sort.Strings(found)
	return " (found bound in: " + strings.Join(found, ", ") + ")"
}

// ReapplyGeneratedPrimitives restores the recorded native adapters on a
// namespace after one of its chunks executed outside the on-demand loader —
// the compiler's eager hybrid-namespace pass (loadPrecompiledBundle) runs
// chunks directly, and those re-Def bootstrap closures over the adapters
// just like a lazy load does. No-op for namespaces without generated
// primitives or that don't exist.
func ReapplyGeneratedPrimitives(name string) {
	reapplyGeneratedPrimitives(name, LookupNS(name))
}

// reapplyGeneratedPrimitives restores the recorded native adapters on a
// namespace after its on-demand source/bundle load executed (which re-Defs
// bootstrap closures over them). No-op for namespaces without generated
// primitives.
func reapplyGeneratedPrimitives(name string, ns *vm.Namespace) {
	if ns == nil {
		return
	}
	genPrimMu.RLock()
	m := genPrimBindings[resolveNSAlias(name)]
	genPrimMu.RUnlock()
	if len(m) == 0 {
		return
	}
	// Same trusted re-registration as init — silence the core-shadow warning.
	// SetRoot-in-place (not Def): the load compiled code against the interned
	// Var objects, and replacing them here would strand that code on stale
	// vars that later redefinitions never reach.
	vm.SetSuppressShadowWarn(true)
	for n, v := range m {
		setPrimitiveRoot(ns, n, v)
	}
	vm.SetSuppressShadowWarn(false)
}
