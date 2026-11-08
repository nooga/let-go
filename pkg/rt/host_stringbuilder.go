/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/vm"
)

// theHostStringBuilderType is the ValueType for the mutable
// java.lang.StringBuilder compat shim — a JVM-interop shim (mutable, not a
// persistent let-go value), so it lives in the compat layer rather than in
// pkg/vm. Motivated by honeysql, whose :clj fast paths build SQL strings via
// (StringBuilder.) / .append / .toString; the same idiom shows up in most
// performance-minded Clojure libraries.
type theHostStringBuilderType struct{}

func (t *theHostStringBuilderType) String() string     { return t.Name() }
func (t *theHostStringBuilderType) Type() vm.ValueType { return vm.TypeType }
func (t *theHostStringBuilderType) Unbox() any         { return nil }
func (t *theHostStringBuilderType) Name() string       { return "java.lang.StringBuilder" }
func (t *theHostStringBuilderType) Box(any) (vm.Value, error) {
	return vm.NIL, fmt.Errorf("java.lang.StringBuilder cannot be boxed")
}

// hostStringBuilderType is the singleton type for hostStringBuilder values.
var hostStringBuilderType = &theHostStringBuilderType{}

type hostStringBuilder struct{ b strings.Builder }

func (s *hostStringBuilder) Type() vm.ValueType { return hostStringBuilderType }
func (s *hostStringBuilder) Unbox() any         { return s }
func (s *hostStringBuilder) HostShim()          {}

// String is the builder's contents, as StringBuilder.toString is, so str,
// format's %s and print show what was built.
func (s *hostStringBuilder) String() string { return s.b.String() }

// appendString renders a value the way StringBuilder.append does: strings
// append their content, characters their rune, nil appends "null" (Java
// semantics — Clojure code that wants "" wraps the argument in str first),
// and everything else its str representation.
func (s *hostStringBuilder) appendString(v vm.Value) {
	switch a := v.(type) {
	case vm.String:
		s.b.WriteString(string(a))
	case vm.Char:
		s.b.WriteRune(rune(a))
	case *vm.Nil:
		s.b.WriteString("null")
	default:
		s.b.WriteString(strValue(v))
	}
}

func (s *hostStringBuilder) InvokeMethod(name vm.Symbol, args []vm.Value) (vm.Value, error) {
	switch string(name) {
	case "append":
		if len(args) == 1 {
			s.appendString(args[0])
			// Java's append returns the builder for chaining.
			return s, nil
		}
	case "toString":
		if len(args) == 0 {
			return vm.String(s.b.String()), nil
		}
	case "length":
		if len(args) == 0 {
			// Rune count, matching the string host methods' index space.
			return vm.Int(utf8.RuneCountInString(s.b.String())), nil
		}
	case "appendCodePoint":
		if len(args) == 1 {
			cp, ok := args[0].(vm.Int)
			if !ok {
				return vm.NIL, noMatchingMethod("java.lang.StringBuilder", "appendCodePoint", args)
			}
			if cp < 0 || cp > unicode.MaxRune {
				return vm.NIL, throwHost(vm.ClassIllegalArgument, "Not a valid Unicode code point: 0x%X", uint32(cp))
			}
			s.b.WriteRune(rune(cp))
			return s, nil
		}
	case "setLength":
		if len(args) == 1 {
			n, ok := args[0].(vm.Int)
			if !ok {
				return vm.NIL, noMatchingMethod("java.lang.StringBuilder", "setLength", args)
			}
			if n < 0 {
				return vm.NIL, throwHost(vm.ClassIndexOutOfBounds, "String index out of range: %d", int64(n))
			}
			// Truncate to n runes, or pad with NUL as Java does.
			runes := []rune(s.b.String())
			if int(n) <= len(runes) {
				runes = runes[:n]
			} else {
				runes = append(runes, make([]rune, int(n)-len(runes))...)
			}
			s.b.Reset()
			s.b.WriteString(string(runes))
			return vm.NIL, nil
		}
	case "charAt":
		if len(args) == 1 {
			i, ok := args[0].(vm.Int)
			if !ok {
				return vm.NIL, noMatchingMethod("java.lang.StringBuilder", "charAt", args)
			}
			runes := []rune(s.b.String())
			if i < 0 || int(i) >= len(runes) {
				return vm.NIL, throwHost(vm.ClassIndexOutOfBounds, "index %d,length %d", int64(i), len(runes))
			}
			return vm.Char(runes[i]), nil
		}
	case "reverse":
		if len(args) == 0 {
			runes := []rune(s.b.String())
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			s.b.Reset()
			s.b.WriteString(string(runes))
			return s, nil
		}
	}
	return vm.NIL, fmt.Errorf("java.lang.StringBuilder has no method .%s/%d", name, len(args))
}

// installHostStringBuilder registers the (StringBuilder.) constructor forms.
// The optional single argument follows java.lang.StringBuilder: an int is a
// capacity hint (ignored — strings.Builder grows on demand), a string is the
// initial content.
func installHostStringBuilder(ns *vm.Namespace) {
	ctor := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		sb := &hostStringBuilder{}
		switch len(vs) {
		case 0:
		case 1:
			switch a := vs[0].(type) {
			case vm.Int: // capacity hint
			case vm.String:
				sb.b.WriteString(string(a))
			default:
				return vm.NIL, fmt.Errorf("StringBuilder. expected an int capacity or initial string, got %s", a.Type().Name())
			}
		default:
			return vm.NIL, fmt.Errorf("StringBuilder. expected 0 or 1 argument, got %d", len(vs))
		}
		return sb, nil
	})
	for _, n := range []string{"StringBuilder.", "->StringBuilder", "java.lang.StringBuilder.", "->java.lang.StringBuilder"} {
		ns.Def(n, ctor)
	}
}
