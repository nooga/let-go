package vm

import (
	"strings"
	"testing"
)

func testBytecodeFnReturningArg(consts *Consts, arity int) *Func {
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_ARG)
	chunk.Append32(0)
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)
	return MakeFunc(arity, false, chunk)
}

func testInvokeZeroChunk(consts *Consts, fn Fn) *CodeChunk {
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_CONST)
	chunk.Append32(consts.Intern(fn))
	chunk.Append(OP_INVOKE)
	chunk.Append32(0)
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)
	return chunk
}

func testTailCallZeroChunk(consts *Consts, fn Fn) *CodeChunk {
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_CONST)
	chunk.Append32(consts.Intern(fn))
	chunk.Append(OP_TAIL_CALL)
	chunk.Append32(0)
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)
	return chunk
}

func TestResolveBytecodeCallSelectsMultiArityAndPreservesCaptures(t *testing.T) {
	consts := NewConsts()
	oneArg := testBytecodeFnReturningArg(consts, 1)
	twoArg := testBytecodeFnReturningArg(consts, 2)
	multi, err := MakeMultiArity([]Value{oneArg, twoArg})
	if err != nil {
		t.Fatal(err)
	}

	target, direct, err := resolveBytecodeCall(NewMetaFn(multi, NIL), []Value{Int(7)})
	if err != nil {
		t.Fatal(err)
	}
	if !direct || target.fn.chunk != oneArg.chunk {
		t.Fatal("metadata-wrapped multi-arity call did not select its bytecode variant")
	}

	captures := []Value{Int(99)}
	closure := &Closure{fn: multi, closedOvers: captures}
	target, direct, err = resolveBytecodeCall(closure, []Value{Int(7), Int(8)})
	if err != nil {
		t.Fatal(err)
	}
	if !direct || target.fn.chunk != twoArg.chunk {
		t.Fatal("closure over multi-arity function did not select its bytecode variant")
	}
	if len(target.closedOvers) != 1 || target.closedOvers[0] != captures[0] {
		t.Fatalf("closure captures were not preserved: got %v, want %v", target.closedOvers, captures)
	}
}

func TestResolveBytecodeCallVariadicPackingDoesNotMutateBorrowedArgs(t *testing.T) {
	consts := NewConsts()
	variadic := MakeFunc(2, true, testBytecodeFnReturningArg(consts, 2).chunk)
	borrowed := []Value{Int(7), Int(11), Int(13)}

	target, direct, err := resolveBytecodeCall(variadic, borrowed)
	if err != nil {
		t.Fatal(err)
	}
	if !direct {
		t.Fatal("variadic bytecode function was not resolved")
	}
	if borrowed[0] != Value(Int(7)) || borrowed[1] != Value(Int(11)) || borrowed[2] != Value(Int(13)) {
		t.Fatalf("resolver mutated borrowed args: got %v", borrowed)
	}
	if len(target.args) != 2 || target.args[0] != Value(Int(7)) {
		t.Fatalf("unexpected packed args: %v", target.args)
	}
	if target.args[1].String() != "(11 13)" {
		t.Fatalf("unexpected packed rest args: %v", target.args[1])
	}
}

func TestInstallBytecodeCallOwnsBorrowedArgs(t *testing.T) {
	consts := NewConsts()
	targetFn := testBytecodeFnReturningArg(consts, 2)
	frame := NewFrame(testBytecodeFnReturningArg(consts, 0).chunk, nil)
	frame.stack[0] = Int(7)
	frame.stack[1] = Int(11)
	frame.sp = 2

	target, direct, err := resolveBytecodeCall(targetFn, frame.stack[:2])
	if err != nil {
		t.Fatal(err)
	}
	if !direct {
		t.Fatal("fixed-arity bytecode function was not resolved")
	}
	installBytecodeCall(frame, target)
	frame.stack[0] = Int(99)
	if frame.args[0] != Value(Int(7)) || frame.args[1] != Value(Int(11)) {
		t.Fatalf("tail transition retained operand-stack-backed args: %v", frame.args)
	}
	ReleaseFrame(frame)
}

