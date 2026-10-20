package audio

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/rt"
)

// read mixes n frames and returns them as interleaved samples.
func read(m *Mixer, n int) []int16 {
	b := make([]byte, 4*n)
	m.Read(b)
	out := make([]int16, 2*n)
	for i := range out {
		out[i] = int16(binary.LittleEndian.Uint16(b[2*i:]))
	}
	return out
}

// constant is a stream source: every sample is v, forever.
type constant int16

func (c constant) Read(p []byte) (int, error) {
	for i := 0; i+1 < len(p); i += 2 {
		binary.LittleEndian.PutUint16(p[i:], uint16(c))
	}
	return len(p) / 2 * 2, nil
}

// A sample with a loop point plays its intro once, then wraps to the loop
// frame with no gap, sample-exact.
func TestSampleLoopsFromItsLoopFrame(t *testing.T) {
	m := NewMixer()
	pcm := make([]int16, 2*(edgeFrames+4))
	for i := range pcm {
		pcm[i] = int16(i / 2) // frame index, both channels
	}
	m.Play(rt.Sound{PCM: pcm, LoopFrom: edgeFrames, Gain: 1})
	read(m, edgeFrames) // past the fade in
	got := read(m, 10)
	want := []int16{220, 221, 222, 223, 220, 221, 222, 223, 220, 221}
	for i, w := range want {
		if got[2*i] != w {
			t.Fatalf("frame %d: got %d, want %d (all: %v)", i, got[2*i], w, left(got))
		}
	}
}

// A one-shot sample ends and leaves the mix.
func TestSampleEndsAndIsDropped(t *testing.T) {
	m := NewMixer()
	m.Play(rt.Sound{PCM: make([]int16, 2*300), LoopFrom: -1, Gain: 1})
	read(m, 400)
	if len(m.voices) != 0 {
		t.Errorf("%d voices still playing after the sample ended", len(m.voices))
	}
}

// A stream and a sample mix by sum, and the stream's handle stops it with a
// fade rather than a cut.
func TestStreamMixesWithASampleAndStopsWithAFade(t *testing.T) {
	m := NewMixer()
	h, _ := m.Stream(constant(1000))
	pcm := make([]int16, 2*2000)
	for i := range pcm {
		pcm[i] = 500
	}
	m.Play(rt.Sound{PCM: pcm, LoopFrom: -1, Gain: 1})
	read(m, edgeFrames+1)
	if got := read(m, 1); got[0] != 1500 || got[1] != 1500 {
		t.Fatalf("stream + sample = %v, want 1500 on both channels", got)
	}
	m.Stop(h)
	prev, worst := 1500, 0
	for _, x := range left(read(m, edgeFrames+10)) {
		worst = max(worst, abs(int(x)-prev))
		prev = int(x)
	}
	if prev != 500 {
		t.Errorf("after the stop fade the mix is %d, want the sample's 500", prev)
	}
	if worst > 10 {
		t.Errorf("stopping stepped the mix by %d in one frame: a click", worst)
	}
}

// A tone lasts its duration and ends at silence.
func TestToneLastsItsDuration(t *testing.T) {
	m := NewMixer()
	m.Play(rt.Sound{LoopFrom: -1, Gain: 1, Tone: rt.Tone{Freq: 441, Dur: 50 * time.Millisecond, Wave: "square"}})
	frames := rt.AudioRate / 20
	out := read(m, frames)
	if len(m.voices) != 0 {
		t.Fatalf("the tone is still playing after %d frames", frames)
	}
	if out[2*(frames/2)] == 0 {
		t.Error("the tone is silent in the middle")
	}
	if last := out[2*(frames-1)]; abs(int(last)) > 200 {
		t.Errorf("the tone ends at %d, not near silence", last)
	}
}

// Gain moves to its target over the fade, not at once.
func TestGainFades(t *testing.T) {
	m := NewMixer()
	h, _ := m.Stream(constant(1000))
	read(m, edgeFrames)
	m.Gain(h, 0, 100*time.Millisecond)
	mid := read(m, rt.AudioRate/20)[2*(rt.AudioRate/20-1)]
	if mid < 400 || mid > 600 {
		t.Errorf("halfway through a 100 ms fade the level is %d, want about 500", mid)
	}
}

// left is the left channel of interleaved samples.
func left(s []int16) []int16 {
	out := make([]int16, len(s)/2)
	for i := range out {
		out[i] = s[2*i]
	}
	return out
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
