/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

// Package jsonenc encodes the closed set of Go values the runtime hands to a
// JSON encoder: nil, bool, string, int, int64, float64, []any and
// map[string]any. Output is byte-identical to encoding/json.Marshal for those
// types, including its sorted map keys, HTML-safe escaping, float formatting
// and invalid-UTF-8 replacement.
//
// It exists so js/emit and the lg -w host page can encode without linking
// encoding/json, whose reflection-based encoder is live from an init()
// regardless of what a program calls. See the lg_no_json build tag.
package jsonenc

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"unicode/utf8"
)

// Marshal returns the JSON encoding of v.
func Marshal(v any) ([]byte, error) {
	return Append(nil, v)
}

// Append appends the JSON encoding of v to dst.
func Append(dst []byte, v any) ([]byte, error) {
	switch v := v.(type) {
	case nil:
		return append(dst, "null"...), nil
	case bool:
		return strconv.AppendBool(dst, v), nil
	case string:
		return AppendString(dst, v), nil
	case int:
		return strconv.AppendInt(dst, int64(v), 10), nil
	case int64:
		return strconv.AppendInt(dst, v, 10), nil
	case float64:
		return appendFloat(dst, v)
	case []any:
		if v == nil {
			return append(dst, "null"...), nil
		}
		dst = append(dst, '[')
		for i, e := range v {
			if i > 0 {
				dst = append(dst, ',')
			}
			var err error
			if dst, err = Append(dst, e); err != nil {
				return nil, err
			}
		}
		return append(dst, ']'), nil
	case map[string]any:
		if v == nil {
			return append(dst, "null"...), nil
		}
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		dst = append(dst, '{')
		for i, k := range keys {
			if i > 0 {
				dst = append(dst, ',')
			}
			dst = AppendString(dst, k)
			dst = append(dst, ':')
			var err error
			if dst, err = Append(dst, v[k]); err != nil {
				return nil, err
			}
		}
		return append(dst, '}'), nil
	default:
		return nil, fmt.Errorf("json: unsupported type: %T", v)
	}
}

// appendFloat follows encoding/json's floatEncoder: shortest round-trip
// digits, plain notation for 1e-6 <= |f| < 1e21, exponent notation outside
// it with a single-digit negative exponent ("1e-7", not "1e-07").
func appendFloat(dst []byte, f float64) ([]byte, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nil, fmt.Errorf("json: unsupported value: %s", strconv.FormatFloat(f, 'g', -1, 64))
	}
	format := byte('f')
	if abs := math.Abs(f); abs != 0 && (abs < 1e-6 || abs >= 1e21) {
		format = 'e'
	}
	dst = strconv.AppendFloat(dst, f, format, -1, 64)
	if format == 'e' {
		n := len(dst)
		if n >= 4 && dst[n-4] == 'e' && dst[n-3] == '-' && dst[n-2] == '0' {
			dst[n-2] = dst[n-1]
			dst = dst[:n-1]
		}
	}
	return dst, nil
}

const hex = "0123456789abcdef"

// AppendString appends s as a JSON string, escaped as encoding/json does with
// HTML escaping on (its Marshal default): <, > and & become \u003c, \u003e
// and \u0026, U+2028 and U+2029 are escaped, and each invalid UTF-8 byte
// becomes a raw U+FFFD, as encoding/json writes it. strconv.Quote is not a
// substitute: it emits \x, \a and \v, which JSON does not accept.
func AppendString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	start := 0
	for i := 0; i < len(s); {
		if c := s[i]; c < utf8.RuneSelf {
			if c >= 0x20 && c != '"' && c != '\\' && c != '<' && c != '>' && c != '&' {
				i++
				continue
			}
			dst = append(dst, s[start:i]...)
			switch c {
			case '"', '\\':
				dst = append(dst, '\\', c)
			case '\b':
				dst = append(dst, '\\', 'b')
			case '\f':
				dst = append(dst, '\\', 'f')
			case '\n':
				dst = append(dst, '\\', 'n')
			case '\r':
				dst = append(dst, '\\', 'r')
			case '\t':
				dst = append(dst, '\\', 't')
			default:
				dst = append(dst, '\\', 'u', '0', '0', hex[c>>4], hex[c&0xf])
			}
			i++
			start = i
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			dst = append(dst, s[start:i]...)
			dst = append(dst, "\ufffd"...)
			i += size
			start = i
			continue
		}
		if r == '\u2028' || r == '\u2029' {
			dst = append(dst, s[start:i]...)
			dst = append(dst, '\\', 'u', '2', '0', '2', hex[r&0xf])
			i += size
			start = i
			continue
		}
		i += size
	}
	dst = append(dst, s[start:]...)
	return append(dst, '"')
}
