//go:build plan9

/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package cli

// rio doesn't interpret ANSI escapes — strip them from the banner.
// Also adjust the quit hint: plan9's rc has no SIGINT, so Ctrl-C does nothing
// — interrupt with Delete or send EOF with Ctrl-D.
const (
	ansiBold     = ""
	ansiBoldCyan = ""
	ansiDim      = ""
	ansiReset    = ""

	bannerQuitHint = "Delete or Ctrl-D to quit"
)
