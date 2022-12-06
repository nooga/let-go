package rt

import (
	"os"

	"github.com/nooga/let-go/pkg/vm"
)

// nolint
func installOsNS() {
	getenv, err := vm.NativeFnType.Box(os.Getenv)

	if err != nil {
		panic("http NS init failed")
	}

	ns := vm.NewNamespace("os")

	// vars
	CurrentNS = ns.Def("*ns*", ns)

	ns.Def("getenv", getenv)
	RegisterNS(ns)
}
