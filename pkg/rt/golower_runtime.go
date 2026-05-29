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

// Polymorphic binary-op helpers used by lower-go when an operand is
// vm.Value-typed (type inference didn't narrow to a primitive). Each
// mirrors the bytecode VM's OP_<X> handler: vm.Int/vm.Int fast path,
// then vm.Num<X> fallback. NumX errors panic to match bytecode behavior
// when no error handler is installed.

func LtValue(a, b vm.Value) bool {
	if ai, ok := a.(vm.Int); ok {
		if bi, ok := b.(vm.Int); ok {
			return int64(ai) < int64(bi)
		}
	}
	r, err := vm.NumLt(a, b)
	if err != nil {
		panic(err)
	}
	return r
}

func LeValue(a, b vm.Value) bool {
	if ai, ok := a.(vm.Int); ok {
		if bi, ok := b.(vm.Int); ok {
			return int64(ai) <= int64(bi)
		}
	}
	r, err := vm.NumLe(a, b)
	if err != nil {
		panic(err)
	}
	return r
}

func GtValue(a, b vm.Value) bool {
	if ai, ok := a.(vm.Int); ok {
		if bi, ok := b.(vm.Int); ok {
			return int64(ai) > int64(bi)
		}
	}
	r, err := vm.NumGt(a, b)
	if err != nil {
		panic(err)
	}
	return r
}

func GeValue(a, b vm.Value) bool {
	if ai, ok := a.(vm.Int); ok {
		if bi, ok := b.(vm.Int); ok {
			return int64(ai) >= int64(bi)
		}
	}
	r, err := vm.NumGe(a, b)
	if err != nil {
		panic(err)
	}
	return r
}

// Arithmetic helpers always delegate to vm.Num<X> to keep overflow,
// BigInt promotion, and ratio semantics in one place (per design call:
// fast-path-then-fallback is for ordering only; arithmetic prefers
// zero-divergence over a single Int/Int branch).

func AddValue(a, b vm.Value) vm.Value {
	r, err := vm.NumAdd(a, b)
	if err != nil {
		panic(err)
	}
	return r
}

func SubValue(a, b vm.Value) vm.Value {
	r, err := vm.NumSub(a, b)
	if err != nil {
		panic(err)
	}
	return r
}

func MulValue(a, b vm.Value) vm.Value {
	r, err := vm.NumMul(a, b)
	if err != nil {
		panic(err)
	}
	return r
}

func QuotValue(a, b vm.Value) vm.Value {
	r, err := vm.NumQuot(a, b)
	if err != nil {
		panic(err)
	}
	return r
}

// BoxRestArgs boxes a variadic rest-args slice into a vm.Value list.
// Used by the Go lowering when lowering :load-arg for the rest arg of
// a variadic function (where the Go param is ...vm.Value).
func BoxRestArgs(args []vm.Value) vm.Value {
	v, _ := vm.ListType.Box(args)
	return v
}
