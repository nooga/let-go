/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"os"
	"strings"
	"testing"
)

// Guards the hand-added build constraint on the generated interop_xxh3.go.
// lginterop doesn't emit build tags, so regenerating the bindings drops the
// `!tinygo` gate. The tinygo-wasi CI job does catch that, but as a runtime
// trap at namespace install with nothing pointing at the build tag; this
// fails in the stock suite, at regeneration time, with the cause named.
// Remove this once lginterop can emit the constraint itself.
func TestInteropXxh3KeepsTinygoBuildTag(t *testing.T) {
	src, err := os.ReadFile("interop_xxh3.go")
	if err != nil {
		t.Fatalf("reading interop_xxh3.go: %v", err)
	}
	firstLine, _, _ := strings.Cut(string(src), "\n")
	if strings.TrimSpace(firstLine) != "//go:build !tinygo" {
		t.Fatalf("interop_xxh3.go must start with `//go:build !tinygo` (got %q).\n"+
			"The tag is hand-added — lginterop doesn't emit build constraints — and is\n"+
			"required because the generated bindings reflect-box typed Go funcs, which\n"+
			"TinyGo can't call. Re-add it after regenerating (see the NOTE in that file).",
			firstLine)
	}
}
