/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package nrepl

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestStartWithPortZeroReportsBoundPort(t *testing.T) {
	t.Chdir(t.TempDir())

	srv := NewNreplServer(nil)
	output := captureStdout(t, func() {
		if err := srv.Start(0); err != nil {
			t.Fatalf("Start(0) failed: %v", err)
		}
	})
	t.Cleanup(srv.Stop)

	if srv.Port() == 0 {
		t.Fatal("Port() returned 0 after Start(0)")
	}

	portFile, err := os.ReadFile(".nrepl-port")
	if err != nil {
		t.Fatalf("reading .nrepl-port: %v", err)
	}
	if got, want := strings.TrimSpace(string(portFile)), strconv.Itoa(srv.Port()); got != want {
		t.Fatalf(".nrepl-port = %q, want %q", got, want)
	}

	if strings.Contains(output, "port 0") || strings.Contains(output, ":0") {
		t.Fatalf("startup output contains port 0: %q", output)
	}
	if !strings.Contains(output, strconv.Itoa(srv.Port())) {
		t.Fatalf("startup output %q does not contain bound port %d", output, srv.Port())
	}
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("creating stdout pipe: %v", err)
	}

	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()

	f()

	if err := w.Close(); err != nil {
		t.Fatalf("closing stdout pipe: %v", err)
	}

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("reading stdout pipe: %v", err)
	}
	return string(out)
}
