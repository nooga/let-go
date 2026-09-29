/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/zeebo/bencode"
)

// fakePodHelperEnv marks the process that plays the pod: TestMain hands a
// test binary started with it set to runFakePod before any test setup. The
// pod runner starts a pod with no arguments, so the pod is a shell wrapper
// that re-executes this test binary with the marker inherited from the
// test's environment.
const fakePodHelperEnv = "LG_TEST_FAKE_POD_HELPER"

// TestPodEDNReplies runs fixtures/pods/edn_readers_test.lg against a fake EDN
// pod (runFakePod) whose describe reply declares a reader implemented by the
// pod's own client-side code.
func TestPodEDNReplies(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	wrapper := filepath.Join(t.TempDir(), "fake-pod")
	script := fmt.Sprintf("#!/bin/sh\nexec %q\n", exe)
	if err := os.WriteFile(wrapper, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(fakePodHelperEnv, "1")
	t.Setenv("LG_TEST_FAKE_POD", wrapper)

	ok, msg := runFileTests("fixtures/pods/edn_readers_test.lg")
	if !ok {
		t.Fatalf("pod EDN fixture failed: %s", msg)
	}
}

// runFakePod is the pod. It speaks the babashka pod protocol (bencode over
// stdio) and exits the process on shutdown or end of input, so nothing but
// protocol ever reaches the wire.
func runFakePod() {
	out := bufio.NewWriter(os.Stdout)
	reply := func(msg map[string]any) {
		bs, err := bencode.EncodeBytes(msg)
		if err != nil {
			panic(err)
		}
		if _, err := out.Write(bs); err != nil {
			panic(err)
		}
		if err := out.Flush(); err != nil {
			panic(err)
		}
	}
	dec := bencode.NewDecoder(os.Stdin)
	for {
		var msg map[string]any
		if err := dec.Decode(&msg); err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			panic(err)
		}
		id, _ := msg["id"].(string)
		switch msg["op"] {
		case "describe":
			reply(map[string]any{
				"format":  "edn",
				"ops":     map[string]any{"shutdown": map[string]any{}},
				"readers": map[string]any{"my/tag": "pod.test/read-tag"},
				"namespaces": []any{map[string]any{
					"name": "pod.test",
					"vars": []any{
						map[string]any{"name": "read-tag", "code": "(defn read-tag [x] [:tagged x])"},
						map[string]any{"name": "tagged"},
						map[string]any{"name": "unknown"},
						map[string]any{"name": "bad-ex-data"},
						map[string]any{"name": "stream"},
						map[string]any{"name": "stream-bad"},
					},
				}},
			})
		case "invoke":
			switch msg["var"] {
			case "pod.test/tagged":
				reply(map[string]any{"id": id, "value": "#my/tag 42", "status": []any{"done"}})
			case "pod.test/unknown":
				reply(map[string]any{"id": id, "value": "#nope/x 1", "status": []any{"done"}})
			case "pod.test/bad-ex-data":
				reply(map[string]any{"id": id, "ex-message": "boom", "ex-data": "#nope/x 1", "status": []any{"done", "error"}})
			case "pod.test/stream":
				reply(map[string]any{"id": id, "value": "1"})
				reply(map[string]any{"id": id, "value": "#my/tag 2"})
				reply(map[string]any{"id": id, "status": []any{"done"}})
			case "pod.test/stream-bad":
				reply(map[string]any{"id": id, "value": "1"})
				reply(map[string]any{"id": id, "value": "#nope/x 2"})
				reply(map[string]any{"id": id, "status": []any{"done"}})
			default:
				reply(map[string]any{"id": id, "ex-message": fmt.Sprintf("unknown var %v", msg["var"]), "status": []any{"done", "error"}})
			}
		case "shutdown":
			os.Exit(0)
		}
	}
}
