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

// InternVar interns (LookupOrAdd) a Var by namespace and symbol name, creating
// the namespace and/or the var if absent, and returns it. Used by gogen-lowered
// `def` forms, which must CREATE the var at runtime — unlike LookupVar, which
// only resolves an already-existing one (returning nil otherwise).
func InternVar(nsName, symName string) *vm.Var {
	ns := NS(nsName)
	v := ns.LookupOrAdd(vm.Symbol(symName))
	out, _ := v.(*vm.Var)
	return out
}

// SetVarRoot sets the root of a var whose handle is typed vm.Value (the
// lowered :def path holds the InternVar result as vm.Value so a (def …) whose
// value is used downstream — e.g. returned from a var-initializer — keeps a
// vm.Value type). SetRoot is a *vm.Var method, so it can't be called on the
// interface directly; this casts. A non-Var handle is a no-op.
func SetVarRoot(v vm.Value, root vm.Value) {
	if vr, ok := v.(*vm.Var); ok {
		vr.SetRoot(root)
	}
}

// ApplyVarMetaV is ApplyVarMeta for a var handle typed vm.Value (see
// SetVarRoot). A non-Var handle is a no-op.
func ApplyVarMetaV(v vm.Value, meta vm.Value) {
	if vr, ok := v.(*vm.Var); ok {
		ApplyVarMeta(vr, meta)
	}
}

// ApplyVarMeta applies def metadata to v exactly as the bytecode defCompiler
// does: it sets the var's meta map, then mirrors the :dynamic / :private flags
// onto the Var. Shared by build-time def lowering (the apply-def-meta! builtin)
// and gogen-lowered `def` forms so both paths reproduce defCompiler's var setup
// (docstrings, type hints, ^:dynamic / ^:private). A nil/NIL meta is a no-op.
func ApplyVarMeta(v *vm.Var, meta vm.Value) {
	if v == nil || meta == nil || meta == vm.NIL {
		return
	}
	v.SetMeta(meta)
	m, ok := meta.(interface {
		ValueAt(vm.Value) vm.Value
	})
	if !ok {
		return
	}
	if vm.IsTruthy(m.ValueAt(vm.Keyword("dynamic"))) {
		v.SetDynamic()
	}
	if vm.IsTruthy(m.ValueAt(vm.Keyword("private"))) {
		v.SetPrivate()
	}
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

// CachedVarFn lazily resolves a cross-namespace Var ONCE — memoizing it through
// the caller-supplied package-level pointer — then returns that var's CURRENT
// fn root for invocation. The lowered Go calls it as
//
//	rt.CachedVarFn(&__v_ns_name, "ns", "name").Invoke(args)
//
// The *Var (not the fn) is what's cached, and its root is re-read on every call,
// so the var stays the override seam (^:dynamic / ^:redef / alter-var-root still
// work). Resolution is lazy (first call), not at package init(): blank-imported
// lowered packages run init() before the bundle replay loads namespaces, where
// LookupVar would return nil. The .(vm.Fn) assertion matches the trampoline's
// contract — every callable vm value implements vm.Fn.
func CachedVarFn(ptr **vm.Var, nsName, symName string) vm.Fn {
	if *ptr == nil {
		*ptr = LookupVar(nsName, symName)
	}
	return (*ptr).Deref().(vm.Fn)
}

// InvokeValue applies a runtime callable using let-go's dynamic invocation path.
func InvokeValue(target vm.Value, args []vm.Value) (vm.Value, error) {
	return InvokeValueEC(nil, target, args)
}

// InvokeValueEC is InvokeValue threaded with the caller's ExecContext, so a
// callable invoked from lowered code resolves dynamic vars (and propagates the
// scope to spawned goroutines) against the running context rather than the
// process root. A nil ec resolves to the root context. Values that are not
// vm.Fn (only the bare invoker interface) take the context-free path.
func InvokeValueEC(ec *vm.ExecContext, target vm.Value, args []vm.Value) (vm.Value, error) {
	if target == vm.NIL {
		return vm.NIL, fmt.Errorf("invoke of nil")
	}
	if fn, ok := target.(vm.Fn); ok {
		return ec.Invoke(fn, args)
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

// EqValue implements Clojure value-equality (the = operator / OP_EQ) for the
// lowered Go path. Unlike Go's ==, it is safe on non-comparable dynamic types
// (vectors, maps, seqs) and compares by value, not interface identity. The
// lowerer routes = through here whenever an operand is vm.Value (it could hold a
// collection); concrete comparable scalars keep native Go ==.
func EqValue(a, b vm.Value) bool {
	return vm.ValueEquals(a, b)
}

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

// DivValue is the runtime target for the :div IR op when the operands aren't
// both native float64 (int/int yields a Ratio, mixed yields Float). Mirrors
// clojure.core// via vm.NumDiv so the lowered Go path matches the interpreter.
func DivValue(a, b vm.Value) vm.Value {
	r, err := vm.NumDiv(a, b)
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
