//go:build tinygo

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
)

// boxReflectFunc under TinyGo cannot reflect a Go function's signature or call
// through it — reflect.Type.{IsVariadic,NumIn,In} and reflect.Value.Call are
// all unimplemented in TinyGo's reflect. Rather than trap at boot when an
// interop namespace (e.g. xxh3) boxes typed Go funcs, install a stub NativeFn:
// it is present in the namespace but returns a clear error if a program ever
// invokes it. Programs that don't call reflect-boxed interop funcs run fine.
func boxReflectFunc(_ *theNativeFnType, fn any, ty reflect.Type) (Value, error) {
	name := ty.String()
	proxy := func(_ []Value) (Value, error) {
		return NIL, fmt.Errorf("native fn %s: Go functions boxed via reflection are not callable under TinyGo (reflect.Value.Call unimplemented)", name)
	}
	return &NativeFn{
		arity:       -1,
		isVariadric: true,
		fn:          fn,
		proxy:       proxy,
	}, nil
}
