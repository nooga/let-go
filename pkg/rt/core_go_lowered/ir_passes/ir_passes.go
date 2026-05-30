package ir_passes

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func replace_refs_BANG_(arg0 vm.Value) (vm.Value, error) {
	var v6 vm.Value
	var callErr error
	_ = v6
	v6, callErr = rt.InvokeValue(rt.LookupVar("ir", "set-refs!").Deref(), []vm.Value{rt.LookupVar("ir.passes", "*current-fn*").Deref(), rt.LookupVar("ir.passes", "*current-inst*").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v6, nil
}
func replace_aux_BANG_(arg0 vm.Value) (vm.Value, error) {
	var v6 vm.Value
	var callErr error
	_ = v6
	v6, callErr = rt.InvokeValue(rt.LookupVar("ir", "set-aux!").Deref(), []vm.Value{rt.LookupVar("ir.passes", "*current-fn*").Deref(), rt.LookupVar("ir.passes", "*current-inst*").Deref(), arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v6, nil
}
func remove_BANG_() (vm.Value, error) {
	var arg__17315_4 vm.Value
	var arg__17322_11 vm.Value
	var v13 vm.Value
	var callErr error
	_, _, _ = arg__17315_4, arg__17322_11, v13
	arg__17315_4, callErr = rt.InvokeValue(rt.LookupVar("ir.zipper", "current-block").Deref(), []vm.Value{rt.LookupVar("ir.passes", "*current-zip*").Deref()})
	if callErr != nil {
		return nil, callErr
	}
	arg__17322_11, callErr = rt.InvokeValue(rt.LookupVar("ir.zipper", "current-block").Deref(), []vm.Value{rt.LookupVar("ir.passes", "*current-zip*").Deref()})
	if callErr != nil {
		return nil, callErr
	}
	v13, callErr = rt.InvokeValue(rt.LookupVar("ir", "remove-inst!").Deref(), []vm.Value{rt.LookupVar("ir.passes", "*current-fn*").Deref(), arg__17322_11, rt.LookupVar("ir.passes", "*current-inst*").Deref()})
	if callErr != nil {
		return nil, callErr
	}
	return v13, nil
}
