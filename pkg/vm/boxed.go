/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

type aBoxedType struct {
	typ     reflect.Type
	methods map[Symbol]*NativeFn
}

func (t *aBoxedType) String() string  { return t.Name() }
func (t *aBoxedType) Type() ValueType { return TypeType }
func (t *aBoxedType) Unbox() any      { return t.typ }

func (t *aBoxedType) Name() string { return "go." + t.typ.String() }
func (t *aBoxedType) Box(value any) (Value, error) {
	if !reflect.TypeOf(value).ConvertibleTo(t.typ) {
		return NIL, NewTypeError(value, "can't be boxed as", t)
	}
	return &Boxed{value: value, typ: t}, nil
}

type Boxed struct {
	value any
	typ   *aBoxedType
	// hash caches the String-based hash lazily. It packs the "computed" flag
	// into bit 32 so a legitimate hash of 0 still reads as cached: 0 means
	// uncomputed, otherwise the low 32 bits hold the hash. Atomic so
	// concurrent Hash() callers don't race (both compute the same value).
	hash atomic.Uint64
}

// Type implements Value
func (n *Boxed) Type() ValueType { return n.typ }

// Unbox implements Value
func (n *Boxed) Unbox() any { return n.value }

func (n *Boxed) String() string {
	return fmt.Sprintf("<%s %v>", n.typ.Name(), n.value)
}

// Hash implements Hashable. Without it, hashValue falls back to
// computeHash -> hashBytes(String()), which formats the boxed value via
// fmt.Sprintf("%v") on EVERY hash — catastrophic when boxed values (e.g.
// SourceInfo) are embedded in IR vectors that get hashed repeatedly during
// lowering (= comparisons, CSE, map/set keys). *Boxed is an immutable pointer,
// so we compute the String-based hash once and cache it. The value is
// identical to the old computeHash fallback, so map/set semantics are
// unchanged — only the repeated Sprintf is eliminated.
func (n *Boxed) Hash() uint32 {
	if packed := n.hash.Load(); packed != 0 {
		return uint32(packed)
	}
	h := hashString(n.String())
	n.hash.Store(uint64(1)<<32 | uint64(h))
	return h
}

func (n *Boxed) InvokeMethod(methodName Symbol, args []Value) (Value, error) {
	if n.typ.methods == nil {
		return NIL, fmt.Errorf("%v doesn't have any methods", n.typ)
	}
	method, ok := n.typ.methods[methodName]
	if !ok {
		return NIL, fmt.Errorf("method %s not found in %v", methodName, n.typ)
	}
	return method.Invoke(append([]Value{n}, args...))
}

func (n *Boxed) ValueAt(key Value) Value {
	return n.ValueAtOr(key, NIL)
}

func (n *Boxed) ValueAtOr(key Value, dflt Value) Value {
	name, ok := key.Unbox().(string)
	if !ok {
		return dflt
	}
	v, err := BoxValue(reflect.ValueOf(n.value).FieldByName(name))
	if err != nil {
		return dflt
	}
	return v
}

// BoxedType is the type of NilValues. Guarded by boxedTypesMu: boxing values
// from multiple goroutines (e.g. parallel lowering) otherwise races on this
// shared map, which is a runtime panic for concurrent writes, not just a
// detector warning.
var (
	boxedTypesMu sync.RWMutex
	BoxedTypes   map[string]*aBoxedType = map[string]*aBoxedType{}
)

func valueType(value any) *aBoxedType {
	reflected := reflect.TypeOf(value)
	// Use the full reflect.Type.String() (e.g. "*xxh3.Hasher") as the cache
	// key.  reflect.Type.Name() is empty for pointer/slice/map types, so
	// distinct pointer types would collide on the empty key and the first
	// boxed pointer would shadow every subsequent one.
	key := reflected.String()

	boxedTypesMu.RLock()
	t, ok := BoxedTypes[key]
	boxedTypesMu.RUnlock()
	if ok {
		return t
	}

	// Build the method table outside the write lock — reflection is the
	// expensive part and needs no mutual exclusion.
	t = &aBoxedType{
		typ:     reflected,
		methods: nil,
	}
	methodc := reflected.NumMethod()
	if methodc > 0 {
		t.methods = map[Symbol]*NativeFn{}
		for i := range methodc {
			m := reflected.Method(i)
			me, err := NativeFnType.Box(m.Func.Interface())
			if err != nil {
				fmt.Println(reflected.Name(), "boxing method failed", err)
				continue
			}
			mef, ok := me.(*NativeFn)
			if !ok {
				fmt.Println(reflected.Name(), "boxed method is not a native fn")
				continue
			}
			t.methods[Symbol(m.Name)] = mef
		}
	}

	// Recheck under the write lock: a concurrent caller may have inserted the
	// same key while we built our table. First writer wins so every caller
	// shares one *aBoxedType per type.
	boxedTypesMu.Lock()
	defer boxedTypesMu.Unlock()
	if existing, ok := BoxedTypes[key]; ok {
		return existing
	}
	BoxedTypes[key] = t
	return t
}

func NewBoxed(value any) *Boxed {
	return &Boxed{value: value, typ: valueType(value)}
}
