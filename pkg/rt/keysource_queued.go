/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"io"
	"sync"
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
	cond    *sync.Cond
	r       io.Reader
	buf     []byte // read but not yet tokenized out
	eof     bool
	started bool
}

// NewQueuedKeySource returns a KeySource that reads keys from r via a background
// goroutine and a buffered queue, so KeyPending is non-blocking on platforms
// that lack a poll/FIONREAD peek. The reader goroutine starts lazily on the
// first ReadKey/KeyPending, so binding this at *keys* without ever consulting it
// (e.g. when api.WithKeySource overrides it) reads nothing.
func NewQueuedKeySource(r io.Reader) KeySource {
	s := &queuedKeySource{r: r}
	s.cond = sync.NewCond(&s.mu)
	return s
}

// queuedReadChunkSize matches native's readChunkSize: a blocking Read returns as
// soon as any data is available, so a larger buffer adds no latency to a single
// keystroke — it just collects held-key bursts in fewer reads.
const queuedReadChunkSize = 256

// start launches the background reader once. Called with s.mu held.
func (s *queuedKeySource) start() {
	if s.started {
		return
	}
	s.started = true
	go s.readLoop()
}

// readLoop blocks on r.Read WITHOUT holding the lock, then locks only to append
// the chunk and wake a waiter. Runs until EOF or a read error, which it records
// as eof so a parked ReadKey unblocks with the nil contract.
func (s *queuedKeySource) readLoop() {
	chunk := make([]byte, queuedReadChunkSize)
	for {
		n, err := s.r.Read(chunk)
		s.mu.Lock()
		if n > 0 {
			s.buf = append(s.buf, chunk[:n]...)
		}
		if err != nil || n == 0 {
			s.eof = true
			s.cond.Broadcast()
			s.mu.Unlock()
			return
		}
		s.cond.Broadcast()
		s.mu.Unlock()
	}
}

// ReadKey blocks until a whole key token is available (or EOF), tokenizing one
// key per call off the shared buffer via scanKey. "" signals end-of-input (the
// read-key nil contract). Blocking here matches native (which blocks in poll);
// callers stay responsive by gating on KeyPending first.
func (s *queuedKeySource) ReadKey() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.start()
	for len(s.buf) == 0 && !s.eof {
		s.cond.Wait()
	}
	if len(s.buf) == 0 {
		return "", nil // EOF
	}
	for {
		status, n := scanKey(s.buf)
		if status == keyNeedMore && n > 1 && !s.eof {
			// A multi-byte token (CSI/SS3/UTF-8) split across reader chunks —
			// short reads are normal over a network transport like 9P/drawterm.
			// Wait for one more append, then re-scan to stitch it, rather than
			// emitting a broken partial that would desync every following key.
			// Bounded: each wait adds bytes (progress) or hits eof (then emit
			// best-effort below). A lone ESC (scanKey → keyNeedMore, 1) is
			// deliberately NOT waited on — it is the Escape key, emitted now
			// rather than parked until the next keypress.
			before := len(s.buf)
			for len(s.buf) == before && !s.eof {
				s.cond.Wait()
			}
			continue
		}
		tok := string(s.buf[:n])
		s.buf = s.buf[n:]
		if len(s.buf) == 0 {
			s.buf = nil // release the backing array once drained
		}
		return tok, nil
	}
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
