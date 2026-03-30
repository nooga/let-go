/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/vm"
	"golang.org/x/term"
)

var termOldState *term.State

// nolint
func installTermNS() {
	ns := vm.NewNamespace("term")
	ns.Refer(CoreNS, "", true)

	// raw-mode! — enter raw terminal mode, returns true
	rawMode, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if termOldState != nil {
			return vm.TRUE, nil // already in raw mode
		}
		old, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return vm.NIL, nil // not a TTY
		}
		termOldState = old
		return vm.TRUE, nil
	})
	ns.Def("raw-mode!", rawMode)

	// restore-mode! — restore terminal to original state
	restoreMode, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if termOldState == nil {
			return vm.NIL, nil
		}
		err := term.Restore(int(os.Stdin.Fd()), termOldState)
		termOldState = nil
		if err != nil {
			return vm.NIL, fmt.Errorf("restore-mode!: %w", err)
		}
		return vm.TRUE, nil
	})
	ns.Def("restore-mode!", restoreMode)

	// read-key — read a single keypress, returns a string
	// Returns single chars, or escape sequences like "\x1b[A" for arrow keys
	readKey, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		buf := make([]byte, 16)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return vm.NIL, nil // EOF or error
		}
		if n == 0 {
			return vm.NIL, nil
		}
		return vm.String(buf[:n]), nil
	})
	ns.Def("read-key", readKey)

	// size — returns [cols rows] or nil if not a TTY
	sizeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		w, h, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return vm.NIL, nil // not a TTY
		}
		return vm.NewPersistentVector([]vm.Value{vm.MakeInt(w), vm.MakeInt(h)}), nil
	})
	ns.Def("size", sizeFn)

	// move-cursor — (move-cursor col row) — 1-based ANSI positioning
	moveCursor, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("move-cursor expects 2 args (col row)")
		}
		col, ok1 := vs[0].(vm.Int)
		row, ok2 := vs[1].(vm.Int)
		if !ok1 || !ok2 {
			return vm.NIL, fmt.Errorf("move-cursor expects integers")
		}
		fmt.Printf("\033[%d;%dH", int(row), int(col))
		return vm.NIL, nil
	})
	ns.Def("move-cursor", moveCursor)

	// clear — clear screen
	clearFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[2J")
		return vm.NIL, nil
	})
	ns.Def("clear", clearFn)

	// clear-line — clear current line
	clearLine, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[2K")
		return vm.NIL, nil
	})
	ns.Def("clear-line", clearLine)

	// hide-cursor — hide terminal cursor
	hideCursor, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[?25l")
		return vm.NIL, nil
	})
	ns.Def("hide-cursor", hideCursor)

	// show-cursor — show terminal cursor
	showCursor, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[?25h")
		return vm.NIL, nil
	})
	ns.Def("show-cursor", showCursor)

	// set-fg — (set-fg r g b) or (set-fg color-code)
	setFg, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		switch len(vs) {
		case 1:
			c, ok := vs[0].(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("set-fg expects integer color code")
			}
			fmt.Printf("\033[38;5;%dm", int(c))
		case 3:
			r, ok1 := vs[0].(vm.Int)
			g, ok2 := vs[1].(vm.Int)
			b, ok3 := vs[2].(vm.Int)
			if !ok1 || !ok2 || !ok3 {
				return vm.NIL, fmt.Errorf("set-fg expects 3 integers (r g b)")
			}
			fmt.Printf("\033[38;2;%d;%d;%dm", int(r), int(g), int(b))
		default:
			return vm.NIL, fmt.Errorf("set-fg expects 1 or 3 args")
		}
		return vm.NIL, nil
	})
	ns.Def("set-fg", setFg)

	// set-bg — (set-bg r g b) or (set-bg color-code)
	setBg, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		switch len(vs) {
		case 1:
			c, ok := vs[0].(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("set-bg expects integer color code")
			}
			fmt.Printf("\033[48;5;%dm", int(c))
		case 3:
			r, ok1 := vs[0].(vm.Int)
			g, ok2 := vs[1].(vm.Int)
			b, ok3 := vs[2].(vm.Int)
			if !ok1 || !ok2 || !ok3 {
				return vm.NIL, fmt.Errorf("set-bg expects 3 integers (r g b)")
			}
			fmt.Printf("\033[48;2;%d;%d;%dm", int(r), int(g), int(b))
		default:
			return vm.NIL, fmt.Errorf("set-bg expects 1 or 3 args")
		}
		return vm.NIL, nil
	})
	ns.Def("set-bg", setBg)

	// reset-style — reset all ANSI attributes
	resetStyle, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[0m")
		return vm.NIL, nil
	})
	ns.Def("reset-style", resetStyle)

	// bold — enable bold
	bold, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[1m")
		return vm.NIL, nil
	})
	ns.Def("bold", bold)

	// underline — enable underline
	underline, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[4m")
		return vm.NIL, nil
	})
	ns.Def("underline", underline)

	// inverse — enable inverse/reverse video
	inverse, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[7m")
		return vm.NIL, nil
	})
	ns.Def("inverse", inverse)

	// write — (write str) — write string at current cursor position, no newline
	writeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("write expects 1 arg")
		}
		var s string
		if str, ok := vs[0].(vm.String); ok {
			s = string(str)
		} else {
			s = vs[0].String()
		}
		fmt.Print(s)
		return vm.NIL, nil
	})
	ns.Def("write", writeFn)

	// write-at — (write-at col row str) — write string at position
	writeAt, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("write-at expects 3 args (col row str)")
		}
		col, ok1 := vs[0].(vm.Int)
		row, ok2 := vs[1].(vm.Int)
		if !ok1 || !ok2 {
			return vm.NIL, fmt.Errorf("write-at expects integers for col and row")
		}
		var s string
		if str, ok := vs[2].(vm.String); ok {
			s = string(str)
		} else {
			s = vs[2].String()
		}
		fmt.Printf("\033[%d;%dH%s", int(row), int(col), s)
		return vm.NIL, nil
	})
	ns.Def("write-at", writeAt)

	// char-width — (char-width str) — returns display width of first char (1 or 2 for CJK)
	charWidth, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("char-width expects 1 arg")
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("char-width expects string")
		}
		r, _ := utf8.DecodeRuneInString(string(s))
		if r == utf8.RuneError {
			return vm.MakeInt(0), nil
		}
		return vm.MakeInt(1), nil
	})
	ns.Def("char-width", charWidth)

	// alternate-screen — switch to alternate screen buffer
	altScreen, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[?1049h")
		return vm.NIL, nil
	})
	ns.Def("alternate-screen", altScreen)

	// main-screen — switch back to main screen buffer
	mainScreen, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[?1049l")
		return vm.NIL, nil
	})
	ns.Def("main-screen", mainScreen)

	// flush — flush stdout
	flushFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		os.Stdout.Sync()
		return vm.NIL, nil
	})
	ns.Def("flush", flushFn)

	RegisterNS(ns)
}
