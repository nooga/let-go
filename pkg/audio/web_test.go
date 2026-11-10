//go:build js && wasm

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"strings"
	"syscall/js"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/rt"
)

// setGlobal sets a JS global for one test, restoring it afterwards.
func setGlobal(t *testing.T, name string, v any) {
	t.Helper()
	g := js.Global()
	old := g.Get(name)
	g.Set(name, v)
	t.Cleanup(func() {
		if old.IsUndefined() {
			g.Delete(name)
		} else {
			g.Set(name, old)
		}
	})
}

// fakeHost installs _lgAudioOpen and returns what the page was handed.
// handle is what the fake returns: an object for the message path, or {}
// for the shared ring, which needs none.
func fakeHost(t *testing.T, handle js.Value) *js.Value {
	t.Helper()
	var spec js.Value
	open := js.FuncOf(func(_ js.Value, a []js.Value) any {
		spec = a[0]
		return handle
	})
	t.Cleanup(open.Release)
	setGlobal(t, "_lgAudioOpen", open)
	return &spec
}

// asWorker makes the page look like a cross-origin isolated worker: no
// document, SharedArrayBuffer available (node has it).
func asWorker(t *testing.T) {
	setGlobal(t, "crossOriginIsolated", true)
	if !js.Global().Get("document").IsUndefined() {
		setGlobal(t, "document", js.Undefined())
	}
}

// asPage makes the page look like a main thread without isolation.
func asPage(t *testing.T) {
	setGlobal(t, "crossOriginIsolated", false)
	setGlobal(t, "document", js.Global().Get("Object").New())
}

func frames(d time.Duration) int64 { return int64(d) * rt.AudioRate / int64(time.Second) }

func eventually(t *testing.T, what string, ok func() bool) {
	t.Helper()
	for deadline := time.Now().Add(2 * time.Second); time.Now().Before(deadline); {
		if ok() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal(what)
}

func TestWebStartNeedsTheAudioHost(t *testing.T) {
	asPage(t)
	setGlobal(t, "_lgAudioOpen", js.Undefined())
	d := &Device{Mixer: NewMixer()}
	_, err := d.start(latency{queue: defaultQueue, frames: defaultDeviceFrames}, false)
	if err == nil || !strings.Contains(err.Error(), "-w-audio") {
		t.Fatalf("err %v, want one naming -w-audio", err)
	}
}

func voices(m *Mixer) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.voices)
}

func tone(t *testing.T, d *Device, dur time.Duration) {
	t.Helper()
	if _, err := d.Mixer.Play(rt.Sound{Tone: rt.Tone{Freq: 440, Dur: dur, Wave: "sine"}, Gain: 0.5}); err != nil {
		t.Fatal(err)
	}
}

func TestWebSharedRingInAWorker(t *testing.T) {
	asWorker(t)
	spec := fakeHost(t, js.Global().Get("Object").New())
	d := &Device{Mixer: NewMixer()}
	b, err := d.start(latency{queue: defaultQueue, frames: defaultDeviceFrames}, false)
	if err != nil {
		t.Fatal(err)
	}
	o := b.(*webOut)
	if !o.shared || spec.Get("data").IsUndefined() || spec.Get("ctrl").IsUndefined() {
		t.Fatalf("worker took the message path; spec %v", spec.Get("data"))
	}
	// Nothing is queued before the page plays: a sound rendered now would be
	// heard stale at the first gesture.
	time.Sleep(30 * time.Millisecond)
	if q := o.queued(); q != 0 {
		t.Fatalf("%d frames queued before the page ran", q)
	}

	// The page starts; the pump fills the worker's queue with what plays now.
	tone(t, d, time.Second)
	ctrl := js.Global().Get("Int32Array").New(spec.Get("ctrl"))
	atomics := js.Global().Get("Atomics")
	atomics.Call("store", ctrl, 3, pageRunning)
	eventually(t, "the pump never filled the running ring", func() bool { return o.queued() >= frames(webQueueWorker) })
	samples := js.Global().Get("Float32Array").New(spec.Get("data"))
	heard := false
	for i := 0; i < 2*int(o.queued()); i++ {
		if samples.Index(i).Float() != 0 {
			heard = true
			break
		}
	}
	if !heard {
		t.Fatal("the ring holds only silence; the tone was not rendered into it")
	}

	// Play out what is queued, as the worklet does, and the pump refills.
	atomics.Call("store", ctrl, 1, ctrl.Index(0).Int())
	eventually(t, "the pump never refilled a drained ring", func() bool { return o.queued() >= o.queue })
}

// messageHost is a fake page handle for the message path.
func messageHost(t *testing.T) (handle js.Value, pushed *int) {
	t.Helper()
	pushed = new(int)
	handle = js.Global().Get("Object").New()
	handle.Set("played", 0)
	handle.Set("state", pageOpening)
	push := js.FuncOf(func(_ js.Value, a []js.Value) any {
		*pushed += a[0].Get("length").Int() / 2
		return nil
	})
	t.Cleanup(push.Release)
	handle.Set("push", push)
	return handle, pushed
}

