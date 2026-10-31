/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/nooga/let-go/pkg/rt"
)

// Options sets the device's latency budget. The device is one per process, so
// these belong to the host embedding let-go, not to any one sound source. Zero
// fields take the defaults in brackets.
type Options struct {
	// Queue is how far the player reads ahead of the device [24 ms on macOS,
	// else 40 ms, stretched to one device buffer where that is longer].
	// Shorter hears a trigger sooner; longer survives a stalled main thread.
	// A set value shorter than one device buffer is an error.
	Queue time.Duration
	// DeviceFrames is the device buffer in frames [1024]: a power of two from
	// 64 to 16384, the range a browser's ScriptProcessor accepts, so a value
	// tuned here carries over to a browser backend. Raise it first if
	// dropouts appear; 512 crackled in a terminal game on macOS, where it is
	// four AudioQueue buffers of 2.9 ms.
	DeviceFrames int64
}

// The defaults follow undertone's native player (its nativeaudio/latency.go).
// Measured 2026-10-26 on a MacBook Air (speakers into its own mic, trigger to
// onset), 40 ms of read-ahead gave 150-158 ms and 24 ms gave 119-128 ms; 24 ms
// over a 512-frame device crackled under a busy main thread. Only macOS has
// been heard at 24 ms, so other platforms keep 40 ms.
const defaultDeviceFrames = 1024

func defaultQueue() time.Duration {
	if runtime.GOOS == "darwin" { // GOOS is "ios" on iOS, which keeps 40 ms
		return 24 * time.Millisecond
	}
	return 40 * time.Millisecond
}

// Two environment variables override Options without a rebuild, for tuning a
// stock tagged lg by ear: LETGO_AUDIO_QUEUE, a Go duration such as "30ms", and
// LETGO_AUDIO_DEVICE_FRAMES. They win over Options, since whoever sets them is
// tuning the program that set Options. A bad value is an error when the device
// opens, not a silent fallback, so a tuning run never measures the default by
// mistake.
const (
	envQueue        = "LETGO_AUDIO_QUEUE"
	envDeviceFrames = "LETGO_AUDIO_DEVICE_FRAMES"
)

// latency is Options resolved against the environment and the defaults.
type latency struct {
	queue  time.Duration
	frames int64
}

func resolve(opt Options, getenv func(string) string) (latency, error) {
	l := latency{queue: opt.Queue, frames: opt.DeviceFrames}
	if v := getenv(envQueue); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 {
			return latency{}, fmt.Errorf("audio: %s=%q is not a positive duration", envQueue, v)
		}
		l.queue = d
	}
	if v := getenv(envDeviceFrames); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return latency{}, fmt.Errorf("audio: %s=%q is not an integer", envDeviceFrames, v)
		}
		l.frames = n
	}
	if l.frames == 0 {
		l.frames = defaultDeviceFrames
	}
	if l.frames < 64 || l.frames > 16384 || l.frames&(l.frames-1) != 0 {
		return latency{}, fmt.Errorf("audio: %d device frames is not a power of two from 64 to 16384", l.frames)
	}
	// With no queue asked for, stretch the default to one device buffer where
	// a buffer is longer, so the default alone is never an error.
	buf := deviceBuffer(l.frames)
	if l.queue == 0 {
		l.queue = max(defaultQueue(), buf)
	}
	// Shorter than one device buffer can't help (the device drains whole
	// buffers), and under one frame it rounds to 0 bytes, which oto reads as
	// "use the half-second default", the silent fallback this rejects.
	if l.queue < buf {
		return latency{}, fmt.Errorf("audio: a queue of %v is shorter than the %d-frame device buffer (%v)", l.queue, l.frames, buf)
	}
	return l, nil
}

// deviceBuffer is frames as the duration oto wants, rounded up: oto converts
// back by truncating, and a duration that divides exactly-but-short lands one
// frame low, turning 1024 into 1023.
func deviceBuffer(frames int64) time.Duration {
	return time.Duration((frames*int64(time.Second) + rt.AudioRate - 1) / rt.AudioRate)
}

// playerBytes is the queue as oto's SetBufferSize wants it: bytes of signed
// 16-bit stereo, aligned to a whole frame.
func playerBytes(queue time.Duration) int {
	const bytesPerFrame = 2 * 2
	n := int(int64(queue) * rt.AudioRate * bytesPerFrame / int64(time.Second))
	return n / bytesPerFrame * bytesPerFrame
}
