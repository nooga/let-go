//go:build plan9 || js || wasip1

/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/vm"
)

// repl is a minimal line-by-line REPL for platforms without readline. The
// chzyer/readline library depends on termios/ioctl, unavailable on plan9 and
// js/wasm, so these builds read from stdin via bufio.Scanner — no completion,
// no syntax highlighting, no in-line editing. (On js the root binary isn't the
// REPL entry at all; this just keeps GOOS=js go build ./... compiling.)
func repl(ctx *compiler.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	prompt := ctx.CurrentNS().Name() + "=> "
	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil && err != io.EOF {
				fmt.Println("prompt failed:", err)
			}
			fmt.Println()
			return
		}
		in := strings.TrimRight(scanner.Text(), "\r\n")
		if in == "" {
			continue
		}
		ctx.SetSource("REPL")
		val, err := runForm(ctx, in)
		if err != nil {
			fmt.Print(vm.FormatError(err))
		} else {
			fmt.Println(val.String())
		}
		prompt = ctx.CurrentNS().Name() + "=> "
	}
}
