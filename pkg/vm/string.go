/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"
)

type theStringType struct {
	zero String
}

func (t *theStringType) String() string  { return t.Name() }
func (t *theStringType) Type() ValueType { return TypeType }
func (t *theStringType) Unbox() any      { return reflect.TypeFor[*theStringType]() }

func (t *theStringType) Name() string { return "let-go.lang.String" }

func (t *theStringType) Box(bare any) (Value, error) {
	raw, ok := bare.(string)
	if !ok {
		return StringType.zero, NewTypeError(bare, "can't be boxed as", t)
	}
	return String(raw), nil
}

// StringType is the type of StringValues
var StringType *theStringType = &theStringType{zero: ""}

// String is boxed int
type String string

func (l String) Conj(value Value) Collection {
	return String(string(l) + value.String())
}

func (l String) RawCount() int {
	return len(l)
}

func (l String) Count() Value {
	return Int(len(l))
}

func (l String) Empty() Collection {
	return String("")
}

// Hash implements Hashable for fast map lookups.
func (l String) Hash() uint32 { return hashString(string(l)) }

// Type implements Value
func (l String) Type() ValueType { return StringType }

// Unbox implements Unbox
func (l String) Unbox() any {
	return string(l)
}

func (l String) InvokeMethod(name Symbol, args []Value) (Value, error) {
	switch name {
	case "replace":
		if len(args) != 2 {
			return NIL, fmt.Errorf("string.replace expected 2 arguments")
		}
		old, ok := args[0].(String)
		if !ok {
			return NIL, fmt.Errorf("string.replace expected string target")
		}
		repl, ok := args[1].(String)
		if !ok {
			return NIL, fmt.Errorf("string.replace expected string replacement")
		}
		return String(strings.ReplaceAll(string(l), string(old), string(repl))), nil
	case "getBytes":
		if len(args) == 0 {
			// No-arg form: let-go String is Go's UTF-8 string, so the
			// default encoding is UTF-8 by construction.
			return NewByteArrayFrom([]byte(string(l))), nil
		}
		if len(args) == 1 {
			enc, ok := args[0].(String)
			if !ok {
				return NIL, fmt.Errorf("string.getBytes encoding must be string, got %s", args[0].Type().Name())
			}
			e := strings.ToUpper(string(enc))
			if e != "UTF-8" && e != "UTF8" {
				return NIL, fmt.Errorf("string.getBytes: only UTF-8 supported, got %q", string(enc))
			}
			return NewByteArrayFrom([]byte(string(l))), nil
		}
		return NIL, fmt.Errorf("string.getBytes expected 0 or 1 argument, got %d", len(args))
	// The java.lang.String surface Clojure libraries reach on their :clj
	// branches (honeysql's entity formatting walks .length/.charAt/.indexOf,
	// its util.str calls .toString/.concat). Indices are RUNE indices — Java's
	// are UTF-16 units, and runes match them for all BMP text while byte
	// indices break on the first non-ASCII character. Case mapping is Go's
	// locale-independent one, which is exactly what code reaching for
	// toUpperCase(Locale/US) wants; a locale argument is accepted and ignored.
	case "toString":
		if len(args) == 0 {
			return l, nil
		}
	case "length":
		if len(args) == 0 {
			return Int(utf8.RuneCountInString(string(l))), nil
		}
	case "isEmpty":
		if len(args) == 0 {
			return Boolean(len(l) == 0), nil
		}
	case "charAt":
		if len(args) == 1 {
			i, ok := args[0].(Int)
			if !ok {
				return NIL, fmt.Errorf("string.charAt index must be an int, got %s", args[0].Type().Name())
			}
			rs := []rune(string(l))
			if i < 0 || int(i) >= len(rs) {
				return NIL, fmt.Errorf("string.charAt: index %d out of bounds for length %d", int(i), len(rs))
			}
			return Char(rs[i]), nil
		}
	case "indexOf":
		if len(args) == 1 || len(args) == 2 {
			rs := []rune(string(l))
			from := 0
			if len(args) == 2 {
				f, ok := args[1].(Int)
				if !ok {
					return NIL, fmt.Errorf("string.indexOf fromIndex must be an int, got %s", args[1].Type().Name())
				}
				// Java clamps rather than throwing on out-of-range fromIndex.
				from = min(max(int(f), 0), len(rs))
			}
			var needle string
			switch a := args[0].(type) {
			case String:
				needle = string(a)
			case Char:
				needle = string(rune(a))
			case Int: // Java's indexOf(int ch) codepoint form
				needle = string(rune(a))
			default:
				return NIL, fmt.Errorf("string.indexOf expected a string or character, got %s", args[0].Type().Name())
			}
			tail := string(rs[from:])
			idx := strings.Index(tail, needle)
			if idx < 0 {
				return Int(-1), nil
			}
			return Int(from + utf8.RuneCountInString(tail[:idx])), nil
		}
	case "concat":
		if len(args) == 1 {
			s, ok := args[0].(String)
			if !ok {
				return NIL, fmt.Errorf("string.concat expected a string, got %s", args[0].Type().Name())
			}
			return l + s, nil
		}
	case "substring":
		if len(args) == 1 || len(args) == 2 {
			rs := []rune(string(l))
			begin, ok := args[0].(Int)
			if !ok {
				return NIL, fmt.Errorf("string.substring beginIndex must be an int, got %s", args[0].Type().Name())
			}
			end := Int(len(rs))
			if len(args) == 2 {
				if end, ok = args[1].(Int); !ok {
					return NIL, fmt.Errorf("string.substring endIndex must be an int, got %s", args[1].Type().Name())
				}
			}
			if begin < 0 || begin > end || int(end) > len(rs) {
				return NIL, fmt.Errorf("string.substring: range [%d, %d) out of bounds for length %d", int(begin), int(end), len(rs))
			}
			return String(rs[begin:end]), nil
		}
	case "startsWith":
		if len(args) == 1 {
			if s, ok := args[0].(String); ok {
				return Boolean(strings.HasPrefix(string(l), string(s))), nil
			}
			return NIL, fmt.Errorf("string.startsWith expected a string, got %s", args[0].Type().Name())
		}
	case "endsWith":
		if len(args) == 1 {
			if s, ok := args[0].(String); ok {
				return Boolean(strings.HasSuffix(string(l), string(s))), nil
			}
			return NIL, fmt.Errorf("string.endsWith expected a string, got %s", args[0].Type().Name())
		}
	case "contains":
		if len(args) == 1 {
			if s, ok := args[0].(String); ok {
				return Boolean(strings.Contains(string(l), string(s))), nil
			}
			return NIL, fmt.Errorf("string.contains expected a string, got %s", args[0].Type().Name())
		}
	case "toUpperCase":
		if len(args) <= 1 { // optional locale argument, ignored
			return String(strings.ToUpper(string(l))), nil
		}
	case "toLowerCase":
		if len(args) <= 1 { // optional locale argument, ignored
			return String(strings.ToLower(string(l))), nil
		}
	case "trim":
		if len(args) == 0 {
			return String(strings.TrimSpace(string(l))), nil
		}
	default:
		return NIL, fmt.Errorf("method %s not found on string", name)
	}
	return NIL, fmt.Errorf("string.%s: wrong number of arguments %d", name, len(args))
}

