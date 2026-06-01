/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

// IFnProtocol and IFnInvoke wire up "invokable values": any value whose type
// satisfies IFnProtocol can be called like a function, with the call dispatched
// to its -invoke method (the IFnInvoke protocol fn). They are set once by the
// runtime after core defines the IFn protocol; until then invoking a non-Fn
// value is a plain type error, as before.
var (
	IFnProtocol *Protocol
	IFnInvoke   Fn
	// IDerefProtocol / IDerefDeref wire "derefable values": any value whose type
	// satisfies IDeref can be used with @ / deref, dispatched to its -deref
	// method. Set once by the runtime after core defines the IDeref protocol.
	IDerefProtocol *Protocol
	IDerefDeref    Fn
)

// AsFn returns fraw as a callable Fn: directly if it already is one, or — when
// fraw's type satisfies the IFn protocol — wrapped in an adapter that routes
// calls to its -invoke method. Reports false when fraw is not callable.
func AsFn(fraw Value) (Fn, bool) {
	if fn, ok := fraw.(Fn); ok {
		return fn, true
	}
	if IFnProtocol != nil && IFnInvoke != nil && IFnProtocol.Satisfies(fraw) {
		return ifnAdapter{recv: fraw}, true
	}
	return nil, false
}

// ifnAdapter makes an IFn-satisfying value callable: a call (recv arg...) is
// dispatched to (-invoke recv [arg...]), the call's arguments passed as a
// single vector so one variadic -invoke method handles every arity.
type ifnAdapter struct{ recv Value }

func (a ifnAdapter) Invoke(args []Value) (Value, error) {
	return IFnInvoke.Invoke([]Value{a.recv, NewArrayVector(args)})
}
func (a ifnAdapter) Arity() int      { return -1 } // any arity (variadic -invoke)
func (a ifnAdapter) Type() ValueType { return a.recv.Type() }
func (a ifnAdapter) Unbox() any      { return a.recv.Unbox() }
func (a ifnAdapter) String() string  { return a.recv.String() }

// AsRef returns v as a Reference: directly if it already is one, or — when v's
// type satisfies the IDeref protocol — wrapped in an adapter that routes Deref
// to its -deref method. Reports false when v is not derefable.
func AsRef(v Value) (Reference, bool) {
	if r, ok := v.(Reference); ok {
		return r, true
	}
	if IDerefProtocol != nil && IDerefDeref != nil && IDerefProtocol.Satisfies(v) {
		return iderefAdapter{recv: v}, true
	}
	return nil, false
}

// iderefAdapter makes an IDeref-satisfying value derefable: Deref dispatches to
// (-deref recv). Reference.Deref has no error return, so an error thrown inside
// -deref surfaces via the VM's panic-recovery path (consistent with other
// Reference implementations).
type iderefAdapter struct{ recv Value }

func (a iderefAdapter) Deref() Value {
	v, err := IDerefDeref.Invoke([]Value{a.recv})
	if err != nil {
		panic(&thrownPanic{err: err})
	}
	return v
}
