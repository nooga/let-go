/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TestTypeConstantsSurviveRepeatedLoad loads one compiled unit twice in one
// process. Each load decodes the unit's defrecord or deftype constant before
// its defining form runs again, and the constant must keep the type that form
// leaves defined. Compiling in process defines the type first, so the first
// load already finds it registered.
func TestTypeConstantsSurviveRepeatedLoad(t *testing.T) {
	for _, tc := range []struct{ name, ns, src string }{
		{"defrecord", "reload.record", `(ns reload.record)
(defrecord Pt [x])
(defmacro pt-const [] (list 'quote (->Pt 7)))
(def ok [(instance? Pt (pt-const)) (= (pt-const) (->Pt 7))])`},
		{"deftype", "reload.dtype", `(ns reload.dtype)
(deftype Box [v])
(defmacro box-const [] (list 'quote (Box. 7)))
(def ok [(instance? Box (box-const)) (= Box (type (box-const)))])`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer rt.CurrentNS.SetRoot(rt.CurrentNS.Deref())
			consts := vm.NewConsts()
			c := NewCompiler(consts, rt.NS(rt.NameCoreNS))
			chunk, _, err := c.CompileMultiple(strings.NewReader(tc.src))
			if err != nil {
				t.Fatalf("compile: %v", err)
			}
			var buf bytes.Buffer
			if err := bytecode.EncodeCompilation(&buf, consts, chunk); err != nil {
				t.Fatalf("encode: %v", err)
			}
			for load := 1; load <= 2; load++ {
				unit, err := rt.DecodeExecUnit(buf.Bytes())
				if err != nil {
					t.Fatalf("load %d: decode: %v", load, err)
				}
				if err := rt.RunExecUnit(unit); err != nil {
					t.Fatalf("load %d: run: %v", load, err)
				}
				ok := nsVar(tc.ns, "ok")
				if got := ok.String(); got != "[true true]" {
					t.Errorf("load %d: got %s, want [true true]", load, got)
				}
			}
		})
	}
}

// compileUnit compiles src in process, which also evaluates it, and returns the
// encoded unit.
func compileUnit(t *testing.T, src string) []byte {
	t.Helper()
	consts := vm.NewConsts()
	c := NewCompiler(consts, rt.NS(rt.NameCoreNS))
	chunk, _, err := c.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	var buf bytes.Buffer
	if err := bytecode.EncodeCompilation(&buf, consts, chunk); err != nil {
		t.Fatalf("encode: %v", err)
	}
	return buf.Bytes()
}

func nsVar(ns, name string) vm.Value {
	return rt.NS(ns).Lookup(vm.Symbol(name)).(*vm.Var).Deref()
}

// Decoding a unit that quotes a deftype instance, without running it, must not
// make a later source deftype keep the old type and its protocol impls.
func TestDecodeAloneLeavesRedefinitionNew(t *testing.T) {
	defer rt.CurrentNS.SetRoot(rt.CurrentNS.Deref())
	data := compileUnit(t, `(ns reload.decode-only)
(defprotocol P (m [x]))
(deftype Box [v] P (m [_] :old-impl))
(defmacro box-const [] (list 'quote (Box. 1)))
(def c (box-const))`)
	if _, err := rt.DecodeExecUnit(data); err != nil {
		t.Fatalf("decode: %v", err)
	}
	compileUnit(t, `(ns reload.decode-only)
(deftype Box [v])
(def r (satisfies? P (->Box 1)))`)
	if got := nsVar("reload.decode-only", "r"); got != vm.FALSE {
		t.Fatalf("the redefined Box still satisfies P: got %v", got)
	}
}

// A unit built from a record with other fields than the live one keeps one
// type for its constant and its defrecord.
func TestFieldChangeBetweenLoads(t *testing.T) {
	defer rt.CurrentNS.SetRoot(rt.CurrentNS.Deref())
	data := compileUnit(t, `(ns reload.fields)
(defrecord Pt [x y z])
(defmacro pt-const [] (list 'quote (->Pt 1 2 3)))
(def ok [(instance? Pt (pt-const)) (= (pt-const) (->Pt 1 2 3))])`)
	compileUnit(t, `(ns reload.fields)
(defrecord Pt [x y])`)
	unit, err := rt.DecodeExecUnit(data)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if err := rt.RunExecUnit(unit); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got := nsVar("reload.fields", "ok").String(); got != "[true true]" {
		t.Fatalf("got %s, want [true true]", got)
	}
}
