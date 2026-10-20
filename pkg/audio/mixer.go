/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

// Package audio is a sound backend for the audio namespace (let-go #255): a
// mixer that implements rt.Audio, and, built with -tags audio_oto, an oto
// device it plays through, bound at the *audio* root when the package is
// imported. Without the tag the package is the mixer alone, which renders
// headlessly for tests and offline mixing.
package audio

import (
	"encoding/binary"
	"io"
	"math"
	"sync"
	"time"

	"github.com/nooga/let-go/pkg/rt"
)

// edgeFrames is the fade a sound starts and stops with, about 5 ms, so a tone
// or a stop never clicks.
const edgeFrames = 220

// Mixer sums every playing sound into one stream of interleaved stereo signed
// 16-bit little-endian frames at rt.AudioRate, read through Read. It is safe
// for one reader and any number of goroutines calling the rt.Audio methods.
type Mixer struct {
	mu     sync.Mutex
	voices []*voice
	last   rt.AudioHandle
	mix    []float64
}

// NewMixer returns a mixer with nothing playing; Read yields silence.
func NewMixer() *Mixer { return &Mixer{} }

type voice struct {
	h rt.AudioHandle

	// one of: pcm (a sample), tone (left > 0), or src (a stream)
	pcm      []int16
	pos      int
	loopFrom int
	tone     rt.Tone
	phase    float64
	left     int
	src      io.Reader
	sbuf     []byte

	gain, target, step float64
	stopping, done     bool
}

func (m *Mixer) add(v *voice, gain float64) rt.AudioHandle {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.last++
	v.h = m.last
	v.target = gain
	v.step = gain / edgeFrames
	m.voices = append(m.voices, v)
	return v.h
}

// Play starts a sample or a tone.
func (m *Mixer) Play(s rt.Sound) (rt.AudioHandle, error) {
	v := &voice{pcm: s.PCM, loopFrom: s.LoopFrom}
	if s.Tone.Freq > 0 {
		v.pcm, v.tone = nil, s.Tone
		v.left = int(s.Tone.Dur.Seconds() * rt.AudioRate)
	}
	return m.add(v, s.Gain), nil
}

// Stream mixes src until it errors or the handle is stopped.
func (m *Mixer) Stream(src io.Reader) (rt.AudioHandle, error) {
	return m.add(&voice{src: src}, 1), nil
}

// Stop fades the sound out over edgeFrames and drops it.
func (m *Mixer) Stop(h rt.AudioHandle) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v := m.find(h); v != nil {
		v.stopping = true
		v.target, v.step = 0, max(v.gain, v.target)/edgeFrames
	}
}

// Gain moves the sound's level to g over fade, or over edgeFrames when fade
// is shorter, so a level change never clicks.
func (m *Mixer) Gain(h rt.AudioHandle, g float64, fade time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v := m.find(h); v != nil && !v.stopping {
		frames := max(edgeFrames, fade.Seconds()*rt.AudioRate)
		v.target, v.step = g, math.Abs(g-v.gain)/frames
	}
}

// Enabled is false for a bare mixer: nothing hears it.
func (m *Mixer) Enabled() bool { return false }

func (m *Mixer) find(h rt.AudioHandle) *voice {
	for _, v := range m.voices {
		if v.h == h {
			return v
		}
	}
	return nil
}

// Read mixes the next len(p)/4 frames.
func (m *Mixer) Read(p []byte) (int, error) {
	n := len(p) / 4
	if n == 0 {
		return 0, io.ErrShortBuffer
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if cap(m.mix) < 2*n {
		m.mix = make([]float64, 2*n)
	}
	mix := m.mix[:2*n]
	clear(mix)
	live := m.voices[:0]
	for _, v := range m.voices {
		v.render(mix)
		if !v.done {
			live = append(live, v)
		}
	}
	clear(m.voices[len(live):])
	m.voices = live
	for i, x := range mix {
		binary.LittleEndian.PutUint16(p[2*i:], uint16(int16(max(-32768, min(32767, x)))))
	}
	return 4 * n, nil
}

// render adds the voice's next frames into mix, interleaved stereo.
func (v *voice) render(mix []float64) {
	frames := len(mix) / 2
	if v.src != nil && cap(v.sbuf) < 4*frames {
		v.sbuf = make([]byte, 4*frames)
	}
	if v.src != nil {
		got, err := io.ReadFull(v.src, v.sbuf[:4*frames])
		if err != nil {
			v.done = true
		}
		frames = got / 4
	}
	for i := 0; i < frames; i++ {
		var l, r float64
		switch {
		case v.src != nil:
			l = float64(int16(binary.LittleEndian.Uint16(v.sbuf[4*i:])))
			r = float64(int16(binary.LittleEndian.Uint16(v.sbuf[4*i+2:])))
		case v.left > 0:
			l = 0.5 * 32767 * wave(v.tone.Wave, v.phase)
			r = l
			v.phase += v.tone.Freq / rt.AudioRate
			v.phase -= math.Floor(v.phase)
			v.left--
			if v.left == edgeFrames {
				v.target, v.step = 0, v.gain/edgeFrames
			}
			if v.left == 0 {
				v.done = true
			}
		default:
			if v.pos >= len(v.pcm) {
				if v.loopFrom < 0 {
					v.done = true
					return
				}
				v.pos = 2 * v.loopFrom
			}
			l, r = float64(v.pcm[v.pos]), float64(v.pcm[v.pos+1])
			v.pos += 2
		}
		if v.gain != v.target {
			if v.gain < v.target {
				v.gain = min(v.target, v.gain+v.step)
			} else {
				v.gain = max(v.target, v.gain-v.step)
			}
		}
		mix[2*i] += l * v.gain
		mix[2*i+1] += r * v.gain
		if v.stopping && v.gain == 0 {
			v.done = true
		}
		if v.done {
			return
		}
	}
}

func wave(name string, phase float64) float64 {
	switch name {
	case "square":
		if phase < 0.5 {
			return 1
		}
		return -1
	case "saw":
		return 2*phase - 1
	case "triangle":
		return 1 - 4*math.Abs(phase-0.5)
	}
	return math.Sin(2 * math.Pi * phase)
}

var _ rt.Audio = (*Mixer)(nil)
