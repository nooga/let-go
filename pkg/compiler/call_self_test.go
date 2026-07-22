package compiler

import (
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallSelf_FibCorrectness(t *testing.T) {
	_, err := Eval(`(defn call-self-fib [n]
		(if (< n 2) n (+ (call-self-fib (- n 1)) (call-self-fib (- n 2)))))`)
	require.NoError(t, err)

	out, err := Eval(`(call-self-fib 10)`)
	require.NoError(t, err)
	assert.Equal(t, 55, out.Unbox())
}

func TestCallSelf_ObservesRedefinedVar(t *testing.T) {
	_, err := Eval(`(defn call-self-redef [n]
		(if (= n 0)
			0
			(do
				(def call-self-redef (fn [_] 99))
				(+ 1 (call-self-redef 0)))))`)
	require.NoError(t, err)

	out, err := Eval(`(call-self-redef 1)`)
	require.NoError(t, err)
	assert.Equal(t, 100, out.Unbox(), "self-call must invoke the Var's replacement root")
}

func TestCallSelf_NotEmittedForLexicallyNestedDefn(t *testing.T) {
	_, err := Eval(`(let [captured 42]
		(defn call-self-capturing [n]
			(if (= n 0)
				(+ 1 (call-self-capturing 1))
				captured)))`)
	require.NoError(t, err)

	out, err := Eval(`(call-self-capturing 0)`)
	require.NoError(t, err)
	assert.Equal(t, 43, out.Unbox())
}

func TestCallSelf_NotEmittedForFunctionNestedInDefInitializer(t *testing.T) {
	_, err := Eval(`(def call-self-container
		[(fn [n]
			(if (= n 0)
				10
				(+ 1 (call-self-container 0))))])`)
	require.NoError(t, err)

	_, err = Eval(`((first call-self-container) 1)`)
	require.Error(t, err, "call must resolve call-self-container through its vector root")
}

func TestCallSelf_EmittedForNonTailSelfRecursion(t *testing.T) {
	_, err := Eval(`(defn call-self-fib2 [n]
		(if (< n 2) n (+ (call-self-fib2 (- n 1)) (call-self-fib2 (- n 2)))))`)
	require.NoError(t, err)

	ns := rt.CurrentNS.Deref().(*vm.Namespace)
	v := ns.Lookup(vm.Symbol("call-self-fib2"))
	require.NotEqual(t, vm.NIL, v)
	fn, ok := v.(*vm.Var).Deref().(*vm.Func)
	require.True(t, ok, "expected *Func root")

	code := fn.Chunk().Code()
	callSelf, invoke, loadVar := 0, 0, 0
	for i := 0; i < len(code); {
		op := code[i] & 0xff
		switch op {
		case vm.OP_CALL_SELF:
			callSelf++
		case vm.OP_INVOKE:
			invoke++
		case vm.OP_LOAD_VAR:
			loadVar++
		}
		i += testOpcodeStride(op)
	}
	assert.Equal(t, 2, callSelf, "both recursive arms should be CALL_SELF")
	assert.Equal(t, 0, invoke, "no generic INVOKE in fib body")
	assert.Equal(t, 0, loadVar, "no LOAD_VAR for self-calls")
}

func TestCallSelf_NotEmittedForTailSelfCall(t *testing.T) {
	// Tail self-calls keep LOAD_VAR + TAIL_CALL (already frame-reusing).
	_, err := Eval(`(defn call-self-countdown [n]
		(if (< n 1) n (call-self-countdown (- n 1))))`)
	require.NoError(t, err)

	ns := rt.CurrentNS.Deref().(*vm.Namespace)
	v := ns.Lookup(vm.Symbol("call-self-countdown"))
	fn := v.(*vm.Var).Deref().(*vm.Func)
	code := fn.Chunk().Code()
	callSelf, tailCall := 0, 0
	for i := 0; i < len(code); {
		op := code[i] & 0xff
		switch op {
		case vm.OP_CALL_SELF:
			callSelf++
		case vm.OP_TAIL_CALL:
			tailCall++
		}
		i += testOpcodeStride(op)
	}
	assert.Equal(t, 0, callSelf)
	assert.Equal(t, 1, tailCall)
}

func TestCallSelf_NotEmittedForNestedLambda(t *testing.T) {
	_, err := Eval(`(defn call-self-outer [n]
		((fn [x] (call-self-outer x)) n))`)
	require.NoError(t, err)

	ns := rt.CurrentNS.Deref().(*vm.Namespace)
	v := ns.Lookup(vm.Symbol("call-self-outer"))
	fn := v.(*vm.Var).Deref().(*vm.Func)
	// Outer chunk should not CALL_SELF; the nested lambda's call stays LOAD_VAR.
	code := fn.Chunk().Code()
	for i := 0; i < len(code); {
		op := code[i] & 0xff
		assert.NotEqual(t, int32(vm.OP_CALL_SELF), op)
		i += testOpcodeStride(op)
	}
}

func testOpcodeStride(op int32) int {
	switch op & 0xff {
	case vm.OP_TRY_PUSH:
		return 3
	case vm.OP_CALL_SELF:
		return 3
	case vm.OP_RECUR:
		return 4
	case vm.OP_LOAD_ARG, vm.OP_BRANCH_TRUE, vm.OP_BRANCH_FALSE, vm.OP_JUMP,
		vm.OP_POP_N, vm.OP_DUP_NTH, vm.OP_INVOKE, vm.OP_LOAD_CLOSEDOVER,
		vm.OP_RECUR_FN, vm.OP_MAKE_MULTI_ARITY, vm.OP_TAIL_CALL,
		vm.OP_LOAD_CONST, vm.OP_LOAD_VAR, vm.OP_FINALLY_END:
		return 2
	default:
		return 1
	}
}
