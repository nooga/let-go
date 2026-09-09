/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// defStaticNS returns (creating if needed) a bare namespace that holds host
// static members (e.g. `Util/hash`, `Long/parseLong`). Unlike DefNSBare it does
// NOT auto-refer clojure.core, so defining a member whose name overlaps a core
// var (`Util/hash` vs `clojure.core/hash`) does not print a shadow WARNING at
// every `lg` startup. These namespaces only carry static members; they never
// resolve core names, so the missing refer is harmless.
func defStaticNS(name string) *vm.Namespace {
	name = resolveNSAlias(name)
	nsMu.RLock()
	if e := nsRegistry[name]; e != nil {
		nsMu.RUnlock()
		return e
	}
	nsMu.RUnlock()

	ns := vm.NewNamespace(name)
	nsMu.Lock()
	nsRegistry[name] = ns
	nsMu.Unlock()
	return ns
}

// installJVMStatics registers JVM static methods and (instance? Class x) markers
// that Clojure libraries reach on their :clj branches. Motivated by metosin/malli
// (registry, regex cache, entry parser, string coercion) but each is a general,
// real host-interop resolution backed by a let-go primitive.
func installJVMStatics(ns *vm.Namespace) {
	// --- (instance? Class x) host-class markers ---
	// Map literals evaluate to vm.MapType (let-go.lang.Map), so match that.
	RegisterHostClass("java.util.Map", vm.MapType)
	RegisterHostClass("CharSequence", vm.StringType)
	RegisterHostClass("Pattern", vm.RegexType) // bare java.util.regex.Pattern
	RegisterHostClass("java.util.AbstractList", vm.Symbol("java.util.AbstractList"))
	RegisterHostClass("java.util.Vector", vm.Symbol("java.util.Vector"))

	// --- clojure.lang.LazilyPersistentVector/createOwning(objectArray) -> vector.
	vecFn := ns.Lookup("vec").(*vm.Var).Deref()
	for _, nm := range []string{"LazilyPersistentVector", "clojure.lang.LazilyPersistentVector"} {
		defStaticNS(nm).Def("createOwning", vecFn)
	}

	// --- clojure.lang.PersistentArrayMap/createWithCheck(flatArray) -> map,
	// throwing on a duplicate key. The array is [k0 v0 k1 v1 ...].
	createWithCheck := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("PersistentArrayMap/createWithCheck expects 1 arg")
		}
		arr, ok := vs[0].(*vm.TypedArray)
		if !ok {
			return vm.NIL, fmt.Errorf("PersistentArrayMap/createWithCheck expects an array")
		}
		n := arr.RawCount()
		if n%2 != 0 {
			return vm.NIL, fmt.Errorf("PersistentArrayMap/createWithCheck: odd-length array (%d)", n)
		}
		var m vm.Associative = vm.EmptyPersistentMap
		for i := 0; i+1 < n; i += 2 {
			k := arr.Get(i)
			before := m.(*vm.PersistentMap).RawCount()
			m = m.Assoc(k, arr.Get(i+1))
			if m.(*vm.PersistentMap).RawCount() == before {
				return vm.NIL, fmt.Errorf("duplicate key: %s", k)
			}
		}
		return m, nil
	})
	for _, nm := range []string{"PersistentArrayMap", "clojure.lang.PersistentArrayMap"} {
		defStaticNS(nm).Def("createWithCheck", createWithCheck)
	}

	// --- java.lang.reflect.Array/newInstance(class, n) -> n-element object array
	// (the class argument is ignored — let-go arrays are untyped).
	defStaticNS("Array").Def("newInstance", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) == 2 {
			if n, ok := vs[1].(vm.Int); ok {
				if n < 0 {
					return vm.NIL, fmt.Errorf("Array/newInstance: negative size %d", n)
				}
				return vm.NewObjectArray(int(n)), nil
			}
		}
		return vm.NIL, fmt.Errorf("Array/newInstance expects (class, int)")
	}))

	// --- clojure.lang.Util/hash + Murmur3/hashLong -> let-go hash; Util/hashCombine
	// -> the standard Clojure combine.
	hashOf := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.Int(int64(vm.HashValue(vs[0]))), nil
	})
	hashCombine := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		a, _ := vs[0].(vm.Int)
		b, _ := vs[1].(vm.Int)
		return vm.Int(int64(a) ^ (int64(b) + 0x9e3779b9 + (int64(a) << 6) + (int64(a) >> 2))), nil
	})
	for _, nm := range []string{"Util", "clojure.lang.Util"} {
		defStaticNS(nm).Def("hash", hashOf)
		defStaticNS(nm).Def("hashCombine", hashCombine)
	}
	for _, nm := range []string{"Murmur3", "clojure.lang.Murmur3"} {
		defStaticNS(nm).Def("hashLong", hashOf)
	}

	// --- number/uuid string parses (Long/parseLong, Float/parseFloat, ...). ---
	parseLong := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("parseLong expects a string")
		}
		// strconv.Atoi (returns int) + MakeInt matches let-go's own parse-long
		// and avoids an int64->int conversion (CodeQL "incorrect conversion").
		n, err := strconv.Atoi(string(s))
		if err != nil {
			return vm.NIL, err
		}
		return vm.MakeInt(n), nil
	})
	parseFloat := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("parseFloat expects a string")
		}
		f, err := strconv.ParseFloat(string(s), 64)
		if err != nil {
			return vm.NIL, err
		}
		return vm.Float(f), nil
	})
	defStaticNS("Long").Def("parseLong", parseLong)
	defStaticNS("Integer").Def("parseInt", parseLong)
	defStaticNS("Float").Def("parseFloat", parseFloat)
	defStaticNS("Double").Def("parseDouble", parseFloat)

	// --- java.util.Locale constants. Code reaches for these to force
	// locale-independent casing (honeysql upper-cases SQL keywords with
	// toUpperCase(Locale/US) to dodge the Turkish-I trap). Go's case mapping
	// is locale-independent already, so the constants are opaque markers the
	// string .toUpperCase/.toLowerCase methods accept and ignore.
	// getDefault/setDefault/forLanguageTag are honest no-ops over the same
	// markers: a default-locale change genuinely cannot affect let-go's case
	// mapping, which is what code testing locale independence asserts.
	localeGetDefault := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.Symbol("java.util.Locale/US"), nil
	})
	localeSetDefault := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("Locale/setDefault expects 1 argument")
		}
		return vm.NIL, nil
	})
	localeForTag := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("Locale/forLanguageTag expects 1 argument")
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("Locale/forLanguageTag expects a string")
		}
		return vm.Symbol("java.util.Locale/" + string(s)), nil
	})
	for _, nm := range []string{"Locale", "java.util.Locale"} {
		lns := defStaticNS(nm)
		for _, c := range []string{"US", "ROOT", "ENGLISH"} {
			lns.Def(c, vm.Symbol("java.util.Locale/"+c))
		}
		lns.Def("getDefault", localeGetDefault)
		lns.Def("setDefault", localeSetDefault)
		lns.Def("forLanguageTag", localeForTag)
	}

	// --- java.net.URLEncoder/encode. Go's url.QueryEscape implements the
	// same application/x-www-form-urlencoded rules Java does (space -> '+').
	// The optional second argument is a charset; let-go strings are UTF-8 by
	// construction, so only UTF-8 is accepted.
	urlEncode := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 && len(vs) != 2 {
			return vm.NIL, fmt.Errorf("URLEncoder/encode expects 1 or 2 arguments, got %d", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("URLEncoder/encode expects a string")
		}
		if len(vs) == 2 {
			enc, ok := vs[1].(vm.String)
			if !ok {
				return vm.NIL, fmt.Errorf("URLEncoder/encode charset must be a string")
			}
			e := strings.ToUpper(string(enc))
			if e != "UTF-8" && e != "UTF8" {
				return vm.NIL, fmt.Errorf("URLEncoder/encode: only UTF-8 supported, got %q", string(enc))
			}
		}
		return vm.String(url.QueryEscape(string(s))), nil
	})
	for _, nm := range []string{"URLEncoder", "java.net.URLEncoder"} {
		defStaticNS(nm).Def("encode", urlEncode)
	}

	// bare UUID/fromString (lang.go registers only the fully-qualified
	// java.util.UUID namespace).
	defStaticNS("UUID").Def("fromString", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("UUID/fromString expects a string")
		}
		u := vm.ParseUUID(string(s))
		if u == nil {
			return vm.NIL, fmt.Errorf("invalid UUID: %q", string(s))
		}
		return u, nil
	}))
}

