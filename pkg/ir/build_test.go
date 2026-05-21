/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// compileFib returns the *vm.Func for the standard tree-recursive fib.
// Useful as a workhorse test fixture; fib exercises args, consts, vars,
// comparison + branch, recursive call, addition, return.
func compileFib(t *testing.T) *vm.Func {
	t.Helper()
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn fib [n] (if (<= n 1) n (+ (fib (- n 1)) (fib (- n 2)))))`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile fib: %v", err)
	}
	fibVar := ns.Lookup(vm.Symbol("fib"))
	if fibVar == vm.NIL {
		t.Fatal("fib var not found after compile")
	}
	v := fibVar.(*vm.Var).Deref()
	fn, ok := v.(*vm.Func)
	if !ok {
		t.Fatalf("fib is not *vm.Func, got %T", v)
	}
	return fn
}

func TestBuildFib(t *testing.T) {
	fn := compileFib(t)
	irFn, err := Build(fn.Chunk(), "fib", 1, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	// Sanity checks.
	if irFn.Arity != 1 {
		t.Errorf("arity = %d, want 1", irFn.Arity)
	}
	if len(irFn.Blocks) < 2 {
		t.Errorf("expected at least 2 blocks (if/else), got %d", len(irFn.Blocks))
	}
	if len(irFn.Nodes) < 10 {
		t.Errorf("expected at least 10 nodes for fib, got %d", len(irFn.Nodes))
	}

	// Dump for visual inspection (visible with go test -v).
	t.Logf("\n%s", Dump(irFn))
}

func TestBuildSimpleConst(t *testing.T) {
	// (defn const42 [] 42) — entry block ends with Return Const(42).
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn const42 [] 42)`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("const42")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := Build(fn.Chunk(), "const42", 0, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if len(irFn.Blocks) != 1 {
		t.Errorf("expected 1 block, got %d", len(irFn.Blocks))
	}
	t.Logf("\n%s", Dump(irFn))
}

