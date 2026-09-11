//go:build !tinygo && !lg_no_http

package rt

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

// serveBody starts a server whose handler answers every request with status
// 200 and the given :body value.
func serveBody(t *testing.T, body vm.Value) *httptest.Server {
	t.Helper()
	fn, err := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
		m := vm.EmptyPersistentMap.Assoc(vm.Keyword("status"), vm.Int(200)).(*vm.PersistentMap)
		return m.Assoc(vm.Keyword("body"), body), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return httptest.NewServer(&Handler{fn: fn.(vm.Fn)})
}

// checkStreamed requires the response headers and the first line to arrive
// while the handler still holds the second line behind gate, then opens the
// gate and reads the rest. A buffering server sends nothing, not even headers,
// until the whole body is produced, so the request runs on its own goroutine
// under the deadline; on timeout the gate is opened so the handler can finish
// and the server can shut down.
func checkStreamed(t *testing.T, server *httptest.Server, gate chan struct{}) {
	t.Helper()
	var once sync.Once
	open := func() { once.Do(func() { close(gate) }) }
	defer server.Close()
	defer open()

	type first struct {
		line string
		err  error
		r    *bufio.Reader
		body io.Closer
	}
	got := make(chan first, 1)
	go func() {
		resp, err := http.Get(server.URL)
		if err != nil {
			got <- first{err: err}
			return
		}
		r := bufio.NewReader(resp.Body)
		line, err := r.ReadString('\n')
		got <- first{line: line, err: err, r: r, body: resp.Body}
	}()

	var res first
	select {
	case res = <-got:
	case <-time.After(2 * time.Second):
		open()
		t.Fatal("the first chunk did not arrive while the second was held; the body was buffered")
	}
	if res.body != nil {
		defer res.body.Close()
	}
	if res.err != nil || res.line != "first\n" {
		t.Fatalf("first chunk: %q %v", res.line, res.err)
	}

	open()
	rest, err := io.ReadAll(res.r)
	if err != nil || string(rest) != "second\n" {
		t.Fatalf("rest of body: %q %v", rest, err)
	}
}

// TestServeStreamsChannelAndLazySeqBodies: a channel or lazy-seq :body is
// written element by element with a flush after each, so a client receives the
// first chunk while the handler is still producing the next
// (nooga/let-go#831).
func TestServeStreamsChannelAndLazySeqBodies(t *testing.T) {
	t.Run("channel", func(t *testing.T) {
		gate := make(chan struct{})
		ch := make(vm.Chan)
		go func() {
			ch <- vm.String("first\n")
			<-gate
			ch <- vm.String("second\n")
			close(ch)
		}()
		checkStreamed(t, serveBody(t, ch), gate)
	})

	t.Run("lazy seq", func(t *testing.T) {
		gate := make(chan struct{})
		second, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
			<-gate
			return vm.NewCons(vm.String("second\n"), vm.EmptyList), nil
		})
		first, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
			return vm.NewCons(vm.String("first\n"), vm.NewLazySeq(second.(vm.Fn))), nil
		})
		checkStreamed(t, serveBody(t, vm.NewLazySeq(first.(vm.Fn))), gate)
	})
}
