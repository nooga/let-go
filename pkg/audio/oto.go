//go:build audio_oto && !js

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package audio

import (
	"github.com/ebitengine/oto/v3"

	"github.com/nooga/let-go/pkg/rt"
)

// start plays the mixer through oto. The queue is l's whether or not it was
// asked for: the native default is the measured one.
func (d *Device) start(l latency, _ bool) (any, error) {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   rt.AudioRate,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
		BufferSize:   deviceBuffer(l.frames),
	})
	if err != nil {
		return nil, err
	}
	<-ready
	// oto opens the device asynchronously and reports a failure here.
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	player := ctx.NewPlayer(d.Mixer)
	player.SetBufferSize(playerBytes(l.queue))
	player.Play()
	// Returned to be kept in Device.backend: oto closes an unreachable player
	// mid-playback (oto#293).
	return player, nil
}
