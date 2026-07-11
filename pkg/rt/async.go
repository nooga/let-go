/*
 * async namespace — let-go's equivalent of clojure.core.async
 *
 * Re-exports core async primitives (go, chan, <!, >!) and adds:
 * - close! — close a channel
 * - buffer / chan with buffer size
 * - timeout — channel that closes after N ms
 * - pipe — connect two channels
 */

package rt

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

// Mult broadcasts values from a source channel to multiple tap channels.
type Mult struct {
	src  vm.Chan
	taps map[vm.Chan]bool
	mu   sync.Mutex
}

// Pub routes values from a source channel to subscribers by topic.
type Pub struct {
	src     vm.Chan
	topicFn vm.Fn
	subs    map[any]vm.Chan
	mu      sync.Mutex
}

func init() { RegisterInstaller(installAsyncNS) }

// --- Shared timeout daemon ---
//
// All (timeout ms) channels are closed by ONE long-lived goroutine instead of
// a goroutine per call. A per-call goroutine is cheap on stock Go but not
// free: each (timeout ms) paid a goroutine spawn, its stack, and a timer. On
// runtimes with heap-allocated fixed-size goroutine stacks (TinyGo), a
// timeout-per-frame animation loop turns that into a large allocation plus a
// GC cycle per frame. The daemon makes timeout calls allocation-free apart
// from the channel and a queue entry.
//
// The daemon runs in the root scope and exits on scope cancellation. Each
// entry carries the Done channel of the scope context that was current when
// its timeout was created, preserving the old per-call semantics: a timeout
// closes early only when ITS OWN context generation is cancelled. A stopping
// daemon therefore closes just the entries whose context is done and hands
// any survivors — timeouts created after a CancelAll installed a fresh
// generation — to a respawned daemon in that new generation.

type timerEntry struct {
	deadline time.Time
	ch       vm.Chan
	done     <-chan struct{} // Done of the scope context at creation time
}

// cancelled reports whether the entry's own context generation is done.
func (e *timerEntry) cancelled() bool {
	select {
	case <-e.done:
		return true
	default:
		return false
	}
}

var (
	timerMu      sync.Mutex
	timerQueue   []timerEntry
	timerKick    = make(chan struct{}, 1)
	timerRunning bool
)

// timeoutChan returns a channel that the timer daemon closes after ms.
func timeoutChan(ms int) vm.Chan {
	ch := make(vm.Chan)
	if ms <= 0 {
		close(ch)
		return ch
	}
	entry := timerEntry{
		deadline: time.Now().Add(time.Duration(ms) * time.Millisecond),
		ch:       ch,
		done:     vm.Goroutines.Context().Done(),
	}
	timerMu.Lock()
	timerQueue = append(timerQueue, entry)
	start := !timerRunning
	if start {
		timerRunning = true
	}
	timerMu.Unlock()
	if start {
		vm.Goroutines.Go(timerDaemon)
	} else {
		select {
		case timerKick <- struct{}{}:
		default:
		}
	}
	return ch
}

func timerDaemon(ctx context.Context) {
	for {
		timerMu.Lock()
		now := time.Now()
		var due []vm.Chan
		keep := timerQueue[:0]
		var next time.Time
		for _, e := range timerQueue {
			if !e.deadline.After(now) || e.cancelled() {
				due = append(due, e.ch)
			} else {
				keep = append(keep, e)
				if next.IsZero() || e.deadline.Before(next) {
					next = e.deadline
				}
			}
		}
		timerQueue = keep
		timerMu.Unlock()
		for _, ch := range due {
			close(ch)
		}

		if next.IsZero() {
			// Idle: park until a new timeout arrives or the scope cancels.
			select {
			case <-timerKick:
				continue
			case <-ctx.Done():
				timerDaemonStop()
				return
			}
		}
		t := time.NewTimer(time.Until(next))
		select {
		case <-t.C:
		case <-timerKick: // an earlier deadline may have been queued
		case <-ctx.Done():
			t.Stop()
			timerDaemonStop()
			return
		}
		t.Stop()
	}
}

