//go:build (audio_oto && !js) || (js && wasm)

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/nooga/let-go/pkg/rt"
)

// Device is the process's sound output: a Mixer played through a backend,
// oto natively (oto.go) or an AudioWorklet in a browser bundle (web.go). The
// device opens on the first sound, so a program that never plays one never
// touches the audio system.
type Device struct {
	*Mixer
	mu      sync.Mutex // guards opt and opened, which Configure reads
	opt     Options
	opened  bool
	once    sync.Once
	err     error
	backend any // what start returns, kept reachable for as long as the process plays
}

var device = &Device{Mixer: NewMixer()}

// Default returns the process's Device, which the package also binds at the
// *audio* root. An embedder calls Stream on it to join a synth to the mix.
func Default() *Device { return device }

func init() { rt.SetAudioRoot(device) }

// Configure sets the device's latency budget. It must come before the first
// sound, since the backend fixes its buffers when it opens; after that it
// returns an error and changes nothing. LETGO_AUDIO_QUEUE and
// LETGO_AUDIO_DEVICE_FRAMES still win over opt (see Options); they are
// checked when the device opens, so an error here is always about opt.
func (d *Device) Configure(opt Options) error {
	if _, err := resolve(opt, noEnv); err != nil {
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

// noEnv lets Configure check the host's Options without the overrides.
func noEnv(string) string { return "" }

func (d *Device) open() error {
	d.once.Do(func() {
		d.mu.Lock()
		d.opened = true
		opt := d.opt
		d.mu.Unlock()
		l, err := resolve(opt, os.Getenv)
		if err != nil {
			// audio/enabled? can only answer false, and a game that checks it
			// turns sound off quietly, so a tuning run would hear silence and
			// never learn why. Say it once here.
			fmt.Fprintln(os.Stderr, err)
			d.err = err
			return
		}
		// A backend whose best queue differs from the native default (a
		// browser's does) needs to know whether this one was asked for.
		asked := opt.Queue != 0 || os.Getenv(envQueue) != ""
		d.backend, d.err = d.start(l, asked)
		if d.err != nil {
			fmt.Fprintln(os.Stderr, d.err)
		}
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
// worked and the backend has not since found it can't play (a browser page
// may only learn that after the device opened).
func (d *Device) Enabled() bool {
	if d.open() != nil {
		return false
	}
	l, ok := d.backend.(interface{ Lost() bool })
	return !ok || !l.Lost()
}

var _ rt.Audio = (*Device)(nil)
