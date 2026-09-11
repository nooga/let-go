//go:build !tinygo && !lg_no_http

package rt

import (
	"context"
	"errors"
	"io"
	"net"
	"net/url"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

func requestOpts(url string, stream bool, timeout vm.Value) vm.Value {
	m := vm.EmptyPersistentMap.
		Assoc(vm.Keyword("url"), vm.String(url)).
		Assoc(vm.Keyword("timeout"), timeout).(*vm.PersistentMap)
	if stream {
		m = m.Assoc(vm.Keyword("as"), vm.Keyword("stream")).(*vm.PersistentMap)
	}
	return m
}

// invokeRequest runs http/request in a child scope that stays open until
// the test finishes, so body reads see only the request's own timeouts.
func invokeRequest(t *testing.T, opts vm.Value) (vm.Value, error) {
	t.Helper()
	scope := vm.Goroutines.Child()
	t.Cleanup(func() { vm.CloseScoped(scope, time.Second) })
	ec := vm.NewExecContext()
	ec.SetScope(scope)
	fn := LookupVar("http", "request").Deref().(vm.Fn)
	return ec.Invoke(fn, []vm.Value{opts})
}

func TestHTTPRequestTimeoutScopes(t *testing.T) {
	t.Run("request scope covers a slow response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(400 * time.Millisecond)
			_, _ = w.Write([]byte("late"))
		}))
		defer server.Close()
		timeouts := vm.EmptyPersistentMap.Assoc(vm.Keyword("request"), vm.Float(0.1))
		started := time.Now()
		_, err := invokeRequest(t, requestOpts(server.URL, false, timeouts))
		if err == nil || !strings.Contains(err.Error(), "http request timeout after 100ms") {
			t.Fatalf("expected request timeout, got %v", err)
		}
		if time.Since(started) > 350*time.Millisecond {
			t.Fatalf("timeout did not fire promptly: %s", time.Since(started))
		}
	})

	t.Run("bare number and timeout_ms mean request scope", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(300 * time.Millisecond)
		}))
		defer server.Close()
		if _, err := invokeRequest(t, requestOpts(server.URL, false, vm.Float(0.05))); err == nil || !strings.Contains(err.Error(), "request timeout") {
			t.Fatalf("bare number: %v", err)
		}
		opts := vm.EmptyPersistentMap.
			Assoc(vm.Keyword("url"), vm.String(server.URL)).
			Assoc(vm.Keyword("timeout_ms"), vm.Int(50))
		if _, err := invokeRequest(t, opts); err == nil || !strings.Contains(err.Error(), "request timeout") {
			t.Fatalf("timeout_ms: %v", err)
		}
	})

	t.Run("stream_read scope fires on a stalled body, not on a slow start", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			flusher := w.(http.Flusher)
			_, _ = w.Write([]byte("first\n"))
			flusher.Flush()
			time.Sleep(500 * time.Millisecond)
			_, _ = w.Write([]byte("second\n"))
		}))
		defer server.Close()
		timeouts := vm.EmptyPersistentMap.
			Assoc(vm.Keyword("request"), vm.Float(5)).
			Assoc(vm.Keyword("stream_read"), vm.Float(0.1))
		v, err := invokeRequest(t, requestOpts(server.URL, true, timeouts))
		if err != nil {
			t.Fatal(err)
		}
		reader := v.(vm.Lookup).ValueAt(vm.Keyword("body")).Unbox().(*LGReader)
		defer reader.Close()
		buf := make([]byte, 64)
		n, err := reader.Read(buf)
		if err != nil || string(buf[:n]) != "first\n" {
			t.Fatalf("first chunk: %q %v", buf[:n], err)
		}
		_, err = reader.Read(buf)
		if err == nil || !strings.Contains(err.Error(), "http stream_read timeout after 100ms") {
			t.Fatalf("expected stream_read timeout, got %v", err)
		}
	})

	t.Run("request scope on a stream bounds the headers and the first chunk only", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			flusher := w.(http.Flusher)
			_, _ = w.Write([]byte("a"))
			flusher.Flush()
			time.Sleep(250 * time.Millisecond)
			_, _ = w.Write([]byte("b"))
		}))
		defer server.Close()
		timeouts := vm.EmptyPersistentMap.Assoc(vm.Keyword("request"), vm.Float(0.1))
		v, err := invokeRequest(t, requestOpts(server.URL, true, timeouts))
		if err != nil {
			t.Fatal(err)
		}
		reader := v.(vm.Lookup).ValueAt(vm.Keyword("body")).Unbox().(*LGReader)
		data, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil || string(data) != "ab" {
			t.Fatalf("later chunks are not bound by the request scope: %q %v", data, err)
		}
	})

	t.Run("a slow first chunk is the request scope, later gaps are stream_read", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			flusher := w.(http.Flusher)
			flusher.Flush() // headers at once, first byte late
			time.Sleep(300 * time.Millisecond)
			_, _ = w.Write([]byte("late-first"))
		}))
		defer server.Close()
		timeouts := vm.EmptyPersistentMap.
			Assoc(vm.Keyword("request"), vm.Float(1)).
			Assoc(vm.Keyword("stream_read"), vm.Float(0.1))
		v, err := invokeRequest(t, requestOpts(server.URL, true, timeouts))
		if err != nil {
			t.Fatal(err)
		}
		reader := v.(vm.Lookup).ValueAt(vm.Keyword("body")).Unbox().(*LGReader)
		data, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil || string(data) != "late-first" {
			t.Fatalf("first chunk slower than stream_read but within request must succeed: %q %v", data, err)
		}
		short := vm.EmptyPersistentMap.
			Assoc(vm.Keyword("request"), vm.Float(0.1)).
			Assoc(vm.Keyword("stream_read"), vm.Float(5))
		v, err = invokeRequest(t, requestOpts(server.URL, true, short))
		if err != nil {
			t.Fatal(err)
		}
		reader = v.(vm.Lookup).ValueAt(vm.Keyword("body")).Unbox().(*LGReader)
		_, err = io.ReadAll(reader)
		_ = reader.Close()
		if err == nil || !strings.Contains(err.Error(), "http request timeout after 100ms") {
			t.Fatalf("first chunk past the request scope must fail as request timeout, got %v", err)
		}
	})

	t.Run("connect scope fires on an unreachable host", func(t *testing.T) {
		timeouts := vm.EmptyPersistentMap.Assoc(vm.Keyword("connect"), vm.Float(0.1))
		started := time.Now()
		_, err := invokeRequest(t, requestOpts("http://10.255.255.1:9/", false, timeouts))
		if err == nil || !strings.Contains(err.Error(), "http connect timeout after 100ms") {
			t.Fatalf("expected connect timeout, got %v", err)
		}
		if time.Since(started) > time.Second {
			t.Fatalf("connect timeout did not fire promptly: %s", time.Since(started))
		}
	})
}

