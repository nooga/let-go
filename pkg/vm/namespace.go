/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
)

type theNamespaceType struct{}

func (t *theNamespaceType) String() string  { return t.Name() }
func (t *theNamespaceType) Type() ValueType { return TypeType }
func (t *theNamespaceType) Unbox() any      { return reflect.TypeFor[*theNamespaceType]() }

func (t *theNamespaceType) Name() string { return "let-go.lang.Namespace" }
func (t *theNamespaceType) Box(fn any) (Value, error) {
	return NIL, NewTypeError(fn, "can't be boxed as", t)
}

var NamespaceType *theNamespaceType = &theNamespaceType{}

// nsLookup is set by the rt package to enable qualified symbol resolution
// across all loaded namespaces (e.g., foo/x looks up ns "foo" globally).
var nsLookup func(name string) *Namespace

// SetNSLookup sets the global namespace lookup function.
func SetNSLookup(fn func(string) *Namespace) {
	nsLookup = fn
}

// Refer entries are immutable once constructed and inserted, so their fields
// (ns, all, only) may be read without holding the owning namespace's lock —
// only the refers MAP itself is mutex-guarded.
type Refer struct {
	ns   *Namespace
	all  bool
	only map[Symbol]bool
}

type Namespace struct {
	// mu guards the four maps below. Governing invariant: NEVER hold one
	// namespace's lock while acquiring another's. Cross-namespace reads
	// (Lookup/Def/FuzzySymbolLookup) snapshot or use the target's own brief
	// accessors, so no lock is ever held across another lock acquisition —
	// the lock graph has no cycles (writers only lock their own namespace).
	mu       sync.RWMutex
	name     string // immutable after construction; Name() is lock-free
	registry map[Symbol]*Var
	refers   map[Symbol]*Refer
	aliases  map[Symbol]*Namespace
	excludes map[Symbol]bool // names excluded from clojure.core auto-refer
}

// coreNamespacePtr is set by the rt package after clojure.core is registered.
// Used by Def to check whether a name shadows core.
var coreNamespacePtr *Namespace

// SetCoreNamespace registers clojure.core for the warn-on-shadow check.
// Called once during rt initialization.
func SetCoreNamespace(ns *Namespace) {
	coreNamespacePtr = ns
}

// lgBaselineNamespaces are the lg-specific extra namespaces (let-go.core,
// let-go.types, …) auto-refer'd alongside clojure.core. Like clojure.core each
// is a lowest-priority resolution baseline: an explicit refer shadows it. Set by
// rt init via AddBaselineNamespace.
var lgBaselineNamespaces []*Namespace

// AddBaselineNamespace registers ns as an auto-refer'd resolution baseline.
func AddBaselineNamespace(ns *Namespace) {
	lgBaselineNamespaces = append(lgBaselineNamespaces, ns)
}

// isBaselineNS reports whether ns is an auto-refer'd resolution baseline
// (clojure.core or an lg-baseline like let-go.core/types) — deferred behind
// explicit refers in Lookup.
func isBaselineNS(ns *Namespace) bool {
	if coreNamespacePtr != nil && ns == coreNamespacePtr {
		return true
	}
	for _, b := range lgBaselineNamespaces {
		if ns == b {
			return true
		}
	}
	return false
}

// --- brief single-op locked accessors -------------------------------------
// Each takes its own lock for exactly one map op and returns a copy/pointer,
// never a live map reference, so callers never iterate an unguarded map.

// localVar reads the namespace's own registry. nil when absent.
func (n *Namespace) localVar(s Symbol) *Var {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.registry[s]
}

// aliasFor reads the namespace's own aliases.
func (n *Namespace) aliasFor(s Symbol) (*Namespace, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	t, ok := n.aliases[s]
	return t, ok
}

// referFor reads a single refer entry by key.
func (n *Namespace) referFor(s Symbol) (*Refer, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	r, ok := n.refers[s]
	return r, ok
}

// refersSnapshot copies the refer pointers so callers can iterate (and follow
// each refer's target namespace) without holding this namespace's lock — the
// Refer values are immutable once inserted.
func (n *Namespace) refersSnapshot() []*Refer {
	n.mu.RLock()
	defer n.mu.RUnlock()
	out := make([]*Refer, 0, len(n.refers))
	for _, r := range n.refers {
		out = append(out, r)
	}
	return out
}

