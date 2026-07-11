/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"runtime"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

func asyncFn(t *testing.T, name string) vm.Fn {
	t.Helper()
	v := NS("async").Lookup(vm.Symbol(name))
	if v == nil {
		t.Fatalf("async/%s not found", name)
	}
	fn, ok := v.(*vm.Var).Deref().(vm.Fn)
	if !ok {
		t.Fatalf("async/%s is not an Fn", name)
	}
	return fn
}

func invoke(t *testing.T, fn vm.Fn, args ...vm.Value) vm.Value {
	t.Helper()
	v, err := fn.Invoke(args)
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	return v
}

// drainChan reads all values from ch until it closes or timeout.
func drainChan(t *testing.T, ch vm.Chan, timeout time.Duration) []vm.Value {
	t.Helper()
	var out []vm.Value
	deadline := time.After(timeout)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, v)
		case <-deadline:
			t.Fatalf("drainChan timed out after %v (got %d values)", timeout, len(out))
			return out
		}
	}
}

// TestAsyncToChanPipe is a correctness regression guard: to-chan! emits a
// collection onto a channel and closes it; pipe forwards src→dst and
// closes dst when src closes.
func TestAsyncToChanPipe(t *testing.T) {
	toChan := asyncFn(t, "to-chan!")
	pipe := asyncFn(t, "pipe")

	src := invoke(t, toChan, vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)})).(vm.Chan)
	dst := make(vm.Chan, 10)
	invoke(t, pipe, src, dst)

	got := drainChan(t, dst, 2*time.Second)
	if len(got) != 3 || got[0] != vm.Int(1) || got[2] != vm.Int(3) {
		t.Fatalf("expected [1 2 3] through pipe, got %v", got)
	}
}

// TestAsyncSplit is a correctness guard for split's routing + close.
func TestAsyncSplit(t *testing.T) {
	toChan := asyncFn(t, "to-chan!")
	split := asyncFn(t, "split")
	isEven, _ := vm.NativeFnType.Wrap(func(a []vm.Value) (vm.Value, error) {
		return vm.Boolean(int(a[0].(vm.Int))%2 == 0), nil
	})

	src := invoke(t, toChan, vm.NewArrayVector([]vm.Value{
		vm.Int(1), vm.Int(2), vm.Int(3), vm.Int(4)})).(vm.Chan)
	res := invoke(t, split, isEven, src)
	pair := res.(vm.ArrayVector)
	trueCh := pair.ValueAt(vm.Int(0)).(vm.Chan)
	falseCh := pair.ValueAt(vm.Int(1)).(vm.Chan)

	// split's channels are unbuffered and it sends interleaved, so read
	// both concurrently to avoid deadlock.
	type res2 struct{ evens, odds []vm.Value }
	ch := make(chan res2, 1)
	go func() {
		var r res2
		r.evens = drainChan(t, trueCh, 2*time.Second)
		ch <- r
	}()
	odds := drainChan(t, falseCh, 2*time.Second)
	evens := (<-ch).evens
	if len(evens) != 2 || len(odds) != 2 {
		t.Fatalf("expected 2 evens + 2 odds, got evens=%v odds=%v", evens, odds)
	}
}

// TestAsyncGoroutinesTrackedAndDrained pins the registry wiring: a
// channel pipeline that blocks forever (no reader on its output) must be
// (a) counted by the VM goroutine registry and (b) released by Drain via
// context cancellation — otherwise it leaks for the life of the process.
func TestAsyncGoroutinesTrackedAndDrained(t *testing.T) {
	// Drain anything left by other tests so the baseline is clean.
	vm.Goroutines.Drain(2 * time.Second)
	base := vm.Goroutines.Live()

	toChan := asyncFn(t, "to-chan!")
	pipe := asyncFn(t, "pipe")

	// Pipe into an UNBUFFERED dst that nobody reads → the pipe goroutine
	// blocks on `dst <- v` forever.
	src := invoke(t, toChan, vm.NewArrayVector([]vm.Value{vm.Int(1), vm.Int(2)})).(vm.Chan)
	dst := make(vm.Chan) // unbuffered, no reader
	invoke(t, pipe, src, dst)

	// Let the goroutines start and block.
	deadline := time.Now().Add(time.Second)
	for vm.Goroutines.Live() <= base && time.Now().Before(deadline) {
		time.Sleep(2 * time.Millisecond)
	}
	if vm.Goroutines.Live() <= base {
		t.Fatalf("expected async goroutines to be tracked (live=%d, base=%d)", vm.Goroutines.Live(), base)
	}

	if !vm.Goroutines.Drain(3 * time.Second) {
		t.Fatal("Drain did not release blocked async goroutines (not ctx-cancellable)")
	}
	if got := vm.Goroutines.Live(); got != 0 {
		t.Fatalf("expected 0 live after drain, got %d", got)
	}
}

