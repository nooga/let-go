package vm

import (
	"sync"
	"testing"
)

func TestFrameArena_NestedSelfInvoke(t *testing.T) {
	// Countdown via OP_INVOKE self-call exercises the arena-threaded path.
	consts := NewConsts()
	body := NewCodeChunk(consts)
	zero := consts.Intern(Int(0))
	one := consts.Intern(Int(1))

	fn := MakeFunc(1, false, body)
	fn.SetName("arena-countdown")
	fnIdx := consts.Intern(fn)

	body.Append(OP_LOAD_ARG)
	body.Append32(0)
	body.Append(OP_LOAD_CONST)
	body.Append32(zero)
	body.Append(OP_EQ)
	br := body.Length()
	body.Append(OP_BRANCH_TRUE)
	body.Append32(0)
	body.Append(OP_LOAD_CONST)
	body.Append32(fnIdx)
	body.Append(OP_LOAD_ARG)
	body.Append32(0)
	body.Append(OP_LOAD_CONST)
	body.Append32(one)
	body.Append(OP_SUB)
	body.Append(OP_INVOKE)
	body.Append32(1)
	body.Append(OP_RETURN)
	ret0 := body.Length()
	body.Update32(br+1, int32(ret0-br))
	body.Append(OP_LOAD_CONST)
	body.Append32(zero)
	body.Append(OP_RETURN)
	body.SetMaxStack(8)

	out, err := fn.Invoke([]Value{Int(5)})
	if err != nil {
		t.Fatal(err)
	}
	if got := out.Unbox().(int); got != 0 {
		t.Fatalf("got %v, want 0", got)
	}
}

func TestFrameArena_ConcurrentTopLevel(t *testing.T) {
	consts := NewConsts()
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_ARG)
	chunk.Append32(0)
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)
	fn := MakeFunc(1, false, chunk)

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				out, err := fn.Invoke([]Value{Int(n)})
				if err != nil {
					t.Errorf("invoke: %v", err)
					return
				}
				if out.Unbox().(int) != n {
					t.Errorf("got %v want %d", out, n)
					return
				}
			}
		}(i)
	}
	wg.Wait()
}

func TestFrameArena_TopLevelInvokeDoesNotAllocate(t *testing.T) {
	consts := NewConsts()
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_ARG)
	chunk.Append32(0)
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)
	fn := MakeFunc(1, false, chunk)
	args := []Value{Int(42)}

	var invokeErr error
	allocs := testing.AllocsPerRun(1000, func() {
		_, invokeErr = fn.Invoke(args)
	})
	if invokeErr != nil {
		t.Fatal(invokeErr)
	}
	if allocs != 0 {
		t.Fatalf("top-level invoke allocated %.2f objects per call, want 0", allocs)
	}
}

func TestFrameArena_ReusesFrames(t *testing.T) {
	arena := &frameArena{}
	consts := NewConsts()
	chunk := NewCodeChunk(consts)
	chunk.Append(OP_LOAD_CONST)
	chunk.Append32(consts.Intern(NIL))
	chunk.Append(OP_RETURN)
	chunk.SetMaxStack(4)

	f1 := newFrameIn(arena, chunk, nil)
	ReleaseFrame(f1)
	if len(arena.free) != 1 {
		t.Fatalf("expected 1 free frame, got %d", len(arena.free))
	}
	f2 := newFrameIn(arena, chunk, nil)
	if len(arena.free) != 0 {
		t.Fatalf("expected freelist empty after get, got %d", len(arena.free))
	}
	if f1 != f2 {
		t.Fatal("arena should hand back the same Frame pointer")
	}
	ReleaseFrame(f2)
}
