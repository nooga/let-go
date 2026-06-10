//go:build linux

/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

/*
 * unix namespace — AF_UNIX stream sockets with SCM_RIGHTS fd passing.
 *
 * Just enough primitives to build a control-plane over a unix socket:
 *   (unix/listen path)         → listener
 *   (unix/accept listener)     → conn (blocks)
 *   (unix/connect path)        → conn
 *   (unix/send conn data)
 *   (unix/send conn data fds)  → nil — fds may be IOHandles or raw Ints
 *   (unix/recv conn max-bytes max-fds) → {:data "..." :fds [IOHandle ...]}
 *   (unix/close listener-or-conn)
 *   (unix/fd handle)           → Int — extract raw fd from IOHandle
 *
 * Received fds are wrapped as IOHandles so they drop straight into
 * spawn-async's stdio slots.
 */

package rt

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/nooga/let-go/pkg/vm"
)

func init() { RegisterInstaller(installUnixNS) }

func installUnixNS() {
	unboxConn := func(v vm.Value) (*net.UnixConn, error) {
		b, ok := v.(*vm.Boxed)
		if !ok {
			return nil, fmt.Errorf("expected unix conn, got %s", v.Type().Name())
		}
		c, ok := b.Unbox().(*net.UnixConn)
		if !ok {
			return nil, fmt.Errorf("expected unix conn, got %T", b.Unbox())
		}
		return c, nil
	}

	unboxListener := func(v vm.Value) (*net.UnixListener, error) {
		b, ok := v.(*vm.Boxed)
		if !ok {
			return nil, fmt.Errorf("expected unix listener, got %s", v.Type().Name())
		}
		l, ok := b.Unbox().(*net.UnixListener)
		if !ok {
			return nil, fmt.Errorf("expected unix listener, got %T", b.Unbox())
		}
		return l, nil
	}

	listenFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("unix/listen expects 1 arg (path)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("unix/listen expected String path")
		}
		_ = os.Remove(string(path))
		l, err := net.ListenUnix("unixpacket", &net.UnixAddr{Name: string(path), Net: "unixpacket"})
		if err != nil {
			return vm.NIL, fmt.Errorf("unix/listen: %v", err)
		}
		return vm.NewBoxed(l), nil
	})

	acceptFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("unix/accept expects 1 arg (listener)")
		}
		l, err := unboxListener(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("unix/accept: %v", err)
		}
		c, err := l.AcceptUnix()
		if err != nil {
			return vm.NIL, fmt.Errorf("unix/accept: %v", err)
		}
		return vm.NewBoxed(c), nil
	})

	connectFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("unix/connect expects 1 arg (path)")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("unix/connect expected String path")
		}
		c, err := net.DialUnix("unixpacket", nil, &net.UnixAddr{Name: string(path), Net: "unixpacket"})
		if err != nil {
			return vm.NIL, fmt.Errorf("unix/connect: %v", err)
		}
		return vm.NewBoxed(c), nil
	})

	extractFd := func(v vm.Value) (int, error) {
		switch fv := v.(type) {
		case vm.Int:
			return int(fv), nil
		case *vm.Boxed:
			if h, ok := fv.Unbox().(*IOHandle); ok {
				f := h.File()
				if f == nil {
					return 0, fmt.Errorf("IOHandle is not file-backed; fd ops require a real file descriptor")
				}
				return int(f.Fd()), nil
			}
			if f, ok := fv.Unbox().(*os.File); ok {
				return int(f.Fd()), nil
			}
			return 0, fmt.Errorf("unsupported fd value %T", fv.Unbox())
		}
		return 0, fmt.Errorf("unsupported fd value %T", v)
	}

	sendFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("unix/send expects 2-3 args (conn data [fds])")
		}
		c, err := unboxConn(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("unix/send: %v", err)
		}
		data, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("unix/send expected String data")
		}
		var oob []byte
		if len(vs) == 3 && vs[2] != vm.NIL {
			fdsSeq, ok := vs[2].(vm.Sequable)
			if !ok {
				return vm.NIL, fmt.Errorf("unix/send expected Sequable fds")
			}
			var fds []int
			for s := fdsSeq.Seq(); s != nil && s != vm.EmptyList; s = s.Next() {
				fd, err := extractFd(s.First())
				if err != nil {
					return vm.NIL, fmt.Errorf("unix/send: %v", err)
				}
				fds = append(fds, fd)
			}
			if len(fds) > 0 {
				oob = syscall.UnixRights(fds...)
			}
		}
		if _, _, err := c.WriteMsgUnix([]byte(string(data)), oob, nil); err != nil {
			return vm.NIL, fmt.Errorf("unix/send: %v", err)
		}
		return vm.NIL, nil
	})

	recvFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("unix/recv expects 3 args (conn max-bytes max-fds)")
		}
		c, err := unboxConn(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("unix/recv: %v", err)
		}
		maxBytes, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("unix/recv expected Int max-bytes")
		}
		maxFds, ok := vs[2].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("unix/recv expected Int max-fds")
		}
		buf := make([]byte, int(maxBytes))
		var oob []byte
		if maxFds > 0 {
			oob = make([]byte, syscall.CmsgSpace(4*int(maxFds)))
		}
		n, oobn, _, _, err := c.ReadMsgUnix(buf, oob)
		if err != nil {
			return vm.NIL, fmt.Errorf("unix/recv: %v", err)
		}
		var fdVals []vm.Value
		if oobn > 0 {
			msgs, err := syscall.ParseSocketControlMessage(oob[:oobn])
			if err != nil {
				return vm.NIL, fmt.Errorf("unix/recv parse: %v", err)
			}
			for _, m := range msgs {
				ff, err := syscall.ParseUnixRights(&m)
				if err != nil {
					return vm.NIL, fmt.Errorf("unix/recv parse rights: %v", err)
				}
				for _, fd := range ff {
					f := os.NewFile(uintptr(fd), fmt.Sprintf("scm-%d", fd))
					fdVals = append(fdVals, vm.NewBoxed(NewIOHandle(f)))
				}
			}
		}
		m := vm.EmptyPersistentMap
		m = m.Assoc(vm.Keyword("data"), vm.String(buf[:n])).(*vm.PersistentMap)
		m = m.Assoc(vm.Keyword("fds"), vm.NewArrayVector(fdVals)).(*vm.PersistentMap)
		return m, nil
	})

	closeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("unix/close expects 1 arg")
		}
		b, ok := vs[0].(*vm.Boxed)
		if !ok {
			return vm.NIL, fmt.Errorf("unix/close expected listener or conn")
		}
		switch v := b.Unbox().(type) {
		case *net.UnixListener:
			_ = v.Close()
		case *net.UnixConn:
			_ = v.Close()
		case *IOHandle:
			_ = v.Close()
		default:
			return vm.NIL, fmt.Errorf("unix/close: unsupported type %T", v)
		}
		return vm.NIL, nil
	})

	fdFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("unix/fd expects 1 arg")
		}
		fd, err := extractFd(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("unix/fd: %v", err)
		}
		return vm.MakeInt(fd), nil
	})

	ns := vm.NewNamespace("unix")
	ns.Def("listen", listenFn)
	ns.Def("accept", acceptFn)
	ns.Def("connect", connectFn)
	ns.Def("send", sendFn)
	ns.Def("recv", recvFn)
	ns.Def("close", closeFn)
	ns.Def("fd", fdFn)
	RegisterNS(ns)
}
