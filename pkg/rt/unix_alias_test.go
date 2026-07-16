//go:build linux

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// unixConnPair returns two connected *net.UnixConn via socketpair(2).
func unixConnPair(t *testing.T) (*net.UnixConn, *net.UnixConn) {
	t.Helper()
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatalf("socketpair: %v", err)
	}
	mk := func(fd int) *net.UnixConn {
		f := os.NewFile(uintptr(fd), "unixpair")
		c, err := net.FileConn(f)
		_ = f.Close() // FileConn dups the fd; drop ours
		if err != nil {
			t.Fatalf("FileConn: %v", err)
		}
		uc, ok := c.(*net.UnixConn)
		if !ok {
			t.Fatalf("expected *net.UnixConn, got %T", c)
		}
		return uc
	}
	return mk(fds[0]), mk(fds[1])
}

// TestUnixBangAliases checks that the canonical bang names (#526) resolve, that
// they carry data end to end, and that the deprecated pre-bang names still
// resolve and delegate.
func TestUnixBangAliases(t *testing.T) {
	for _, name := range []string{"write!", "read!", "close!", "send", "recv", "close"} {
		_ = nsFn(t, "unix", name) // fails the test if unregistered or not an Fn
	}

	a, b := unixConnPair(t)
	defer func() { _ = a.Close() }()
	defer func() { _ = b.Close() }()
	av, bv := vm.NewBoxed(a), vm.NewBoxed(b)

	// write! on one end → read! on the other returns the bytes under :data.
	if _, err := nsFn(t, "unix", "write!").Invoke([]vm.Value{av, vm.String("ping")}); err != nil {
		t.Fatalf("unix/write!: %v", err)
	}
	res, err := nsFn(t, "unix", "read!").Invoke([]vm.Value{bv, vm.MakeInt(16), vm.MakeInt(0)})
	if err != nil {
		t.Fatalf("unix/read!: %v", err)
	}
	m, ok := res.(vm.Lookup)
	if !ok {
		t.Fatalf("unix/read! returned %T, want a map", res)
	}
	if got := m.ValueAt(vm.Keyword("data")); got != vm.String("ping") {
		t.Fatalf("unix/read! :data = %#v, want \"ping\"", got)
	}

	// The deprecated alias still delegates (and prints a one-time notice).
	if _, err := nsFn(t, "unix", "send").Invoke([]vm.Value{av, vm.String("x")}); err != nil {
		t.Fatalf("unix/send (deprecated alias): %v", err)
	}
}
