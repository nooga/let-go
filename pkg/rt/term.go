//go:build !js && !plan9 && !wasip1

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"errors"
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

// nativeKeySource is the native *keys* binding: term/read-key and key-pending?
// read from stdin, with the SIGWINCH wake-byte self-pipe (setupWinch) breaking
// a blocked read on resize. The body is the pre-seam read-key/key-pending?
// logic, unchanged — it just lives behind the KeySource interface now so the
// term ops, embedders (api.WithKeySource), and tests share one source.
type nativeKeySource struct{}

// keyBuf holds bytes read from stdin but not yet handed out as keys. A single
// raw read can carry several keys (held-key auto-repeat / queued input) or a
// multi-byte escape sequence; ReadKey tokenizes one key per call off this
// buffer (see scanKey) and refills via readRaw only when it's empty. Native
// has a single stdin and reads aren't concurrent, so a package-level buffer is
// safe.
var keyBuf []byte

// readChunkSize is the per-read stdin buffer. os.Stdin.Read returns as soon as
// any data is available (it never waits to fill), so a larger buffer adds no
// latency to a single keystroke — it just grabs more bytes when a burst is
// already queued, cutting syscalls and ReadKey refills. Correctness doesn't
// depend on it: a token straddling the buffer is stitched by the refill path
// (see scanKey). 256 covers held-key bursts and back-to-back escape
// sequences with headroom.
const readChunkSize = 256

func (s nativeKeySource) ReadKey() (string, error) {
	for {
		if len(keyBuf) == 0 {
			chunk, err := s.readRaw()
			if err != nil {
				return "", err
			}
			if len(chunk) == 0 {
				return "", nil // EOF / nil contract
			}
			if isWinchWake(chunk) {
				return "\x07", nil // SIGWINCH wake — synthetic, not tokenized
			}
			keyBuf = chunk
		}
		status, n := scanKey(keyBuf)
		// If the front token is split (escape sequence or UTF-8 rune) and more
		// bytes are actually waiting on stdin, pull them in and scan again
		// rather than emitting a broken partial. Gating on rawPending keeps a
		// genuinely truncated token from blocking a refill that would never
		// complete.
		if status == keyNeedMore && s.rawPending() {
			chunk, err := s.readRaw()
			if err != nil {
				return "", err
			}
			if len(chunk) != 0 && !isWinchWake(chunk) {
				keyBuf = append(keyBuf, chunk...)
				continue
			}
		}
		tok := string(keyBuf[:n])
		keyBuf = keyBuf[n:]
		if len(keyBuf) == 0 {
			keyBuf = nil // release the backing array once drained
		}
		return tok, nil
	}
}

func isWinchWake(b []byte) bool {
	return len(b) == 1 && b[0] == '\x07'
}

// readRaw does one blocking poll+read, returning raw stdin bytes, nil on EOF,
// or BEL ("\x07") when only the SIGWINCH wake fired. This is the pre-tokenizer
// read-key body, except it now keeps input as bytes until ReadKey emits the
// final token string.
func (nativeKeySource) readRaw() ([]byte, error) {
	setupWinch() // idempotent; armed on first read-key

	if winchPipeR == nil {
		// Pipe setup failed; plain blocking read (preserves prior behavior).
		buf := make([]byte, readChunkSize)
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return nil, nil
		}
		return buf[:n], nil
	}

	stdinFd := int(os.Stdin.Fd())
	winchFd := int(winchPipeR.Fd())
	fds := []unix.PollFd{
		{Fd: int32(stdinFd), Events: unix.POLLIN},
		{Fd: int32(winchFd), Events: unix.POLLIN},
	}

	// Retry on EINTR — golang.org/x/sys/unix.Poll doesn't auto-retry.
	// SIGWINCH (which we install) and Go runtime preemption (SIGURG) can
	// both interrupt the blocking poll(2). Pre-PR os.Stdin.Read got
	// transparent EINTR retry via Go's ignoringEINTR; preserve that here.
	for {
		_, err := unix.Poll(fds, -1)
		if err == unix.EINTR {
			continue
		}
		if err != nil {
			return nil, nil
		}
		break
	}

	// Drain the wake pipe unconditionally if it fired — keeps it quiescent
	// for the next cycle even when stdin also has data.
	if fds[1].Revents&unix.POLLIN != 0 {
		drain := make([]byte, 16)
		_, _ = winchPipeR.Read(drain)
	}

	// Disconnect / error on stdin (POLLHUP / POLLERR / POLLNVAL). These are
	// always reported in Revents regardless of requested Events. Some kernels
	// set POLLIN alongside (Linux on pipe EOF: POLLIN|POLLHUP → falls into the
	// POLLIN branch below, Read returns n==0); some don't (Linux on closed pty
	// master: POLLHUP only). Without explicit handling, the latter would fall
	// through to the BEL return below, the caller would re-enter read-key,
	// Poll would re-fire immediately, and the loop would spin at 100% CPU
	// emitting BEL forever — EOF never surfaced. Routing through os.Stdin.Read
	// here too surfaces n==0 as the nil contract callers had pre-PR.
	if fds[0].Revents&(unix.POLLHUP|unix.POLLERR|unix.POLLNVAL) != 0 {
		buf := make([]byte, readChunkSize)
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return nil, nil
		}
		return buf[:n], nil
	}

	// Prefer real input — user keys shouldn't queue behind a resize wake.
	if fds[0].Revents&unix.POLLIN != 0 {
		buf := make([]byte, readChunkSize)
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return nil, nil
		}
		return buf[:n], nil
	}

	// Only the wake fired — return BEL so the loop's term/size diff runs.
	return []byte{'\x07'}, nil
}

