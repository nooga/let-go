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
	subs    map[interface{}]vm.Chan
	mu      sync.Mutex
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
		go func() {
			time.Sleep(time.Duration(int(ms)) * time.Millisecond)
			close(ch)
		}()
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
		go func() {
			for v := range src {
				dst <- v
			}
			if shouldClose {
				close(dst)
			}
		}()
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
		go func() {
			for s := seq.Seq(); s != nil; s = s.Next() {
				ch <- s.First()
			}
			if shouldClose {
				close(ch)
			}
		}()
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
			go func(c vm.Chan) {
				for v := range c {
					out <- v
				}
				done <- struct{}{}
			}(ch)
		}
		// Close output when all inputs are done
		go func() {
			for range count {
				<-done
			}
			close(out)
		}()
		return out, nil
	})

	// reduce — async reduce: (async/reduce f init ch) → channel with result
	reducef, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
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
		go func() {
			acc := init
			for v := range ch {
				result, err := fn.Invoke([]vm.Value{acc, v})
				if err != nil {
					break
				}
				acc = result
			}
			out <- acc
			close(out)
		}()
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
		go func() {
			acc := coll
			for v := range ch {
				if assoc, ok := acc.(vm.Collection); ok {
					acc = assoc.Conj(v)
				}
			}
			out <- acc
			close(out)
		}()
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
		go func() {
			for s := seq.Seq(); s != nil; s = s.Next() {
				ch <- s.First()
			}
			close(ch)
		}()
		return ch, nil
	})

	// alts! — select on multiple channel operations
	// (alts! [ch1 ch2 [ch3 val]]) → [val port]
	// Each entry is either a channel (take) or [channel value] (put).
	altsf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
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

		chosen, value, ok := reflect.Select(cases)
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

	// promise-chan — channel that caches and replays one value to all takers
	promiseChan, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		ch := make(vm.Chan)
		p := vm.NewPromise()
		// Deliver goroutine: first put delivers, then replay forever
		go func() {
			v, ok := <-ch
			if ok {
				p.Deliver(v)
			} else {
				p.Deliver(vm.NIL)
			}
		}()
		// Wrap in a boxed value that acts like a channel but reads from promise
		// Simpler approach: return a special channel that replays
		out := make(vm.Chan)
		go func() {
			// Wait for the value
			val := p.Deref()
			// Now replay it to anyone who reads
			for {
				out <- val
			}
		}()
		// We need both ends: write to ch, read from out
		// Package them together. Actually, simplest: use a proxy channel
		proxy := make(vm.Chan)
		delivered := make(chan struct{}, 1)
		var cached vm.Value
		go func() {
			// First value written becomes the cached value
			v, ok := <-proxy
			if ok {
				cached = v
			} else {
				cached = vm.NIL
			}
			close(delivered)
		}()
		go func() {
			<-delivered
			// Now serve reads forever
			for {
				proxy <- cached
			}
		}()
		// Close the unused channels
		close(ch)
		close(out)
		return proxy, nil
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
		go func() {
			for v := range src {
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
		}()
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
	pubf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
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
		p := &Pub{src: src, topicFn: topicFn, subs: make(map[interface{}]vm.Chan)}
		go func() {
			for v := range src {
				topic, err := topicFn.Invoke([]vm.Value{v})
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
		}()
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
	splitf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
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
		go func() {
			for v := range src {
				result, err := pred.Invoke([]vm.Value{v})
				if err != nil || !vm.IsTruthy(result) {
					falseCh <- v
				} else {
					trueCh <- v
				}
			}
			close(trueCh)
			close(falseCh)
		}()
		return vm.NewArrayVector([]vm.Value{trueCh, falseCh}), nil
	})

	// async/map — apply f to values taken from multiple channels simultaneously
	mapf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
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
		go func() {
			for {
				args := make([]vm.Value, len(chs))
				allOk := true
				for i, ch := range chs {
					v, ok := <-ch
					if !ok {
						allOk = false
						break
					}
					args[i] = v
				}
				if !allOk {
					break
				}
				result, err := fn.Invoke(args)
				if err != nil {
					break
				}
				out <- result
			}
			close(out)
		}()
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
		go func() {
			count := int(n)
			for i := 0; i < count; i++ {
				v, ok := <-ch
				if !ok {
					break
				}
				out <- v
			}
			close(out)
		}()
		return out, nil
	})

	ns := vm.NewNamespace("async")
	ns.Refer(CoreNS, "", true)

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
