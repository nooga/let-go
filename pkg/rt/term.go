//go:build !js && !plan9

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"unicode/utf8"
	"unsafe"

	"github.com/nooga/let-go/pkg/vm"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

// SIGWINCH wake-byte plumbing for term/read-key.
//
// Without this, term/read-key sits in a blocking os.Stdin.Read until input
// arrives, so terminal resizes are invisible to a game loop sitting at the
// prompt — the loop's (term/size) check only runs after the next keypress
// even though the OS has already delivered SIGWINCH.
//
// Design: a SIGWINCH handler writes one BEL byte (0x07) into a self-pipe.
// read-key uses unix.Poll to wait on stdin OR the pipe. When the pipe fires,
// read-key returns BEL as if the user typed it — caller-side parse-key
// equivalents see :unknown, the game loop's (term/size) check fires on the
// next iteration and triggers a redraw at the new dimensions. No caller-side
// wiring needed beyond what's already in place for keypress-driven resize
// detection.
//
// The wake-byte pattern generalizes to any blocking read where an external
// event needs to be surfaced — only the signaling source changes.
var (
	winchPipeR     *os.File
	winchSetupOnce sync.Once
)

func setupWinch() {
	winchSetupOnce.Do(func() {
		r, w, err := os.Pipe()
		if err != nil {
			return // best-effort; read-key falls back to plain stdin
		}
		winchPipeR = r
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, unix.SIGWINCH)
		go func() {
			for range sigCh {
				// One BEL per SIGWINCH. Storms collapse — the read side
				// drains whatever's pending in a single Read regardless.
				_, _ = w.Write([]byte{0x07})
			}
		}()
	})
}

type winsize struct {
	rows, cols, xpix, ypix uint16
}

// tiocswinsz is the ioctl number for TIOCSWINSZ. Defined per-OS in
// term_pty_{linux,other}.go since the value differs across platforms.
// openPtyPair() also lives there — Linux uses /dev/ptmx + TIOCGPTN; darwin
// would need posix_openpt + grantpt/unlockpt. Container runtimes live on
// Linux so the "other" variant returns an error.

func fdFromHandle(v vm.Value) (int, error) {
	b, ok := v.(*vm.Boxed)
	if !ok {
		return -1, fmt.Errorf("expected IOHandle")
	}
	if h, ok := b.Unbox().(*IOHandle); ok {
		f := h.File()
		if f == nil {
			return -1, fmt.Errorf("IOHandle is not file-backed; raw-mode terminal ops require a real file descriptor")
		}
		return int(f.Fd()), nil
	}
	if f, ok := b.Unbox().(*os.File); ok {
		return int(f.Fd()), nil
	}
	return -1, fmt.Errorf("expected IOHandle, got %T", b.Unbox())
}

var termOldState *term.State

func init() { RegisterInstaller(installTermNS) }

