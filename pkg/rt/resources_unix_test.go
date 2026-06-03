//go:build unix

/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"path/filepath"
	"syscall"
	"testing"
)

// TestFSResourceProviderRejectsFIFO: a non-regular file (FIFO) must not be
// exposed as a resource — opening it would block io/slurp on the pipe. Lives
// in a unix-tagged file because syscall.Mkfifo is not available on all
// supported platforms (e.g. Plan 9, Windows).
func TestFSResourceProviderRejectsFIFO(t *testing.T) {
	root := t.TempDir()
	fifo := filepath.Join(root, "pipe")
	if err := syscall.Mkfifo(fifo, 0644); err != nil {
		t.Skipf("mkfifo unavailable: %v", err)
	}
	p := NewFSResourceProvider([]string{root})
	if _, ok := p.Open("pipe"); ok {
		t.Errorf("fifo: expected not found")
	}
}