func TestSuspendedCallArityValidatesContinuation(t *testing.T) {
	valid := NewCodeChunk(NewConsts())
	valid.Append(OP_INVOKE)
	valid.Append32(1)
	valid.SetMaxStack(4)
	frame := NewFrame(valid, nil)
	frame.sp = 2
	if arity, err := suspendedCallArity(frame); err != nil || arity != 1 {
		t.Fatalf("valid continuation: arity=%d err=%v", arity, err)
	}

	tests := []struct {
		name  string
		code  []int32
		ip    int
		sp    int
		match string
	}{
		{name: "negative ip", code: []int32{OP_INVOKE, 0}, ip: -1, sp: 1, match: "out of bounds"},
		{name: "ip past end", code: []int32{OP_INVOKE, 0}, ip: 2, sp: 1, match: "out of bounds"},
		{name: "missing arity", code: []int32{OP_INVOKE}, ip: 0, sp: 1, match: "missing arity"},
		{name: "wrong opcode", code: []int32{OP_NOOP, 0}, ip: 0, sp: 1, match: "NOOP"},
		{name: "invalid stack depth", code: []int32{OP_INVOKE, 0}, ip: 0, sp: 5, match: "stack depth 5 is out of bounds"},
		{name: "negative arity", code: []int32{OP_INVOKE, -1}, ip: 0, sp: 1, match: "arity -1"},
		{name: "arity exceeds stack", code: []int32{OP_TAIL_CALL, 2}, ip: 0, sp: 2, match: "exceeds stack depth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk := NewCodeChunk(NewConsts())
			chunk.Append(tt.code...)
			chunk.SetMaxStack(4)
			candidate := NewFrame(chunk, nil)
			candidate.ip = tt.ip
			candidate.sp = tt.sp
			_, err := suspendedCallArity(candidate)
			if err == nil || !strings.Contains(err.Error(), tt.match) {
				t.Fatalf("got %v, want error containing %q", err, tt.match)
			}
			ReleaseFrame(candidate)
		})
	}
	ReleaseFrame(frame)
}

func TestTailCallClosureArityErrorMatchesInvokeAndIsCatchable(t *testing.T) {
	consts := NewConsts()
	closure := &Closure{fn: testBytecodeFnReturningArg(consts, 1)}

	run := func(chunk *CodeChunk) error {
		frame := NewFrame(chunk, nil)
		frame.ec = RootExecContext
		_, err := frame.Run()
		ReleaseFrame(frame)
		return err
	}
	invokeErr := run(testInvokeZeroChunk(consts, closure))
	tailErr := run(testTailCallZeroChunk(consts, closure))
	if invokeErr == nil || tailErr == nil {
		t.Fatalf("expected arity errors: invoke=%v tail=%v", invokeErr, tailErr)
	}
	if got, want := innermostMessage(tailErr), innermostMessage(invokeErr); got != want {
		t.Fatalf("tail-call arity error differs from invoke: got %q, want %q", got, want)
	}

	catching := NewCodeChunk(consts)
	catching.Append(OP_TRY_PUSH)
	catching.Append32(9) // catch at the final OP_RETURN
	catching.Append32(0)
	catching.Append(OP_LOAD_CONST)
	catching.Append32(consts.Intern(closure))
	catching.Append(OP_TAIL_CALL)
	catching.Append32(0)
	catching.Append(OP_TRY_POP)
	catching.Append(OP_RETURN)
	catching.Append(OP_RETURN)
	catching.SetMaxStack(4)
	frame := NewFrame(catching, nil)
	frame.ec = RootExecContext
	result, err := frame.Run()
	ReleaseFrame(frame)
	if err != nil {
		t.Fatalf("caught tail-call arity error escaped: %v", err)
	}
	ex, ok := result.(*ExInfo)
	if !ok || !strings.Contains(ex.Message(), "expected 1 args, got 0") {
		t.Fatalf("unexpected caught value: %v", result)
	}
}

