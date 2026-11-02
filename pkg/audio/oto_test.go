//go:build audio_oto && !js

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// These use a fresh Device rather than Default, and never play, so no test
// opens a sound device.

// clearEnv keeps a developer's exported tuning values out of the tests.
func clearEnv(t *testing.T) {
	t.Setenv(envQueue, "")
	t.Setenv(envDeviceFrames, "")
}

func TestConfigureRejectsABadBudgetUpFront(t *testing.T) {
	clearEnv(t)
	d := &Device{Mixer: NewMixer()}
	if err := d.Configure(Options{DeviceFrames: 1000}); err == nil {
		t.Fatal("1000 frames accepted")
	}
	if d.opt != (Options{}) {
		t.Fatalf("a refused budget was kept: %+v", d.opt)
	}
}

func TestConfigureBeforeOpenIsKept(t *testing.T) {
	clearEnv(t)
	d := &Device{Mixer: NewMixer()}
	opt := Options{Queue: 60 * time.Millisecond, DeviceFrames: 2048}
	if err := d.Configure(opt); err != nil {
		t.Fatal(err)
	}
	if d.opt != opt {
		t.Fatalf("got %+v", d.opt)
	}
}

func TestConfigureIgnoresTheEnvironment(t *testing.T) {
	// 30 ms is valid at the default 1024 frames but not at 2048. The override
	// is checked at open, so the host's own valid budget is kept.
	t.Setenv(envQueue, "")
	t.Setenv(envDeviceFrames, "2048")
	d := &Device{Mixer: NewMixer()}
	opt := Options{Queue: 30 * time.Millisecond}
	if err := d.Configure(opt); err != nil {
		t.Fatal(err)
	}
	if d.opt != opt {
		t.Fatalf("got %+v", d.opt)
	}
}

func TestConfigureAfterOpenChangesNothing(t *testing.T) {
	clearEnv(t)
	d := &Device{Mixer: NewMixer(), opened: true}
	if err := d.Configure(Options{DeviceFrames: 2048}); err == nil {
		t.Fatal("configured an open device")
	}
	if d.opt != (Options{}) {
		t.Fatalf("got %+v", d.opt)
	}
}

// The bad override fails in resolve, before oto is touched, so Enabled here
// opens no sound device either.
func TestABadOverrideFailsTheOpenAndSaysWhy(t *testing.T) {
	clearEnv(t)
	t.Setenv(envQueue, "10ms")
	d := &Device{Mixer: NewMixer()}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr := os.Stderr
	os.Stderr = w
	enabled := d.Enabled()
	os.Stderr = stderr
	w.Close()
	said, _ := io.ReadAll(r)

	if enabled {
		t.Fatal("audio enabled under a bad override")
	}
	if !strings.Contains(string(said), "shorter than") {
		t.Fatalf("stderr %q does not give the reason", said)
	}
	// The failed open is final: a later Configure is refused.
	if err := d.Configure(Options{}); err == nil {
		t.Fatal("configured a device whose open already failed")
	}
}
