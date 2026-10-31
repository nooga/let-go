//go:build audio_oto && !js

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"io"
	"runtime"
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

// The latency budget follows undertone's native player (its
// nativeaudio/latency.go): a 1024-frame device buffer, which oto takes as a
// duration and truncates back, so it is rounded up to survive the round trip,
// and a short player read-ahead in place of oto's 500 ms. Measured 2026-10-26
// on a MacBook Air (speakers into its own mic, trigger to onset), 40 ms of
// read-ahead gave 150-158 ms and 24 ms gave 119-128 ms; 24 ms over a
// 512-frame device crackled under a busy main thread. Only macOS has been
// heard at 24 ms, so other platforms keep 40 ms.
const deviceFrames = 1024

const deviceBuffer = time.Duration((deviceFrames*int64(time.Second) + rt.AudioRate - 1) / rt.AudioRate)

func readAheadMS() int {
	if runtime.GOOS == "darwin" { // GOOS is "ios" on iOS, which keeps 40 ms
		return 24
	}
	return 40
}

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
		d.player.SetBufferSize(rt.AudioRate * 4 * readAheadMS() / 1000)
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
