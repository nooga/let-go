package hash

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func xxh3_hasher() (vm.Value, error) {
	var arg__212_1 vm.Value
	var arg__215_4 vm.Value
	var v5 vm.Value
	var callErr error
	_, _, _ = arg__212_1, arg__215_4, v5
	arg__212_1, callErr = rt.InvokeValue(rt.LookupVar("xxh3", "New").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__215_4, callErr = rt.InvokeValue(rt.LookupVar("xxh3", "New").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v5, callErr = rt.InvokeValue(rt.LookupVar("hash", "hasher-map").Deref(), []vm.Value{arg__215_4})
	if callErr != nil {
		return nil, callErr
	}
	return v5, nil
}
func xxh3_hasher_seed(arg0 vm.Value) (vm.Value, error) {
	var arg__219_2 vm.Value
	var arg__224_5 vm.Value
	var v6 vm.Value
	var callErr error
	_, _, _ = arg__219_2, arg__224_5, v6
	arg__219_2, callErr = rt.InvokeValue(rt.LookupVar("xxh3", "NewSeed").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__224_5, callErr = rt.InvokeValue(rt.LookupVar("xxh3", "NewSeed").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v6, callErr = rt.InvokeValue(rt.LookupVar("hash", "hasher-map").Deref(), []vm.Value{arg__224_5})
	if callErr != nil {
		return nil, callErr
	}
	return v6, nil
}
