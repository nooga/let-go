package rt

import (
	"net/http"

	"github.com/nooga/let-go/pkg/vm"
)

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

	ns := vm.NewNamespace("http")

	// vars
	CurrentNS = ns.Def("*ns*", ns)

	ns.Def("handle", handle)
	ns.Def("serve", serve)
	RegisterNS(ns)
}
