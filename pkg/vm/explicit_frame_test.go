package vm

import "testing"

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

func TestChildFrameForSelectsMultiArityAndPreservesCaptures(t *testing.T) {
	consts := NewConsts()
	oneArg := testBytecodeFnReturningArg(consts, 1)
	twoArg := testBytecodeFnReturningArg(consts, 2)
	multi, err := MakeMultiArity([]Value{oneArg, twoArg})
	if err != nil {
		t.Fatal(err)
	}

	child, err := childFrameFor(NewMetaFn(multi, NIL), []Value{Int(7)}, RootExecContext)
	if err != nil {
		t.Fatal(err)
	}
	if child == nil || child.code != oneArg.chunk {
		t.Fatal("metadata-wrapped multi-arity call did not select its bytecode variant")
	}
	ReleaseFrame(child)

	captures := []Value{Int(99)}
	closure := &Closure{fn: multi, closedOvers: captures}
	child, err = childFrameFor(closure, []Value{Int(7), Int(8)}, RootExecContext)
	if err != nil {
		t.Fatal(err)
	}
	if child == nil || child.code != twoArg.chunk {
		t.Fatal("closure over multi-arity function did not select its bytecode variant")
	}
	if len(child.closedOvers) != 1 || child.closedOvers[0] != captures[0] {
		t.Fatalf("closure captures were not preserved: got %v, want %v", child.closedOvers, captures)
	}
	ReleaseFrame(child)
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
