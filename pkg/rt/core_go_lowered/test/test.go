package test

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func testing_contexts_str() (vm.Value, error) {
	var arg__31473_6 vm.Value
	var arg__31481_14 vm.Value
	var arg__31482_15 vm.Value
	var arg__31490_23 vm.Value
	var arg__31498_31 vm.Value
	var arg__31499_32 vm.Value
	var v33 vm.Value
	var callErr error
	_, _, _, _, _, _, _ = arg__31473_6, arg__31481_14, arg__31482_15, arg__31490_23, arg__31498_31, arg__31499_32, v33
	arg__31473_6, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "interpose").Deref(), []vm.Value{vm.String(" > "), rt.LookupVar("clojure.test", "*testing-contexts*").Deref()})
	if callErr != nil {
		return nil, callErr
	}
	arg__31481_14, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "interpose").Deref(), []vm.Value{vm.String(" > "), rt.LookupVar("clojure.test", "*testing-contexts*").Deref()})
	if callErr != nil {
		return nil, callErr
	}
	arg__31482_15, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), arg__31481_14})
	if callErr != nil {
		return nil, callErr
	}
	arg__31490_23, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "interpose").Deref(), []vm.Value{vm.String(" > "), rt.LookupVar("clojure.test", "*testing-contexts*").Deref()})
	if callErr != nil {
		return nil, callErr
	}
	arg__31498_31, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "interpose").Deref(), []vm.Value{vm.String(" > "), rt.LookupVar("clojure.test", "*testing-contexts*").Deref()})
	if callErr != nil {
		return nil, callErr
	}
	arg__31499_32, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "apply").Deref(), []vm.Value{rt.LookupVar("clojure.core", "str").Deref(), arg__31498_31})
	if callErr != nil {
		return nil, callErr
	}
	v33, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "str").Deref(), []vm.Value{arg__31499_32})
	if callErr != nil {
		return nil, callErr
	}
	return v33, nil
}
func apply_template(arg0 vm.Value, arg1 vm.Value, arg2 vm.Value) (vm.Value, error) {
	var arg__31545_4 vm.Value
	var arg__31553_7 vm.Value
	var v8 vm.Value
	var callErr error
	_, _, _ = arg__31545_4, arg__31553_7, v8
	arg__31545_4, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zipmap").Deref(), []vm.Value{arg0, arg2})
	if callErr != nil {
		return nil, callErr
	}
	arg__31553_7, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "zipmap").Deref(), []vm.Value{arg0, arg2})
	if callErr != nil {
		return nil, callErr
	}
	v8, callErr = rt.InvokeValue(rt.LookupVar("clojure.walk", "postwalk-replace").Deref(), []vm.Value{arg__31553_7, arg1})
	if callErr != nil {
		return nil, callErr
	}
	return v8, nil
}
