/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

// queuedKeySource is a KeySource for hosts without a non-blocking input peek: a
// background goroutine blocks on an io.Reader and appends to a shared queue,
// ReadKey tokenizes one key per call off it (via scanKey), and KeyPending peeks
// it. This is the shape Plan 9 needs — it has no poll(2)/FIONREAD, so a
// non-blocking key-pending? can only be answered from a buffer a concurrent
// reader fills, where native (term.go) reads stdin synchronously inside ReadKey.
// It is platform-neutral, so it is also available to embedders whose input is a
// plain blocking reader (a pipe, a socket) — see NewQueuedKeySource.
//
// State lives on the instance, not a package global like native's keyBuf: the
// goroutine consumes the reader, so it must start only when THIS source is the
// one being read. If an embedder rebinds *keys* via api.WithKeySource, this
// source's ReadKey/KeyPending are never called and the reader never starts — no
// bytes stolen from the source that replaced it.
type queuedKeySource struct {
	mu      sync.Mutex
	r       io.Reader
	buf     []byte        // read but not yet tokenized out
	err     error         // sticky reader error; nil for a clean EOF
	done    bool          // reader finished (EOF or error) — no more bytes coming
	started bool          // reader goroutine launched
	notify  chan struct{} // buffered(1) wakeup; reader signals on append/done
	grace   time.Duration // inter-byte wait for stitching a split escape sequence
}

// escGrace is the default bound on how long ReadKey waits for the rest of an
// incomplete escape (or UTF-8) sequence before emitting what it has. It's the
// classic terminal "ESC delay": long enough to stitch a sequence split across
// reads on a slow transport (drawterm/9P), short enough that a bare Escape key
// still registers promptly. Only an incomplete front token pays it — a complete
// key returns with zero added latency. Without it, distinguishing a bare ESC
// from the start of a split arrow is impossible on a host with no FIONREAD peek
// (a lone ESC would either be emitted too eagerly, breaking split arrows, or
// waited on forever).
const escGrace = 40 * time.Millisecond

// escGraceEnv overrides escGrace, in milliseconds, for slow/high-latency
// transports — drawterm over a WAN link, where the 40ms default can lapse
// mid-sequence and split an arrow into a stray ESC + [A. Read once per source at
// construction.
const escGraceEnv = "LETGO_ESC_GRACE_MS"

// resolveGrace returns the inter-byte escape grace, honoring escGraceEnv when it
// holds a positive integer and falling back to escGrace otherwise.
func resolveGrace() time.Duration {
	if v := os.Getenv(escGraceEnv); v != "" {
		if ms, err := strconv.Atoi(v); err == nil && ms > 0 {
			return time.Duration(ms) * time.Millisecond
		}
	}
	return escGrace
}

// NewQueuedKeySource returns a KeySource that reads keys from r via a background
// goroutine and a buffered queue, so KeyPending is non-blocking on platforms
// that lack a poll/FIONREAD peek. The reader goroutine starts lazily on the
// first ReadKey/KeyPending, so binding this at *keys* without ever consulting it
// (e.g. when api.WithKeySource overrides it) reads nothing.
func NewQueuedKeySource(r io.Reader) KeySource {
	return &queuedKeySource{r: r, notify: make(chan struct{}, 1), grace: resolveGrace()}
}

// queuedReadChunkSize matches native's readChunkSize: a blocking Read returns as
// soon as any data is available, so a larger buffer adds no latency to a single
// keystroke — it just collects held-key bursts in fewer reads.
const queuedReadChunkSize = 256

// start launches the reader goroutine once. The caller must hold s.mu: start
// flips s.started, so the guard against a double launch relies on the lock. The
// goroutine it spawns (readLoop) then does its own finer-grained locking and
// does not hold s.mu while blocked on Read.
func (s *queuedKeySource) start() {
	if s.started {
		return
	}
	s.started = true
	go s.readLoop()
}

