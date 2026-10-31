/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * The audio namespace — a handle-based host capability (let-go #255, Audio).
 * The guest names a sound; the host plays it and hands back a handle the guest
 * can stop or fade later. Fire-and-forget: no playback state (position,
 * level, whether a sound has ended) flows back into the guest. A guest can
 * still observe its handles, audio/enabled? and refused calls, so it replays
 * identically with or without sound only if it does not branch on those.
 *
 * Bound at the *audio* root by the host, the same binding discipline as
 * *surface*: a backend that registers itself (pkg/audio, built with
 * -tags audio_oto), an Audio via api.WithAudio for embedders and tests, and
 * nopAudio otherwise, so every audio/ call is a silent no-op on a host with no
 * sound.
 *
 * Stream is Go-only: a native module that synthesizes continuously (a tracker,
 * a generative engine) joins the host's mixer as a stream instead of opening a
 * device of its own, and gets an ordinary handle back. lg code never sees it.
 */

package rt

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

// AudioRate is the rate of every PCM buffer and stream on the seam:
// interleaved stereo signed 16-bit frames at 44.1 kHz.
const AudioRate = 44100

// AudioHandle names one playing sound. Zero is never issued, so a nop
// backend's zero handle is safe to stop or fade.
type AudioHandle int64

// Sound is what audio/tone and audio/sample ask the host to play.
type Sound struct {
	// PCM is interleaved stereo at AudioRate. Empty for a tone.
	PCM []int16
	// LoopFrom is the frame playback returns to at the end, for an intro
	// that plays once before a loop; -1 plays once.
	LoopFrom int
	// Tone, when Freq is set, synthesizes a tone instead of playing PCM.
	Tone Tone
	Gain float64
}

// Tone is a synthesized note: Wave is "sine", "square", "saw" or "triangle".
type Tone struct {
	Freq float64
	Dur  time.Duration
	Wave string
}

// Audio is the host seam for the audio namespace.
type Audio interface {
	Play(Sound) (AudioHandle, error)
	// Stream mixes src, interleaved stereo signed 16-bit little-endian PCM at
	// AudioRate, until it returns an error or the handle is stopped. src is
	// read on the host's audio goroutine in small slices, so it must not
	// block.
	Stream(src io.Reader) (AudioHandle, error)
	Stop(AudioHandle)
	// Gain moves a sound's level to g over fade. A backend may smooth a zero
	// or very short fade over a few milliseconds so the change does not
	// click; pkg/audio's mixer does.
	Gain(h AudioHandle, g float64, fade time.Duration)
	// Enabled reports whether sound can be heard: a backend is built in and
	// it opened a device.
	Enabled() bool
}

type nopAudio struct{}

func (nopAudio) Play(Sound) (AudioHandle, error)          { return 0, nil }
func (nopAudio) Stream(io.Reader) (AudioHandle, error)    { return 0, nil }
func (nopAudio) Stop(AudioHandle)                         {}
func (nopAudio) Gain(AudioHandle, float64, time.Duration) {}
func (nopAudio) Enabled() bool                            { return false }

// SetAudioRoot binds a as the process-wide audio backend: the root of
// *audio*, which api.WithAudio and (binding [*audio* ...]) still override. A
// backend package calls it from init.
func SetAudioRoot(a Audio) {
	if v := LookupCoreVar("*audio*"); v != nil {
		v.SetRoot(vm.NewBoxed(a))
	}
}

// AudioVia resolves the current *audio* binding, nopAudio when unbound. A
// ctx-aware native (vm.NewCtxNativeFn) uses it to reach Stream.
func AudioVia(ec *vm.ExecContext) Audio {
	if v := LookupCoreVar("*audio*"); v != nil {
		if b, ok := ec.Deref(v).(*vm.Boxed); ok {
			if a, ok := b.Unbox().(Audio); ok {
				return a
			}
		}
	}
	return nopAudio{}
}

func init() { RegisterInstaller(installAudioNS) }

func installAudioNS() {
	ns := vm.NewNamespace("audio")
	def := func(name string, fn func(Audio, []vm.Value) (vm.Value, error)) {
		ns.Def(name, vm.NewCtxNativeFn(name, func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
			return fn(AudioVia(ec), vs)
		}))
	}

	// (audio/tone {:freq 440 :ms 120 :wave :square :gain 0.3}) -> handle
	def("tone", func(a Audio, vs []vm.Value) (vm.Value, error) {
		m, err := audioOpts("audio/tone", vs, 1)
		if err != nil {
			return vm.NIL, err
		}
		g, err := startGain("audio/tone", m)
		if err != nil {
			return vm.NIL, err
		}
		freq := optFloat(m, "freq", 0)
		dur, ok := msDuration(optFloat(m, "ms", 0))
		if !(freq > 0 && !math.IsInf(freq, 1)) || !ok || dur <= 0 {
			return vm.NIL, fmt.Errorf("audio/tone needs a finite positive :freq and :ms")
		}
		s := Sound{LoopFrom: -1, Gain: g, Tone: Tone{Freq: freq, Dur: dur, Wave: optName(m, "wave", "sine")}}
		switch s.Tone.Wave {
		case "sine", "square", "saw", "triangle":
		default:
			return vm.NIL, fmt.Errorf("audio/tone: unknown :wave %q", s.Tone.Wave)
		}
		return audioHandle(a.Play(s))
	})

	// (audio/sample pcm) or (audio/sample pcm {:loop-from frames :gain g}) -> handle.
	// pcm is a byte-array of interleaved stereo signed 16-bit LE at 44.1 kHz.
	def("sample", func(a Audio, vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("audio/sample expects 1 or 2 args (pcm [opts]), got %d", len(vs))
		}
		ta, ok := vs[0].(*vm.TypedArray)
		if !ok || ta.Kind() != vm.ArrayByte {
			return vm.NIL, fmt.Errorf("audio/sample expects a byte-array of 16-bit stereo PCM")
		}
		b := ta.Unbox().([]byte)
		if len(b) == 0 || len(b)%4 != 0 {
			return vm.NIL, fmt.Errorf("audio/sample: %d bytes is not a whole number of stereo frames", len(b))
		}
		m, err := audioOpts("audio/sample", vs[1:], 0)
		if err != nil {
			return vm.NIL, err
		}
		pcm := make([]int16, len(b)/2)
		for i := range pcm {
			pcm[i] = int16(uint16(b[2*i]) | uint16(b[2*i+1])<<8)
		}
		g, err := startGain("audio/sample", m)
		if err != nil {
			return vm.NIL, err
		}
		s := Sound{PCM: pcm, LoopFrom: -1, Gain: g}
		if v := optFloat(m, "loop-from", -1); v >= 0 {
			if v >= float64(len(pcm)/2) {
				return vm.NIL, fmt.Errorf("audio/sample: :loop-from %v is past the last frame", v)
			}
			s.LoopFrom = int(v)
		}
		return audioHandle(a.Play(s))
	})

	// (audio/stop h) -> nil. A nil h, the handle a host with no sound
	// returns, is a no-op, so a game needs no guard to run silent.
	def("stop", func(a Audio, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("audio/stop expects 1 arg (handle), got %d", len(vs))
		}
		if vs[0] == vm.NIL {
			return vm.NIL, nil
		}
		h, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("audio/stop expects a handle")
		}
		a.Stop(AudioHandle(h))
		return vm.NIL, nil
	})

	// (audio/gain h g) or (audio/gain h g fade-ms) -> nil
	def("gain", func(a Audio, vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("audio/gain expects 2 or 3 args (handle gain [fade-ms]), got %d", len(vs))
		}
		h, ok := vs[0].(vm.Int)
		ok = ok || vs[0] == vm.NIL
		g, ok2 := vm.ToFloat(vs[1])
		if !ok || !ok2 || !validGain(g) {
			return vm.NIL, fmt.Errorf("audio/gain expects a handle and a gain of at least 0")
		}
		var fade time.Duration
		if len(vs) == 3 {
			ms, ok := vm.ToFloat(vs[2])
			if ok {
				fade, ok = msDuration(ms)
			}
			if !ok {
				return vm.NIL, fmt.Errorf("audio/gain: fade-ms must be finite and at least 0")
			}
		}
		if vs[0] != vm.NIL {
			a.Gain(AudioHandle(h), g, fade)
		}
		return vm.NIL, nil
	})

	// (audio/enabled?) -> bool
	def("enabled?", func(a Audio, vs []vm.Value) (vm.Value, error) {
		return vm.Boolean(a.Enabled()), nil
	})

	RegisterNS(ns)
}