func TestWebMessagePathOnAPage(t *testing.T) {
	asPage(t)
	handle, pushed := messageHost(t)
	spec := fakeHost(t, handle)
	d := &Device{Mixer: NewMixer()}
	b, err := d.start(latency{queue: defaultQueue, frames: defaultDeviceFrames}, false)
	if err != nil {
		t.Fatal(err)
	}
	o := b.(*webOut)
	if o.shared || !spec.Get("data").IsUndefined() {
		t.Fatal("an unisolated page was handed a shared ring")
	}
	handle.Set("state", pageRunning)
	eventually(t, "the pump never pushed the message queue", func() bool { return int64(*pushed) >= frames(webQueueMessage) })
	// The worklet reports it played everything; the pump tops up.
	before := *pushed
	handle.Set("played", before)
	eventually(t, "the pump never pushed after the worklet played", func() bool { return *pushed >= before+int(frames(webQueueMessage)) })
}

// Before the first gesture nothing drains the ring. A sound triggered then
// must still end on time, not wait at the mixer's stopped clock and burst out
// with every other sound once playback starts.
func TestWebStoppedPageKeepsTheMixersClock(t *testing.T) {
	asWorker(t)
	fakeHost(t, js.Global().Get("Object").New())
	d := &Device{Mixer: NewMixer()}
	if _, err := d.start(latency{queue: defaultQueue, frames: defaultDeviceFrames}, false); err != nil {
		t.Fatal(err)
	}
	tone(t, d, 100*time.Millisecond)
	eventually(t, "a 100 ms tone was still playing in the mixer before the page ran", func() bool { return voices(d.Mixer) == 0 })
}

// A running page whose position reports are late (a busy VM on its thread
// delays them) must not look stalled: the pump waits, discarding nothing that
// is still due to play.
func TestWebLateReportsDiscardNothing(t *testing.T) {
	asPage(t)
	handle, _ := messageHost(t)
	handle.Set("state", pageRunning)
	fakeHost(t, handle)
	d := &Device{Mixer: NewMixer()}
	if _, err := d.start(latency{queue: defaultQueue, frames: defaultDeviceFrames}, false); err != nil {
		t.Fatal(err)
	}
	tone(t, d, 300*time.Millisecond)
	time.Sleep(500 * time.Millisecond) // played never moves
	if voices(d.Mixer) == 0 {
		t.Fatal("the tone was discarded while the page ran with its reports late")
	}
}

// The page may only learn it can't play after the device opened (a worker's
// relay, a blocked worklet module); audio/enabled? must then turn false.
func TestWebLatePageFailureDisables(t *testing.T) {
	asWorker(t)
	spec := fakeHost(t, js.Global().Get("Object").New())
	d := &Device{Mixer: NewMixer()}
	if !d.Enabled() {
		t.Fatal("not enabled before the page reported anything")
	}
	ctrl := js.Global().Get("Int32Array").New(spec.Get("ctrl"))
	js.Global().Get("Atomics").Call("store", ctrl, 3, pageFailed)
	eventually(t, "audio/enabled? stayed true after the page failed", func() bool { return !d.Enabled() })
}

func TestWebAskedQueueWins(t *testing.T) {
	asWorker(t)
	fakeHost(t, js.Global().Get("Object").New())
	d := &Device{Mixer: NewMixer()}
	b, err := d.start(latency{queue: 200 * time.Millisecond, frames: defaultDeviceFrames}, true)
	if err != nil {
		t.Fatal(err)
	}
	if o := b.(*webOut); o.queue != frames(200*time.Millisecond) {
		t.Fatalf("queue %d frames, want the asked 200 ms (%d)", o.queue, frames(200*time.Millisecond))
	}
}

// countReader counts the bytes read from it and returns silence.
type countReader struct{ n int }

func (c *countReader) Read(p []byte) (int, error) {
	clear(p)
	c.n += len(p)
	return len(p), nil
}

// A long suspension (a background tab whose pump ran once a minute) renders at
// most one ring of discarded audio in one call, not the whole gap, and the
// discard path never needs push's float32 buffer.
func TestWebDiscardIsBounded(t *testing.T) {
	src := &countReader{}
	o := &webOut{src: src, cap: 4 * webSlice}
	o.discard(10*o.cap, func() time.Time { return time.Time{} })
	if got, want := src.n, o.cap*4; got != want {
		t.Fatalf("discard rendered %d bytes, want one ring (%d)", got, want)
	}
	if o.f32 != nil {
		t.Fatal("discard allocated the float32 buffer")
	}
}

func TestWebDiscardStopsAtRenderBudget(t *testing.T) {
	src := &countReader{}
	o := &webOut{src: src, cap: 4 * webSlice}
	o.discard(o.cap, func() time.Time {
		// Advance the clock to the deadline after the first slice.
		if src.n == 0 {
			return time.Time{}
		}
		return time.Time{}.Add(webBudget)
	})
	if got, want := src.n, webSlice*4; got != want {
		t.Fatalf("discard rendered %d bytes, want one slice (%d)", got, want)
	}
}
