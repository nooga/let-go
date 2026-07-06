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
