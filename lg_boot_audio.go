//go:build audio_oto

package main

// The oto sound backend is opt-in: build with -tags audio_oto. Importing it
// binds the *audio* root, so the audio namespace plays; without it every
// audio/ call is a silent no-op.
import _ "github.com/nooga/let-go/pkg/audio"
