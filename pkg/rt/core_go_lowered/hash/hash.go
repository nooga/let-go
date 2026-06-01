package hash

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func xxh3_hasher() (vm.Value, error) {
	var arg__217_1 vm.Value
	var arg__220_4 vm.Value
	var v5 vm.Value
	var callErr error
	_, _, _ = arg__217_1, arg__220_4, v5
	arg__217_1, callErr = rt.InvokeValue(rt.LookupVar("xxh3", "New").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	arg__220_4, callErr = rt.InvokeValue(rt.LookupVar("xxh3", "New").Deref(), []vm.Value{})
	if callErr != nil {
		return nil, callErr
	}
	v5, callErr = rt.InvokeValue(rt.LookupVar("hash", "hasher-map").Deref(), []vm.Value{arg__220_4})
	if callErr != nil {
		return nil, callErr
	}
	return v5, nil
}
func xxh3_hasher_seed(arg0 vm.Value) (vm.Value, error) {
	var arg__224_2 vm.Value
	var arg__229_5 vm.Value
	var v6 vm.Value
	var callErr error
	_, _, _ = arg__224_2, arg__229_5, v6
	arg__224_2, callErr = rt.InvokeValue(rt.LookupVar("xxh3", "NewSeed").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	arg__229_5, callErr = rt.InvokeValue(rt.LookupVar("xxh3", "NewSeed").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	v6, callErr = rt.InvokeValue(rt.LookupVar("hash", "hasher-map").Deref(), []vm.Value{arg__229_5})
	if callErr != nil {
		return nil, callErr
	}
	return v6, nil
}