// excludedLocked reports whether the symbol is in the exclude set.
func (n *Namespace) excludedLocked(s Symbol) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.excludes[s]
}

// cacheAlias rewrites an alias to its freshly-resolved target. Extracted from
// Lookup so the write is a single guarded op rather than inline under a read.
func (n *Namespace) cacheAlias(s Symbol, target *Namespace) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.aliases[s] = target
}

// Exclude marks a symbol as excluded from clojure.core auto-refer.
// Called from the ns macro for :refer-clojure :exclude [...].
func (n *Namespace) Exclude(name string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.excludes[Symbol(name)] = true
}

// IsExcluded reports whether the symbol is in the exclude set.
func (n *Namespace) IsExcluded(name Symbol) bool {
	return n.excludedLocked(name)
}

func (n *Namespace) Type() ValueType { return NamespaceType }

// Unbox implements Unbox
func (n *Namespace) Unbox() any {
	return nil
}

func NewNamespace(name string) *Namespace {
	return &Namespace{
		name:     name,
		registry: map[Symbol]*Var{},
		refers:   map[Symbol]*Refer{},
		aliases:  map[Symbol]*Namespace{},
		excludes: map[Symbol]bool{},
	}
}

func (n *Namespace) RegistrySize() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return len(n.registry)
}

// isShadowingCoreRefer reports whether name `s` is currently visible
// unqualified in namespace `n` via a refer of clojure.core.
//
// Refer entries are keyed by namespace name (e.g. "clojure.core"), not
// by symbol — so we look up that single entry, then check whether `s`
// is in scope via :refer :all or :refer :only.
func isShadowingCoreRefer(n *Namespace, s Symbol) bool {
	for _, ref := range n.refersSnapshot() {
		if ref == nil || ref.ns != coreNamespacePtr {
			continue
		}
		if ref.all {
			return true
		}
		if ref.only != nil && ref.only[s] {
			return true
		}
	}
	return false
}

func (n *Namespace) Def(name string, val Value) *Var {
	s := Symbol(name)
	// Warn-on-core-shadow: emit Clojure-parity warning when a non-core
	// namespace defines a name that is currently REFERRED in from
	// clojure.core (i.e. previously visible in this ns unqualified),
	// unless explicitly excluded via (:refer-clojure :exclude).
	//
	// Clojure JVM only warns on shadow-of-refer, not on raw name overlap:
	//   (ns foo (:refer-clojure :only [defn]))
	//   (defn reset! [x] x)  ;; no warning — reset! was never refered in
	//
	// Stdlib Go-side ns.Def calls (e.g. profile/reset!) build namespaces
	// that don't auto-refer clojure.core, so they correctly stay silent.
	// User code that uses the default (ns ...) form gets clojure.core
	// auto-refered :all, so it does warn on shadow. The check reads core's
	// and this ns's maps via brief accessors (no lock held across the write
	// below); the warning itself is best-effort so a benign TOCTOU is fine.
	if coreNamespacePtr != nil && n != coreNamespacePtr && !n.excludedLocked(s) {
		if isShadowingCoreRefer(n, s) {
			if existing := coreNamespacePtr.localVar(s); existing != nil && !existing.isPrivate {
				// Only warn the first time we shadow in this ns; subsequent
				// re-defs of our own var don't re-warn.
				if n.localVar(s) == nil {
					fmt.Fprintf(os.Stderr,
						"WARNING: %s already refers to: #'clojure.core/%s in namespace: %s, being replaced by: #'%s/%s\n",
						name, name, n.name, n.name, name)
				}
			}
		}
	}
	va := NewVar(n, n.name, name)
	va.SetRoot(val)
	if val.Type() == NativeFnType {
		val.(*NativeFn).SetName(name)
	}
	if f, ok := val.(*Func); ok {
		f.SetName(name)
	}
	n.mu.Lock()
	n.registry[s] = va
	n.mu.Unlock()
	return va
}

// LookupLocal checks only the namespace's own registry, not refers or aliases.
func (n *Namespace) LookupLocal(symbol Symbol) *Var {
	return n.localVar(symbol)
}

