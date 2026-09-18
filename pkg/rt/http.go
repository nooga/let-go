//go:build !tinygo && !lg_no_http

/*
 * Copyright (c) 2022-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

// rawString extracts a raw Go string from a Value without quoting.
func rawString(v vm.Value) string {
	if s, ok := v.(vm.String); ok {
		return string(s)
	}
	if kw, ok := v.(vm.Keyword); ok {
		return string(kw)
	}
	return v.String()
}

// extractURL gets a URL string from a String value or a URL record.
func extractURL(v vm.Value) (string, error) {
	if s, ok := v.(vm.String); ok {
		return string(s), nil
	}
	if r, ok := v.(*vm.Record); ok {
		raw := r.ValueAt(vm.Keyword("raw"))
		if raw != vm.NIL {
			if s, ok := raw.(vm.String); ok {
				return string(s), nil
			}
		}
		return "", fmt.Errorf("URL record missing :raw field")
	}
	return "", fmt.Errorf("expected String or URL, got %s", v.Type().Name())
}

type Handler struct {
	fn vm.Fn
}

// methodToLG maps an HTTP method to a lowercase keyword, Ring-style
// (:get, :post, …). Extension methods (PATCH, etc.) pass through lowercased
// rather than falling back to an empty keyword.
func methodToLG(method string) vm.Keyword {
	return vm.Keyword(strings.ToLower(method))
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, request *http.Request) {
	url := request.URL

	scheme := "http"
	if request.TLS != nil {
		scheme = "https"
	}

	defer request.Body.Close()
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		resp.WriteHeader(500)
		_, err := resp.Write(fmt.Appendf(nil, "%s", err))
		if err != nil {
			_ = WriteToErr(nil, fmt.Sprintln("HTTP Error while writing error 500", err))
		}
		return
	}

	var headers vm.Value = vm.NIL
	var contentType string
	if len(request.Header) > 0 {
		hs := vm.EmptyPersistentMap
		for k, v := range request.Header {
			hs = hs.Assoc(vm.String(strings.ToLower(k)), vm.String(strings.Join(v, ","))).(*vm.PersistentMap)
		}
		headers = hs
		contentType = request.Header.Get("Content-Type")
	}

	req := httpRequestMapping.StructToRecord(HTTPRequest{
		RequestMethod: methodToLG(request.Method),
		Scheme:        scheme,
		URI:           url.RequestURI(),
		Path:          url.Path,
		QueryString:   url.RawQuery,
		Body:          string(bodyBytes),
		RemoteAddr:    request.RemoteAddr,
		ServerAddr:    request.Host,
		ServerPort:    url.Port(),
		ContentType:   contentType,
		Headers:       headers,
	})

	res, err := h.fn.Invoke([]vm.Value{req})
	if err != nil {
		resp.WriteHeader(500)
		_, err := resp.Write(fmt.Appendf(nil, "%s", err))
		if err != nil {
			_ = WriteToErr(nil, fmt.Sprintln("HTTP Error while writing error 500", err))
		}
		return
	}

	ress, ok := res.(vm.Lookup)
	if !ok {
		resp.WriteHeader(500)
		_, err := resp.Write([]byte("handler returned malformed response"))
		if err != nil {
			_ = WriteToErr(nil, fmt.Sprintln("HTTP Error while writing error 500", err))
		}
		return
	}
	head := resp.Header()
	respHeaders := ress.ValueAt(vm.Keyword("headers"))
	if respHeaders != vm.NIL {
		if sq, ok := respHeaders.(vm.Sequable); ok {
			for s := sq.Seq(); s != nil; s = s.Next() {
				entry := s.First()
				// An empty map seqs to one nil entry; skip it, as the
				// client header loops do.
				if entry == vm.NIL {
					continue
				}
				// Use Sequable to get key/value from any vector type
				eSeq, ok := entry.(vm.Sequable)
				if !ok {
					continue
				}
				es := eSeq.Seq()
				k := es.First()
				v := es.Next().First()
				head.Add(rawString(k), rawString(v))
			}
		}
	}
	status := ress.ValueAt(vm.Keyword("status"))
	if status == vm.NIL {
		status = vm.Int(200)
	}
	body := ress.ValueAt(vm.Keyword("body"))
	resp.WriteHeader(int(status.(vm.Int)))
	if streamResponseBody(resp, request, body) {
		return
	}
	respBody, bodyErr := coerceResponseBody(body)
	if bodyErr != nil {
		_ = WriteToErr(nil, fmt.Sprintln("HTTP Error coercing body:", bodyErr))
		return
	}
	if respBody != nil {
		_, err = resp.Write(respBody)
	}
	if err != nil {
		_ = WriteToErr(nil, fmt.Sprintln("HTTP Error while writing error 500", err))
	}
}

// streamResponseBody writes a channel or lazy sequence body incrementally,
// flushing after every element, so a handler can hold or pace a response
// mid-body (SSE fixtures, long polls). Strings and readers keep the buffered
// path. Iteration stops when the client goes away. Reports whether it handled
// the body.
func streamResponseBody(resp http.ResponseWriter, request *http.Request, body vm.Value) bool {
	flusher, _ := resp.(http.Flusher)
	writeChunk := func(v vm.Value) bool {
		var chunk string
		if str, ok := v.(vm.String); ok {
			chunk = string(str)
		} else {
			chunk = v.String()
		}
		if _, err := resp.Write([]byte(chunk)); err != nil {
			return false
		}
		if flusher != nil {
			flusher.Flush()
		}
		return request.Context().Err() == nil
	}
	switch b := body.(type) {
	case vm.Chan:
		if flusher != nil {
			flusher.Flush()
		}
		for {
			select {
			case v, ok := <-b:
				if !ok || !writeChunk(v) {
					return true
				}
			case <-request.Context().Done():
				return true
			}
		}
	case *vm.LazySeq:
		if flusher != nil {
			flusher.Flush()
		}
		for s := b.Seq(); s != nil; s = s.Next() {
			if !writeChunk(s.First()) {
				return true
			}
		}
		return true
	}
	return false
}

// defaultStopTimeout bounds http/stop's graceful drain before it falls back
// to closing connections outright.
const defaultStopTimeout = 5 * time.Second

// lgServer is the handle behind http/start: a bound listener, the server
// draining it, and a done channel that closes once the server has fully
// stopped. "Fully" matters: Shutdown makes Serve return ErrServerClosed at
// once, before in-flight requests finish, so done is owned by the stop path
// (closed after Shutdown or the Close fallback returns) and by the serving
// goroutine only when Serve fails on its own.
type lgServer struct {
	srv      *http.Server
	ln       net.Listener
	done     chan struct{}
	err      error
	stopOnce sync.Once
	doneOnce sync.Once
}

func (s *lgServer) finish(err error) {
	s.doneOnce.Do(func() {
		s.err = err
		close(s.done)
	})
}

// startServer binds addr now, so a bad address is the caller's error, then
// serves on a goroutine of scope. A second goroutine ties the server to the
// scope's context: cancelling the scope stops the server the same way
// http/stop does.
func startServer(scope *vm.Scope, ctx context.Context, handler vm.Fn, addr string) (*lgServer, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	s := &lgServer{
		srv:  &http.Server{Handler: &Handler{fn: handler}},
		ln:   ln,
		done: make(chan struct{}),
	}
	scope.Go(func(_ context.Context) {
		err := s.srv.Serve(ln)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.finish(err)
		}
	})
	scope.Go(func(_ context.Context) {
		select {
		case <-ctx.Done():
			stopServer(s, defaultStopTimeout)
		case <-s.done:
		}
	})
	return s, nil
}

// stopServer drains the server gracefully for up to timeout, then closes
// whatever is still open. Idempotent: a second caller returns at once and
// observes completion through waitServer.
func stopServer(s *lgServer, timeout time.Duration) {
	s.stopOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := s.srv.Shutdown(ctx); err != nil {
			_ = s.srv.Close()
		}
		s.finish(nil)
	})
}

// waitServer blocks until the server has stopped: nil after http/stop (or
// a scope cancel), else the error Serve returned.
func waitServer(s *lgServer) error {
	<-s.done
	return s.err
}

func serverRecord(s *lgServer) vm.Value {
	port := 0
	if tcp, ok := s.ln.Addr().(*net.TCPAddr); ok {
		port = tcp.Port
	}
	return httpServerMapping.StructToRecord(HTTPServer{
		Addr:   s.ln.Addr().String(),
		Port:   port,
		Server: vm.NewBoxed(s),
	})
}

// unboxServer accepts the record http/start returned or its :server value.
func unboxServer(v vm.Value) (*lgServer, error) {
	asServer := func(v vm.Value) *lgServer {
		if b, ok := v.(*vm.Boxed); ok {
			if s, ok := b.Unbox().(*lgServer); ok {
				return s
			}
		}
		return nil
	}
	if s := asServer(v); s != nil {
		return s, nil
	}
	// A record: the :server field. Checked after the boxed case because a
	// *vm.Boxed is also a Lookup (by reflection) and rejects the probe.
	if _, boxed := v.(*vm.Boxed); !boxed {
		if l, ok := v.(vm.Lookup); ok {
			if s := asServer(l.ValueAt(vm.Keyword("server"))); s != nil {
				return s, nil
			}
		}
	}
	return nil, vm.NewExecutionError(fmt.Sprintf("expected an http server (from http/start), got %s", v.Type().Name()))
}

func init() { RegisterInstaller(installHttpNS) }

// nolint
func installHttpNS() {
	// serveArgs validates the (handler addr) pair http/serve and http/start
	// share.
	serveArgs := func(name string, vs []vm.Value) (vm.Fn, string, error) {
		if len(vs) != 2 {
			return nil, "", vm.NewExecutionError(name + " expects 2 args (handler, addr)")
		}
		handlerFunc, ok := vs[0].(vm.Fn)
		if !ok {
			return nil, "", vm.NewExecutionError(name + " expected handler function as Fn")
		}
		addr, ok := vs[1].(vm.String)
		if !ok {
			return nil, "", vm.NewExecutionError(name + " expected listen address as String")
		}
		return handlerFunc, string(addr), nil
	}

	// http/serve — (http/serve handler addr)
	// Ring-style: handler is a fn that takes a request map, returns a response map.
	// Blocks until the server stops: on a bind error, a scope cancel, or an
	// http/stop from another goroutine.
	serve := vm.NewCtxNativeFn("http/serve", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		handlerFunc, addr, err := serveArgs("serve", vs)
		if err != nil {
			return vm.NIL, err
		}
		s, err := startServer(ec.Scope(), ec.Context(), handlerFunc, addr)
		if err != nil {
			return vm.NIL, err
		}
		if err := waitServer(s); err != nil {
			return vm.NIL, err
		}
		return vm.NIL, nil
	})

	// http/start — (http/start handler addr) -> http/Server record
	// Binds now and serves in the background; the record carries :addr,
	// :port (resolved, so ":0" works) and the :server handle for stop/wait.
	start := vm.NewCtxNativeFn("http/start", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		handlerFunc, addr, err := serveArgs("start", vs)
		if err != nil {
			return vm.NIL, err
		}
		s, err := startServer(ec.Scope(), ec.Context(), handlerFunc, addr)
		if err != nil {
			return vm.NIL, err
		}
		return serverRecord(s), nil
	})

	// http/stop — (http/stop server) or (http/stop server timeout-ms)
	// Graceful shutdown, then Close once timeout-ms (default 5000) elapses.
	stop, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, vm.NewExecutionError("stop expects 1-2 args (server, timeout-ms)")
		}
		s, err := unboxServer(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		timeout := defaultStopTimeout
		if len(vs) == 2 {
			ms, ok := vs[1].(vm.Int)
			if !ok || ms < 0 {
				return vm.NIL, vm.NewExecutionError("stop timeout-ms must be a non-negative integer")
			}
			timeout = time.Duration(ms) * time.Millisecond
		}
		stopServer(s, timeout)
		return vm.NIL, nil
	})
	if err != nil {
		panic("http NS init failed")
	}

	// http/wait — (http/wait server)
	// Blocks until the server has stopped; nil on a clean stop.
	wait, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, vm.NewExecutionError("wait expects 1 arg (server)")
		}
		s, err := unboxServer(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		if err := waitServer(s); err != nil {
			return vm.NIL, err
		}
		return vm.NIL, nil
	})
	if err != nil {
		panic("http NS init failed")
	}

	// HTTP client: build response record from http.Response
	// Default: body is a string. With :as :stream in opts, body is an io/reader.
	buildResponseMap := func(resp *http.Response, asStream bool) (vm.Value, error) {
		hs := vm.EmptyPersistentMap
		for k, v := range resp.Header {
			hs = hs.Assoc(vm.String(strings.ToLower(k)), vm.String(strings.Join(v, ","))).(*vm.PersistentMap)
		}
		var body vm.Value
		if asStream {
			body = vm.NewBoxed(newLGReader(resp.Body, resp.Body))
		} else {
			defer resp.Body.Close()
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return vm.NIL, err
			}
			body = vm.String(bodyBytes)
		}
		return httpResponseMapping.StructToRecord(HTTPResponse{
			Status:  resp.StatusCode,
			Body:    body,
			Headers: hs,
		}), nil
	}

	// Check if opts map has :as :stream
	isStreamOpt := func(opts vm.Value) bool {
		if opts == nil || opts == vm.NIL {
			return false
		}
		if l, ok := opts.(vm.Lookup); ok {
			v := l.ValueAt(vm.Keyword("as"))
			return v == vm.Keyword("stream")
		}
		return false
	}

	// http/get — (http/get url) or (http/get url opts)
	httpGet := vm.NewCtxNativeFn("http/get", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("http/get expects 1-2 args")
		}
		urlStr, err := extractURL(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		req, err := http.NewRequestWithContext(ec.Context(), "GET", urlStr, nil)
		if err != nil {
			return vm.NIL, err
		}
		if len(vs) == 2 {
			if opts, ok := vs[1].(vm.Lookup); ok {
				hdrs := opts.ValueAt(vm.Keyword("headers"))
				if hdrs != vm.NIL {
					if sq, ok := hdrs.(vm.Sequable); ok {
						for s := sq.Seq(); s != nil; s = s.Next() {
							entry := s.First()
							if entry == vm.NIL {
								continue
							}
							eSeq, ok := entry.(vm.Sequable)
							if !ok {
								continue
							}
							es := eSeq.Seq()
							k := es.First()
							v := es.Next().First()
							req.Header.Set(rawString(k), rawString(v))
						}
					}
				}
			}
		}
		var opts vm.Value
		if len(vs) == 2 {
			opts = vs[1]
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return vm.NIL, err
		}
		return buildResponseMap(resp, isStreamOpt(opts))
	})

	// http/post — (http/post url body) or (http/post url body opts)
	httpPost := vm.NewCtxNativeFn("http/post", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) < 2 || len(vs) > 3 {
			return vm.NIL, fmt.Errorf("http/post expects 2-3 args")
		}
		urlStr, err := extractURL(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		var bodyStr string
		if s, ok := vs[1].(vm.String); ok {
			bodyStr = string(s)
		} else {
			bodyStr = vs[1].String()
		}
		req, err := http.NewRequestWithContext(ec.Context(), "POST", urlStr, strings.NewReader(bodyStr))
		if err != nil {
			return vm.NIL, err
		}
		if len(vs) == 3 {
			if opts, ok := vs[2].(vm.Lookup); ok {
				ct := opts.ValueAt(vm.Keyword("content-type"))
				if ct != vm.NIL {
					req.Header.Set("Content-Type", rawString(ct))
				}
				hdrs := opts.ValueAt(vm.Keyword("headers"))
				if hdrs != vm.NIL {
					if sq, ok := hdrs.(vm.Sequable); ok {
						for s := sq.Seq(); s != nil; s = s.Next() {
							entry := s.First()
							if entry == vm.NIL {
								continue
							}
							eSeq, ok := entry.(vm.Sequable)
							if !ok {
								continue
							}
							es := eSeq.Seq()
							k := es.First()
							v := es.Next().First()
							req.Header.Set(rawString(k), rawString(v))
						}
					}
				}
			}
		}
		var postOpts vm.Value
		if len(vs) == 3 {
			postOpts = vs[2]
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return vm.NIL, err
		}
		return buildResponseMap(resp, isStreamOpt(postOpts))
	})

	// http/request — (http/request {:method :get :url "..." :headers {...} :body "..."})
	httpRequest := vm.NewCtxNativeFn("http/request", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("http/request expects 1 arg (options map)")
		}
		opts, ok := vs[0].(vm.Lookup)
		if !ok {
			return vm.NIL, fmt.Errorf("http/request expected a map")
		}
		method := "GET"
		if m := opts.ValueAt(vm.Keyword("method")); m != vm.NIL {
			ms := m.String()
			if ms[0] == ':' {
				ms = ms[1:]
			}
			method = strings.ToUpper(ms)
		}
		urlVal := opts.ValueAt(vm.Keyword("url"))
		if urlVal == vm.NIL {
			return vm.NIL, fmt.Errorf("http/request requires :url")
		}
		reqURL, err := extractURL(urlVal)
		if err != nil {
			return vm.NIL, err
		}
		var bodyReader io.Reader
		if b := opts.ValueAt(vm.Keyword("body")); b != vm.NIL {
			if s, ok := b.(vm.String); ok {
				bodyReader = strings.NewReader(string(s))
			} else {
				bodyReader = strings.NewReader(b.String())
			}
		}
		req, err := http.NewRequestWithContext(ec.Context(), method, reqURL, bodyReader)
		if err != nil {
			return vm.NIL, err
		}
		hdrs := opts.ValueAt(vm.Keyword("headers"))
		if hdrs != vm.NIL {
			if sq, ok := hdrs.(vm.Sequable); ok {
				for s := sq.Seq(); s != nil; s = s.Next() {
					entry := s.First()
					if entry == vm.NIL {
						continue
					}
					eSeq, ok := entry.(vm.Sequable)
					if !ok {
						continue
					}
					es := eSeq.Seq()
					k := es.First()
					v := es.Next().First()
					req.Header.Set(rawString(k), rawString(v))
				}
			}
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return vm.NIL, err
		}
		return buildResponseMap(resp, isStreamOpt(vs[0]))
	})

	ns := vm.NewNamespace("http")

	// Intentional shadows of clojure.core names — suppress warn-on-shadow.
	ns.Exclude("get")

	// The server fns run under the caller's scope too: a cancelled scope
	// stops the server, so they carry the same meta as the clients.
	clientMeta := vm.EmptyPersistentMap.Assoc(vm.Keyword("scope-cancellation"), vm.TRUE)
	ns.Def("serve", serve).SetMeta(clientMeta)
	ns.Def("start", start).SetMeta(clientMeta)
	ns.Def("stop", stop)
	ns.Def("wait", wait)
	ns.Def("get", httpGet).SetMeta(clientMeta)
	ns.Def("post", httpPost).SetMeta(clientMeta)
	ns.Def("request", httpRequest).SetMeta(clientMeta)
	RegisterNS(ns)
}
