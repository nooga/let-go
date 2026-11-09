//go:build js && wasm

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"sync/atomic"
	"syscall/js"
	"time"

	"github.com/nooga/let-go/pkg/rt"
)

// The browser backend: a pump goroutine renders the mixer in small slices
// into a ring, and an AudioWorklet on the page copies frames out of it. The
// page side (the AudioContext, the worklet, the gesture unlock) lives in the
// bundle's lg-host-core.js, behind _lgAudioOpen, so a Go program only ever
// fills a ring. lg -w -w-audio links this file in.
//
// The pump renders ahead instead of on demand: Go wasm runs every goroutine
// on one thread without preemption, so rendering inside the worklet's refill
// deadline overruns it on a phone and crackles (what oto's own js backend
// does). Rendering a slice at a time within a budget, then sleeping, lets the
// queue absorb a slow slice or a busy VM. This design comes from undertone's
// webaudio package, which plays tesserae's browser build; measured there on
// an iPhone over HTTPS, clean with 0 underruns (2026-10).
//
// Where SharedArrayBuffer is available (a cross-origin isolated page, which
// the bundle needs for its worker mode anyway) the ring is shared memory the
// pump writes straight into. Elsewhere slices reach the worklet as
// transferred messages, and the worklet reports back what it has played.

const (
	webSlice  = 512                  // frames rendered per step
	webBudget = 3 * time.Millisecond // longest render run before yielding
)

// Default queues per path, used unless one was asked for (Options.Queue,
// LETGO_AUDIO_QUEUE). They follow what undertone measured for its browser
// output (2026-10-31): 60 ms with a shared ring on the page's thread, 80 ms
// over messages, whose arrival is less even, and 120 ms in a worker, where
// the pump shares one thread with the whole VM.
const (
	webQueueShared  = 60 * time.Millisecond
	webQueueMessage = 80 * time.Millisecond
	webQueueWorker  = 120 * time.Millisecond
)

// webOut is the pump's side of one ring.
type webOut struct {
	src   io.Reader
	queue int64 // frames to keep queued
	cap   int   // ring size in frames, a whole number of slices

	shared  bool
	ctrl    js.Value   // shared: Int32Array of [write, read, underruns, page state]
	views   []js.Value // shared: one Uint8Array per slice-sized slot of the ring
	atomics js.Value

	handle js.Value // message path: the page's handle, with push(), played and state
	sent   int64    // message path: frames posted

	lost     atomic.Bool // the page reported it can't play
	raw, f32 []byte      // one slice of mixer frames, and as float32
}

// The page's state, which lg-host-core.js publishes in ctrl[3] or
// handle.state. The pump fills the ring only while the page plays it.
const (
	pageOpening = 0 // the worklet is loading
	pageRunning = 1 // the context runs and the worklet is connected
	pageStopped = 2 // suspended: before the first gesture, or interrupted
	pageFailed  = 3 // no AudioContext or worklet; nothing will ever play
)

func (d *Device) start(l latency, asked bool) (any, error) {
	g := js.Global()
	open := g.Get("_lgAudioOpen")
	if open.Type() != js.TypeFunction {
		return nil, errors.New("audio: this page has no audio host; build the bundle with lg -w -w-audio")
	}
	worker := g.Get("document").IsUndefined()
	shared := g.Get("crossOriginIsolated").Truthy() && !g.Get("SharedArrayBuffer").IsUndefined()
	if worker && !shared {
		return nil, errors.New("audio: playing from a worker needs a cross-origin isolated page (SharedArrayBuffer)")
	}
	queue := l.queue
	if !asked {
		switch {
		case worker:
			queue = webQueueWorker
		case shared:
			queue = webQueueShared
		default:
			queue = webQueueMessage
		}
	}
	// About a second, in whole slices, so a slice never wraps the shared ring.
	o := &webOut{src: d.Mixer, shared: shared, cap: webSlice * int(math.Ceil(float64(rt.AudioRate)/webSlice))}
	// Between one slice and the ring less one: a target the ring can't hold
	// would never be reached, and the pump would spin.
	o.queue = min(max(int64(queue)*rt.AudioRate/int64(time.Second), webSlice), int64(o.cap-webSlice))

	spec := map[string]any{"sampleRate": rt.AudioRate, "cap": o.cap}
	if shared {
		sab := g.Get("SharedArrayBuffer")
		samples, ctrl := sab.New(o.cap*8), sab.New(16)
		o.ctrl = g.Get("Int32Array").New(ctrl)
		o.atomics = g.Get("Atomics")
		data := g.Get("Uint8Array").New(samples)
		o.views = make([]js.Value, o.cap/webSlice)
		for i := range o.views {
			o.views[i] = data.Call("subarray", i*webSlice*8, (i+1)*webSlice*8)
		}
		spec["data"], spec["ctrl"] = samples, ctrl
	}
	// In a worker this is a relay to the page, which reports a failure later
	// through the state cell; on the page it fails here.
	h := open.Invoke(spec)
	if !worker && (h.IsNull() || h.IsUndefined()) {
		return nil, errors.New("audio: the page could not open audio (no AudioContext for this rate, or no AudioWorklet: it needs HTTPS or localhost)")
	}
	if !shared {
		o.handle = h
	}
	go o.pump()
	return o, nil
}

