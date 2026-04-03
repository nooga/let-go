//go:build js && wasm

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"syscall/js"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/vm"
)

// SharedArrayBuffer layout (Int32Array indices):
//   [0] key ready flag (Atomics.wait/notify target)
//   [1] key byte count
//   [6] terminal cols
//   [7] terminal rows
// Uint8Array view at byte offset 8, length 16: raw key bytes

func installTermNS() {
	// Set *in-wasm* so user code can detect WASM environment
	CoreNS.Lookup("*in-wasm*").(*vm.Var).SetRoot(vm.TRUE)

	ns := vm.NewNamespace("term")
	ns.Refer(CoreNS, "", true)

	// raw-mode! — no-op in WASM (xterm.js is always raw)
	rawMode, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.TRUE, nil
	})
	ns.Def("raw-mode!", rawMode)

	// restore-mode! — no-op in WASM
	restoreMode, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.TRUE, nil
	})
	ns.Def("restore-mode!", restoreMode)

	// read-key — blocks via Atomics.wait on SharedArrayBuffer
	readKey, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		atomics := js.Global().Get("Atomics")
		keyInt32 := js.Global().Get("_lgKeyInt32")
		keyUint8 := js.Global().Get("_lgKeyUint8")

		if keyInt32.IsUndefined() || keyUint8.IsUndefined() {
			return vm.NIL, fmt.Errorf("read-key: terminal input not available (no SharedArrayBuffer)")
		}

		// Flush output before blocking
		js.Global().Call("_lgFlush")

		// Block until a key is ready
		atomics.Call("wait", keyInt32, 0, 0)

		// Read key length and bytes
		keyLen := atomics.Call("load", keyInt32, 1).Int()
		if keyLen <= 0 || keyLen > 16 {
			atomics.Call("store", keyInt32, 0, 0)
			return vm.NIL, nil
		}

		keyBytes := make([]byte, keyLen)
		for i := 0; i < keyLen; i++ {
			keyBytes[i] = byte(keyUint8.Index(i).Int())
		}

		// Reset flag for next key
		atomics.Call("store", keyInt32, 0, 0)

		return vm.String(keyBytes), nil
	})
	ns.Def("read-key", readKey)

	// size — read from SharedArrayBuffer
	sizeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		keyInt32 := js.Global().Get("_lgKeyInt32")
		if keyInt32.IsUndefined() {
			return vm.NewPersistentVector([]vm.Value{vm.MakeInt(80), vm.MakeInt(24)}), nil
		}
		atomics := js.Global().Get("Atomics")
		w := atomics.Call("load", keyInt32, 6).Int()
		h := atomics.Call("load", keyInt32, 7).Int()
		if w == 0 {
			w = 80
		}
		if h == 0 {
			h = 24
		}
		return vm.NewPersistentVector([]vm.Value{vm.MakeInt(w), vm.MakeInt(h)}), nil
	})
	ns.Def("size", sizeFn)

	// --- Output functions — identical to native, just emit ANSI via fmt.Print ---
	// xterm.js handles all ANSI escape sequences natively.

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

	clearFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[2J")
		return vm.NIL, nil
	})
	ns.Def("clear", clearFn)

	clearLine, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[2K")
		return vm.NIL, nil
	})
	ns.Def("clear-line", clearLine)

	hideCursor, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[?25l")
		return vm.NIL, nil
	})
	ns.Def("hide-cursor", hideCursor)

	showCursor, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[?25h")
		return vm.NIL, nil
	})
	ns.Def("show-cursor", showCursor)

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

	resetStyle, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[0m")
		return vm.NIL, nil
	})
	ns.Def("reset-style", resetStyle)

	bold, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[1m")
		return vm.NIL, nil
	})
	ns.Def("bold", bold)

	underline, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[4m")
		return vm.NIL, nil
	})
	ns.Def("underline", underline)

	inverse, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[7m")
		return vm.NIL, nil
	})
	ns.Def("inverse", inverse)

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

	altScreen, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[?1049h")
		return vm.NIL, nil
	})
	ns.Def("alternate-screen", altScreen)

	mainScreen, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fmt.Print("\033[?1049l")
		return vm.NIL, nil
	})
	ns.Def("main-screen", mainScreen)

	// flush — flush buffered output to xterm.js via postMessage
	flushFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		flush := js.Global().Get("_lgFlush")
		if !flush.IsUndefined() {
			flush.Invoke()
		}
		return vm.NIL, nil
	})
	ns.Def("flush", flushFn)

	RegisterNS(ns)
}
