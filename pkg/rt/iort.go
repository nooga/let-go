/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/nooga/let-go/pkg/vm"
)

// IOHandle is the runtime's I/O abstraction. It wraps an io.Writer and/or
// io.Reader so the same type backs file handles, stdin/stdout/stderr,
// bytes.Buffer-style captures, and embedder-supplied sinks. File-backed
// handles also retain the *os.File so callers that need Fd()/Sync() still
// have it.
type IOHandle struct {
	name   string
	file   *os.File // set only when this handle wraps a file
	writer io.Writer
	reader io.Reader
	bufrd  *bufio.Reader
}

// NewIOHandle wraps an *os.File. Preserves all existing call sites and
// keeps file-specific operations (Fd, Sync) available via File().
func NewIOHandle(f *os.File) *IOHandle {
	return &IOHandle{name: f.Name(), file: f, writer: f, reader: f}
}

// NewWriterHandle wraps an arbitrary io.Writer. No reads.
func NewWriterHandle(name string, w io.Writer) *IOHandle {
	return &IOHandle{name: name, writer: w}
}

// NewReaderHandle wraps an arbitrary io.Reader. No writes.
func NewReaderHandle(name string, r io.Reader) *IOHandle {
	return &IOHandle{name: name, reader: r}
}

// File returns the underlying *os.File or nil when the handle wraps an
// arbitrary writer/reader. Callers that need Fd()/Sync() (e.g. term/raw
// mode) must handle nil.
func (h *IOHandle) File() *os.File { return h.file }

// Writer / Reader are the interface-shaped accessors the print/error
// refactor consults.
func (h *IOHandle) Writer() io.Writer    { return h.writer }
func (h *IOHandle) ReaderRaw() io.Reader { return h.reader }

// Write is the convenience write-string path used by the print fns and
// (write! handle x). Returns an error if the handle isn't writable.
func (h *IOHandle) Write(s string) (int, error) {
	if h.writer == nil {
		return 0, fmt.Errorf("IOHandle %q is not writable", h.name)
	}
	return io.WriteString(h.writer, s)
}

func (h *IOHandle) String() string {
	return fmt.Sprintf("#<IOHandle %s>", h.name)
}

// Reader returns a buffered reader over the underlying io.Reader, lazily
// constructed on first use. Returns nil for write-only handles.
func (h *IOHandle) Reader() *bufio.Reader {
	if h.reader == nil {
		return nil
	}
	if h.bufrd == nil {
		h.bufrd = bufio.NewReader(h.reader)
	}
	return h.bufrd
}

