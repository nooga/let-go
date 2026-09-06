/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"bytes"
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TestDefMetaSurvivesLGBRoundTrip guards the fix for a real gap: defCompiler
// used to apply :doc/:private/:dynamic/:arglists/:file/:line/:column directly
// to the *vm.Var object live in the COMPILING process (e.g. cmd/lgbgen). That
// mutation never reached the .lgb wire format, so a namespace precompiled into
// core_compiled.lgb and loaded by a normal binary got a brand-new, meta-less
// stub Var at decode (see Namespace.DefStub) — every builtin's :doc/:arglists
// silently vanished in a normal build even though the same code compiled fresh
// (source-loaded) kept it.
//
// defCompiler now ALSO emits bytecode that invokes the already-existing
// apply-def-meta! builtin (the same one the Go-lowered/gogen def path uses) on
// the resolved var, mirroring how set-macro! already survives this exact
// round-trip for :macro. This test proves the metadata reaches a var that
// was NEVER touched by the original compiling process — only decoded and run,
// as core_compiled.lgb is in a shipped binary.
func TestDefMetaSurvivesLGBRoundTrip(t *testing.T) {
	const varName = "lgb-roundtrip-meta-test-fn"

	consts := vm.NewConsts()
	c := NewCompiler(consts, rt.NS(rt.NameCoreNS))
	chunk, err := c.Compile(`(defn ^{:private true :dynamic true} ` + varName + ` "a test doc" ([x] x) ([x y] (+ x y)))`)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	var buf bytes.Buffer
	if err := bytecode.EncodeCompilation(&buf, consts, chunk); err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	// Simulate decoding in a completely separate process: every var
	// reference resolves to a FRESH, never-before-seen *vm.Var — except
	// apply-def-meta! itself, which (like every other builtin) is registered
	// into "core" by a Go init() that always runs before any decode, so it
	// resolves to the real, working native fn.
	var decodedVar *vm.Var
	resolve := func(nsName, name string) *vm.Var {
		if name == "apply-def-meta!" {
			return rt.CoreNS.Lookup(vm.Symbol(name)).(*vm.Var)
		}
		v := vm.NewVar(nil, nsName, name)
		if name == varName {
			decodedVar = v
		}
		return v
	}

	unit, err := bytecode.DecodeToExecUnit(&buf, resolve)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if decodedVar == nil {
		t.Fatalf("resolver never saw a var-ref for %q", varName)
	}
	if decodedVar.Meta() != vm.NIL {
		t.Fatalf("decoded var already has meta before bytecode ran: %v", decodedVar.Meta())
	}

	f := vm.NewFrame(unit.MainChunk, nil)
	if _, err := f.RunProtected(); err != nil {
		t.Fatalf("running decoded chunk failed: %v", err)
	}
	vm.ReleaseFrame(f)

	if !decodedVar.IsPrivate() {
		t.Errorf(":private did not survive the round trip")
	}
	if !decodedVar.IsDynamic() {
		t.Errorf(":dynamic did not survive the round trip")
	}
	meta := decodedVar.Meta()
	m, ok := meta.(interface{ ValueAt(vm.Value) vm.Value })
	if !ok {
		t.Fatalf("decoded var has no meta map after running: %v", meta)
	}
	if doc := m.ValueAt(vm.Keyword("doc")); doc != vm.String("a test doc") {
		t.Errorf(":doc did not survive the round trip, got %v", doc)
	}
	if priv := m.ValueAt(vm.Keyword("private")); !vm.IsTruthy(priv) {
		t.Errorf(":private meta entry did not survive the round trip, got %v", priv)
	}
	arglists := m.ValueAt(vm.Keyword("arglists"))
	if arglists == vm.NIL {
		t.Fatalf(":arglists did not survive the round trip")
	}
	if got, want := arglists.String(), "([x] [x y])"; got != want {
		t.Errorf(":arglists = %s, want %s", got, want)
	}
	for key, want := range map[vm.Keyword]vm.Value{
		vm.Keyword("file"):   vm.String("<default>"),
		vm.Keyword("line"):   vm.Int(1),
		vm.Keyword("column"): vm.Int(1),
	} {
		if got := m.ValueAt(key); got != want {
			t.Errorf("%s = %v, want %v", key, got, want)
		}
	}
}

func TestDefMetaUsesCompactBundleConstant(t *testing.T) {
	consts := vm.NewConsts()
	c := NewCompiler(consts, rt.NS(rt.NameCoreNS))
	if _, err := c.Compile(`(defn ^:private compact-meta-constant "a test doc" [x] x)`); err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	var compact vm.DefMetaPairs
	for _, value := range consts.Values() {
		switch value := value.(type) {
		case *vm.PersistentMap:
			if value.ValueAt(vm.Keyword("doc")) == vm.String("a test doc") {
				t.Fatalf("def metadata remained a persistent-map constant: %v", value)
			}
		case vm.DefMetaPairs:
			for i := 0; i+1 < len(value); i += 2 {
				if value[i] == vm.Keyword("doc") && value[i+1] == vm.String("a test doc") {
					compact = value
				}
			}
		}
	}
	if compact == nil {
		t.Fatal("no compact def-metadata constant found")
	}
}

func TestCompactDefMetaRejectsNonMapSequence(t *testing.T) {
	entries := vm.ArrayVector{
		vm.MapEntry{Key: vm.Keyword("doc"), Value: vm.String("not a map")},
	}
	if compact, ok := compactDefMeta(entries); ok {
		t.Fatalf("non-map metadata sequence compacted as %v", compact)
	}
}
