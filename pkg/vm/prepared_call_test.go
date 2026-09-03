package vm

import (
	"testing"
)

// testBranchTailOrReturnFn builds a unary fn: truthy arg -> return the arg,
// falsy arg -> tail-call tail (zero args). The tail call exercises
// installBytecodeCall's same-frame rebind on the PreparedCall's owned frame.
func testBranchTailOrReturnFn(consts *Consts, tail Fn) *Func {
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_ARG)
	chunk.Append32(0)
	chunk.Append(OP_BRANCH_TRUE)
	chunk.Append32(7) // -> the LOAD_ARG/RETURN pair below
	chunk.Append(OP_LOAD_CONST)
	chunk.Append32(consts.Intern(tail))
	chunk.Append(OP_TAIL_CALL)
	chunk.Append32(0)
	chunk.Append(OP_RETURN)
	chunk.Append(OP_LOAD_ARG)
	chunk.Append32(0)
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)
	return MakeFunc(1, false, chunk)
}

func testConstFn(consts *Consts, v Value) *Func {
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_CONST)
	chunk.Append32(consts.Intern(v))
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)
	return MakeFunc(0, false, chunk)
}

func TestPreparedCallRepeatedInvocation(t *testing.T) {
	consts := NewConsts()
	identity := testBytecodeFnReturningArg(consts, 1)
	pc := RootExecContext.PrepareCall(identity, 1)
	if pc == nil {
		t.Fatal("PrepareCall rejected a plain fixed-arity Func")
	}
	defer pc.Release()
	for i := 0; i < 3; i++ {
		v, err := pc.Call1(Int(i))
		if err != nil {
			t.Fatal(err)
		}
		if v != Int(i) {
			t.Fatalf("call %d returned %v", i, v)
		}
	}
}

func TestPreparedCallPreservesClosureCaptures(t *testing.T) {
	consts := NewConsts()
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_CLOSEDOVER)
	chunk.Append32(0)
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)
	closure := &Closure{fn: MakeFunc(1, false, chunk), closedOvers: []Value{Int(99)}}

	pc := RootExecContext.PrepareCall(closure, 1)
	if pc == nil {
		t.Fatal("PrepareCall rejected a closure over a fixed-arity Func")
	}
	defer pc.Release()
	for i := 0; i < 2; i++ {
		v, err := pc.Call1(NIL)
		if err != nil {
			t.Fatal(err)
		}
		if v != Int(99) {
			t.Fatalf("call %d lost the capture: %v", i, v)
		}
	}
}

// A tail call in the callee rebinds the owned frame's code via
// installBytecodeCall; the next prepared call must run the original callable
// again, not the tail target.
func TestPreparedCallResetsAfterTailCallRebind(t *testing.T) {
	consts := NewConsts()
	tail := testConstFn(consts, Int(42))
	fn := testBranchTailOrReturnFn(consts, tail)

	pc := RootExecContext.PrepareCall(fn, 1)
	if pc == nil {
		t.Fatal("PrepareCall rejected a plain fixed-arity Func")
	}
	defer pc.Release()

	v, err := pc.Call1(NIL) // falsy -> tail-calls into testConstFn
	if err != nil {
		t.Fatal(err)
	}
	if v != Int(42) {
		t.Fatalf("tail-call path returned %v", v)
	}
	v, err = pc.Call1(Int(7)) // truthy -> must run fn, not the rebound tail target
	if err != nil {
		t.Fatal(err)
	}
	if v != Int(7) {
		t.Fatalf("prepared frame stayed rebound to the tail target: %v", v)
	}
}

