//go:build audio_oto && !js

package audio

import (
	"io"
	"os"
	"strings"
	"testing"
)

// os.Pipe is not implemented on js/wasm, so this stays on the oto build.
// The bad override fails in resolve, before oto is touched, so Enabled here
// opens no sound device either. It swaps the process-wide os.Stderr, so it
// must never call t.Parallel.
func TestABadOverrideFailsTheOpenAndSaysWhy(t *testing.T) {
	clearEnv(t)
	t.Setenv(envQueue, "10ms")
	d := &Device{Mixer: NewMixer()}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr := os.Stderr
	os.Stderr = w
	enabled := d.Enabled()
	os.Stderr = stderr
	w.Close()
	said, _ := io.ReadAll(r)

	if enabled {
		t.Fatal("audio enabled under a bad override")
	}
	if !strings.Contains(string(said), "shorter than") {
		t.Fatalf("stderr %q does not give the reason", said)
	}
	// The failed open is final: a later Configure is refused.
	if err := d.Configure(Options{}); err == nil {
		t.Fatal("configured a device whose open already failed")
	}
}
