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

func TestWebSharedRingInAWorker(t *testing.T) {
	asWorker(t)
	spec := fakeHost(t, js.Global().Get("Object").New())
	d := &Device{Mixer: NewMixer()}
	if _, err := d.Mixer.Play(rt.Sound{Tone: rt.Tone{Freq: 440, Dur: time.Second, Wave: "sine"}, Gain: 0.5}); err != nil {
		t.Fatal(err)
	}
	b, err := d.start(latency{queue: defaultQueue, frames: defaultDeviceFrames}, false)
	if err != nil {
		t.Fatal(err)
	}
	o := b.(*webOut)
	if !o.shared || spec.Get("data").IsUndefined() || spec.Get("ctrl").IsUndefined() {
		t.Fatalf("worker took the message path; spec %v", spec.Get("data"))
	}
	// Primed before the page saw the ring: at least the worker's queue.
	if q := o.queued(); q < frames(webQueueWorker) {
		t.Fatalf("handed over %d frames queued, want the %v worker queue (%d)", q, webQueueWorker, frames(webQueueWorker))
	}
	samples := js.Global().Get("Float32Array").New(spec.Get("data"))
	heard := false
	for i := 0; i < 2*int(o.queued()); i++ {
		if samples.Index(i).Float() != 0 {
			heard = true
			break
		}
	}
	if !heard {
		t.Fatal("the primed ring holds only silence; the tone was not rendered into it")
	}

	// Play out what is queued, as the worklet does, and the pump refills.
	ctrl := js.Global().Get("Int32Array").New(spec.Get("ctrl"))
	js.Global().Get("Atomics").Call("store", ctrl, 1, ctrl.Index(0).Int())
	eventually(t, "the pump never refilled a drained ring", func() bool { return o.queued() >= o.queue })
}

func TestWebMessagePathOnAPage(t *testing.T) {
	asPage(t)
	var pushed int
	handle := js.Global().Get("Object").New()
	handle.Set("played", 0)
	push := js.FuncOf(func(_ js.Value, a []js.Value) any {
		pushed += a[0].Get("length").Int() / 2
		return nil
	})
	t.Cleanup(push.Release)
	handle.Set("push", push)
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
	if int64(pushed) < frames(webQueueMessage) {
		t.Fatalf("pushed %d frames, want the %v message queue (%d)", pushed, webQueueMessage, frames(webQueueMessage))
	}
	// The worklet reports it played everything; the pump tops up.
	before := pushed
	handle.Set("played", pushed)
	eventually(t, "the pump never pushed after the worklet played", func() bool { return pushed >= before+int(frames(webQueueMessage)) })
}

// Before the first gesture nothing drains the ring. A sound triggered then
// must still end on time, not wait at the mixer's stopped clock and burst out
// with every other sound once playback starts.
func TestWebStalledConsumerKeepsTheMixersClock(t *testing.T) {
	asWorker(t)
	fakeHost(t, js.Global().Get("Object").New())
	d := &Device{Mixer: NewMixer()}
	if _, err := d.start(latency{queue: defaultQueue, frames: defaultDeviceFrames}, false); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Mixer.Play(rt.Sound{Tone: rt.Tone{Freq: 440, Dur: 100 * time.Millisecond, Wave: "sine"}, Gain: 0.5}); err != nil {
		t.Fatal(err)
	}
	eventually(t, "a 100 ms tone was still playing in the mixer with the consumer stalled", func() bool {
		d.Mixer.mu.Lock()
		defer d.Mixer.mu.Unlock()
		return len(d.Mixer.voices) == 0
	})
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
