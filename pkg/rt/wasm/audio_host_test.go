/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package wasm

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestAudioWorkletFailureStaysFailedAfterClose(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node not on PATH; cannot run the browser audio host")
	}

	_, audio, found := strings.Cut(lgHostCoreJS, "// --- Audio (")
	if !found {
		t.Fatal("audio section missing from lg-host-core.js")
	}
	_, audio, found = strings.Cut(audio, "const LG_AUDIO_WORKLET")
	if !found {
		t.Fatal("audio worklet missing from lg-host-core.js")
	}
	audio, _, found = strings.Cut(audio, "// --- LetGoHost request/eval")
	if !found {
		t.Fatal("audio section end missing from lg-host-core.js")
	}

	const driver = `
globalThis.document = { addEventListener() {} };
let closes = 0;
globalThis.window = { AudioContext: class {
  state = 'suspended';
  audioWorklet = { addModule: () => Promise.reject(new Error('blocked worklet')) };
  close() {
    closes++;
    this.state = 'closed';
    setImmediate(() => this.onstatechange?.());
    return Promise.resolve();
  }
}};
const ctrl = new SharedArrayBuffer(16);
const shared = lgAudio.open({ sampleRate: 44100, cap: 512, data: new SharedArrayBuffer(4096), ctrl });
const message = lgAudio.open({ sampleRate: 44100, cap: 512 });
setImmediate(() => setImmediate(() => {
  const workerState = Atomics.load(new Int32Array(ctrl), 3);
  if (closes !== 2 || workerState !== 3 || shared.state !== 3 || message.state !== 3) {
    throw new Error('closes and audio states: ' + [closes, workerState, shared.state, message.state]);
  }
}));
`
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node)
	cmd.Stdin = strings.NewReader("const LG_AUDIO_WORKLET" + audio + driver)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("node: %v\n%s", err, out)
	}
}
