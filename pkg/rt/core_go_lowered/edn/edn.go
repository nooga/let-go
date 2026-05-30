package edn

import (
	rt "github.com/nooga/let-go/pkg/rt"
	vm "github.com/nooga/let-go/pkg/vm"
)

func write_string(arg0 vm.Value) (vm.Value, error) {
	var v2 vm.Value
	var callErr error
	_ = v2
	v2, callErr = rt.InvokeValue(rt.LookupVar("clojure.core", "pr-str").Deref(), []vm.Value{arg0})
	if callErr != nil {
		return nil, callErr
	}
	return v2, nil
}
