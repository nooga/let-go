package io

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func slurp_lines(arg0 vm.Value) (vm.Value, error) {
	var arg__228_2 vm.Value
	var arg__233_5 vm.Value
	var v6 vm.Value
	var callErr error
	_, _, _ = arg__228_2, arg__233_5, v6
	arg__228_2, callErr = rt.InvokeValue(rt.LookupVar("io", "read-lines").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__233_5, callErr = rt.InvokeValue(rt.LookupVar("io", "read-lines").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "vec").Deref(), []vm.Value{arg__233_5})
	if callErr != nil {
		return nil, callErr
	}
	return v6, nil
}
