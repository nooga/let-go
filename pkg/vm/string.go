/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"regexp"
	"regexp/syntax"
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
			r, ok := RuneAt(string(l), int(i))
			if !ok {
				return NIL, fmt.Errorf("string.charAt: index %d out of bounds for length %d", int(i), utf8.RuneCountInString(string(l)))
			}
			return Char(r), nil
		}
	case "indexOf":
		if len(args) == 1 || len(args) == 2 {
			str := string(l)
			from, bfrom := 0, 0
			if len(args) == 2 {
				f, ok := args[1].(Int)
				if !ok {
					return NIL, fmt.Errorf("string.indexOf fromIndex must be an int, got %s", args[1].Type().Name())
				}
				// Java clamps rather than throwing on out-of-range fromIndex.
				from = max(int(f), 0)
				if bfrom, ok = RuneOffset(str, from); !ok {
					from, bfrom = utf8.RuneCountInString(str), len(str)
				}
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
			if !ByteSearchable(needle) {
				// Search the []rune round-trip, as before: it turns each
				// invalid byte of the string into U+FFFD, which an invalid
				// or U+FFFD needle would match differently on the raw bytes.
				tail := string([]rune(str)[from:])
				idx := strings.Index(tail, needle)
				if idx < 0 {
					return Int(-1), nil
				}
				return Int(from + utf8.RuneCountInString(tail[:idx])), nil
			}
			tail := str[bfrom:]
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
			str := string(l)
			begin, ok := args[0].(Int)
			if !ok {
				return NIL, fmt.Errorf("string.substring beginIndex must be an int, got %s", args[0].Type().Name())
			}
			bbegin, okBegin := RuneOffset(str, int(begin))
			bend, okEnd := len(str), true
			var end Int
			if len(args) == 2 {
				if end, ok = args[1].(Int); !ok {
					return NIL, fmt.Errorf("string.substring endIndex must be an int, got %s", args[1].Type().Name())
				}
				bend, okEnd = RuneOffset(str, int(end))
			}
			if !okBegin || !okEnd || bbegin > bend {
				n := utf8.RuneCountInString(str)
				if len(args) == 1 {
					end = Int(n)
				}
				return NIL, fmt.Errorf("string.substring: range [%d, %d) out of bounds for length %d", int(begin), int(end), n)
			}
			sub := str[bbegin:bend]
			if !utf8.ValidString(sub) {
				// Slicing []rune turned invalid bytes into U+FFFD; keep that.
				return String([]rune(sub)), nil
			}
			return String(sub), nil
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
	case "strip":
		if len(args) == 0 {
			return String(strings.TrimFunc(string(l), JavaIsWhitespace)), nil
		}
	case "stripLeading":
		if len(args) == 0 {
			return String(strings.TrimLeftFunc(string(l), JavaIsWhitespace)), nil
		}
	case "stripTrailing":
		if len(args) == 0 {
			return String(strings.TrimRightFunc(string(l), JavaIsWhitespace)), nil
		}
	case "isBlank":
		if len(args) == 0 {
			return Boolean(strings.IndexFunc(string(l), func(r rune) bool { return !JavaIsWhitespace(r) }) < 0), nil
		}
	case "lastIndexOf":
		if len(args) == 1 || len(args) == 2 {
			return l.lastIndexOf(args)
		}
	case "codePoints":
		if len(args) == 0 {
			cps := make(CodePointStream, 0, len(l))
			for _, r := range string(l) {
				cps = append(cps, int64(r))
			}
			return cps, nil
		}
	case "codePointCount":
		if len(args) == 2 {
			begin, ok1 := args[0].(Int)
			end, ok2 := args[1].(Int)
			if !ok1 || !ok2 {
				return NIL, fmt.Errorf("string.codePointCount indices must be ints")
			}
			// Indices are rune indices already, so the count is the span.
			n := utf8.RuneCountInString(string(l))
			if begin < 0 || int(end) > n || begin > end {
				return NIL, fmt.Errorf("string.codePointCount: range [%d, %d) out of bounds for length %d", int(begin), int(end), n)
			}
			return end - begin, nil
		}
	case "split":
		if len(args) == 1 || len(args) == 2 {
			return l.split(args)
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

// lastIndexOf is java.lang.String.lastIndexOf over rune indices: the needle is
// a string, a character or a code point, and the optional fromIndex is the
// last index a match may start at.
func (l String) lastIndexOf(args []Value) (Value, error) {
	var needle []rune
	switch a := args[0].(type) {
	case String:
		needle = []rune(string(a))
	case Char:
		needle = []rune{rune(a)}
	case Int:
		needle = []rune{rune(a)}
	default:
		return NIL, fmt.Errorf("string.lastIndexOf expected a string or character, got %s", args[0].Type().Name())
	}
	hay := []rune(string(l))
	from := len(hay) - len(needle)
	if len(args) == 2 {
		f, ok := args[1].(Int)
		if !ok {
			return NIL, fmt.Errorf("string.lastIndexOf fromIndex must be an int, got %s", args[1].Type().Name())
		}
		from = min(from, int(f))
	}
	for i := from; i >= 0; i-- {
		if string(hay[i:i+len(needle)]) == string(needle) {
			return Int(i), nil
		}
	}
	return Int(-1), nil
}

// split is java.lang.String.split: the separator is a regex, a positive
// limit caps the number of parts, zero drops trailing empty strings, and a
// negative limit keeps them. A zero-width match at index 0 yields no leading
// empty string. It returns an array, as the JVM's String[] is one.
func (l String) split(args []Value) (Value, error) {
	pat, ok := args[0].(String)
	if !ok {
		return NIL, fmt.Errorf("string.split expected a regex string, got %s", args[0].Type().Name())
	}
	limit := 0
	if len(args) == 2 {
		n, ok := args[1].(Int)
		if !ok {
			return NIL, fmt.Errorf("string.split limit must be an int, got %s", args[1].Type().Name())
		}
		limit = int(n)
	}
	re, err := regexp.Compile(string(pat))
	if err != nil {
		return NIL, fmt.Errorf("string.split: %w", err)
	}
	str := string(l)
	if str == "" {
		return NewObjectArrayFrom([]Value{String("")}), nil
	}
	var parts []Value
	start := 0
	ms, err := javaMatches(re, string(pat), str)
	if err != nil {
		return NIL, fmt.Errorf("string.split: %w", err)
	}
	for _, m := range ms {
		if limit > 0 && len(parts) == limit-1 {
			break
		}
		if m[1] == 0 {
			continue // a zero-width match at the start splits nothing off
		}
		parts = append(parts, String(str[start:m[0]]))
		start = m[1]
	}
	parts = append(parts, String(str[start:]))
	if limit == 0 {
		for len(parts) > 0 && parts[len(parts)-1] == String("") {
			parts = parts[:len(parts)-1]
		}
	}
	return NewObjectArrayFrom(parts), nil
}

// javaMatches returns the matches of re in str as Matcher.find enumerates
// them. Go's FindAllStringIndex drops an empty match that touches the end of
// the previous match; Java reports it, so "abc" split on "a*" sees [0,1) and
// then [1,1). The dropped empties are restored here.
//
// Whether the pattern matches empty at such a position e is tested by
// anchoring \A(?s:.)(pat) at the rune before e, so group 1 is the match
// starting at e and its assertions see the same left context: RE2's \b, \B
// and (?m)^ look back one rune, and \A and ^ cannot match there since e > 0.
// pat is spliced in as its parsed form, which closes any \Q quote that the
// source leaves open to the end of the pattern.
func javaMatches(re *regexp.Regexp, pat, str string) ([][]int, error) {
	ms := re.FindAllStringIndex(str, -1)
	var at *regexp.Regexp
	emptyAt := func(e int) (bool, error) {
		if at == nil {
			parsed, err := syntax.Parse(pat, syntax.Perl)
			if err != nil {
				return false, err
			}
			if at, err = regexp.Compile(`\A(?s:.)(` + parsed.String() + `)`); err != nil {
				return false, err
			}
		}
		_, size := utf8.DecodeLastRuneInString(str[:e])
		loc := at.FindStringSubmatchIndex(str[e-size:])
		return loc != nil && loc[2] == size && loc[3] == size, nil
	}
	out := make([][]int, 0, len(ms))
	prevEnd := -1 // end of the previous non-empty match
	for _, m := range ms {
		if prevEnd >= 0 && m[0] != prevEnd {
			ok, err := emptyAt(prevEnd)
			if err != nil {
				return nil, err
			}
			if ok {
				out = append(out, []int{prevEnd, prevEnd})
			}
		}
		out = append(out, m)
		prevEnd = -1
		if m[1] > m[0] {
			prevEnd = m[1]
		}
	}
	if prevEnd >= 0 {
		ok, err := emptyAt(prevEnd)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, []int{prevEnd, prevEnd})
		}
	}
	return out, nil
}