// nolint
func installTermNS() {
	ns := vm.NewNamespace("term")
	ns.Refer(CoreNS, "", true)

	// raw-mode! — enter raw terminal mode
	//   (term/raw-mode!)         → TRUE, stores global state (stdin)
	//   (term/raw-mode! handle)  → opaque saved state (Boxed), caller restores it
	rawMode, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 1 {
			fd, err := fdFromHandle(vs[0])
			if err != nil {
				return vm.NIL, fmt.Errorf("raw-mode!: %w", err)
			}
			old, err := term.MakeRaw(fd)
			if err != nil {
				return vm.NIL, fmt.Errorf("raw-mode!: %w", err)
			}
			return vm.NewBoxed(old), nil
		}
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
	//   (term/restore-mode!)              → restore global stdin state
	//   (term/restore-mode! handle saved) → restore specific fd
	restoreMode, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 2 {
			fd, err := fdFromHandle(vs[0])
			if err != nil {
				return vm.NIL, fmt.Errorf("restore-mode!: %w", err)
			}
			b, ok := vs[1].(*vm.Boxed)
			if !ok {
				return vm.NIL, fmt.Errorf("restore-mode!: expected saved-state")
			}
			st, ok := b.Unbox().(*term.State)
			if !ok {
				return vm.NIL, fmt.Errorf("restore-mode!: expected saved-state, got %T", b.Unbox())
			}
			if err := term.Restore(fd, st); err != nil {
				return vm.NIL, fmt.Errorf("restore-mode!: %w", err)
			}
			return vm.TRUE, nil
		}
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
	// Returns single chars, or escape sequences like "\x1b[A" for arrow keys.
	// Also returns BEL ("\x07") when the terminal is resized (SIGWINCH) so
	// the game loop's term/size check can fire without waiting for input;
	// see setupWinch above for the design.
	readKey, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		setupWinch() // idempotent; armed on first read-key

		if winchPipeR == nil {
			// Pipe setup failed; plain blocking read (preserves prior behavior).
			buf := make([]byte, 16)
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				return vm.NIL, nil
			}
			return vm.String(buf[:n]), nil
		}

		stdinFd := int(os.Stdin.Fd())
		winchFd := int(winchPipeR.Fd())
		fds := []unix.PollFd{
			{Fd: int32(stdinFd), Events: unix.POLLIN},
			{Fd: int32(winchFd), Events: unix.POLLIN},
		}

		// Retry on EINTR — golang.org/x/sys/unix.Poll doesn't auto-retry.
		// SIGWINCH (which we install) and Go runtime preemption (SIGURG)
		// can both interrupt the blocking poll(2). Pre-PR os.Stdin.Read
		// got transparent EINTR retry via Go's ignoringEINTR; preserve
		// that contract here.
		for {
			_, err := unix.Poll(fds, -1)
			if err == unix.EINTR {
				continue
			}
			if err != nil {
				return vm.NIL, nil
			}
			break
		}

		// Drain the wake pipe unconditionally if it fired — keeps it
		// quiescent for the next cycle even when stdin also has data.
		if fds[1].Revents&unix.POLLIN != 0 {
			drain := make([]byte, 16)
			_, _ = winchPipeR.Read(drain)
		}

		// Disconnect / error on stdin (POLLHUP / POLLERR / POLLNVAL).
		// These are always reported in Revents regardless of requested
		// Events. Some kernels set POLLIN alongside (Linux on pipe EOF:
		// POLLIN|POLLHUP → falls into the POLLIN branch below, Read
		// returns n==0); some don't (Linux on closed pty master:
		// POLLHUP only). Without explicit handling, the latter would
		// fall through to the BEL return below, the caller would re-
		// enter read-key, Poll would re-fire immediately, and the loop
		// would spin at 100% CPU emitting BEL forever — EOF never
		// surfaced. Routing through os.Stdin.Read here too surfaces
		// n==0 as vm.NIL, the same end-of-input contract callers had
		// pre-PR.
		if fds[0].Revents&(unix.POLLHUP|unix.POLLERR|unix.POLLNVAL) != 0 {
			buf := make([]byte, 16)
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				return vm.NIL, nil
			}
			return vm.String(buf[:n]), nil
		}

		// Prefer real input — user keys shouldn't queue behind a resize wake.
		if fds[0].Revents&unix.POLLIN != 0 {
			buf := make([]byte, 16)
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				return vm.NIL, nil
			}
			return vm.String(buf[:n]), nil
		}

		// Only the wake fired — return BEL so the loop's term/size diff runs.
		return vm.String("\x07"), nil
	})
	ns.Def("read-key", readKey)

	// key-pending? — true if stdin has at least one byte buffered (the
	// OS has received input from the TTY but read-key has not yet
	// consumed it), OR if a SIGWINCH wake byte is queued in the winch
	// self-pipe. Uses the FIONREAD ioctl on both fds; non-consuming.
	// fionread is defined per-OS alongside tiocswinsz.
	//
	// Including the winch pipe matches what callers see on a future
	// read-key call (resize wakes surface as BEL there) and keeps the
	// (when (key-pending?) (read-key)) idiom honest for animation /
	// poll loops that want to break out on input OR on resize.
	// Returns false if both checks come up empty or error out.
	keyPendingFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		// Arm the winch handler if it hasn't been already — a caller
		// that hits key-pending? before any read-key would otherwise
		// miss SIGWINCH wakes entirely.
		setupWinch()

		var n int32
		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
			uintptr(os.Stdin.Fd()),
			uintptr(fionread),
			uintptr(unsafe.Pointer(&n)))
		if errno == 0 && n > 0 {
			return vm.TRUE, nil
		}
		if winchPipeR != nil {
			var pn int32
			_, _, perrno := syscall.Syscall(syscall.SYS_IOCTL,
				uintptr(winchPipeR.Fd()),
				uintptr(fionread),
				uintptr(unsafe.Pointer(&pn)))
			if perrno == 0 && pn > 0 {
				return vm.TRUE, nil
			}
		}
		return vm.FALSE, nil
	})
	ns.Def("key-pending?", keyPendingFn)

	// size — returns [cols rows] or nil if not a TTY
	//   (term/size)         → stdout winsize
	//   (term/size handle)  → arbitrary fd winsize
	sizeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fd := int(os.Stdout.Fd())
		if len(vs) == 1 {
			f, err := fdFromHandle(vs[0])
			if err != nil {
				return vm.NIL, fmt.Errorf("size: %w", err)
			}
			fd = f
		}
		w, h, err := term.GetSize(fd)
		if err != nil {
			return vm.NIL, nil // not a TTY
		}
		return vm.NewPersistentVector([]vm.Value{vm.MakeInt(w), vm.MakeInt(h)}), nil
	})
	ns.Def("size", sizeFn)

	// set-size — (term/set-size handle cols rows) — TIOCSWINSZ
	setSize, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("set-size expects 3 args (handle cols rows)")
		}
		fd, err := fdFromHandle(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("set-size: %w", err)
		}
		cols, ok1 := vs[1].(vm.Int)
		rows, ok2 := vs[2].(vm.Int)
		if !ok1 || !ok2 {
			return vm.NIL, fmt.Errorf("set-size: expected ints (cols rows)")
		}
		ws := winsize{rows: uint16(rows), cols: uint16(cols)}
		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
			uintptr(fd), uintptr(tiocswinsz), uintptr(unsafe.Pointer(&ws)))
		if errno != 0 {
			return vm.NIL, fmt.Errorf("set-size: %v", errno)
		}
		return vm.NIL, nil
	})
	ns.Def("set-size", setSize)

	// tty? — (term/tty? handle) → bool
	ttyPred, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("tty? expects 1 arg (handle)")
		}
		fd, err := fdFromHandle(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("tty?: %w", err)
		}
		if term.IsTerminal(fd) {
			return vm.TRUE, nil
		}
		return vm.FALSE, nil
	})
	ns.Def("tty?", ttyPred)

	// open-pty — (term/open-pty) → {:master IOHandle :slave IOHandle :slave-path "/dev/pts/N"}
	openPty, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("open-pty expects 0 args")
		}
		master, slave, slavePath, err := openPtyPair()
		if err != nil {
			return vm.NIL, fmt.Errorf("open-pty: %w", err)
		}
		m := vm.EmptyPersistentMap
		m = m.Assoc(vm.Keyword("master"), vm.NewBoxed(NewIOHandle(master))).(*vm.PersistentMap)
		m = m.Assoc(vm.Keyword("slave"), vm.NewBoxed(NewIOHandle(slave))).(*vm.PersistentMap)
		m = m.Assoc(vm.Keyword("slave-path"), vm.String(slavePath)).(*vm.PersistentMap)
		return m, nil
	})
	ns.Def("open-pty", openPty)

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
