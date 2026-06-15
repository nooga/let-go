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

// TestInterfaceType tests that an interface type with a method is correctly
// emitted as Go source code.
func TestInterfaceType(t *testing.T) {
	// Build: type Shape interface { Area() vm.Value }
	// Using: (gogen/type-decl "Shape" (gogen/interface-type [(gogen/iface-method "Area" [] ["vm.Value"])]))

	// First build the method: Area() vm.Value
	valueType := must(t)(cType(vm.String("vm.Value")))
	resultField := must(t)(cResult(valueType))
	resultFields := vm.NewArrayVector([]vm.Value{resultField})

	areaMethod := must(t)(cIfaceMethod(vm.String("Area"), vm.NewArrayVector(nil), resultFields))
	methods := vm.NewArrayVector([]vm.Value{areaMethod})

	// Build the interface type
	interfaceType := must(t)(cInterfaceType(methods))

	// Build the type declaration
	typeDecl := must(t)(cTypeDecl(vm.String("Shape"), interfaceType))

	got := render(t, typeDecl)

	// Assert the output contains the expected Go syntax
	if !strings.Contains(got, "type Shape interface") {
		t.Errorf("got %q, expected to contain 'type Shape interface'", got)
	}
	if !strings.Contains(got, "Area()") {
		t.Errorf("got %q, expected to contain 'Area()'", got)
	}
	if !strings.Contains(got, "vm.Value") {
		t.Errorf("got %q, expected to contain 'vm.Value'", got)
	}
}

// TestInterfaceTypeMultipleMethods tests an interface with multiple methods.
func TestInterfaceTypeMultipleMethods(t *testing.T) {
	// Build: type Reader interface { Read([]byte) (int, error); Close() error }

	// Build Read method: Read([]byte) (int, error)
	byteSliceType := must(t)(cType(vm.String("[]byte")))
	readParam := must(t)(cParam(vm.String("p"), byteSliceType))
	readParamFields := vm.NewArrayVector([]vm.Value{readParam})

	intType := must(t)(cType(vm.String("int")))
	errType := must(t)(cType(vm.String("error")))
	readResult1 := must(t)(cResult(intType))
	readResult2 := must(t)(cResult(errType))
	readResultFields := vm.NewArrayVector([]vm.Value{readResult1, readResult2})

	readMethod := must(t)(cIfaceMethod(vm.String("Read"), readParamFields, readResultFields))

	// Build Close method: Close() error
	closeMethod := must(t)(cIfaceMethod(vm.String("Close"), vm.NewArrayVector(nil), vm.NewArrayVector([]vm.Value{
		must(t)(cResult(must(t)(cType(vm.String("error"))))),
	})))

	methods := vm.NewArrayVector([]vm.Value{readMethod, closeMethod})

	// Build the interface type
	interfaceType := must(t)(cInterfaceType(methods))

	// Build the type declaration
	typeDecl := must(t)(cTypeDecl(vm.String("Reader"), interfaceType))

	got := render(t, typeDecl)

	// Assert the output contains the expected Go syntax
	if !strings.Contains(got, "type Reader interface") {
		t.Errorf("got %q, expected to contain 'type Reader interface'", got)
	}
	if !strings.Contains(got, "Read(") {
		t.Errorf("got %q, expected to contain 'Read('", got)
	}
	if !strings.Contains(got, "[]byte") {
		t.Errorf("got %q, expected to contain '[]byte'", got)
	}
	if !strings.Contains(got, "Close()") {
		t.Errorf("got %q, expected to contain 'Close()'", got)
	}
	if !strings.Contains(got, "error") {
		t.Errorf("got %q, expected to contain 'error'", got)
	}
}
