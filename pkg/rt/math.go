/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/nooga/let-go/pkg/vm"
)

// toDouble coerces a Value to float64, returning an error if not numeric.
func toDouble(v vm.Value) (float64, error) {
	f, ok := vm.ToFloat(v)
	if !ok {
		return 0, fmt.Errorf("expected numeric value, got %s", v.Type().Name())
	}
	return f, nil
}

// toLong coerces a Value to int64, returning an error if not numeric.
func toLong(v vm.Value) (int64, error) {
	switch n := v.(type) {
	case vm.Int:
		return int64(n), nil
	case vm.Float:
		return int64(n), nil
	case *vm.BigInt:
		v, ok := n.ToInt64()
		if !ok {
			return 0, fmt.Errorf("bigint too large for long")
		}
		return v, nil
	}
	return 0, fmt.Errorf("expected numeric value, got %s", v.Type().Name())
}

// mathFn1 wraps a unary float64→float64 Go math function.
func mathFn1(f func(float64) float64) vm.Value {
	v, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.Float(f(a)), nil
	})
	return v
}

// mathFn2 wraps a binary (float64,float64)→float64 Go math function.
func mathFn2(f func(float64, float64) float64) vm.Value {
	v, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		b, err := toDouble(vs[1])
		if err != nil {
			return vm.NIL, err
		}
		return vm.Float(f(a, b)), nil
	})
	return v
}