// scriptedBody plays one step per Read, so a test can make a read finish
// before or after its deadline deterministically.
type scriptedBody struct {
	steps []func(p []byte) (int, error)
	i     int
}

func (b *scriptedBody) Read(p []byte) (int, error) {
	if b.i >= len(b.steps) {
		return 0, io.EOF
	}
	step := b.steps[b.i]
	b.i++
	return step(p)
}

func (b *scriptedBody) Close() error { return nil }

// lateChunk returns data only after d, so a shorter deadline expires while
// the read is in progress but the read still succeeds.
func lateChunk(d time.Duration, data string) func([]byte) (int, error) {
	return func(p []byte) (int, error) {
		time.Sleep(d)
		return copy(p, data), nil
	}
}

// TestGapReaderReportsOnlyTheTimeoutThatFired covers nooga/let-go#857 item 2:
// the fired state belongs to the read whose deadline expired, a clean end of
// stream is never turned into a timeout, and a later failure caused by that
// timeout's cancellation reports the scope and duration that actually fired.
func TestGapReaderReportsOnlyTheTimeoutThatFired(t *testing.T) {
	noop := func() {}
	buf := make([]byte, 16)

	t.Run("end of stream after a racing deadline stays end of stream", func(t *testing.T) {
		g := &gapReader{body: &scriptedBody{steps: []func([]byte) (int, error){
			lateChunk(40*time.Millisecond, "a"),
			func([]byte) (int, error) { return 0, io.EOF },
		}}, first: 10 * time.Millisecond, cancel: noop}
		if n, err := g.Read(buf); n != 1 || err != nil {
			t.Fatalf("the first read must deliver its data: %d %v", n, err)
		}
		if _, err := g.Read(buf); err != io.EOF {
			t.Fatalf("a clean end of stream must stay io.EOF, got %v", err)
		}
	})

	t.Run("a later cancellation reports the scope and duration that fired", func(t *testing.T) {
		g := &gapReader{body: &scriptedBody{steps: []func([]byte) (int, error){
			lateChunk(40*time.Millisecond, "a"),
			func([]byte) (int, error) { return 0, context.Canceled },
		}}, first: 10 * time.Millisecond, cancel: noop}
		_, _ = g.Read(buf)
		_, err := g.Read(buf)
		if err == nil || err.Error() != "http request timeout after 10ms" {
			t.Fatalf("expected the request timeout that fired, got %v", err)
		}
	})

	t.Run("an unrelated error passes through unchanged", func(t *testing.T) {
		reset := errors.New("connection reset by peer")
		g := &gapReader{body: &scriptedBody{steps: []func([]byte) (int, error){
			func([]byte) (int, error) { return 0, reset },
		}}, first: time.Second, cancel: noop}
		if _, err := g.Read(buf); err != reset {
			t.Fatalf("expected the reset error unchanged, got %v", err)
		}
	})
}