// installMathStatics registers java.lang.Math. Go's math package is close but
// not a drop-in: Math/round returns a long and rounds half toward POSITIVE
// infinity (JVM: floor(x + 0.5)), where Go's math.Round returns a float64 and
// rounds half away from zero, so -2.5 differs (-2 on the JVM, -3 in Go).
// Math/abs preserves long-in/long-out, and Math/getExponent has no Go
// equivalent at all. Non-finite inputs propagate rather than raising, as on
// the JVM.
func installMathStatics() {
	ns := defStaticNS("Math")
	for _, nm := range []string{"Math", "java.lang.Math"} {
		if nm != "Math" {
			ns = defStaticNS(nm)
		}
		mathNS := ns

		unary := func(name string, fn func(float64) float64) {
			mathNS.Def(name, mustWrap(func(vs []vm.Value) (vm.Value, error) {
				if len(vs) != 1 {
					return vm.NIL, fmt.Errorf("%s expects 1 arg", name)
				}
				f, ok := vm.ToFloat(vs[0])
				if !ok {
					return vm.NIL, fmt.Errorf("%s expected number, got %s", name, vs[0].Type().Name())
				}
				return vm.Float(fn(float64(f))), nil
			}))
		}
		unary("floor", math.Floor)
		unary("ceil", math.Ceil)
		unary("sqrt", math.Sqrt)
		unary("log", math.Log)
		unary("exp", math.Exp)

		// Go's math.Pow follows IEEE 754, which differs from java.lang.Math on
		// three cases involving a base of magnitude 1: IEEE says pow(1, y) is 1
		// for every y, while Java returns NaN for a NaN or infinite exponent.
		mathNS.Def("pow", mustWrap(func(vs []vm.Value) (vm.Value, error) {
			if len(vs) != 2 {
				return vm.NIL, fmt.Errorf("pow expects 2 args")
			}
			bv, ok1 := vm.ToFloat(vs[0])
			ev, ok2 := vm.ToFloat(vs[1])
			if !ok1 || !ok2 {
				return vm.NIL, fmt.Errorf("pow expected numbers")
			}
			b, e := float64(bv), float64(ev)
			switch {
			case e == 0:
				// Java: a zero exponent is 1.0 even for a NaN base.
				return vm.Float(1), nil
			case math.IsNaN(e):
				return vm.Float(math.NaN()), nil
			case math.Abs(b) == 1 && math.IsInf(e, 0):
				return vm.Float(math.NaN()), nil
			}
			return vm.Float(math.Pow(b, e)), nil
		}))

		// long in -> long out, double in -> double out, matching the JVM's
		// overloads. Returning a float for an integer argument would break
		// integer callers.
		mathNS.Def("abs", mustWrap(func(vs []vm.Value) (vm.Value, error) {
			if len(vs) != 1 {
				return vm.NIL, fmt.Errorf("abs expects 1 arg")
			}
			if i, ok := vs[0].(vm.Int); ok {
				if int64(i) < 0 {
					return vm.MakeInt(int(-int64(i))), nil
				}
				return i, nil
			}
			f, ok := vm.ToFloat(vs[0])
			if !ok {
				return vm.NIL, fmt.Errorf("abs expected number, got %s", vs[0].Type().Name())
			}
			return vm.Float(math.Abs(float64(f))), nil
		}))

		// JVM Math.round: half rounds toward POSITIVE infinity, so -2.5 is -2
		// (math.Round would give -3). NaN is 0, and out-of-range saturates to
		// Long.MAX_VALUE / Long.MIN_VALUE rather than wrapping.
		//
		// Do NOT implement this as floor(x + 0.5): the addition rounds before
		// floor runs, so 0.49999999999999994 would give 1 instead of 0 and
		// 4503599627370497.0 would gain one. Java itself abandoned that formula
		// for this reason (JDK-6430675). Compare the fraction against 0.5
		// instead, which never perturbs the input.
		mathNS.Def("round", mustWrap(func(vs []vm.Value) (vm.Value, error) {
			if len(vs) != 1 {
				return vm.NIL, fmt.Errorf("round expects 1 arg")
			}
			f, ok := vm.ToFloat(vs[0])
			if !ok {
				return vm.NIL, fmt.Errorf("round expected number, got %s", vs[0].Type().Name())
			}
			x := float64(f)
			if math.IsNaN(x) {
				return vm.MakeInt(0), nil
			}
			r := math.Floor(x)
			if x-r >= 0.5 {
				r++
			}
			// JVM overloads on the argument type: round(float) returns an int
			// and saturates at Integer.MAX_VALUE, round(double) returns a long.
			if _, isF32 := vs[0].(vm.Float32); isF32 {
				if r >= 2147483647.0 {
					return vm.MakeInt(math.MaxInt32), nil
				}
				if r <= -2147483648.0 {
					return vm.MakeInt(math.MinInt32), nil
				}
				return vm.MakeInt(int(r)), nil
			}
			// float64(math.MaxInt64) rounds up to 2^63, so compare against that
			// boundary rather than the constant itself.
			if r >= 9223372036854775808.0 {
				return longCompatValue(math.MaxInt64), nil
			}
			if r <= -9223372036854775808.0 {
				return longCompatValue(math.MinInt64), nil
			}
			return longCompatValue(int64(r)), nil
		}))

		mathNS.Def("scalb", mustWrap(func(vs []vm.Value) (vm.Value, error) {
			if len(vs) != 2 {
				return vm.NIL, fmt.Errorf("scalb expects 2 args")
			}
			f, ok1 := vm.ToFloat(vs[0])
			e, ok2 := vm.ToInt(vs[1])
			if !ok1 || !ok2 {
				return vm.NIL, fmt.Errorf("scalb expected (double, int)")
			}
			return vm.Float(math.Ldexp(float64(f), int(e))), nil
		}))

		// JVM Math.getExponent: unbiased exponent as an int. Zero and
		// subnormals give MIN_EXPONENT-1 (-1023); NaN and the infinities give
		// MAX_EXPONENT+1 (1024). Go has no equivalent, so read the biased
		// exponent bits directly.
		mathNS.Def("getExponent", mustWrap(func(vs []vm.Value) (vm.Value, error) {
			if len(vs) != 1 {
				return vm.NIL, fmt.Errorf("getExponent expects 1 arg")
			}
			f, ok := vm.ToFloat(vs[0])
			if !ok {
				return vm.NIL, fmt.Errorf("getExponent expected number, got %s", vs[0].Type().Name())
			}
			x := float64(f)
			// JVM overloads here too: the float form uses the single-precision
			// bias (127), so zero gives -127 and NaN/infinity give 128.
			if _, isF32 := vs[0].(vm.Float32); isF32 {
				if math.IsNaN(x) || math.IsInf(x, 0) {
					return vm.MakeInt(128), nil
				}
				biased := int((math.Float32bits(float32(x)) >> 23) & 0xFF)
				if biased == 0 { // zero or subnormal
					return vm.MakeInt(-127), nil
				}
				return vm.MakeInt(biased - 127), nil
			}
			if math.IsNaN(x) || math.IsInf(x, 0) {
				return vm.MakeInt(1024), nil
			}
			biased := int((math.Float64bits(x) >> 52) & 0x7FF)
			if biased == 0 { // zero or subnormal
				return vm.MakeInt(-1023), nil
			}
			return vm.MakeInt(biased - 1023), nil
		}))
	}
}
