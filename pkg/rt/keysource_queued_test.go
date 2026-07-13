/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"errors"
	"io"
	"reflect"
	"strings"
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

// testGrace is a short inter-byte grace so timing-dependent cases resolve fast.
const testGrace = 15 * time.Millisecond

// newSource builds a queued source over r with the test grace. Constructs the
// unexported type directly (same package) so tests can shorten the grace.
func newSource(r io.Reader) *queuedKeySource {
	return &queuedKeySource{r: r, notify: make(chan struct{}, 1), grace: testGrace}
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
	return newSource(&chunkReader{ch: ch})
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

// TestQueuedKeySourceStitchesAcrossReads drives the wait-for-more path
// deterministically: an arrow arrives byte-by-byte over an open channel — the
// worst case the review called out (ESC, then [, then A as separate reads).
// ReadKey must stitch it into one token, not emit ESC early and split the rest.
func TestQueuedKeySourceStitchesAcrossReads(t *testing.T) {
	for _, tc := range []struct {
		name   string
		chunks []string
		want   string
	}{
		{"esc bracket A separately", []string{"\x1b", "[", "A"}, "\x1b[A"},
		{"esc then bracket-A", []string{"\x1b", "[A"}, "\x1b[A"},
		{"esc-bracket then A", []string{"\x1b[", "A"}, "\x1b[A"},
		{"utf8 byte by byte", []string{"\xe2", "\x8c", "\x98"}, "⌘"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ch := make(chan []byte) // unbuffered: each send blocks until read
			ks := newSource(&chunkReader{ch: ch})
			res := make(chan string, 1)
			go func() {
				k, _ := ks.ReadKey()
				res <- k
			}()
			for _, c := range tc.chunks {
				ch <- []byte(c) // each byte arrives within the inter-byte grace
			}
			select {
			case got := <-res:
				if got != tc.want {
					t.Errorf("stitched token = %q, want %q", got, tc.want)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("ReadKey did not return after the sequence completed")
			}
			close(ch)
		})
	}
}

// TestQueuedKeySourceBareEscEmitsAfterGrace: a lone ESC with no completing byte
// (channel stays open, so no EOF short-circuit) must resolve to the Escape key
// once the inter-byte grace expires — not hang, and not split.
func TestQueuedKeySourceBareEscEmitsAfterGrace(t *testing.T) {
	ch := make(chan []byte, 1)
	ks := newSource(&chunkReader{ch: ch})
	ch <- []byte("\x1b") // no close — reader blocks after this, done stays false

	res := make(chan string, 1)
	go func() {
		k, _ := ks.ReadKey()
		res <- k
	}()
	select {
	case got := <-res:
		if got != "\x1b" {
			t.Errorf("bare ESC = %q, want %q", got, "\x1b")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("bare ESC never resolved — likely parked waiting for more bytes")
	}
	close(ch)
}

// TestQueuedKeySourceReaderErrorPropagates: a non-EOF reader error must surface
// through ReadKey's error return after buffered bytes drain, not be swallowed as
// a clean EOF.
func TestQueuedKeySourceReaderErrorPropagates(t *testing.T) {
	wantErr := errors.New("console read failed")
	ks := newSource(&errAfterReader{data: []byte("ab"), err: wantErr})

	// The bytes that arrived before the error still tokenize.
	for _, want := range []string{"a", "b"} {
		if k, err := ks.ReadKey(); k != want || err != nil {
			t.Fatalf("ReadKey = (%q, %v), want (%q, nil)", k, err, want)
		}
	}
	// Then the retained error surfaces.
	if k, err := ks.ReadKey(); k != "" || err != wantErr {
		t.Fatalf("ReadKey at error = (%q, %v), want (\"\", %v)", k, err, wantErr)
	}
}

// errAfterReader yields data once, then returns err on the next read.
type errAfterReader struct {
	data []byte
	err  error
	done bool
}

func (r *errAfterReader) Read(p []byte) (int, error) {
	if !r.done {
		r.done = true
		return copy(p, r.data), nil
	}
	return 0, r.err
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

// TestQueuedKeySourceGraceEnvOverride: NewQueuedKeySource honors
// LETGO_ESC_GRACE_MS when it holds a positive integer, and falls back to the
// default for empty/invalid values. Reads the unexported grace directly (same
// package) since it isn't otherwise observable.
func TestQueuedKeySourceGraceEnvOverride(t *testing.T) {
	t.Setenv("LETGO_ESC_GRACE_MS", "120")
	if g := NewQueuedKeySource(strings.NewReader("")).(*queuedKeySource).grace; g != 120*time.Millisecond {
		t.Fatalf("grace = %v, want 120ms", g)
	}
	t.Setenv("LETGO_ESC_GRACE_MS", "nonsense")
	if g := NewQueuedKeySource(strings.NewReader("")).(*queuedKeySource).grace; g != escGrace {
		t.Fatalf("invalid override: grace = %v, want default %v", g, escGrace)
	}
}
