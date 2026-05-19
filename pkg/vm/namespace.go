/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"strings"
)

type theNamespaceType struct{}

func (t *theNamespaceType) String() string     { return t.Name() }
func (t *theNamespaceType) Type() ValueType    { return TypeType }
func (t *theNamespaceType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theNamespaceType) Name() string { return "let-go.lang.Namespace" }
func (t *theNamespaceType) Box(fn interface{}) (Value, error) {
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
}

func (n *Namespace) Type() ValueType { return NamespaceType }

// Unbox implements Unbox
func (n *Namespace) Unbox() interface{} {
	return nil
}

func NewNamespace(name string) *Namespace {
	return &Namespace{
		name:     name,
		registry: map[Symbol]*Var{},
		refers:   map[Symbol]*Refer{},
		aliases:  map[Symbol]*Namespace{},
	}
}

func (n *Namespace) RegistrySize() int { return len(n.registry) }

func (n *Namespace) Def(name string, val Value) *Var {
	s := Symbol(name)
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

func (n *Namespace) LookupOrAdd(symbol Symbol) Value {
	val, ok := n.registry[symbol]
	if !ok {
		return n.Def(string(symbol), NIL)
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
		if v == nil || v.isPrivate {
			return NIL
		}
		return v
	}
	// Fallback: direct namespace lookup from global registry
	if nsLookup != nil {
		if target := nsLookup(string(sns.(Symbol))); target != nil {
			v := target.registry[sym.(Symbol)]
			if v != nil && !v.isPrivate {
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
