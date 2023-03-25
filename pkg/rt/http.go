/*
 * Copyright (c) 2022 Marcin Gasperowicz <xnooga@gmail.com>
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
	req := vm.Map{}

	req[vm.Keyword("request-method")] = methodToLG(request.Method)
	url := request.URL

	if request.TLS == nil {
		req[vm.Keyword("scheme")] = vm.Keyword("http")
	} else {
		req[vm.Keyword("scheme")] = vm.Keyword("https")
	}
	req[vm.Keyword("uri")] = vm.String(url.RequestURI())
	req[vm.Keyword("query-string")] = vm.String(url.RawQuery)
	defer request.Body.Close()
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		resp.WriteHeader(500)
		_, err := resp.Write([]byte(fmt.Sprintf("%s", err)))
		if err != nil {
			fmt.Println("HTTP Error while writing error 500", err)
		}
		return
	}
	req[vm.Keyword("body")] = vm.String(bytes)
	req[vm.Keyword("remote-addr")] = vm.String(request.RemoteAddr)
	req[vm.Keyword("server-addr")] = vm.String(request.Host)
	req[vm.Keyword("server-port")] = vm.String(url.Port())

	if len(request.Header) > 0 {
		hs := vm.Map{}
		for k, v := range request.Header {
			hs[vm.String(strings.ToLower(k))] = vm.String(strings.Join(v, ","))
		}
		req[vm.Keyword("headers")] = hs
	}

	res, err := h.fn.Invoke([]vm.Value{req})
	if err != nil {
		resp.WriteHeader(500)
		_, err := resp.Write([]byte(fmt.Sprintf("%s", err)))
		if err != nil {
			fmt.Println("HTTP Error while writing error 500", err)
		}
		return
	}

	ress, ok := res.(vm.Map)
	if !ok {
		resp.WriteHeader(500)
		_, err := resp.Write([]byte("handler returned malformed response"))
		if err != nil {
			fmt.Println("HTTP Error while writing error 500", err)
		}
		return
	}
	head := resp.Header()
	headers, ok := ress[vm.Keyword("headers")]
	if ok {
		for k, v := range headers.(vm.Map) {
			head.Add(k.Unbox().(string), v.Unbox().(string))
		}
	}
	status, ok := ress[vm.Keyword("status")]
	if !ok {
		status = vm.Int(200)
	}
	body, ok := ress[vm.Keyword("body")]
	if !ok {
		body = vm.String("")
	}
	resp.WriteHeader(int(status.(vm.Int)))
	_, err = resp.Write([]byte(body.(vm.String)))
	if err != nil {
		fmt.Println("HTTP Error while writing error 500", err)
	}
}

// nolint
func installHttpNS() {
	// FIXME, this should box the function directly
	handle, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		fnm := vs[1].Unbox().(func(interface{}))
		var fn func(w http.ResponseWriter, r *http.Request) interface{}
		fnm(&fn)
		http.HandleFunc(vs[0].Unbox().(string), func(w http.ResponseWriter, r *http.Request) {
			fn(w, r)
		})
		return vm.NIL, nil
	})

	if err != nil {
		panic("http NS init failed")
	}

	serve, err := vm.NativeFnType.Box(http.ListenAndServe)
	if err != nil {
		panic("http NS init failed")
	}

	serve2, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, vm.NewExecutionError("serve expects 2 args")
		}
		addr, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, vm.NewExecutionError("serve expected listen address as String")
		}
		handlerFunc, ok := vs[0].(vm.Fn)
		if !ok {
			return vm.NIL, vm.NewExecutionError("serve expected handler function as Fn")
		}
		handler := &Handler{fn: handlerFunc}
		http.ListenAndServe(string(addr), handler)
		return vm.NIL, nil
	})
	if err != nil {
		panic("http NS init failed")
	}

	ns := vm.NewNamespace("http")

	// vars
	CurrentNS = ns.Def("*ns*", ns)

	ns.Def("handle", handle)
	ns.Def("serve", serve)
	ns.Def("serve2", serve2)
	RegisterNS(ns)
}