// Lost reports that the page can't play, which it may only learn after
// start returned (in a worker, or when the worklet fails to load).
func (o *webOut) Lost() bool { return o.lost.Load() }

func (o *webOut) state() int {
	if o.shared {
		return o.atomics.Call("load", o.ctrl, 3).Int()
	}
	return o.handle.Get("state").Int()
}

// queued counts frames handed over and not yet played. On the message path
// that includes slices still in flight to the worklet; counting only the
// worklet's own fill misses them and overfills its ring.
func (o *webOut) queued() int64 {
	if o.shared {
		w, r := o.atomics.Call("load", o.ctrl, 0).Int(), o.atomics.Call("load", o.ctrl, 1).Int()
		return int64((w - r + o.cap) % o.cap)
	}
	return o.sent - int64(o.handle.Get("played").Int())
}

// fill renders slices until the queue is full or the budget runs out.
func (o *webOut) fill() {
	t0 := time.Now()
	for o.queued() < o.queue && time.Since(t0) < webBudget {
		o.push()
	}
}

// pump keeps the queue full while the page plays it. Sleeping hands the
// thread back to the page (or the worker's VM).
//
// While the page doesn't play (before its first gesture, or suspended) the
// mixer's clock would stand still, since it only advances as it renders, and
// every sound triggered meanwhile would start at the same instant and burst
// out together on resume. So the pump renders and discards in step with the
// wall clock instead, as a native device keeps playing: a sound nobody could
// hear is gone, not saved up. Nothing is queued before the first gesture, so
// the first sound heard is current; after a suspension the queue already
// handed over (at most Queue) still plays on resume.
//
// The page reports its state rather than the pump inferring it from the
// worklet's position: a busy VM on the page's thread delays the position
// reports, which would look like a stall and discard audio still due to play.
func (o *webOut) pump() {
	last := time.Now()
	for {
		now := time.Now()
		switch o.state() {
		case pageRunning:
			o.fill()
		case pageFailed:
			o.lost.Store(true)
			return
		default:
			o.discard(int(int64(now.Sub(last)) * rt.AudioRate / int64(time.Second)))
		}
		last = now
		if o.state() == pageRunning && o.queued() < o.queue {
			time.Sleep(time.Millisecond) // behind: yield, come straight back
		} else {
			time.Sleep(8 * time.Millisecond) // about a frame
		}
	}
}

// discard renders frames of the mix and drops them, at most one ring's worth
// (about a second) and no longer than webBudget per call, as fill renders. A
// background tab can run the pump once a minute, and rendering all of that in
// one go stalls the VM's thread; the time past either bound is dropped, so a
// sound still playing resumes where the mixer is rather than where the wall
// clock is.
func (o *webOut) discard(frames int) {
	if o.raw == nil {
		o.raw = make([]byte, webSlice*4)
	}
	frames = min(frames, o.cap)
	t0 := time.Now()
	for ; frames > 0 && time.Since(t0) < webBudget; frames -= webSlice {
		n := min(frames, webSlice)
		io.ReadFull(o.src, o.raw[:n*4])
	}
}

// push renders one slice and hands it over. The mixer never ends or errors,
// so neither does this.
func (o *webOut) push() {
	if o.raw == nil {
		o.raw = make([]byte, webSlice*4)
	}
	if o.f32 == nil {
		o.f32 = make([]byte, webSlice*8)
	}
	io.ReadFull(o.src, o.raw)
	for i := 0; i < webSlice*2; i++ {
		s := float32(int16(binary.LittleEndian.Uint16(o.raw[2*i:]))) / 32768
		binary.LittleEndian.PutUint32(o.f32[4*i:], math.Float32bits(s))
	}
	if o.shared {
		w := o.atomics.Call("load", o.ctrl, 0).Int()
		js.CopyBytesToJS(o.views[w/webSlice], o.f32)
		o.atomics.Call("store", o.ctrl, 0, (w+webSlice)%o.cap)
		return
	}
	// A new Float32Array each time: transferring it detaches its buffer.
	a := js.Global().Get("Float32Array").New(webSlice * 2)
	js.CopyBytesToJS(js.Global().Get("Uint8Array").New(a.Get("buffer")), o.f32)
	o.handle.Call("push", a)
	o.sent += webSlice
}
