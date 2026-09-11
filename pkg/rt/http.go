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
	"sync/atomic"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

// requestTimeouts are the three client timeout scopes an http/request may
// carry under :timeout — {:connect s :request s :stream_read s} (seconds,
// numbers) or a bare number meaning :request — plus the older :timeout_ms
// (request scope, milliseconds). Zero means "not set".
type requestTimeouts struct {
	connect    time.Duration
	request    time.Duration
	streamRead time.Duration
}

func secondsValue(v vm.Value) time.Duration {
	switch n := v.(type) {
	case vm.Int:
		return time.Duration(float64(n) * float64(time.Second))
	case vm.Float:
		return time.Duration(float64(n) * float64(time.Second))
	default:
		return 0
	}
}

func timeoutsFromOpts(opts vm.Lookup) requestTimeouts {
	var t requestTimeouts
	if ms := opts.ValueAt(vm.Keyword("timeout_ms")); ms != vm.NIL {
		if n, ok := ms.(vm.Int); ok && n > 0 {
			t.request = time.Duration(int64(n)) * time.Millisecond
		}
	}
	tv := opts.ValueAt(vm.Keyword("timeout"))
	if tv == vm.NIL {
		return t
	}
	if d := secondsValue(tv); d > 0 {
		t.request = d
		return t
	}
	if l, ok := tv.(vm.Lookup); ok {
		if d := secondsValue(l.ValueAt(vm.Keyword("connect"))); d > 0 {
			t.connect = d
		}
		if d := secondsValue(l.ValueAt(vm.Keyword("request"))); d > 0 {
			t.request = d
		}
		if d := secondsValue(l.ValueAt(vm.Keyword("stream_read"))); d > 0 {
			t.streamRead = d
		}
	}
	return t
}

// transportsByConnectTimeout keeps one pooled transport per connect timeout
// so requests still reuse connections.
var transportsByConnectTimeout sync.Map

func clientForConnectTimeout(d time.Duration) *http.Client {
	if d <= 0 {
		return http.DefaultClient
	}
	if c, ok := transportsByConnectTimeout.Load(d); ok {
		return c.(*http.Client)
	}
	base, _ := http.DefaultTransport.(*http.Transport)
	tr := base.Clone()
	dialer := &net.Dialer{Timeout: d, KeepAlive: 30 * time.Second}
	tr.DialContext = dialer.DialContext
	c := &http.Client{Transport: tr}
	actual, _ := transportsByConnectTimeout.LoadOrStore(d, c)
	return actual.(*http.Client)
}

// timeoutError names the scope that fired so callers can classify it.
type timeoutError struct {
	scope string
	after time.Duration
}

func (e *timeoutError) Error() string {
	return fmt.Sprintf("http %s timeout after %s", e.scope, e.after)
}

func describeTimeout(err error, t requestTimeouts, connectOnly bool) error {
	if err == nil {
		return nil
	}
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() && t.connect > 0 {
		var opErr *net.OpError
		if errors.As(err, &opErr) && opErr.Op == "dial" {
			return &timeoutError{"connect", t.connect}
		}
	}
	if errors.Is(err, context.DeadlineExceeded) && t.request > 0 && !connectOnly {
		return &timeoutError{"request", t.request}
	}
	return err
}

// cancelOnClose releases the request context when a streamed body closes.
type cancelOnClose struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c *cancelOnClose) Close() error {
	err := c.ReadCloser.Close()
	c.cancel()
	return err
}

// gapReader enforces a stream_read timeout: each Read must complete within
// the gap, otherwise the request context is cancelled and the read fails.
type gapReader struct {
	body   io.ReadCloser
	gap    time.Duration
	cancel context.CancelFunc
	fired  atomic.Bool
}

func (g *gapReader) Read(p []byte) (int, error) {
	timer := time.AfterFunc(g.gap, func() { g.fired.Store(true); g.cancel() })
	n, err := g.body.Read(p)
	timer.Stop()
	if err != nil && g.fired.Load() {
		return n, &timeoutError{"stream_read", g.gap}
	}
	return n, err
}

func (g *gapReader) Close() error {
	err := g.body.Close()
	g.cancel()
	return err
}

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

func init() { RegisterInstaller(installHttpNS) }

// nolint
func installHttpNS() {
	// http/serve — (http/serve handler addr)
	// Ring-style: handler is a fn that takes a request map, returns a response map.
	serve, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, vm.NewExecutionError("serve expects 2 args (handler, addr)")
		}
		handlerFunc, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, vm.NewExecutionError("serve expected handler function as Fn")
		}
		addr, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, vm.NewExecutionError("serve expected listen address as String")
		}
		handler := &Handler{fn: handlerFunc}
		err := http.ListenAndServe(string(addr), handler)
		if err != nil {
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
		timeouts := timeoutsFromOpts(opts)
		asStream := isStreamOpt(vs[0])
		ctx := ec.Context()
		var cancel context.CancelFunc
		if timeouts.request > 0 && !asStream {
			ctx, cancel = context.WithTimeout(ctx, timeouts.request)
		} else {
			ctx, cancel = context.WithCancel(ctx)
		}
		req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
		if err != nil {
			cancel()
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
		// A streamed request applies :request to reaching the headers only;
		// the body is governed by :stream_read gaps and scope cancellation.
		var headerTimer *time.Timer
		var headersTimedOut atomic.Bool
		if asStream && timeouts.request > 0 {
			headerTimer = time.AfterFunc(timeouts.request, func() { headersTimedOut.Store(true); cancel() })
		}
		resp, err := clientForConnectTimeout(timeouts.connect).Do(req)
		if headerTimer != nil {
			headerTimer.Stop()
		}
		if err != nil {
			cancel()
			if headersTimedOut.Load() {
				return vm.NIL, &timeoutError{"request", timeouts.request}
			}
			return vm.NIL, describeTimeout(err, timeouts, false)
		}
		if asStream {
			body := io.ReadCloser(resp.Body)
			if timeouts.streamRead > 0 {
				body = &gapReader{body: resp.Body, gap: timeouts.streamRead, cancel: cancel}
			} else {
				body = &cancelOnClose{ReadCloser: resp.Body, cancel: cancel}
			}
			resp.Body = body
			return buildResponseMap(resp, true)
		}
		defer cancel()
		result, err := buildResponseMap(resp, false)
		if err != nil {
			return vm.NIL, describeTimeout(err, timeouts, false)
		}
		return result, nil
	})

	ns := vm.NewNamespace("http")

	// Intentional shadows of clojure.core names — suppress warn-on-shadow.
	ns.Exclude("get")

	ns.Def("serve", serve)
	clientMeta := vm.EmptyPersistentMap.Assoc(vm.Keyword("scope-cancellation"), vm.TRUE)
	ns.Def("get", httpGet).SetMeta(clientMeta)
	ns.Def("post", httpPost).SetMeta(clientMeta)
	ns.Def("request", httpRequest).SetMeta(clientMeta)
	RegisterNS(ns)
}