// First implements Seq
func (l String) First() Value {
	for _, r := range l {
		return Char(r)
	}
	return NIL
}

// More implements Seq
func (l String) More() Seq {
	r := l.Next()
	if r == nil {
		return EmptyList
	}
	return r
}

// Next implements Seq
func (l String) Next() Seq {
	if len(l) <= 1 {
		return nil
	}
	ret := EmptyList
	s := []rune(l)
	for i := len(s) - 1; i >= 1; i-- {
		ret = ret.Conj(Char(s[i])).(*List)
	}
	return ret
}

// Cons implements Seq
func (l String) Cons(val Value) Seq {
	return NIL
}

func (l String) Seq() Seq {
	if len(l) == 0 {
		return nil
	}
	ret := EmptyList
	s := []rune(l)
	for i := len(s) - 1; i >= 0; i-- {
		ret = ret.Conj(Char(s[i])).(*List)
	}
	return ret
}

// Nth implements Indexed: positional access by integer index.
func (l String) Nth(i int) Value { return l.ValueAt(Int(i)) }

func (l String) ValueAt(key Value) Value {
	return l.ValueAtOr(key, NIL)
}

func (l String) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	numkey, ok := key.(Int)
	if !ok || numkey < 0 || int(numkey) >= len(l) {
		// len(l) is a byte count ≥ rune count, so negative/oversized keys
		// bail here without walking; in-range-by-bytes keys verify below.
		return dflt
	}
	// Walk runes in place instead of materializing []rune(l) — a single
	// (nth s i) used to allocate the WHOLE rune slice (O(n) bytes per
	// char lookup; 1.4GB in one lgbgen --target=go run).
	i := int(numkey)
	n := 0
	for _, r := range string(l) {
		if n == i {
			return Char(r)
		}
		n++
	}
	return dflt
}

// String returns the EDN/Clojure-readable form: surrounded by double quotes
// with EDN-conformant escapes. NOTE: Go's %q emits \xNN / \a / \v for control
// chars, which the reader rejects ("unknown escape sequence \x") and which are
// not valid EDN — so %q did not round-trip its own output (e.g. ANSI-colored
// strings like "\x1b[32m..."). Escapes here match the reader's accepted set
// (reader.go): \" \\ \t \r \n \b \f, and \uXXXX for everything else below
// 0x20 plus DEL; printable Unicode passes through as raw UTF-8.
func (l String) String() string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range string(l) {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\t':
			b.WriteString(`\t`)
		case '\r':
			b.WriteString(`\r`)
		case '\n':
			b.WriteString(`\n`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		default:
			if r < 0x20 || r == 0x7f {
				fmt.Fprintf(&b, `\u%04X`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
	return b.String()
}
