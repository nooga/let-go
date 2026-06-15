/*
 * Copyright (c) 2022-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package nrepl

import (
	"net"
	"os"
	"strconv"
	"testing"
)

// With -p 0 the OS picks a free ephemeral port. The server must report the
// bound port (not the requested 0) via Port(), the banner, and .nrepl-port,
// otherwise no client can discover where it's listening.
func TestStartResolvesEphemeralPort(t *testing.T) {
	t.Chdir(t.TempDir()) // Start writes .nrepl-port in cwd

	n := NewNreplServer(nil)
	if err := n.Start(0); err != nil {
		t.Fatalf("Start(0): %v", err)
	}
	defer n.Stop()

	bound := n.listener.Addr().(*net.TCPAddr).Port
	if bound == 0 {
		t.Fatal("listener bound to port 0")
	}
	if n.Port() != bound {
		t.Errorf("Port() = %d, want bound port %d", n.Port(), bound)
	}

	data, err := os.ReadFile(".nrepl-port")
	if err != nil {
		t.Fatalf("read .nrepl-port: %v", err)
	}
	if got := string(data); got != strconv.Itoa(bound) {
		t.Errorf(".nrepl-port = %q, want %q", got, strconv.Itoa(bound))
	}
}

// A fixed port must still round-trip unchanged.
func TestStartFixedPort(t *testing.T) {
	t.Chdir(t.TempDir())

	// Grab a free port, release it, then ask the server for it explicitly.
	probe, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("probe listen: %v", err)
	}
	want := probe.Addr().(*net.TCPAddr).Port
	probe.Close()

	n := NewNreplServer(nil)
	if err := n.Start(want); err != nil {
		t.Fatalf("Start(%d): %v", want, err)
	}
	defer n.Stop()

	if n.Port() != want {
		t.Errorf("Port() = %d, want %d", n.Port(), want)
	}
}
