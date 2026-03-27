/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"strings"
)

// compileErrorLike is satisfied by compiler.CompileError without importing it.
type compileErrorLike interface {
	error
	InnermostMessage() string
	InnermostSource() *SourceInfo
}

// FormatError produces a user-friendly error display with source snippets,
// inspired by Rust/Elm-style error reporting.
func FormatError(err error) string {
	// Handle compile errors
	if ce, ok := err.(compileErrorLike); ok {
		return formatCompileError(ce)
	}

	var b strings.Builder

	// Collect frames from the error chain
	type frame struct {
		msg    string
		source *SourceInfo
		fnName string
	}

	var frames []frame
	current := err
	for current != nil {
		switch e := current.(type) {
		case *ExecutionError:
			name := strings.TrimPrefix(e.message, "calling ")
			frames = append(frames, frame{msg: e.message, source: e.source, fnName: name})
			current = e.cause
		case *TypeError:
			frames = append(frames, frame{msg: current.Error()})
			current = e.cause
		default:
			frames = append(frames, frame{msg: current.Error()})
			current = nil
		}
	}

	if len(frames) == 0 {
		return err.Error()
	}

	// Root cause is the last frame
	root := frames[len(frames)-1]

	// Error header
	fmt.Fprintf(&b, "\x1b[1;31merror:\x1b[0m %s\n", root.msg)

	// Source snippet for the deepest frame that has source info
	for i := len(frames) - 1; i >= 0; i-- {
		if frames[i].source != nil {
			writeSnippet(&b, frames[i].source)
			break
		}
	}

	// Stack trace (if more than one frame with a name)
	hasTrace := false
	for i := 0; i < len(frames)-1; i++ {
		if frames[i].fnName != "" {
			hasTrace = true
			break
		}
	}
	if hasTrace {
		b.WriteString("\n\x1b[1mstack trace:\x1b[0m\n")
		for i := len(frames) - 2; i >= 0; i-- {
			f := frames[i]
			loc := "<unknown>"
			if f.source != nil {
				loc = f.source.String()
			}
			fmt.Fprintf(&b, "  at %s (%s)\n", f.fnName, loc)
		}
	}

	return b.String()
}

func formatCompileError(ce compileErrorLike) string {
	var b strings.Builder
	msg := ce.InnermostMessage()
	src := ce.InnermostSource()

	fmt.Fprintf(&b, "\x1b[1;31merror:\x1b[0m %s\n", msg)
	if src != nil {
		writeSnippet(&b, src)
	}
	return b.String()
}

func writeSnippet(b *strings.Builder, info *SourceInfo) {
	line := SourceRegistry.GetLine(info.File, info.Line)
	if line == "" {
		fmt.Fprintf(b, "  \x1b[1;34m-->\x1b[0m %s\n", info.String())
		return
	}

	lineNum := info.Line + 1
	width := len(fmt.Sprintf("%d", lineNum))
	padding := strings.Repeat(" ", width)

	fmt.Fprintf(b, "  \x1b[1;34m-->\x1b[0m %s\n", info.String())
	fmt.Fprintf(b, " %s \x1b[1;34m|\x1b[0m\n", padding)
	fmt.Fprintf(b, " \x1b[1;34m%d\x1b[0m \x1b[1;34m|\x1b[0m %s\n", lineNum, line)

	// Position indicator
	col := info.Column
	if col < 0 {
		col = 0
	}
	pointer := strings.Repeat(" ", col) + "\x1b[1;31m^^^\x1b[0m"
	fmt.Fprintf(b, " %s \x1b[1;34m|\x1b[0m %s\n", padding, pointer)
}