// --- promise-chan semantics -------------------------------------------

func promiseChanOps(t *testing.T) (mk, put, take, closef vm.Fn) {
	t.Helper()
	return asyncFn(t, "promise-chan"), asyncFn(t, ">!"), asyncFn(t, "<!"), asyncFn(t, "close!")
}

// TestPromiseChanPutThenTake: a value put is replayed to many takers.
func TestPromiseChanPutThenTake(t *testing.T) {
	mk, put, take, _ := promiseChanOps(t)
	pc := invoke(t, mk)
	invoke(t, put, pc, vm.Int(42))
	for i := 0; i < 3; i++ {
		if got := invoke(t, take, pc); got != vm.Int(42) {
			t.Fatalf("take %d: expected 42, got %v", i, got)
		}
	}
}

// TestPromiseChanTakeBeforePutNoStealRace: a taker parked BEFORE the put
// must (a) receive the value and (b) NOT consume it — later takers still
// see it. This is the race the old single-channel design lost.
func TestPromiseChanTakeBeforePutNoStealRace(t *testing.T) {
	mk, put, take, _ := promiseChanOps(t)
	pc := invoke(t, mk)

	got := make(chan vm.Value, 1)
	go func() {
		v, err := take.Invoke([]vm.Value{pc})
		if err != nil {
			got <- vm.NIL
			return
		}
		got <- v
	}()

	// Give the taker time to park, then deliver.
	time.Sleep(50 * time.Millisecond)
	invoke(t, put, pc, vm.Int(7))

	select {
	case v := <-got:
		if v != vm.Int(7) {
			t.Fatalf("parked taker: expected 7, got %v", v)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("parked taker never woke after put")
	}

	// Value survived: a later taker still sees it.
	if v := invoke(t, take, pc); v != vm.Int(7) {
		t.Fatalf("later taker: expected 7 (replay), got %v", v)
	}
}

// TestPromiseChanSubsequentPutsDropped: only the first put wins.
func TestPromiseChanSubsequentPutsDropped(t *testing.T) {
	mk, put, take, _ := promiseChanOps(t)
	pc := invoke(t, mk)
	invoke(t, put, pc, vm.Int(1))
	invoke(t, put, pc, vm.Int(2)) // dropped
	if v := invoke(t, take, pc); v != vm.Int(1) {
		t.Fatalf("expected first put 1 to win, got %v", v)
	}
}

// TestPromiseChanCloseEmptyYieldsNil: closing without a value → takers
// get nil, not a hang.
func TestPromiseChanCloseEmptyYieldsNil(t *testing.T) {
	mk, _, take, closef := promiseChanOps(t)
	pc := invoke(t, mk)
	invoke(t, closef, pc)
	if v := invoke(t, take, pc); v != vm.NIL {
		t.Fatalf("expected nil from closed-empty promise-chan, got %v", v)
	}
}

// TestPromiseChanValueSurvivesClose: a delivered value keeps being served
// even after close.
func TestPromiseChanValueSurvivesClose(t *testing.T) {
	mk, put, take, closef := promiseChanOps(t)
	pc := invoke(t, mk)
	invoke(t, put, pc, vm.Int(99))
	invoke(t, closef, pc) // no-op for the cached value
	if v := invoke(t, take, pc); v != vm.Int(99) {
		t.Fatalf("expected cached 99 after close, got %v", v)
	}
}

// --- go-block / <! / >! drainability ----------------------------------

func coreOrAsyncFn(t *testing.T, name string) vm.Fn {
	t.Helper()
	v := NS("async").Lookup(vm.Symbol(name))
	if v == nil {
		t.Fatalf("%s not found", name)
	}
	fn, ok := v.(*vm.Var).Deref().(vm.Fn)
	if !ok {
		t.Fatalf("%s is not an Fn", name)
	}
	return fn
}

// TestGoBlockParkedInTakeIsDrainable pins the fix: a (go (<! ch)) parked
// on an empty channel runs in a registry-tracked goroutine, but <! used
// to block on a bare `<-ch` with no view of the registry context — so
// Drain could not release it. After wiring <! to the registry context,
// CancelAll/Drain unblocks the parked take and the go-block exits.
func TestGoBlockParkedInTakeIsDrainable(t *testing.T) {
	vm.Goroutines.Drain(2 * time.Second)
	base := vm.Goroutines.Live()

	goStar := coreOrAsyncFn(t, "go*")
	take := asyncFn(t, "<!")
	ch := make(vm.Chan) // nobody ever delivers

	thunk, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
		return take.Invoke([]vm.Value{ch})
	})
	invoke(t, goStar, thunk) // spawns the go-block goroutine; it parks in <!

	deadline := time.Now().Add(time.Second)
	for vm.Goroutines.Live() <= base && time.Now().Before(deadline) {
		time.Sleep(2 * time.Millisecond)
	}
	if vm.Goroutines.Live() <= base {
		t.Fatalf("go-block goroutine not tracked (live=%d base=%d)", vm.Goroutines.Live(), base)
	}
	// Live++ happens before the goroutine body runs; give it time to
	// actually reach and park in the channel op (reading the current
	// registry context) before we Drain.
	time.Sleep(100 * time.Millisecond)

	if !vm.Goroutines.Drain(3 * time.Second) {
		t.Fatal("Drain could not release a go-block parked in <! (not ctx-aware)")
	}
	if got := vm.Goroutines.Live(); got != 0 {
		t.Fatalf("expected 0 live after drain, got %d", got)
	}
}

