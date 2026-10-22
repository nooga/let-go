//go:build audio_oto && !js

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"io"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"

	"github.com/nooga/let-go/pkg/rt"
)

// Device is the process's sound output: a Mixer played through oto. The
// device opens on the first sound, so a program that never plays one never
// touches the audio system.
type Device struct {
	*Mixer
	once   sync.Once
	err    error
	player *oto.Player // kept in a field: oto closes an unreachable player mid-playback (oto#293)
}

// The latency budget measured for undertone's native player (its
// nativeaudio/latency.go): a 512-frame device buffer, which oto takes as a
// duration and truncates back, so it is rounded up to survive the round trip,
// and 40 ms of player read-ahead in place of oto's 500. Tap to sound measured
// about 130 ms, against about 820 ms with oto's defaults.
const deviceBuffer = time.Duration((512*8*int64(time.Second) + rt.AudioRate*4 - 1) / (rt.AudioRate * 4))

var device = &Device{Mixer: NewMixer()}

// Default returns the process's Device, which the package also binds at the
// *audio* root. An embedder calls Stream on it to join a synth to the mix.
func Default() *Device { return device }

func init() { rt.SetAudioRoot(device) }

func (d *Device) open() error {
	d.once.Do(func() {
		ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
			SampleRate:   rt.AudioRate,
			ChannelCount: 2,
			Format:       oto.FormatSignedInt16LE,
			BufferSize:   deviceBuffer,
		})
		if err != nil {
			d.err = err
			return
		}
		<-ready
		// oto opens the device asynchronously and reports a failure here.
		if err := ctx.Err(); err != nil {
			d.err = err
			return
		}
		d.player = ctx.NewPlayer(d.Mixer)
		d.player.SetBufferSize(rt.AudioRate * 4 * 40 / 1000)
		d.player.Play()
	})
	return d.err
}

func (d *Device) Play(s rt.Sound) (rt.AudioHandle, error) {
	if err := d.open(); err != nil {
		return 0, err
	}
	return d.Mixer.Play(s)
}

func (d *Device) Stream(src io.Reader) (rt.AudioHandle, error) {
	if err := d.open(); err != nil {
		return 0, err
	}
	return d.Mixer.Stream(src)
}

// Enabled opens the device if it is not open yet, and reports whether that
// worked.
func (d *Device) Enabled() bool { return d.open() == nil }

var _ rt.Audio = (*Device)(nil)