func (s nativeKeySource) KeyPending() bool {
	if len(keyBuf) > 0 {
		return true // tokens still buffered from a prior multi-key read
	}
	return s.rawPending()
}

// rawPending reports whether stdin (or the SIGWINCH wake-pipe) has bytes the
// buffer hasn't consumed yet — the OS-level half of KeyPending, without the
// keyBuf check. ReadKey uses it to decide whether refilling a split token
// will actually make progress.
func (nativeKeySource) rawPending() bool {
	// Arm the winch handler if it hasn't been already — a caller that hits
	// key-pending? before any read-key would otherwise miss SIGWINCH wakes.
	setupWinch()

	var n int32
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(os.Stdin.Fd()),
		uintptr(fionread),
		uintptr(unsafe.Pointer(&n)))
	if errno == 0 && n > 0 {
		return true
	}
	if winchPipeR != nil {
		var pn int32
		_, _, perrno := syscall.Syscall(syscall.SYS_IOCTL,
			uintptr(winchPipeR.Fd()),
			uintptr(fionread),
			uintptr(unsafe.Pointer(&pn)))
		if perrno == 0 && pn > 0 {
			return true
		}
	}
	return false
}

// mouseEventMap renders a decoded SGR mouse report as the tagged map read-key
// hands to Clojure: {:type :mouse :action :press :button :left :x 10 :y 5
// :mods #{:shift}}. :mods is an empty set when no modifier is held.
func mouseEventMap(ev MouseEvent) vm.Value {
	var mods []vm.Value
	if ev.Shift {
		mods = append(mods, vm.Keyword("shift"))
	}
	if ev.Ctrl {
		mods = append(mods, vm.Keyword("ctrl"))
	}
	if ev.Meta {
		mods = append(mods, vm.Keyword("meta"))
	}
	m := vm.EmptyPersistentMap
	m = m.Assoc(vm.Keyword("type"), vm.Keyword("mouse")).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("action"), vm.Keyword(ev.Action)).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("button"), vm.Keyword(ev.Button)).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("x"), vm.MakeInt(ev.X)).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("y"), vm.MakeInt(ev.Y)).(*vm.PersistentMap)
	m = m.Assoc(vm.Keyword("mods"), vm.NewPersistentSet(mods)).(*vm.PersistentMap)
	return m
}

func init() { RegisterInstaller(installTermNS) }

