//go:build !plan9

/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

const (
	ansiBold     = "\x1b[1m"
	ansiBoldCyan = "\x1b[1;36m"
	ansiDim      = "\x1b[90m"
	ansiReset    = "\x1b[0m"

	bannerQuitHint = "Ctrl-C to quit"
)