// signal wakes a waiting ReadKey without blocking; the buffered channel
// coalesces bursts, and ReadKey re-checks the buffer under the lock after each
// wake, so a dropped (already-pending) signal is never a missed update.
func (s *queuedKeySource) signal() {
	select {
	case s.notify <- struct{}{}:
	default:
	}
}

// readLoop blocks on r.Read WITHOUT holding the lock, then locks only to append
// the chunk and wake a waiter. On a read error it retains the error (draining
// any bytes that came with it) and stops; io.EOF is a clean end, recorded as a
// nil error so ReadKey returns the "" contract. A (0, nil) read is not EOF (the
// io.Reader contract) — it loops and reads again.
func (s *queuedKeySource) readLoop() {
	chunk := make([]byte, queuedReadChunkSize)
	for {
		n, err := s.r.Read(chunk)
		s.mu.Lock()
		if n > 0 {
			s.buf = append(s.buf, chunk[:n]...)
		}
		if err != nil {
			if err != io.EOF {
				s.err = err
			}
			s.done = true
			s.mu.Unlock()
			s.signal()
			return
		}
		s.mu.Unlock()
		if n > 0 {
			s.signal()
		}
	}
}

// ReadKey blocks until a whole key token is available (or the reader ends),
// tokenizing one key per call off the shared buffer via scanKey. It returns ""
// with a nil error at a clean EOF (the read-key nil contract), or "" with the
// retained error if the reader failed. Blocking here matches native (which
// blocks in poll); callers stay responsive by gating on KeyPending first.
func (s *queuedKeySource) ReadKey() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.start()

	for len(s.buf) == 0 && !s.done {
		s.condWait(0) // no deadline: block until data or reader end
	}
	if len(s.buf) == 0 {
		return "", s.err // clean EOF (err nil) or a real reader failure
	}

	var deadline time.Time // set on the first incomplete front token
	for {
		status, n := scanKey(s.buf)
		if status == keyReady || s.done {
			// Complete token, or the reader ended so no completion is coming —
			// emit best-effort (a lone ESC at EOF is the Escape key; a truncated
			// sequence surfaces rather than hangs).
			return s.take(n), nil
		}
		// keyNeedMore and the reader is still live: a front token split across
		// reads. Wait for more bytes up to the grace period, then re-scan; if the
		// grace expires with no completion, emit what we have. This is what makes
		// a bare ESC vs. a split arrow decidable without a FIONREAD peek. The
		// grace is inter-byte — each byte that arrives resets it — so a sequence
		// dribbling in over a slow transport keeps stitching, while a truly bare
		// ESC (nothing more within one grace window) emits as the Escape key.
		if deadline.IsZero() {
			deadline = time.Now().Add(s.grace)
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return s.take(n), nil
		}
		before := len(s.buf)
		s.condWait(remaining)
		if len(s.buf) > before {
			deadline = time.Time{} // new bytes — restart the inter-byte window
		}
	}
}

// take slices the front n bytes off the buffer as a token, releasing the backing
// array once drained. Called with s.mu held.
func (s *queuedKeySource) take(n int) string {
	tok := string(s.buf[:n])
	s.buf = s.buf[n:]
	if len(s.buf) == 0 {
		s.buf = nil
	}
	return tok
}

// condWait releases s.mu, waits for a reader wake (or the timeout, if nonzero),
// then re-acquires it — a condition-variable wait done over a channel so it can
// be time-bounded. Named for the sync.Cond.Wait pattern: the lock is dropped for
// the duration of the wait, not held across it. Call with s.mu held; returns
// holding it.
func (s *queuedKeySource) condWait(timeout time.Duration) {
	s.mu.Unlock()
	if timeout > 0 {
		select {
		case <-s.notify:
		case <-time.After(timeout):
		}
	} else {
		<-s.notify
	}
	s.mu.Lock()
}

// KeyPending reports whether a key is buffered and ready, without consuming it.
// Non-blocking and eof-blind (false once the queue drains at end-of-input),
// mirroring native's FIONREAD-based rawPending so a per-tick input poll doesn't
// busy-drain end-of-input forever.
func (s *queuedKeySource) KeyPending() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.start()
	return len(s.buf) > 0
}
