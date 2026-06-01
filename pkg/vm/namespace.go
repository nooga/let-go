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

type Refer struct {
	ns   *Namespace
	all  bool
	only map[Symbol]bool
}

type Namespace struct {
	name     string
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

// Exclude marks a symbol as excluded from clojure.core auto-refer.
// Called from the ns macro for :refer-clojure :exclude [...].
func (n *Namespace) Exclude(name string) {
	n.excludes[Symbol(name)] = true
}

// IsExcluded reports whether the symbol is in the exclude set.
func (n *Namespace) IsExcluded(name Symbol) bool {
	return n.excludes[name]
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

func (n *Namespace) RegistrySize() int { return len(n.registry) }

// isShadowingCoreRefer reports whether name `s` is currently visible
// unqualified in namespace `n` via a refer of clojure.core.
//
// Refer entries are keyed by namespace name (e.g. "clojure.core"), not
// by symbol — so we look up that single entry, then check whether `s`
// is in scope via :refer :all or :refer :only.
func isShadowingCoreRefer(n *Namespace, s Symbol) bool {
	for _, ref := range n.refers {
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
	// auto-refered :all, so it does warn on shadow.
	if coreNamespacePtr != nil && n != coreNamespacePtr && !n.excludes[s] {
		if isShadowingCoreRefer(n, s) {
			if existing, ok := coreNamespacePtr.registry[s]; ok && existing != nil && !existing.isPrivate {
				// Only warn the first time we shadow in this ns; subsequent
				// re-defs of our own var don't re-warn.
				if _, alreadyDefined := n.registry[s]; !alreadyDefined {
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
	n.registry[s] = va
	return va
}

// LookupLocal checks only the namespace's own registry, not refers or aliases.
func (n *Namespace) LookupLocal(symbol Symbol) *Var {
	return n.registry[symbol]
}

// DefStub creates a var with NIL root without triggering the warn-on-shadow
// check. Intended for bundle decoders that pre-populate var references
// before the namespace's own chunk runs (which would Def them properly).
// Do NOT use DefStub to intentionally suppress warnings for new code; use
// Namespace.Exclude (via :refer-clojure :exclude) instead.
func (n *Namespace) DefStub(name string) *Var {
	s := Symbol(name)
	va := NewVar(n, n.name, name)
	va.SetRoot(NIL)
	n.registry[s] = va
	return va
}

func (n *Namespace) LookupOrAdd(symbol Symbol) Value {
	val, ok := n.registry[symbol]
	if !ok {
		// Intern an UNBOUND var (no root) rather than Def(NIL): the compiler
		// calls this while compiling a `(def x v)` form before it runs, and
		// `defonce` must be able to tell that interned-but-unrun state from a
		// var that has actually been assigned. Deref still yields NIL.
		va := NewVar(n, n.name, string(symbol))
		n.registry[symbol] = va
		return va
	}
	return val
}

func (n *Namespace) Lookup(symbol Symbol) Value {
	sns, sym := symbol.Namespaced()
	if sns == NIL {
		v := n.registry[sym.(Symbol)]
		if v == nil {
			for _, ref := range n.refers {
				v = ref.ns.registry[sym.(Symbol)]
				if v != nil {
					if v.isPrivate {
						return NIL
					}
					return v
				}
			}
		}
		if v == nil {
			return NIL
		}
		return v
	}
	// Alias-qualified resolution via aliases
	if target, ok := n.aliases[sns.(Symbol)]; ok {
		v := target.registry[sym.(Symbol)]
		if v == nil && nsLookup != nil {
			// Alias may point to a placeholder namespace created before source
			// load completed. Re-resolve by name so runtime loader can
			// materialize the namespace on demand, then retry the symbol lookup.
			if loaded := nsLookup(target.Name()); loaded != nil {
				target = loaded
				n.aliases[sns.(Symbol)] = loaded
				v = target.registry[sym.(Symbol)]
			}
		}
		if v == nil || v.isPrivate {
			return NIL
		}
		return v
	}
	// Fallback: direct namespace lookup from global registry
	if nsLookup != nil {
		if target := nsLookup(string(sns.(Symbol))); target != nil {
			v := target.registry[sym.(Symbol)]
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
	if refer, ok := n.refers[sns.(Symbol)]; ok {
		v := refer.ns.registry[sym.(Symbol)]
		if v == nil || v.isPrivate {
			return NIL
		}
		if !refer.all {
			if refer.only == nil {
				return NIL
			}
			if _, ok := refer.only[sym.(Symbol)]; !ok {
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
	n.refers[Symbol(ns.Name())] = &Refer{
		ns:   ns,
		all:  false,
		only: set,
	}
}

// Alias creates a symbol alias to another namespace in this namespace.
func (n *Namespace) Alias(alias Symbol, target *Namespace) {
	n.aliases[alias] = target
}

// ImportVar links a var from another namespace into this namespace under the given alias.
// Returns true when the var exists and is not private.
func (n *Namespace) ImportVar(from *Namespace, name Symbol, alias Symbol) bool {
	v := from.registry[name]
	if v == nil || v.isPrivate {
		return false
	}
	n.registry[alias] = v
	return true
}

// ResolveAlias returns the namespace for the given alias, or nil.
func (n *Namespace) ResolveAlias(alias Symbol) *Namespace {
	return n.aliases[alias]
}

func (n *Namespace) Name() string {
	return n.name
}

func (n *Namespace) String() string {
	return fmt.Sprintf("<ns %s>", n.Name())
}

func FuzzySymbolLookup(ns *Namespace, s Symbol, lookupPrivate bool) []Symbol {
	ret := []Symbol{}
	for _, r := range ns.refers {
		ret = append(ret, FuzzySymbolLookup(r.ns, s, false)...)
	}
	for k := range ns.registry {
		if strings.HasPrefix(string(k), string(s)) {
			if ns.registry[k].isPrivate && !lookupPrivate {
				continue
			}
			ret = append(ret, k)
		}
	}
	return ret
}