// TestGoBlockParkedInPutIsDrainable: same, for a put on a full/unread
// channel.
func TestGoBlockParkedInPutIsDrainable(t *testing.T) {
	vm.Goroutines.Drain(2 * time.Second)
	base := vm.Goroutines.Live()

	goStar := coreOrAsyncFn(t, "go*")
	put := asyncFn(t, ">!")
	ch := make(vm.Chan) // unbuffered, no reader

	thunk, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
		return put.Invoke([]vm.Value{ch, vm.Int(1)})
	})
	invoke(t, goStar, thunk)

	deadline := time.Now().Add(time.Second)
	for vm.Goroutines.Live() <= base && time.Now().Before(deadline) {
		time.Sleep(2 * time.Millisecond)
	}
	if vm.Goroutines.Live() <= base {
		t.Fatalf("go-block goroutine not tracked (live=%d base=%d)", vm.Goroutines.Live(), base)
	}
	// Live++ happens before the goroutine body runs; give it time to
	// actually reach and park in the channel op (reading the current
	// registry context) before we Drain.
	time.Sleep(100 * time.Millisecond)

	if !vm.Goroutines.Drain(3 * time.Second) {
		t.Fatal("Drain could not release a go-block parked in >! (not ctx-aware)")
	}
	if got := vm.Goroutines.Live(); got != 0 {
		t.Fatalf("expected 0 live after drain, got %d", got)
	}
}

// TestTimeoutSurvivesCancelAll covers the generation race: a timeout created
// right after CancelAll must keep its deadline, even when the previous
// generation's daemon has not yet noticed its cancellation and still owns the
// queue. Before the per-entry generation check, the stopping daemon closed
// such a timeout immediately.
func TestTimeoutSurvivesCancelAll(t *testing.T) {
	prev := runtime.GOMAXPROCS(1) // widen the CancelAll→daemon-stop race window
	defer runtime.GOMAXPROCS(prev)

	oldGen := timeoutChan(5000) // parks the daemon on a long deadline
	vm.Goroutines.CancelAll()
	fresh := timeoutChan(200) // new generation, enqueued before the old daemon stops

	select {
	case <-oldGen:
	case <-time.After(2 * time.Second):
		t.Fatal("pre-CancelAll timeout should close promptly on cancellation")
	}

	select {
	case <-fresh:
		t.Fatal("post-CancelAll timeout closed early — it belongs to the fresh generation")
	case <-time.After(80 * time.Millisecond):
	}

	select {
	case <-fresh:
	case <-time.After(2 * time.Second):
		t.Fatal("post-CancelAll timeout never closed at its deadline")
	}
}
