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
	req := vm.EmptyPersistentMap
	req = req.Assoc(vm.Keyword("request-method"), methodToLG(request.Method)).(*vm.PersistentMap)
	url := request.URL

	if request.TLS == nil {
		req = req.Assoc(vm.Keyword("scheme"), vm.Keyword("http")).(*vm.PersistentMap)
	} else {
		req = req.Assoc(vm.Keyword("scheme"), vm.Keyword("https")).(*vm.PersistentMap)
	}
	req = req.Assoc(vm.Keyword("uri"), vm.String(url.RequestURI())).(*vm.PersistentMap)
	req = req.Assoc(vm.Keyword("query-string"), vm.String(url.RawQuery)).(*vm.PersistentMap)
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
	req = req.Assoc(vm.Keyword("body"), vm.String(bytes)).(*vm.PersistentMap)
	req = req.Assoc(vm.Keyword("remote-addr"), vm.String(request.RemoteAddr)).(*vm.PersistentMap)
	req = req.Assoc(vm.Keyword("server-addr"), vm.String(request.Host)).(*vm.PersistentMap)
	req = req.Assoc(vm.Keyword("server-port"), vm.String(url.Port())).(*vm.PersistentMap)

	if len(request.Header) > 0 {
		hs := vm.EmptyPersistentMap
		for k, v := range request.Header {
			hs = hs.Assoc(vm.String(strings.ToLower(k)), vm.String(strings.Join(v, ","))).(*vm.PersistentMap)
		}
		req = req.Assoc(vm.Keyword("headers"), hs).(*vm.PersistentMap)
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
	headers := ress.ValueAt(vm.Keyword("headers"))
	if headers != vm.NIL {
		if sq, ok := headers.(vm.Sequable); ok {
			for s := sq.Seq(); s != nil && s != vm.EmptyList; s = s.Next() {
				entry := s.First().(vm.ArrayVector)
				head.Add(entry[0].Unbox().(string), entry[1].Unbox().(string))
			}
		}
	}
	status := ress.ValueAt(vm.Keyword("status"))
	if status == vm.NIL {
		status = vm.Int(200)
	}
	body := ress.ValueAt(vm.Keyword("body"))
	if body == vm.NIL {
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

	ns.Def("handle", handle)
	ns.Def("serve", serve)
	ns.Def("serve2", serve2)
	RegisterNS(ns)
}
