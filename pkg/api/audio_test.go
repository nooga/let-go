package api_test

import (
	"io"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/api"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// recordAudio keeps what the guest asked to play, the audio dual of
// captureSurface.
type recordAudio struct {
	played  []rt.Sound
	stopped []rt.AudioHandle
	gains   []float64
	fades   []time.Duration
}

func (r *recordAudio) Play(s rt.Sound) (rt.AudioHandle, error) {
	r.played = append(r.played, s)
	return rt.AudioHandle(len(r.played)), nil
}
func (r *recordAudio) Stream(io.Reader) (rt.AudioHandle, error) { return 0, nil }
func (r *recordAudio) Stop(h rt.AudioHandle)                    { r.stopped = append(r.stopped, h) }
func (r *recordAudio) Gain(h rt.AudioHandle, g float64, fade time.Duration) {
	r.gains, r.fades = append(r.gains, g), append(r.fades, fade)
}
func (r *recordAudio) Enabled() bool { return true }

// An rt.Audio passed with api.WithAudio receives the guest's tones and
// samples, and the handles it returns come back to operate on.
func TestWithAudio(t *testing.T) {
	rec := &recordAudio{}
	lg, err := api.NewLetGo("withaudio-test", api.WithAudio(rec))
	if err != nil {
		t.Fatal(err)
	}
	v, err := lg.Run(`(let [h (audio/tone {:freq 440 :ms 120 :wave :square :gain 0.3})]
	                     (audio/gain h 0.5 200)
	                     (audio/stop h)
	                     (audio/sample (byte-array [1 0 2 0 3 0 4 0]) {:loop-from 1})
	                     (audio/enabled?))`)
	if err != nil {
		t.Fatal(err)
	}
	if v != vm.TRUE {
		t.Errorf("enabled? = %v, want true", v)
	}
	if len(rec.played) != 2 {
		t.Fatalf("played %d sounds, want 2", len(rec.played))
	}
	tone := rec.played[0].Tone
	if tone.Freq != 440 || tone.Dur != 120*time.Millisecond || tone.Wave != "square" || rec.played[0].Gain != 0.3 {
		t.Errorf("tone = %+v gain %v", tone, rec.played[0].Gain)
	}
	if len(rec.stopped) != 1 || rec.stopped[0] != 1 {
		t.Errorf("stopped %v, want the tone's handle 1", rec.stopped)
	}
	if len(rec.gains) != 1 || rec.gains[0] != 0.5 || rec.fades[0] != 200*time.Millisecond {
		t.Errorf("gain %v over %v, want 0.5 over 200ms", rec.gains, rec.fades)
	}
	s := rec.played[1]
	if len(s.PCM) != 4 || s.PCM[0] != 1 || s.PCM[3] != 4 || s.LoopFrom != 1 {
		t.Errorf("sample = %+v, want 4 samples looping from frame 1", s)
	}
}

// With no backend bound, audio/ calls are silent no-ops that return nil
// handles, which stop and gain take as no-ops too, and bad arguments still
// fail, so native dev catches them.
func TestAudioNoopWithoutBackend(t *testing.T) {
	lg, err := api.NewLetGo("audio-noop-test")
	if err != nil {
		t.Fatal(err)
	}
	v, err := lg.Run(`(let [h (audio/tone {:freq 440 :ms 50})]
	                    [h (audio/gain h 0.5 100) (audio/stop h) (audio/stop 3) (audio/enabled?)])`)
	if err != nil {
		t.Fatal(err)
	}
	if got := v.String(); got != "[nil nil nil nil false]" {
		t.Errorf("got %s, want [nil nil nil nil false]", got)
	}
	for _, bad := range []string{
		`(audio/tone {:freq 0 :ms 50})`,
		`(audio/tone {:freq 440 :ms 50 :wave :kazoo})`,
		`(audio/sample (byte-array [1 2 3]))`,
		`(audio/sample (byte-array [1 0 2 0]) {:loop-from 1})`,
		`(audio/gain 1 -1)`,
		`(audio/gain nil -1)`,
		`(audio/stop :h)`,
	} {
		if _, err := lg.Run(bad); err == nil {
			t.Errorf("%s: want an error", bad)
		}
	}
}
