/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "fmt"

// Protocol defines a named set of methods with type-based dispatch.
type Protocol struct {
	name    string
	methods []Symbol                     // method names
	impls   map[ValueType]*PersistentMap // type → {method-name → fn}
	nilImpl *PersistentMap               // implementation for nil
}

func NewProtocol(name string, methods []Symbol) *Protocol {
	return &Protocol{
		name:    name,
		methods: methods,
		impls:   make(map[ValueType]*PersistentMap),
	}
}

func (p *Protocol) Type() ValueType   { return ProtocolType }
func (p *Protocol) Unbox() any        { return p }
func (p *Protocol) String() string    { return fmt.Sprintf("<protocol %s>", p.name) }
func (p *Protocol) Name() string      { return p.name }
func (p *Protocol) Methods() []Symbol { return p.methods }

// Extend adds implementations for a type.
// implMap is a PersistentMap of {method-keyword → fn}.
func (p *Protocol) Extend(vt ValueType, implMap *PersistentMap) {
	p.impls[vt] = implMap
}

// ExtendNil adds implementations for nil.
func (p *Protocol) ExtendNil(implMap *PersistentMap) {
	p.nilImpl = implMap
}

// Lookup finds the implementation of a method for a given value's type.
func (p *Protocol) Lookup(methodName Symbol, target Value) (Fn, bool) {
	key := Keyword(methodName)

	if target == NIL {
		if p.nilImpl != nil {
			v := p.nilImpl.ValueAt(key)
			if v != NIL {
				if fn, ok := v.(Fn); ok {
					return fn, true
				}
			}
		}
		return nil, false
	}

	vt := target.Type()
	if fn, ok := p.lookupIn(p.impls[vt], key); ok {
		return fn, true
	}
	// Fall back to a default extended onto Object (AnyType) — Clojure's
	// (extend-type Object ...) universal default. Also covers a partial impl
	// that is missing this particular method.
	if vt != AnyType {
		if fn, ok := p.lookupIn(p.impls[AnyType], key); ok {
			return fn, true
		}
	}
	return nil, false
}

// lookupIn pulls a method fn out of one type's impl map, if present.
func (p *Protocol) lookupIn(implMap *PersistentMap, key Value) (Fn, bool) {
	if implMap == nil {
		return nil, false
	}
	v := implMap.ValueAt(key)
	if v == NIL {
		return nil, false
	}
	fn, ok := v.(Fn)
	return fn, ok
}

// Satisfies returns true if the given value's type has an implementation, or a
// universal default was extended onto Object (AnyType).
func (p *Protocol) Satisfies(target Value) bool {
	if target == NIL {
		return p.nilImpl != nil || p.impls[AnyType] != nil
	}
	if _, ok := p.impls[target.Type()]; ok {
		return ok
	}
	return p.impls[AnyType] != nil
}

// ProtocolFn is a function that dispatches on the first arg's type via a protocol.
type ProtocolFn struct {
	protocol   *Protocol
	methodName Symbol
	name       string
}

func NewProtocolFn(protocol *Protocol, methodName Symbol) *ProtocolFn {
	return &ProtocolFn{
		protocol:   protocol,
		methodName: methodName,
		name:       string(methodName),
	}
}

func (f *ProtocolFn) Type() ValueType { return FuncType }
func (f *ProtocolFn) Unbox() any      { return f }
func (f *ProtocolFn) String() string {
	return fmt.Sprintf("<protocol-fn %s/%s>", f.protocol.name, f.name)
}
func (f *ProtocolFn) Arity() int { return -1 }

func (f *ProtocolFn) Invoke(args []Value) (Value, error) {
	if len(args) == 0 {
		return NIL, fmt.Errorf("protocol fn %s requires at least one argument", f.name)
	}

	impl, ok := f.protocol.Lookup(f.methodName, args[0])
	if !ok {
		typeName := "nil"
		if args[0] != NIL {
			typeName = args[0].Type().Name()
		}
		return NIL, fmt.Errorf("no implementation of protocol %s method %s for type %s",
			f.protocol.name, f.name, typeName)
	}

	return impl.Invoke(args)
}

// Protocol type metadata

type theProtocolType struct{}

func (t *theProtocolType) String() string  { return t.Name() }
func (t *theProtocolType) Type() ValueType { return TypeType }
func (t *theProtocolType) Unbox() any      { return nil }
func (t *theProtocolType) Name() string    { return "let-go.lang.Protocol" }
func (t *theProtocolType) Box(bare any) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var ProtocolType *theProtocolType = &theProtocolType{}
