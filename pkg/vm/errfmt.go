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
	var goStack string
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
		case *GoPanicError:
			// An unexpected Go panic (e.g. syscall/js under WASM, nil deref in
			// a builtin). Carries the Go stack — the only traceback that pins
			// the actual crash site, since native fns have no .lg source.
			goStack = e.GoStack()
			frames = append(frames, frame{msg: e.Error()})
			current = nil
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
	fmt.Fprintf(&b, ansiBoldRed+"error:"+ansiReset+" %s\n", root.msg)

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
		b.WriteString("\n" + ansiBold + "stack trace:" + ansiReset + "\n")
		for i := len(frames) - 2; i >= 0; i-- {
			f := frames[i]
			loc := "<unknown>"
			if f.source != nil {
				loc = f.source.String()
			}
			fmt.Fprintf(&b, "  at %s (%s)\n", f.fnName, loc)
		}
	}

	// For an unexpected Go panic, surface the let-go-relevant Go frames so the
	// real crash site (the .go file:line) is visible — otherwise the only
	// traceback is the source-less let-go frames above.
	if origin := letGoStackFrames(goStack); origin != "" {
		b.WriteString("\n" + ansiBold + "go panic origin:" + ansiReset + "\n")
		b.WriteString(origin)
	}

	return b.String()
}

// letGoStackFrames extracts the let-go-relevant func+location pairs from a
// debug.Stack() dump (those inside this repo), skipping the panic-recover
// machinery, so a Go panic's actual origin (e.g. term_wasm.go:91) is shown.
// Returns "" when there are no in-repo frames (or no stack at all).
func letGoStackFrames(stack string) string {
	if stack == "" {
		return ""
	}
	const repo = "nooga/let-go/"
	lines := strings.Split(stack, "\n")
	var out strings.Builder
	shown := 0
	for i := 0; i+1 < len(lines) && shown < 8; i++ {
		fn := strings.TrimSpace(lines[i])
		if !strings.Contains(fn, repo) || strings.Contains(fn, "recoverThrownPanic") {
			continue
		}
		loc := strings.TrimSpace(lines[i+1])
		if idx := strings.Index(loc, repo); idx >= 0 {
			loc = loc[idx+len(repo):]
		}
		if sp := strings.IndexByte(loc, ' '); sp >= 0 {
			loc = loc[:sp]
		}
		fmt.Fprintf(&out, "  %s\n     %s\n", fn, loc)
		shown++
		i++ // consume the location line
	}
	return out.String()
}

func formatCompileError(ce compileErrorLike) string {
	var b strings.Builder
	msg := ce.InnermostMessage()
	src := ce.InnermostSource()

	fmt.Fprintf(&b, ansiBoldRed+"error:"+ansiReset+" %s\n", msg)
	if src != nil {
		writeSnippet(&b, src)
	}
	return b.String()
}

func writeSnippet(b *strings.Builder, info *SourceInfo) {
	line := SourceRegistry.GetLine(info.File, info.Line)
	if line == "" {
		fmt.Fprintf(b, "  "+ansiBoldBlue+"-->"+ansiReset+" %s\n", info.String())
		return
	}

	lineNum := info.Line + 1
	width := len(fmt.Sprintf("%d", lineNum))
	padding := strings.Repeat(" ", width)

	fmt.Fprintf(b, "  "+ansiBoldBlue+"-->"+ansiReset+" %s\n", info.String())
	fmt.Fprintf(b, " %s "+ansiBoldBlue+"|"+ansiReset+"\n", padding)
	fmt.Fprintf(b, " "+ansiBoldBlue+"%d"+ansiReset+" "+ansiBoldBlue+"|"+ansiReset+" %s\n", lineNum, line)

	// Position indicator
	col := max(info.Column, 0)
	pointer := strings.Repeat(" ", col) + ansiBoldRed + "^^^" + ansiReset
	fmt.Fprintf(b, " %s "+ansiBoldBlue+"|"+ansiReset+" %s\n", padding, pointer)
}
