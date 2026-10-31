//go:build audio_oto && !js

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/ebitengine/oto/v3"

	"github.com/nooga/let-go/pkg/rt"
)

// Device is the process's sound output: a Mixer played through oto. The
// device opens on the first sound, so a program that never plays one never
// touches the audio system.
type Device struct {
	*Mixer
	mu     sync.Mutex // guards opt and opened, which Configure reads
	opt    Options
	opened bool
	once   sync.Once
	err    error
	player *oto.Player // kept in a field: oto closes an unreachable player mid-playback (oto#293)
}

var device = &Device{Mixer: NewMixer()}

// Default returns the process's Device, which the package also binds at the
// *audio* root. An embedder calls Stream on it to join a synth to the mix.
func Default() *Device { return device }

func init() { rt.SetAudioRoot(device) }

// Configure sets the device's latency budget. It must come before the first
// sound, since oto fixes the device buffer when it opens; after that it
// returns an error and changes nothing. LETGO_AUDIO_QUEUE and
// LETGO_AUDIO_DEVICE_FRAMES still win over opt (see Options).
func (d *Device) Configure(opt Options) error {
	if _, err := resolve(opt, os.Getenv); err != nil {
		return err
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.opened {
		return errors.New("audio: the device is already open; configure it before the first sound")
	}
	d.opt = opt
	return nil
}

func (d *Device) open() error {
	d.once.Do(func() {
		d.mu.Lock()
		d.opened = true
		opt := d.opt
		d.mu.Unlock()
		l, err := resolve(opt, os.Getenv)
		if err != nil {
			d.err = err
			return
		}
		ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
			SampleRate:   rt.AudioRate,
			ChannelCount: 2,
			Format:       oto.FormatSignedInt16LE,
			BufferSize:   deviceBuffer(l.frames),
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
		d.player.SetBufferSize(playerBytes(l.queue))
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
