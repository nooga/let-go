//go:build plan9

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/vm"
)

// Plan 9 terminal support. Input is real: raw-mode! holds /dev/consctl open with
// "rawon" (Plan 9 has no termios; consctl controls the console), and read-key /
// key-pending? read through a queuedKeySource (keysource_queued.go) — a
// background goroutine feeding a byte queue, because Plan 9 lacks the Unix
// poll(2)/FIONREAD non-blocking peek. The ANSI output functions route through
// *out* (WriteToOut) like native. size reads $COLS/$LINES (fallback 80x24); real
// window geometry (/dev/wctl pixels) is a deferred follow-up. open-pty/set-size
// stay unsupported. mouse decoding is native-only.

// consctl holds /dev/consctl open while the terminal is in raw mode. Plan 9
// reverts to cooked mode as soon as this fd closes, so raw-mode! keeps it open
// and restore-mode! writes "rawoff" then closes it. Guarded by consctlMu.
var (
	consctl   *os.File
	consctlMu sync.Mutex
)

// readEnvInt reads /env/<name> as a live file and parses it as an int, returning
// def when the file is missing or unparseable. Plan 9's environment is a
// filesystem, so reading the file each call picks up updates that os.Getenv's
// startup snapshot misses: alacritty9 (and drawterm hosts) rewrite /env/COLS and
// /env/LINES on every window resize, so the live read is what lets term/size
// track the real terminal size instead of a stale startup value.
func readEnvInt(name string, def int) int {
	b, err := os.ReadFile("/env/" + name)
	if err != nil {
		return def
	}
	// /env values can carry a trailing NUL (Plan 9 stores them NUL-terminated)
	// and/or a newline; trim before parsing.
	if n, err := strconv.Atoi(strings.TrimRight(string(b), "\x00\n\r\t ")); err == nil {
		return n
	}
	return def
}

// stdinIsConsole reports whether fd 0 is the window's console (/dev/cons). Plan
// 9 keeps /dev/cons and /dev/consctl present in the namespace even when stdin is
// redirected, so opening consctl always succeeds — checking the actual path of
// fd 0 is how we tell a real interactive console from `lg <somepipe`, and it
// keeps raw-mode! from flipping a console we aren't reading. It's the plan9
// analogue of native's "is stdin a TTY" gate.
func stdinIsConsole() bool {
	p, err := syscall.Fd2path(0)
	if err != nil {
		return false
	}
	// The cons device is /dev/cons in a normal namespace, but it can surface
	// under another path: 9front reports the raw device path #c/cons when /dev
	// isn't bound over the cons device, and a console reached through a different
	// mount carries that mount's prefix (…/dev/cons). Match those forms rather
	// than one hard-coded string, while still rejecting the pipe (#|/data) and
	// redirected-file stdins this gate exists to catch.
	return p == "#c/cons" || strings.HasSuffix(p, "/dev/cons")
}

func init() { RegisterInstaller(installTermNS) }

