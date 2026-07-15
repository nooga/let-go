//go:build !plan9 && !js && !wasip1

/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/chzyer/readline"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func isCompletionTerminator(r rune) bool {
	switch r {
	case '(', ')', '[', ']', '{', '}', '"', '\\', '\'', '@', '`', '~', ';', '#':
		return true
	}
	return unicode.IsSpace(r)
}

// completer implements readline's AutoCompleter interface.
// readline passes line as []rune and pos as a rune index; we stay in runes
// throughout so non-ASCII input and mid-line cursors are handled correctly.
type completer struct {
	ctx *compiler.Context
}

func (c *completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	if pos > len(line) {
		pos = len(line)
	}
	head := line[:pos]

	start := pos
	for start > 0 && !isCompletionTerminator(head[start-1]) {
		start--
	}
	prefix := string(head[start:])

	symbols := rt.FuzzyNamespacedSymbolLookup(c.ctx.CurrentNS(), vm.Symbol(prefix))
	for _, s := range symbols {
		newLine = append(newLine, []rune(string(s)+" "))
	}
	length = pos - start
	return
}

// ANSI color codes for syntax highlighting. ansiReset and ansiBold are
// declared in lg_ansi.go; the rest are repl-only.
const (
	ansiMagenta = "\x1b[35m"
	ansiYellow  = "\x1b[33m"
	ansiBlue    = "\x1b[34m"
	ansiCyan    = "\x1b[36m"
)

var tokenColors = map[compiler.TokenKind]string{
	compiler.TokenNumber:      ansiMagenta,
	compiler.TokenPunctuation: ansiYellow,
	compiler.TokenKeyword:     ansiBlue,
	compiler.TokenString:      ansiCyan,
	compiler.TokenSpecial:     ansiBold,
}

// syntaxHighlighter paints tokens with ANSI escape codes. Token positions
// from LispReader are rune offsets (reader bumps pos per ReadRune), so we
// slice the rune buffer directly rather than the UTF-8 byte string.
type syntaxHighlighter struct{}

func (s *syntaxHighlighter) Paint(line []rune, _ int) []rune {
	if len(line) == 0 {
		return line
	}

	reader := compiler.NewLispReaderTokenizing(strings.NewReader(string(line)), "syntax")
	reader.Read() //nolint:errcheck // partial parse is expected mid-edit

	var out strings.Builder
	out.Grow(len(line) + 32)
	cursor := 0
	for _, t := range reader.Tokens {
		if t.End <= t.Start || t.End > len(line) {
			continue
		}
		color, ok := tokenColors[t.Kind]
		if !ok {
			continue
		}
		if t.Start > cursor {
			out.WriteString(string(line[cursor:t.Start]))
		}
		out.WriteString(color)
		out.WriteString(string(line[t.Start:t.End]))
		out.WriteString(ansiReset)
		cursor = t.End
	}
	if cursor < len(line) {
		out.WriteString(string(line[cursor:]))
	}
	return []rune(out.String())
}

// historyFile returns a per-user history path, falling back to disabling
// history if no suitable directory exists.
func historyFile() string {
	if dir, err := os.UserHomeDir(); err == nil && dir != "" {
		return filepath.Join(dir, ".lg_history")
	}
	return ""
}

func repl(ctx *compiler.Context) {
	prompt := ctx.CurrentNS().Name() + "=> "

	// Configure readline
	config := &readline.Config{
		Prompt:          prompt,
		EOFPrompt:       "exit",
		AutoComplete:    &completer{ctx: ctx},
		InterruptPrompt: "^C",
		HistoryFile:     historyFile(),
		Painter:         &syntaxHighlighter{},
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		fmt.Println("failed to initialize readline:", err)
		return
	}
	defer rl.Close()

	for {
		// Update prompt in case namespace changed
		prompt = ctx.CurrentNS().Name() + "=> "
		rl.SetPrompt(prompt)

		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt, etc.
			if err == readline.ErrInterrupt {
				// Ctrl-C on an empty line quits (matches banner hint);
				// otherwise just discard the in-progress line.
				if line == "" {
					break
				}
				continue
			}
			fmt.Println()
			break
		}

		if line == "" {
			continue
		}

		ctx.SetSource("REPL")
		val, err := runForm(ctx, line)
		if err != nil {
			fmt.Print(vm.FormatError(err))
		} else {
			fmt.Println(val.String())
		}
	}
}
