/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import "testing"

// cljs/JVM-standard mutable deftype: ^:mutable fields are lexically in scope in
// method bodies (bare reads) and mutated with (set! field val). Read-after-write
// must observe the new value within the same method.
func TestMutableDeftype_BareReadAndSetBang(t *testing.T) {
	src := `(do
              (defprotocol Counter (bump [this]) (value [this]))
              (deftype Ctr [^:mutable n]
                Counter
                (bump [this] (set! n (inc n)) n)
                (value [this] n))
              (def c (->Ctr 0))
              (bump c) (bump c) (bump c)
              (= 3 (value c)))`
	val, err := runSource(t, src)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if val.String() != "true" {
		t.Fatalf("expected true (3 bumps), got %v", val)
	}
}

// A let-bound local shadows a field of the same name inside its scope.
func TestMutableDeftype_LetShadowsField(t *testing.T) {
	src := `(do
              (defprotocol P (m [this]))
              (deftype T [^:mutable n]
                P
                (m [this] (set! n 10) (let [n 99] n)))
              (= 99 (m (->T 0))))`
	val, err := runSource(t, src)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if val.String() != "true" {
		t.Fatalf("let local should shadow field n (want 99), got %v", val)
	}
}

// A let init expression on the RHS still sees the field (outer scope).
func TestMutableDeftype_LetInitSeesField(t *testing.T) {
	src := `(do
              (defprotocol P (m [this]))
              (deftype T [v]
                P
                (m [this] (let [x v] x)))
              (= 5 (m (->T 5))))`
	val, err := runSource(t, src)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if val.String() != "true" {
		t.Fatalf("let init RHS should read field v (want 5), got %v", val)
	}
}

// A fn parameter shadows a field of the same name in the fn body.
func TestMutableDeftype_FnParamShadowsField(t *testing.T) {
	src := `(do
              (defprotocol P (m [this]))
              (deftype T [v]
                P
                (m [this] ((fn [v] v) 7)))
              (= 7 (m (->T 100))))`
	val, err := runSource(t, src)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if val.String() != "true" {
		t.Fatalf("fn param should shadow field v (want 7), got %v", val)
	}
}

// An immutable field is also lexically in scope for bare reads in method bodies.
func TestMutableDeftype_ImmutableFieldReadable(t *testing.T) {
	src := `(do
              (defprotocol Boxed (unbox [this]))
              (deftype Box [v] Boxed (unbox [this] v))
              (= 42 (unbox (->Box 42))))`
	val, err := runSource(t, src)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if val.String() != "true" {
		t.Fatalf("expected true, got %v", val)
	}
}
