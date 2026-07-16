/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"bytes"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
	"github.com/zeebo/bencode"
)

// nsFn looks up a namespace var and returns it as an Fn, following the
// async_test.go / urlparam_test.go convention for exercising an rt namespace
// from within package rt (a package rt test cannot import pkg/compiler to
// eval lg source — that would be an import cycle).
func nsFn(t *testing.T, ns, name string) vm.Fn {
	t.Helper()
	vr, ok := NS(ns).Lookup(vm.Symbol(name)).(*vm.Var)
	if !ok {
		t.Fatalf("%s/%s not registered", ns, name)
	}
	fn, ok := vr.Deref().(vm.Fn)
	if !ok {
		t.Fatalf("%s/%s is not an Fn", ns, name)
	}
	return fn
}

// dialPair opens a listener, net/dials it from lg, and returns the lg-side
// conn value together with the accepted server-side net.Conn. cleanup closes
// everything. Used by the bencode tests to drive a real TCP round-trip.
func dialPair(t *testing.T) (connV vm.Value, srv net.Conn, cleanup func()) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := l.Addr().(*net.TCPAddr)
	accepted := make(chan net.Conn, 1)
	acceptErr := make(chan error, 1)
	go func() {
		c, err := l.Accept()
		if err != nil {
			acceptErr <- err
			return
		}
		accepted <- c
	}()

	connV, err = nsFn(t, "net", "dial").Invoke([]vm.Value{vm.String(addr.IP.String()), vm.Int(addr.Port)})
	if err != nil {
		l.Close()
		t.Fatalf("net/dial: %v", err)
	}
	select {
	case srv = <-accepted:
	case err := <-acceptErr:
		l.Close()
		t.Fatalf("accept: %v", err)
	}
	cleanup = func() {
		nsFn(t, "net", "close!").Invoke([]vm.Value{connV})
		srv.Close()
		l.Close()
	}
	return connV, srv, cleanup
}

// TestBencodeWriteRead round-trips an nREPL-shaped message: write! a request
// map from lg (decoded on the Go side), then read! a response dict with a
// nested list back into an lg map with String keys and a vector value.
func TestBencodeWriteRead(t *testing.T) {
	connV, srv, cleanup := dialPair(t)
	defer cleanup()

	bWrite := nsFn(t, "bencode", "write!")
	bRead := nsFn(t, "bencode", "read!")

	req := vm.NewPersistentMap([]vm.Value{
		vm.String("op"), vm.String("eval"),
		vm.String("code"), vm.String("(+ 1 1)"),
	})
	if _, err := bWrite.Invoke([]vm.Value{connV, req}); err != nil {
		t.Fatalf("bencode/write!: %v", err)
	}
	var got map[string]interface{}
	if err := bencode.NewDecoder(srv).Decode(&got); err != nil {
		t.Fatalf("server decode: %v", err)
	}
	if got["op"] != "eval" || got["code"] != "(+ 1 1)" {
		t.Fatalf("server decoded %#v", got)
	}

	resp, err := bencode.EncodeBytes(map[string]interface{}{
		"status": []interface{}{"done"},
		"id":     "1",
		"ns":     "user",
	})
	if err != nil {
		t.Fatalf("encode resp: %v", err)
	}
	if _, err := srv.Write(resp); err != nil {
		t.Fatalf("server write: %v", err)
	}

	v, err := bRead.Invoke([]vm.Value{connV})
	if err != nil {
		t.Fatalf("bencode/read!: %v", err)
	}
	m, ok := v.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("bencode/read! returned %T, want *vm.PersistentMap", v)
	}
	if id := m.ValueAt(vm.String("id")); id != vm.String("1") {
		t.Fatalf("id = %v, want \"1\"", id)
	}
	st := m.ValueAt(vm.String("status"))
	vec, ok := st.(vm.ArrayVector)
	if !ok {
		t.Fatalf("status = %T, want a vector", st)
	}
	if vec.Count() != vm.Int(1) || vec.ValueAt(vm.Int(0)) != vm.String("done") {
		t.Fatalf("status = %v, want [\"done\"]", st)
	}
}