func installTermNS() {
	// rio doesn't render ANSI escapes — flip *ansi?* so user code (e.g.
	// the test runner's PASS/FAIL printer) can avoid emitting them.
	if v, ok := CoreNS.Lookup("*ansi?*").(*vm.Var); ok {
		v.SetRoot(vm.FALSE)
	}

	ns := vm.NewNamespace("term")
	ns.Refer(CoreNS, "", true)

	stubNil, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.NIL, nil
	})

	// Bind the plan9 key source (background stdin reader + queue) at the *keys*
	// root — the input dual of the os.Stdout *out* default in iort.go. read-key
	// and key-pending? consult the bound source, so api.WithKeySource / (binding
	// [*keys* …]) transparently overrides it for tests and embedders.
	CoreNS.Lookup("*keys*").(*vm.Var).SetRoot(vm.NewBoxed(NewQueuedKeySource(os.Stdin)))

	// raw-mode! — enter raw mode by opening /dev/consctl and writing "rawon",
	// holding the fd (Plan 9 reverts to cooked mode when it closes). Only the
	// 0-arg global form is supported; xsofy never passes a handle. Returns nil
	// (not an error) when stdin isn't the console — mirrors native's not-a-TTY
	// path, and avoids raw-moding a console we aren't reading (piped stdin).
	rawMode, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("raw-mode!: per-handle raw mode not supported on plan9")
		}
		if !stdinIsConsole() {
			return vm.NIL, nil // stdin redirected — consctl would flip the wrong console
		}
		consctlMu.Lock()
		defer consctlMu.Unlock()
		if consctl != nil {
			return vm.TRUE, nil // already raw
		}
		f, err := os.OpenFile("/dev/consctl", os.O_WRONLY, 0)
		if err != nil {
			return vm.NIL, nil // no console control available
		}
		if _, err := f.WriteString("rawon"); err != nil {
			f.Close()
			return vm.NIL, nil
		}
		consctl = f
		return vm.TRUE, nil
	})
	ns.Def("raw-mode!", rawMode)

	// restore-mode! — write "rawoff" and close /dev/consctl, reverting to cooked
	// mode. Idempotent: a no-op (nil) when raw mode was never entered.
	restoreMode, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		consctlMu.Lock()
		defer consctlMu.Unlock()
		if consctl == nil {
			return vm.NIL, nil
		}
		_, _ = consctl.WriteString("rawoff")
		err := consctl.Close()
		consctl = nil
		if err != nil {
			return vm.NIL, fmt.Errorf("restore-mode!: %w", err)
		}
		return vm.TRUE, nil
	})
	ns.Def("restore-mode!", restoreMode)

	// read-key — read one keypress through the *keys* source. Returns single
	// chars, or escape sequences like "\x1b[A" for arrows; "" is the
	// end-of-input nil contract. Ctx-aware so api.WithKeySource is honored.
	// Mouse-report decoding is native-only — plan9 has no mouse path.
	readKey := vm.NewCtxNativeFn("read-key", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		s, err := boundKeySource(ec).ReadKey()
		if err != nil {
			return vm.NIL, err
		}
		if s == "" {
			return vm.NIL, nil
		}
		return vm.String(s), nil
	})
	ns.Def("read-key", readKey)

	// key-pending? — true if a key is buffered and ready, without consuming it.
	// Plan 9 has no poll/FIONREAD, so the bound source answers from its
	// background-reader queue. Non-blocking; eof-blind so the
	// (when (key-pending?) (read-key)) idiom doesn't busy-spin at end-of-input.
	keyPendingFn := vm.NewCtxNativeFn("key-pending?", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if boundKeySource(ec).KeyPending() {
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})
	ns.Def("key-pending?", keyPendingFn)

	// size — [cols rows] read live from /env/COLS and /env/LINES (80x24 fallback).
	// alacritty9 and drawterm hosts publish the terminal size there and rewrite
	// it on every resize; reading the files each call (not os.Getenv, which
	// snapshots at process start) makes term/size track the current window, so
	// callers that diff size per tick actually see resizes. A bare rio console
	// leaves the vars unset and gets the 80x24 fallback. (Real per-window pixel
	// geometry via /dev/wctl is a separate, larger follow-up.)
	sizeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.NewPersistentVector([]vm.Value{
			vm.MakeInt(readEnvInt("COLS", 80)),
			vm.MakeInt(readEnvInt("LINES", 24)),
		}), nil
	})
	ns.Def("size", sizeFn)

	ns.Def("set-size", stubNil)

	// tty? — true when fd 0 is the window console (/dev/cons), the same gate
	// raw-mode! uses. The native arity takes a handle; plan9 answers for stdin.
	ttyPred, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if stdinIsConsole() {
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})
	ns.Def("tty?", ttyPred)

	openPty, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.NIL, fmt.Errorf("open-pty: not supported on plan9")
	})
	ns.Def("open-pty", openPty)

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

	ansi := func(name, seq string) vm.Value {
		return vm.NewCtxNativeFn(name, func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
			return vm.NIL, WriteToOut(ec, seq)
		})
	}
	ns.Def("clear", ansi("clear", "\033[2J"))
	ns.Def("clear-line", ansi("clear-line", "\033[2K"))
	ns.Def("hide-cursor", ansi("hide-cursor", "\033[?25l"))
	ns.Def("show-cursor", ansi("show-cursor", "\033[?25h"))
	ns.Def("reset-style", ansi("reset-style", "\033[0m"))
	ns.Def("bold", ansi("bold", "\033[1m"))
	ns.Def("underline", ansi("underline", "\033[4m"))
	ns.Def("inverse", ansi("inverse", "\033[7m"))
	ns.Def("alternate-screen", ansi("alternate-screen", "\033[?1049h"))
	ns.Def("main-screen", ansi("main-screen", "\033[?1049l"))

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

	// flush — sync the active *out* binding, mirroring native term.go (#308).
	//
	// Plan 9 has no fsync on a console/pipe/device: Sync there fails — the local
	// console reports EINVAL ("bad arg in system call"), while /dev/stdout over
	// drawterm reports EPERM ("permission denied") — and flushing such an fd is a
	// no-op either way. xsofy calls flush every frame, so for a console/pipe/
	// device stdout (the hot path) we skip the syscall entirely — an early no-op
	// — rather than issue a failing Fwstat per frame.
	//
	// os.Stdout.Stat().IsRegular() can't make this call on Plan 9: a console
	// Stats as "regular" through Go's FileMode (Plan 9 has no device-type bits),
	// so IsRegular misclassifies the console as a syncable file and the flush
	// then really Syncs it and errors. Classify by the fd's namespace path
	// instead — a real file is syncable; the console (/dev/cons, /dev/stdout) and
	// any #-prefixed device/pipe are not. A genuine regular-file or
	// buffered-writer *out* still gets a real Sync with its errors surfaced; the
	// EINVAL/EPERM a console-backed *out* raises (e.g. a rebound handle that slips
	// past the fast path) is swallowed below.
	stdoutSyncable := false
	if p, e := syscall.Fd2path(int(os.Stdout.Fd())); e == nil {
		stdoutSyncable = !strings.HasPrefix(p, "#") &&
			p != "/dev/cons" && p != "/dev/stdout"
	}
	flushFn := vm.NewCtxNativeFn("flush", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		h := resolveIOHandleVar(ec, "*out*")
		if (h == nil || h.File() == os.Stdout) && !stdoutSyncable {
			return vm.NIL, nil // console/pipe stdout: no fsync, no syscall
		}
		var err error
		if h != nil {
			err = h.Sync()
		} else {
			err = os.Stdout.Sync()
		}
		if errors.Is(err, syscall.EINVAL) || errors.Is(err, syscall.EPERM) {
			err = nil // console/device-backed *out* — fsync is a meaningless no-op
		}
		return vm.NIL, err
	})
	ns.Def("flush", flushFn)

	RegisterNS(ns)
}
