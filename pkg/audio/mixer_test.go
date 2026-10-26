package audio

import (
	"bytes"
	"encoding/binary"
	"math"
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

// A one-shot sample that ends on a loud frame fades to silence instead of
// cutting off: no frame steps far from the one before, and the last is near 0.
func TestSampleFadesAtItsEnd(t *testing.T) {
	m := NewMixer()
	pcm := make([]int16, 2*300)
	for i := range pcm {
		pcm[i] = 10000
	}
	m.Play(rt.Sound{PCM: pcm, LoopFrom: -1, Gain: 1})
	out := left(read(m, 400))
	prev, worst := 0, 0
	for _, x := range out {
		worst = max(worst, abs(int(x)-prev))
		prev = int(x)
	}
	if worst > 100 {
		t.Errorf("the sample stepped the mix by %d in one frame: a click", worst)
	}
	if last := out[299]; last > 100 {
		t.Errorf("the sample's last frame is %d, not near silence", last)
	}
}

// A stream that runs dry partway through a mix plays every frame it returned
// before it leaves the mix.
func TestStreamPlaysItsFinalShortBuffer(t *testing.T) {
	m := NewMixer()
	pcm := make([]byte, 4*300)
	for i := 0; i < len(pcm); i += 2 {
		binary.LittleEndian.PutUint16(pcm[i:], 1000)
	}
	m.Stream(bytes.NewReader(pcm))
	out := left(read(m, 400))
	if out[250] != 1000 || out[299] != 1000 {
		t.Errorf("frames 250 and 299 are %d and %d, want the stream's 1000", out[250], out[299])
	}
	if out[300] != 0 || len(m.voices) != 0 {
		t.Errorf("after the stream ran dry: frame 300 is %d with %d voices, want 0 and none", out[300], len(m.voices))
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

// A Gain call during a tone's end fade does not undo it: the tone still ends
// at silence instead of cutting off at the new level.
func TestToneEndFadeOutlastsAGain(t *testing.T) {
	m := NewMixer()
	h, _ := m.Play(rt.Sound{LoopFrom: -1, Gain: 1, Tone: rt.Tone{Freq: 441, Dur: 50 * time.Millisecond, Wave: "square"}})
	frames := rt.AudioRate / 20
	read(m, frames-edgeFrames/2) // into the end fade
	m.Gain(h, 1, 0)
	out := read(m, edgeFrames)
	if last := out[2*(edgeFrames/2-1)]; abs(int(last)) > 200 {
		t.Errorf("the tone ends at %d after a Gain in its end fade, not near silence", last)
	}
}

// Play refuses a sound it cannot render rather than panicking in the mix.
func TestPlayRefusesWhatItCannotRender(t *testing.T) {
	m := NewMixer()
	for name, s := range map[string]rt.Sound{
		"odd PCM":           {PCM: make([]int16, 3), LoopFrom: -1, Gain: 1},
		"empty loop":        {LoopFrom: 0, Gain: 1},
		"loop past the end": {PCM: make([]int16, 4), LoopFrom: 2, Gain: 1},
		"NaN gain":          {PCM: make([]int16, 4), LoopFrom: -1, Gain: math.NaN()},
		"infinite gain":     {PCM: make([]int16, 4), LoopFrom: -1, Gain: math.Inf(1)},
		"infinite tone":     {LoopFrom: -1, Gain: 1, Tone: rt.Tone{Freq: math.Inf(1), Dur: time.Second}},
	} {
		if _, err := m.Play(s); err == nil {
			t.Errorf("%s: Play accepted it", name)
		}
	}
	if len(m.voices) != 0 {
		t.Errorf("%d refused sounds were added to the mix", len(m.voices))
	}
}

// Gain ignores a level that is not finite, so the mix never turns NaN.
func TestGainIgnoresNaN(t *testing.T) {
	m := NewMixer()
	h, _ := m.Stream(constant(1000))
	read(m, edgeFrames)
	m.Gain(h, math.NaN(), 0)
	if got := read(m, 10)[18]; got != 1000 {
		t.Errorf("after Gain(NaN) the level is %d, want 1000", got)
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
