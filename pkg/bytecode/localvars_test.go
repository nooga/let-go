/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package bytecode

import (
	"bytes"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// TestLocalVarTableRoundtrip verifies that a chunk's local-variable debug table
// (slot -> original name) survives the bundle encode/decode round-trip, so crash
// traces from bundle-loaded code (e.g. under WASM) can name local variables.
func TestLocalVarTableRoundtrip(t *testing.T) {
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk.SetMaxStack(1)
	chunk.AddLocalVar(0, "x")
	chunk.AddLocalVar(1, "y")

	b := NewModuleBuilder()
	b.AddChunk(chunk)
	m := b.Build()

	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	decoded, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	lv := decoded.Chunks[0].LocalVars
	if len(lv) != 2 {
		t.Fatalf("expected 2 local vars, got %d (%+v)", len(lv), lv)
	}
	if lv[0].Slot != 0 || lv[0].Name != "x" {
		t.Errorf("lv[0]: got %+v want {0 x}", lv[0])
	}
	if lv[1].Slot != 1 || lv[1].Name != "y" {
		t.Errorf("lv[1]: got %+v want {1 y}", lv[1])
	}
}

// TestLocalVarTableReconstructsOnChunk verifies the decoded vm.CodeChunk carries
// the local-variable table back (not just the raw ChunkData), so a running VM
// can resolve slot -> name.
func TestLocalVarTableReconstructsOnChunk(t *testing.T) {
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append(vm.OP_LOAD_CONST, 0, vm.OP_RETURN)
	chunk.SetMaxStack(1)
	chunk.AddLocalVar(2, "acc")

	b := NewModuleBuilder()
	b.AddChunk(chunk)
	m := b.Build()

	var buf bytes.Buffer
	if err := Encode(&buf, m); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	unit, err := DecodeToExecUnit(&buf, func(ns, name string) *vm.Var { return nil })
	if err != nil {
		t.Fatalf("DecodeToExecUnit: %v", err)
	}
	lv := unit.MainChunk.LocalVars()
	if len(lv) != 1 || lv[0].Slot != 2 || lv[0].Name != "acc" {
		t.Fatalf("reconstructed chunk local vars wrong: %+v", lv)
	}
}
