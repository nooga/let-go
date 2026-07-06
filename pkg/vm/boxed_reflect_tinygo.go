//go:build tinygo

/*
 * TinyGo does not implement reflect.Type.Method / reflect.Value.Call, so
 * boxed Go values can't expose their reflected method tables to let-go code
 * under TinyGo. The VM still works; a boxed value's methods just won't resolve
 * unless hand-shimmed here.
 *
 * We hand-register the handful of methods let-go programs actually reach for
 * via the `.` interop operator. Today that's `time.Time.Sub` — let-go's `now`
 * returns a boxed time.Time and the idiomatic monotonic-ms clock is
 * `(quot (.Sub (now) epoch) 1000000)` (see xsofy ui.lg). Add cases here as new
 * boxed-method needs surface under TinyGo.
 */

package vm

import (
	"reflect"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

func reflectMethods(t reflect.Type) map[Symbol]*NativeFn {
	if t == timeType {
		// (.Sub a b) — method value semantics: receiver is the first arg.
		sub, _ := NativeFnType.Wrap(func(vs []Value) (Value, error) {
			a, ok := vs[0].Unbox().(time.Time)
			if !ok {
				return NIL, NewTypeError(vs[0], "is not a time for .Sub on", NativeFnType)
			}
			b, ok := vs[1].Unbox().(time.Time)
			if !ok {
				return NIL, NewTypeError(vs[1], "is not a time for .Sub on", NativeFnType)
			}
			return Int(int64(a.Sub(b))), nil // Duration is int64 ns, matching the reflect path
		})
		return map[Symbol]*NativeFn{Symbol("Sub"): sub.(*NativeFn)}
	}
	return nil
}