// nolint
func installMathNS() {
	ns := vm.NewNamespace("math")
	ns.Refer(CoreNS, "", true)

	// Constants
	ns.Def("E", vm.Float(math.E))
	ns.Def("PI", vm.Float(math.Pi))

	// Trigonometric
	ns.Def("sin", mathFn1(math.Sin))
	ns.Def("cos", mathFn1(math.Cos))
	ns.Def("tan", mathFn1(math.Tan))
	ns.Def("asin", mathFn1(math.Asin))
	ns.Def("acos", mathFn1(math.Acos))
	ns.Def("atan", mathFn1(math.Atan))
	ns.Def("atan2", mathFn2(math.Atan2))

	// Hyperbolic
	ns.Def("sinh", mathFn1(math.Sinh))
	ns.Def("cosh", mathFn1(math.Cosh))
	ns.Def("tanh", mathFn1(math.Tanh))

	// Exponential & logarithmic
	ns.Def("exp", mathFn1(math.Exp))
	ns.Def("log", mathFn1(math.Log))
	ns.Def("log10", mathFn1(math.Log10))
	ns.Def("expm1", mathFn1(math.Expm1))
	ns.Def("log1p", mathFn1(math.Log1p))
	ns.Def("pow", mathFn2(math.Pow))

	// Rounding — ceil, floor, rint return double (matching Clojure)
	ns.Def("ceil", mathFn1(math.Ceil))
	ns.Def("floor", mathFn1(math.Floor))
	ns.Def("rint", mathFn1(math.RoundToEven))

	// round returns long (matching Clojure)
	roundf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.MakeInt(int(math.Round(a))), nil
	})
	ns.Def("round", roundf)

	// Root functions
	ns.Def("sqrt", mathFn1(math.Sqrt))
	ns.Def("cbrt", mathFn1(math.Cbrt))

	// Exact arithmetic — throw on overflow (operates on longs)
	addExact, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toLong(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		b, err := toLong(vs[1])
		if err != nil {
			return vm.NIL, err
		}
		r := a + b
		// Overflow: signs of a and b are same, but result sign differs
		if (a^b) >= 0 && (a^r) < 0 {
			return vm.NIL, fmt.Errorf("integer overflow")
		}
		return vm.MakeInt(int(r)), nil
	})
	ns.Def("add-exact", addExact)

	subtractExact, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toLong(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		b, err := toLong(vs[1])
		if err != nil {
			return vm.NIL, err
		}
		r := a - b
		// Overflow: signs of a and b differ, and result sign differs from a
		if (a^b) < 0 && (a^r) < 0 {
			return vm.NIL, fmt.Errorf("integer overflow")
		}
		return vm.MakeInt(int(r)), nil
	})
	ns.Def("subtract-exact", subtractExact)

	multiplyExact, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toLong(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		b, err := toLong(vs[1])
		if err != nil {
			return vm.NIL, err
		}
		r := a * b
		if a != 0 && r/a != b {
			return vm.NIL, fmt.Errorf("integer overflow")
		}
		return vm.MakeInt(int(r)), nil
	})
	ns.Def("multiply-exact", multiplyExact)

	incrementExact, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toLong(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		if a == math.MaxInt64 {
			return vm.NIL, fmt.Errorf("integer overflow")
		}
		return vm.MakeInt(int(a + 1)), nil
	})
	ns.Def("increment-exact", incrementExact)

	decrementExact, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toLong(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		if a == math.MinInt64 {
			return vm.NIL, fmt.Errorf("integer overflow")
		}
		return vm.MakeInt(int(a - 1)), nil
	})
	ns.Def("decrement-exact", decrementExact)

	negateExact, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toLong(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		if a == math.MinInt64 {
			return vm.NIL, fmt.Errorf("integer overflow")
		}
		return vm.MakeInt(int(-a)), nil
	})
	ns.Def("negate-exact", negateExact)

	// floor-div and floor-mod — integer division rounding toward negative infinity
	floorDiv, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toLong(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		b, err := toLong(vs[1])
		if err != nil {
			return vm.NIL, err
		}
		if b == 0 {
			return vm.NIL, fmt.Errorf("divide by zero")
		}
		d := a / b
		// Adjust for floored division: if signs differ and remainder is non-zero
		if (a^b) < 0 && d*b != a {
			d--
		}
		return vm.MakeInt(int(d)), nil
	})
	ns.Def("floor-div", floorDiv)

	floorMod, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toLong(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		b, err := toLong(vs[1])
		if err != nil {
			return vm.NIL, err
		}
		if b == 0 {
			return vm.NIL, fmt.Errorf("divide by zero")
		}
		r := a % b
		if r != 0 && (r > 0) != (b > 0) {
			r += b
		}
		return vm.MakeInt(int(r)), nil
	})
	ns.Def("floor-mod", floorMod)

	// Floating-point utilities
	ns.Def("copy-sign", mathFn2(math.Copysign))
	ns.Def("hypot", mathFn2(math.Hypot))

	// signum — returns -1.0, 0.0, or 1.0
	signumf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		if math.IsNaN(a) {
			return vm.Float(math.NaN()), nil
		}
		if a > 0 {
			return vm.Float(1.0), nil
		}
		if a < 0 {
			return vm.Float(-1.0), nil
		}
		return vm.Float(a), nil // preserve +0.0 / -0.0
	})
	ns.Def("signum", signumf)

	ns.Def("IEEE-remainder", mathFn2(math.Remainder))

	// ulp — unit in last place
	ulpf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		if math.IsNaN(a) {
			return vm.Float(math.NaN()), nil
		}
		if math.IsInf(a, 0) {
			return vm.Float(math.Inf(1)), nil
		}
		a = math.Abs(a)
		next := math.Nextafter(a, math.Inf(1))
		return vm.Float(next - a), nil
	})
	ns.Def("ulp", ulpf)

	// get-exponent — unbiased exponent of a double
	getExponent, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		if math.IsNaN(a) || math.IsInf(a, 0) {
			return vm.MakeInt(1024), nil // matches Java Float.MAX_EXPONENT + 1
		}
		if a == 0 {
			return vm.MakeInt(-1023), nil // matches Java Float.MIN_EXPONENT - 1
		}
		_, exp := math.Frexp(a)
		return vm.MakeInt(exp - 1), nil // Frexp returns 0.5*2^exp, Java uses 1.x*2^exp
	})
	ns.Def("get-exponent", getExponent)

	// scalb — d × 2^scaleFactor
	scalbf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		sf, err := toLong(vs[1])
		if err != nil {
			return vm.NIL, err
		}
		return vm.Float(math.Ldexp(a, int(sf))), nil
	})
	ns.Def("scalb", scalbf)

	// next-up — adjacent double toward +Inf
	nextUpf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.Float(math.Nextafter(a, math.Inf(1))), nil
	})
	ns.Def("next-up", nextUpf)

	// next-down — adjacent double toward -Inf
	nextDownf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.Float(math.Nextafter(a, math.Inf(-1))), nil
	})
	ns.Def("next-down", nextDownf)

	// next-after — adjacent float toward direction
	ns.Def("next-after", mathFn2(math.Nextafter))

	// Angle conversion
	toRadians, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.Float(a * math.Pi / 180.0), nil
	})
	ns.Def("to-radians", toRadians)

	toDegrees, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		a, err := toDouble(vs[0])
		if err != nil {
			return vm.NIL, err
		}
		return vm.Float(a * 180.0 / math.Pi), nil
	})
	ns.Def("to-degrees", toDegrees)

	// random — pseudorandom double in [0.0, 1.0)
	// Note: Clojure's Math/random is a Java static method; we expose it as (math/random)
	randomf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		return vm.Float(rand.Float64()), nil
	})
	ns.Def("random", randomf)

	// abs — works on both long and double, returns same type
	absf, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("wrong number of arguments %d", len(vs))
		}
		switch n := vs[0].(type) {
		case vm.Int:
			if int(n) < 0 {
				return vm.MakeInt(-int(n)), nil
			}
			return n, nil
		case vm.Float:
			return vm.Float(math.Abs(float64(n))), nil
		case *vm.BigInt:
			return vm.NumAbs(n)
		}
		return vm.NIL, fmt.Errorf("expected numeric value, got %s", vs[0].Type().Name())
	})
	ns.Def("abs", absf)

	RegisterNS(ns)
}
