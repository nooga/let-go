package rt

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"
)

// An underrun reads as silence of the full length asked for, never as a
// short read, which would end the mixer's voice.
func TestFeedRingUnderrunIsSilence(t *testing.T) {
	f := newFeedRing(16, 0)
	if closed, cancelled := f.write(context.Background(), []byte{1, 2, 3, 4}); closed || cancelled {
		t.Fatal("write into an empty ring should complete")
	}
	p := bytes.Repeat([]byte{9}, 12)
	n, err := f.Read(p)
	if n != 12 || err != nil {
		t.Fatalf("Read = %d, %v; want 12, nil", n, err)
	}
	if want := []byte{1, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0}; !bytes.Equal(p, want) {
		t.Errorf("read %v, want %v", p, want)
	}
	if n, err := f.Read(p); n != 12 || err != nil {
		t.Errorf("Read of an empty, open ring = %d, %v; want 12, nil", n, err)
	}
}

// Data written across the end of the buffer comes back in order.
func TestFeedRingWraps(t *testing.T) {
	f := newFeedRing(8, 0)
	ctx := context.Background()
	f.write(ctx, []byte{1, 2, 3, 4, 5, 6})
	p := make([]byte, 4)
	f.Read(p)
	f.write(ctx, []byte{7, 8, 9, 10, 11, 12})
	p = make([]byte, 8)
	f.Read(p)
	if want := []byte{5, 6, 7, 8, 9, 10, 11, 12}; !bytes.Equal(p, want) {
		t.Errorf("read %v, want %v", p, want)
	}
}

// A write bigger than the free space waits for the reader, then finishes.
func TestFeedRingWriteWaitsForTheReader(t *testing.T) {
	f := newFeedRing(4, 0)
	done := make(chan struct{})
	go func() {
		f.write(context.Background(), []byte{1, 2, 3, 4, 5, 6, 7, 8})
		close(done)
	}()
	select {
	case <-done:
		t.Fatal("write of 8 bytes into a 4-byte ring returned before anything was read")
	case <-time.After(50 * time.Millisecond):
	}
	got := make([]byte, 0, 8)
	for len(got) < 8 {
		p := make([]byte, 4)
		f.Read(p)
		got = append(got, p...)
		time.Sleep(5 * time.Millisecond)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("write still waiting after the reader drained the ring")
	}
	if want := []byte{1, 2, 3, 4, 5, 6, 7, 8}; !bytes.Equal(got[:8], want) {
		t.Errorf("read %v, want %v", got[:8], want)
	}
}

// Closing the ring (audio/stop) releases a waiting writer.
func TestFeedRingCloseReleasesAWriter(t *testing.T) {
	f := newFeedRing(4, 0)
	res := make(chan bool)
	go func() {
		closed, _ := f.write(context.Background(), make([]byte, 8))
		res <- closed
	}()
	time.Sleep(20 * time.Millisecond)
	f.StreamEnded()
	select {
	case closed := <-res:
		if !closed {
			t.Error("write returned without reporting the ring closed")
		}
	case <-time.After(time.Second):
		t.Fatal("close did not release the waiting write")
	}
}

// Cancelling the writer's scope releases it and says so.
func TestFeedRingWriteIsCancellable(t *testing.T) {
	f := newFeedRing(4, 0)
	ctx, cancel := context.WithCancel(context.Background())
	res := make(chan bool)
	go func() {
		_, cancelled := f.write(ctx, make([]byte, 8))
		res <- cancelled
	}()
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case cancelled := <-res:
		if !cancelled {
			t.Error("write returned without reporting the cancel")
		}
	case <-time.After(time.Second):
		t.Fatal("cancel did not release the waiting write")
	}
}

// A closed ring plays out what it holds, then reports EOF so the voice ends.
func TestFeedRingDrainsThenEOF(t *testing.T) {
	f := newFeedRing(8, 0)
	f.write(context.Background(), []byte{1, 2, 3, 4})
	f.StreamEnded()
	p := make([]byte, 4)
	if n, err := f.Read(p); n != 4 || err != nil || p[0] != 1 {
		t.Fatalf("first Read = %d, %v, %v; want the 4 fed bytes", n, err, p)
	}
	if _, err := f.Read(p); err != io.EOF {
		t.Errorf("Read after drain = %v, want EOF", err)
	}
}

// A paused stream isn't drained, so a write that doesn't fit drops the rest
// instead of waiting: a writer that is waiting when the pause comes returns,
// and one that starts during the pause doesn't wait at all.
func TestFeedRingHoldDropsInsteadOfWaiting(t *testing.T) {
	f := newFeedRing(4, 0)
	res := make(chan bool)
	go func() {
		closed, cancelled := f.write(context.Background(), make([]byte, 8))
		res <- closed || cancelled
	}()
	time.Sleep(20 * time.Millisecond)
	f.StreamHeld(true)
	select {
	case odd := <-res:
		if odd {
			t.Error("a held write reported closed or cancelled; want a plain return")
		}
	case <-time.After(time.Second):
		t.Fatal("Hold(true) did not release the waiting write")
	}
	done := make(chan struct{})
	go func() { f.write(context.Background(), make([]byte, 8)); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("a write into a full, held ring waited")
	}
}