// TestBencodeConversions pins the conversion table's edge cases: keyword keys
// become name strings, Ints round-trip, and nil/bool/float error on write.
func TestBencodeConversions(t *testing.T) {
	connV, srv, cleanup := dialPair(t)
	defer cleanup()

	bWrite := nsFn(t, "bencode", "write!")
	bRead := nsFn(t, "bencode", "read!")

	req := vm.NewPersistentMap([]vm.Value{
		vm.Keyword("op"), vm.String("clone"),
		vm.String("n"), vm.Int(42),
	})
	if _, err := bWrite.Invoke([]vm.Value{connV, req}); err != nil {
		t.Fatalf("bencode/write!: %v", err)
	}
	var got map[string]interface{}
	if err := bencode.NewDecoder(srv).Decode(&got); err != nil {
		t.Fatalf("server decode: %v", err)
	}
	if got["op"] != "clone" {
		t.Fatalf("keyword key not stringified: %#v", got)
	}
	if got["n"] != int64(42) {
		t.Fatalf("int value = %#v (%T), want int64(42)", got["n"], got["n"])
	}

	// Int round-trips back into lg as an Int.
	resp, _ := bencode.EncodeBytes(map[string]interface{}{"n": int64(7)})
	if _, err := srv.Write(resp); err != nil {
		t.Fatalf("server write: %v", err)
	}
	v, err := bRead.Invoke([]vm.Value{connV})
	if err != nil {
		t.Fatalf("bencode/read!: %v", err)
	}
	if n := v.(*vm.PersistentMap).ValueAt(vm.String("n")); n != vm.Int(7) {
		t.Fatalf("round-trip int = %v (%T), want Int(7)", n, n)
	}

	// nil, bool and float have no bencode representation → write! rejects them.
	for _, bad := range []vm.Value{vm.NIL, vm.TRUE, vm.Float(1.5)} {
		if _, err := bWrite.Invoke([]vm.Value{connV, bad}); err == nil {
			t.Fatalf("bencode/write! accepted %s, want error", bad.Type().Name())
		}
	}
}

// TestBencodeTimeout verifies a timed read against a silent server errors with
// "timeout", and that the read deadline is cleared so a later read still works.
func TestBencodeTimeout(t *testing.T) {
	connV, srv, cleanup := dialPair(t)
	defer cleanup()

	bRead := nsFn(t, "bencode", "read!")

	_, err := bRead.Invoke([]vm.Value{connV, vm.Int(100)})
	if err == nil {
		t.Fatal("bencode/read! with timeout against a silent server: expected error")
	}
	if !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("timeout error missing \"timeout\": %v", err)
	}

	// Deadline cleared: a subsequent blocking read succeeds.
	resp, _ := bencode.EncodeBytes(map[string]interface{}{"status": []interface{}{"done"}})
	if _, err := srv.Write(resp); err != nil {
		t.Fatalf("server write: %v", err)
	}
	v, err := bRead.Invoke([]vm.Value{connV})
	if err != nil {
		t.Fatalf("bencode/read! after timeout: %v", err)
	}
	m, ok := v.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("bencode/read! returned %T, want *vm.PersistentMap", v)
	}
	if _, ok := m.ValueAt(vm.String("status")).(vm.ArrayVector); !ok {
		t.Fatalf("status = %T, want a vector", m.ValueAt(vm.String("status")))
	}
}

