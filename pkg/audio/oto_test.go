//go:build audio_oto && !js

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"testing"
	"time"
)

// These use a fresh Device rather than Default, and never play, so no test
// opens a sound device.

func TestConfigureRejectsABadBudgetUpFront(t *testing.T) {
	d := &Device{Mixer: NewMixer()}
	if err := d.Configure(Options{DeviceFrames: 1000}); err == nil {
		t.Fatal("1000 frames accepted")
	}
	if d.opt != (Options{}) {
		t.Fatalf("a refused budget was kept: %+v", d.opt)
	}
}

func TestConfigureBeforeOpenIsKept(t *testing.T) {
	d := &Device{Mixer: NewMixer()}
	opt := Options{Queue: 60 * time.Millisecond, DeviceFrames: 2048}
	if err := d.Configure(opt); err != nil {
		t.Fatal(err)
	}
	if d.opt != opt {
		t.Fatalf("got %+v", d.opt)
	}
}

func TestConfigureAfterOpenChangesNothing(t *testing.T) {
	d := &Device{Mixer: NewMixer(), opened: true}
	if err := d.Configure(Options{DeviceFrames: 2048}); err == nil {
		t.Fatal("configured an open device")
	}
	if d.opt != (Options{}) {
		t.Fatalf("got %+v", d.opt)
	}
}
