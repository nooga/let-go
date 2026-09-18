//go:build !js

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

/*
 * net namespace — minimal TCP client and owned listener primitives.
 *
 *   (net/dial host port)          → conn (Boxed). 3s connect timeout.
 *   (net/write! conn str-or-bytes) → nil. Raw bytes on the wire.
 *   (net/read! conn max-bytes)    → byte-array, nil on clean EOF.
 *   (net/close! conn)             → nil.
 *   (net/listen host port)        → listener; port zero asks the OS.
 *   (net/local-address listener) → {:host string :port integer}.
 *   (net/accept listener)         → conn, compatible with bencode.
 *   (net/close! listener)         → nil; accepted connections remain open.
 *
 * Native only — js/wasm gets a stub (net_wasm.go), because Go's GOOS=js net
 * stack is an in-process fake that would connect to nothing. Plain TCP here.
 * The Boxed conn also carries a lazily-created bencode.Decoder so the bencode
 * namespace (net.go's sibling bencode.go) can frame messages over the same
 * buffered stream.
 */

package rt

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/nooga/let-go/pkg/vm"
	"github.com/zeebo/bencode"
)

func init() { RegisterInstaller(installNetNS) }

// netConn wraps a live TCP connection. dec is the persistent bencode.Decoder
// the bencode namespace attaches on first read (see bencode.go); it wraps a
// bufio.Reader, so read-ahead bytes past one message must survive between
// reads — hence it lives on the conn, not per-call. readMu serializes that
// stateful read path so concurrent bencode/read! calls can't corrupt the
// decoder; writes are not locked, so a blocked read never stalls a writer.
type netConn struct {
	conn   net.Conn
	dec    *bencode.Decoder
	readMu sync.Mutex
}

type netListener struct{ listener net.Listener }

func unboxNetListener(v vm.Value) (*netListener, error) {
	b, ok := v.(*vm.Boxed)
	if !ok {
		return nil, fmt.Errorf("expected net listener")
	}
	l, ok := b.Unbox().(*netListener)
	if !ok {
		return nil, fmt.Errorf("expected net listener, got %T", b.Unbox())
	}
	return l, nil
}

func unboxNetConn(v vm.Value) (*netConn, error) {
	b, ok := v.(*vm.Boxed)
	if !ok {
		return nil, fmt.Errorf("expected net conn, got %s", v.Type().Name())
	}
	c, ok := b.Unbox().(*netConn)
	if !ok {
		return nil, fmt.Errorf("expected net conn, got %T", b.Unbox())
	}
	return c, nil
}

func installNetNS() {
	listenFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("net/listen expects 2 args (host port)")
		}
		host, ok := vs[0].(vm.String)
		if !ok || host == "" {
			return vm.NIL, fmt.Errorf("net/listen expected nonempty String host")
		}
		port, ok := vs[1].(vm.Int)
		if !ok || port < 0 || port > 65535 {
			return vm.NIL, fmt.Errorf("net/listen expected Int port between 0 and 65535")
		}
		listener, err := net.Listen("tcp", net.JoinHostPort(string(host), strconv.Itoa(int(port))))
		if err != nil {
			return vm.NIL, fmt.Errorf("net/listen: %w", err)
		}
		return vm.NewBoxed(&netListener{listener: listener}), nil
	})
	addressFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("net/local-address expects 1 arg (listener)")
		}
		l, err := unboxNetListener(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("net/local-address: %w", err)
		}
		address, ok := l.listener.Addr().(*net.TCPAddr)
		if !ok {
			return vm.NIL, fmt.Errorf("net/local-address: expected TCP address")
		}
		return vm.EmptyPersistentMap.Assoc(vm.Keyword("host"), vm.String(address.IP.String())).(*vm.PersistentMap).Assoc(vm.Keyword("port"), vm.Int(address.Port)), nil
	})
	acceptFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("net/accept expects 1 arg (listener)")
		}
		l, err := unboxNetListener(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("net/accept: %w", err)
		}
		conn, err := l.listener.Accept()
		if err != nil {
			return vm.NIL, fmt.Errorf("net/accept: %w", err)
		}
		return vm.NewBoxed(&netConn{conn: conn}), nil
	})
	dialFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("net/dial expects 2 args (host port)")
		}
		host, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("net/dial expected String host")
		}
		port, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("net/dial expected Int port")
		}
		addr := net.JoinHostPort(string(host), strconv.Itoa(int(port)))
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err != nil {
			return vm.NIL, fmt.Errorf("net/dial: %v", err)
		}
		return vm.NewBoxed(&netConn{conn: conn}), nil
	})

	writeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("net/write! expects 2 args (conn data)")
		}
		c, err := unboxNetConn(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("net/write!: %v", err)
		}
		data, ok := asBytes(vs[1])
		if !ok {
			return vm.NIL, fmt.Errorf("net/write! expected String or byte-array data")
		}
		if _, err := c.conn.Write(data); err != nil {
			return vm.NIL, fmt.Errorf("net/write!: %v", err)
		}
		return vm.NIL, nil
	})

	readFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("net/read! expects 2 args (conn max-bytes)")
		}
		c, err := unboxNetConn(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("net/read!: %v", err)
		}
		maxBytes, ok := vs[1].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("net/read! expected Int max-bytes")
		}
		if maxBytes <= 0 {
			return vm.NIL, fmt.Errorf("net/read! expected positive Int max-bytes, got %d", int(maxBytes))
		}
		buf := make([]byte, int(maxBytes))
		n, err := c.conn.Read(buf)
		if n > 0 {
			// Return the bytes read; a coincident EOF surfaces on the next read.
			return vm.NewByteArrayFrom(buf[:n]), nil
		}
		if err == io.EOF {
			return vm.NIL, nil
		}
		if err != nil {
			return vm.NIL, fmt.Errorf("net/read!: %v", err)
		}
		return vm.NewByteArrayFrom(nil), nil
	})

	closeFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("net/close! expects 1 arg (conn or listener)")
		}
		if l, err := unboxNetListener(vs[0]); err == nil {
			_ = l.listener.Close()
			return vm.NIL, nil
		}
		c, err := unboxNetConn(vs[0])
		if err != nil {
			return vm.NIL, fmt.Errorf("net/close!: %v", err)
		}
		_ = c.conn.Close()
		return vm.NIL, nil
	})

	ns := vm.NewNamespace("net")
	ns.Def("listen", listenFn)
	ns.Def("local-address", addressFn)
	ns.Def("accept", acceptFn)
	ns.Def("dial", dialFn)
	ns.Def("write!", writeFn)
	ns.Def("read!", readFn)
	ns.Def("close!", closeFn)
	RegisterNS(ns)
}
