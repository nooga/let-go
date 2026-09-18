//go:build !tinygo && !lg_no_http

package rt

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

// okHandler answers every request with 200 "ok".
func okHandler(t *testing.T) vm.Fn {
	t.Helper()
	fn, err := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
		return vm.NewPersistentMap([]vm.Value{
			vm.Keyword("status"), vm.Int(200),
			vm.Keyword("body"), vm.String("ok"),
		}), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return fn.(vm.Fn)
}

// gatedHandler answers 200 "ok" only after gate is closed, so a request can
// be held in flight across a stop.
func gatedHandler(t *testing.T, gate chan struct{}) vm.Fn {
	t.Helper()
	fn, err := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
		<-gate
		return vm.NewPersistentMap([]vm.Value{
			vm.Keyword("status"), vm.Int(200),
			vm.Keyword("body"), vm.String("ok"),
		}), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return fn.(vm.Fn)
}

func getBody(t *testing.T, url string) (int, string, error) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, string(b), nil
}

// waitReturned reports whether waitServer returns within d.
func waitReturned(s *lgServer, d time.Duration) (error, bool) {
	done := make(chan error, 1)
	go func() { done <- waitServer(s) }()
	select {
	case err := <-done:
		return err, true
	case <-time.After(d):
		return nil, false
	}
}

func TestServerStartServesAndStops(t *testing.T) {
	scope := vm.Goroutines.Child()
	s, err := startServer(scope, scope.Context(), okHandler(t), "127.0.0.1:0")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	url := "http://" + s.ln.Addr().String() + "/"
	code, body, err := getBody(t, url)
	if err != nil || code != 200 || body != "ok" {
		t.Fatalf("GET before stop: code=%d body=%q err=%v", code, body, err)
	}
	stopServer(s, 5*time.Second)
	if _, _, err := getBody(t, url); err == nil {
		t.Fatalf("GET after stop succeeded; expected a connection error")
	}
	if err, ok := waitReturned(s, time.Second); !ok || err != nil {
		t.Fatalf("wait after stop: returned=%v err=%v", ok, err)
	}
}

func TestServerRecordExposesAddrAndPort(t *testing.T) {
	scope := vm.Goroutines.Child()
	s, err := startServer(scope, scope.Context(), okHandler(t), "127.0.0.1:0")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer stopServer(s, time.Second)
	rec := serverRecord(s).(vm.Lookup)
	addr := rec.ValueAt(vm.Keyword("addr"))
	port := rec.ValueAt(vm.Keyword("port"))
	if !strings.HasPrefix(addr.(vm.String).String(), "\"127.0.0.1:") {
		t.Fatalf("addr: %v", addr)
	}
	if p := int(port.(vm.Int)); p <= 0 || !strings.HasSuffix(string(addr.(vm.String)), ":"+strconv.Itoa(p)) {
		t.Fatalf("port %d does not match addr %v", p, addr)
	}
	got, err := unboxServer(serverRecord(s))
	if err != nil || got != s {
		t.Fatalf("unboxServer(record): %v %v", got, err)
	}
	got, err = unboxServer(rec.ValueAt(vm.Keyword("server")))
	if err != nil || got != s {
		t.Fatalf("unboxServer(boxed): %v %v", got, err)
	}
	if _, err := unboxServer(vm.String("nope")); err == nil {
		t.Fatalf("unboxServer accepted a string")
	}
}

func TestServerStartRejectsBadAddress(t *testing.T) {
	scope := vm.Goroutines.Child()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	if _, err := startServer(scope, scope.Context(), okHandler(t), ln.Addr().String()); err == nil {
		t.Fatalf("start on a held port succeeded")
	}
	if _, err := startServer(scope, scope.Context(), okHandler(t), "127.0.0.1:99999"); err == nil {
		t.Fatalf("start on an invalid port succeeded")
	}
}

func TestServerStopIsIdempotent(t *testing.T) {
	scope := vm.Goroutines.Child()
	s, err := startServer(scope, scope.Context(), okHandler(t), "127.0.0.1:0")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	stopServer(s, time.Second)
	stopServer(s, time.Second)
	if err, ok := waitReturned(s, time.Second); !ok || err != nil {
		t.Fatalf("wait: returned=%v err=%v", ok, err)
	}
}

func TestServerStopsOnScopeCancel(t *testing.T) {
	scope := vm.Goroutines.Child()
	s, err := startServer(scope, scope.Context(), okHandler(t), "127.0.0.1:0")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	scope.Cancel()
	if err, ok := waitReturned(s, time.Second); !ok || err != nil {
		t.Fatalf("wait after cancel: returned=%v err=%v", ok, err)
	}
	if _, _, err := getBody(t, "http://"+s.ln.Addr().String()+"/"); err == nil {
		t.Fatalf("GET after scope cancel succeeded")
	}
}

func TestServerStopDrainsActiveRequest(t *testing.T) {
	scope := vm.Goroutines.Child()
	gate := make(chan struct{})
	s, err := startServer(scope, scope.Context(), gatedHandler(t, gate), "127.0.0.1:0")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	url := "http://" + s.ln.Addr().String() + "/"
	type result struct {
		code int
		err  error
	}
	got := make(chan result, 1)
	go func() {
		code, _, err := getBody(t, url)
		got <- result{code, err}
	}()
	// Let the request reach the handler before stopping.
	time.Sleep(100 * time.Millisecond)
	go stopServer(s, 5*time.Second)
	if _, ok := waitReturned(s, 300*time.Millisecond); ok {
		t.Fatalf("wait returned while a request was still in flight")
	}
	close(gate)
	r := <-got
	if r.err != nil || r.code != 200 {
		t.Fatalf("in-flight request: code=%d err=%v", r.code, r.err)
	}
	if err, ok := waitReturned(s, time.Second); !ok || err != nil {
		t.Fatalf("wait after drain: returned=%v err=%v", ok, err)
	}
}

func TestServerStopForcesCloseAfterTimeout(t *testing.T) {
	scope := vm.Goroutines.Child()
	gate := make(chan struct{})
	defer close(gate)
	s, err := startServer(scope, scope.Context(), gatedHandler(t, gate), "127.0.0.1:0")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	url := "http://" + s.ln.Addr().String() + "/"
	got := make(chan error, 1)
	go func() {
		_, _, err := getBody(t, url)
		got <- err
	}()
	time.Sleep(100 * time.Millisecond)
	stopServer(s, 50*time.Millisecond)
	if err, ok := waitReturned(s, time.Second); !ok || err != nil {
		t.Fatalf("wait after forced close: returned=%v err=%v", ok, err)
	}
	select {
	case err := <-got:
		if err == nil {
			t.Fatalf("held request completed; expected it to be cut off")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("held request did not fail after Close")
	}
}

func TestServerWaitReportsServeError(t *testing.T) {
	// Closing the listener behind the server's back makes Serve return a
	// non-ErrServerClosed error, which wait must surface.
	scope := vm.Goroutines.Child()
	s, err := startServer(scope, scope.Context(), okHandler(t), "127.0.0.1:0")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	s.ln.Close()
	err, ok := waitReturned(s, time.Second)
	if !ok {
		t.Fatalf("wait did not return after the listener died")
	}
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		t.Fatalf("expected the Serve error, got %v", err)
	}
}
