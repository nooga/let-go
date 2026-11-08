/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"unicode"
)

// JavaIsWhitespace reports whether r is whitespace by
// java.lang.Character/isWhitespace: a Unicode space, line or paragraph
// separator other than the non-breaking ones, or one of the ASCII controls
// \t \n \u000B \f \r \u001C-\u001F. String .strip and .isBlank use the same
// rule, so it differs from Go's unicode.IsSpace (which counts U+0085 and the
// non-breaking spaces).
func JavaIsWhitespace(r rune) bool {
	switch r {
	case ' ', ' ', ' ':
		return false
	case '\t', '\n', '\u000B', '\f', '\r', '\u001C', '\u001D', '\u001E', '\u001F':
		return true
	}
	return unicode.In(r, unicode.Zs, unicode.Zl, unicode.Zp)
}

// theCodePointStreamType is the type of what String .codePoints returns: the
// java.util.stream.IntStream surface Clojure code reaches for, which is
// .toArray (usually fed to vec) and .count.
type theCodePointStreamType struct{}

func (t *theCodePointStreamType) String() string  { return t.Name() }
func (t *theCodePointStreamType) Type() ValueType { return TypeType }
func (t *theCodePointStreamType) Unbox() any      { return nil }
func (t *theCodePointStreamType) Name() string    { return "java.util.stream.IntStream" }
func (t *theCodePointStreamType) Box(bare any) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

// CodePointStreamType is the singleton type for CodePointStream values.
var CodePointStreamType = &theCodePointStreamType{}

// CodePointStream holds the code points of a string.
type CodePointStream []int64

func (c CodePointStream) Type() ValueType { return CodePointStreamType }
func (c CodePointStream) Unbox() any      { return []int64(c) }
func (c CodePointStream) String() string  { return "#<java.util.stream.IntStream>" }
func (c CodePointStream) HostShim()       {}

func (c CodePointStream) InvokeMethod(name Symbol, args []Value) (Value, error) {
	switch name {
	case "toArray":
		if len(args) == 0 {
			return NewIntArrayFrom(append([]int64(nil), c...)), nil
		}
	case "count":
		if len(args) == 0 {
			return MakeInt(len(c)), nil
		}
	default:
		return NIL, fmt.Errorf("java.util.stream.IntStream has no method .%s", name)
	}
	return NIL, fmt.Errorf("IntStream.%s: wrong number of arguments %d", name, len(args))
}