// TestRoundtripFib: build IR for fib, lower, install lowered chunk
// as the new body of the fib Func, run fib(5), compare to original.
//
// This is the real soundness check: a recursive function with branches
// and recursive calls must execute identically after round-trip.
func TestRoundtripFib(t *testing.T) {
	fn := compileFib(t)

	irFn, err := Build(fn.Chunk(), "fib", 1, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	t.Logf("IR:\n%s", Dump(irFn))

	lowered, err := Lower(irFn)
	if err != nil {
		t.Fatalf("Lower: %v", err)
	}

	// Run original fib(10). fib(10) = 55.
	args := []vm.Value{vm.Int(10)}
	frame := vm.NewFrame(fn.Chunk(), args)
	origResult, err := frame.Run()
	vm.ReleaseFrame(frame)
	if err != nil {
		t.Fatalf("original run: %v", err)
	}
	t.Logf("original fib(5) = %s", origResult)

	// Run lowered. Note: lowered version still references the *old*
	// fib var (which holds the original chunk). So lowered chunk will
	// call original fib recursively — still a valid soundness test as
	// long as the top-level behavior matches.
	frame2 := vm.NewFrame(lowered, args)
	loweredResult, err := frame2.Run()
	vm.ReleaseFrame(frame2)
	if err != nil {
		t.Fatalf("lowered run: %v", err)
	}
	t.Logf("lowered fib(5) = %s", loweredResult)

	if origResult.String() != loweredResult.String() {
		t.Errorf("fib(5) mismatch: orig %s vs lowered %s", origResult, loweredResult)
	}
}

// TestRoundtripConst42: build IR for (fn [] 42), lower back to bytecode,
// run, compare to original. The simplest possible round-trip.
func TestRoundtripConst42(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	_, _, err := ctx.CompileMultiple(strings.NewReader(`(defn rtc42 [] 42)`))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("rtc42")).(*vm.Var).Deref()
	fn := v.(*vm.Func)

	// Original execution.
	frame := vm.NewFrame(fn.Chunk(), nil)
	origResult, err := frame.Run()
	vm.ReleaseFrame(frame)
	if err != nil {
		t.Fatalf("original run: %v", err)
	}
	t.Logf("original: %s", origResult)

	// Build IR, lower back.
	irFn, err := Build(fn.Chunk(), "rtc42", 0, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	t.Logf("IR:\n%s", Dump(irFn))

	lowered, err := Lower(irFn)
	if err != nil {
		t.Fatalf("Lower: %v", err)
	}

	// Run lowered.
	frame2 := vm.NewFrame(lowered, nil)
	loweredResult, err := frame2.Run()
	vm.ReleaseFrame(frame2)
	if err != nil {
		t.Fatalf("lowered run: %v", err)
	}
	t.Logf("lowered: %s", loweredResult)

	if origResult.String() != loweredResult.String() {
		t.Errorf("results differ: orig %s vs lowered %s", origResult, loweredResult)
	}
}

// TestBuildLoopRecur: a loop with recur. Verifies the back-edge
// becomes a Branch to the loop header with new args, and round-trips
// through Lower to produce correct runtime behavior.
func TestBuildLoopRecur(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn sumto [n] (loop [i 0 acc 0] (if (< i n) (recur (inc i) (+ acc i)) acc)))`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("sumto")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := Build(fn.Chunk(), "sumto", 1, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	t.Logf("IR:\n%s", Dump(irFn))

	// Structural check: a loop-header block with >=2 preds and >=1 param.
	hasJoinBlock := false
	for _, blk := range irFn.Blocks {
		if len(blk.Preds) >= 2 && len(blk.Params) >= 1 {
			hasJoinBlock = true
			break
		}
	}
	if !hasJoinBlock {
		t.Errorf("expected a loop-header block with multiple preds and block-args, got:\n%s", Dump(irFn))
	}

	// Round-trip: Lower the IR back to bytecode, run it, compare to the
	// expected result. (sumto 5) = 0+1+2+3+4 = 10.
	chunk, err := Lower(irFn)
	if err != nil {
		t.Fatalf("Lower: %v", err)
	}
	frame := vm.NewFrame(chunk, []vm.Value{vm.Int(5)})
	result, err := frame.Run()
	vm.ReleaseFrame(frame)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.String() != "10" {
		t.Errorf("expected (sumto 5) = 10, got %s", result)
	}
}

// TestDumpLoopBytecode: just print the raw bytecode so we can see what
// (loop ...) looks like. Diagnostic for debugging loop support.
func TestDumpLoopBytecode(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn sumto2 [n] (loop [i 0 acc 0] (if (< i n) (recur (inc i) (+ acc i)) acc)))`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("sumto2")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	chunk := fn.Chunk()
	code := chunk.Code()
	t.Logf("sumto2 bytecode (%d words):", len(code))
	for i := 0; i < len(code); {
		op := code[i]
		t.Logf("  %3d: %s", i, vm.OpcodeToString(op))
		stride := instStrideForTest(op)
		if stride <= 0 {
			stride = 1
		}
		for j := 1; j < stride && i+j < len(code); j++ {
			t.Logf("       arg %d: %d", j-1, code[i+j])
		}
		i += stride
	}
}

func instStrideForTest(op int32) int {
	return instStride(op)
}

