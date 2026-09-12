//go:build !tinygo && !lg_no_http

/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

func TestHTTPScopeCancellationIsCatchable(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	scope := vm.Goroutines.Child()
	defer func() {
		close(release)
		server.Close()
		vm.CloseScoped(scope, time.Second)
	}()
	lg, err := NewLetGo("http-cancellation-test")
	if err != nil {
		t.Fatal(err)
	}
	v, err := lg.Run(fmt.Sprintf(`(fn [] (try (http/request {:url %q}) (catch e :caught)))`, server.URL))
	if err != nil {
		t.Fatal(err)
	}
	ec := vm.NewExecContext()
	ec.SetScope(scope)
	var result vm.Value
	var invokeErr error
	scope.Go(func(context.Context) {
		result, invokeErr = ec.Invoke(v.(vm.Fn), nil)
	})
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("compiled request did not reach server")
	}
	if !scope.Shutdown(time.Second) {
		t.Fatal("compiled HTTP call did not drain on cancellation")
	}
	if invokeErr != nil || result != vm.Keyword("caught") {
		t.Fatalf("cancellation did not reach Lisp catch: result=%v err=%v", result, invokeErr)
	}
}
