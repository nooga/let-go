/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// Lowered Go calls a value through rt.InvokeValueEC and its fixed-arity
// siblings. A value whose type implements the IFn protocol is callable on the
// VM, so it has to be callable there too: (filter pred xs) with a reify
// predicate reaches these helpers once core's filter is lowered.
func TestLoweredInvokeCallsIFnProtocolValues(t *testing.T) {
	adder := mustEval(t, `(reify IFn (-invoke [_ args] (apply + 100 args)))`)
	if _, ok := adder.(vm.Fn); ok {
		t.Fatalf("test premise: a reify IFn value is %T, already a vm.Fn", adder)
	}

	got, err := rt.InvokeValueEC(nil, adder, []vm.Value{vm.Int(1), vm.Int(2)})
	if err != nil || got != vm.Value(vm.Int(103)) {
		t.Fatalf("InvokeValueEC = %v, %v; want 103", got, err)
	}
	got, err = rt.InvokeValueEC(nil, adder, nil)
	if err != nil || got != vm.Value(vm.Int(100)) {
		t.Fatalf("InvokeValueEC with no args = %v, %v; want 100", got, err)
	}
	got, err = rt.InvokeValueEC1(nil, adder, vm.Int(1))
	if err != nil || got != vm.Value(vm.Int(101)) {
		t.Fatalf("InvokeValueEC1 = %v, %v; want 101", got, err)
	}
	got, err = rt.InvokeValueEC4(nil, adder, vm.Int(1), vm.Int(2), vm.Int(3), vm.Int(4))
	if err != nil || got != vm.Value(vm.Int(110)) {
		t.Fatalf("InvokeValueEC4 = %v, %v; want 110", got, err)
	}
}

// A var whose root is an IFn-protocol value is called through the var by
// lowered code (rt.CachedVarFnEC), as the VM calls it.
func TestCachedVarFnCallsIFnProtocolValue(t *testing.T) {
	mustEval(t, `(def lowered-invoke-ifn-var (reify IFn (-invoke [_ args] (count args))))`)
	var ptr *vm.Var
	fn := rt.CachedVarFnEC(nil, &ptr, rt.NameCoreNS, "lowered-invoke-ifn-var")
	got, err := fn.Invoke([]vm.Value{vm.Int(7), vm.Int(8)})
	if err != nil || got != vm.Value(vm.Int(2)) {
		t.Fatalf("call through the var = %v, %v; want 2", got, err)
	}
}

// A value that is not callable reports the VM's error, not a Go type name or
// a separate message for nil.
func TestLoweredInvokeOfNonFunctionMatchesVM(t *testing.T) {
	for _, tc := range []struct {
		src    string
		target vm.Value
	}{
		{"(5 1)", vm.Int(5)},
		{"(nil 1)", vm.NIL},
	} {
		_, vmErr := Eval(tc.src)
		if vmErr == nil {
			t.Fatalf("test premise: %s should fail on the VM", tc.src)
		}
		_, err := rt.InvokeValueEC(nil, tc.target, []vm.Value{vm.Int(1)})
		if err == nil {
			t.Fatalf("InvokeValueEC of %s succeeded", tc.target)
		}
		if strings.Contains(err.Error(), "vm.Int") || !strings.Contains(err.Error(), "is not a function") {
			t.Fatalf("InvokeValueEC of %s: %q; want the VM's wording (%q)", tc.target, err, vmErr)
		}
		if !strings.Contains(vmErr.Error(), "is not a function") {
			t.Fatalf("test premise: the VM's error for %s is %q", tc.src, vmErr)
		}
	}
}