// timerDaemonStop runs when the daemon's own context generation is cancelled.
// It closes the entries belonging to cancelled generations (unblocking their
// takers, matching the old per-goroutine ctx.Done behavior) but NOT entries
// created after a CancelAll installed a fresh generation — those keep their
// deadlines under a respawned daemon. Without the split, a timeout created in
// the race window between CancelAll and the old daemon noticing would be
// closed immediately by a cancellation that predates it.
func timerDaemonStop() {
	timerMu.Lock()
	var closeNow []vm.Chan
	keep := timerQueue[:0]
	for _, e := range timerQueue {
		if e.cancelled() {
			closeNow = append(closeNow, e.ch)
		} else {
			keep = append(keep, e)
		}
	}
	timerQueue = keep
	respawn := len(keep) > 0
	if !respawn {
		timerRunning = false
	}
	timerMu.Unlock()
	for _, ch := range closeNow {
		close(ch)
	}
	if respawn {
		// Survivors belong to a newer generation; serve them from a fresh
		// daemon spawned under the current scope context. If yet another
		// CancelAll raced in, the new daemon's first pass closes the newly
		// cancelled entries via the per-entry check.
		vm.Goroutines.Go(timerDaemon)
	}
}

// --- Buffer policies (core.async buffer / dropping-buffer / sliding-buffer) ---
//
// vm.Chan is a raw Go channel, which cannot express drop-on-full semantics
// itself, so the policy lives beside the channel (chanPolicy) and is honored
// by every put path (>! in lang.go, offer! here). Per clojure.core.async:
//   (buffer n)          — fixed: puts park when full
//   (dropping-buffer n) — puts always complete; when full the NEW val is
//                         dropped (no transfer)
//   (sliding-buffer n)  — puts always complete; when full the OLDEST
//                         buffered val is dropped
// Buffer constructors return a marker [:policy n] that (chan buf) consumes.

const (
	bufFixed    = 0
	bufDropping = 1
	bufSliding  = 2
)

// chanPolicy maps channels created with a dropping/sliding buffer to their
// policy. Fixed/plain channels are absent (zero-lookup default). Entries are
// removed on close!.
var chanPolicy sync.Map // vm.Chan → int

func policyOf(ch vm.Chan) int {
	if p, ok := chanPolicy.Load(ch); ok {
		return p.(int)
	}
	return bufFixed
}

func bufMarker(policy vm.Keyword, n vm.Int) vm.Value {
	return vm.NewArrayVector([]vm.Value{vm.Keyword("buffer"), policy, n})
}

// asBufMarker recognizes the [:buffer :policy n] marker vectors produced by
// buffer/dropping-buffer/sliding-buffer.
func asBufMarker(v vm.Value) (policy int, n int, ok bool) {
	vec, isVec := v.(vm.ArrayVector)
	if !isVec || len(vec) != 3 {
		return 0, 0, false
	}
	if kw, isKw := vec[0].(vm.Keyword); !isKw || kw != vm.Keyword("buffer") {
		return 0, 0, false
	}
	size, isInt := vec[2].(vm.Int)
	if !isInt {
		return 0, 0, false
	}
	switch vec[1] {
	case vm.Keyword("fixed"):
		return bufFixed, int(size), true
	case vm.Keyword("dropping"):
		return bufDropping, int(size), true
	case vm.Keyword("sliding"):
		return bufSliding, int(size), true
	}
	return 0, 0, false
}

// putWithPolicy sends v on ch honoring its buffer policy. Returns
// (accepted, cancelled). Parity with core.async: a put on a CLOSED channel
// returns accepted=false instead of panicking (>!! "returns true unless the
// channel is already closed"); dropping/sliding puts always complete.
func putWithPolicy(ctx context.Context, ch vm.Chan, v vm.Value) (accepted bool, cancelled bool) {
	defer func() {
		if r := recover(); r != nil {
			// send on closed channel → put returns false (core.async parity)
			accepted, cancelled = false, false
		}
	}()
	switch policyOf(ch) {
	case bufDropping:
		select {
		case ch <- v:
		default: // full → drop the NEW value, put still "completes"
		}
		return true, false
	case bufSliding:
		for {
			select {
			case ch <- v:
				return true, false
			default:
			}
			// full → evict the OLDEST, then retry the send (bounded spin:
			// each lap either sends or evicts; contention only with other
			// producers on the same channel)
			select {
			case <-ch:
			default:
			}
			select {
			case ch <- v:
				return true, false
			case <-ctx.Done():
				return false, true
			default: // raced with another producer; loop
			}
		}
	default:
		select {
		case ch <- v:
			return true, false
		case <-ctx.Done():
			return false, true
		}
	}
}

