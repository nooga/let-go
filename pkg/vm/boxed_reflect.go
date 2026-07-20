//go:build !tinygo

package vm

import (
	"fmt"
	"reflect"
)

func reflectMethods(t reflect.Type) map[Symbol]*NativeFn {
	n := t.NumMethod()
	if n == 0 {
		return nil
	}
	out := map[Symbol]*NativeFn{}
	for i := 0; i < n; i++ {
		m := t.Method(i)
		me, err := NativeFnType.Box(m.Func.Interface())
		if err != nil {
			fmt.Println(t.Name(), "boxing method failed", err)
			continue
		}
		mef, ok := me.(*NativeFn)
		if !ok {
			fmt.Println(t.Name(), "boxed method is not a native fn")
			continue
		}
		out[Symbol(m.Name)] = mef
	}
	return out
}

// methodLookupError explains a failed InvokeMethod. On stock Go, reflect sees
// the full method set, so a type with no methods truly has none and a missing
// name is a genuine typo.
func methodLookupError(bt *aBoxedType, name Symbol) error {
	if bt.typ.NumMethod() == 0 {
		return fmt.Errorf("%v doesn't have any methods", bt)
	}
	return fmt.Errorf("method %s not found in %v", name, bt)
}

// RegisterBoxedMethods is a no-op on reflect-capable builds: the full method
// table comes from reflection. The TinyGo build (boxed_reflect_tinygo.go)
// consults this registry because reflect.Type.Method is unimplemented there.
func RegisterBoxedMethods(t reflect.Type, ms map[Symbol]*NativeFn) {}
