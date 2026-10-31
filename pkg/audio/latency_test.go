/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"strings"
	"testing"
	"time"
)

func env(kv map[string]string) func(string) string {
	return func(k string) string { return kv[k] }
}

func TestDefaultsAreTheSeamsMeasuredBudget(t *testing.T) {
	l, err := resolve(Options{}, env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if l.frames != 1024 || l.queue != defaultQueue() {
		t.Fatalf("got %d frames, %v queue; want 1024, %v", l.frames, l.queue, defaultQueue())
	}
}

func TestDeviceBufferSurvivesOtosTruncatingRoundTrip(t *testing.T) {
	for _, frames := range []int64{64, 512, 1024, 16384} {
		// oto turns the duration back into bytes of stereo float32 by truncating.
		const otoBytesPerFrame = 2 * 4
		perSecond := int64(44100) * otoBytesPerFrame
		bytes := int64(deviceBuffer(frames)) * perSecond / int64(time.Second)
		if got := bytes / otoBytesPerFrame; got != frames {
			t.Errorf("%d frames came back as %d", frames, got)
		}
	}
}

func TestOptionsSetTheBudget(t *testing.T) {
	l, err := resolve(Options{Queue: 60 * time.Millisecond, DeviceFrames: 2048}, env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if l.frames != 2048 || l.queue != 60*time.Millisecond {
		t.Fatalf("got %+v", l)
	}
}

func TestEnvironmentWinsOverOptions(t *testing.T) {
	l, err := resolve(Options{Queue: 60 * time.Millisecond, DeviceFrames: 2048},
		env(map[string]string{envQueue: "30ms", envDeviceFrames: "512"}))
	if err != nil {
		t.Fatal(err)
	}
	if l.frames != 512 || l.queue != 30*time.Millisecond {
		t.Fatalf("got %+v", l)
	}
}

func TestDefaultQueueStretchesToALongerDeviceBuffer(t *testing.T) {
	l, err := resolve(Options{DeviceFrames: 4096}, env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if l.queue != deviceBuffer(4096) {
		t.Fatalf("queue %v, want one 4096-frame buffer %v", l.queue, deviceBuffer(4096))
	}
}

func TestBadValuesAreErrorsNotDefaults(t *testing.T) {
	cases := []struct {
		name string
		opt  Options
		env  map[string]string
		want string
	}{
		{"queue env not a duration", Options{}, map[string]string{envQueue: "fast"}, envQueue},
		{"queue env not positive", Options{}, map[string]string{envQueue: "0s"}, envQueue},
		{"frames env not an int", Options{}, map[string]string{envDeviceFrames: "1k"}, envDeviceFrames},
		{"frames not a power of two", Options{DeviceFrames: 1000}, nil, "power of two"},
		{"frames too small", Options{DeviceFrames: 32}, nil, "power of two"},
		{"frames too large", Options{DeviceFrames: 32768}, nil, "power of two"},
		{"queue under one buffer", Options{Queue: 10 * time.Millisecond}, nil, "shorter than"},
		{"negative queue", Options{Queue: -time.Millisecond}, nil, "shorter than"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := resolve(c.opt, env(c.env))
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("err %v, want one mentioning %q", err, c.want)
			}
		})
	}
}

func TestPlayerBytesAlignToAFrame(t *testing.T) {
	for _, q := range []time.Duration{24 * time.Millisecond, 40 * time.Millisecond, 23219955} {
		if n := playerBytes(q); n%4 != 0 || n <= 0 {
			t.Errorf("%v gave %d bytes", q, n)
		}
	}
}
