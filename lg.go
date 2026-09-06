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
	if code := cli.Main(version, commit); code != 0 {
		os.Exit(code)
	}
}
