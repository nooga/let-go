/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TestCompilerRecordsLocalVarNames verifies the compiler records local binding
// names (slot -> name) onto the chunk's debug table, so they can be serialized
// and surfaced in crash traces.
func TestCompilerRecordsLocalVarNames(t *testing.T) {
	src := `(let [foo 1 bar 2] (+ foo bar))`
	cp := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := NewCompiler(cp, ns)

	chunk, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	names := map[string]bool{}
	for _, lv := range chunk.LocalVars() {
		names[lv.Name] = true
	}
	if !names["foo"] || !names["bar"] {
		t.Fatalf("expected local names foo and bar in chunk debug table, got %+v", chunk.LocalVars())
	}
}

// TestCompilerRecordsFnParamNames verifies function parameter names are captured
// onto the fn body chunk — the key case for naming locals in crash traces.
func TestCompilerRecordsFnParamNames(t *testing.T) {
	src := `(fn [aa bb] (+ aa bb))`
	cp := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := NewCompiler(cp, ns)
	_, result, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	fn, ok := result.(*vm.Func)
	if !ok {
		t.Fatalf("expected *vm.Func, got %T", result)
	}
	names := map[string]bool{}
	for _, lv := range fn.Chunk().LocalVars() {
		names[lv.Name] = true
	}
	if !names["aa"] || !names["bb"] {
		t.Fatalf("expected param names aa, bb on fn chunk, got %+v", fn.Chunk().LocalVars())
	}
}
