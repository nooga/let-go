/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"context"
	"io"
	"sync"
	"time"
)

// AudioSources is implemented by a backend that can hand back the source it
// is mixing for a stream handle. audio/feed! uses it to reach the ring behind
// a stream lg opened, so the backend, not a registry here, decides how long
// that ring lives.
type AudioSources interface {
	Source(AudioHandle) io.Reader
}

// StreamNotify is implemented by a stream source that wants to hear about its
// stream's lifecycle. A host calls StreamHeld(true) on Pause and
// StreamHeld(false) on Resume, and StreamEnded once on Stop or when it drops
// the stream, without holding any lock of its own. A source that doesn't
// implement it is never called, so an io.Closer passed to Stream is not
// closed for its owner.
type StreamNotify interface {
	StreamHeld(held bool)
	StreamEnded()
}

// feedRing is the source behind a stream lg code opens with audio/stream:
// audio/feed! writes PCM in, and the host's mixer reads it out on its audio
// goroutine through Audio.Stream.
type feedRing struct {
	mu           sync.Mutex
	buf          []byte
	r, n         int
	closed, held bool
	space        chan struct{} // a wakeup for a waiting writer: the reader freed room
	wake         chan struct{} // closed and replaced on StreamEnded and StreamHeld, waking every waiting writer
	lastRead     time.Time
	stall        time.Duration
}

// newFeedRing holds size bytes, which play for d. A writer gives up waiting
// once nothing has read the ring for twice that, or a second at least: the
// host stopped draining it (a lost device, a mixer nobody pumps), and a guest
// waiting on it would hang for good.
func newFeedRing(size int, d time.Duration) *feedRing {
	return &feedRing{
		buf: make([]byte, size), space: make(chan struct{}, 1), wake: make(chan struct{}),
		lastRead: time.Now(), stall: max(time.Second, 2*d),
	}
}

// Read never blocks and never comes back short. The mixer ends a stream's
// voice on any short read, so whatever has not been fed yet plays as silence:
// an underrun is a gap in the sound, not the end of the stream. Read reports
// EOF only once the ring is closed and empty.
func (f *feedRing) Read(p []byte) (int, error) {
	f.mu.Lock()
	if f.closed && f.n == 0 {
		f.mu.Unlock()
		return 0, io.EOF
	}
	f.lastRead = time.Now()
	k := min(len(p), f.n)
	for i := 0; i < k; {
		c := copy(p[i:k], f.buf[f.r:])
		i += c
		f.r = (f.r + c) % len(f.buf)
	}
	f.n -= k
	f.mu.Unlock()
	clear(p[k:])
	if k > 0 {
		select {
		case f.space <- struct{}{}:
		default:
		}
	}
	return len(p), nil
}

// write copies b into the ring, waiting for room while it is full, so a
// guest that feeds faster than the device plays is held back by the
// device's clock. A paused stream isn't drained, so while held write drops
// b instead of waiting for a resume the waiting guest could never issue,
// and drops what is left once the host has stopped reading (newFeedRing).
// It returns early when the ring is closed (the host is done with the
// stream) or ctx is cancelled, and says which.
func (f *feedRing) write(ctx context.Context, b []byte) (closed, cancelled bool) {
	for {
		f.mu.Lock()
		if f.closed {
			f.mu.Unlock()
			return true, false
		}
		if f.held || time.Since(f.lastRead) > f.stall {
			f.mu.Unlock()
			return false, false
		}
		for len(b) > 0 && f.n < len(f.buf) {
			w := (f.r + f.n) % len(f.buf)
			end := len(f.buf)
			if w < f.r {
				end = f.r
			}
			c := copy(f.buf[w:end], b)
			f.n += c
			b = b[c:]
		}
		wake, wait := f.wake, f.stall-time.Since(f.lastRead)
		f.mu.Unlock()
		if len(b) == 0 {
			return false, false
		}
		t := time.NewTimer(max(0, wait))
		select {
		case <-f.space:
		case <-wake:
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return false, true
		}
		t.Stop()
	}
}

// StreamHeld implements StreamNotify.
func (f *feedRing) StreamHeld(held bool) {
	f.mu.Lock()
	f.held = held
	f.broadcast()
	f.mu.Unlock()
}

// StreamEnded implements StreamNotify: the host is done with the stream,
// however it was stopped, so a writer waiting on it returns. What the ring
// still holds plays out through the stop fade, then Read reports EOF.
func (f *feedRing) StreamEnded() {
	f.mu.Lock()
	if !f.closed {
		f.closed = true
		f.broadcast()
	}
	f.mu.Unlock()
}

// broadcast wakes every waiting writer; f.mu must be held.
func (f *feedRing) broadcast() {
	close(f.wake)
	f.wake = make(chan struct{})
}
