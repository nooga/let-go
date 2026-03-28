/*
 * Copyright (c) 2022-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"io"
	"net/http"
	"strings"

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

func methodToLG(scheme string) vm.Keyword {
	return map[string]vm.Keyword{
		"GET":     vm.Keyword("get"),
		"POST":    vm.Keyword("post"),
		"PUT":     vm.Keyword("put"),
		"DELETE":  vm.Keyword("delete"),
		"HEAD":    vm.Keyword("head"),
		"OPTIONS": vm.Keyword("options"),
	}[scheme]
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
		_, err := resp.Write([]byte(fmt.Sprintf("%s", err)))
		if err != nil {
			fmt.Println("HTTP Error while writing error 500", err)
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
		RequestMethod: string(methodToLG(request.Method)),
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
		_, err := resp.Write([]byte(fmt.Sprintf("%s", err)))
		if err != nil {
			fmt.Println("HTTP Error while writing error 500", err)
		}
		return
	}

	ress, ok := res.(vm.Lookup)
	if !ok {
		resp.WriteHeader(500)
		_, err := resp.Write([]byte("handler returned malformed response"))
		if err != nil {
			fmt.Println("HTTP Error while writing error 500", err)
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
				ks := k.String()
				if k.Type() == vm.KeywordType {
					ks = ks[1:]
				}
				head.Add(ks, v.String())
			}
		}
	}
	status := ress.ValueAt(vm.Keyword("status"))
	if status == vm.NIL {
		status = vm.Int(200)
	}
	body := ress.ValueAt(vm.Keyword("body"))
	resp.WriteHeader(int(status.(vm.Int)))
	respBody, bodyErr := coerceResponseBody(body)
	if bodyErr != nil {
		fmt.Println("HTTP Error coercing body:", bodyErr)
		return
	}
	if respBody != nil {
		_, err = resp.Write(respBody)
	}
	if err != nil {
		fmt.Println("HTTP Error while writing error 500", err)
	}
}

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
	httpGet, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("http/get expects 1-2 args")
		}
		urlStr, err := extractURL(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		req, err := http.NewRequest("GET", urlStr, nil)
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
	httpPost, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
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
		req, err := http.NewRequest("POST", urlStr, strings.NewReader(bodyStr))
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
	httpRequest, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
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
		req, err := http.NewRequest(method, reqURL, bodyReader)
		if err != nil {
			return vm.NIL, err
		}
		hdrs := opts.ValueAt(vm.Keyword("headers"))
		if hdrs != vm.NIL {
			if sq, ok := hdrs.(vm.Sequable); ok {
				for s := sq.Seq(); s != nil; s = s.Next() {
					entry := s.First()
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

	if err != nil {
		panic("http NS init failed")
	}

	ns := vm.NewNamespace("http")

	ns.Def("serve", serve)
	ns.Def("get", httpGet)
	ns.Def("post", httpPost)
	ns.Def("request", httpRequest)
	RegisterNS(ns)
}
