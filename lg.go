/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

// Command lg is the let-go command line. Everything it does lives in
// github.com/nooga/let-go/pkg/cli so a third-party module can build a binary
// with the same full feature set — repl, resolver, -c, -b, -w — around its own
// generated interop packages. See docs/guide/custom-lg.md.
//
// Only the ldflags contract stays here: goreleaser and the Makefile both set
// -X main.version / -X main.commit, and ldflags can only reach the package
// that declares the variable.
package main

import (
	"os"

	"github.com/nooga/let-go/pkg/cli"
)

// Set by goreleaser via ldflags
var (
	version = "dev"
	commit  = "none"
)

func main() {
	// Exit only on failure, as this has always done: a zero code returns from
	// main normally instead of going through os.Exit.
	//
	// bootMain is the build-tagged boot seam: identity by default, and under
	// -tags glplat_ebiten a trampoline that hands the main thread to ebiten,
	// which demands it, while the lg program runs on a goroutine. It wraps
	// cli.Main from here rather than living inside pkg/cli so the importable
	// CLI keeps no graphics dependency.
	if code := bootMain(func() int { return cli.Main(version, commit) }); code != 0 {
		os.Exit(code)
	}
}
