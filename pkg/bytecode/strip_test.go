/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package bytecode

import (
	"bytes"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestStripDebugRemovesDebugSections(t *testing.T) {
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append32(int(vm.OP_LOAD_CONST))
	chunk.Append32(consts.Intern(vm.Int(42)))
	chunk.Append32(int(vm.OP_RETURN))
	chunk.SetMaxStack(1)
	chunk.AddSourceInfoAt(0, vm.SourceInfo{File: "test.lg", Line: 1, Column: 1, EndLine: 1, EndColumn: 5})
	chunk.AddLocalVar(0, "x")

	var enc bytes.Buffer
	if err := EncodeCompilation(&enc, consts, chunk); err != nil {
		t.Fatal(err)
	}
	fat := enc.Bytes()

	slim, err := StripDebug(fat)
	if err != nil {
		t.Fatal(err)
	}
	if len(slim) >= len(fat) {
		t.Fatalf("stripped bundle not smaller: %d >= %d", len(slim), len(fat))
	}

	m, err := Decode(bytes.NewReader(slim))
	if err != nil {
		t.Fatalf("stripped bundle does not decode: %v", err)
	}
	for i, c := range m.Chunks {
		if len(c.SourceMap) != 0 {
			t.Errorf("chunk %d retains %d source map entries", i, len(c.SourceMap))
		}
		if len(c.LocalVars) != 0 {
			t.Errorf("chunk %d retains %d localvar entries", i, len(c.LocalVars))
		}
	}
	if m.Flags&FlagLocalVars != 0 {
		t.Error("FlagLocalVars still set on stripped bundle")
	}
	if m.Flags&FlagDebugSplit == 0 {
		t.Error("FlagDebugSplit not set on stripped bundle")
	}

	// The stripped bundle must still execute.
	unit, err := DecodeToExecUnitBytes(slim, func(ns, name string) *vm.Var { return nil })
	if err != nil {
		t.Fatal(err)
	}
	f := vm.NewFrame(unit.MainChunk, nil)
	out, err := f.Run()
	vm.ReleaseFrame(f)
	if err != nil {
		t.Fatal(err)
	}
	if out != vm.Int(42) {
		t.Fatalf("stripped bundle ran to %v, want 42", out)
	}
}

func TestSplitDebugCompanionRestoresDebugInfo(t *testing.T) {
	encode := func(value vm.Value, file string) []byte {
		t.Helper()
		consts := vm.NewConsts()
		chunk := vm.NewCodeChunk(consts)
		chunk.Append32(int(vm.OP_LOAD_CONST))
		chunk.Append32(consts.Intern(value))
		chunk.Append32(int(vm.OP_RETURN))
		chunk.SetMaxStack(1)
		chunk.AddSourceInfoAt(0, vm.SourceInfo{File: file, Line: 4, Column: 2, EndLine: 4, EndColumn: 8})
		chunk.AddLocalVar(0, "answer")
		var buf bytes.Buffer
		if err := EncodeCompilation(&buf, consts, chunk); err != nil {
			t.Fatal(err)
		}
		return buf.Bytes()
	}

	fat := encode(vm.Int(42), "answer.lg")
	slim, debug, err := SplitDebug(fat)
	if err != nil {
		t.Fatal(err)
	}
	if len(debug) == 0 {
		t.Fatal("SplitDebug returned an empty companion")
	}

	unit, err := DecodeToExecUnitBytesWithDebug(slim, debug, func(ns, name string) *vm.Var { return nil })
	if err != nil {
		t.Fatal(err)
	}
	source := unit.MainChunk.LookupSource(0)
	if source == nil || source.File != "answer.lg" || source.Line != 4 || source.Column != 2 {
		t.Fatalf("source map was not restored: %#v", source)
	}
	locals := unit.MainChunk.LocalVars()
	if len(locals) != 1 || locals[0].Slot != 0 || locals[0].Name != "answer" {
		t.Fatalf("local vars were not restored: %#v", locals)
	}

	otherSlim, _, err := SplitDebug(encode(vm.Int(43), "other.lg"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeToExecUnitBytesWithDebug(otherSlim, debug, nil); err == nil ||
		!bytes.Contains([]byte(err.Error()), []byte("SHA-256 mismatch")) {
		t.Fatalf("mismatched companion error = %v, want SHA-256 mismatch", err)
	}

	// A stale sidecar beside a newly-built, unstripped artifact is harmless.
	plain, err := DecodeToExecUnitBytesWithDebug(fat, debug, nil)
	if err != nil {
		t.Fatalf("unstripped LGB rejected stale companion: %v", err)
	}
	if source := plain.MainChunk.LookupSource(0); source == nil || source.File != "answer.lg" {
		t.Fatalf("unstripped source map changed: %#v", source)
	}
}

func TestSplitDebugPreservesCompression(t *testing.T) {
	consts := vm.NewConsts()
	chunk := vm.NewCodeChunk(consts)
	chunk.Append32(int(vm.OP_LOAD_CONST))
	chunk.Append32(consts.Intern(vm.Int(42)))
	chunk.Append32(int(vm.OP_RETURN))
	chunk.SetMaxStack(1)
	chunk.AddSourceInfoAt(0, vm.SourceInfo{File: "compressed.lg", Line: 3, Column: 1})

	var encoded bytes.Buffer
	if err := EncodeCompilationCompressed(&encoded, consts, chunk, true); err != nil {
		t.Fatal(err)
	}
	slim, debug, err := SplitDebug(encoded.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	m, err := Decode(bytes.NewReader(slim))
	if err != nil {
		t.Fatal(err)
	}
	if m.Flags&FlagCompressed == 0 || m.Flags&FlagDebugSplit == 0 {
		t.Fatalf("split compressed flags = 0x%04x, want compression and split-debug", m.Flags)
	}
	unit, err := DecodeToExecUnitBytesWithDebug(slim, debug, nil)
	if err != nil {
		t.Fatal(err)
	}
	if source := unit.MainChunk.LookupSource(0); source == nil || source.File != "compressed.lg" {
		t.Fatalf("compressed companion source map was not restored: %#v", source)
	}
}

func TestStripDebugBundleFormat(t *testing.T) {
	consts := vm.NewConsts()
	mk := func(v vm.Value) *vm.CodeChunk {
		c := vm.NewCodeChunk(consts)
		c.Append32(int(vm.OP_LOAD_CONST))
		c.Append32(consts.Intern(v))
		c.Append32(int(vm.OP_RETURN))
		c.SetMaxStack(1)
		c.AddSourceInfoAt(0, vm.SourceInfo{File: "ns.lg", Line: 1, Column: 1, EndLine: 1, EndColumn: 3})
		c.AddLocalVar(0, "y")
		return c
	}
	nsChunks := map[string]*vm.CodeChunk{"a.ns": mk(vm.Int(1)), "main": mk(vm.Int(7))}

	var enc bytes.Buffer
	if err := EncodeBundleOrdered(&enc, consts, nsChunks, []string{"a.ns", "main"}); err != nil {
		t.Fatal(err)
	}
	slim, err := StripDebug(enc.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	m, err := Decode(bytes.NewReader(slim))
	if err != nil {
		t.Fatalf("stripped bundle does not decode: %v", err)
	}
	if len(m.NSTable) != 2 {
		t.Fatalf("NS table lost in strip: %v", m.NSTable)
	}
	for i, c := range m.Chunks {
		if len(c.SourceMap) != 0 || len(c.LocalVars) != 0 {
			t.Errorf("chunk %d retains debug sections", i)
		}
	}

	// Idempotent: stripping a stripped bundle is a no-op.
	again, err := StripDebug(slim)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(slim, again) {
		t.Errorf("strip is not idempotent: %d -> %d bytes", len(slim), len(again))
	}

	unit, err := DecodeToExecUnitBytes(slim, func(ns, name string) *vm.Var { return nil })
	if err != nil {
		t.Fatal(err)
	}
	f := vm.NewFrame(unit.MainChunk, nil)
	out, err := f.Run()
	vm.ReleaseFrame(f)
	if err != nil {
		t.Fatal(err)
	}
	if out != vm.Int(7) {
		t.Fatalf("stripped bundle main ran to %v, want 7", out)
	}
}
