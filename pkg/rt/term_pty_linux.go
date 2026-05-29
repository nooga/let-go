//go:build linux

/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	tiocswinsz = 0x5414
	tiocgptn   = 0x80045430
	tiocsptlck = 0x40045431
	fionread   = 0x541B
)

func openPtyPair() (master *os.File, slave *os.File, slavePath string, err error) {
	m, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, "", fmt.Errorf("open ptmx: %w", err)
	}
	var zero int32
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		m.Fd(), uintptr(tiocsptlck), uintptr(unsafe.Pointer(&zero))); errno != 0 {
		m.Close()
		return nil, nil, "", fmt.Errorf("unlockpt: %v", errno)
	}
	var n int32
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		m.Fd(), uintptr(tiocgptn), uintptr(unsafe.Pointer(&n))); errno != 0 {
		m.Close()
		return nil, nil, "", fmt.Errorf("ptsname: %v", errno)
	}
	slavePath = fmt.Sprintf("/dev/pts/%d", n)
	s, err := os.OpenFile(slavePath, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		m.Close()
		return nil, nil, "", fmt.Errorf("open slave: %w", err)
	}
	return m, s, slavePath, nil
}
