/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func tempFile(t *testing.T, name string) *os.File {
	t.Helper()
	f, err := os.Create(filepath.Join(t.TempDir(), name))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })
	return f
}

// TestIOHandleProcessWriter: ProcessWriter yields an *os.File for file-backed
// handles (so os/exec passes the descriptor through) and the raw writer for
// everything else (so captures still get a pipe).
func TestIOHandleProcessWriter(t *testing.T) {
	t.Run("std stream handle resolves the current file", func(t *testing.T) {
		first, second := tempFile(t, "first"), tempFile(t, "second")
		cur := first
		h := newStdStreamHandle(first, func() *os.File { return cur })
		if got := h.ProcessWriter(); got != first {
			t.Fatalf("ProcessWriter = %T %v, want the first file", got, got)
		}
		cur = second
		if got := h.ProcessWriter(); got != second {
			t.Fatalf("ProcessWriter after swap = %T %v, want the second file", got, got)
		}
	})
	t.Run("plain file handle yields its file", func(t *testing.T) {
		f := tempFile(t, "plain")
		if got := NewIOHandle(f).ProcessWriter(); got != f {
			t.Fatalf("ProcessWriter = %T %v, want the file", got, got)
		}
	})
	t.Run("writer handle yields the writer itself", func(t *testing.T) {
		buf := &bytes.Buffer{}
		if got := NewWriterHandle("buf", buf).ProcessWriter(); got != buf {
			t.Fatalf("ProcessWriter = %T %v, want the buffer", got, got)
		}
	})
	t.Run("reader-only handle yields nil", func(t *testing.T) {
		if got := NewReaderHandle("r", strings.NewReader("")).ProcessWriter(); got != nil {
			t.Fatalf("ProcessWriter = %T %v, want nil", got, got)
		}
	})
}

// TestExecStarPassesStdoutDescriptor: with *out*/*err* at their root
// bindings, the child of os/exec* writes to the CURRENT os.Stdout/os.Stderr
// descriptors themselves, not to a pipe copied into them. /proc/self/fd
// makes that observable: the child reads back its own fd targets.
func TestExecStarPassesStdoutDescriptor(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("relies on /proc/self/fd")
	}
	outF, errF := tempFile(t, "stdout"), tempFile(t, "stderr")
	savedOut, savedErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = outF, errF
	t.Cleanup(func() { os.Stdout, os.Stderr = savedOut, savedErr })

	fn := LookupVar("os", "exec*").Deref().(vm.Fn)
	args := []vm.Value{vm.String("sh"), vm.String("-c"),
		vm.String("readlink /proc/self/fd/1; readlink /proc/self/fd/2 >&2")}
	code, err := vm.NewExecContext().Invoke(fn, args)
	if err != nil {
		t.Fatalf("os/exec*: %v", err)
	}
	if code != vm.Int(0) {
		t.Fatalf("os/exec* exit = %v, want 0", code)
	}
	for _, c := range []struct {
		stream string
		f      *os.File
	}{{"stdout", outF}, {"stderr", errF}} {
		data, err := os.ReadFile(c.f.Name())
		if err != nil {
			t.Fatal(err)
		}
		got := strings.TrimSpace(string(data))
		if got != c.f.Name() {
			t.Errorf("child %s fd -> %q, want %q (a pipe means the descriptor was not passed through)", c.stream, got, c.f.Name())
		}
	}
}
