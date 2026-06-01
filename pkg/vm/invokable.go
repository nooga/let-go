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
