/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"unicode"

	"github.com/nooga/let-go/pkg/vm"
)

// installCharacterStatics registers the java.lang.Character statics. The JVM
// overloads each of them on char and on int code point (the int form is how
// supplementary-plane characters are classified), so each takes either; the
// case mappings return the kind they were given, as the JVM's do.
func installCharacterStatics(ns *vm.Namespace) {
	predicate := func(name string, pred func(rune) bool) {
		ns.Def(name, mustWrap(func(vs []vm.Value) (vm.Value, error) {
			if len(vs) != 1 {
				return vm.NIL, fmt.Errorf("Character/%s expects 1 arg", name)
			}
			r, _, err := characterArg(name, vs[0])
			if err != nil {
				return vm.NIL, err
			}
			return vm.Boolean(pred(r)), nil
		}))
	}
	predicate("isDigit", unicode.IsDigit)
	predicate("isLetter", unicode.IsLetter)
	predicate("isLetterOrDigit", func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) })
	predicate("isLowerCase", unicode.IsLower)
	predicate("isUpperCase", unicode.IsUpper)
	predicate("isWhitespace", vm.JavaIsWhitespace)

	caseMap := func(name string, to func(rune) rune) {
		ns.Def(name, mustWrap(func(vs []vm.Value) (vm.Value, error) {
			if len(vs) != 1 {
				return vm.NIL, fmt.Errorf("Character/%s expects 1 arg", name)
			}
			r, isChar, err := characterArg(name, vs[0])
			if err != nil {
				return vm.NIL, err
			}
			if isChar {
				return vm.Char(to(r)), nil
			}
			return vm.MakeInt(int(to(r))), nil
		}))
	}
	caseMap("toLowerCase", unicode.ToLower)
	caseMap("toUpperCase", unicode.ToUpper)

	ns.Def("digit", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("Character/digit expects 2 args")
		}
		r, _, err := characterArg("digit", vs[0])
		if err != nil {
			return vm.NIL, err
		}
		radix, ok := vm.ToInt64(vs[1])
		if !ok {
			return vm.NIL, fmt.Errorf("Character/digit radix must be an integer, got %s", vs[1].Type().Name())
		}
		return vm.MakeInt(characterDigit(r, int(radix))), nil
	}))
}

// characterArg reads a char or an int code point, reporting which it was.
func characterArg(name string, v vm.Value) (rune, bool, error) {
	if c, ok := v.(vm.Char); ok {
		return rune(c), true, nil
	}
	if cp, ok := vm.ToInt64(v); ok {
		return rune(cp), false, nil
	}
	return 0, false, fmt.Errorf("Character/%s expected character or code point, got %s", name, v.Type().Name())
}

// characterDigit is java.lang.Character/digit: the value of r as a digit in
// radix, or -1. Decimal digits count in every script, and Latin letters
// (ASCII or fullwidth) stand for 10 through 35.
func characterDigit(r rune, radix int) int {
	if radix < 2 || radix > 36 {
		return -1
	}
	v := -1
	switch {
	case unicode.IsDigit(r):
		// Unicode keeps every decimal digit set in a contiguous run from
		// zero, so a digit's value is its distance from the run's start.
		start := r
		for start > 0 && unicode.IsDigit(start-1) {
			start--
		}
		v = int(r-start) % 10
	case 'a' <= r && r <= 'z':
		v = int(r-'a') + 10
	case 'A' <= r && r <= 'Z':
		v = int(r-'A') + 10
	case 'ａ' <= r && r <= 'ｚ':
		v = int(r-'ａ') + 10
	case 'Ａ' <= r && r <= 'Ｚ':
		v = int(r-'Ａ') + 10
	}
	if v >= radix {
		return -1
	}
	return v
}
