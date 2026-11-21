package api_test

import (
	"encoding/binary"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/api"
	"github.com/nooga/let-go/pkg/audio"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// recordAudio keeps what the guest asked to play, the audio dual of
// captureSurface.
type recordAudio struct {
	played  []rt.Sound
	stopped []rt.AudioHandle
	paused  []rt.AudioHandle
	resumed []rt.AudioHandle
	gains   []float64
	fades   []time.Duration
}

func (r *recordAudio) Play(s rt.Sound) (rt.AudioHandle, error) {
	r.played = append(r.played, s)
	return rt.AudioHandle(len(r.played)), nil
}
func (r *recordAudio) Stream(io.Reader) (rt.AudioHandle, error) { return 0, nil }
func (r *recordAudio) Stop(h rt.AudioHandle)                    { r.stopped = append(r.stopped, h) }
func (r *recordAudio) Pause(h rt.AudioHandle)                   { r.paused = append(r.paused, h) }
func (r *recordAudio) Resume(h rt.AudioHandle)                  { r.resumed = append(r.resumed, h) }
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
	                     (audio/pause h)
	                     (audio/resume h)
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
	if len(rec.paused) != 1 || rec.paused[0] != 1 || len(rec.resumed) != 1 || rec.resumed[0] != 1 {
		t.Errorf("paused %v and resumed %v, want the tone's handle 1 in each", rec.paused, rec.resumed)
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
// handles, which stop, pause, resume and gain take as no-ops too, and bad arguments still
// fail, so native dev catches them.
func TestAudioNoopWithoutBackend(t *testing.T) {
	lg, err := api.NewLetGo("audio-noop-test")
	if err != nil {
		t.Fatal(err)
	}
	v, err := lg.Run(`(let [h (audio/tone {:freq 440 :ms 50})]
	                    [h (audio/gain h 0.5 100) (audio/pause h) (audio/resume h) (audio/stop h) (audio/stop 3) (audio/enabled?)])`)
	if err != nil {
		t.Fatal(err)
	}
	if got := v.String(); got != "[nil nil nil nil nil nil false]" {
		t.Errorf("got %s, want [nil nil nil nil nil nil false]", got)
	}
	for _, bad := range []string{
		`(audio/tone {:freq 0 :ms 50})`,
		`(audio/tone {:freq 440 :ms 50 :wave :kazoo})`,
		`(audio/sample (byte-array [1 2 3]))`,
		`(audio/sample (byte-array [1 0 2 0]) {:loop-from 1})`,
		`(audio/gain 1 -1)`,
		`(audio/gain nil -1)`,
		`(audio/tone {:freq 440 :ms 50 :gain -1})`,
		`(audio/sample (byte-array [-24 3 -24 3]) {:loop-from 0 :gain -1})`,
		`(audio/stop :h)`,
		`(audio/pause :h)`,
		`(audio/resume 1 2)`,
		`(audio/tone {:freq ##Inf :ms 50})`,
		`(audio/tone {:freq 440 :ms ##Inf})`,
		`(audio/tone {:freq 440 :ms 1e300})`,
		`(audio/sample (byte-array [1 0 2 0]) {:loop-from ##Inf})`,
		`(audio/sample (byte-array [1 0 2 0]) {:loop-from 1e300})`,
		`(audio/gain 1 0.5 ##Inf)`,
		`(audio/gain 1 0.5 ##NaN)`,
		`(audio/stream {:buffer-ms 0})`,
		`(audio/stream {:buffer-ms ##Inf})`,
		`(audio/stream {:buffer-ms 1e9})`,
		`(audio/feed! 1 (byte-array [1 2 3]))`,
		`(audio/feed! 1 [1 0 2 0])`,
		`(audio/feed! :h (byte-array 4))`,
		`(audio/feed! 1)`,
	} {
		if _, err := lg.Run(bad); err == nil {
			t.Errorf("%s: want an error", bad)
		}
	}
}

// With no backend, a stream is nil like any other handle, and feeding it is
// a no-op, so a guest that streams runs silent without a guard.
func TestAudioStreamNoopWithoutBackend(t *testing.T) {
	lg, err := api.NewLetGo("audio-stream-noop-test")
	if err != nil {
		t.Fatal(err)
	}
	v, err := lg.Run(`(let [h (audio/stream)] [h (audio/feed! h (byte-array 8)) (audio/stop h)])`)
	if err != nil {
		t.Fatal(err)
	}
	if got := v.String(); got != "[nil nil nil]" {
		t.Errorf("got %s, want [nil nil nil]", got)
	}
}

// pcmOf is n stereo frames of the constant sample x, the byte-array literal
// a guest would build.
func pcmOf(n int, x int16) string {
	lo, hi := int8(uint16(x)), int8(uint16(x)>>8)
	return fmt.Sprintf("(byte-array (take %d (cycle [%d %d %d %d])))", 4*n, lo, hi, lo, hi)
}

// leftAt reads frame i's left sample from mixer output.
func leftAt(p []byte, i int) int16 { return int16(binary.LittleEndian.Uint16(p[4*i:])) }

// PCM a guest feeds plays through the mixer; an underrun plays silence
// without ending the stream, and the next feed is heard.
func TestAudioStreamFeedsTheMixer(t *testing.T) {
	m := audio.NewMixer()
	lg, err := api.NewLetGo("audio-stream-mixer-test", api.WithAudio(m))
	if err != nil {
		t.Fatal(err)
	}
	h, err := lg.Run(fmt.Sprintf(`(let [h (audio/stream {:buffer-ms 50})] (audio/feed! h %s) h)`, pcmOf(400, 1000)))
	if err != nil {
		t.Fatal(err)
	}
	p := make([]byte, 4*600)
	m.Read(p)
	// The first 220 frames carry the mixer's fade-in, so look past it.
	if got := leftAt(p, 300); got != 1000 {
		t.Errorf("frame 300 = %d, want the fed 1000", got)
	}
	if got := leftAt(p, 500); got != 0 {
		t.Errorf("frame 500 = %d, want silence once the feed ran out", got)
	}
	// A read that finds nothing fed at all must not end the stream either.
	m.Read(p)
	if _, err := lg.Run(fmt.Sprintf(`(audio/feed! %s %s)`, h, pcmOf(100, -500))); err != nil {
		t.Fatal(err)
	}
	m.Read(p)
	if got := leftAt(p, 50); got != -500 {
		t.Errorf("after the underrun, frame 50 = %d, want the newly fed -500", got)
	}
}

// feed! waits while the stream's buffer is full, so the mixer's reads pace
// the guest; stopping the stream releases a feed! that is waiting.
func TestAudioStreamFeedWaitsAndStopReleases(t *testing.T) {
	m := audio.NewMixer()
	lg, err := api.NewLetGo("audio-stream-pace-test", api.WithAudio(m))
	if err != nil {
		t.Fatal(err)
	}
	h, err := lg.Run(`(audio/stream {:buffer-ms 1})`)
	if err != nil {
		t.Fatal(err)
	}
	feed := func() chan error {
		done := make(chan error, 1)
		go func() {
			_, err := lg.Run(fmt.Sprintf(`(audio/feed! %s %s)`, h, pcmOf(1000, 7)))
			done <- err
		}()
		return done
	}

	done := feed()
	select {
	case <-done:
		t.Fatal("feeding 1000 frames into a 44-frame buffer returned before the mixer read any")
	case <-time.After(50 * time.Millisecond):
	}
	p := make([]byte, 4*256)
	deadline := time.After(2 * time.Second)
drain:
	for {
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
			break drain
		case <-deadline:
			t.Fatal("feed! still waiting while the mixer kept reading")
		default:
			m.Read(p)
			time.Sleep(time.Millisecond)
		}
	}

	// A second instance on the same mixer stops it while this feed! waits.
	done = feed()
	time.Sleep(20 * time.Millisecond)
	other, err := api.NewLetGo("audio-stream-stop-test", api.WithAudio(m))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := other.Run(fmt.Sprintf(`(audio/stop %s)`, h)); err != nil {
		t.Fatal(err)
	}
	waitFeed(t, m, done, "audio/stop did not release the waiting feed!")

	// The stopped stream's handle no longer reaches a ring: feed! drops.
	done = feed()
	waitFeed(t, m, done, "feed! on a stopped stream waited")
}

// waitFeed reads the mixer, as a device would, until done reports or a
// second passes, so a feed! that is paced by reads can finish.
func waitFeed(t *testing.T, m *audio.Mixer, done chan error, msg string) {
	t.Helper()
	p := make([]byte, 4*256)
	deadline := time.After(time.Second)
	for {
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
			return
		case <-deadline:
			t.Fatal(msg)
		default:
			m.Read(p)
			time.Sleep(time.Millisecond)
		}
	}
}