// Close closes the underlying file if file-backed, or any wrapped writer/
// reader that implements io.Closer. No-op for non-closable handles like
// bytes.Buffer or os.Stdout.
func (h *IOHandle) Close() error {
	if h.file != nil {
		return h.file.Close()
	}
	if c, ok := h.writer.(io.Closer); ok {
		return c.Close()
	}
	if c, ok := h.reader.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// Sync flushes/syncs the handle. For file-backed handles, calls File.Sync.
// For writers that implement a Flush() error method (e.g. bufio.Writer),
// calls that. Otherwise no-op (buffers are eager).
func (h *IOHandle) Sync() error {
	if h.file != nil {
		return h.file.Sync()
	}
	if f, ok := h.writer.(interface{ Flush() error }); ok {
		return f.Flush()
	}
	return nil
}

// getIOHandle extracts an *IOHandle from a Boxed value or wraps a raw *os.File.
func getIOHandle(v vm.Value) (*IOHandle, error) {
	b, ok := v.(*vm.Boxed)
	if !ok {
		return nil, fmt.Errorf("expected IOHandle, got %s", v.Type().Name())
	}
	switch u := b.Unbox().(type) {
	case *IOHandle:
		return u, nil
	case *os.File:
		return NewIOHandle(u), nil
	}
	return nil, fmt.Errorf("expected IOHandle, got %T", b.Unbox())
}

// nolint
func installIOBuiltins(ns *vm.Namespace) {
	// open — (open path mode) → IOHandle
	openf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("open expects 1-2 args")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("open expected String path")
		}
		mode := vm.Keyword("read")
		if len(vs) == 2 {
			m, ok := vs[1].(vm.Keyword)
			if !ok {
				return vm.NIL, fmt.Errorf("open expected Keyword mode")
			}
			mode = m
		}
		var flag int
		switch mode {
		case "read":
			flag = os.O_RDONLY
		case "write":
			flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		case "append":
			flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
		case "rw":
			flag = os.O_RDWR | os.O_CREATE
		default:
			return vm.NIL, fmt.Errorf("open: unknown mode %s", mode)
		}
		f, err := os.OpenFile(string(path), flag, 0644)
		if err != nil {
			return vm.NIL, err
		}
		h := NewIOHandle(f)
		return vm.NewBoxed(h), nil
	})

	// close! — (close! handle-or-chan) — works on IO handles and channels
	closef, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("close! expects 1 arg")
		}
		// Channel?
		if ch, ok := vs[0].(vm.Chan); ok {
			close(ch)
			return vm.NIL, nil
		}
		// IO handle
		h, err := getIOHandle(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.NIL, h.Close()
	})

	// read-line — (read-line handle) → String or nil at EOF
	readLine, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("read-line expects 1 arg")
		}
		h, err := getIOHandle(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		r := h.Reader()
		if r == nil {
			return vm.NIL, fmt.Errorf("read-line: handle %q is not readable", h.name)
		}
		line, err := r.ReadString('\n')
		if err != nil {
			if len(line) > 0 {
				// Return what we have even on EOF
				if line[len(line)-1] == '\n' {
					line = line[:len(line)-1]
				}
				return vm.String(line), nil
			}
			return vm.NIL, nil // EOF
		}
		// Strip trailing newline
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		return vm.String(line), nil
	})

	// write! — (write! handle str)
	writef, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("write! expects 2 args")
		}
		h, err := getIOHandle(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		var s string
		if str, ok := vs[1].(vm.String); ok {
			s = string(str)
		} else {
			s = vs[1].String()
		}
		_, err = h.Write(s)
		return vm.NIL, err
	})

	// flush! — (flush! handle)
	flushf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("flush! expects 1 arg")
		}
		h, err := getIOHandle(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.NIL, h.Sync()
	})

	// read-bytes — (read-bytes handle n) → String or nil at EOF
	readBytes, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("read-bytes expects 2 args")
		}
		h, err := getIOHandle(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		n, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("read-bytes expected Int count")
		}
		if h.reader == nil {
			return vm.NIL, fmt.Errorf("read-bytes: handle is not readable")
		}
		buf := make([]byte, int(n))
		nread, err := h.reader.Read(buf)
		if nread == 0 {
			return vm.NIL, nil // EOF
		}
		return vm.String(buf[:nread]), nil
	})

	// file-exists? — (file-exists? path)
	fileExists, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("file-exists? expects 1 arg")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("file-exists? expected String path")
		}
		_, err := os.Stat(string(path))
		if err != nil {
			if os.IsNotExist(err) {
				return vm.FALSE, nil
			}
			return vm.NIL, err
		}
		return vm.TRUE, nil
	})

	// delete-file — (delete-file path)
	deleteFile, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("delete-file expects 1 arg")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("delete-file expected String path")
		}
		return vm.NIL, os.Remove(string(path))
	})

	// mkdir — (mkdir path)
	mkdirf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("mkdir expects 1 arg")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("mkdir expected String path")
		}
		return vm.NIL, os.MkdirAll(string(path), 0755)
	})

	// *in*, *out*, *err* — stdin, stdout, stderr as IOHandle
	stdinHandle := vm.NewBoxed(NewIOHandle(os.Stdin))
	stdoutHandle := vm.NewBoxed(NewIOHandle(os.Stdout))
	stderrHandle := vm.NewBoxed(NewIOHandle(os.Stderr))

	ns.Def("open", openf)
	ns.Def("close!", closef)
	ns.Def("read-line", readLine)
	ns.Def("write!", writef)
	ns.Def("flush!", flushf)
	ns.Def("read-bytes", readBytes)
	ns.Def("file-exists?", fileExists)
	ns.Def("delete-file", deleteFile)
	ns.Def("mkdir", mkdirf)
	ns.Def("*in*", stdinHandle)
	ns.Def("*out*", stdoutHandle)
	ns.Def("*err*", stderrHandle)
}

// resolveIOHandleVar looks up a var (e.g. "*out*") in the core namespace
// and unwraps its current binding to an *IOHandle. Returns nil if the var
// isn't installed yet (e.g. during early init) or the current binding
// doesn't unwrap to a handle / *os.File. Callers must fall back to a sane
// default in that case.
//
// Resolution path: ns.LookupLocal -> Var.Deref (lock-free atomic load,
// respects (binding [...] ...) per pkg/vm/var.go:95-103) -> Boxed.Unbox.
func resolveIOHandleVar(varName string) *IOHandle {
	ns := lookupNSCached(NameCoreNS)
	if ns == nil {
		return nil
	}
	v := ns.LookupLocal(vm.Symbol(varName))
	if v == nil {
		return nil
	}
	b, ok := v.Deref().(*vm.Boxed)
	if !ok {
		return nil
	}
	switch u := b.Unbox().(type) {
	case *IOHandle:
		return u
	case *os.File:
		return NewIOHandle(u)
	}
	return nil
}

// WriteToOut writes s through the current dynamic binding of *out*. Falls
// back to os.Stdout if *out* isn't installed yet or doesn't resolve to a
// handle (early-boot conditions only).
func WriteToOut(s string) error {
	if h := resolveIOHandleVar("*out*"); h != nil {
		_, err := h.Write(s)
		return err
	}
	_, err := io.WriteString(os.Stdout, s)
	return err
}

// WriteToErr writes s through the current dynamic binding of *err*. Falls
// back to os.Stderr.
func WriteToErr(s string) error {
	if h := resolveIOHandleVar("*err*"); h != nil {
		_, err := h.Write(s)
		return err
	}
	_, err := io.WriteString(os.Stderr, s)
	return err
}

// LookupCoreVar returns the *Var for a name in the core namespace, or nil
// if not installed. Exported for packages (e.g. nrepl) that need to
// PushBinding/PopBinding on the core I/O vars to implement scoped capture.
func LookupCoreVar(varName string) *vm.Var {
	ns := lookupNSCached(NameCoreNS)
	if ns == nil {
		return nil
	}
	return ns.LookupLocal(vm.Symbol(varName))
}
