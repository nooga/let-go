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

// jsLoadInt does an Atomics.load and returns the int value, falling back to 0
// when the result is undefined/null (e.g. an uninitialized SAB slot or a
// non-cross-origin-isolated context). This guards the `.Int()` call, which
// otherwise panics with "syscall/js: call of Value.Int on undefined".
func jsLoadInt(atomics, arr js.Value, idx int) int {
	v := atomics.Call("load", arr, idx)
	if v.IsUndefined() || v.IsNull() {
		return 0
	}
	return v.Int()
}

func init() { RegisterInstaller(installTermNS) }

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
		keyLen := jsLoadInt(atomics, keyInt32, 1)
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

	// key-pending? — true if the key-ready flag at SAB[0] is set, i.e.
	// the JS side has written a key that read-key has not yet consumed.
	// Lets animation / poll loops break out on user input without
	// entering read-key (which would Atomics.wait if no key is queued).
	keyPendingFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		keyInt32 := js.Global().Get("_lgKeyInt32")
		atomics := js.Global().Get("Atomics")
		if keyInt32.IsUndefined() || atomics.IsUndefined() {
			return vm.FALSE, nil
		}
		v := atomics.Call("load", keyInt32, 0)
		if v.IsUndefined() || v.IsNull() {
			return vm.FALSE, nil
		}
		if v.Int() != 0 {
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})
	ns.Def("key-pending?", keyPendingFn)

	// size — read from SharedArrayBuffer
	sizeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		keyInt32 := js.Global().Get("_lgKeyInt32")
		atomics := js.Global().Get("Atomics")
		if keyInt32.IsUndefined() || atomics.IsUndefined() {
			return vm.NewPersistentVector([]vm.Value{vm.MakeInt(80), vm.MakeInt(24)}), nil
		}
		w := jsLoadInt(atomics, keyInt32, 6)
		h := jsLoadInt(atomics, keyInt32, 7)
		if w == 0 {
			w = 80
		}
		if h == 0 {
			h = 24
		}
		return vm.NewPersistentVector([]vm.Value{vm.MakeInt(w), vm.MakeInt(h)}), nil
	})
	ns.Def("size", sizeFn)

	// --- Output functions — identical to native; route ANSI through *out* ---
	// (WriteToOut) so (binding [*out* …]) is honored. xterm.js handles ANSI
	// natively. *out*'s root is os.Stdout, so the served bundle's output still
	// flows through the existing fs.writeSync path until the host-writer lands.

	moveCursor := vm.NewCtxNativeFn("move-cursor", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("move-cursor expects 2 args (col row)")
		}
		col, ok1 := vs[0].(vm.Int)
		row, ok2 := vs[1].(vm.Int)
		if !ok1 || !ok2 {
			return vm.NIL, fmt.Errorf("move-cursor expects integers")
		}
		return vm.NIL, WriteToOut(ec, fmt.Sprintf("\033[%d;%dH", int(row), int(col)))
	})
	ns.Def("move-cursor", moveCursor)

	clearFn := vm.NewCtxNativeFn("clear", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[2J")
	})
	ns.Def("clear", clearFn)

	clearLine := vm.NewCtxNativeFn("clear-line", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[2K")
	})
	ns.Def("clear-line", clearLine)

	hideCursor := vm.NewCtxNativeFn("hide-cursor", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[?25l")
	})
	ns.Def("hide-cursor", hideCursor)

	showCursor := vm.NewCtxNativeFn("show-cursor", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[?25h")
	})
	ns.Def("show-cursor", showCursor)

	setFg := vm.NewCtxNativeFn("set-fg", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		var seq string
		switch len(vs) {
		case 1:
			c, ok := vs[0].(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("set-fg expects integer color code")
			}
			seq = fmt.Sprintf("\033[38;5;%dm", int(c))
		case 3:
			r, ok1 := vs[0].(vm.Int)
			g, ok2 := vs[1].(vm.Int)
			b, ok3 := vs[2].(vm.Int)
			if !ok1 || !ok2 || !ok3 {
				return vm.NIL, fmt.Errorf("set-fg expects 3 integers (r g b)")
			}
			seq = fmt.Sprintf("\033[38;2;%d;%d;%dm", int(r), int(g), int(b))
		default:
			return vm.NIL, fmt.Errorf("set-fg expects 1 or 3 args")
		}
		return vm.NIL, WriteToOut(ec, seq)
	})
	ns.Def("set-fg", setFg)

	setBg := vm.NewCtxNativeFn("set-bg", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		var seq string
		switch len(vs) {
		case 1:
			c, ok := vs[0].(vm.Int)
			if !ok {
				return vm.NIL, fmt.Errorf("set-bg expects integer color code")
			}
			seq = fmt.Sprintf("\033[48;5;%dm", int(c))
		case 3:
			r, ok1 := vs[0].(vm.Int)
			g, ok2 := vs[1].(vm.Int)
			b, ok3 := vs[2].(vm.Int)
			if !ok1 || !ok2 || !ok3 {
				return vm.NIL, fmt.Errorf("set-bg expects 3 integers (r g b)")
			}
			seq = fmt.Sprintf("\033[48;2;%d;%d;%dm", int(r), int(g), int(b))
		default:
			return vm.NIL, fmt.Errorf("set-bg expects 1 or 3 args")
		}
		return vm.NIL, WriteToOut(ec, seq)
	})
	ns.Def("set-bg", setBg)

	resetStyle := vm.NewCtxNativeFn("reset-style", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[0m")
	})
	ns.Def("reset-style", resetStyle)

	bold := vm.NewCtxNativeFn("bold", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[1m")
	})
	ns.Def("bold", bold)

	underline := vm.NewCtxNativeFn("underline", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[4m")
	})
	ns.Def("underline", underline)

	inverse := vm.NewCtxNativeFn("inverse", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[7m")
	})
	ns.Def("inverse", inverse)

	writeFn := vm.NewCtxNativeFn("write", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("write expects 1 arg")
		}
		var s string
		if str, ok := vs[0].(vm.String); ok {
			s = string(str)
		} else {
			s = vs[0].String()
		}
		return vm.NIL, WriteToOut(ec, s)
	})
	ns.Def("write", writeFn)

	writeAt := vm.NewCtxNativeFn("write-at", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
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
		return vm.NIL, WriteToOut(ec, fmt.Sprintf("\033[%d;%dH%s", int(row), int(col), s))
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

	altScreen := vm.NewCtxNativeFn("alternate-screen", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[?1049h")
	})
	ns.Def("alternate-screen", altScreen)

	mainScreen := vm.NewCtxNativeFn("main-screen", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[?1049l")
	})
	ns.Def("main-screen", mainScreen)

	// flush — sync the active *out* binding, then drive the JS-side transport
	// flush so the default screen path still reaches xterm.js. A rebound
	// (buffered/embedder) *out* is flushed by Sync(); _lgFlush then no-ops on
	// the empty screen buffer. The default-root *out* wraps os.Stdout, whose
	// Sync() is a no-op in wasm, so _lgFlush stays the real screen flush.
	flushFn := vm.NewCtxNativeFn("flush", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if h := resolveIOHandleVar(ec, "*out*"); h != nil {
			if err := h.Sync(); err != nil {
				return vm.NIL, err
			}
		}
		flush := js.Global().Get("_lgFlush")
		if !flush.IsUndefined() {
			flush.Invoke()
		}
		return vm.NIL, nil
	})
	ns.Def("flush", flushFn)

	RegisterNS(ns)
}
