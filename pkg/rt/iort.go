/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nooga/let-go/pkg/vm"
)

// IOHandle wraps an *os.File with optional buffered reader for line-based reads.
type IOHandle struct {
	File   *os.File
	reader *bufio.Reader
}

func NewIOHandle(f *os.File) *IOHandle {
	return &IOHandle{File: f}
}

func (h *IOHandle) String() string {
	return fmt.Sprintf("#<IOHandle %s>", h.File.Name())
}

func (h *IOHandle) Reader() *bufio.Reader {
	if h.reader == nil {
		h.reader = bufio.NewReader(h.File)
	}
	return h.reader
}

func (h *IOHandle) Close() error {
	return h.File.Close()
}

// getIOHandle extracts an *IOHandle from a Boxed value or wraps a raw *os.File.
func getIOHandle(v vm.Value) (*IOHandle, error) {
	b, ok := v.(*vm.Boxed)
	if !ok {
		return nil, fmt.Errorf("expected IOHandle, got %s", v.Type().Name())
	}
	if h, ok := b.Unbox().(*IOHandle); ok {
		return h, nil
	}
	if f, ok := b.Unbox().(*os.File); ok {
		return NewIOHandle(f), nil
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
		line, err := h.Reader().ReadString('\n')
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
		_, err = h.File.WriteString(s)
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
		return vm.NIL, h.File.Sync()
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
		buf := make([]byte, int(n))
		nread, err := h.File.Read(buf)
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
