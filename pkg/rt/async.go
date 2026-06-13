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
		close(ch)
		return vm.NIL, nil
	})

	// chan with optional buffer size: (chan) or (chan n)
	chanBuf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 0 {
			return make(vm.Chan), nil
		}
		n, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("chan expected Int buffer size")
		}
		return make(vm.Chan, int(n)), nil
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
		ch := make(vm.Chan)
		vm.Goroutines.Go(func(ctx context.Context) {
			t := time.NewTimer(time.Duration(int(ms)) * time.Millisecond)
			defer t.Stop()
			select {
			case <-t.C:
			case <-ctx.Done():
			}
			close(ch)
		})
		return ch, nil
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
	// (alts! [ch1 ch2 [ch3 val]]) → [val port]
	// Each entry is either a channel (take) or [channel value] (put).
	altsf := vm.NewCtxNativeFn("alts!", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("alts! expects 1 arg (vector of ports)")
		}
		seq, ok := vs[0].(vm.Sequable)
		if !ok {
			return vm.NIL, fmt.Errorf("alts! expected sequable of ports")
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

		// Append a recv on the registry context's Done channel so an
		// alts! parked on its ports — e.g. inside a (go ...) block — is
		// released by a CancelAll/Drain on shutdown. If that case wins,
		// return nil (no port chosen).
		cancelIdx := len(cases)
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ec.Context().Done()),
		})

		chosen, value, ok := reflect.Select(cases)
		if chosen == cancelIdx {
			return vm.NIL, nil
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

	// offer! — non-blocking put, returns true if accepted, false if not
	offerf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("offer! expects 2 args")
		}
		ch, ok := vs[0].(vm.Chan)
		if !ok {
			return vm.NIL, fmt.Errorf("offer! expected Chan")
		}
		select {
		case ch <- vs[1]:
			return vm.TRUE, nil
		default:
			return vm.FALSE, nil
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
	ns.Def("offer!", offerf)
	ns.Def("poll!", pollf)
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
