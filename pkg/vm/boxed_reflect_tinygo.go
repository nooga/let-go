//go:build tinygo

/*
 * TinyGo does not implement reflect.Type.Method / reflect.Value.Call, so
 * boxed Go values can't expose their reflected method tables to let-go code
 * under TinyGo. The VM still works; a boxed value's methods just won't resolve
 * unless hand-shimmed here.
 *
 * We hand-register the handful of methods let-go programs actually reach for
 * via the `.` interop operator. Today that's two methods on time.Time, which
 * let-go's `now` returns boxed: `Sub`, for the monotonic-ms clock idiom
 * `(quot (.Sub (now) epoch) 1000000)` (see xsofy ui.lg), and `UnixMilli`, which
 * core's `inst-ms` calls. Add cases here as new boxed-method needs surface
 * under TinyGo.
 */

package vm

import (
	"fmt"
	"reflect"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

// extraBoxedMethods holds hand-shimmed method tables registered by other
// packages, which can't be special-cased here without import cycles — e.g.
// pkg/rt's *LGWriter. Populated at init time, read-only afterwards.
var extraBoxedMethods = map[reflect.Type]map[Symbol]*NativeFn{}

// RegisterBoxedMethods installs a hand-shimmed method table for t. TinyGo
// can't reflect method tables (reflect.Type.Method is unimplemented), so any
// boxed type whose methods let-go code calls via `.` needs an entry.
func RegisterBoxedMethods(t reflect.Type, ms map[Symbol]*NativeFn) {
	extraBoxedMethods[t] = ms
}

func reflectMethods(t reflect.Type) map[Symbol]*NativeFn {
	if ms, ok := extraBoxedMethods[t]; ok {
		return ms
	}
	if t == timeType {
		return timeMethods()
	}
	return nil
}

// timeMethods shims time.Time's methods with method value semantics: the
// receiver is the first arg. Results match what the reflect path returns.
func timeMethods() map[Symbol]*NativeFn {
	return map[Symbol]*NativeFn{
		// (.Sub a b): a Duration, int64 nanoseconds.
		Symbol("Sub"): wrapTimeMethod("Sub", func(vs []Value) (Value, error) {
			a, ok := vs[0].Unbox().(time.Time)
			if !ok {
				return NIL, NewTypeError(vs[0], "is not a time for .Sub on", NativeFnType)
			}
			b, ok := vs[1].Unbox().(time.Time)
			if !ok {
				return NIL, NewTypeError(vs[1], "is not a time for .Sub on", NativeFnType)
			}
			return Int(int64(a.Sub(b))), nil
		}),
		// (.UnixMilli t): epoch milliseconds, an int64.
		Symbol("UnixMilli"): wrapTimeMethod("UnixMilli", func(vs []Value) (Value, error) {
			a, ok := vs[0].Unbox().(time.Time)
			if !ok {
				return NIL, NewTypeError(vs[0], "is not a time for .UnixMilli on", NativeFnType)
			}
			return Int(a.UnixMilli()), nil
		}),
	}
}

// wrapTimeMethod fails loud: a discarded error would install a nil method that
// only blows up at call time, off in the tinygo build.
func wrapTimeMethod(name string, f func(vs []Value) (Value, error)) *NativeFn {
	fn, err := NativeFnType.Wrap(f)
	if err != nil {
		panic(fmt.Errorf("reflectMethods: wrapping time.%s: %w", name, err))
	}
	return fn.(*NativeFn)
}

// methodLookupError explains a failed InvokeMethod. Under TinyGo we cannot
// enumerate a type's methods (reflect.Type.Method is unimplemented), so a type
// that actually has methods reaches here with an empty table. NumMethod is
// implemented, so we can still tell "has methods we can't reflect" apart from
// "genuinely no methods" and report the former honestly instead of claiming the
// value has no methods at all.
func methodLookupError(bt *aBoxedType, name Symbol) error {
	if bt.typ.NumMethod() > 0 {
		return fmt.Errorf("method %s on %v is unavailable under TinyGo: reflect.Type.Method "+
			"is unimplemented, so boxed Go methods must be hand-shimmed in "+
			"boxed_reflect_tinygo.go", name, bt)
	}
	return fmt.Errorf("%v doesn't have any methods", bt)
}
