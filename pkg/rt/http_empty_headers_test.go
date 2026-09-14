//go:build !tinygo && !lg_no_http

package rt

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func mapOf(kvs ...vm.Value) vm.Value {
	m := vm.Value(vm.EmptyPersistentMap)
	for i := 0; i+1 < len(kvs); i += 2 {
		m = m.(*vm.PersistentMap).Assoc(kvs[i], kvs[i+1])
	}
	return m
}

// TestHTTPClientsAcceptEmptyHeaders: an empty :headers map is a valid option
// and must not panic any client (nooga/let-go#828).
func TestHTTPClientsAcceptEmptyHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()
	empty := vm.Value(vm.EmptyPersistentMap)
	opts := mapOf(vm.Keyword("headers"), empty)
	cases := []struct {
		name string
		args []vm.Value
	}{
		{"get", []vm.Value{vm.String(server.URL), opts}},
		{"post", []vm.Value{vm.String(server.URL), vm.String("body"), opts}},
		{"request", []vm.Value{mapOf(vm.Keyword("url"), vm.String(server.URL), vm.Keyword("headers"), empty)}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v vm.Value
			var err error
			func() {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("panicked on empty headers: %v", r)
					}
				}()
				fn := LookupVar("http", c.name).Deref().(vm.Fn)
				v, err = vm.NewExecContext().Invoke(fn, c.args)
			}()
			if err != nil {
				t.Fatalf("http/%s: %v", c.name, err)
			}
			if status := v.(vm.Lookup).ValueAt(vm.Keyword("status")); status.String() != "200" {
				t.Fatalf("http/%s: status %v", c.name, status)
			}
		})
	}
}