// An error unwind must not poison the owned frame for subsequent calls.
func TestPreparedCallReusableAfterError(t *testing.T) {
	consts := NewConsts()
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_ARG)
	chunk.Append32(0)
	chunk.Append(OP_BRANCH_TRUE)
	chunk.Append32(7) // truthy -> throw the arg
	chunk.Append(OP_LOAD_ARG)
	chunk.Append32(0)
	chunk.Append(OP_RETURN)
	chunk.Append(OP_NOOP)
	chunk.Append(OP_NOOP)
	chunk.Append(OP_LOAD_ARG)
	chunk.Append32(0)
	chunk.Append(OP_THROW)
	chunk.SetMaxStack(4)
	fn := MakeFunc(1, false, chunk)

	pc := RootExecContext.PrepareCall(fn, 1)
	if pc == nil {
		t.Fatal("PrepareCall rejected a plain fixed-arity Func")
	}
	defer pc.Release()

	if _, err := pc.Call1(Int(1)); err == nil {
		t.Fatal("throwing call did not surface an error")
	}
	v, err := pc.Call1(FALSE)
	if err != nil {
		t.Fatal(err)
	}
	if v != FALSE {
		t.Fatalf("call after error returned %v", v)
	}
}

func TestPreparedCallRejectsNonBytecodeAndVariadicTargets(t *testing.T) {
	native, err := NativeFnType.Wrap(func(_ []Value) (Value, error) { return NIL, nil })
	if err != nil {
		t.Fatal(err)
	}
	if pc := RootExecContext.PrepareCall(native.(Fn), 1); pc != nil {
		t.Fatal("PrepareCall accepted a NativeFn")
	}

	consts := NewConsts()
	variadic := MakeFunc(1, true, NewCodeChunk(consts))
	if pc := RootExecContext.PrepareCall(variadic, 1); pc != nil {
		t.Fatal("PrepareCall accepted a variadic Func")
	}
}

// Arities without a matching CallN method must not prepare: CallN could not
// populate their argument slots (arity 0 would panic, higher arities would
// pass Go nil interfaces into bytecode).
func TestPreparedCallRejectsUnsupportedArities(t *testing.T) {
	consts := NewConsts()
	nullary := testConstFn(consts, Int(1))
	if pc := RootExecContext.PrepareCall(nullary, 0); pc != nil {
		t.Fatal("PrepareCall accepted arity 0 with no Call0 entry point")
	}
	ternary := testBytecodeFnReturningArg(consts, 3)
	if pc := RootExecContext.PrepareCall(ternary, 3); pc != nil {
		t.Fatal("PrepareCall accepted arity 3 with no Call3 entry point")
	}
}

// A CallN whose N disagrees with the prepared arity must error, not panic or
// invoke with unpopulated argument slots.
func TestPreparedCallArityMismatchErrors(t *testing.T) {
	consts := NewConsts()

	unary := RootExecContext.PrepareCall(testBytecodeFnReturningArg(consts, 1), 1)
	if unary == nil {
		t.Fatal("PrepareCall rejected a plain unary Func")
	}
	defer unary.Release()
	if _, err := unary.Call2(Int(1), Int(2)); err == nil {
		t.Fatal("Call2 on a unary preparation did not error")
	}
	if v, err := unary.Call1(Int(7)); err != nil || v != Int(7) {
		t.Fatalf("unary preparation unusable after mismatch: %v, %v", v, err)
	}

	binary := RootExecContext.PrepareCall(testBytecodeFnReturningArg(consts, 2), 2)
	if binary == nil {
		t.Fatal("PrepareCall rejected a plain binary Func")
	}
	defer binary.Release()
	if _, err := binary.Call1(Int(1)); err == nil {
		t.Fatal("Call1 on a binary preparation did not error")
	}
}

func TestPreparedCallCall2RepeatedInvocation(t *testing.T) {
	consts := NewConsts()
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_ARG)
	chunk.Append32(1)
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)
	second := MakeFunc(2, false, chunk)

	pc := RootExecContext.PrepareCall(second, 2)
	if pc == nil {
		t.Fatal("PrepareCall rejected a plain binary Func")
	}
	defer pc.Release()
	for i := 0; i < 3; i++ {
		v, err := pc.Call2(Int(-1), Int(i))
		if err != nil {
			t.Fatal(err)
		}
		if v != Int(i) {
			t.Fatalf("call %d returned %v", i, v)
		}
	}
}
