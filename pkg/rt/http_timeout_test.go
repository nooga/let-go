//go:build !tinygo && !lg_no_http

package rt

import (
	"io"
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

func TestLineSeqPropagatesReadErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher := w.(http.Flusher)
		_, _ = w.Write([]byte("first\n"))
		flusher.Flush()
		time.Sleep(500 * time.Millisecond)
		_, _ = w.Write([]byte("second\n"))
	}))
	defer server.Close()
	timeouts := vm.EmptyPersistentMap.Assoc(vm.Keyword("stream_read"), vm.Float(0.1))
	v, err := invokeRequest(t, requestOpts(server.URL, true, timeouts))
	if err != nil {
		t.Fatal(err)
	}
	body := v.(vm.Lookup).ValueAt(vm.Keyword("body"))
	lineSeq := LookupVar("io", "line-seq").Deref().(vm.Fn)
	ec := vm.NewExecContext()
	seqVal, err := ec.Invoke(lineSeq, []vm.Value{body})
	if err != nil {
		t.Fatal(err)
	}
	first := LookupVar("core", "first").Deref().(vm.Fn)
	rest := LookupVar("core", "rest").Deref().(vm.Fn)
	head, err := ec.Invoke(first, []vm.Value{seqVal})
	if err != nil || head.String() != `"first"` {
		t.Fatalf("first line: %v %v", head, err)
	}
	tail, err := ec.Invoke(rest, []vm.Value{seqVal})
	if err != nil {
		t.Fatal(err)
	}
	_, err = ec.Invoke(first, []vm.Value{tail})
	if err == nil || !strings.Contains(err.Error(), "stream_read timeout") {
		t.Fatalf("line-seq must surface the read error, got %v", err)
	}
}