// TestRequestTimeoutWinsOverASlowDial covers nooga/let-go#857 item 1: with
// both scopes set, a dial still running when the request deadline expires is
// a request timeout with the request duration, not a connect timeout.
func TestRequestTimeoutWinsOverASlowDial(t *testing.T) {
	for _, stream := range []bool{false, true} {
		timeouts := vm.EmptyPersistentMap.
			Assoc(vm.Keyword("connect"), vm.Float(5)).
			Assoc(vm.Keyword("request"), vm.Float(0.1))
		started := time.Now()
		_, err := invokeRequest(t, requestOpts("http://10.255.255.1:9/", stream, timeouts))
		if err == nil || !strings.Contains(err.Error(), "http request timeout after 100ms") {
			t.Fatalf("stream=%v: expected the request timeout, got %v", stream, err)
		}
		if d := time.Since(started); d > 2*time.Second {
			t.Fatalf("stream=%v: the request deadline did not cut the dial short: %s", stream, d)
		}
	}
}

// timeoutNetErr is a net.Error that reports a timeout, like a dial whose
// deadline passed.
type timeoutNetErr struct{}

func (timeoutNetErr) Error() string   { return "i/o timeout" }
func (timeoutNetErr) Timeout() bool   { return true }
func (timeoutNetErr) Temporary() bool { return true }

// TestDescribeTimeoutPrefersTheRequestScopeThatFired covers nooga/let-go#857
// item 1 deterministically. When the request deadline expires during a dial,
// the transport may return either the context error or the dial's own
// timeout; the dial timeout must still be reported as the request scope.
func TestDescribeTimeoutPrefersTheRequestScopeThatFired(t *testing.T) {
	dialTimeout := &url.Error{Op: "Get", URL: "http://example.invalid/", Err: &net.OpError{Op: "dial", Net: "tcp", Err: timeoutNetErr{}}}
	scopes := requestTimeouts{connect: 10 * time.Second, request: 3 * time.Second}
	if got := describeTimeout(dialTimeout, scopes, true); got == nil || got.Error() != "http request timeout after 3s" {
		t.Fatalf("request deadline fired during the dial: expected the request timeout, got %v", got)
	}
	if got := describeTimeout(dialTimeout, scopes, false); got == nil || got.Error() != "http connect timeout after 10s" {
		t.Fatalf("the dialer's own deadline: expected the connect timeout, got %v", got)
	}
}
