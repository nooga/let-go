//go:build tinygo && wasm

package rt

import "github.com/nooga/let-go/pkg/vm"

// The lginterop-generated interop_xxh3.go is `!tinygo` (reflect-boxed), so
// under TinyGo the `xxh3` namespace is registered from the Wrap-based
// adapters in hash_xxh3_wrap.go instead. Same names, same values.
func installXxh3NS() {
	ns := vm.NewNamespace("xxh3")
	ns.Def("Hash", mustWrap(xxh3Hash))
	ns.Def("HashSeed", mustWrap(xxh3HashSeed))
	ns.Def("HashString", mustWrap(xxh3HashString))
	ns.Def("HashStringSeed", mustWrap(xxh3HashStringSeed))
	RegisterNS(ns)
}

func init() { RegisterInstaller(installXxh3NS) }