// audioHandle turns a backend's result into the guest's: an Int handle, or
// nil from a host with no sound.
func audioHandle(h AudioHandle, err error) (vm.Value, error) {
	if err != nil {
		return vm.NIL, err
	}
	if h == 0 {
		return vm.NIL, nil
	}
	return vm.Int(h), nil
}

// audioOpts reads an optional options map; need is how many args are required.
func audioOpts(name string, vs []vm.Value, need int) (*vm.PersistentMap, error) {
	if len(vs) < need || len(vs) > 1 {
		return nil, fmt.Errorf("%s: expected an options map", name)
	}
	if len(vs) == 0 || vs[0] == vm.NIL {
		return nil, nil
	}
	m, ok := vs[0].(*vm.PersistentMap)
	if !ok {
		return nil, fmt.Errorf("%s: options must be a map", name)
	}
	return m, nil
}

// validGain is a level a mixer can ramp to: finite and at least 0.
func validGain(g float64) bool { return g >= 0 && !math.IsInf(g, 1) }

// maxMs is the longest time.Duration, in ms.
const maxMs = float64(math.MaxInt64 / int64(time.Millisecond))

// msDuration converts ms to a Duration, refusing NaN, a negative, ##Inf or
// anything a Duration cannot hold rather than converting it arch-dependently.
func msDuration(ms float64) (time.Duration, bool) {
	if !(ms >= 0 && ms <= maxMs) {
		return 0, false
	}
	return time.Duration(ms * float64(time.Millisecond)), true
}

// startGain reads a sound's starting :gain, 1 by default.
func startGain(name string, m *vm.PersistentMap) (float64, error) {
	g := optFloat(m, "gain", 1)
	if !validGain(g) {
		return 0, fmt.Errorf("%s: :gain must be finite and at least 0", name)
	}
	return g, nil
}

func optFloat(m *vm.PersistentMap, k string, def float64) float64 {
	if m == nil {
		return def
	}
	if f, ok := vm.ToFloat(m.ValueAt(vm.Keyword(k))); ok {
		return f
	}
	return def
}

func optName(m *vm.PersistentMap, k, def string) string {
	if m == nil {
		return def
	}
	switch v := m.ValueAt(vm.Keyword(k)).(type) {
	case vm.Keyword:
		return string(v)
	case vm.String:
		return string(v)
	}
	return def
}