// PromiseChan implements clojure.core.async's promise-chan: a channel
// that caches the FIRST value put to it and replays that value to every
// taker, forever. Subsequent puts are dropped; closing without a value
// makes takers receive nil; closing after a value is a no-op (the value
// keeps being served).
//
// Unlike a raw vm.Chan it STORES the value rather than transferring it,
// which is what makes the semantics correct: with a single raw channel a
// taker parked before the first put could steal that put before it was
// cached, so later takers would never see it. Storing the value behind a
// latch removes that race entirely.
//
// Dispatched to by >! / <! / close! when they receive a boxed
// *PromiseChan instead of a vm.Chan. Methods are unexported so the
// reflective method-boxing in vm.NewBoxed skips them.
type PromiseChan struct {
	mu     sync.Mutex
	value  vm.Value
	set    bool
	closed bool
	ready  chan struct{} // closed once set or closed; latch for takers
}

func newPromiseChan() *PromiseChan {
	return &PromiseChan{value: vm.NIL, ready: make(chan struct{})}
}

// put caches v if no value has been delivered and the chan is open;
// otherwise it is dropped (first put wins).
func (p *PromiseChan) put(v vm.Value) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.set || p.closed {
		return
	}
	p.value = v
	p.set = true
	close(p.ready)
}

// take returns the cached value, blocking until one is delivered or the
// chan is closed. ctx (the registry context) lets a blocked take be
// drained on shutdown; cancellation returns nil.
func (p *PromiseChan) take(ctx context.Context) vm.Value {
	p.mu.Lock()
	if p.set || p.closed {
		v := p.value
		p.mu.Unlock()
		return v
	}
	p.mu.Unlock()
	select {
	case <-p.ready:
	case <-ctx.Done():
		return vm.NIL
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.value // the delivered value, or NIL if closed empty
}

// doClose marks the chan closed. No-op once a value is set, so the value
// keeps being served to takers.
func (p *PromiseChan) doClose() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.set || p.closed {
		return
	}
	p.closed = true
	close(p.ready)
}

// asPromiseChan returns the *PromiseChan a value wraps, if any.
func asPromiseChan(v vm.Value) (*PromiseChan, bool) {
	b, ok := v.(*vm.Boxed)
	if !ok {
		return nil, false
	}
	pc, ok := b.Unbox().(*PromiseChan)
	return pc, ok
}

// ctxSend sends v on ch, aborting if ctx is cancelled (e.g. a registry
// Drain on shutdown / between bench iterations). Returns false if the
// send was abandoned due to cancellation.
func ctxSend(ctx context.Context, ch vm.Chan, v vm.Value) bool {
	select {
	case ch <- v:
		return true
	case <-ctx.Done():
		return false
	}
}

// ctxRecv receives from ch. ok is false when ch is closed; live is false
// when ctx was cancelled before a value arrived (caller should return).
func ctxRecv(ctx context.Context, ch vm.Chan) (v vm.Value, ok bool, live bool) {
	select {
	case rv, rok := <-ch:
		return rv, rok, true
	case <-ctx.Done():
		return vm.NIL, false, false
	}
}

