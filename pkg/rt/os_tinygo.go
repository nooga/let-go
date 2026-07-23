//go:build tinygo

/*
 * TinyGo stub: the full os namespace uses reflect.Type.IsVariadic via
 * Box() to register os.Getenv / exec.Command / os.TempDir, and TinyGo's
 * reflect does not implement IsVariadic. We register a minimal subset
 * using Wrap (no reflect) so apps like xsofy can at least call os/exit.
 */

package rt

import (
	"fmt"
	"os"

	"github.com/nooga/let-go/pkg/vm"
)

// os.go's init() that registers this installer is //go:build !tinygo, so the
// tinygo build must register its own — without this the os namespace is never
// installed and os/getenv, os/exit, os/args resolve to nil (e.g. xsofy's
// seed boot-param crashed "nil is not a function" calling os/getenv).
func init() { RegisterInstaller(installOsNS) }

func installOsNS() {
	// mustWrap (not a discarded error): a failed Wrap here would otherwise Def a
	// nil fn and surface as a confusing "nil is not a function" at call time,
	// far from the cause. Fail loud at namespace init instead.
	exitFn := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/exit expects 1 arg")
		}
		code, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("os/exit expected Int")
		}
		os.Exit(int(code))
		return vm.NIL, nil
	})

	// Real os.Getenv, not a blind "". TinyGo's WASI/native targets support it, so
	// env-based config (e.g. xsofy's seed boot-param) actually works there; the
	// wasm target has no environ and os.Getenv returns "", same as the old stub.
	getenvFn := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/getenv expects 1 arg")
		}
		name, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/getenv expected String")
		}
		return vm.String(os.Getenv(string(name))), nil
	})

	ns := vm.NewNamespace("os")
	ns.Def("exit", exitFn)
	ns.Def("getenv", getenvFn)
	ns.Def("args", vm.NewPersistentVector(nil))
	RegisterNS(ns)
}

func lineSeparator() string { return "\n" }

func mustWrap(fn func([]vm.Value) (vm.Value, error)) vm.Value {
	v, err := vm.NativeFnType.Wrap(fn)
	if err != nil {
		panic(err)
	}
	return v
}
