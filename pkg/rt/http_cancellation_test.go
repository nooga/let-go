//go:build !tinygo && !lg_no_http

/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

func scopedHTTPArgs(name, url string, stream bool) []vm.Value {
	opts := vm.EmptyPersistentMap.Assoc(vm.Keyword("url"), vm.String(url))
	if stream {
		opts = opts.(vm.Associative).Assoc(vm.Keyword("as"), vm.Keyword("stream"))
	}
	switch name {
	case "get":
		return []vm.Value{vm.String(url), opts}
	case "post":
		return []vm.Value{vm.String(url), vm.String("payload"), opts}
	default:
		return []vm.Value{opts}
	}
}

func TestHTTPClientScopeCancellation(t *testing.T) {
	for _, name := range []string{"get", "post", "request"} {
		for _, phase := range []string{"headers", "body", "stream"} {
			t.Run(name+"/"+phase, func(t *testing.T) {
				started := make(chan struct{})
				canceled := make(chan struct{})
				release := make(chan struct{})
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Drain uploads so net/http can detect a peer disconnect.
					_, _ = io.Copy(io.Discard, r.Body)
					if phase != "headers" {
						w.WriteHeader(http.StatusOK)
						w.(http.Flusher).Flush()
					}
					close(started)
					select {
					case <-r.Context().Done():
						close(canceled)
					case <-release:
					}
				}))
				scope := vm.Goroutines.Child()
				t.Cleanup(func() {
					scope.Cancel()
					close(release)
					server.Close()
					vm.CloseScoped(scope, time.Second)
				})
				ec := vm.NewExecContext()
				ec.SetScope(scope)
				fn := LookupVar("http", name).Deref().(vm.Fn)
				type outcome struct {
					value vm.Value
					err   error
				}
				done := make(chan outcome, 1)
				reading := make(chan struct{})
				scope.Go(func(context.Context) {
					v, err := ec.Invoke(fn, scopedHTTPArgs(name, server.URL, phase == "stream"))
					if err == nil && phase == "stream" {
						body := v.(vm.Lookup).ValueAt(vm.Keyword("body")).Unbox().(*LGReader)
						defer body.Close()
						close(reading)
						_, err = io.ReadAll(body)
					}
					done <- outcome{v, err}
				})
				select {
				case <-started:
				case <-time.After(2 * time.Second):
					t.Fatal("request did not reach server")
				}
				if phase == "stream" {
					select {
					case <-reading:
					case <-time.After(2 * time.Second):
						t.Fatal("stream response not returned")
					}
				}
				if !scope.Shutdown(time.Second) {
					t.Fatalf("cancellation left %d HTTP worker(s) live", scope.LiveTree())
				}
				result := <-done
				if !errors.Is(result.err, context.Canceled) {
					t.Fatalf("want context.Canceled, got value=%v err=%v", result.value, result.err)
				}
				if phase != "stream" && result.value != vm.NIL {
					t.Fatalf("failed native call returned %v instead of nil", result.value)
				}
				select {
				case <-canceled:
				case <-time.After(time.Second):
					t.Fatal("server did not observe request cancellation")
				}
				if scope.LiveTree() != 0 {
					t.Fatal("scope not drained")
				}
			})
		}
	}
}

func TestHTTPClientScopeIsolationAndStreamLifetime(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		_, _ = io.WriteString(w, "complete response")
	}))
	defer server.Close()
	for _, name := range []string{"get", "post", "request"} {
		t.Run(name, func(t *testing.T) {
			canceled := vm.Goroutines.Child()
			defer vm.CloseScoped(canceled, time.Second)
			canceled.Cancel()
			sibling := vm.Goroutines.Child()
			defer vm.CloseScoped(sibling, time.Second)
			ec := vm.NewExecContext()
			ec.SetScope(sibling)
			fn := LookupVar("http", name).Deref().(vm.Fn)
			for _, stream := range []bool{false, true} {
				v, err := ec.Invoke(fn, scopedHTTPArgs(name, server.URL, stream))
				if err != nil {
					t.Fatal(err)
				}
				body := v.(vm.Lookup).ValueAt(vm.Keyword("body"))
				var content string
				if stream {
					reader := body.Unbox().(*LGReader)
					data, err := io.ReadAll(reader)
					_ = reader.Close()
					if err != nil {
						t.Fatal(err)
					}
					content = string(data)
				} else {
					content = string(body.(vm.String))
				}
				if content != "complete response" {
					t.Fatalf("body = %q", content)
				}
			}
		})
	}
}

func TestHTTPClientScopeCancelsUnreadStream(t *testing.T) {
	for _, name := range []string{"get", "post", "request"} {
		t.Run(name, func(t *testing.T) {
			canceled := make(chan struct{})
			release := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = io.Copy(io.Discard, r.Body)
				w.WriteHeader(http.StatusOK)
				w.(http.Flusher).Flush()
				select {
				case <-r.Context().Done():
					close(canceled)
				case <-release:
				}
			}))
			scope := vm.Goroutines.Child()
			var response vm.Value
			var requestErr error
			t.Cleanup(func() {
				scope.Cancel()
				close(release)
				server.Close()
				vm.CloseScoped(scope, time.Second)
				if requestErr == nil && response != nil {
					body := response.(vm.Lookup).ValueAt(vm.Keyword("body")).Unbox().(*LGReader)
					_ = body.Close()
				}
			})
			ec := vm.NewExecContext()
			ec.SetScope(scope)
			fn := LookupVar("http", name).Deref().(vm.Fn)
			scope.Go(func(context.Context) {
				response, requestErr = ec.Invoke(fn, scopedHTTPArgs(name, server.URL, true))
			})
			if !scope.Await(2 * time.Second) {
				t.Fatal("stream response did not return")
			}
			if requestErr != nil {
				t.Fatal(requestErr)
			}
			// No body read or close may cause the disconnect being asserted.
			select {
			case <-canceled:
				t.Fatal("stream canceled before scope cancellation")
			default:
			}
			if !scope.Shutdown(time.Second) || scope.LiveTree() != 0 {
				t.Fatal("scope did not drain")
			}
			select {
			case <-canceled:
			case <-time.After(time.Second):
				t.Fatal("scope cancellation did not disconnect unread, unclosed stream")
			}
		})
	}
}