// nolint
func installTermNS() {
	// Bind the native key source (stdin + SIGWINCH wake-pipe) at the *keys*
	// root, the input dual of the os.Stdout *out* default in iort.go.
	CoreNS.Lookup("*keys*").(*vm.Var).SetRoot(vm.NewBoxed(nativeKeySource{}))

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

	// read-key — read a single keypress through the *keys* source.
	// Returns single chars, or escape sequences like "\x1b[A" for arrow keys.
	// nativeKeySource also returns BEL ("\x07") on resize (SIGWINCH) so the
	// game loop's term/size check can fire without waiting for input; see
	// setupWinch above. Ctx-aware so api.WithKeySource / (binding [*keys* …])
	// is honored; "" is the end-of-input nil contract.
	readKey := vm.NewCtxNativeFn("read-key", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		s, err := boundKeySource(ec).ReadKey()
		if err != nil {
			return vm.NIL, err
		}
		if s == "" {
			return vm.NIL, nil
		}
		// Hybrid return: plain keys stay strings (no caller churn), an SGR
		// mouse report becomes a tagged map. Mouse reports only arrive once
		// enable-mouse has been called, so string-only callers are unaffected.
		if ev, ok := decodeSGRMouse(s); ok {
			return mouseEventMap(ev), nil
		}
		return vm.String(s), nil
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
	keyPendingFn := vm.NewCtxNativeFn("key-pending?", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if boundKeySource(ec).KeyPending() {
			return vm.TRUE, nil
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

	// clear — clear screen
	clearFn := vm.NewCtxNativeFn("clear", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[2J")
	})
	ns.Def("clear", clearFn)

	// clear-line — clear current line
	clearLine := vm.NewCtxNativeFn("clear-line", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[2K")
	})
	ns.Def("clear-line", clearLine)

	// hide-cursor — hide terminal cursor
	hideCursor := vm.NewCtxNativeFn("hide-cursor", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[?25l")
	})
	ns.Def("hide-cursor", hideCursor)

	// show-cursor — show terminal cursor
	showCursor := vm.NewCtxNativeFn("show-cursor", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[?25h")
	})
	ns.Def("show-cursor", showCursor)

	// set-fg — (set-fg r g b) or (set-fg color-code)
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

	// set-bg — (set-bg r g b) or (set-bg color-code)
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

	// reset-style — reset all ANSI attributes
	resetStyle := vm.NewCtxNativeFn("reset-style", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[0m")
	})
	ns.Def("reset-style", resetStyle)

	// bold — enable bold
	bold := vm.NewCtxNativeFn("bold", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[1m")
	})
	ns.Def("bold", bold)

	// underline — enable underline
	underline := vm.NewCtxNativeFn("underline", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[4m")
	})
	ns.Def("underline", underline)

	// inverse — enable inverse/reverse video
	inverse := vm.NewCtxNativeFn("inverse", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[7m")
	})
	ns.Def("inverse", inverse)

	// write — (write str) — write string at current cursor position, no newline
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

	// write-at — (write-at col row str) — write string at position
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
	altScreen := vm.NewCtxNativeFn("alternate-screen", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[?1049h")
	})
	ns.Def("alternate-screen", altScreen)

	// main-screen — switch back to main screen buffer
	mainScreen := vm.NewCtxNativeFn("main-screen", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[?1049l")
	})
	ns.Def("main-screen", mainScreen)

	// enable-mouse — turn on SGR mouse reporting (1000 button tracking + 1006
	// SGR extended coords). Reports then arrive through read-key as tagged maps
	// (see decodeSGRMouse). Click-only: drag/motion modes (1002/1003) are not
	// enabled. The same sequence works through xterm.js for the WASM build.
	// Callers must pair this with disable-mouse in their cleanup, alongside
	// restore-mode!, or the terminal is left emitting click escapes.
	enableMouse := vm.NewCtxNativeFn("enable-mouse", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[?1000;1006h")
	})
	ns.Def("enable-mouse", enableMouse)

	// disable-mouse — turn off the reporting that enable-mouse switched on.
	disableMouse := vm.NewCtxNativeFn("disable-mouse", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		return vm.NIL, WriteToOut(ec, "\033[?1000;1006l")
	})
	ns.Def("disable-mouse", disableMouse)

	// flush — sync the active *out* binding so flush hits the same sink the
	// term/* bytes did. File-backed handles fsync; buffered embedder writers
	// flush; otherwise no-op. Falls back to os.Stdout only at early boot
	// (before *out* is installed). Pre-#223 this synced os.Stdout directly,
	// which diverged once the term/* ops started honoring *out*.
	flushFn := vm.NewCtxNativeFn("flush", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		var (
			err        error
			fileBacked bool
		)
		if h := resolveIOHandleVar(ec, "*out*"); h != nil {
			err = h.Sync()
			fileBacked = h.File() != nil
		} else {
			err = os.Stdout.Sync()
			fileBacked = true
		}
		// fsync on terminals and pipes returns ENOTTY, EBADF, or EINVAL
		// depending on the platform; flushing those fds is a no-op, so swallow
		// them. Only for a file-backed *out* -- an embedder writer's Sync goes
		// through Flush() and must surface its own errors. A regular file's
		// fsync never returns these, so real I/O errors (EIO, ENOSPC, etc.)
		// still propagate.
		if fileBacked && (errors.Is(err, syscall.ENOTTY) ||
			errors.Is(err, syscall.EBADF) ||
			errors.Is(err, syscall.EINVAL)) {
			err = nil
		}
		return vm.NIL, err
	})
	ns.Def("flush", flushFn)

	RegisterNS(ns)
}
