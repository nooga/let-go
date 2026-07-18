/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// BenchmarkProtocolFnInvoke measures a full protocol-method call on an
// extended type: ProtocolFn.Invoke → Protocol.Lookup → impl-map fetch →
// native fn invoke. The lookup path is where the per-call keyword lookup
// cost lives.
func BenchmarkProtocolFnInvoke(b *testing.B) {
	ident, err := NativeFnType.Wrap(func(vs []Value) (Value, error) {
		return vs[0], nil
	})
	if err != nil {
		b.Fatal(err)
	}
	impl := NewArrayMap([]Value{Keyword("-frob"), ident})
	p := NewProtocol("Frobber", []Symbol{Symbol("-frob")})
	p.Extend(IntType, impl)
	pf := NewProtocolFn(p, Symbol("-frob"))
	args := []Value{Int(7)}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, err := pf.Invoke(args)
		if err != nil || v != Int(7) {
			b.Fatalf("invoke: %v %v", v, err)
		}
	}
}

// BenchmarkProtocolLookupFallback measures the miss-then-AnyType-default
// path, which consulted two impl maps per call.
func BenchmarkProtocolLookupFallback(b *testing.B) {
	ident, err := NativeFnType.Wrap(func(vs []Value) (Value, error) {
		return vs[0], nil
	})
	if err != nil {
		b.Fatal(err)
	}
	impl := NewArrayMap([]Value{Keyword("-frob"), ident})
	p := NewProtocol("Frobber", []Symbol{Symbol("-frob")})
	p.Extend(AnyType, impl)
	target := String("s")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn, ok := p.Lookup(Symbol("-frob"), target)
		if !ok || fn == nil {
			b.Fatal("lookup failed")
		}
	}
}
