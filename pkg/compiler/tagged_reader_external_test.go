/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/api"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func readAngleFragment(input compiler.TaggedRawInput) (vm.Value, error) {
	opener, err := input.ReadRune()
	if err != nil {
		return vm.NIL, err
	}
	if opener != '[' {
		return vm.NIL, fmt.Errorf("raw fragment requires [")
	}
	var body strings.Builder
	for {
		ch, err := input.ReadRune()
		if err != nil {
			return vm.NIL, err
		}
		if ch == ']' {
			return vm.String(body.String()), nil
		}
		body.WriteRune(ch)
	}
}

func TestExternalPackageCanImplementRawTaggedReader(t *testing.T) {
	registry := compiler.NewTaggedReaderRegistry()
	if err := registry.RegisterRaw("raw", readAngleFragment); err != nil {
		t.Fatal(err)
	}
	reader := compiler.NewLispReaderWithTaggedReaders(strings.NewReader("#raw[hello] 42"), "probe.lg", registry)
	got, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if got != vm.String("hello") {
		t.Fatalf("raw fragment = %q, want hello", got)
	}
	next, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if next != vm.Int(42) {
		t.Fatalf("next form = %v, want 42", next)
	}
}

func TestExternalPackageCompilerUsesRegisteredDataReader(t *testing.T) {
	registry := compiler.NewTaggedReaderRegistry()
	if err := registry.RegisterData("external", func(value vm.Value) (vm.Value, error) {
		if value.String() != "[1 2 3]" {
			t.Fatalf("data reader input = %v, want [1 2 3]", value)
		}
		return vm.String("compiled-custom-data"), nil
	}); err != nil {
		t.Fatal(err)
	}

	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS)).SetTaggedReaders(registry)
	_, got, err := ctx.CompileMultiple(strings.NewReader("#external [1 2 3]"))
	if err != nil {
		t.Fatal(err)
	}
	if got != vm.String("compiled-custom-data") {
		t.Fatalf("compiled custom data reader result = %v, want compiled-custom-data", got)
	}
}

func TestTruncatedTaggedLiteralIsIncompleteThroughPublicAPIs(t *testing.T) {
	registry := compiler.NewTaggedReaderRegistry()
	if err := registry.RegisterData("external", func(value vm.Value) (vm.Value, error) {
		return value, nil
	}); err != nil {
		t.Fatal(err)
	}

	_, err := compiler.NewLispReaderWithTaggedReaders(strings.NewReader("#external"), "external.lg", registry).Read()
	if err == nil {
		t.Fatal("truncated tagged literal returned nil error")
	}
	if !compiler.IsErrorEOF(err) {
		t.Fatalf("compiler.IsErrorEOF(%v) = false", err)
	}
	if !api.IsIncomplete(err) {
		t.Fatalf("api.IsIncomplete(%v) = false", err)
	}
}
