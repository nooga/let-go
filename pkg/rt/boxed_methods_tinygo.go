//go:build tinygo

/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"reflect"

	"github.com/nooga/let-go/pkg/vm"
)

// Hand-shimmed *LGWriter methods for TinyGo, where boxed method dispatch
// can't reflect (reflect.Type.Method is unimplemented; see
// pkg/vm/boxed_reflect_tinygo.go). Embedded core's print-method
// implementations call (.WriteString w ...) on every collection print, so
// without these any pr-str/pprint of a map or vector dies at call time.

func lgWriterArg(v vm.Value, method string) (*LGWriter, error) {
	if b, ok := v.(*vm.Boxed); ok {
		if w, ok := b.Unbox().(*LGWriter); ok {
			return w, nil
		}
	}
	return nil, fmt.Errorf("receiver of .%s is not an LGWriter", method)
}

func lgWriterMethod(name string, fn func(w *LGWriter, args []vm.Value) (vm.Value, error)) *vm.NativeFn {
	wrapped, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		w, werr := lgWriterArg(vs[0], name)
		if werr != nil {
			return vm.NIL, werr
		}
		return fn(w, vs[1:])
	})
	// Fail loud at init, matching the tinygo-lane convention: a discarded
	// error would install a nil method that only blows up at call time.
	if err != nil {
		panic(fmt.Errorf("boxed_methods_tinygo: wrapping LGWriter.%s: %w", name, err))
	}
	return wrapped.(*vm.NativeFn)
}

func init() {
	vm.RegisterBoxedMethods(reflect.TypeOf(&LGWriter{}), map[vm.Symbol]*vm.NativeFn{
		vm.Symbol("WriteString"): lgWriterMethod("WriteString", func(w *LGWriter, args []vm.Value) (vm.Value, error) {
			s, ok := args[0].(vm.String)
			if !ok {
				return vm.NIL, fmt.Errorf(".WriteString wants a string argument")
			}
			n, err := w.WriteString(string(s))
			return vm.Int(n), err
		}),
		vm.Symbol("String"): lgWriterMethod("String", func(w *LGWriter, _ []vm.Value) (vm.Value, error) {
			return vm.String(w.String()), nil
		}),
		vm.Symbol("Flush"): lgWriterMethod("Flush", func(w *LGWriter, _ []vm.Value) (vm.Value, error) {
			return vm.NIL, w.Flush()
		}),
	})
}