// TestNet exercises the net namespace against a real net.Listener: dial,
// write! (string and byte-array with a high byte, verbatim on the wire),
// read! (bytes back, then nil on EOF), and close!.
func TestNet(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	addr := l.Addr().(*net.TCPAddr)

	type accepted struct {
		c   net.Conn
		err error
	}
	ch := make(chan accepted, 1)
	go func() {
		c, err := l.Accept()
		ch <- accepted{c, err}
	}()

	dial := nsFn(t, "net", "dial")
	write := nsFn(t, "net", "write!")
	read := nsFn(t, "net", "read!")
	closeFn := nsFn(t, "net", "close!")

	connV, err := dial.Invoke([]vm.Value{vm.String(addr.IP.String()), vm.Int(addr.Port)})
	if err != nil {
		t.Fatalf("net/dial: %v", err)
	}

	a := <-ch
	if a.err != nil {
		t.Fatalf("accept: %v", a.err)
	}
	srv := a.c
	defer srv.Close()

	// write! a String — raw bytes on the wire.
	if _, err := write.Invoke([]vm.Value{connV, vm.String("hello")}); err != nil {
		t.Fatalf("net/write! string: %v", err)
	}
	buf := make([]byte, 5)
	if _, err := io.ReadFull(srv, buf); err != nil {
		t.Fatalf("server read string: %v", err)
	}
	if string(buf) != "hello" {
		t.Fatalf("server got %q, want %q", buf, "hello")
	}

	// write! a byte-array containing a >127 byte — must arrive verbatim.
	payload := []byte{0x00, 0xFF, 0x7F, 0x80}
	if _, err := write.Invoke([]vm.Value{connV, vm.NewByteArrayFrom(payload)}); err != nil {
		t.Fatalf("net/write! bytes: %v", err)
	}
	buf2 := make([]byte, len(payload))
	if _, err := io.ReadFull(srv, buf2); err != nil {
		t.Fatalf("server read bytes: %v", err)
	}
	if !bytes.Equal(buf2, payload) {
		t.Fatalf("server got % x, want % x", buf2, payload)
	}

	// server writes back → net/read! returns a byte-array with the raw bytes.
	back := []byte{1, 2, 200, 3}
	if _, err := srv.Write(back); err != nil {
		t.Fatalf("server write: %v", err)
	}
	got, err := read.Invoke([]vm.Value{connV, vm.Int(64)})
	if err != nil {
		t.Fatalf("net/read!: %v", err)
	}
	gotBytes, ok := asBytes(got)
	if !ok {
		t.Fatalf("net/read! did not return a byte-array: %T", got)
	}
	if !bytes.Equal(gotBytes, back) {
		t.Fatalf("net/read! got % x, want % x", gotBytes, back)
	}

	// server closes → net/read! returns nil on EOF.
	srv.Close()
	eofV, err := read.Invoke([]vm.Value{connV, vm.Int(64)})
	if err != nil {
		t.Fatalf("net/read! after close: %v", err)
	}
	if eofV != vm.NIL {
		t.Fatalf("net/read! at EOF returned %v, want nil", eofV)
	}

	// close! succeeds.
	if _, err := closeFn.Invoke([]vm.Value{connV}); err != nil {
		t.Fatalf("net/close!: %v", err)
	}
}

// TestNetReadNegativeSize verifies a non-positive max-bytes is a normal
// net/read! argument error, not an internal makeslice panic.
func TestNetReadNegativeSize(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	addr := l.Addr().(*net.TCPAddr)
	go func() {
		if c, err := l.Accept(); err == nil {
			defer c.Close()
			io.Copy(io.Discard, c)
		}
	}()

	dial := nsFn(t, "net", "dial")
	read := nsFn(t, "net", "read!")
	connV, err := dial.Invoke([]vm.Value{vm.String(addr.IP.String()), vm.Int(addr.Port)})
	if err != nil {
		t.Fatalf("net/dial: %v", err)
	}
	defer nsFn(t, "net", "close!").Invoke([]vm.Value{connV})

	for _, n := range []int64{-1, 0} {
		_, err := read.Invoke([]vm.Value{connV, vm.Int(n)})
		if err == nil {
			t.Fatalf("net/read! max-bytes=%d: expected error, got nil", n)
		}
		if !strings.Contains(err.Error(), "net/read!") {
			t.Fatalf("net/read! max-bytes=%d error missing fn prefix: %v", n, err)
		}
	}
}

// TestNetDialClosedPort verifies dialing a port nothing listens on errors,
// with the fn name in the message.
func TestNetDialClosedPort(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := l.Addr().(*net.TCPAddr)
	l.Close() // free the port so the dial below is refused

	dial := nsFn(t, "net", "dial")
	_, err = dial.Invoke([]vm.Value{vm.String("127.0.0.1"), vm.Int(addr.Port)})
	if err == nil {
		t.Fatal("net/dial to a closed port: expected error, got nil")
	}
	if !strings.Contains(err.Error(), "net/dial") {
		t.Fatalf("net/dial error missing fn prefix: %v", err)
	}
}