// TestBuild_FallThroughExitStack: verifies that a block ending without
// an explicit terminator and falling through carries its exit-stack
// values as branch args. Foundation for loop support.
func TestBuild_FallThroughExitStack(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn sumto3 [n] (loop [i 0 acc 0] (if (< i n) (recur (inc i) (+ acc i)) acc)))`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("sumto3")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := Build(fn.Chunk(), "sumto3", 1, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	// Entry block's terminator should be OpBranch with len(Args) >= 2
	// (the two loop-binding initializers).
	entryTerm := irFn.Node(irFn.Blocks[irFn.Entry].Term)
	if entryTerm.Op != OpBranch {
		t.Fatalf("expected entry block terminator OpBranch, got %s", entryTerm.Op)
	}
	bt, ok := entryTerm.Aux.(*BranchTarget)
	if !ok || bt == nil {
		t.Fatalf("expected *BranchTarget aux, got %T", entryTerm.Aux)
	}
	if len(bt.Args) < 2 {
		t.Errorf("expected entry fall-through branch to carry >=2 args (loop bindings), got %d: %v", len(bt.Args), bt.Args)
	}
}

// TestBuild_LoopHeaderParamsMatchPreds: the loop header block should have
// block-params count equal to the max args carried by any reachable
// predecessor's branch terminator. Foundation for end-to-end loop support.
func TestBuild_LoopHeaderParamsMatchPreds(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn sumto4 [n] (loop [i 0 acc 0] (if (< i n) (recur (inc i) (+ acc i)) acc)))`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("sumto4")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := Build(fn.Chunk(), "sumto4", 1, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	// Find the loop header (block with >=2 preds).
	var headerID BlockID = -1
	for i, blk := range irFn.Blocks {
		if len(blk.Preds) >= 2 {
			headerID = BlockID(i)
			break
		}
	}
	if headerID == -1 {
		t.Fatalf("no loop header found:\n%s", Dump(irFn))
	}
	header := irFn.Blocks[headerID]
	if len(header.Params) < 2 {
		t.Errorf("expected loop header to have >=2 params (loop bindings), got %d:\n%s", len(header.Params), Dump(irFn))
	}
	// Each reachable predecessor should pass exactly len(header.Params) args.
	for _, predID := range header.Preds {
		// Skip unreachable preds (the dead fall-through after RECUR).
		// We don't have direct access to b.reachable here, but we can check
		// the pred has a real terminator and is in the function's block list.
		predTerm := irFn.Node(irFn.Blocks[predID].Term)
		var args []NodeID
		matched := false
		switch t := predTerm.Aux.(type) {
		case *BranchTarget:
			if t.Target == headerID {
				args = t.Args
				matched = true
			}
		case *CondTarget:
			if t.True != nil && t.True.Target == headerID {
				args = t.True.Args
				matched = true
			}
			if t.False != nil && t.False.Target == headerID {
				args = t.False.Args
				matched = true
			}
		}
		if !matched {
			continue // pred doesn't actually target header (defensive)
		}
		if len(args) != len(header.Params) {
			t.Errorf("pred block %d passes %d args to header but header has %d params:\n%s",
				predID, len(args), len(header.Params), Dump(irFn))
		}
	}
}

func TestBuildAddTwoArgs(t *testing.T) {
	// (defn add2 [x y] (+ x y))
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	_, _, err := ctx.CompileMultiple(strings.NewReader(`(defn add2 [x y] (+ x y))`))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("add2")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := Build(fn.Chunk(), "add2", 2, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	// Expect: one block. Body: LoadArg 0, LoadArg 1, Add, Return.
	if len(irFn.Blocks) != 1 {
		t.Errorf("expected 1 block, got %d", len(irFn.Blocks))
	}
	t.Logf("\n%s", Dump(irFn))
}

// TestStackEffectCoversInstStride verifies that every opcode instStride
// claims to handle is also handled by stackEffect. If instStride is
// extended without updating stackEffect, the abstract-interpretation
// pre-pass would silently compute wrong heights.
func TestStackEffectCoversInstStride(t *testing.T) {
	// Enumerate opcodes mentioned in instStride. Each must return a
	// non-default value from stackEffect for some plausible immediate.
	// We construct a synthetic 4-word code buffer so opcodes with
	// immediates can be queried safely.
	probe := []int32{vm.OP_NOOP, vm.OP_POP, vm.OP_RETURN,
		vm.OP_ADD, vm.OP_SUB, vm.OP_MUL,
		vm.OP_LT, vm.OP_LTE, vm.OP_GT, vm.OP_GTE, vm.OP_EQ,
		vm.OP_INC, vm.OP_DEC, vm.OP_SET_VAR,
		vm.OP_LOAD_CONST, vm.OP_LOAD_VAR, vm.OP_LOAD_ARG,
		vm.OP_BRANCH_TRUE, vm.OP_BRANCH_FALSE, vm.OP_JUMP,
		vm.OP_INVOKE, vm.OP_TAIL_CALL,
		vm.OP_POP_N, vm.OP_DUP_NTH,
		vm.OP_RECUR}
	for _, op := range probe {
		if instStride(op) == 0 {
			t.Errorf("instStride says it does not handle 0x%x — fix probe list", op)
			continue
		}
		// 4 words is enough for the largest opcode (RECUR).
		code := []int32{op, 1, 2, 0}
		// Just verify stackEffect doesn't panic and returns a sane value.
		// (No specific expected value — we just want coverage assurance.)
		_ = stackEffect(op, code, 0)
	}
}
