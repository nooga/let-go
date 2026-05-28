//go:build plan9

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
)

// Plan 9 has no termios/ioctl. raw-mode/size/pty are stubbed; the ANSI output
// functions still run via fmt.Print (rio won't render escapes, but they're
// harmless). If you need real terminal control on plan9, wire in /dev/cons here.

func init() { RegisterInstaller(installTermNS) }

func installTermNS() {
	// rio doesn't render ANSI escapes — flip *ansi?* so user code (e.g.
	// the test runner's PASS/FAIL printer) can avoid emitting them.
	if v, ok := CoreNS.Lookup("*ansi?*").(*vm.Var); ok {
		v.SetRoot(vm.FALSE)
	}

	ns := vm.NewNamespace("term")
	ns.Refer(CoreNS, "", true)

	stubTrue, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.TRUE, nil
	})
	stubNil, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.NIL, nil
	})

	ns.Def("raw-mode!", stubTrue)
	ns.Def("restore-mode!", stubTrue)
	ns.Def("read-key", stubNil)

	sizeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.NewPersistentVector([]vm.Value{vm.MakeInt(80), vm.MakeInt(24)}), nil
	})
	ns.Def("size", sizeFn)

	ns.Def("set-size", stubNil)

	ttyPred, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.FALSE, nil
	})
	ns.Def("tty?", ttyPred)

	openPty, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.NIL, fmt.Errorf("open-pty: not supported on plan9")
	})
	ns.Def("open-pty", openPty)

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

	ansi := func(seq string) vm.Value {
		fn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
			fmt.Print(seq)
			return vm.NIL, nil
		})
		return fn
	}
	ns.Def("clear", ansi("\033[2J"))
	ns.Def("clear-line", ansi("\033[2K"))
	ns.Def("hide-cursor", ansi("\033[?25l"))
	ns.Def("show-cursor", ansi("\033[?25h"))
	ns.Def("reset-style", ansi("\033[0m"))
	ns.Def("bold", ansi("\033[1m"))
	ns.Def("underline", ansi("\033[4m"))
	ns.Def("inverse", ansi("\033[7m"))
	ns.Def("alternate-screen", ansi("\033[?1049h"))
	ns.Def("main-screen", ansi("\033[?1049l"))

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

	flushFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		os.Stdout.Sync()
		return vm.NIL, nil
	})
	ns.Def("flush", flushFn)

	RegisterNS(ns)
}
