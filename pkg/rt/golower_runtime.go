/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

type invoker interface {
	Invoke([]vm.Value) (vm.Value, error)
}

// LookupVar resolves a runtime Var by namespace and symbol name.
func LookupVar(nsName, symName string) *vm.Var {
	ns := NS(nsName)
	if ns == nil {
		return nil
	}
	v := ns.Lookup(vm.Symbol(symName))
	if v == vm.NIL {
		return nil
	}
	out, _ := v.(*vm.Var)
	return out
}

// InvokeValue applies a runtime callable using let-go's dynamic invocation path.
func InvokeValue(target vm.Value, args []vm.Value) (vm.Value, error) {
	if target == vm.NIL {
		return vm.NIL, fmt.Errorf("invoke of nil")
	}
	inv, ok := target.(invoker)
	if !ok {
		return vm.NIL, fmt.Errorf("%T does not implement Invoke", target)
	}
	return inv.Invoke(args)
}

// BoxNativeFn wraps a Go function literal as a let-go callable.
// Lowered Go closures use this at the runtime boundary after capturing
// outer Go locals directly.
func BoxNativeFn(fn any) vm.Value {
	v, err := vm.NativeFnType.Box(fn)
	if err != nil {
		panic(err)
	}
	return v
}

// MakeNativeMultiArity combines native-lowered function branches into a
// runtime multi-arity callable.
func MakeNativeMultiArity(fns []vm.Value) vm.Value {
	v, err := vm.MakeMultiArity(fns)
	if err != nil {
		panic(err)
	}
	return v
}
