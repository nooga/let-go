package audio

import (
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/rt"
)

// ramp is a sample whose every frame holds its own index, on both channels.
func ramp(frames int) []int16 {
	pcm := make([]int16, 2*frames)
	for i := range pcm {
		pcm[i] = int16(i / 2)
	}
	return pcm
}

// near reports whether a and b differ by at most one: where a fade lands to
// the frame can depend on float rounding.
func near(a, b int) bool { return abs(a-b) <= 1 }

// A paused sample fades out, stays silent however long it is held, and
// resumes from the frame it stopped on.
func TestPauseHoldsASampleAndResumesFromItsFrame(t *testing.T) {
	m := NewMixer()
	h, _ := m.Play(rt.Sound{PCM: ramp(2000), LoopFrom: -1, Gain: 1})
	read(m, edgeFrames+30) // fade in, then 30 frames at full level
	m.Pause(h)
	out := left(read(m, edgeFrames+10))
	if out[edgeFrames+9] != 0 {
		t.Fatalf("after the pause fade the mix is %d, want silence", out[edgeFrames+9])
	}
	for i, x := range left(read(m, 5000)) {
		if x != 0 {
			t.Fatalf("paused frame %d is %d, want silence", i, x)
		}
	}
	m.Resume(h)
	read(m, edgeFrames)
	// The sample advanced through the fade in, the 30 frames and the pause
	// fade, then held; the resume fade moved it on by edgeFrames.
	want := edgeFrames + 30 + edgeFrames + edgeFrames
	if got := int(read(m, 1)[0]); !near(got, want) {
		t.Errorf("after resuming the sample is at frame %d, want %d", got, want)
	}
}

// counter is a stream that counts the frames it is asked for.
type counter struct{ frames int }

func (c *counter) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	c.frames += len(p) / 4
	return len(p) / 4 * 4, nil
}

// A paused stream is read only to the end of its pause fade, and not again
// until it resumes, so a synth behind it keeps its place.
func TestPausedStreamIsNotRead(t *testing.T) {
	m := NewMixer()
	src := &counter{}
	h, _ := m.Stream(src)
	read(m, edgeFrames+100)
	m.Pause(h)
	read(m, 1000)
	if want := 2*edgeFrames + 100; !near(src.frames, want) {
		t.Fatalf("read %d frames by the end of the pause fade, want %d", src.frames, want)
	}
	held := src.frames
	read(m, 4000)
	if src.frames != held {
		t.Fatalf("the paused stream was read on to frame %d from %d", src.frames, held)
	}
	m.Resume(h)
	read(m, 100)
	if src.frames != held+100 {
		t.Errorf("after resuming, read %d frames, want %d", src.frames-held, 100)
	}
}

// Pausing and resuming a loud stream fade rather than cut.
func TestPauseAndResumeDoNotClick(t *testing.T) {
	m := NewMixer()
	h, _ := m.Stream(constant(10000))
	read(m, edgeFrames+1)
	m.Pause(h)
	a := left(read(m, edgeFrames+10))
	m.Resume(h)
	b := left(read(m, edgeFrames+10))
	prev, worst := 10000, 0
	for _, x := range append(a, b...) {
		worst = max(worst, abs(int(x)-prev))
		prev = int(x)
	}
	if worst > 100 {
		t.Errorf("pause or resume stepped the mix by %d in one frame: a click", worst)
	}
	if prev != 10000 {
		t.Errorf("after resuming the mix is %d, want the stream's 10000", prev)
	}
}

// A gain set while paused is the level the sound resumes to, and does not
// make the paused sound audible.
func TestGainWhilePausedSetsTheResumeLevel(t *testing.T) {
	m := NewMixer()
	h, _ := m.Stream(constant(10000))
	read(m, edgeFrames+1)
	m.Pause(h)
	read(m, edgeFrames+10)
	m.Gain(h, 0.5, 0)
	if x := read(m, 100)[0]; x != 0 {
		t.Fatalf("a gain on a paused sound made it audible: %d", x)
	}
	m.Resume(h)
	read(m, edgeFrames)
	if x := read(m, 1)[0]; x != 5000 {
		t.Errorf("resumed at %d, want 5000 from the gain set while paused", x)
	}
}

// A paused sound can still be stopped, and then leaves the mix.
func TestStopWhilePaused(t *testing.T) {
	m := NewMixer()
	h, _ := m.Stream(constant(1000))
	read(m, edgeFrames+1)
	m.Pause(h)
	read(m, edgeFrames+10)
	m.Stop(h)
	m.Resume(h)
	read(m, 10)
	if len(m.voices) != 0 {
		t.Errorf("%d voices left after stopping a paused stream", len(m.voices))
	}
}

// A paused tone keeps the rest of its duration for after it resumes.
func TestPausedToneKeepsItsRemainingLength(t *testing.T) {
	m := NewMixer()
	h, _ := m.Play(rt.Sound{LoopFrom: -1, Gain: 1, Tone: rt.Tone{Freq: 441, Dur: 50 * time.Millisecond, Wave: "square"}})
	read(m, 1000)
	m.Pause(h)
	read(m, 10000) // far past the tone's 2205 frames
	if len(m.voices) != 1 {
		t.Fatal("the paused tone ended while held")
	}
	m.Resume(h)
	// 1000 frames played and edgeFrames in the pause fade; the rest of the
	// 2205, the resume fade included, follow.
	rest := 2205 - 1000 - edgeFrames
	read(m, rest-10)
	if len(m.voices) != 1 {
		t.Fatal("the resumed tone ended early")
	}
	read(m, 20)
	if len(m.voices) != 0 {
		t.Error("the resumed tone outlasted its duration")
	}
}

// Pause and Resume on an unknown handle, a repeated Pause, and a Resume of a
// sound that isn't paused change nothing.
func TestPauseAndResumeAreIdempotent(t *testing.T) {
	m := NewMixer()
	h, _ := m.Stream(constant(1000))
	read(m, edgeFrames+1)
	m.Resume(h)
	m.Pause(99)
	if x := read(m, 1)[0]; x != 1000 {
		t.Fatalf("an unpaused stream is at %d, want 1000", x)
	}
	m.Pause(h)
	read(m, edgeFrames/2)
	m.Pause(h) // must not restart the fade from the halfway level as "the" level
	read(m, edgeFrames)
	m.Resume(h)
	read(m, edgeFrames)
	if x := read(m, 1)[0]; x != 1000 {
		t.Errorf("after a double pause the stream resumed to %d, want 1000", x)
	}
}
