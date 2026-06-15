/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// TestEmbedField verifies that an embedded (anonymous) struct field is emitted
// without a field name, so the embedded type is promoted. This is the mechanism
// by which a generated record struct embeds vm.RecordBase.
func TestEmbedField(t *testing.T) {
	// Build: type Square struct { vm.RecordBase; side vm.Value }
	rbType := must(t)(cType(vm.String("vm.RecordBase")))
	embed := must(t)(cEmbedField(rbType))

	sideType := must(t)(cType(vm.String("vm.Value")))
	sideField := must(t)(cFieldDecl(vm.String("side"), sideType, vm.String("")))

	fields := vm.NewArrayVector([]vm.Value{embed, sideField})
	structType := must(t)(cStructType(fields))
	typeDecl := must(t)(cTypeDecl(vm.String("Square"), structType))

	got := render(t, typeDecl)

	if !strings.Contains(got, "type Square struct") {
		t.Errorf("got %q, expected to contain 'type Square struct'", got)
	}
	// Embedded field: the type appears with no preceding field name.
	if !strings.Contains(got, "vm.RecordBase") {
		t.Errorf("got %q, expected to contain embedded 'vm.RecordBase'", got)
	}
	if strings.Contains(got, "RecordBase vm.RecordBase") {
		t.Errorf("embedded field was emitted as a named field, not anonymous:\n%s", got)
	}
	if !strings.Contains(got, "side vm.Value") {
		t.Errorf("got %q, expected to contain 'side vm.Value'", got)
	}
}