// lookupViaRefers searches for a symbol in the namespace's refers, honoring
// :all/:only filters and baseline-namespace precedence. Returns nil if not found.
// Baseline namespaces (clojure.core, let-go.core, etc.) are lowest-priority:
// explicit refers shadow them. Private vars are always hidden.
func (n *Namespace) lookupViaRefers(sym Symbol) *Var {
	var baselineHit *Var
	for _, ref := range n.refersSnapshot() {
		v := ref.ns.localVar(sym)
		if v == nil || v.isPrivate {
			continue
		}
		isBaseline := isBaselineNS(ref.ns)
		// Baseline namespaces are always the full auto-refer :all; explicit
		// refers honor :refer [syms] — an :only refer contributes a symbol
		// only when it is listed.
		if !isBaseline && !ref.all && (ref.only == nil || !ref.only[sym]) {
			continue
		}
		if isBaseline {
			baselineHit = v
			continue
		}
		return v
	}
	return baselineHit
}

// DefStub interns an UNBOUND var (no root) without triggering the
// warn-on-shadow check. Intended for bundle decoders that pre-populate var
// references before the namespace's own chunk runs (which Defs the ones that
// have an init). Leaving the root unset is load-bearing: a no-init (def x) /
// (declare x) emits only OP_LOAD_CONST, so the stub is the var that survives
// into the namespace, and it must report (bound? (var x)) == false to match
// the source-compiled path (LookupOrAdd). A NIL root would make precompiled /
// LGB-loaded code wrongly appear bound, breaking def/declare promise semantics.
// Do NOT use DefStub to intentionally suppress warnings for new code; use
// Namespace.Exclude (via :refer-clojure :exclude) instead.
func (n *Namespace) DefStub(name string) *Var {
	s := Symbol(name)
	va := NewVar(n, n.name, name)
	n.mu.Lock()
	n.registry[s] = va
	n.mu.Unlock()
	return va
}

func (n *Namespace) LookupOrAdd(symbol Symbol) Value {
	if v := n.localVar(symbol); v != nil {
		return v
	}
	// Intern an UNBOUND var (no root) rather than Def(NIL): the compiler
	// calls this while compiling a `(def x v)` form before it runs, and
	// `defonce` must be able to tell that interned-but-unrun state from a
	// var that has actually been assigned. Deref still yields NIL.
	n.mu.Lock()
	defer n.mu.Unlock()
	// Re-check under the write lock: a concurrent LookupOrAdd may have
	// interned it between our RLock read and acquiring the Lock.
	if v, ok := n.registry[symbol]; ok {
		return v
	}
	va := NewVar(n, n.name, string(symbol))
	n.registry[symbol] = va
	return va
}

func (n *Namespace) Lookup(symbol Symbol) Value {
	noteLookup(n.name, string(symbol))
	sns, sym, hasNS := symbol.NamespacedRaw()
	if !hasNS {
		if v := n.localVar(sym); v != nil {
			return v
		}
		// Unqualified miss: search refers via the baseline-aware helper.
		//
		// Precedence (Clojure semantics): an explicit refer — (:require [lib
		// :refer :all]) or :refer [syms], and (use lib) — SHADOWS the auto-refer
		// baseline namespaces (clojure.core and let-go.core) that RegisterNS
		// installs into every ns. lookupViaRefers returns the first matching
		// explicit refer (honoring :all / :only) and falls back to a baseline
		// hit only when no explicit refer provides the symbol.
		if v := n.lookupViaRefers(sym); v != nil {
			return v
		}
		return NIL
	}
	// Alias-qualified resolution via aliases
	if target, ok := n.aliasFor(sns); ok {
		v := target.localVar(sym)
		if v == nil && nsLookup != nil {
			// Alias may point to a placeholder namespace created before source
			// load completed. Re-resolve by name so runtime loader can
			// materialize the namespace on demand, then retry the symbol lookup.
			if loaded := nsLookup(target.Name()); loaded != nil {
				target = loaded
				n.cacheAlias(sns, loaded)
				v = target.localVar(sym)
			}
		}
		// If direct lookup missed, try target's refers (follows target's `:refer`s)
		if v == nil {
			v = target.lookupViaRefers(sym)
		}
		if v == nil || v.isPrivate {
			return NIL
		}
		return v
	}
	// Fallback: direct namespace lookup from global registry
	if nsLookup != nil {
		if target := nsLookup(string(sns)); target != nil {
			v := target.localVar(sym)
			// A private var is visible to a fully-qualified reference only from
			// within its own namespace — `my.ns/-priv` is legal inside my.ns
			// (e.g. a macro that expands to a qualified call to a private helper
			// in the same ns).
			if v != nil && (!v.isPrivate || target == n) {
				return v
			}
		}
	}
	// Fallback via refers
	if refer, ok := n.referFor(sns); ok {
		v := refer.ns.localVar(sym)
		if v == nil || v.isPrivate {
			return NIL
		}
		if !refer.all {
			if refer.only == nil {
				return NIL
			}
			if _, ok := refer.only[sym]; !ok {
				return NIL
			}
		}
		return v
	}
	return NIL
}

