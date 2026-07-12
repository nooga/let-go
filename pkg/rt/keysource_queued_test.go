/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"io"
	"reflect"
	"testing"
	"time"
)

// chunkReader hands out one queued chunk per Read (splitting a chunk across
// Reads only if the caller's buffer is smaller), blocks when the queue is empty,
// and returns io.EOF once the channel is closed. It lets a test control exactly
// how input bytes are split across reads — the thing the queued source has to
// stitch back together.
type chunkReader struct {
	ch  chan []byte
	rem []byte
}

func (c *chunkReader) Read(p []byte) (int, error) {
	if len(c.rem) == 0 {
		b, ok := <-c.ch
		if !ok {
			return 0, io.EOF
		}
		c.rem = b
	}
	n := copy(p, c.rem)
	c.rem = c.rem[n:]
	return n, nil
}

// feed returns a queued source over the given chunks, delivered in order then
// EOF. All chunks are queued up front, so ReadKey sees end-of-input after the
// last one.
func feed(chunks ...[]byte) KeySource {
	ch := make(chan []byte, len(chunks))
	for _, c := range chunks {
		ch <- c
	}
	close(ch)
	return NewQueuedKeySource(&chunkReader{ch: ch})
}

// readAll drains keys until the EOF nil contract ("") and returns the tokens.
// Guarded by a deadline so a stuck ReadKey fails loudly instead of hanging CI.
func readAll(t *testing.T, ks KeySource) []string {
	t.Helper()
	var got []string
	done := make(chan []string, 1)
	go func() {
		var out []string
		for {
			k, err := ks.ReadKey()
			if err != nil {
				t.Errorf("ReadKey error: %v", err)
				break
			}
			if k == "" {
				break
			}
			out = append(out, k)
		}
		done <- out
	}()
	select {
	case got = <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("readAll timed out — ReadKey blocked")
	}
	return got
}

func TestQueuedKeySourceTokenizes(t *testing.T) {
	cases := []struct {
		name   string
		chunks [][]byte
		want   []string
	}{
		{"held-key burst", [][]byte{[]byte("llll")}, []string{"l", "l", "l", "l"}},
		{"arrow one chunk", [][]byte{[]byte("\x1b[A")}, []string{"\x1b[A"}},
		{"lone esc", [][]byte{[]byte("\x1b")}, []string{"\x1b"}},
		{"ss3 f1", [][]byte{[]byte("\x1bOP")}, []string{"\x1bOP"}},
		{"mixed keys and csi", [][]byte{[]byte("l\x1b[Bx")}, []string{"l", "\x1b[B", "x"}},
		// A CSI arrow split across two reads must come back as one token.
		{"split arrow", [][]byte{[]byte("\x1b["), []byte("A")}, []string{"\x1b[A"}},
		// A multi-byte rune split across reads must never be torn.
		{"split utf8", [][]byte{{0xe2}, {0x8c, 0x98}}, []string{"⌘"}},
		// A truncated CSI at EOF is emitted best-effort rather than lost.
		{"truncated csi at eof", [][]byte{[]byte("\x1b[")}, []string{"\x1b["}},
		{"empty then eof", nil, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := readAll(t, feed(tc.chunks...))
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("tokens = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestQueuedKeySourceStitchesWhileBlocked drives the wait-for-more path
// deterministically: the partial CSI arrives while ReadKey is already parked, so
// ReadKey must wait for the completing byte and stitch, not emit "\x1b[".
func TestQueuedKeySourceStitchesWhileBlocked(t *testing.T) {
	ch := make(chan []byte) // unbuffered: each send blocks until the reader takes it
	ks := NewQueuedKeySource(&chunkReader{ch: ch})

	res := make(chan string, 1)
	go func() {
		k, _ := ks.ReadKey()
		res <- k
	}()

	ch <- []byte("\x1b[") // reader appends; ReadKey wakes, sees keyNeedMore, waits
	ch <- []byte("A")     // completes the sequence; ReadKey stitches and returns

	select {
	case got := <-res:
		if got != "\x1b[A" {
			t.Errorf("stitched token = %q, want %q", got, "\x1b[A")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ReadKey did not return after the sequence completed")
	}
	close(ch)
}

// TestQueuedKeySourceKeyPending verifies the non-blocking, eof-blind peek:
// false before input, true once a byte is buffered, false again once drained at
// end-of-input (so a poll loop doesn't busy-spin on EOF).
func TestQueuedKeySourceKeyPending(t *testing.T) {
	ch := make(chan []byte, 1)
	ks := NewQueuedKeySource(&chunkReader{ch: ch})

	// Starts the reader; nothing queued yet → not pending.
	if ks.KeyPending() {
		t.Fatal("KeyPending true before any input")
	}

	ch <- []byte("x")
	if !eventually(func() bool { return ks.KeyPending() }) {
		t.Fatal("KeyPending never became true after input")
	}

	if k, _ := ks.ReadKey(); k != "x" {
		t.Fatalf("ReadKey = %q, want %q", k, "x")
	}
	if ks.KeyPending() {
		t.Fatal("KeyPending true after the only key was consumed")
	}

	// EOF must not read as pending (else xsofy's input poll busy-drains it).
	close(ch)
	if eventually(func() bool { return ks.KeyPending() }) {
		t.Fatal("KeyPending true at EOF")
	}
}

func eventually(cond func() bool) bool {
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return cond()
}