func TestDescendedFrameLifecycle(t *testing.T) {
	oldAllocAttr := allocAttrEnabled
	allocAttrEnabled = true
	attrMu.Lock()
	oldAttrStack := attrStack
	attrStack = nil
	attrMu.Unlock()
	t.Cleanup(func() {
		attrMu.Lock()
		attrStack = oldAttrStack
		attrMu.Unlock()
		allocAttrEnabled = oldAllocAttr
	})

	consts := NewConsts()
	var parentChunk, childChunk *CodeChunk
	checkerValue, err := NativeFnType.Wrap(func(_ []Value) (Value, error) {
		attrMu.Lock()
		defer attrMu.Unlock()
		if len(attrStack) != 2 {
			return NIL, NewExecutionError("allocation attribution stack did not include parent and child")
		}
		if attrStack[0].code != parentChunk || attrStack[1].code != childChunk {
			return NIL, NewExecutionError("allocation attribution stack has the wrong frame order")
		}
		return Int(1), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	checker := checkerValue.(Fn)
	childChunk = testInvokeZeroChunk(consts, checker)
	child := MakeFunc(0, false, childChunk)
	parentChunk = testInvokeZeroChunk(consts, child)

	root := NewFrame(parentChunk, nil)
	root.ec = RootExecContext
	if _, err := root.Run(); err != nil {
		t.Fatal(err)
	}
	ReleaseFrame(root)

	attrMu.Lock()
	defer attrMu.Unlock()
	if len(attrStack) != 0 {
		t.Fatalf("allocation attribution stack leaked %d frame(s)", len(attrStack))
	}
}

func TestZeroArityClosureTailCallDescendsInDispatchLoop(t *testing.T) {
	oldAllocAttr := allocAttrEnabled
	allocAttrEnabled = true
	attrMu.Lock()
	oldAttrStack := attrStack
	attrStack = nil
	attrMu.Unlock()
	t.Cleanup(func() {
		attrMu.Lock()
		attrStack = oldAttrStack
		attrMu.Unlock()
		allocAttrEnabled = oldAllocAttr
	})

	consts := NewConsts()
	checkerValue, err := NativeFnType.Wrap(func(_ []Value) (Value, error) {
		attrMu.Lock()
		defer attrMu.Unlock()
		if len(attrStack) != 2 {
			return NIL, NewExecutionError("zero-arity tail call did not enter exactly one child frame")
		}
		if attrStack[1].parent != attrStack[0] {
			return NIL, NewExecutionError("zero-arity closure re-entered Frame.Run instead of linking its parent")
		}
		return Int(1), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	childChunk := testInvokeZeroChunk(consts, checkerValue.(Fn))
	closure := &Closure{fn: MakeFunc(0, false, childChunk)}
	root := NewFrame(testTailCallZeroChunk(consts, closure), nil)
	root.ec = RootExecContext
	if _, err := root.Run(); err != nil {
		t.Fatal(err)
	}
	ReleaseFrame(root)

	attrMu.Lock()
	defer attrMu.Unlock()
	if len(attrStack) != 0 {
		t.Fatalf("allocation attribution stack leaked %d frame(s)", len(attrStack))
	}
}

func TestOpcodeProfileSeparatesDescendedFrames(t *testing.T) {
	ResetProfile()
	ProfilingEnabled.Store(true)
	t.Cleanup(func() {
		ProfilingEnabled.Store(false)
		ResetProfile()
	})

	consts := NewConsts()
	childChunk := NewCodeChunk(consts)
	childChunk.Append(OP_LOAD_CONST)
	childChunk.Append32(consts.Intern(NIL))
	childChunk.Append(OP_RETURN)
	childChunk.SetMaxStack(4)
	parentChunk := testInvokeZeroChunk(consts, MakeFunc(0, false, childChunk))

	root := NewFrame(parentChunk, nil)
	root.ec = RootExecContext
	if _, err := root.Run(); err != nil {
		t.Fatal(err)
	}
	ReleaseFrame(root)

	var initialLoads, crossFrameLoads, invokeReturns uint64
	for _, pair := range PairSnapshot() {
		switch {
		case pair.Prev == 0 && pair.Curr == uint8(OP_LOAD_CONST):
			initialLoads = pair.Count
		case pair.Prev == uint8(OP_INVOKE) && pair.Curr == uint8(OP_LOAD_CONST):
			crossFrameLoads = pair.Count
		case pair.Prev == uint8(OP_INVOKE) && pair.Curr == uint8(OP_RETURN):
			invokeReturns = pair.Count
		}
	}
	if initialLoads != 2 {
		t.Fatalf("first opcode was not reset for both frames: got %d initial LOAD_CONST pairs", initialLoads)
	}
	if crossFrameLoads != 0 {
		t.Fatalf("profiler joined parent OP_INVOKE to child LOAD_CONST %d time(s)", crossFrameLoads)
	}
	if invokeReturns != 1 {
		t.Fatalf("parent profiler state was not restored: got %d INVOKE→RETURN pairs", invokeReturns)
	}
}

func TestDescendedErrorReleasesChildFrame(t *testing.T) {
	framePoolMu.Lock()
	oldPool := framePoolStack
	framePoolStack = nil
	framePoolMu.Unlock()
	t.Cleanup(func() {
		framePoolMu.Lock()
		clear(framePoolStack)
		framePoolStack = oldPool
		framePoolMu.Unlock()
	})

	consts := NewConsts()
	badChunk := NewCodeChunk(consts)
	badChunk.Append(OP_LOAD_ARG)
	badChunk.Append32(0)
	badChunk.Append(OP_RETURN)
	badChunk.SetMaxStack(4)
	parentChunk := testInvokeZeroChunk(consts, MakeFunc(0, false, badChunk))

	root := NewFrame(parentChunk, nil)
	root.ec = RootExecContext
	if _, err := root.Run(); err == nil {
		t.Fatal("expected descended child to fail")
	}
	if root.parent != nil {
		t.Fatal("root retained a parent link after failed descent")
	}

	framePoolMu.Lock()
	pooled := append([]*Frame(nil), framePoolStack...)
	framePoolMu.Unlock()
	if len(pooled) != 1 {
		t.Fatalf("failed child was not returned exactly once: pool has %d frames", len(pooled))
	}
	if pooled[0].parent != nil || pooled[0].code != nil || pooled[0].args != nil {
		t.Fatal("failed child retained execution references after pooling")
	}
	ReleaseFrame(root)
}

// rawPanicFn panics with a plain Go value when invoked. It deliberately does
// NOT wrap itself in RecoverPanic the way NativeFn does: ec.Invoke's default
// arm calls it bare, so the panic unwinds through the live dispatch loop —
// the path releasePanickedFrames (attribution on) and the documented
// leak-to-GC trade (attribution off) exist for.
type rawPanicFn struct{ payload string }

func (p *rawPanicFn) String() string                { return "<raw-panic-fn>" }
func (p *rawPanicFn) Type() ValueType               { return NativeFnType }
func (p *rawPanicFn) Unbox() any                    { return p }
func (p *rawPanicFn) Arity() int                    { return 0 }
func (p *rawPanicFn) Invoke([]Value) (Value, error) { panic(p.payload) }

// A Go panic escaping a native while the chain is mid-descent must be
// preserved as-is, and the pool must see the currently-intended behavior in
// BOTH configurations: with attribution on, runChainProtected releases the
// suspended chain's pooled frames; with attribution off (the default), Run is
// deliberately defer-free and the frames leak to GC — pre-#645 parity, per
// Run's doc comment. A future change that reintroduces a must-release-on-panic
// invariant on the default path has to update this test consciously.
func TestPanicUnwindingDescendedFramePinsPoolBehavior(t *testing.T) {
	framePoolMu.Lock()
	oldPool := framePoolStack
	framePoolStack = nil
	framePoolMu.Unlock()
	t.Cleanup(func() {
		framePoolMu.Lock()
		clear(framePoolStack)
		framePoolStack = oldPool
		framePoolMu.Unlock()
	})

	// child (bytecode, arity 1): invoke its argument — the raw-panicking fn.
	consts := NewConsts()
	childChunk := NewCodeChunk(consts)
	childChunk.Append(OP_LOAD_ARG)
	childChunk.Append32(0)
	childChunk.Append(OP_INVOKE)
	childChunk.Append32(0)
	childChunk.Append(OP_RETURN)
	childChunk.SetMaxStack(4)
	child := MakeFunc(1, false, childChunk)

	// root: descend into child with the panicking fn as the argument, so the
	// panic fires while root is a suspended parent of a live pooled frame.
	rootChunk := NewCodeChunk(consts)
	rootChunk.Append(OP_LOAD_CONST)
	rootChunk.Append32(consts.Intern(child))
	rootChunk.Append(OP_LOAD_ARG)
	rootChunk.Append32(0)
	rootChunk.Append(OP_INVOKE)
	rootChunk.Append32(1)
	rootChunk.Append(OP_RETURN)
	rootChunk.SetMaxStack(4)

	root := NewFrame(rootChunk, []Value{&rawPanicFn{payload: "raw-panic-719"}})
	root.ec = RootExecContext

	recovered := func() (r any) {
		defer func() { r = recover() }()
		_, _ = root.Run()
		return nil
	}()
	if recovered != "raw-panic-719" {
		t.Fatalf("panic was not preserved through the dispatch loop: got %v", recovered)
	}

	framePoolMu.Lock()
	pooled := append([]*Frame(nil), framePoolStack...)
	framePoolMu.Unlock()
	if allocAttrEnabled {
		if len(pooled) != 1 {
			t.Fatalf("attribution on: panicked child not released exactly once: pool has %d frames", len(pooled))
		}
	} else {
		if len(pooled) != 0 {
			t.Fatalf("default config: panic path released %d frame(s); the intended behavior is leak-to-GC (update this test if that changed on purpose)", len(pooled))
		}
	}
	ReleaseFrame(root)
}