func (n *Namespace) Refer(ns *Namespace, alias string, all bool) {
	nom := ns.Name()
	if alias != "" {
		nom = alias
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.refers[Symbol(nom)] = &Refer{
		all:  all,
		ns:   ns,
		only: nil,
	}
}

// ReferList refers only selected symbols from the given namespace into this namespace.
func (n *Namespace) ReferList(ns *Namespace, symbols []Symbol) {
	set := make(map[Symbol]bool, len(symbols))
	for _, s := range symbols {
		set[s] = true
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.refers[Symbol(ns.Name())] = &Refer{
		ns:   ns,
		all:  false,
		only: set,
	}
}

// Alias creates a symbol alias to another namespace in this namespace.
//
// Collision guard: re-pointing an existing alias at a DIFFERENT namespace is
// almost always a bug — either two `:as` clauses in one ns picking the same
// short name for different libs, or (historically) cross-namespace alias-table
// contamination during dependency loading. Rather than letting the last writer
// silently win, warn so the conflict is visible. Re-aliasing to the SAME target
// (idempotent ns reload) and first-time aliasing are silent.
func (n *Namespace) Alias(alias Symbol, target *Namespace) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if prev, ok := n.aliases[alias]; ok && prev != nil && target != nil && prev != target {
		fmt.Fprintf(os.Stderr,
			"WARNING: alias %s in namespace %s already points to %s, being repointed to %s\n",
			string(alias), n.name, string(prev.Name()), string(target.Name()))
	}
	n.aliases[alias] = target
}

// ImportVar links a var from another namespace into this namespace under the given alias.
// Returns true when the var exists and is not private.
func (n *Namespace) ImportVar(from *Namespace, name Symbol, alias Symbol) bool {
	v := from.localVar(name)
	if v == nil || v.isPrivate {
		return false
	}
	n.mu.Lock()
	n.registry[alias] = v
	n.mu.Unlock()
	return true
}

// ResolveAlias returns the namespace for the given alias, or nil.
func (n *Namespace) ResolveAlias(alias Symbol) *Namespace {
	t, _ := n.aliasFor(alias)
	return t
}

func (n *Namespace) Name() string {
	return n.name
}

func (n *Namespace) String() string {
	return fmt.Sprintf("<ns %s>", n.Name())
}

// registrySnapshot copies the registry so FuzzySymbolLookup can scan it (and
// recurse into refers) without holding the lock across other-namespace reads.
func (n *Namespace) registrySnapshot() map[Symbol]*Var {
	n.mu.RLock()
	defer n.mu.RUnlock()
	out := make(map[Symbol]*Var, len(n.registry))
	for k, v := range n.registry {
		out[k] = v
	}
	return out
}

// PublicVars snapshots the namespace's own non-private interned vars, keyed by
// their unqualified symbol. Backs the ns-publics core fn. isPrivate is a plain
// field read, so filtering under the read lock introduces no lock-order risk.
func (n *Namespace) PublicVars() map[Symbol]*Var {
	n.mu.RLock()
	defer n.mu.RUnlock()
	out := make(map[Symbol]*Var, len(n.registry))
	for k, v := range n.registry {
		if v.isPrivate {
			continue
		}
		out[k] = v
	}
	return out
}

func FuzzySymbolLookup(ns *Namespace, s Symbol, lookupPrivate bool) []Symbol {
	ret := []Symbol{}
	for _, r := range ns.refersSnapshot() {
		ret = append(ret, FuzzySymbolLookup(r.ns, s, false)...)
	}
	for k, v := range ns.registrySnapshot() {
		if strings.HasPrefix(string(k), string(s)) {
			if v.isPrivate && !lookupPrivate {
				continue
			}
			ret = append(ret, k)
		}
	}
	return ret
}