// nolint
func installAsyncNS() {
	// Look up the core builtins to re-export
	coreNS := nsRegistry[NameCoreNS]

	// close! — close a channel
	closeChan, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("close! expects 1 arg")
		}
		if pc, ok := asPromiseChan(vs[0]); ok {
			pc.doClose()
			return vm.NIL, nil
		}
		ch, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("close! expected Chan")
		}
		chanPolicy.Delete(ch)
		func() {
			// core.async close! is a no-op on an already-closed channel
			defer func() { _ = recover() }()
			close(ch)
		}()
		return vm.NIL, nil
	})

	// chan: (chan) unbuffered, (chan n) fixed buffer, (chan (buffer n)) /
	// (chan (dropping-buffer n)) / (chan (sliding-buffer n)) per core.async.
	chanBuf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return make(vm.Chan), nil
		}
		if n, ok := vs[0].(vm.Int); ok {
			return make(vm.Chan, int(n)), nil
		}
		if policy, n, ok := asBufMarker(vs[0]); ok {
			if n < 1 {
				return vm.NIL, fmt.Errorf("chan buffer size must be >= 1")
			}
			ch := make(vm.Chan, n)
			if policy != bufFixed {
				chanPolicy.Store(ch, policy)
			}
			return ch, nil
		}
		return vm.NIL, fmt.Errorf("chan expected Int size or buffer marker")
	})

	// Buffer constructors (markers consumed by chan; see asBufMarker)
	fixedBuf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("buffer expects 1 arg (n)")
		}
		n, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("buffer expected Int size")
		}
		return bufMarker(vm.Keyword("fixed"), n), nil
	})
	droppingBuf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("dropping-buffer expects 1 arg (n)")
		}
		n, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("dropping-buffer expected Int size")
		}
		return bufMarker(vm.Keyword("dropping"), n), nil
	})
	slidingBuf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("sliding-buffer expects 1 arg (n)")
		}
		n, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("sliding-buffer expected Int size")
		}
		return bufMarker(vm.Keyword("sliding"), n), nil
	})

	// timeout — returns a channel that closes after n milliseconds
	timeout, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("timeout expects 1 arg (ms)")
		}
		ms, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("timeout expected Int milliseconds")
		}
		return timeoutChan(int(ms)), nil
	})

	// pipe — take from src, put on dst, close dst when src closes
	// (pipe src dst) or (pipe src dst close?)
	pipe, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("pipe expects 2-3 args")
		}
		src, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("pipe expected Chan src")
		}
		dst, ok := vs[1].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("pipe expected Chan dst")
		}
		shouldClose := true
		if len(vs) == 3 {
			shouldClose = vm.IsTruthy(vs[2])
		}
		vm.Goroutines.Go(func(ctx context.Context) {
			for {
				v, ok, live := ctxRecv(ctx, src)
				if !live {
					return
				}
				if !ok {
					break
				}
				if !ctxSend(ctx, dst, v) {
					return
				}
			}
			if shouldClose {
				close(dst)
			}
		})
		return dst, nil
	})

	// onto-chan! — put all items from coll onto ch, then close
	// (onto-chan! ch coll) or (onto-chan! ch coll close?)
	ontoChan, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("onto-chan! expects 2-3 args")
		}
		ch, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("onto-chan! expected Chan")
		}
		shouldClose := true
		if len(vs) == 3 {
			shouldClose = vm.IsTruthy(vs[2])
		}
		seq, ok := vs[1].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("onto-chan! expected Sequable")
		}
		vm.Goroutines.Go(func(ctx context.Context) {
			for s := seq.Seq(); s != nil; s = s.Next() {
				if !ctxSend(ctx, ch, s.First()) {
					return
				}
			}
			if shouldClose {
				close(ch)
			}
		})
		return ch, nil
	})

	// merge — take from multiple channels, put onto one output channel
	// (merge chs) or (merge chs buf-size)
	mergef, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("merge expects 1-2 args")
		}
		seq, ok := vs[0].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("merge expected sequable of channels")
		}
		bufSize := 0
		if len(vs) == 2 {
			if n, ok := vs[1].(vm.Int); ok {
				bufSize = int(n)
			}
		}
		out := make(vm.Chan, bufSize)
		// Count channels and start goroutines
		done := make(chan struct{})
		count := 0
		for s := seq.Seq(); s != nil; s = s.Next() {
			ch, ok := s.First().(vm.Chan)
			if !ok {
				continue
			}
			count++
			c := ch
			vm.Goroutines.Go(func(ctx context.Context) {
				for {
					v, ok, live := ctxRecv(ctx, c)
					if !live {
						return
					}
					if !ok {
						break
					}
					if !ctxSend(ctx, out, v) {
						return
					}
				}
				select {
				case done <- struct{}{}:
				case <-ctx.Done():
				}
			})
		}
		// Close output when all inputs are done
		vm.Goroutines.Go(func(ctx context.Context) {
			for range count {
				select {
				case <-done:
				case <-ctx.Done():
					return
				}
			}
			close(out)
		})
		return out, nil
	})

	// reduce — async reduce: (async/reduce f init ch) → channel with result
	reducef := vm.NewCtxNativeFn("reduce", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("async/reduce expects 3 args")
		}
		fn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("async/reduce expected Fn")
		}
		init := vs[1]
		ch, ok := vs[2].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("async/reduce expected Chan")
		}
		out := make(vm.Chan, 1)
		// Convey the caller's bindings into the loop goroutine (like future).
		childEc := ec.Child()
		vm.Goroutines.Go(func(ctx context.Context) {
			acc := init
			for {
				v, ok, live := ctxRecv(ctx, ch)
				if !live {
					return
				}
				if !ok {
					break
				}
				result, err := childEc.Invoke(fn, []vm.Value{acc, v})
				if err != nil {
					break
				}
				acc = result
			}
			out <- acc // out is buffered (cap 1); never blocks
			close(out)
		})
		return out, nil
	})

	// into — async into: (async/into coll ch) → channel with result
	intof, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("async/into expects 2 args")
		}
		coll := vs[0]
		ch, ok := vs[1].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("async/into expected Chan")
		}
		out := make(vm.Chan, 1)
		vm.Goroutines.Go(func(ctx context.Context) {
			acc := coll
			for {
				v, ok, live := ctxRecv(ctx, ch)
				if !live {
					return
				}
				if !ok {
					break
				}
				if assoc, ok := acc.(vm.Collection); ok {
					acc = assoc.Conj(v)
				}
			}
			out <- acc // out is buffered (cap 1); never blocks
			close(out)
		})
		return out, nil
	})

	// to-chan! — create a channel with items from coll
	toChan, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("to-chan! expects 1 arg")
		}
		seq, ok := vs[0].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("to-chan! expected Sequable")
		}
		ch := make(vm.Chan)
		vm.Goroutines.Go(func(ctx context.Context) {
			for s := seq.Seq(); s != nil; s = s.Next() {
				if !ctxSend(ctx, ch, s.First()) {
					return
				}
			}
			close(ch)
		})
		return ch, nil
	})

	// alts! — select on multiple channel operations
	// (alts! [ch1 ch2 [ch3 val]] :default v) → [val port]
	// Each entry is either a channel (take) or [channel value] (put).
	// With :default, alts! never parks: if no port is immediately ready it
	// returns [v :default]. This is the non-blocking form the single-threaded
	// wasm event loop needs, since a parked take freezes the whole runtime.
	altsf := vm.NewCtxNativeFn("alts!", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("alts! expects at least 1 arg (vector of ports)")
		}
		seq, ok := vs[0].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("alts! expected sequable of ports")
		}

		// Trailing options as keyword/value pairs. Only :default is honored.
		var hasDefault bool
		var defaultVal vm.Value = vm.NIL
		for i := 1; i+1 < len(vs); i += 2 {
			if kw, ok := vs[i].(vm.Keyword); ok && kw == vm.Keyword("default") {
				hasDefault = true
				defaultVal = vs[i+1]
			}
		}

		var cases []reflect.SelectCase
		var ports []vm.Value // parallel array: the channel value for each case

		for s := seq.Seq(); s != nil; s = s.Next() {
			item := s.First()

			// [ch val] — put operation
			if vec, ok := item.(vm.Sequable); ok {
				vs := vec.Seq()
				first := vs.First()
				if ch, ok := first.(vm.Chan); ok {
					nxt := vs.Next()
					if nxt != nil {
						// It's a put: [ch val]
						val := nxt.First()
						cases = append(cases, reflect.SelectCase{
							Dir:  reflect.SelectSend,
							Chan: reflect.ValueOf((chan vm.Value)(ch)),
							Send: reflect.ValueOf(val),
						})
						ports = append(ports, ch)
						continue
					}
				}
			}

			// Plain channel — take operation
			if ch, ok := item.(vm.Chan); ok {
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf((chan vm.Value)(ch)),
				})
				ports = append(ports, ch)
				continue
			}

			return vm.NIL, fmt.Errorf("alts! expected channel or [channel value], got %s", item.Type().Name())
		}

		if len(cases) == 0 {
			return vm.NIL, fmt.Errorf("alts! requires at least one port")
		}

		// With :default, add a SelectDefault case so reflect.Select returns
		// immediately when no port is ready — never parking. The ctx Done
		// case is unneeded here: the non-blocking path can't block, so there
		// is nothing for a cancel to release.
		//
		// Without :default, append a recv on the registry context's Done
		// channel so an alts! parked on its ports — e.g. inside a (go ...)
		// block — is released by a CancelAll/Drain on shutdown. If that case
		// wins, return nil (no port chosen).
		sentinelIdx := len(cases)
		if hasDefault {
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
		} else {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ec.Context().Done()),
			})
		}

		chosen, value, ok := reflect.Select(cases)
		if chosen == sentinelIdx {
			if hasDefault {
				return vm.NewArrayVector([]vm.Value{defaultVal, vm.Keyword("default")}), nil
			}
			return vm.NIL, nil // ctx cancelled
		}
		port := ports[chosen]

		var result vm.Value
		if cases[chosen].Dir == reflect.SelectRecv {
			if ok {
				result = value.Interface().(vm.Value)
			} else {
				result = vm.NIL // channel closed
			}
		} else {
			// Put operation — result is true if successful
			result = vm.TRUE
		}

		return vm.NewArrayVector([]vm.Value{result, port}), nil
	})

	// offer! — non-blocking put: true if accepted, nil if not (core.async
	// returns nil, not false, when the offer can't complete). On dropping/
	// sliding channels a put can always complete immediately.
	offerf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("offer! expects 2 args")
		}
		ch, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("offer! expected Chan")
		}
		if vs[1] == vm.NIL {
			// core.async: nil values are not allowed on channels (same as >!).
			// Keeping the channel nil-free is what makes poll!'s value-or-nil
			// contract unambiguous.
			return vm.NIL, fmt.Errorf("offer! can't put nil on chan")
		}
		switch policyOf(ch) {
		case bufDropping, bufSliding:
			if accepted, _ := putWithPolicy(context.Background(), ch, vs[1]); accepted {
				return vm.TRUE, nil
			}
			return vm.NIL, nil // closed
		default:
			ok := func() (ok bool) {
				defer func() {
					if recover() != nil {
						ok = false // offer on closed channel → nil, not panic
					}
				}()
				select {
				case ch <- vs[1]:
					return true
				default:
					return false
				}
			}()
			if ok {
				return vm.TRUE, nil
			}
			return vm.NIL, nil
		}
	})

	// poll! — non-blocking take, returns value or nil
	pollf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("poll! expects 1 arg")
		}
		ch, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("poll! expected Chan")
		}
		select {
		case v, ok := <-ch:
			if ok {
				return v, nil
			}
			return vm.NIL, nil
		default:
			return vm.NIL, nil
		}
	})

	// promise-chan — a channel that caches the first value put and
	// replays it to every taker (see PromiseChan). Returned boxed; >! /
	// <! / close! dispatch to it. No goroutine needed — the value is
	// stored behind a latch, so there is no parked-taker steal race.
	promiseChan, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("promise-chan expects 0 args")
		}
		return vm.NewBoxed(newPromiseChan()), nil
	})

	// mult — create a mult (broadcast) from a source channel
	multf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("mult expects 1 arg")
		}
		src, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("mult expected Chan")
		}
		m := &Mult{src: src, taps: make(map[vm.Chan]bool)}
		vm.Goroutines.Go(func(ctx context.Context) {
			for {
				v, ok, live := ctxRecv(ctx, src)
				if !live {
					return
				}
				if !ok {
					break
				}
				m.mu.Lock()
				for ch, closeCh := range m.taps {
					select {
					case ch <- v:
					default:
						// drop if tap is full
					}
					_ = closeCh
				}
				m.mu.Unlock()
			}
			// Source closed — close all taps that requested it
			m.mu.Lock()
			for ch, shouldClose := range m.taps {
				if shouldClose {
					close(ch)
				}
			}
			m.mu.Unlock()
		})
		return vm.NewBoxed(m), nil
	})

	// tap — add a channel to a mult
	tapf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("tap expects 2-3 args")
		}
		b, ok := vs[0].(*vm.Boxed)
		if !ok {
			return vm.NIL, fmt.Errorf("tap expected Mult")
		}
		m, ok := b.Unbox().(*Mult)
		if !ok {
			return vm.NIL, fmt.Errorf("tap expected Mult")
		}
		ch, ok := vs[1].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("tap expected Chan")
		}
		shouldClose := true
		if len(vs) == 3 {
			shouldClose = vm.IsTruthy(vs[2])
		}
		m.mu.Lock()
		m.taps[ch] = shouldClose
		m.mu.Unlock()
		return ch, nil
	})

	// untap — remove a channel from a mult
	untapf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("untap expects 2 args")
		}
		b, ok := vs[0].(*vm.Boxed)
		if !ok {
			return vm.NIL, fmt.Errorf("untap expected Mult")
		}
		m, ok := b.Unbox().(*Mult)
		if !ok {
			return vm.NIL, fmt.Errorf("untap expected Mult")
		}
		ch, ok := vs[1].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("untap expected Chan")
		}
		m.mu.Lock()
		delete(m.taps, ch)
		m.mu.Unlock()
		return vm.NIL, nil
	})

	// untap-all — remove all taps from a mult
	untapAllf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("untap-all expects 1 arg")
		}
		b, ok := vs[0].(*vm.Boxed)
		if !ok {
			return vm.NIL, fmt.Errorf("untap-all expected Mult")
		}
		m, ok := b.Unbox().(*Mult)
		if !ok {
			return vm.NIL, fmt.Errorf("untap-all expected Mult")
		}
		m.mu.Lock()
		m.taps = make(map[vm.Chan]bool)
		m.mu.Unlock()
		return vm.NIL, nil
	})

	// pub — create a pub from a source channel with a topic fn
	pubf := vm.NewCtxNativeFn("pub", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("pub expects 2 args (ch, topic-fn)")
		}
		src, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("pub expected Chan")
		}
		topicFn, ok := vs[1].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("pub expected Fn")
		}
		p := &Pub{src: src, topicFn: topicFn, subs: make(map[any]vm.Chan)}
		// Convey the caller's bindings into the loop goroutine (like future).
		childEc := ec.Child()
		vm.Goroutines.Go(func(ctx context.Context) {
			for {
				v, ok, live := ctxRecv(ctx, src)
				if !live {
					return
				}
				if !ok {
					break
				}
				topic, err := childEc.Invoke(topicFn, []vm.Value{v})
				if err != nil {
					continue
				}
				key := topic.Unbox()
				p.mu.Lock()
				if ch, ok := p.subs[key]; ok {
					select {
					case ch <- v:
					default:
					}
				}
				p.mu.Unlock()
			}
			// Source closed — close all sub channels
			p.mu.Lock()
			for _, ch := range p.subs {
				close(ch)
			}
			p.mu.Unlock()
		})
		return vm.NewBoxed(p), nil
	})

	// sub — subscribe to a topic on a pub
	subf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("sub expects 3 args (pub, topic, ch)")
		}
		b, ok := vs[0].(*vm.Boxed)
		if !ok {
			return vm.NIL, fmt.Errorf("sub expected Pub")
		}
		p, ok := b.Unbox().(*Pub)
		if !ok {
			return vm.NIL, fmt.Errorf("sub expected Pub")
		}
		topic := vs[1].Unbox()
		ch, ok := vs[2].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("sub expected Chan")
		}
		p.mu.Lock()
		p.subs[topic] = ch
		p.mu.Unlock()
		return ch, nil
	})

	// unsub — unsubscribe from a topic
	unsubf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("unsub expects 2 args (pub, topic)")
		}
		b, ok := vs[0].(*vm.Boxed)
		if !ok {
			return vm.NIL, fmt.Errorf("unsub expected Pub")
		}
		p, ok := b.Unbox().(*Pub)
		if !ok {
			return vm.NIL, fmt.Errorf("unsub expected Pub")
		}
		topic := vs[1].Unbox()
		p.mu.Lock()
		delete(p.subs, topic)
		p.mu.Unlock()
		return vm.NIL, nil
	})

	// split — route values from ch into two channels based on predicate
	splitf := vm.NewCtxNativeFn("split", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 4 {
			return vm.NIL, fmt.Errorf("split expects 2-4 args")
		}
		pred, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("split expected Fn predicate")
		}
		src, ok := vs[1].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("split expected Chan")
		}
		trueCh := make(vm.Chan)
		falseCh := make(vm.Chan)
		// Convey the caller's bindings into the loop goroutine (like future).
		childEc := ec.Child()
		vm.Goroutines.Go(func(ctx context.Context) {
			for {
				v, ok, live := ctxRecv(ctx, src)
				if !live {
					return
				}
				if !ok {
					break
				}
				result, err := childEc.Invoke(pred, []vm.Value{v})
				if err != nil || !vm.IsTruthy(result) {
					if !ctxSend(ctx, falseCh, v) {
						return
					}
				} else {
					if !ctxSend(ctx, trueCh, v) {
						return
					}
				}
			}
			close(trueCh)
			close(falseCh)
		})
		return vm.NewArrayVector([]vm.Value{trueCh, falseCh}), nil
	})

	// async/map — apply f to values taken from multiple channels simultaneously
	mapf := vm.NewCtxNativeFn("map", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("async/map expects 2 args (f, chs)")
		}
		fn, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, fmt.Errorf("async/map expected Fn")
		}
		seq, ok := vs[1].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("async/map expected sequable of channels")
		}
		var chs []vm.Chan
		for s := seq.Seq(); s != nil; s = s.Next() {
			ch, ok := s.First().(vm.Chan)
			if !ok {
				return vm.NIL, fmt.Errorf("async/map expected channels")
			}
			chs = append(chs, ch)
		}
		out := make(vm.Chan)
		// Convey the caller's bindings into the loop goroutine (like future).
		childEc := ec.Child()
		vm.Goroutines.Go(func(ctx context.Context) {
			for {
				args := make([]vm.Value, len(chs))
				allOk := true
				for i, ch := range chs {
					v, ok, live := ctxRecv(ctx, ch)
					if !live {
						return
					}
					if !ok {
						allOk = false
						break
					}
					args[i] = v
				}
				if !allOk {
					break
				}
				result, err := childEc.Invoke(fn, args)
				if err != nil {
					break
				}
				if !ctxSend(ctx, out, result) {
					return
				}
			}
			close(out)
		})
		return out, nil
	})

	// async/take — take n values from ch, put on new channel
	takef, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("async/take expects 2 args (n, ch)")
		}
		n, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("async/take expected Int")
		}
		ch, ok := vs[1].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("async/take expected Chan")
		}
		out := make(vm.Chan)
		vm.Goroutines.Go(func(ctx context.Context) {
			count := int(n)
			for range count {
				v, ok, live := ctxRecv(ctx, ch)
				if !live {
					return
				}
				if !ok {
					break
				}
				if !ctxSend(ctx, out, v) {
					return
				}
			}
			close(out)
		})
		return out, nil
	})

	ns := vm.NewNamespace("async")
	ns.Refer(CoreNS, "", true)

	// Intentional shadows of clojure.core names — suppress warn-on-shadow.
	for _, n := range []string{
		"go*", ">!", "<!", "chan", "close!", "split", "reduce",
		">!!", "<!!", "map", "take", "merge", "into",
	} {
		ns.Exclude(n)
	}

	// Re-export core primitives (extract root value from Var)
	ns.Def("go*", coreNS.Lookup("go*").(*vm.Var).Deref())
	ns.Def(">!", coreNS.Lookup(">!").(*vm.Var).Deref())
	ns.Def("<!", coreNS.Lookup("<!").(*vm.Var).Deref())

	// New async-specific fns
	ns.Def("chan", chanBuf)
	ns.Def("close!", closeChan)
	ns.Def("timeout", timeout)
	ns.Def("pipe", pipe)
	ns.Def("onto-chan!", ontoChan)
	ns.Def("to-chan!", toChan)
	ns.Def("alts!", altsf)
	// alts!! — blocking variant. In let-go all channel ops block the calling
	// goroutine (there is no separate parking machinery), so like >!!/<!!
	// below it is the same fn under the core.async-parity name.
	ns.Def("alts!!", altsf)
	ns.Def("offer!", offerf)
	ns.Def("poll!", pollf)
	ns.Def("buffer", fixedBuf)
	ns.Def("dropping-buffer", droppingBuf)
	ns.Def("sliding-buffer", slidingBuf)
	ns.Def("promise-chan", promiseChan)
	ns.Def("mult", multf)
	ns.Def("tap", tapf)
	ns.Def("untap", untapf)
	ns.Def("untap-all", untapAllf)
	ns.Def("pub", pubf)
	ns.Def("sub", subf)
	ns.Def("unsub", unsubf)
	ns.Def("split", splitf)
	ns.Def("map", mapf)
	ns.Def("take", takef)
	ns.Def("merge", mergef)
	ns.Def("reduce", reducef)
	ns.Def("into", intof)

	// Blocking aliases (in let-go all ops block, so these are identical)
	ns.Def(">!!", coreNS.Lookup(">!").(*vm.Var).Deref())
	ns.Def("<!!", coreNS.Lookup("<!").(*vm.Var).Deref())

	RegisterNS(ns)
}
