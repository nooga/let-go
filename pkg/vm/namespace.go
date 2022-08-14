/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
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

type Refer struct {
	ns  *Namespace
	all bool
}

type Namespace struct {
	name     string
	registry map[Symbol]*Var
	refers   map[Symbol]*Refer
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
	}
}

func (n *Namespace) Def(name string, val Value) *Var {
	s := Symbol(name)
	va := NewVar(n, n.name, name)
	va.SetRoot(val)
	if val.Type() == NativeFnType {
		val.(*NativeFn).SetName(name)
	}
	if val.Type() == FuncType {
		val.(*Func).SetName(name)
	}
	n.registry[s] = va
	return va
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
					return v
				}
			}
		}
		if v == nil {
			return NIL
		}
		return v
	}
	refer := n.refers[sns.(Symbol)]
	if refer == nil {
		return NIL
	}
	return refer.ns.registry[sym.(Symbol)]
}

func (n *Namespace) Refer(ns *Namespace, alias string, all bool) {
	nom := ns.Name()
	if alias != "" {
		nom = alias
	}
	n.refers[Symbol(nom)] = &Refer{
		all: all,
		ns:  ns,
	}
}

func (n *Namespace) Name() string {
	return n.name
}

func (n *Namespace) String() string {
	return fmt.Sprintf("<ns %s>", n.Name())
}

func FuzzySymbolLookup(ns *Namespace, s Symbol) []Symbol {
	ret := []Symbol{}
	for _, r := range ns.refers {
		ret = append(ret, FuzzySymbolLookup(r.ns, s)...)
	}
	for k := range ns.registry {
		if strings.HasPrefix(string(k), string(s)) {
			ret = append(ret, k)
		}
	}
	return ret
}