// A guest that pauses its stream and keeps feeding from the same loop must
// not hang: nothing drains a paused stream, so feed! drops what doesn't fit.
// A stop from Go, not through audio/stop, releases a waiting feed! too.
func TestAudioStreamPauseAndGoStopDontHang(t *testing.T) {
	m := audio.NewMixer()
	lg, err := api.NewLetGo("audio-stream-pause-test", api.WithAudio(m))
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := lg.Run(fmt.Sprintf(`(let [h (audio/stream {:buffer-ms 1})]
		                                 (audio/pause h)
		                                 (audio/feed! h %s)
		                                 (audio/resume h))`, pcmOf(1000, 7)))
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("feed! on a paused stream hung")
	}

	h, err := lg.Run(`(audio/stream {:buffer-ms 1})`)
	if err != nil {
		t.Fatal(err)
	}
	done = make(chan error, 1)
	go func() {
		_, err := lg.Run(fmt.Sprintf(`(audio/feed! %s %s)`, h, pcmOf(1000, 7)))
		done <- err
	}()
	time.Sleep(20 * time.Millisecond)
	m.Stop(rt.AudioHandle(h.(vm.Int)))
	waitFeed(t, m, done, "a stop from Go did not release the waiting feed!")
}

// A host that stops reading (a lost device, a mixer nobody pumps) must not
// hang a guest: a stop still releases a waiting feed! without any read, and a
// feed! into a stream nobody drains gives up once the stall timeout passes.
func TestAudioStreamUnreadDoesntHang(t *testing.T) {
	m := audio.NewMixer()
	lg, err := api.NewLetGo("audio-stream-unread-test", api.WithAudio(m))
	if err != nil {
		t.Fatal(err)
	}
	feed := func(h any) chan error {
		done := make(chan error, 1)
		go func() {
			_, err := lg.Run(fmt.Sprintf(`(audio/feed! %s %s)`, h, pcmOf(1000, 7)))
			done <- err
		}()
		return done
	}
	wait := func(done chan error, within time.Duration, msg string) {
		t.Helper()
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(within):
			t.Fatal(msg)
		}
	}

	h, err := lg.Run(`(audio/stream {:buffer-ms 1})`)
	if err != nil {
		t.Fatal(err)
	}
	done := feed(h)
	time.Sleep(20 * time.Millisecond)
	m.Stop(rt.AudioHandle(h.(vm.Int)))
	// Well inside the 1 s stall floor, so the stop is what released it.
	wait(done, 300*time.Millisecond, "a stop with nothing reading did not release the waiting feed!")

	h, err = lg.Run(`(audio/stream {:buffer-ms 1})`)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	done = feed(h)
	select {
	case <-done:
		t.Fatal("feed! into an unread stream returned before the stall timeout")
	case <-time.After(500 * time.Millisecond):
	}
	wait(done, 2*time.Second, "feed! into a stream nobody reads never gave up")
	if d := time.Since(start); d < time.Second {
		t.Errorf("feed! gave up after %v, want at least the 1 s stall floor", d)
	}
}
