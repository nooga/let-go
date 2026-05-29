//go:build !linux && !js

/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"
)

// On non-Linux hosts we don't implement pty opening — lgcr targets Linux and
// the term ns can still be used for terminal UI (raw-mode, size, colors).
const tiocswinsz = 0x5414 // unused; placeholder so term.go compiles
// FIONREAD on Darwin/BSD: _IOR('f', 127, int). Same value across Darwin,
// FreeBSD, NetBSD, OpenBSD. Used by key-pending? to peek the stdin
// buffer without consuming.
const fionread = 0x4004667F

func openPtyPair() (*os.File, *os.File, string, error) {
	return nil, nil, "", fmt.Errorf("open-pty is only supported on Linux")
}
